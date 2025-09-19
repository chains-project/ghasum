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
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type (
	manifest struct {
		Runs runs `yaml:"runs"`
	}

	runs struct {
		Steps []step `yaml:"steps"`
	}

	workflow struct {
		Jobs map[string]job `yaml:"jobs"`
	}

	job struct {
		Uses  string `yaml:"uses"`
		Steps []step `yaml:"steps"`
	}

	step struct {
		Uses string `yaml:"uses"`
	}
)

func parseManifest(data []byte) (manifest, error) {
	var m manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("could not parse manifest: %v", err)
	}

	return m, nil
}

func parseWorkflow(data []byte) (workflow, error) {
	var w workflow
	if err := yaml.Unmarshal(data, &w); err != nil {
		return w, fmt.Errorf("could not parse workflow: %v", err)
	}

	return w, nil
}

func parseUses(uses string) (GitHubAction, error) {
	var a GitHubAction

	// uses values that don't fit the [GitHubAction] model
	switch {
	case strings.HasPrefix(uses, "./"):
		return a, ErrLocalAction
	case strings.HasPrefix(uses, "docker://"):
		return a, ErrDockerUses
	}

	// split "uses" into "repo"@"ref"
	i := strings.IndexRune(uses, '@')
	if strings.Count(uses, "@") != 1 {
		return a, ErrInvalidUses
	}

	repo := uses[:i]
	a.Ref = uses[i+1:]

	// split "repo" into "owner"/"project[/path]"
	i = strings.IndexRune(repo, '/')
	if i <= 0 || i == len(repo)-1 {
		return a, ErrInvalidUsesRepo
	}

	a.Owner = repo[:i]
	project := repo[i+1:]

	// split "project" into "project"[/"path"]
	i = strings.IndexRune(project, '/')
	if i == 0 || i == len(project)-1 {
		return a, ErrInvalidUsesPath
	} else if i > 0 && i < len(project)-1 {
		a.Path = project[i+1:]
		project = project[:i]
	}

	a.Project = project
	return a, nil
}
