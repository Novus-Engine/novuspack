# NovusPack Go Implementation - API Version 1

This directory contains the Go implementation of NovusPack API version 1.

## Module Path

```text
github.com/novus-engine/novuspack/api/go/v1
```

## Installation

```go
import "github.com/novus-engine/novuspack/api/go/v1"
```

## Usage Example

```go
package main

import (
    "github.com/novus-engine/novuspack/api/go/v1"
)

func main() {
    // Create a new package
    pkg := novuspack.NewPackage()
    // ... use package
}
```

## Building

From this directory:

```bash
go build ./...
```

Or from the repository root:

```bash
make test-go-v1
```

## Testing

### Unit Tests

```bash
make test
```

### BDD Tests

The BDD tests use shared feature files from `../../../features/`:

```bash
make bdd
```

## BDD Test Infrastructure

The BDD test infrastructure is located in `bdd/` and uses:

- [godog](https://github.com/cucumber/godog) for BDD test execution
- Shared feature files from the repository root `features/` directory
- Step definitions in `bdd/steps/`

## Module Configuration

The `go.mod` file defines:

- Module path: `github.com/novus-engine/novuspack/api/go/v1`
- Go version: 1.25
- Dependencies: See `go.mod` for complete list

## Related Documentation

- [Go Implementation Overview](../README.md)
- [Main README](../../../README.md)
- [Technical Specifications](../../../docs/tech_specs/)
- [Versioning Policy](../../../docs/VERSIONING.md)
