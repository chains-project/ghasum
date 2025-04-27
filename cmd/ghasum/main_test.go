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

package main

import (
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	commands := map[string]func(){
		"ghasum": main,
	}

	testscript.Main(m, commands)
}

func TestCli(t *testing.T) {
	t.Parallel()

	params := testscript.Params{
		Dir: "../../testdata",
	}

	testscript.Run(t, params)
}

func TestExitCodes(t *testing.T) {
	t.Parallel()

	if got := exitCodeSuccess; got != 0 {
		t.Fatalf("Exit code for success must be 0 (got %d)", got)
	}

	exitCodes := []int{
		exitCodeSuccess,
		exitCodeError,
		exitCodeUsage,
		exitCodeFailure,
	}

	for i, a := range exitCodes {
		for j, b := range exitCodes {
			if i != j && a == b {
				t.Fatalf("Exit codes must be unique (%d and %d are identical)", i, j)
			}
		}
	}
}

func TestCommands(t *testing.T) {
	t.Parallel()

	for name, fn := range commands {
		name, fn := name, fn
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if fn == nil {
				t.Fatal("Command should not be nil")
			}
		})
	}
}

func TestHelpers(t *testing.T) {
	t.Parallel()

	for name, fn := range helpers {
		name, fn := name, fn
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if fn == nil {
				t.Fatal("Helper should not be nil")
			}

			got := fn()
			if want := "ghasum " + name; !strings.Contains(got, want) {
				t.Errorf("Help is missing substring %q", want)
			}

			if !strings.HasSuffix(got, "\n") {
				t.Errorf("Help should end with a trailing newline")
			}

			if strings.Contains(got, "<command>") {
				if want := "The available commands are:"; !strings.Contains(got, want) {
					t.Errorf("Command accepts a command but does not list them (missing %q)", want)
				}
			}

			if strings.Contains(got, "[flags]") {
				if want := "The available flags are:"; !strings.Contains(got, want) {
					t.Errorf("Command accepts flags but does not describe them (missing %q)", want)
				}
			}
		})
	}
}
