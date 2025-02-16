// Copyright 2023-2025 Eric Cornelissen
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
	"fmt"
	"os"
)

type (
	// A Command is a function that performs a ghasum command.
	Command func(args []string) error

	// A Helper is a function that returns the help text for a ghasum command.
	Helper func() string
)

const (
	cmdNameCache   = "cache"
	cmdNameHelp    = "help"
	cmdNameInit    = "init"
	cmdNameUpdate  = "update"
	cmdNameVerify  = "verify"
	cmdNameVersion = "version"
)

const (
	exitCodeSuccess = iota
	exitCodeError
	exitCodeUsage
	exitCodeFailure
)

const (
	flagNameCache   = "cache"
	flagNameForce   = "force"
	flagNameNoCache = "no-cache"
	flagNameNoEvict = "no-evict"
	flagNameOffline = "offline"
)

var (
	errCache      = errors.New("cache error (using -cache or -no-cache may avoid this error)")
	errFailure    = errors.New("")
	errUsage      = errors.New("")
	errUnexpected = errors.New("an unexpected error occurred")
)

var commands = map[string]Command{
	cmdNameCache:   cmdCache,
	cmdNameHelp:    cmdHelp,
	cmdNameInit:    cmdInit,
	cmdNameUpdate:  cmdUpdate,
	cmdNameVerify:  cmdVerify,
	cmdNameVersion: cmdVersion,
}

var helpers = map[string]Helper{
	cmdNameCache:   helpCache,
	cmdNameHelp:    help,
	cmdNameInit:    helpInit,
	cmdNameUpdate:  helpUpdate,
	cmdNameVerify:  helpVerify,
	cmdNameVersion: helpVersion,
}

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) < 2 {
		fmt.Println(help())
		return exitCodeSuccess
	}

	command := os.Args[1]
	fn, ok := commands[command]
	if !ok {
		fmt.Println(help())
		return exitCodeUsage
	}

	err := fn(os.Args[2:])
	switch {
	case err == nil:
		return exitCodeSuccess
	case errors.Is(err, errUsage):
		helpFn := helpers[command]
		fmt.Println(helpFn())
		return exitCodeUsage
	case errors.Is(err, errFailure):
		fmt.Println(err)
		return exitCodeFailure
	default:
		fmt.Fprintln(os.Stderr, err)
		return exitCodeError
	}
}
