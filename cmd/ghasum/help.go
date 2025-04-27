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
	"errors"
	"flag"
	"fmt"
	"os"
)

func cmdHelp(argv []string) error {
	flags := flag.NewFlagSet(cmdNameHelp, flag.ContinueOnError)

	flags.Usage = func() { fmt.Fprintln(os.Stderr) }
	if err := flags.Parse(argv); err != nil {
		return errUsage
	}

	args := flags.Args()
	switch len(args) {
	case 0:
		fmt.Print(help())
		return nil
	case 1:
		return helpFor(args[0])
	default:
		return errors.New("you can ask help for only one command at the time")
	}
}

func helpFor(command string) error {
	fn, ok := helpers[command]
	if !ok {
		return fmt.Errorf(`unknown command %q (see "ghasum help")`, command)
	}

	fmt.Print(fn())
	return nil
}

func help() string {
	return `usage: ghasum <command> [arguments]

Checksums manager for the GitHub Action ecosystem. Track and verify checksums
for the GitHub Actions used in a project to avoid using Actions that changed.

The available commands are:

    cache     Manage the ghasum cache.
    init      Initialize ghasum for a repository.
    update    Update the checksums for a repository.
    verify    Verify the checksums for a repository.
    version   Print the ghasum version.

Use "ghasum help <command>" for more information about a command.
`
}
