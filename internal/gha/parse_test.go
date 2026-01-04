// Copyright 2024-2026 Eric Cornelissen
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
	"strings"
	"testing"
	"testing/quick"
)

func TestParseUses(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			uses string
			want GitHubAction
		}

		testCases := map[string]TestCase{
			"versioned action, branch ref": {
				uses: "foo/bar@main",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Ref:     "main",
				},
			},
			"versioned action, moving tag ref": {
				uses: "foo/bar@v1",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Ref:     "v1",
				},
			},
			"versioned action, specific version tag": {
				uses: "foo/baz@v3.1.4",
				want: GitHubAction{
					Owner:   "foo",
					Project: "baz",
					Ref:     "v3.1.4",
				},
			},
			"versioned action, non-standard ref": {
				uses: "hello/world@random-ref",
				want: GitHubAction{
					Owner:   "hello",
					Project: "world",
					Ref:     "random-ref",
				},
			},
			"versioned action, commit ref": {
				uses: "hallo/wereld@35dd46a3b3dfbb14198f8d19fb083ce0832dce4a",
				want: GitHubAction{
					Owner:   "hallo",
					Project: "wereld",
					Ref:     "35dd46a3b3dfbb14198f8d19fb083ce0832dce4a",
				},
			},
			"versioned action, subdirectory": {
				uses: "foo/bar/baz@v2",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Path:    "baz",
					Ref:     "v2",
				},
			},
			"reusable workflow": {
				uses: "octo-org/another-repo/.github/workflows/workflow.yml@v1",
				want: GitHubAction{
					Owner:   "octo-org",
					Project: "another-repo",
					Path:    ".github/workflows/workflow.yml",
					Ref:     "v1",
				},
			},
			"uppercase in owner": {
				uses: "Foo/bar/baz@v42",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Path:    "baz",
					Ref:     "v42",
				},
			},
			"uppercase in project": {
				uses: "foo/Bar/baz@v42",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Path:    "baz",
					Ref:     "v42",
				},
			},
			"uppercase in path": {
				uses: "foo/bar/Baz@v42",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Path:    "Baz",
					Ref:     "v42",
				},
			},
			"uppercase in ref": {
				uses: "foo/bar/baz@V42",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Path:    "baz",
					Ref:     "V42",
				},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := parseUses(tt.uses)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := got.Owner, tt.want.Owner; got != want {
					t.Errorf("Incorrect owner (got %q, want %q)", got, want)
				}

				if got, want := got.Project, tt.want.Project; got != want {
					t.Errorf("Incorrect project (got %q, want %q)", got, want)
				}

				if got, want := got.Path, tt.want.Path; got != want {
					t.Errorf("Incorrect path (got %q, want %q)", got, want)
				}

				if got, want := got.Ref, tt.want.Ref; got != want {
					t.Errorf("Incorrect ref (got %q, want %q)", got, want)
				}

				if got, want := got.Kind, tt.want.Kind; got != want {
					t.Errorf("Incorrect kind (got %q, want %q)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			uses string
			want error
		}

		testCases := map[string]TestCase{
			"an action in the same repository as the workflow": {
				uses: "./.github/actions/hello-world-action",
				want: ErrLocalAction,
			},
			"a reusable workflow in the same repository as the workflow": {
				uses: "./.github/workflow/reusable.yml",
				want: ErrLocalAction,
			},
			"a Docker Hub action": {
				uses: "docker://alpine:3.8",
				want: ErrDockerUses,
			},
			"a GitHub Package Container registry action": {
				uses: "docker://ghcr.io/OWNER/IMAGE_NAME",
				want: ErrDockerUses,
			},
			"plain string": {
				uses: "foobar",
				want: ErrInvalidUses,
			},
			"only a /, no @": {
				uses: "foo/bar",
				want: ErrInvalidUses,
			},
			"one extra @ before /": {
				uses: "f@o/bar@baz",
				want: ErrInvalidUses,
			},
			"one extra @ after /": {
				uses: "foo/b@r@baz",
				want: ErrInvalidUses,
			},
			"one extra @ in path": {
				uses: "foo/bar/b@z@ref",
				want: ErrInvalidUses,
			},
			"only a @, no /": {
				uses: "foo@bar",
				want: ErrInvalidUsesRepo,
			},
			"empty repository name (no path)": {
				uses: "foo/@bar",
				want: ErrInvalidUsesRepo,
			},
			"empty path": {
				uses: "foo/bar/@baz",
				want: ErrInvalidUsesPath,
			},
			"empty repository name and path": {
				uses: "foo//@bar",
				want: ErrInvalidUsesPath,
			},
			"empty repository name with path": {
				uses: "foo//bar@baz",
				want: ErrInvalidUsesPath,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_, err := parseUses(tt.uses)
				if err == nil {
					t.Fatal("Unexpected success")
				}

				if got, want := err, tt.want; !errors.Is(got, want) {
					t.Errorf("Incorrect error (got %q, want %q)", got, want)
				}
			})
		}
	})

	t.Run("Arbitrary values", func(t *testing.T) {
		t.Parallel()

		constructive := func(owner, project, path, ref string) bool {
			if len(owner) == 0 || len(project) == 0 || len(ref) == 0 {
				return true
			}

			repo := fmt.Sprintf("%s/%s", owner, project)
			if strings.Count(repo, "/") != 1 {
				return true
			}

			if len(path) > 0 {
				repo = fmt.Sprintf("%s/%s", repo, path)
			}

			if strings.ContainsRune(repo, '@') || strings.ContainsRune(ref, '@') {
				return true
			}

			uses := fmt.Sprintf("%s@%s", repo, ref)

			action, err := parseUses(uses)
			if err != nil {
				return false
			}

			return action.Owner == strings.ToLower(owner) &&
				action.Project == strings.ToLower(project) &&
				action.Path == path &&
				action.Ref == ref
		}

		if err := quick.Check(constructive, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}

		noPanic := func(uses string) bool {
			_, _ = parseUses(uses)
			return true
		}

		if err := quick.Check(noPanic, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}
	})
}

func TestParseWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   string
			want workflow
		}

		testCases := map[string]TestCase{
			"workflow with no jobs": {
				in: workflowWithNoJobs,
				want: workflow{
					Jobs: map[string]job{},
				},
			},
			"workflow with a job with no steps": {
				in: workflowWithJobNoSteps,
				want: workflow{
					Jobs: map[string]job{
						"no-steps": {},
					},
				},
			},
			"workflow with a 'uses:' job": {
				in: workflowWithJobUses,
				want: workflow{
					Jobs: map[string]job{
						"uses": {
							Uses: "reusable/workflow/.github/workflows/workflow.yml@v1",
						},
					},
				},
			},
			"workflow with one job and steps": {
				in: workflowWithJobWithSteps,
				want: workflow{
					Jobs: map[string]job{
						"only-job": {
							Steps: []step{
								{
									Uses: "foo/bar@v1",
								},
								{
									Uses: "",
								},
								{
									Uses: "foo/baz@v2",
								},
							},
						},
					},
				},
			},
			"workflow with multiple jobs and steps": {
				in: workflowWithJobsWithSteps,
				want: workflow{
					Jobs: map[string]job{
						"job-a": {
							Steps: []step{
								{
									Uses: "foo/bar@v1",
								},
							},
						},
						"job-b": {
							Steps: []step{
								{
									Uses: "",
								},
								{
									Uses: "foo/baz@v2",
								},
							},
						},
					},
				},
			},
			"workflow with 'owner/repo/path@v'-style 'uses:' value": {
				in: workflowWithNestedActions,
				want: workflow{
					Jobs: map[string]job{
						"only-job": {
							Steps: []step{
								{
									Uses: "nested/action/1@v1",
								},
								{
									Uses: "nested/action/2@v1",
								},
							},
						},
					},
				},
			},
			"workflow with an invalid step 'uses:' value": {
				in: workflowWithInvalidStepUses,
				want: workflow{
					Jobs: map[string]job{
						"job": {
							Steps: []step{
								{
									Uses: "this-is-not-an-action",
								},
							},
						},
					},
				},
			},
			"workflow with an invalid job 'uses:' value": {
				in: workflowWithInvalidJobUses,
				want: workflow{
					Jobs: map[string]job{
						"job": {
							Uses: "this-is-not-a-reusable-workflow",
						},
					},
				},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := parseWorkflow([]byte(tt.in))
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got.Jobs), len(tt.want.Jobs); got != want {
					t.Fatalf("Incorrect jobs length (got %d, want %d)", got, want)
				}

				for name, job := range got.Jobs {
					want, ok := tt.want.Jobs[name]
					if !ok {
						t.Errorf("Got unwanted job %q", name)
						continue
					}

					if got, want := len(job.Steps), len(want.Steps); got != want {
						t.Errorf("Incorrect steps length for job %q (got %d, want %d)", name, got, want)
						continue
					}

					for i, step := range job.Steps {
						want := want.Steps[i]

						if got, want := step.Uses, want.Uses; got != want {
							t.Errorf("Incorrect uses for step %d of job %q (got %q, want %q)", i, name, got, want)
						}
					}
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string]string{
			"workflow with YAML syntax error": yamlWithSyntaxError,
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				if _, err := parseWorkflow([]byte(tt)); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})

	t.Run("Arbitrary values", func(t *testing.T) {
		t.Parallel()

		noPanic := func(w []byte) bool {
			_, _ = parseWorkflow(w)
			return true
		}

		if err := quick.Check(noPanic, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}
	})
}

func TestParseManifest(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   string
			want manifest
		}

		testCases := map[string]TestCase{
			"manifest with no steps": {
				in:   manifestWithNoSteps,
				want: manifest{},
			},
			"manifest with steps": {
				in: manifestWithStep,
				want: manifest{
					Runs: runs{
						Steps: []step{
							{Uses: "foo/bar@v1"},
							{Uses: ""},
							{Uses: "foo/baz@v2"},
						},
					},
				},
			},
			"manifest with 'owner/repo/path@v'-style 'uses:' value": {
				in: manifestWithNestedActions,
				want: manifest{
					Runs: runs{
						Steps: []step{
							{Uses: "nested/action/1@v1"},
							{Uses: "nested/action/2@v1"},
						},
					},
				},
			},
			"manifest with invalid 'uses:' value": {
				in: manifestWithInvalidUses,
				want: manifest{
					Runs: runs{
						Steps: []step{
							{Uses: "this-is-not-an-action"},
						},
					},
				},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := parseManifest([]byte(tt.in))
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got.Runs.Steps), len(tt.want.Runs.Steps); got != want {
					t.Fatalf("Incorrect steps length (got %d, want %d)", got, want)
				}

				for i, got := range got.Runs.Steps {
					want := tt.want.Runs.Steps[i]
					if got, want := got.Uses, want.Uses; got != want {
						t.Errorf("Incorrect uses for step %d (got %q, want %q)", i, got, want)
					}
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string]string{
			"manifest with YAML syntax error": yamlWithSyntaxError,
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				if _, err := parseManifest([]byte(tt)); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})

	t.Run("Arbitrary values", func(t *testing.T) {
		t.Parallel()

		noPanic := func(w []byte) bool {
			_, _ = parseManifest(w)
			return true
		}

		if err := quick.Check(noPanic, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}
	})
}
