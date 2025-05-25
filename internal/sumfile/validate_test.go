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

package sumfile

import (
	"testing"
)

func TestHasEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Non-empty examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string][]Entry{
			"no entries": {},
			"one ID parts": {
				{
					Checksum: "checksum",
					ID:       []string{"foobar"},
				},
			},
			"multiple ID parts": {
				{
					Checksum: "checksum",
					ID:       []string{"foo", "bar"},
				},
			},
		}

		for name, entries := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got := hasMissing(entries)
				if got {
					t.Fatal("Unexpected positive result")
				}
			})
		}
	})

	t.Run("Empty examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string][]Entry{
			"empty checksum": {
				{
					Checksum: "",
					ID:       []string{"foobar"},
				},
			},
			"empty id array": {
				{
					Checksum: "not-empty",
					ID:       []string{},
				},
			},
			"empty id part": {
				{
					Checksum: "not-empty",
					ID:       []string{""},
				},
			},
		}

		for name, entries := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got := hasMissing(entries)
				if !got {
					t.Fatal("Unexpected negative result")
				}
			})
		}
	})
}

func TestHasDuplicates(t *testing.T) {
	t.Parallel()

	t.Run("No duplicates examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string][]Entry{
			"no entries": {},
			"one part": {
				{ID: []string{"foo"}},
				{ID: []string{"bar"}},
			},
			"two parts": {
				{ID: []string{"foo", "bar"}},
				{ID: []string{"hello", "world"}},
			},
			"two parts, first differs": {
				{ID: []string{"bar", "foo"}},
				{ID: []string{"baz", "foo"}},
			},
			"two parts, second differs": {
				{ID: []string{"foo", "bar"}},
				{ID: []string{"foo", "baz"}},
			},
		}

		for name, entries := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got := hasDuplicates(entries)
				if got {
					t.Fatal("Unexpected positive result")
				}
			})
		}
	})

	t.Run("Duplicate examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string][]Entry{
			"one part": {
				{ID: []string{"foobar"}},
				{ID: []string{"foobar"}},
			},
			"multiple parts": {
				{ID: []string{"foo", "bar"}},
				{ID: []string{"foo", "bar"}},
			},
		}

		for name, entries := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got := hasDuplicates(entries)
				if !got {
					t.Fatal("Unexpected negative result")
				}
			})
		}
	})
}
