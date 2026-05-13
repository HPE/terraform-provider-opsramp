# Contributing to terraform-provider-opsramp

Thank you for your interest in contributing!

## Before You Start

- This repository has been created as part of 
[HPE Open Source community](https://www.hpe.com/us/en/open-source.html).
- HPE employees contributing as part of their role are strictly required to comply 
with HPE Open Program Office (OPO) and HPE Open Source Review Board (OSRB) 
requirements. 
- **External (non-HPE) contributors must sign the HPE Individual Contributor License 
Agreement (ICLA) before any pull request can be merged.** See [Contributor License 
Agreement](#contributor-license-agreement) below.
- **Do not sign any third-party CLA** presented to you without first having it 
reviewed by the HPE Open Source Review Board (OSRB).
Contact [osrb-admins@hpe.com](mailto:osrb-admins@hpe.com).

## Contributor License Agreement

This project uses the **HPE Individual Contributor License Agreement (ICLA)**,
an Apache-style CLA that covers copyright and patent grants necessary for an
Apache-2.0 licensed project. A CLA captures both the copyright
license and the patent grant from each contributor, and provides legal
protections for both HPE and the contributor. 

### Process for external contributors

1. Before opening your first pull request, request the HPE Individual CLA form
   by emailing [osrb-admins@hpe.com](mailto:osrb-admins@hpe.com).
2. Complete, sign, and return the form to
   [osrb-admins@hpe.com](mailto:osrb-admins@hpe.com).
3. OSRB will confirm once your name has been added to HPE's on-file list.
4. Open your pull request. Reference your CLA confirmation in the PR
   description.
5. PRs from contributors without a signed ICLA on file will not be merged.

The ICLA is executed once and covers all future contributions to this project.

## Coding Standards

- Follow existing code patterns documented in [README.md](README.md#architecture-conventions).
- All new or modified HPE-authored source files must include the SPDX header:

  ```go
  // SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
  // SPDX-License-Identifier: Apache-2.0
  ```

- Run `go build ./...` and `go vet ./...` before submitting a pull request.

## Pull Request Process

1. Fork the repository and create a feature branch from `main`.
2. Make your changes, add tests where applicable.
3. Ensure your ICLA is on file with OSRB (external contributors) or that you
   are an HPE employee contributing under your employment agreement.
4. Open a pull request against `main` and describe what was changed and why.
5. Address any review feedback.

## License

By contributing, you agree that your contributions will be licensed under the
[Apache License, Version 2.0](LICENSE).
