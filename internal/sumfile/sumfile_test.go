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
	"slices"
	"testing"
	"testing/quick"
)

func TestAnyVersion(t *testing.T) {
	t.Parallel()

	decodable := func(version Version, entries []Entry) bool {
		version = (version % VersionLatest) + 1 // normalize version

		encoded, err := Encode(version, entries)
		if err != nil {
			return true
		}

		decoded, err := Decode(encoded)
		if err != nil {
			return false
		}

		return SetEqual(decoded, entries)
	}

	if err := quick.Check(decodable, nil); err != nil {
		t.Errorf("decode(encode(x)) errored for: %v", err)
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		sumfile string
		want    []Entry
	}

	testCases := map[string]TestCase{
		"basic example": {
			sumfile: `version 1

actions/checkout@v4.2.0 e6ng7MJDyAPkTZ/6d/plZK2YhZRzJZvxhYAPUPpNAzc=
`,
			want: []Entry{
				{
					Checksum: "e6ng7MJDyAPkTZ/6d/plZK2YhZRzJZvxhYAPUPpNAzc=",
					ID:       []string{"actions/checkout", "v4.2.0"},
				},
			},
		},
		"windows newlines": {
			sumfile: "version 1\r\n\r\nactions/checkout@v4.2.0 e6ng7MJDyAPkTZ/6d/plZK2YhZRzJZvxhYAPUPpNAzc=\r\n",
			want: []Entry{
				{
					Checksum: "e6ng7MJDyAPkTZ/6d/plZK2YhZRzJZvxhYAPUPpNAzc=",
					ID:       []string{"actions/checkout", "v4.2.0"},
				},
			},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := Decode(tt.sumfile)
			if err != nil {
				t.Fatalf("Want no error, got %v", err)
			}

			if got, want := len(got), len(tt.want); got != want {
				t.Fatalf("Want %d entries, got %d entries", got, want)
			}

			for i, got := range got {
				want := tt.want[i]

				if got, want := got.Checksum, want.Checksum; got != want {
					t.Errorf("Unexpected %dth checksum (got %q, want %q)", i, got, want)
				}

				if got, want := got.ID, want.ID; !slices.Equal(got, want) {
					t.Errorf("Unexpected %dth checksum (got %q, want %q)", i, got, want)
				}
			}
		})
	}
}

func TestNoChecksums(t *testing.T) {
	t.Parallel()

	t.Run("Decode", func(t *testing.T) {
		t.Parallel()

		entries, err := Decode("version 1\n")
		if err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}

		if got, want := len(entries), 0; got != want {
			t.Errorf("Incorrect result count (got %d, want %d)", got, want)
		}
	})

	t.Run("Encode", func(t *testing.T) {
		t.Parallel()

		if _, err := Encode(1, []Entry{}); err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}
	})
}

func TestUnknownVersion(t *testing.T) {
	t.Parallel()

	t.Run("Decode", func(t *testing.T) {
		t.Parallel()

		if _, err := Decode("version 0\n"); err == nil {
			t.Fatal("Unexpected success")
		}
	})

	t.Run("Encode", func(t *testing.T) {
		t.Parallel()

		if _, err := Encode(0, []Entry{}); err == nil {
			t.Fatal("Unexpected success")
		}
	})
}

func TestDecodeCorruptFile(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"",
		" ",
		"version",
		"version ",
		"not a version",
		"version 1",
		`version 1
version 2
`,
		`version 2
version 1
`,
		`version 1
version 1
`,
		`version 1
duplicate header
duplicate header
`,
		`version 1

duplicate checksum
duplicate checksum
`,
		`version 1

missing final newline`,
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			if _, err := Decode(tc); err == nil {
				t.Fatal("Unexpected success")
			}
		})
	}
}
