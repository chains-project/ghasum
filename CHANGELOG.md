<!-- SPDX-License-Identifier: CC0-1.0 -->

# Changelog

All notable changes to `ghasum` will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to [Semantic
Versioning].

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Enhancements

- Add the `ghasum list` subcommand to get a nested list of GitHub Actions
  dependencies for the target.
- Include archived status in the `ghasum list` output.
- Request a bug report when a panic occurs.

### Security

- Upgrade Go to `v1.24.4`.

## [v0.5.2] - 2025-05-29

### Bugs

- Fix verifying a non-job target identified by the absolute path on Windows.
- Fix various cases where files or directories could not found on Windows.

## [v0.5.1] - 2025-05-25

### Enhancements

- Report redundant checksums on verification of an entire repository.

### Bugs

- Fix errors for actions with a `Dockerfile` manifest.
- Fix unexpected error due to Windows-style newlines in the sumfile.

## [v0.5.0] - 2025-05-03

### Enhancements

- Correct typo in the `ghasum help verify` output.
- Correct typo in the `ghasum verify` output.
- Enable cache eviction on `ghasum init`.
- Ensure `ghasum verify` outcome is linked to `gha.sum` content.
- Include reusable workflows in `gha.sum`, including transitive actions used in
  reusable workflows.

## [v0.4.0] - 2025-04-27

### Enhancements

- Include transitive actions in `gha.sum` (opt-out available).
- Improve performance of cloning repositories at a commit.

### Bugs

- Fix errors for `uses:` values with local actions.
- Fix errors for `uses:` values with Docker Hub actions.

### Security

- Upgrade Go to `v1.24.2`.

### Miscellaneous

- Improve reproducibility by using `-trimpath` for release builds.

## [v0.3.0] - 2025-01-25

### Enhancements

- Improve behavior for sumfiles with duplicate entries.
- Add `-offline` verification support
- Add cache eviction support.
- Make `ghasum update` preserve existing checksums by default.

### Bugs

- Fix behavior for sumfiles with duplicate entries.

### Security

- Upgrade Go to `v1.23.5`.

### Miscellaneous

- Improve reproducibility by using `-trimpath`.

## [v0.2.0] - 2024-03-21

### Enhancements

- Support validating a single workflow.
- Support validating a single job in a workflow.
- Make `ghasum update` error if the `gha.sum` file is corrupted.

### Bugs

- Unlock `gha.sum` if an error occurs during updating.
- Correct parsing uses values with multiple `@` characters.

### Security

- Upgrade Go to `v1.22.1`.

## [v0.1.0] - 2024-02-17

Initial release.
