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
	"strings"
	"testing"
	"testing/quick"
)

func TestToUrl(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   Repository
			want string
		}

		testCases := map[string]TestCase{
			"Example without ref": {
				in: Repository{
					Owner:   "chains-project",
					Project: "ghasum",
				},
				want: "https://github.com/chains-project/ghasum",
			},
			"Example with ref": {
				in: Repository{
					Owner:   "with",
					Project: "ref",
					Ref:     "v1",
				},
				want: "https://github.com/with/ref",
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got := toUrl(&tt.in)
				if want := tt.want; got != want {
					t.Errorf("Incorrect result (got %q, wan %q)", got, want)
				}
			})
		}
	})

	t.Run("Arbitrary", func(t *testing.T) {
		t.Parallel()

		isGitHubUrl := func(repo Repository) bool {
			url := toUrl(&repo)
			return strings.HasPrefix(url, "https://github.com/")
		}

		if err := quick.Check(isGitHubUrl, nil); err != nil {
			t.Errorf("Missing GitHub URL for: %v", err)
		}

		containsOwner := func(repo Repository) bool {
			url := toUrl(&repo)
			return strings.Contains(url, repo.Owner)
		}

		if err := quick.Check(containsOwner, nil); err != nil {
			t.Errorf("Missing repository owner for: %v", err)
		}

		containsProject := func(repo Repository) bool {
			url := toUrl(&repo)
			return strings.Contains(url, repo.Project)
		}

		if err := quick.Check(containsProject, nil); err != nil {
			t.Errorf("Missing repository project for: %v", err)
		}
	})
}
