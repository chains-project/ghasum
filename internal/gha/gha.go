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
	"io/fs"
	"path/filepath"
)

// A GitHubAction identifies a specific version of a GitHub Action.
type GitHubAction struct {
	// Owner is the GitHub user or organization that owns the repository that
	// houses the GitHub Action.
	Owner string

	// Project is the name of the GitHub repository (excluding the owner) that
	// houses the GitHub Action.
	Project string

	// Path is the path of the action inside the GitHub repository.
	Path string

	// Ref is the git ref (branch, tag, commit SHA), also known as version, of the
	// GitHub Action.
	Ref string

	// Kind is the [ActionKind] of the GitHub Action.
	Kind ActionKind
}

// ActionKind identifies the type of reusable component in GitHub Action.
type ActionKind uint8

const (
	_ ActionKind = iota

	// Action represent a GitHub Actions component that is an "action".
	//
	// An action is a path in a repository that has an Action manifest, i.e. a
	// file named either action.yml, action.yaml, or Dockerfile. These are used
	// in the `uses:` value of steps.
	Action

	// ReusableWorkflow represent a GitHub Actions component that is a "reusable
	// workflow".
	//
	// A reusable workflow is a workflow in a repository with the appropriate
	// workflow trigger. These are used in the `uses:` value of workflow jobs.
	ReusableWorkflow
)

// WorkflowsPath is the relative path to the GitHub Actions workflow directory.
var WorkflowsPath = filepath.Join(".github", "workflows")

// RepoActions extracts the GitHub Actions used in the repository at the given
// file system hierarchy.
func RepoActions(repo fs.FS) ([]GitHubAction, error) {
	rawWorkflows, err := workflowsInRepo(repo)
	if err != nil {
		return nil, err
	}

	workflows := make([]workflow, len(rawWorkflows))
	for i, rawWorkflow := range rawWorkflows {
		w, parseErr := parseWorkflow(rawWorkflow.content)
		if parseErr != nil {
			return nil, fmt.Errorf("%v for %q", parseErr, rawWorkflow.path)
		}

		workflows[i] = w
	}

	actions, err := actionsInWorkflows(workflows)
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// WorkflowActions extracts the GitHub Actions used in the specified workflow at
// the given file system hierarchy.
func WorkflowActions(repo fs.FS, path string) ([]GitHubAction, error) {
	data, err := workflowInRepo(repo, path)
	if err != nil {
		return nil, err
	}

	w, err := parseWorkflow(data)
	if err != nil {
		return nil, err
	}

	actions, err := actionsInWorkflows([]workflow{w})
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// JobActions extracts the GitHub Actions used in the specified job in the
// specified workflow at the given file system hierarchy.
func JobActions(repo fs.FS, path, name string) ([]GitHubAction, error) {
	data, err := workflowInRepo(repo, path)
	if err != nil {
		return nil, err
	}

	w, err := parseWorkflow(data)
	if err != nil {
		return nil, err
	}

	for job := range w.Jobs {
		if job != name {
			delete(w.Jobs, job)
		}
	}

	if len(w.Jobs) == 0 {
		return nil, fmt.Errorf("job %q not found in workflow %q", name, path)
	}

	actions, err := actionsInWorkflows([]workflow{w})
	if err != nil {
		return nil, err
	}

	return actions, nil
}

// ManifestActions extracts the GitHub Actions used in the manifest in the
// specified directory in the given file system hierarchy.
func ManifestActions(repo fs.FS, path string) ([]GitHubAction, error) {
	data, err := manifestInRepo(repo, path)
	if errors.Is(err, ErrDockerfileManifest) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	m, err := parseManifest(data)
	if err != nil {
		return nil, fmt.Errorf("could not parse manifest at: %v", err)
	}

	actions, err := actionsInManifest(m)
	if err != nil {
		return nil, fmt.Errorf("could not extract actions from manifest at: %v", err)
	}

	return actions, nil
}

func (a GitHubAction) String() string {
	if a.Path == "" {
		return fmt.Sprintf("%s/%s@%s", a.Owner, a.Project, a.Ref)
	} else {
		return fmt.Sprintf("%s/%s/%s@%s", a.Owner, a.Project, a.Path, a.Ref)
	}
}
