// Copyright 2024 Eric Cornelissen
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
	"io/fs"
	"path"
)

// A GitHubAction identifies a specific version of a GitHub Action.
type GitHubAction struct {
	// Owner is the GitHub user or organization that owns the repository that
	// houses the GitHub Action.
	Owner string

	// Project is the name of the GitHub repository (excluding the owner) that
	// houses the GitHub Action.
	Project string

	// Ref is the git ref (branch, tag, commit SHA), also known as version, of the
	// GitHub Action.
	Ref string
}

// WorkflowsPath is the relative path to the GitHub Actions workflow directory.
var WorkflowsPath = path.Join(".github", "workflows")

// RepoActions extracts the GitHub RepoActions used in the repository at the
// given file system hierarchy.
func RepoActions(repo fs.FS) ([]GitHubAction, error) {
	rawWorkflows, err := workflowsInRepo(repo)
	if err != nil {
		return nil, err
	}

	workflows := make([]workflow, len(rawWorkflows))
	for i, rawWorkflow := range rawWorkflows {
		w, parseErr := parseWorkflow(rawWorkflow)
		if parseErr != nil {
			return nil, parseErr
		}

		workflows[i] = w
	}

	actions, err := actionsInWorkflows(workflows)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
