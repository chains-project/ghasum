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
	"bytes"
	"errors"
	"testing"
	"testing/quick"

	"github.com/liamg/memoryfs"
)

func TestActionsInManifest(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			manifest manifest
			want     int
		}

		testCases := map[string]TestCase{
			"no steps": {
				manifest: manifest{
					Runs: runs{
						Steps: []step{},
					},
				},
				want: 0,
			},
			"one step without uses": {
				manifest: manifest{
					Runs: runs{
						Steps: []step{
							{},
						},
					},
				},
				want: 0,
			},
			"one step with uses": {
				manifest: manifest{
					Runs: runs{
						Steps: []step{
							{
								Uses: "foo/bar@v1",
							},
						},
					},
				},
				want: 1,
			},
			"multiple step with unique uses": {
				manifest: manifest{
					Runs: runs{
						Steps: []step{
							{
								Uses: "foo/bar@v1",
							},
							{
								Uses: "foo/baz@v2",
							},
						},
					},
				},
				want: 2,
			},
			"multiple steps with duplicate uses": {
				manifest: manifest{
					Runs: runs{
						Steps: []step{
							{
								Uses: "foo/bar@v1",
							},
							{
								Uses: "foo/bar@v1",
							},
						},
					},
				},
				want: 1,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := actionsInManifest(tt.manifest)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got), tt.want; got != want {
					t.Errorf("Incorrect result length (got %d, want %d)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			manifest manifest
		}

		testCases := map[string]TestCase{
			"invalid uses value": {
				manifest: manifest{
					Runs: runs{
						Steps: []step{
							{
								Uses: "this isn't an action",
							},
						},
					},
				},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				if _, err := actionsInManifest(tt.manifest); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})

	t.Run("Arbitrary", func(t *testing.T) {
		t.Parallel()

		unique := func(m manifest) bool {
			actions, err := actionsInManifest(m)
			if err != nil {
				return true
			}

			seen := make(map[GitHubAction]struct{}, 0)
			for _, action := range actions {
				if _, ok := seen[action]; ok {
					return false
				}

				seen[action] = struct{}{}
			}

			return true
		}

		if err := quick.Check(unique, nil); err != nil {
			t.Errorf("Duplicate value detected for: %v", err)
		}
	})
}

func TestActionsInWorkflows(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   []workflow
			want int
		}

		testCases := map[string]TestCase{
			"no jobs": {
				in: []workflow{
					{
						Jobs: map[string]job{},
					},
				},
				want: 0,
			},
			"one job without steps": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{},
							},
						},
					},
				},
				want: 0,
			},
			"multiple jobs without steps": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example-a": {
								Steps: []step{},
							},
							"example-b": {
								Steps: []step{},
							},
						},
					},
				},
				want: 0,
			},
			"one job with a step without uses": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{},
								},
							},
						},
					},
				},
				want: 0,
			},
			"one job with one step": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
						},
					},
				},
				want: 1,
			},
			"multiple jobs with one unique step each": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example-a": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
							"example-b": {
								Steps: []step{
									{
										Uses: "foo/baz@v1",
									},
								},
							},
						},
					},
				},
				want: 2,
			},
			"one job with multiple unique steps": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
									{
										Uses: "foo/baz@v2",
									},
								},
							},
						},
					},
				},
				want: 2,
			},
			"multiple jobs with multiple unique steps": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example-a": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
									{
										Uses: "hello/world@v2",
									},
								},
							},
							"example-b": {
								Steps: []step{
									{
										Uses: "foo/baz@v1",
									},
									{
										Uses: "hallo/wereld@v2",
									},
								},
							},
						},
					},
				},
				want: 4,
			},
			"one jobs with duplicate steps": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
									{
										Uses: "foo/bar@v1",
									},
								},
							},
						},
					},
				},
				want: 1,
			},
			"multiple jobs with duplicate step between them": {
				in: []workflow{
					{

						Jobs: map[string]job{
							"example-a": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
							"example-b": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
						},
					},
				},
				want: 1,
			},
			"job with local reusable workflow": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Uses: "./.github/workflows/workflow-2.yml",
							},
						},
					},
				},
				want: 1,
			},
			"job with external reusable workflow": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Uses: "octo-org/another-repo/.github/workflows/workflow.yml@v1",
							},
						},
					},
				},
				want: 1,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := actionsInWorkflows(tt.in)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got), tt.want; got != want {
					t.Errorf("Incorrect result length (got %d, want %d)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in []workflow
		}

		testCases := map[string]TestCase{
			"invalid job uses value": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Uses: "this isn't a reusable workflow",
							},
						},
					},
				},
			},
			"invalid step uses value": {
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "this isn't an action",
									},
								},
							},
						},
					},
				},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				if _, err := actionsInWorkflows(tt.in); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})

	t.Run("Arbitrary", func(t *testing.T) {
		t.Parallel()

		unique := func(workflows []workflow) bool {
			actions, err := actionsInWorkflows(workflows)
			if err != nil {
				return true
			}

			seen := make(map[GitHubAction]struct{}, 0)
			for _, action := range actions {
				if _, ok := seen[action]; ok {
					return false
				}

				seen[action] = struct{}{}
			}

			return true
		}

		if err := quick.Check(unique, nil); err != nil {
			t.Errorf("Duplicate value detected for: %v", err)
		}
	})
}

func TestWorkflowsInRepo(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			fs   map[string]mockFsEntry
			want []workflowFile
		}

		testCases := map[string]TestCase{
			".yml workflow": {
				fs: map[string]mockFsEntry{
					".github": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflows": {
								Dir: true,
								Children: map[string]mockFsEntry{
									"example.yml": {
										Content: []byte(workflowWithJobsWithSteps),
									},
								},
							},
						},
					},
				},
				want: []workflowFile{
					{
						content: []byte(workflowWithJobsWithSteps),
						path:    ".github/workflows/example.yml",
					},
				},
			},
			".yaml workflow": {
				fs: map[string]mockFsEntry{
					".github": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflows": {
								Dir: true,
								Children: map[string]mockFsEntry{
									"example.yaml": {
										Content: []byte(workflowWithJobsWithSteps),
									},
								},
							},
						},
					},
				},
				want: []workflowFile{
					{
						content: []byte(workflowWithJobsWithSteps),
						path:    ".github/workflows/example.yaml",
					},
				},
			},
			"non-workflow file": {
				fs: map[string]mockFsEntry{
					".github": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflows": {
								Dir: true,
								Children: map[string]mockFsEntry{
									"greeting.txt": {
										Content: []byte("Hello world!"),
									},
								},
							},
						},
					},
				},
				want: []workflowFile{},
			},
			"nested directory": {
				fs: map[string]mockFsEntry{
					".github": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflows": {
								Dir: true,
								Children: map[string]mockFsEntry{
									"nested": {
										Dir: true,
										Children: map[string]mockFsEntry{
											"workflow.yml": {},
										},
									},
								},
							},
						},
					},
				},
				want: []workflowFile{},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				repo, err := mockRepo(tt.fs)
				if err != nil {
					t.Fatalf("Could not initialize file system: %+v", err)
				}

				got, err := workflowsInRepo(repo)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got), len(tt.want); got != want {
					t.Fatalf("Incorrect result length (got %d, want %d)", got, want)
				}

				for i, got := range got {
					if got, want := got.content, tt.want[i].content; !bytes.Equal(got, want) {
						t.Errorf("Incorrect content for workflow %d (got %s, want %s)", i, got, want)
					}

					if got, want := got.path, tt.want[i].path; got != want {
						t.Errorf("Incorrect path for workflow %d (got %s, want %s)", i, got, want)
					}
				}
			})
		}
	})

	t.Run("No actions", func(t *testing.T) {
		t.Parallel()

		repo := memoryfs.New()
		if _, err := workflowsInRepo(repo); err == nil {
			t.Fatal("Unexpected success")
		}
	})
}

func TestManifestInRepo(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			fs   map[string]mockFsEntry
			dir  string
			want []byte
		}

		testCases := map[string]TestCase{
			".yml manifest in root": {
				fs: map[string]mockFsEntry{
					"action.yml": {
						Content: []byte(manifestWithStep),
					},
				},
				dir:  "",
				want: []byte(manifestWithStep),
			},
			".yaml manifest in root": {
				fs: map[string]mockFsEntry{
					"action.yaml": {
						Content: []byte(manifestWithStep),
					},
				},
				dir:  "",
				want: []byte(manifestWithStep),
			},
			".yml manifest, nested": {
				fs: map[string]mockFsEntry{
					"nested": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"action.yml": {
								Content: []byte(manifestWithStep),
							},
						},
					},
				},
				dir:  "nested",
				want: []byte(manifestWithStep),
			},
			".yaml manifest, nested": {
				fs: map[string]mockFsEntry{
					"nested": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"action.yaml": {
								Content: []byte(manifestWithStep),
							},
						},
					},
				},
				dir:  "nested",
				want: []byte(manifestWithStep),
			},
			".yml and .yaml": {
				fs: map[string]mockFsEntry{
					"action.yaml": {
						Content: []byte(yamlWithSyntaxError),
					},
					"action.yml": {
						Content: []byte(manifestWithStep),
					},
				},
				dir:  "",
				want: []byte(manifestWithStep),
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				repo, err := mockRepo(tt.fs)
				if err != nil {
					t.Fatalf("Could not initialize file system: %+v", err)
				}

				got, err := manifestInRepo(repo, tt.dir)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if want := tt.want; !bytes.Equal(got, want) {
					t.Errorf("Incorrect content for the manifest (got %s, want %s)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			fs   map[string]mockFsEntry
			dir  string
			want error
		}

		testCases := map[string]TestCase{
			".yml manifest in different dir": {
				fs: map[string]mockFsEntry{
					"action.yml": {
						Content: []byte(manifestWithStep),
					},
				},
				dir:  "nested",
				want: ErrNoManifest,
			},
			".yaml manifest in different dir": {
				fs: map[string]mockFsEntry{
					"action.yaml": {
						Content: []byte(manifestWithStep),
					},
				},
				dir:  "nested",
				want: ErrNoManifest,
			},
			"Dockerfile manifest": {
				fs: map[string]mockFsEntry{
					"Dockerfile": {
						Content: []byte(manifestDockerfile),
					},
				},
				dir:  "",
				want: ErrDockerfileManifest,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				repo, err := mockRepo(tt.fs)
				if err != nil {
					t.Fatalf("Could not initialize file system: %+v", err)
				}

				_, err = manifestInRepo(repo, tt.dir)
				if err == nil {
					t.Fatal("Unexpected success")
				}

				if got, want := err, tt.want; !errors.Is(got, want) {
					t.Errorf("Unexpected error (got %q, want %q)", got, want)
				}
			})
		}
	})
}
