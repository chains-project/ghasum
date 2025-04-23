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

package gha

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"path"
	"path/filepath"
	"slices"
)

type workflowFile struct {
	path    string
	content []byte
}

func actionsInManifest(manifest manifest) ([]GitHubAction, error) {
	unique := make(map[string]GitHubAction, 0)
	err := actionsInSteps(manifest.Runs.Steps, unique)
	if err != nil {
		return nil, err
	}

	return slices.Collect(maps.Values(unique)), nil
}

func actionsInWorkflows(workflows []workflow) ([]GitHubAction, error) {
	unique := make(map[string]GitHubAction, 0)
	for _, workflow := range workflows {
		for _, job := range workflow.Jobs {
			err := actionsInSteps(job.Steps, unique)
			if err != nil {
				return nil, err
			}
		}
	}

	return slices.Collect(maps.Values(unique)), nil
}

func actionsInSteps(steps []step, m map[string]GitHubAction) error {
	for _, step := range steps {
		uses := step.Uses
		if uses == "" {
			continue
		}

		action, err := parseUses(uses)
		if errors.Is(err, ErrLocalAction) || errors.Is(err, ErrDockerUses) {
			continue
		} else if err != nil {
			return err
		}

		id := fmt.Sprintf("%s%s%s%s", action.Owner, action.Project, action.Path, action.Ref)
		m[id] = action
	}

	return nil
}

func workflowsInRepo(repo fs.FS) ([]workflowFile, error) {
	workflows := make([]workflowFile, 0)
	walk := func(entryPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if entryPath == WorkflowsPath {
				return nil
			} else {
				return fs.SkipDir
			}
		}

		if ext := path.Ext(entryPath); ext != ".yml" && ext != ".yaml" {
			return nil
		}

		data, err := workflowInRepo(repo, entryPath)
		if err != nil {
			return err
		}

		workflows = append(workflows, workflowFile{
			content: data,
			path:    entryPath,
		})

		return nil
	}

	if err := fs.WalkDir(repo, WorkflowsPath, walk); err != nil {
		return nil, fmt.Errorf("failed to find workflows: %v", err)
	}

	return workflows, nil
}

func workflowInRepo(repo fs.FS, path string) ([]byte, error) {
	file, err := repo.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open workflow at %q: %v", path, err)
	}

	data, _ := io.ReadAll(file)
	return data, nil
}

func manifestInRepo(repo fs.FS, dir string) ([]byte, error) {
	path := filepath.Join(dir, "action.yml")
	file, err := repo.Open(path)
	if err != nil {
		path = filepath.Join(dir, "action.yaml")
		if file, err = repo.Open(path); err != nil {
			return nil, fmt.Errorf("could not open manifest (action.yml or action.yaml) at %s: %v", dir, err)
		}
	}

	data, _ := io.ReadAll(file)
	return data, nil
}
