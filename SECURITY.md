<!-- SPDX-License-Identifier: CC0-1.0 -->

# Security Policy

The maintainers of the `ghasum` project take security issues seriously. We
appreciate your efforts to responsibly disclose your findings. Due to the
non-funded and open-source nature of the project, we take a best-efforts
approach when it comes to engaging with security reports.

This document should be considered expired after 2026-06-01. If you are reading
this after that date, try to find an up-to-date version in the official source
repository.

## Supported Versions

Only the latest release of the project is supported with security updates.

### Threat Model

The program considers the Go runtime and CLI arguments to be trusted. The
content of the target repository are assumed to be correct but are otherwise
untrusted. All other input and external content is considered untrusted. Any
violation of availability, confidentiality, or integrity is considered an issue.

The project considers project maintainers and the GitHub infrastructure to be
trusted. Any action performed by any other GitHub user against the repository is
considered untrusted.

## Reporting a Vulnerability

To report a security issue in the latest release or development head, either:

- [Report it through GitHub][new github advisory], or
- Send an email to [security@ericcornelissen.dev] with the terms "SECURITY" and
  "ghasum" in the subject line.

Please do not open a regular issue or Pull Request in the public repository.

To report a security issue in an older version - i.e. the latest release isn't
affected - please report it publicly. For example, as a regular issue in the
public repository. If in doubt, report the issue privately.

[new github advisory]: https://github.com/chains-project/ghasum/security/advisories/new
[security@ericcornelissen.dev]: mailto:security@ericcornelissen.dev?subject=SECURITY%20%28ghasum%29

### What to Include in a Report

Try to include as many of the following items as possible in a security report:

- An explanation of the problem
- A proof of concept exploit
- A suggested severity
- Relevant [CWE] identifiers
- The latest affected version
- The earliest affected version
- A suggested patch
- An automated regression test

[cwe]: https://cwe.mitre.org/

## Advisories

| ID               | Date       | Affected version(s) | Patched version(s) |
| :--------------- | :--------- | :------------------ | :----------------- |
| -                | -          | -                   | -                  |

## Acknowledgments

We would like to publicly thank the following reporters:

- _None yet_
