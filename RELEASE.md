<!-- SPDX-License-Identifier: CC0-1.0 -->

# Release Guidelines

To release a new version of the `ghasum` project follow these steps (using
v1.2.3 as an example):

1. Make sure that your local copy of the repository is up-to-date, sync:

   ```shell
   git checkout main
   git pull origin main
   ```

   Or clone:

   ```shell
   git clone git@github.com:chains-project/ghasum.git
   ```

1. Create a [git tag] for the new version and push it:

   ```shell
   git tag v1.2.3
   git push origin v1.2.3
   ```

   > **Note**: At this point, the continuous delivery automation may pick up and
   > complete the release process. If not, or only partially, continue following
   > the remaining steps.

1. Create pre-compiled binaries - with checksums - for various targets using:

   ```shell
   go run tasks.go build-all
   ```

1. Create a [GitHub Release] for the [git tag] of the new release. The release
   title should be "Release {_version_}" (e.g. "Release v1.2.3"). The release
   text should be a description of the changes since the previous release. The
   release artifact should be the pre-compiled binaries, including checksums,
   from the previous step.

[git tag]: https://git-scm.com/book/en/v2/Git-Basics-Tagging
[github release]: https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository
