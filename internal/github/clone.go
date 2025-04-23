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

package github

import (
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// A Repository represents a GitHub repository.
type Repository struct {
	// Owner is the name of the user or organization that owns the project.
	Owner string

	// Project is the name of the project to clone.
	Project string

	// Ref is the reference to check out.
	Ref string
}

// Clone will clone the given repository at the exact ref from GitHub into the
// given directory. Note that the git index will be omitted.
func Clone(dir string, repo *Repository) error {
	if err := clone(dir, repo); err != nil {
		return err
	}

	if err := os.RemoveAll(path.Join(dir, ".git")); err != nil {
		return fmt.Errorf("could not remove git index: %v", err)
	}

	return nil
}

func clone(dir string, repo *Repository) error {
	if err := cloneAtTag(dir, repo); err == nil {
		return nil
	}

	if err := cloneAtBranch(dir, repo); err == nil {
		return nil
	}

	return cloneAtCommit(dir, repo)
}

func cloneAtBranch(dir string, repo *Repository) error {
	opts := git.CloneOptions{
		URL:           toUrl(repo),
		Depth:         1,
		SingleBranch:  true,
		Tags:          git.NoTags,
		ReferenceName: plumbing.NewBranchReferenceName(repo.Ref),
	}

	_, err := git.PlainClone(dir, false, &opts)
	if err != nil {
		return fmt.Errorf("could not clone %q (as branch) from %q: %v", repo.Ref, opts.URL, err)
	}

	return nil
}

func cloneAtCommit(dir string, repo *Repository) error {
	repository, err := git.PlainInit(dir, false)
	if err != nil {
		return fmt.Errorf("could not initialize a repository for %s/%s: %v", repo.Owner, repo.Project, err)
	}

	remote := "origin"
	url := toUrl(repo)

	remoteCfg := config.RemoteConfig{
		Name: remote,
		URLs: []string{url},
	}
	if _, err = repository.CreateRemote(&remoteCfg); err != nil {
		return fmt.Errorf("could not set remote of the repository to %q: %v", url, err)
	}

	remoteRef := repo.Ref
	localRef := "refs/heads/tmp"

	fetchOpts := git.FetchOptions{
		Depth:      1,
		RemoteName: remote,
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("%s:%s", remoteRef, localRef)),
		},
	}
	if err = repository.Fetch(&fetchOpts); err != nil {
		return fmt.Errorf("could not fetch commits from %q: %v", url, err)
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return fmt.Errorf("could not obtain worktree for %s/%s: %v", repo.Owner, repo.Project, err)
	}

	checkoutOpts := git.CheckoutOptions{
		Hash: plumbing.NewHash(repo.Ref),
	}
	if err = worktree.Checkout(&checkoutOpts); err != nil {
		return fmt.Errorf("could not checkout ref %q for %s/%s: %v", repo.Ref, repo.Owner, repo.Project, err)
	}

	return nil
}

func cloneAtTag(dir string, repo *Repository) error {
	opts := git.CloneOptions{
		URL:           toUrl(repo),
		Depth:         1,
		SingleBranch:  true,
		Tags:          git.NoTags,
		ReferenceName: plumbing.NewTagReferenceName(repo.Ref),
	}

	_, err := git.PlainClone(dir, false, &opts)
	if err != nil {
		return fmt.Errorf("could not clone %q (as tag) from %q: %v", repo.Ref, opts.URL, err)
	}

	return nil
}

func toUrl(repo *Repository) (url string) {
	return fmt.Sprintf("https://github.com/%s/%s", repo.Owner, repo.Project)
}
