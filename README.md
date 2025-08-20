<!-- SPDX-License-Identifier: CC-BY-4.0 -->

# `ghasum`

Checksums for GitHub Actions.

Compute and verify checksums for all GitHub Actions in a project to guarantee
that the Actions you choose to include haven't changed since. `ghasum` gives
better integrity guarantees than pinning Actions by commit hash and is also more
user friendly as well.

## Usage

To start using `ghasum` navigate to a project that use GitHub Actions and run:

```shell
ghasum init
```

Commit the `gha.sum` file that is created so that the checksums can be verified
in the future. To verify run:

```shell
ghasum verify
```

For further help with using `ghasum` run:

```shell
ghasum help
```

## Integration

To use `ghasum` in your GitHub Actions workflows there are two options. One is
to create and use a local action that runs `ghasum` (recommended) and the other
is to run `ghasum` inline in every job.

### Local Action (Recommended)

To integrate `ghasum` as a local action into your GitHub Actions workflows you
have to create the local action and then use it in every job in every workflow.

1. <details>

   <summary>Create a local action to run <code>ghasum</code>:</summary>

   Create the file `.github/actions/ghasum/action.yml` and copy the following
   content into the file. Then fill in the `ghasum` version and checksums.

   ```yaml
   name: ghasum
   description: Verify checksums of actions

   inputs:
     checksum:
       description: The checksum of the ghasum checksums file
       required: false
       default: 0d9ca91...     # Set the 'checksums-sha512.txt' file's checksum.
     version:
       description: The version of ghasum to use
       required: false
       default: vX.Y.Z         # Set the ghasum version.

   runs:
     using: composite
     steps:
       # Unix
       - name: Initialize ghasum directory
         if: runner.os == 'macOS' || runner.os == 'Linux'
         shell: bash
         run: mkdir -p /tmp/ghasum
       - name: Download ghasum checksums
         if: runner.os == 'macOS' || runner.os == 'Linux'
         shell: bash
         working-directory: /tmp/ghasum
         env:
           CHECKSUM: ${{ inputs.checksum }}
           GH_TOKEN: ${{ github.token }}
           VERSION: ${{ inputs.version }}
         run: |
           ARTIFACT='checksums-sha512.txt'
           gh release download "$VERSION" --repo chains-project/ghasum --pattern "$ARTIFACT"
           echo "$CHECKSUM  $ARTIFACT" | shasum -a 256 -c -

       # Windows
       - name: Initialize ghasum directory
         if: runner.os == 'Windows'
         shell: pwsh
         run: mkdir C:\ghasum
       - name: Download ghasum checksums
         if: runner.os == 'Windows'
         shell: pwsh
         working-directory: C:\ghasum
         env:
           CHECKSUM: ${{ inputs.checksum }}
           GH_TOKEN: ${{ github.token }}
           VERSION: ${{ inputs.version }}
         run: |
           $ARTIFACT = "checksums-sha512.txt"
           gh release download "$env:VERSION" --repo chains-project/ghasum --pattern "$ARTIFACT"
           if ((Get-FileHash -Algorithm SHA256 "$ARTIFACT").Hash -ne "$env:CHECKSUM") {
             Write-Error 'Checksum mismatch!'
             exit 1
           } else {
             Write-Host 'Checksum match'
           }

       # macOS
       - name: Pick the ghasum CLI (amd64)
         if: runner.os == 'macOS' && runner.arch == 'X64'
         id: pick-macos-amd64
         shell: bash
         run: echo 'artifact=ghasum_darwin_amd64.tar.gz' >>"$GITHUB_OUTPUT"
       - name: Pick the ghasum CLI (arm64)
         if: runner.os == 'macOS' && runner.arch == 'ARM64'
         id: pick-macos-arm64
         shell: bash
         run: echo 'artifact=ghasum_darwin_arm64.tar.gz' >>"$GITHUB_OUTPUT"
       - name: Download the ghasum CLI
         if: runner.os == 'macOS'
         shell: bash
         working-directory: /tmp/ghasum
         env:
           ARTIFACT: ${{ steps.pick-macos-amd64.outputs.artifact || steps.pick-macos-arm64.outputs.artifact }}
           GH_TOKEN: ${{ github.token }}
           VERSION: ${{ inputs.version }}
         run: |
           gh release download "$VERSION" --repo chains-project/ghasum --pattern "$ARTIFACT"
           shasum --check --ignore-missing checksums-sha512.txt
           tar -xf "$ARTIFACT"
       - name: Verify the action checksums
         if: runner.os == 'macOS'
         shell: bash
         env:
           JOB: ${{ github.job }}
           WORKFLOW: ${{ github.workflow_ref }}
         run: |
           WORKFLOW=$(echo "$WORKFLOW" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
           /tmp/ghasum/ghasum verify -cache /Users/runner/work/_actions -no-evict -offline "$WORKFLOW:$JOB"

       # Linux
       - name: Pick the ghasum CLI (amd64)
         if: runner.os == 'Linux' && runner.arch == 'X64'
         id: pick-linux-amd64
         shell: bash
         run: echo 'artifact=ghasum_linux_amd64.tar.gz' >>"$GITHUB_OUTPUT"
       - name: Pick the ghasum CLI (arm64)
         if: runner.os == 'Linux' && runner.arch == 'ARM64'
         id: pick-linux-arm64
         shell: bash
         run: echo 'artifact=ghasum_linux_arm64.tar.gz' >>"$GITHUB_OUTPUT"
       - name: Download the ghasum CLI
         if: runner.os == 'Linux'
         shell: bash
         working-directory: /tmp/ghasum
         env:
           ARTIFACT: ${{ steps.pick-linux-amd64.outputs.artifact || steps.pick-linux-arm64.outputs.artifact }}
           GH_TOKEN: ${{ github.token }}
           VERSION: ${{ inputs.version }}
         run: |
           gh release download "$VERSION" --repo chains-project/ghasum --pattern "$ARTIFACT"
           shasum --check --ignore-missing checksums-sha512.txt
           tar -xf "$ARTIFACT"
       - name: Verify the action checksums
         if: runner.os == 'Linux'
         shell: bash
         env:
           JOB: ${{ github.job }}
           WORKFLOW: ${{ github.workflow_ref }}
         run: |
           WORKFLOW=$(echo "$WORKFLOW" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
           /tmp/ghasum/ghasum verify -cache /home/runner/work/_actions -no-evict -offline "$WORKFLOW:$JOB"

       # Windows
       - name: Pick the ghasum CLI (amd64)
         if: runner.os == 'Windows' && runner.arch == 'X64'
         id: pick-windows-amd64
         shell: pwsh
         run: |
           'artifact=ghasum_windows_amd64.zip' >>"$env:GITHUB_OUTPUT"
       - name: Pick the ghasum CLI (arm64)
         if: runner.os == 'Windows' && runner.arch == 'ARM64'
         id: pick-windows-arm64
         shell: pwsh
         run: |
           'artifact=ghasum_windows_arm64.zip' >>"$env:GITHUB_OUTPUT"
       - name: Download the ghasum CLI
         if: runner.os == 'Windows'
         shell: pwsh
         working-directory: C:\ghasum
         env:
           ARTIFACT: ${{ steps.pick-windows-amd64.outputs.artifact || steps.pick-windows-arm64.outputs.artifact }}
           GH_TOKEN: ${{ github.token }}
           VERSION: ${{ inputs.version }}
         run: |
           gh release download "$env:VERSION" --repo chains-project/ghasum --pattern "$env:ARTIFACT"
           $line = Get-Content checksums-sha512.txt | Where-Object { $_ -match "\b$env:ARTIFACT$" }
           if (-not $line) {
             Write-Error 'Checksum missing'
             exit 2
           } else {
             if ($line -match "^([a-fA-F0-9]+)  $env:ARTIFACT$") {
               $want = $matches[1]
               $got = (Get-FileHash -Path $env:ARTIFACT -Algorithm SHA512).Hash
               if ($got.ToLower() -ne $want.ToLower()) {
                 Write-Error 'Checksum mismatch'
                 exit 1
               } else {
                 Write-Host 'Checksum match'
                 Expand-Archive -Path "$env:ARTIFACT" -DestinationPath .
               }
             } else {
               Write-Error 'Checksums malformed'
               exit 2
             }
           }
       - name: Verify the action checksums
         if: runner.os == 'Windows'
         shell: pwsh
         env:
           JOB: ${{ github.job }}
           WORKFLOW: ${{ github.workflow_ref }}
         run: |
           $WorkflowParts = $env:WORKFLOW -split '@'
           $WorkflowPath = ($WorkflowParts[0] -split '/')[2..4] -join '/'
           if (Test-Path -Path 'C:\a\_actions') {
             C:\ghasum\ghasum.exe verify -cache C:\a\_actions -no-evict -offline "${WorkflowPath}:$env:JOB"
           } else {
             C:\ghasum\ghasum.exe verify -cache D:\a\_actions -no-evict -offline "${WorkflowPath}:$env:JOB"
           }

       # Cleanup
       - name: Cleanup (Unix)
         if: runner.os == 'macOS' || runner.os == 'Linux'
         shell: bash
         run: rm -rf /tmp/ghasum
       - name: Cleanup (Windows)
         if: runner.os == 'Windows'
         shell: pwsh
         run: Remove-Item -Recurse -Force -Path C:\ghasum
   ```

   </details>

2. <details>

   <summary>Use the local action in your workflows:</summary>

   ```yaml
   jobs:
     example:
       steps:
       # The repository has to be checked out before verifying checksums because
       #  it requires access to the content in .github/workflows. Because this
       #  action is run before the checksums are verified it should be pinned to
       #  a commit SHA.
       - name: Checkout repository
         uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

       # Verify the checksums with ghasum through the local action.
       - name: Verify action checksums
         uses: ./.github/actions/ghasum

       # The rest of your job ...
   ```

   </details>

### Inline

<details>

<summary>For Ubuntu runners</summary>

```yaml
job:
  runs-on: ubuntu-24.04 # Also 'ubuntu-latest'
  steps:
  # The repository has to be checked out before verifying checksums because it
  #  requires access to the content in .github/workflows. Because this action is
  #  run before the checksums are verified it should be pinned to a commit SHA.
  - name: Checkout repository
    uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

  # Verify the action checksums with ghasum.
  - name: Verify action checksums
    env:
      VERSION: vX.Y.Z                # Set the ghasum version.
      CHECKSUM: f5f2ff0...           # Set the ghasum binary checksum.
      GH_TOKEN: ${{ github.token }}  # Required for the GitHub CLI (`gh`).
      JOB: ${{ github.job }}
      WORKFLOW: ${{ github.workflow_ref }}
    run: |
      # Download the ghasum CLI
      ARTIFACT="ghasum_linux_amd64.tar.gz"
      gh release download "${VERSION}" --repo chains-project/ghasum --pattern "${ARTIFACT}"
      echo "${CHECKSUM}  ${ARTIFACT}" | shasum -a 512 -c -
      tar -xf "${ARTIFACT}"

      # Verify the action checksums
      WORKFLOW=$(echo "${WORKFLOW}" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
      ./ghasum verify -cache /home/runner/work/_actions -no-evict -offline "${WORKFLOW}:${JOB}"

  # The rest of your job ...
```

</details>

<details>

<summary>For macOS runners</summary>

For newer ARM-based runners:

```yaml
job:
  runs-on: macos-15 # Also 'macos-latest'
  steps:
  # The repository has to be checked out before verifying checksums because it
  #  requires access to the content in .github/workflows. Because this action is
  #  run before the checksums are verified it should be pinned to a commit SHA.
  - name: Checkout repository
    uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

  # Verify the action checksums with ghasum.
  - name: Verify action checksums
    env:
      VERSION: vX.Y.Z                # Set the ghasum version.
      CHECKSUM: 94a5919...           # Set the ghasum binary checksum.
      GH_TOKEN: ${{ github.token }}  # Required for the GitHub CLI (`gh`).
      JOB: ${{ github.job }}
      WORKFLOW: ${{ github.workflow_ref }}
    run: |
      # Download the ghasum CLI
      ARTIFACT="ghasum_darwin_arm64.tar.gz"
      gh release download "${VERSION}" --repo chains-project/ghasum --pattern "${ARTIFACT}"
      echo "${CHECKSUM}  ${ARTIFACT}" | shasum -a 512 -c -
      tar -xf "${ARTIFACT}"

      # Verify the action checksums
      WORKFLOW=$(echo "${WORKFLOW}" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
      ./ghasum verify -cache /Users/runner/work/_actions -no-evict -offline "${WORKFLOW}:${JOB}"

  # The rest of your job ...
```

For older Intel-based runners:

```yaml
job:
  runs-on: macos-13
  steps:
  # The repository has to be checked out before verifying checksums because it
  #  requires access to the content in .github/workflows. Because this action is
  #  run before the checksums are verified it should be pinned to a commit SHA.
  - name: Checkout repository
    uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

  # Verify the action checksums with ghasum.
  - name: Verify action checksums
    env:
      VERSION: vX.Y.Z                # Set the ghasum version.
      CHECKSUM: 3414193...           # Set the ghasum binary checksum.
      GH_TOKEN: ${{ github.token }}  # Required for the GitHub CLI (`gh`).
      JOB: ${{ github.job }}
      WORKFLOW: ${{ github.workflow_ref }}
    run: |
      # Download the ghasum CLI
      ARTIFACT="ghasum_darwin_amd64.tar.gz"
      gh release download "${VERSION}" --repo chains-project/ghasum --pattern "${ARTIFACT}"
      echo "${CHECKSUM}  ${ARTIFACT}" | shasum -a 512 -c -
      tar -xf "${ARTIFACT}"

      # Verify the action checksums
      WORKFLOW=$(echo "${WORKFLOW}" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
      ./ghasum verify -cache /Users/runner/work/_actions -no-evict -offline "${WORKFLOW}:${JOB}"

  # The rest of your job ...
```

</details>

<details>

<summary>For Windows runners</summary>

```yaml
job:
  runs-on: windows-2025 # Also 'windows-latest'
  steps:
  # The repository has to be checked out before verifying checksums because it
  #  requires access to the content in .github/workflows. Because this action is
  #  run before the checksums are verified it should be pinned to a commit SHA.
  - name: Checkout repository
    uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

  # Verify the action checksums with ghasum.
  - name: Verify action checksums
    env:
      VERSION: vX.Y.Z                # Set the ghasum version.
      CHECKSUM: e3d49db...           # Set the ghasum binary checksum.
      GH_TOKEN: ${{ github.token }}  # Required for the GitHub CLI (`gh`).
      JOB: ${{ github.job }}
      WORKFLOW: ${{ github.workflow_ref }}
    run: |
      # Download the ghasum CLI
      $ARTIFACT = "ghasum_windows_amd64.zip"
      gh release download "$env:VERSION" --repo chains-project/ghasum --pattern "$ARTIFACT"
      if ((Get-FileHash -Algorithm SHA512 "$ARTIFACT").Hash -ne $env:CHECKSUM) {
          Write-Error "Checksum mismatch!"
          exit 1
      }
      Expand-Archive -Path "$ARTIFACT" -DestinationPath .

      # Verify the action checksums
      $WorkflowParts = $env:WORKFLOW -split '@'
      $WorkflowPath = ($WorkflowParts[0] -split '/')[2..4] -join '/'
      .\ghasum.exe verify -cache C:\a\_actions -no-evict -offline "${WorkflowPath}:${env:JOB}"

  # The rest of your job ...
```

</details>

<details>

<summary>For ARM-based Ubuntu runners</summary>

```yaml
job:
  runs-on: ubuntu-24.04-arm
  steps:
  # The repository has to be checked out before verifying checksums because it
  #  requires access to the content in .github/workflows. Because this action is
  #  run before the checksums are verified it should be pinned to a commit SHA.
  - name: Checkout repository
    uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

  # Verify the action checksums with ghasum.
  - name: Verify action checksums
    env:
      VERSION: vX.Y.Z                # Set the ghasum version.
      CHECKSUM: 8a5c3d8...           # Set the ghasum binary checksum.
      GH_TOKEN: ${{ github.token }}  # Required for the GitHub CLI (`gh`).
      JOB: ${{ github.job }}
      WORKFLOW: ${{ github.workflow_ref }}
    run: |
      # Download the ghasum CLI
      ARTIFACT="ghasum_linux_arm64.tar.gz"
      gh release download "${VERSION}" --repo chains-project/ghasum --pattern "${ARTIFACT}"
      echo "${CHECKSUM}  ${ARTIFACT}" | shasum -a 512 -c -
      tar -xf "${ARTIFACT}"

      # Verify the action checksums
      WORKFLOW=$(echo "${WORKFLOW}" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
      ./ghasum verify -cache /home/runner/work/_actions -no-evict -offline "${WORKFLOW}:${JOB}"

  # The rest of your job ...
```

</details>

<details>

<summary>For ARM-based Windows runners</summary>

```yaml
job:
  runs-on: windows-11-arm
  steps:
  # The repository has to be checked out before verifying checksums because it
  #  requires access to the content in .github/workflows. Because this action is
  #  run before the checksums are verified it should be pinned to a commit SHA.
  - name: Checkout repository
    uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2

  # Verify the action checksums with ghasum.
  - name: Verify action checksums
    env:
      VERSION: vX.Y.Z                # Set the ghasum version.
      CHECKSUM: 3114a13...           # Set the ghasum binary checksum.
      GH_TOKEN: ${{ github.token }}  # Required for the GitHub CLI (`gh`).
      JOB: ${{ github.job }}
      WORKFLOW: ${{ github.workflow_ref }}
    run: |
      # Download the ghasum CLI
      $ARTIFACT = "ghasum_windows_arm64.zip"
      gh release download "$env:VERSION" --repo chains-project/ghasum --pattern "$ARTIFACT"
      if ((Get-FileHash -Algorithm SHA512 "$ARTIFACT").Hash -ne $env:CHECKSUM) {
          Write-Error "Checksum mismatch!"
          exit 1
      }
      Expand-Archive -Path "$ARTIFACT" -DestinationPath .

      # Verify the action checksums
      $WorkflowParts = $env:WORKFLOW -split '@'
      $WorkflowPath = ($WorkflowParts[0] -split '/')[2..4] -join '/'
      .\ghasum.exe verify -cache C:\a\_actions -no-evict -offline "${WorkflowPath}:${env:JOB}"

  # The rest of your job ...
```

</details>

### Using Go

In a Go-based project, you can use the Go toolchain to integrate `ghasum`. To
get started, add `ghasum` as a tool dependency:

```shell
go get -tool github.com/chains-project/ghasum/cmd/ghasum
```

Initialize `ghasum` for your project if you haven't already:

```shell
go run github.com/chains-project/ghasum/cmd/ghasum init
```

Verify the setup succeeded:

```shell
go run github.com/chains-project/ghasum/cmd/ghasum verify
```

And update your workflows to verify the checksums at runtime. You can follow
either the Local Action or Inline approach. Like with the other approaches you
must first use `actions/checkout`. Moreover, because this runs from source you
must also first use `actions/setup-go`. Both should be pinned to a commit SHA as
they're run before checksum verification.

<details>

<summary>For Ubuntu/macOS runners</summary>

```yaml
- uses: actions/checkout # @commit-sha
- uses: actions/setup-go # @commit-sha
- name: Verify action checksums
  env:
    JOB: ${{ github.job }}
    WORKFLOW: ${{ github.workflow_ref }}
  run: |
    WORKFLOW=$(echo "${WORKFLOW}" | cut -d '@' -f 1 | cut -d '/' -f 3-5)
    go run github.com/chains-project/ghasum/cmd/ghasum verify \
      -cache /home/runner/work/_actions -no-evict -offline "${WORKFLOW}:${JOB}"
```

</details>

<details>

<summary>For Windows runners</summary>

```yaml
- uses: actions/checkout # @commit-sha
- uses: actions/setup-go # @commit-sha
- name: Verify action checksums
  env:
    JOB: ${{ github.job }}
    WORKFLOW: ${{ github.workflow_ref }}
  run: |
    $WorkflowParts = $env:WORKFLOW -split '@'
    $WorkflowPath = ($WorkflowParts[0] -split '/')[2..4] -join '/'
    go run github.com/chains-project/ghasum/cmd/ghasum verify \
      -cache C:\a\_actions -no-evict -offline "${WorkflowPath}:${env:JOB}"
```

</details>

## Recommendations

When using ghasum it is recommended to pin all Actions to version tags. If
Actions are benign, these won't change over time. Major version tags or branch
refs are expected to change over time as changes are made to the Action, which
results in failing verification by ghasum. Commit SHAs do not have to be used
because the benefits they provide are covered by ghasum.

If an Action misbehaves - moving version refs after publishing - it is
recommended to use commit SHAs instead to avoid failing verification by ghasum.

```yaml
# Recommended: exact version tags
- uses: actions/checkout@v4.1.1

# Possible alternative: commit SHAs
- uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1

# Discouraged: major version refs
- uses: actions/checkout@v4

# Discouraged: branches
- uses: actions/checkout@main
```

## Benefits

- Pins transitive (composite) GitHub Actions.
- Prevents using actions that have changed since you started using them. Avoids
  the impact of supply chain attacks such as [CVE-2025-30066] (`tj-actions`).
- Prevents using [impostor commits].
- Reveals your GitHub Actions dependency hierarchy with `ghasum list`, even
  without integrating `ghasum`.
- Protects against git commit SHA hash collisions (more details below).

[impostor commits]: https://www.chainguard.dev/unchained/what-the-fork-imposter-commits-in-github-actions-and-ci-cd
[CVE-2025-30066]: https://github.com/advisories/GHSA-mrrh-fwg8-r2c3

## Limitations

- Requires manual intervention when an Action is updated.
- The hashing algorithm used for checksums is not configurable.
- `ghasum` does not (yet, [#216]) handle Docker-based [unpinnable actions].
- Checksums do not provide protection against code-based [unpinnable actions].

Some of these limitations may be addressed by Github's Immutable Actions
initiative, see [github/roadmap#592] for more information.

[#216]: https://github.com/chains-project/ghasum/issues/216
[github/roadmap#592]: https://github.com/github/roadmap/issues/592
[unpinnable actions]: https://www.paloaltonetworks.com/blog/prisma-cloud/unpinnable-actions-github-security/

## Background

The dependency ecosystem for GitHub Actions is fully reliant on git. The version
of an Action to use is specified using a git ref (branch or tag) or commit SHA.
Git refs provide no integrity guarantees. And while commit SHAs do provide some
integrity guarantees, since they're based on the older SHA1 hash the guarantees
are not optimal.

Besides being older and having better, modern algorithms available, SHA1 is
vulnerable to attacks, including [SHAttered] and [SHAmbles]. This means it is
possible for a motivated and well-funded adversary to mount an attack on the
GitHub Actions ecosystem. Note that GitHub does have [protections in place] to
detect such attacks, but from what is publicly available this is limited to the
[SHAttered] attack.

This project is a response to that theoretical attack - providing a way to get,
record, and validate checksums for GitHub Actions dependencies using a more
secure hashing algorithm. As an added benefit, it can also be used as an
alternative to in-workflow commit SHA.

[protections in place]: https://github.blog/2017-03-20-sha-1-collision-detection-on-github-com/
[shattered]: https://shattered.io/
[shambles]: https://sha-mbles.github.io/

### Git's hash function transition

The Git project has a [hash function transition] objective with the goal of
migrating from SHA-1 to SHA-256. This discussion was started around the time of
the SHAttered attack and has gradually been developed over time but is, as of
writing, still experimental. The transition would eliminate the need for this
project from a security perspective, but it could remain useful due to its other
perks.

[hash function transition]: https://git-scm.com/docs/hash-function-transition

## License

This software is available under the Apache License 2.0 license, see [LICENSE]
for the full license text. The contents of documentation are licensed under the
[CC BY 4.0] license.

[cc by 4.0]: https://creativecommons.org/licenses/by/4.0/
[LICENSE]: ./LICENSE
