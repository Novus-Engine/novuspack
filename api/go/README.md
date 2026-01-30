# NovusPack Go Implementation

This directory contains the Go implementation of the NovusPack API.

## Structure

The Go implementation follows Go's major version import path conventions:

- **v1**: Code is directly in `api/go/` (no version subdirectory)
- **v2**: Will be in `api/go/v2/` when available (future)

The v1 implementation is self-contained with:

- Go module (`go.mod`)
- Source code
- BDD test infrastructure
- Build configuration

Future major versions (v2+) will be in versioned subdirectories with their own `go.mod` files.

## Module Paths

- **v1**: `github.com/novus-engine/novuspack/api/go`
- **v2**: `github.com/novus-engine/novuspack/api/go/v2` (when available)

This follows Go's major version import path conventions, allowing multiple major versions to coexist.

## Installation

```go
import "github.com/novus-engine/novuspack/api/go"
```

## Usage Example

```go
package main

import (
    "github.com/novus-engine/novuspack/api/go"
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

Unit tests provide comprehensive coverage of the implementation.

Coverage targets are:

- **Overall coverage**: 95% target
- **Per-package minimum**: 92% coverage
- **Per-functions bare minimum**: 80%
- **Exceptions**: Some internal helper functions may have lower coverage (85%+) due to testing limitations with system-level operations (see `internal/helpers.go` for details)

To check coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

#### Unit Test Files

All `*.go` non-test files MUST have a corresponding `*_test.go` file that contains the tests for the base go file.

### BDD Tests

The BDD tests use shared feature files from `../../../features/`:

```bash
make bdd
```

BDD tests validate feature-level behavior and ensure specification compliance.

### BDD Test Infrastructure

The BDD test infrastructure is located in `_bdd/` and uses:

- [godog](https://github.com/cucumber/godog) for BDD test execution
- Shared feature files from the repository root `features/` directory
- Step definitions in `_bdd/steps/`

### Running Tests from Repository Root

From the repository root:

- `make test-go` - Run all Go tests
- `make test-go-v1` - Run Go v1 tests
- `make bdd-go` - Run all Go BDD tests
- `make bdd-go-v1` - Run Go v1 BDD tests

## Module Configuration

The `go.mod` file defines:

- Module path: `github.com/novus-engine/novuspack/api/go`
- Go version: 1.25
- Dependencies: See `go.mod` for complete list

## Dependencies and External Libraries

The Go implementation follows NovusPack's policy of minimizing external library dependencies.
All external libraries must use compatible open-source licenses throughout their entire dependency chain.

### Current Dependencies

The Go v1 implementation currently uses:

- **`github.com/samber/lo`**: Used for collection operations (filtering, mapping, searching, aggregation).
  See [samber/lo Usage Standards](../../docs/implementations/go/samber_lo_usage.md) for guidelines.
- **`github.com/cucumber/godog`**: Used for BDD test execution.
- **`github.com/goccy/go-yaml`**: Used for YAML parsing in tests.

For complete guidelines on external library usage, license compatibility requirements, verification procedures, and adding new dependencies, see the [Contributing Guide - External Library Usage](../../CONTRIBUTING.md#63-external-library-usage) section.

For detailed license information for all dependencies, see [DEPENDENCY_LICENSES.md](DEPENDENCY_LICENSES.md).

## Documentation Validation

The repository provides comprehensive documentation validation tools.
See [scripts/README.md](../../scripts/README.md) for detailed information about all available scripts and their options.

### Quick Reference

From the repository root, you can run:

- **`make docs-check [PATHS="file1.md,dir1,file2.md"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`** - Run all documentation validation checks

  - Runs markdown linting, Go code blocks validation, heading numbering, signatures index, links, requirement references, and coverage audits
  - Use `PATHS` to check specific files/directories (comma-separated)
  - Use `VERBOSE=1` for detailed output
  - Use `OUTPUT` to save reports to a file (applies to validate-links, validate-heading-numbering, validate-go-code-blocks)
  - Use `CHECK_COVERAGE=1` to check that all requirements reference tech specs (validate-links only)

- **`make validate-links [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"] [CHECK_COVERAGE=1]`** - Validate all internal markdown links and anchors

- **`make validate-go-code-blocks [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`** - Validate Go code blocks in tech specs follow conventions

- **`make validate-heading-numbering [PATHS="file1.md,dir1"] [VERBOSE=1] [OUTPUT="file.txt"]`** - Validate markdown heading numbering consistency

- **`make validate-api-go-defs-index [PATHS="file1.md,dir1"] [VERBOSE=1]`** - Validate that all Go definitions in tech specs are in the Go API definitions index

- **`make validate-req-references [PATHS="file1.feature,dir1"] [VERBOSE=1]`** - Validate REQ references in feature files

- **`make audit-feature-coverage [PATHS="file1.feature,dir1"] [VERBOSE=1]`** - Check requirements have BDD feature coverage

- **`make audit-requirements-coverage [PATHS="file1.md,dir1"] [VERBOSE=1]`** - Check tech specs are referenced by requirements

- **`make audit-coverage`** - Run both coverage audits

For complete documentation on all available options and usage examples, see [scripts/README.md](../../scripts/README.md).

## Related Documentation

- [Main README](../../README.md) - Project overview
- [Contributing Guide](../../CONTRIBUTING.md) - Contribution guidelines including external library policy
- [Technical Specifications](../../docs/tech_specs/) - API specifications
- [Versioning Policy](../../docs/specs_versioning.md) - Versioning strategy
- [Scripts Documentation](../../scripts/README.md) - Complete documentation for all validation scripts
