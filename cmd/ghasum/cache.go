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

	"github.com/chains-project/ghasum/internal/cache"
)

func cmdCache(argv []string) error {
	var (
		flags     = flag.NewFlagSet(cmdNameCache, flag.ContinueOnError)
		flagCache = flags.String(flagNameCache, "", "")
	)

	flags.Usage = func() { fmt.Fprintln(os.Stderr) }
	if err := flags.Parse(argv); err != nil {
		return errUsage
	}

	args := flags.Args()
	if len(args) < 1 {
		return errUsage
	} else if len(args) > 1 {
		return errors.New("only one command can be run at the time")
	}

	c, err := cache.New(
		cache.WithLocation(*flagCache),
		cache.WithEviction(false),
	)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	var msg string
	command := args[0]
	switch command {
	case "clear":
		err = c.Clear()
		msg = "Ok"
	case "evict":
		var cnt uint
		cnt, err = c.Evict()
		msg = fmt.Sprintf("Ok (%d evicted)", cnt)
	case "path":
		msg = c.Path()
	default:
		return fmt.Errorf(`unknown command %q (see "ghasum help cache")`, command)
	}

	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	fmt.Println(msg)
	return nil
}

func helpCache() string {
	return `usage: ghasum cache [flags] <command>

Utilities for managing the ghasum cache. This cache is where ghasum stores and
looks up repositories it needs to do its job. The maximum age of entries in the
cache is 5 days, after which it will be evicted.

The available commands are:

    clear   Remove all data from the cache.
    evict   Remove old data from the cache.
    path    Show the path to the cache.

The available flags are:

    -cache dir
        The location of the cache directory. Defaults to a directory named
        .ghasum/ in the user's home directory.
`
}
