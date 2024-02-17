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

package sumfile

import (
	"reflect"
	"testing"
)

func SetEqual(got, want []Entry) bool {
OUTER_GOT:
	for _, got := range got {
		for _, want := range want {
			if reflect.DeepEqual(got, want) {
				continue OUTER_GOT
			}
		}

		return false
	}

OUTER_WANT:
	for _, want := range want {
		for _, got := range got {
			if reflect.DeepEqual(got, want) {
				continue OUTER_WANT
			}
		}

		return false
	}

	return true
}

func TestSetEqual(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name string
		a    []Entry
		b    []Entry
		want bool
	}

	testCases := []TestCase{
		{
			name: "identical",
			a: []Entry{
				{
					Checksum: "bar",
					ID:       []string{"foo"},
				},
			},
			b: []Entry{
				{
					Checksum: "bar",
					ID:       []string{"foo"},
				},
			},
			want: true,
		},
		{
			name: "in a but not in b",
			a: []Entry{
				{
					Checksum: "bar",
					ID:       []string{"foo"},
				},
			},
			b:    []Entry{},
			want: false,
		},
		{
			name: "in b but not in a",
			a:    []Entry{},
			b: []Entry{
				{
					Checksum: "bar",
					ID:       []string{"foo"},
				},
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got, want := SetEqual(tc.a, tc.b), tc.want; got != want {
				t.Errorf("Wrong result (got %t, want %t)", got, want)
			}
		})
	}
}
