// Copyright 2025 Eric Cornelissen
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

	"github.com/chains-project/ghasum/internal/cache"
	"github.com/chains-project/ghasum/internal/ghasum"
)

func cmdList(argv []string) error {
	var (
		flags            = flag.NewFlagSet(cmdNameUpdate, flag.ContinueOnError)
		flagCache        = flags.String(flagNameCache, "", "")
		flagNoCache      = flags.Bool(flagNameNoCache, false, "")
		flagNoEvict      = flags.Bool(flagNameNoEvict, false, "")
		flagNoTransitive = flags.Bool(flagNameNoTransitive, false, "")
	)

	flags.Usage = func() { fmt.Fprintln(os.Stderr) }
	if err := flags.Parse(argv); err != nil {
		return errUsage
	}

	args := flags.Args()
	if len(args) > 1 {
		return errUsage
	}

	target, err := getTarget(args)
	if err != nil {
		return err
	}

	c, err := cache.New(*flagCache, *flagNoCache)
	if err != nil {
		return errors.Join(errCache, err)
	}

	if !*flagNoEvict {
		if err = c.Evict(); err != nil {
			return errors.Join(errCache, err)
		}
	}

	repo, err := os.OpenRoot(target)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	cfg := ghasum.Config{
		Repo:       repo.FS(),
		Path:       target,
		Cache:      c,
		Transitive: !(*flagNoTransitive),
	}

	list, err := ghasum.List(&cfg)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	fmt.Print(list)
	return nil
}

func helpList() string {
	return `usage: ghasum list [flags] [target]

List the GitHub Actions dependencies for the target. If no target is provided it
will default to the current working directory.

The available flags are:

    -cache dir
        The location of the cache directory. This is where ghasum stores and
        looks up repositories it needs.
        Defaults to a directory named .ghasum in the user's home directory.
    -no-cache
        Disable the use of the cache. Makes the -cache flag ineffective.
    -no-evict
        Disable cache eviction.
    -no-transitive
        Do not include transitive actions.
`
}
