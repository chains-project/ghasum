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
	"path"
	"path/filepath"
	"strings"

	"github.com/chains-project/ghasum/internal/cache"
	"github.com/chains-project/ghasum/internal/ghasum"
)

func cmdVerify(argv []string) error {
	var (
		flags            = flag.NewFlagSet(cmdNameVerify, flag.ContinueOnError)
		flagCache        = flags.String(flagNameCache, "", "")
		flagNoCache      = flags.Bool(flagNameNoCache, false, "")
		flagNoEvict      = flags.Bool(flagNameNoEvict, false, "")
		flagNoTransitive = flags.Bool(flagNameNoTransitive, false, "")
		flagOffline      = flags.Bool(flagNameOffline, false, "")
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

	var job string
	if i := strings.LastIndexByte(target, 0x3A); i > 1 {
		job = target[i+1:]
		target = target[0:i]
	}

	stat, err := os.Stat(target)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	var workflow string
	if !stat.IsDir() {
		repo := path.Join(path.Dir(target), "..", "..")
		workflow, _ = filepath.Rel(repo, target)
		workflow = strings.ReplaceAll(workflow, string(filepath.Separator), "/")
		target = repo
	}

	c, err := cache.New(
		cache.WithLocation(*flagCache),
		cache.WithEviction(!*flagNoEvict),
		cache.WithEphemeralCache(*flagNoCache),
	)
	if err != nil {
		return errors.Join(errCache, err)
	}

	repo, err := os.OpenRoot(target)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	cfg := ghasum.Config{
		Repo:       repo.FS(),
		Path:       target,
		Workflow:   workflow,
		Job:        job,
		Cache:      c,
		Offline:    *flagOffline,
		Transitive: !(*flagNoTransitive),
	}

	report, err := ghasum.Verify(&cfg)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	err = reportVerify(&report)
	return err
}

func reportVerify(report *ghasum.VerifyReport) error {
	if cnt := len(report.Problems); cnt > 0 {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("%d problem(s) occurred during validation:\n", cnt))
		for _, problem := range report.Problems {
			sb.WriteString(fmt.Sprintf("  %s\n", problem))
		}

		return errors.Join(errFailure, errors.New(sb.String()))
	}

	if report.Total == 1 {
		fmt.Println("Ok (verified 1 action)")
	} else {
		fmt.Printf("Ok (verified %d actions)\n", report.Total)
	}

	return nil
}

func helpVerify() string {
	return `usage: ghasum verify [flags] [target]

Verify the Actions in the target against the stored checksums. If no target is
provided it will default to the current working directory. If the checksums do
not match this command will error with a non-zero exit code. If ghasum is not
yet initialized this command errors (see "ghasum help init").

The target can be either a directory or a file. If it is a directory it must be
the root of a repository (that is, it should contain the .github directory). For
example:

    ghasum verify my-project

In this case checksums will be verified for every workflow in the repository. If
it is a file it must be a workflow file in a repository. For example:

    ghasum verify my-project/.github/workflows/workflow.yml

In this case checksums will be verified for all jobs in the given workflow. If
it is a file it may specify a job by using a ":job" suffix. For example:

    ghasum verify my-project/.github/workflows/workflow.yml:job-key

In this case checksums will be verified only for the given job in the workflow.

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
        Do not verify checksums for transitive actions.
    -offline
        Run without fetching repositories from the internet, verify exclusively
        against the cache. If the cache is missing an entry it causes an error.
`
}
