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
	"slices"
	"testing"

	"github.com/liamg/memoryfs"
)

func TestNoWorkflows(t *testing.T) {
	t.Parallel()

	repo := memoryfs.New()
	if _, err := RepoActions(repo); err == nil {
		t.Fatal("Unexpected success")
	}
}

func TestFaultyWorkflow(t *testing.T) {
	t.Parallel()

	fs := map[string]mockFsEntry{
		".github": {
			Dir: true,
			Children: map[string]mockFsEntry{
				"workflows": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflow.yaml": {
							Content: []byte(workflowWithJobWithSteps),
						},
						"syntax-error.yml": {
							Content: []byte(yamlWithSyntaxError),
						},
					},
				},
			},
		},
	}

	repo, err := mockRepo(fs)
	if err != nil {
		t.Fatalf("Could not initialize file system: %+v", err)
	}

	if _, err := RepoActions(repo); err == nil {
		t.Fatal("Unexpected success")
	}
}

func TestFaultyUses(t *testing.T) {
	t.Parallel()

	t.Run("job", func(t *testing.T) {
		t.Parallel()

		fs := map[string]mockFsEntry{
			".github": {
				Dir: true,
				Children: map[string]mockFsEntry{
					"workflows": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflow.yaml": {
								Content: []byte(workflowWithJobWithSteps),
							},
							"invalid-uses.yml": {
								Content: []byte(workflowWithInvalidJobUses),
							},
						},
					},
				},
			},
		}

		repo, err := mockRepo(fs)
		if err != nil {
			t.Fatalf("Could not initialize file system: %+v", err)
		}

		if _, err := RepoActions(repo); err == nil {
			t.Fatal("Unexpected success")
		}
	})

	t.Run("step", func(t *testing.T) {
		t.Parallel()

		fs := map[string]mockFsEntry{
			".github": {
				Dir: true,
				Children: map[string]mockFsEntry{
					"workflows": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflow.yaml": {
								Content: []byte(workflowWithJobWithSteps),
							},
							"invalid-uses.yml": {
								Content: []byte(workflowWithInvalidStepUses),
							},
						},
					},
				},
			},
		}

		repo, err := mockRepo(fs)
		if err != nil {
			t.Fatalf("Could not initialize file system: %+v", err)
		}

		if _, err := RepoActions(repo); err == nil {
			t.Fatal("Unexpected success")
		}
	})
}

func TestRealisticRepository(t *testing.T) {
	t.Parallel()

	workflows := map[string]mockFsEntry{
		".github": {
			Dir: true,
			Children: map[string]mockFsEntry{
				"workflows": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"job-uses.yml": {
							Content: []byte(workflowWithJobUses),
						},
						"multiple-jobs.yml": {
							Content: []byte(workflowWithJobsWithSteps),
						},
						"nested": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"foo.bar": {
									Content: []byte("foobar"),
								},
							},
						},
						"nested-action.yml": {
							Content: []byte(workflowWithNestedActions),
						},
						"not-a-workflow.txt": {
							Content: []byte("Hello world!"),
						},
						"one-job.yaml": {
							Content: []byte(workflowWithJobWithSteps),
						},
					},
				},
			},
		},
	}

	repo, err := mockRepo(workflows)
	if err != nil {
		t.Fatalf("Could not initialize file system: %+v", err)
	}

	got, err := RepoActions(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	want := []GitHubAction{
		{
			Owner:   "foo",
			Project: "bar",
			Ref:     "v1",
			Kind:    Action,
		},
		{
			Owner:   "foo",
			Project: "baz",
			Ref:     "v2",
			Kind:    Action,
		},
		{
			Owner:   "nested",
			Project: "action",
			Path:    "1",
			Ref:     "v1",
			Kind:    Action,
		},
		{
			Owner:   "nested",
			Project: "action",
			Path:    "2",
			Ref:     "v1",
			Kind:    Action,
		},
		{
			Owner:   "reusable",
			Project: "workflow",
			Path:    ".github/workflows/workflow.yml",
			Ref:     "v1",
			Kind:    ReusableWorkflow,
		},
	}

	if got, want := len(got), len(want); got != want {
		t.Errorf("Incorrect result length (got %d, want %d)", got, want)
	}

	for _, got := range got {
		if !slices.Contains(want, got) {
			t.Errorf("Unwanted value found %v", got)
		}
	}

	for _, want := range want {
		if !slices.Contains(got, want) {
			t.Errorf("Wanted value missing %v", want)
		}
	}
}

func TestWorkflowActions(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		fs       map[string]mockFsEntry
		workflow string
		wantErr  bool
	}

	testCases := []TestCase{
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithNoJobs),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobNoSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobUses),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobWithSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobsWithSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithNestedActions),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(yamlWithSyntaxError),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  true,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithInvalidJobUses),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  true,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithInvalidStepUses),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  true,
		},
		{
			fs:       map[string]mockFsEntry{},
			workflow: ".github/workflows/workflow.yml",
			wantErr:  true,
		},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			t.Parallel()

			repo, err := mockRepo(tt.fs)
			if err != nil {
				t.Fatalf("Could not initialize file system: %+v", err)
			}

			_, err = WorkflowActions(repo, tt.workflow)
			if err == nil && tt.wantErr {
				t.Error("Unexpected success")
			} else if err != nil && !tt.wantErr {
				t.Errorf("Unexpected failure (got %v)", err)
			}
		})
	}
}

func TestJobActions(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		fs       map[string]mockFsEntry
		workflow string
		job      string
		wantErr  bool
	}

	testCases := []TestCase{
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobNoSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "no-steps",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobWithSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "only-job",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobsWithSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "job-a",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobsWithSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "job-b",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithNestedActions),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "only-job",
			wantErr:  false,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithNoJobs),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "anything",
			wantErr:  true,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithJobWithSteps),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "missing",
			wantErr:  true,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(yamlWithSyntaxError),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "anything",
			wantErr:  true,
		},
		{
			fs: map[string]mockFsEntry{
				".github": {
					Dir: true,
					Children: map[string]mockFsEntry{
						"workflows": {
							Dir: true,
							Children: map[string]mockFsEntry{
								"workflow.yml": {
									Content: []byte(workflowWithInvalidStepUses),
								},
							},
						},
					},
				},
			},
			workflow: ".github/workflows/workflow.yml",
			job:      "job",
			wantErr:  true,
		},
		{
			fs:       map[string]mockFsEntry{},
			workflow: ".github/workflows/workflow.yml",
			job:      "anything",
			wantErr:  true,
		},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			t.Parallel()

			repo, err := mockRepo(tt.fs)
			if err != nil {
				t.Fatalf("Could not initialize file system: %+v", err)
			}

			_, err = JobActions(repo, tt.workflow, tt.job)
			if err == nil && tt.wantErr {
				t.Error("Unexpected success")
			} else if err != nil && !tt.wantErr {
				t.Errorf("Unexpected failure (got %v)", err)
			}
		})
	}
}

func TestManifestActions(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		fs      map[string]mockFsEntry
		path    string
		wantErr bool
	}

	testCases := []TestCase{
		{
			fs: map[string]mockFsEntry{
				"action.yml": {
					Content: []byte(manifestWithNoSteps),
				},
			},
			path:    "",
			wantErr: false,
		},
		{
			fs: map[string]mockFsEntry{
				"action.yml": {
					Content: []byte(manifestWithStep),
				},
			},
			path:    "",
			wantErr: false,
		},
		{
			fs: map[string]mockFsEntry{
				"action.yml": {
					Content: []byte(manifestWithNestedActions),
				},
			},
			path:    "",
			wantErr: false,
		},
		{
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
			path:    "nested",
			wantErr: false,
		},
		{
			fs: map[string]mockFsEntry{
				"action.yaml": {
					Content: []byte(manifestWithStep),
				},
			},
			path:    "",
			wantErr: false,
		},
		{
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
			path:    "nested",
			wantErr: false,
		},
		{
			fs: map[string]mockFsEntry{
				"action.yml": {
					Content: []byte(yamlWithSyntaxError),
				},
			},
			path:    "",
			wantErr: true,
		},
		{
			fs: map[string]mockFsEntry{
				"action.yml": {
					Content: []byte(manifestWithInvalidUses),
				},
			},
			path:    "",
			wantErr: true,
		},
		{
			fs:      map[string]mockFsEntry{},
			path:    "",
			wantErr: true,
		},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			t.Parallel()

			repo, err := mockRepo(tt.fs)
			if err != nil {
				t.Fatalf("Could not initialize file system: %+v", err)
			}

			_, err = ManifestActions(repo, tt.path)
			if err == nil && tt.wantErr {
				t.Error("Unexpected success")
			} else if err != nil && !tt.wantErr {
				t.Errorf("Unexpected failure (got %v)", err)
			}
		})
	}
}
