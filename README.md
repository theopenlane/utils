[![Build status](https://badge.buildkite.com/a3a38b934ca2bb7fc771e19bc5a986a1452fa2962e4e1c63bf.svg?branch=main)](https://buildkite.com/theopenlane/utils)
[![Go Reference](https://pkg.go.dev/badge/github.com/theopenlane/utils.svg)](https://pkg.go.dev/github.com/theopenlane/utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/theopenlane/utils)](https://goreportcard.com/report/github.com/theopenlane/utils)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=theopenlane_utils&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=theopenlane_util)


# utils

Utilities for working within the openlane ecosystem, high level overview of packages:

- cache: redis client interface
- cli: cli helper utilities for printing rows, tables
- contextx: The contextx package provides helper functions for managing context values, particularly for request-scoped data. It uses generics to simplify the handling of context keys.
- dumper: The dumper package is a utility for dumping HTTP request contents, useful for debugging and logging purposes.
- envparse: struct default parsing utility
- gravatar: create sweet robot avatars based on your email
- keygen: The keygen package provides utilities for generating and validating keys. It includes a key generator that can generate keys of any length and a key validator that can validate keys of any length.
- keyring: package keyring allows for quick and easy access to the system keyring service
- marionette: The marionette package is a task manager with scheduling, backoff, and future scheduling capabilities. It is designed as a temporary solution until an external state management system is implemented.
- passwd: The passwd package provides cryptographic utilities for handling passwords.
- rout: rout is a semi-centralized method of handling and surfacing user facing errors
- slack: very minimal slack functions for sending messages
- slice: a nice big juicy slice of functions for working with slices
- sqlite: sqlite client interface
- testutils: test utilities!
- ulids: The ulids package is a lightweight wrapper around the github.com/oklog/ulid package. It provides common functionality such as checking if a ULID is null or zero and includes a process-global, cryptographically random, monotonic, and thread-safe ULID generation mechanism.

## Contributing

See the [contributing](.github/CONTRIBUTING.md) guide for more information.