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

### `ghasum update`

If the checksum file does not exist the process shall exit immediately with an
error.

If the checksum file exists the process shall obtain a lock on it, if this is
not possible to process shall exit immediately (it means the file may be edited
by another process leading to an inconsistent state).

If the file lock is obtained, the process shall first read it and parse it
completely to extract the sumfile version. If this fails the process shall exit
immediately unless the `-force` flag is used (see details below). Else it shall
compute checksums (see [Computing Checksums]) for all new actions used in the
repository (see [Collecting Actions]) using the same hashing algorithm as was
used for the existing checksums. New actions also include new versions of a
previously used actions. Additionally, it should remove any entry which is no
longer in use. No existing checksum for a used action shall be updated. It shall
then store them in a sumfile (see [Storing Checksums]) using the same sumfile
version as before and releases the lock. In short, updating will only add new
and remove old checksums from an existing sumfile.

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
a non-zero exit code, for usability all values should be compared (and all
mismatches reported) before exiting.

The "target" can be one of a: a repository, a workflow, or a job. If the target
is a repository, all actions used in all jobs in all workflows in the repository
will be considered. If the target is a workflow, only actions used in all jobs
in the workflow will be considered. If the target is a job, only actions used in
the job will be considered.

Redundant checksums are ignored by this process.

## Procedures

### Collecting Actions

To determine the set of actions a target depends on, first find all `uses:`
entries in the target. For a repository this covers all workflows in the
workflows directory, otherwise it covers only the target.

For each `uses:` value, excluding the list below, it is added to the set. If the
`-no-transitive` option is NOT set the repository declared by the `uses:` value
is fetched. The action manifest at the path specified in the `uses:` value is
parsed for additional `uses:` values. For each of these transitive `uses:`
values, this process is repeated.

The following `uses:` values are to be excluded from the set of actions a
repository depends on.

- Actions in the same repository as the workflow ("local actions"). Example:

  ```yaml
  steps:
  - uses: ./.github/actions/hello-world-action
  ```

- Docker Hub Actions ([#216]). Examples:

  ```yaml
  steps:
  - uses: docker://alpine:3.8
  - uses: docker://ghcr.io/OWNER/IMAGE_NAME
  - uses: docker://gcr.io/cloud-builders/gradle
  ```

- Reusable workflows ([#215]). Examples:

  ```yaml
  jobs:
    call-workflow-1-in-local-repo:
      uses: octo-org/this-repo/.github/workflows/workflow-1.yml@172239021f7ba04fe7327647b213799853a9eb89
    call-workflow-2-in-local-repo:
      uses: ./.github/workflows/workflow-2.yml
    call-workflow-in-another-repo:
      uses: octo-org/another-repo/.github/workflows/workflow.yml@v1
  ```

[#215]: https://github.com/chains-project/ghasum/issues/215
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

The user is able to control the usage of the cache using the `-cache <dir>` and
`-no-cache` flags. Additionally, the `ghasum cache` command can be used to
manage the cache.

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

- _action manifest_ is the file `action.yml` or `action.yaml` (mutually
  exclusive).
- _checksum file_ is the file `.github/workflows/gha.sum`.
- _workflows directory_ is the directory `.github/workflows`.

[collecting actions]: #collecting-actions
[computing checksums]: #computing-checksums
[storing checksums]: #storing-checksums
[sumfile versions]: #sumfile-versions
