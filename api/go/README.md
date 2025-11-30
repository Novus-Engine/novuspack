# NovusPack Go Implementation

This directory contains the Go implementation of the NovusPack API.

## Structure

The Go implementation uses versioned directories for each major API version:

- `v1/` - API version 1 implementation
- `v2/` - API version 2 implementation (future)

Each version directory is self-contained with its own:

- Go module (`go.mod`)
- Source code
- BDD test infrastructure
- Build configuration

## Module Paths

- **v1**: `github.com/novus-engine/novuspack/api/go/v1`
- **v2**: `github.com/novus-engine/novuspack/api/go/v2` (when available)

This follows Go's major version import path conventions, allowing multiple major versions to coexist.

## Building and Testing

Each version directory has its own `Makefile` with build and test targets.

From the repository root:

- `make test-go` - Run all Go tests
- `make test-go-v1` - Run Go v1 tests
- `make bdd-go` - Run all Go BDD tests
- `make bdd-go-v1` - Run Go v1 BDD tests

From a version directory (e.g., `api/go/v1/`):

- `make test` - Run unit tests
- `make bdd` - Run BDD tests

## Documentation

- [v1 Documentation](v1/README.md) - API version 1 specific documentation

## Related Documentation

- [Main README](../../README.md) - Project overview
- [Technical Specifications](../../docs/tech_specs/) - API specifications
- [Versioning Policy](../../docs/VERSIONING.md) - Versioning strategy
