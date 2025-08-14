// Copyright 2024-2025 Eric Cornelissen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ghasum

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/chains-project/ghasum/internal/checksum"
	"github.com/chains-project/ghasum/internal/gha"
	"github.com/chains-project/ghasum/internal/github"
	"github.com/chains-project/ghasum/internal/sumfile"
)

var ghasumPath = path.Join(gha.WorkflowsPath, "gha.sum")

func clear(file *os.File) error {
	if _, err := file.Seek(0, 0); err != nil {
		return errors.Join(ErrSumfileWrite, err)
	}

	if err := file.Truncate(0); err != nil {
		return errors.Join(ErrSumfileWrite, err)
	}

	return nil
}

func clone(cfg *Config, action *gha.GitHubAction) (string, error) {
	actionDir := path.Join(cfg.Cache.Path(), action.Owner, action.Project, action.Ref)
	if _, err := os.Stat(actionDir); err != nil {
		if cfg.Offline {
			return actionDir, fmt.Errorf("missing %q from cache", action)
		}

		repo := github.Repository{
			Owner:   action.Owner,
			Project: action.Project,
			Ref:     action.Ref,
		}

		err := github.Clone(actionDir, &repo)
		if err != nil {
			return actionDir, fmt.Errorf("clone failed: %v", err)
		}
	}

	return actionDir, nil
}

func compare(got, want []sumfile.Entry, reportRedundant bool) []Problem {
	toMap := func(entries []sumfile.Entry) map[string]string {
		m := make(map[string]string, len(entries))
		for _, entry := range entries {
			key := fmt.Sprintf("%s@%s", entry.ID[0], entry.ID[1])
			m[key] = entry.Checksum
		}

		return m
	}

	cmp := func(got, want map[string]string) []Problem {
		problems := make([]Problem, 0)
		for key, got := range got {
			want, ok := want[key]
			if !ok {
				p := fmt.Sprintf("no checksum found for %q", key)
				problems = append(problems, Problem(p))
				continue
			}

			if got != want {
				p := fmt.Sprintf("checksum mismatch for %q", key)
				problems = append(problems, Problem(p))
			}
		}

		if reportRedundant {
			for key := range want {
				if _, ok := got[key]; !ok {
					p := fmt.Sprintf("redundant checksum for %q", key)
					problems = append(problems, Problem(p))
				}
			}
		}

		return problems
	}

	return cmp(toMap(got), toMap(want))
}

func find(cfg *Config) (tree, error) {
	var (
		actions []gha.GitHubAction
		err     error
	)

	switch {
	case cfg.Job != "":
		actions, err = gha.JobActions(cfg.Repo, cfg.Workflow, cfg.Job)
	case cfg.Workflow != "":
		actions, err = gha.WorkflowActions(cfg.Repo, cfg.Workflow)
	default:
		actions, err = gha.RepoActions(cfg.Repo)
	}

	root := tree{}
	if err != nil {
		return root, fmt.Errorf("could not find GitHub Actions: %v", err)
	}

	parents := make([]*tree, len(actions))
	for i := range actions {
		parents[i] = &root
	}

	for i := 0; i < len(actions); i++ {
		action := actions[i]
		actionDir, err := clone(cfg, &action)
		if err != nil {
			return root, err
		}

		subtree := tree{value: &action}
		if cfg.Transitive {
			repo, _ := os.OpenRoot(actionDir)

			var transitive []gha.GitHubAction
			switch action.Kind {
			case gha.Action:
				transitive, err = gha.ManifestActions(repo.FS(), action.Path)
				if err != nil {
					return root, fmt.Errorf("action manifest parsing failed for %s: %v", action, err)
				}
			case gha.ReusableWorkflow:
				transitive, err = gha.WorkflowActions(repo.FS(), action.Path)
				if err != nil {
					return root, fmt.Errorf("reusable workflow parsing failed for %s: %v", action, err)
				}
			}

			for _, action := range transitive {
				actions = append(actions, action)
				parents = append(parents, &subtree)
			}
		}

		parents[i].add(&subtree)
	}

	return root, nil
}

func compute(cfg *Config, actions tree, algo checksum.Algo) ([]sumfile.Entry, error) {
	if err := cfg.Cache.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize cache: %v", err)
	} else {
		defer cfg.Cache.Cleanup()
	}

	list := slices.Collect(actions.All())
	entries := make(map[string]sumfile.Entry, len(list))
	for _, action := range list {
		actionDir, err := clone(cfg, &action)
		if err != nil {
			return nil, err
		}

		id := fmt.Sprintf("%s%s%s", action.Owner, action.Project, action.Ref)
		if _, ok := entries[id]; !ok {
			checksum, err := checksum.Compute(actionDir, algo)
			if err != nil {
				return nil, fmt.Errorf("could not compute checksum for %q: %v", action, err)
			}

			entries[id] = sumfile.Entry{
				ID:       []string{fmt.Sprintf("%s/%s", action.Owner, action.Project), action.Ref},
				Checksum: strings.Replace(checksum, "h1:", "", 1),
			}
		}
	}

	return slices.Collect(maps.Values(entries)), nil
}

func create(base string) (*os.File, error) {
	fullGhasumPath := path.Join(base, ghasumPath)

	if _, err := os.Stat(fullGhasumPath); err == nil {
		return nil, ErrInitialized
	}

	file, err := os.OpenFile(fullGhasumPath, os.O_CREATE|os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return nil, errors.Join(ErrSumfileCreate, err)
	}

	return file, nil
}

func decode(stored []byte) ([]sumfile.Entry, error) {
	checksums, err := sumfile.Decode(string(stored))
	if err != nil {
		return nil, errors.Join(ErrSumfileDecode, err)
	}

	return checksums, nil
}

func encode(version sumfile.Version, checksums []sumfile.Entry) (string, error) {
	content, err := sumfile.Encode(version, checksums)
	if err != nil {
		return "", errors.Join(ErrSumfileEncode, err)
	}

	return content, nil
}

func list(cfg *Config, t *tree) string {
	var b strings.Builder

	action := t.value
	root := action == nil

	if !root {
		b.WriteString(action.String())
		b.WriteString(" (")
		b.WriteString(action.Kind.String())
		if !cfg.Offline {
			isArchived, err := github.Archived(&github.Repository{
				Owner:   action.Owner,
				Project: action.Project,
			})
			if err == nil && isArchived {
				b.WriteString(", archived")
			}
		}
		b.WriteString(")\n")
	}

	ordered := slices.SortedFunc(
		slices.Values(t.children),
		func(a, b *tree) int {
			return strings.Compare(a.value.String(), b.value.String())
		},
	)

	for _, children := range ordered {
		for line := range strings.Lines(list(cfg, children)) {
			if !root {
				b.WriteString("  ")
			}
			b.WriteString(line)
		}
	}

	return b.String()
}

func open(base string) (*os.File, error) {
	fullGhasumPath := path.Join(base, ghasumPath)

	file, err := os.OpenFile(fullGhasumPath, os.O_RDWR, os.ModeExclusive)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrNotInitialized
	} else if err != nil {
		return nil, errors.Join(ErrSumfileOpen, err)
	}

	if err := os.Chmod(fullGhasumPath, fs.ModeExclusive); err != nil {
		return file, errors.Join(ErrSumfileUnlock, err)
	}

	return file, nil
}

func remove(base string) error {
	fullGhasumPath := path.Join(base, ghasumPath)
	if err := os.Remove(fullGhasumPath); err != nil {
		return errors.Join(ErrSumfileRemove, err)
	}

	return nil
}

func unlock(base string) error {
	fullGhasumPath := path.Join(base, ghasumPath)
	if err := os.Chmod(fullGhasumPath, fs.ModePerm); err != nil {
		return errors.Join(ErrSumfileUnlock, err)
	}

	return nil
}

func version(stored []byte) (sumfile.Version, error) {
	version, err := sumfile.DecodeVersion(string(stored))
	if err != nil {
		return version, errors.Join(ErrSumfileDecode, err)
	}

	return version, nil
}

func write(file *os.File, content string) error {
	if _, err := file.WriteString(content); err != nil {
		return errors.Join(ErrSumfileWrite, err)
	}

	return nil
}
