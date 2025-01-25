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

package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.3.0"

func cmdVersion(argv []string) error {
	var (
		flags = flag.NewFlagSet(cmdNameVersion, flag.ContinueOnError)
	)

	flags.Usage = func() { fmt.Fprintln(os.Stderr) }
	if err := flags.Parse(argv); err != nil {
		return errUsage
	}

	args := flags.Args()
	if len(args) != 0 {
		return errUsage
	}

	fmt.Printf("v%s\n", version)
	return nil
}

func helpVersion() string {
	return `usage: ghasum version

Prints the version of ghasum.`
}
