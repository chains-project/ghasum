<!-- SPDX-License-Identifier: CC-BY-4.0 -->

# Specification of `ghasum`

The specification aims to clarify how `ghasum` operates. Any discrepancy with
the implementation or ambiguity in the specification can be reported as a bug.
There is no guarantee on whether the specification or implementation is correct.

## Actions

### `ghasum init`

If the checksum file exists the process shall exit immediately with an error.

If the checksum file does not exist the process creates it immediately and
obtains a lock on it. If this is not possible the process should exit
immediately (it means either 1. the file has been created since it was checked
and so is not owned by us, or 2. the file could not be created and so cannot be
initialized).

If the file lock is obtained, the process will compute checksums (see [Computing
Checksums]) for all actions used in the repository (see [Collecting Actions])
using the best available hashing algorithm. Then it stores them in a sumfile
(see [Storing Checksums]) using the latest sumfile version. Finally the process
will releases the lock on the file.

If the process fails an attempt should be made to remove the created file (if
removing fails the error is ignored).

### `ghasum list`

Regardless of the existence of the checksum file, the process will find all
actions used by the target (see [Collecting Actions]) and report them in a
hierarchical (i.e., showing transitive dependency relations) to the user.

### `ghasum update`

If the checksum file does not exist the process shall exit immediately with an
error.

If the checksum file exists the process shall obtain a lock on it, if this is
not possible the process shall exit immediately. Otherwise the file could be
changed by another process potentially leading to an inconsistent sumfile.

If the file lock is obtained, the process shall first read it and parse it
completely to extract the sumfile version. If this fails the process shall exit
immediately unless the `-force` flag is used (see details below). Else it shall
compute checksums (see [Computing Checksums]) for all new actions used in the
repository (see [Collecting Actions]) using the same hashing algorithm as was
used for the existing checksums. New actions also include new versions of a
previously used actions. Additionally, it should remove any entry which is no
longer in use. No existing checksum for a used action shall be updated. It shall
then store the updated set in the checksum file (see [Storing Checksums]) using
the same sumfile version as before and releases the lock. In short, updating
will only add new and remove old checksums from an existing sumfile.

With the `-force` flag the process will ignore errors in the sumfile and fix
those while updating. It will also update existing checksums that are incorrect.
If the sumfile version can still be determined from sumfile it will be used,
otherwise the latest available version is used instead. This option is disabled
by default to avoid unknowingly fixing syntax or other errors in a sumfile,
which is an important fact to know about from a security perspective.

This process does not verify any of the checksums currently in the sumfile.

### `ghasum verify`

If the checksum file does not exist the process shall exit immediately with an
error.

If the checksum file exists the process shall read and parse it fully. If this
fails the process shall exit immediately. Else it shall recompute the checksums
(see [Computing Checksums]) for all actions in the target (see [Collecting
Actions]) using the same hashing algorithm as was used for the stored checksums.
It shall compare the computed checksums against the stored checksums.

If any of the checksums does not match or is missing the process shall exit with
a non-zero exit code. For usability all values should be compared (and all
mismatches reported) before exiting.

The "target" can be one of a: a repository, a workflow, or a job. If the target
is a repository, all actions used in all jobs in all workflows in the repository
will be considered. If the target is a workflow, only actions used in all jobs
in the workflow will be considered. If the target is a job, only actions used in
the job will be considered.

If the target is a repository, the checksum file must also be checked for
redundant checksums. I.e., any entry that is present in the checksum file but
was not checked against an action used by the repository must be reported as
redundant and cause the process to exit with a non-zero exit code. If the target
is not a repository then redundant checksums must be ignored.

The `-offline` flag can be used to verify strictly against the cache without
fetching any missing repositories.

## Procedures

### Collecting Actions

To determine the set of actions a target depends on, first find all `uses:`
entries in the target, both at the step- and job-level. For a repository this,
covers all workflows in the workflows directory, otherwise it covers only the
target.

For each `uses:` value it is added to the set. Any Local action or workflow is
not itself included but must be parsed for actions. If the `-no-transitive`
option is set this constitutes the set of actions. Otherwise, transitive actions
must be collected too, recursively.

To collect transitive actions for step-level `uses:` values, the action manifest
of the declared action is parsed for (transitive) actions. This concerns the
manifest at the path declared in the `uses:` value. If multiple actions from the
same repository are used, each action's manifest must be handled.

To collect transitive actions for job-level `uses:` values, the workflow of the
declared reusable workflow is parsed for (transitive) actions. This concerns
only the workflow at the path declared in the `uses:` value. If multiple
reusable workflow from the same repository are used, each workflow must be
handled.

The resulting is a collection of Action identifiers, `<owner>/<project>@<ref>`.
While GitHub Actions is case-insensitive when resolving `<owner>/<project>`,
these must NOT be normalized (as that would break on case sensitive OSes).

Docker Hub Actions, as seen in the example below, are excluded from the set of
actions the repository depends on (see [#216]).

```yaml
steps:
- uses: docker://alpine:3.8
- uses: docker://ghcr.io/OWNER/IMAGE_NAME
- uses: docker://gcr.io/cloud-builders/gradle
```

[#216]: https://github.com/chains-project/ghasum/issues/216

### Computing Checksums

To compute checksums `ghasum` will pull the repository of an action, either at
a specific ref or checking out the ref after pulling, remove the git index (i.e.
the `.git/` directory) and compute a deterministic hash over the files in the
repository, recursing through nested directories.

The hash is not configurable and the only available algorithm is SHA256.

For this process a local cache may be used. The cache will contain repositories
to avoid having to fetch them again. The cache does not contain checksums, which
will always be recomputed.

The user is able to control the usage of the cache using the flags:

- `-cache <dir>` for the location of the cache,
- `-no-cache` to disable the cache for the current execution, and
- `-no-evict` to disable eviction of cache entries.

Additionally, the `ghasum cache` command can be used to manage the cache.

### Storing Checksums

To store checksums `ghasum` uses the checksum file. This file tracks the version
of this file, checksums, and additional metadata. The version of the file and
additional metadata are all stored as _headers_. The way in which checksums are
stored depends on the version of the file, see [Sumfile Versions].

## Sumfile Versions

A checksum must always contain a header named _version_ which states the version
of the sumfile. Additional non-empty lines are considered headers. A header is
interpreted as `<name> <value>`. The first empty line marks the end of the
headers, the following line marks the start of the body of the sumfile. A
sumfile must always end with a final newline. There is no support for comments
in a sumfile.

At a high level a `ghasum` sumfile looks like:

```text
version 1
<header-2-name> <header-2-value>
...

<body>
```

Every header `<name>` and every entry in the `<body>` of the sumfile must have a
unique name/identifier. If two entries have the same identifier the sumfile must
be rejected as corrupt and the program exit with a non-zero exit code.

### Version 1

Sumfile version 1 expects at least one header, namely `version 1`. Any other
headers in the file are ignored. All checksums are stored on a separate line, no
additional empty lines are allowed.

```text
version 1
<optional headers>

<id-1> <checksum-1>
...
<id-n> <checksum-n>
```

## Definitions

- _action manifest_ is the file `action.yml`, `action.yaml`, or `Dockerfile`.
- _checksum file_ is the file `.github/workflows/gha.sum`.
- _workflows directory_ is the directory `.github/workflows`.

[collecting actions]: #collecting-actions
[computing checksums]: #computing-checksums
[storing checksums]: #storing-checksums
[sumfile versions]: #sumfile-versions
