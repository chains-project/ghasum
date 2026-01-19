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
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/quick"
)

func TestVersion1(t *testing.T) {
	t.Parallel()

	correct := func(entries []Entry) bool {
		if err := validV1(entries); err != nil {
			return true
		}

		encoded, _ := encodeV1(entries)
		lines := strings.Split(encoded, "\n")

		decoded, err := decodeV1(lines[:len(lines)-1])
		if err != nil {
			return true // Ignore errors, tested separately
		}

		return SetEqual(decoded, entries)
	}

	if err := quick.Check(correct, nil); err != nil {
		t.Errorf("decode(encode(x)) != x for: %v", err)
	}

	decodable := func(entries []Entry) bool {
		if err := validV1(entries); err != nil {
			return true
		}

		encoded, _ := encodeV1(entries)
		lines := strings.Split(encoded, "\n")

		_, err := decodeV1(lines[:len(lines)-1])
		return err == nil
	}

	if err := quick.Check(decodable, nil); err != nil {
		t.Errorf("decode(encode(x)) errored for: %v", err)
	}

	deterministic := func(entries []Entry) bool {
		got1, err1 := encodeV1(entries)
		got2, err2 := encodeV1(entries)
		return got1 == got2 && ((err1 == nil) == (err2 == nil))
	}

	if err := quick.Check(deterministic, nil); err != nil {
		t.Errorf("encode(x) != encode(x) for: %v", err)
	}
}

func TestDecodeV1(t *testing.T) {
	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			content []string
			want    []Entry
		}

		testCases := map[string]TestCase{
			"no checksums": {
				content: []string{},
				want:    []Entry{},
			},
			"one checksum": {
				content: []string{
					"foo bar",
				},
				want: []Entry{
					{
						Checksum: "bar",
						ID:       []string{"foo"},
					},
				},
			},
			"one multi-part ID checksum": {
				content: []string{
					"foo@bar foobar",
				},
				want: []Entry{
					{
						Checksum: "foobar",
						ID:       []string{"foo", "bar"},
					},
				},
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := decodeV1(tt.content)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got), len(tt.want); got != want {
					t.Fatalf("Incorrect result length (got %d, want %d)", got, want)
				}

				for i, got := range got {
					want := tt.want[i]

					if got, want := got.Checksum, want.Checksum; got != want {
						t.Fatalf("Incorrect checksum %d (got %q, want %q)", i, got, want)
					}

					if got, want := got.ID, want.ID; !slices.Equal(got, want) {
						t.Fatalf("Incorrect id %d (got %v, want %v)", i, got, want)
					}
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			content []string
			want    int
		}

		testCases := map[string]TestCase{
			"no id-checksum separator": {
				content: []string{
					"foobar",
				},
				want: 3,
			},
			"no checksum": {
				content: []string{
					"foobar ",
				},
				want: 3,
			},
			"no id": {
				content: []string{
					" foobar",
				},
				want: 3,
			},
			"on a later line": {
				content: []string{
					"foo bar",
					"syntax-error",
				},
				want: 4,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				_, err := decodeV1(tt.content)
				if err == nil {
					t.Fatal("Unexpected success")
				}

				if got, want := err.Error(), fmt.Sprintf("line %d", tt.want); !strings.Contains(got, want) {
					t.Errorf("Incorrect line number (got %q, want %q)", got, want)
				}
			})
		}
	})
}

func TestEncodeV1(t *testing.T) {
	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			content []Entry
			want    string
		}

		testCases := map[string]TestCase{
			"no checksums": {
				content: []Entry{},
				want:    ``,
			},
			"one checksum": {
				content: []Entry{
					{
						Checksum: "bar",
						ID:       []string{"foo"},
					},
				},
				want: `foo bar
`,
			},
			"one multi-part ID checksum": {
				content: []Entry{
					{
						Checksum: "foobar",
						ID:       []string{"foo", "bar"},
					},
				},
				want: `foo@bar foobar
`,
			},
			"order": {
				content: []Entry{
					{
						Checksum: "bb",
						ID:       []string{"b"},
					},
					{
						Checksum: "aa",
						ID:       []string{"a"},
					},
				},
				want: `a aa
b bb
`,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				got, err := encodeV1(tt.content)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if want := tt.want; got != want {
					t.Fatalf("Incorrect result (got %q, want %q)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		testCases := map[string][]Entry{
			"checksum with newline": {
				{
					ID:       []string{"anything"},
					Checksum: "Hello\nworld!",
				},
			},
			"checksum with carriage return": {
				{
					ID:       []string{"anything"},
					Checksum: "Hello\rworld!",
				},
			},
			"checksum with space": {
				{
					ID:       []string{"anything"},
					Checksum: "Hello world!",
				},
			},
			"ID part with newline": {
				{
					ID:       []string{"Hello\nworld!"},
					Checksum: "anything",
				},
			},
			"ID part with carriage return": {
				{
					ID:       []string{"Hello\rworld!"},
					Checksum: "anything",
				},
			},
			"ID part with space": {
				{
					ID:       []string{"Hello world!"},
					Checksum: "anything",
				},
			},
			"ID part with '@'": {
				{
					ID:       []string{"foo@bar"},
					Checksum: "anything",
				},
			},
		}

		for name, entries := range testCases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				if _, err := encodeV1(entries); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})
}
