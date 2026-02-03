# NovusPack

[![Docs Check][badge-docs-check]][workflow-docs-check]
[![Go CI][badge-go-ci]][workflow-go-ci]
[![Go BDD][badge-go-bdd]][workflow-go-bdd]
[![Python Lint][badge-python-lint]][workflow-python-lint]
[![License][badge-license]][license-file]

## Overview

NovusPack is a multi-language library and format (`.nvpk`) for building secure, index-driven package archives with per-file compression, encryption, signatures, and streaming access.
The current implementation focus is Go (v1), with shared specifications and feature files defining consistent behavior across future languages.

The repository supports multiple language implementations (Go, Rust, Zig, and future languages) while maintaining shared specifications, feature files, and documentation.

## Quick Links

- 📚 **Technical specifications**: [`docs/tech_specs/_main.md`](docs/tech_specs/_main.md)
- 🧾 **Requirements**: [`docs/requirements/README.md`](docs/requirements/README.md)
- 🧪 **Shared feature files (BDD)**: [`features/`](features/)
- 🧩 **Go implementation (v1)**: [`api/go/`](api/go/)
- 🧰 **Python validation tooling**: [`scripts/`](scripts/)
- 🤝 **Contributing guide**: [`CONTRIBUTING.md`](CONTRIBUTING.md)

## Language and Tooling Status

### Go API (implementation)

- **Code**: [`api/go/`](api/go/)
- **CI**: [Go CI][workflow-go-ci], [Go BDD][workflow-go-bdd]
- **Local**: See [Testing Standards - Running Tests](#running-tests).

### Python (validation tooling)

- **Code**: [`scripts/`](scripts/)
- **CI**: [Python Lint][workflow-python-lint]
- **Lint**: `make python-lint`

### Documentation (Markdown, specs, requirements)

- **Docs and specs**: [`docs/`](docs/)
- **CI**: [Documentation Checks][workflow-docs-check]
- **Checks**: `make docs-check`

## Why NovusPack Exists

NovusPack was created to support packages that are metadata-heavy and fast to access without extracting everything up front.
Many archive formats are primarily optimized for maximum compression and bulk extraction workflows.
That is a poor fit when you want to distribute a package that mostly contains raw or lightly compressed data, while still selectively compressing or encrypting specific files.

NovusPack focuses on a format and API that support:

- Fast, directory-like access to individual files at runtime.
- A flat, index-driven internal structure to keep lookups cheap.
- Per-file compression and per-file encryption, so only the files that need protection or size reduction pay the cost.
- Rich metadata so packages can carry descriptive and operational information without relying on external sidecar files.

This came out of thinking about game content distribution for indie developers.
The goal is to make it straightforward to ship many raw text or HTML-like assets as-is, while selectively compressing and encrypting sensitive binaries or content.
In practice, this lets an application treat a Novus package almost like a directory, and only extract or stream the specific files it needs as it runs.

## Key Features

- **Quantum-Safe Encryption**: ML-KEM (CRYSTALS-Kyber) for key exchange and file encryption.
- **Multiple Digital Signatures**: Support for ML-DSA, SLH-DSA, PGP, and X.509 signatures (signature validation is deferred to v2).
- **Advanced Compression**: Multiple algorithms (Zstd, LZ4, LZMA) with per-file selection.
- **Full Package Compression**: Optional package-level compression with separate metadata and data compression for optimal storage efficiency.
- **Streaming Interface**: Memory-efficient file streaming for large files.
- **Unified Format**: Single `.nvpk` format supporting both encrypted and unencrypted files.
- **Buffer Management**: Intelligent buffer pooling optimized for encrypted content.
- **Symlink Support**: Native symlink support within packages for flexible file organization.
- **Multi-Layer Deduplication**: Efficient content deduplication with automatic symlink creation for duplicates.

## Symlinks and Deduplication

NovusPack provides powerful capabilities for managing file relationships and optimizing storage through symlinks and deduplication.

### Symlink Support

NovusPack supports symbolic links within packages, allowing you to:

- **Add symlinks directly** to packages, preserving filesystem structure.
- **Follow symlinks** when adding files (default behavior) or preserve them as symlinks.
- **Convert duplicate paths to symlinks** automatically during deduplication.
- **Validate symlink targets** to ensure they exist and remain within package boundaries.

This is particularly useful for game content where you might want multiple paths pointing to the same asset, or when preserving complex directory structures with symlinks.

**Note on Windows**: Creating symlinks on Windows requires special privileges (typically administrator rights or Developer Mode enabled).
By default, NovusPack extracts symlinks as regular file copies on Windows to avoid privilege requirements.
You can explicitly enable symlink preservation on Windows if your application has the necessary privileges.

### Multi-Layer Deduplication

NovusPack uses a three-layer deduplication system for optimal performance:

- **Layer 1 (Size Check)**: Instant elimination of files with different sizes.
- **Layer 2 (CRC32 Check)**: Fast comparison using existing checksums.
- **Layer 3 (SHA256 Check)**: Cryptographic hash-on-demand for collision resistance.

Deduplication works at multiple processing levels (raw content, processed content after compression, and final stored content), allowing you to eliminate duplicates at the most appropriate stage.

### Symlinks and Deduplication Integration

When deduplication detects duplicate content, you can choose how to handle it:

- **Hard Link Behavior** (default): Add the duplicate path to the existing FileEntry, sharing the same content.
- **Symlink Creation**: Automatically create symlinks pointing to the primary path when duplicates are found.
- **Automatic Conversion**: Enable `AutoConvertToSymlinks` to automatically convert duplicate paths to symlinks during deduplication.

This integration means you can store multiple references to the same content efficiently, with symlinks providing a clear indication of the relationship while maintaining the benefits of deduplication.

## Full Package Compression

NovusPack supports optional full package compression that compresses `FileEntry` metadata (LZ4), file data (Zstd, LZ4, or LZMA), and the file index (LZ4).
It keeps the header, metadata index (an uncompressed index that points to compressed blocks), package comment, and signatures uncompressed for direct access.

This enables selective decompression of metadata and individual files without requiring full package decompression.

For detailed information, see the [Package Compression API](docs/tech_specs/api_package_compression.md) and [Package File Format - Package Compression](docs/tech_specs/package_file_format.md#3-package-compression) specifications.

## File Format

NovusPack uses the `.nvpk` (Novus Package) format, which provides:

- **Encryption**: Quantum-safe ML-KEM with AES-256-GCM compatibility.
- **Signatures**: Multiple signature support for package integrity verification.
- **Compression**: Per-file compression algorithm selection.
- **Key Sizes**: ML-KEM keys ranging from 800 to 1,568 bytes based on security level.

## Spec Highlights

The repository contains detailed, language-agnostic technical specifications.
The root `README.md` is intentionally high-level.
For deeper details, see these canonical documents:

- 🏷️ **Metadata system**: [`docs/tech_specs/metadata.md`](docs/tech_specs/metadata.md) and [`docs/tech_specs/api_metadata.md`](docs/tech_specs/api_metadata.md)
- 🧬 **File type system**: [`docs/tech_specs/file_type_system.md`](docs/tech_specs/file_type_system.md)
- 🧱 **Transformation pipelines**: [`docs/tech_specs/api_file_mgmt_transform_pipelines.md`](docs/tech_specs/api_file_mgmt_transform_pipelines.md)
- 🛡️ **Validation and safe extraction defaults**: [`docs/tech_specs/file_validation.md`](docs/tech_specs/file_validation.md) and [`docs/tech_specs/api_file_mgmt_extraction.md`](docs/tech_specs/api_file_mgmt_extraction.md)
- ✍️ **Write strategies (SafeWrite vs FastWrite)**: [`docs/tech_specs/api_writing.md`](docs/tech_specs/api_writing.md)

## Use Cases

- Software distribution and updates
- Secure data archival with long-term storage
- Content delivery with integrity verification
- Enterprise backup solutions
- Open source project distribution
- Media file packaging with selective encryption
- Document management with access control

## Repository Structure

```text
novuspack/
├── api/                      # Language-specific implementations
│   └── go/                   # Go implementation
│       └── v1/               # API version 1
│           ├── bdd/          # BDD test infrastructure
│           ├── go.mod        # Go module
│           └── README.md     # Implementation-specific docs
├── features/                 # Shared Gherkin feature files (all implementations)
├── docs/                     # Shared documentation and specifications
│   ├── tech_specs/           # API specifications (language-agnostic)
│   └── requirements/         # Requirements documentation
└── README.md                 # This file
```

## Architecture

### Design Principles

- Shared resources at root.
  - Feature files, documentation, and specifications are shared across all implementations to ensure consistency.
- Language-specific code in `api/`.
  - Each language implementation is self-contained in versioned directories (for example, `api/go/`, and `api/rust/v1/`).
- Feature parity.
  - All implementations at the same version number must have identical feature sets and pass the same tests.
- Independent development.
  - Each language can be developed, tested, and released independently.

### Versioning

NovusPack uses a two-tier versioning system:

- **API Version Tags**: Unified tags (e.g., `v1.0.2`) indicate all implementations have the same feature set.
- **Language-Specific Tags**: Implementation-specific tags (e.g., `go/v1.0.2`) when code versions differ.

For complete versioning policy, see [Versioning Documentation](docs/specs_versioning.md).

## Documentation

See [Quick Links](#quick-links) and [Spec Highlights](#spec-highlights) for the canonical documentation entry points.

## AI-Assisted Development

This repository is designed to facilitate AI-assisted development with comprehensive tooling that provides clear, actionable feedback to help AI agents be self-correcting.

The repository provides a comprehensive validation and feedback system through Make targets and Python validation scripts that clearly identify issues, call out specification violations, and provide actionable feedback to enable iterative improvement.

This does however put particular emphasis on making sure the repository technical specifications, feature files, and requirements references are correct and up-to-date; a spec change can have cascading effects that requires everything to be in-sync for all validations to pass.

For detailed information on how the Make targets and Python validation scripts are set up to help AI agents be self-correcting, see the [AI-Assisted Development](CONTRIBUTING.md#3-ai-assisted-development) section in the Contributing Guide.

## Testing Standards

### Coverage Requirements

- **Overall coverage target**: 95%.
- **Per-package minimum**: 92% coverage.
- **Exceptions**: Some internal helper functions may have lower coverage (85%+) due to testing limitations with system-level operations that are difficult to test without mocking.

### Test Types

- **Unit Tests**: Comprehensive component-level testing with high coverage.
- **BDD Tests**: Feature-level validation using shared Gherkin feature files.

### Running Tests

From the repository root:

- `make test` - Run unit tests (currently delegates to `api/go/`)
- `make bdd` - Run Go BDD tests (writes output into `tmp/`)
- `make bdd-ci` - Run Go BDD tests in CI mode (tag-filtered)
- `make coverage` - Generate `coverage.out`
- `make coverage-report` - Show coverage in the terminal
- `make coverage-html` - Generate an HTML coverage report
- `make ci` - Run the full local CI suite (docs + Python lint + Go CI)

For domain-scoped BDD runs, use `make -C api/go bdd-domain BDD_DOMAIN='@domain:core'`.

## Project Status

NovusPack is in active development.
The current focus is the Go v1 API in [`api/go/`](api/go/) and ongoing refinement of the language-agnostic technical specifications in [`docs/tech_specs/`](docs/tech_specs/).

This repository currently provides:

- Comprehensive technical specifications and requirements that are still being refined.
- Shared BDD feature files in [`features/`](features/) that define expected behavior across implementations.
- A Go implementation scaffold and growing implementation in [`api/go/`](api/go/), with unit tests and BDD infrastructure.
- Validation tooling and CI workflows that keep documentation, specs, and code references consistent.

Active refinement areas include:

- Closing gaps between specs and the Go implementation as the v1 API is implemented.
- Tightening and clarifying technical specs as edge cases are discovered during implementation.
- Expanding BDD and unit test coverage to lock in behavior as the API stabilizes.

If you want stability today, treat the Go API as pre-1.0 and expect breaking changes until the v1 surface is declared stable.

## Contributing

This project is actively implementing the Go v1 API and refining the technical specifications.
For information on how to contribute, see the [Contributing Guide](CONTRIBUTING.md).

### External Library Policy

NovusPack follows a policy of minimizing external library dependencies across all language implementations.
When external libraries are necessary, they must use compatible open-source licenses throughout their entire dependency chain.

For complete guidelines on external library usage, license compatibility requirements, and verification procedures, see the [Contributing Guide - External Library Usage](CONTRIBUTING.md#63-external-library-usage) section.

## Why Go for the v1 Implementation

Go was selected as the initial (v1) implementation language because it provides a good balance of developer velocity, performance, and safety.

- **Easy to write and maintain**: Simple language features and excellent tooling make it straightforward to ship and iterate on initial prototyping.
- **Fast**: Strong baseline performance with low operational overhead for typical archive workloads.
- **Type safe**: Static typing helps keep the API surface predictable and reduces entire classes of bugs.
- **Memory safe** (mostly): Automatic memory management avoids many memory corruption issues common in lower-level languages.

For Go development details, see [`api/go/README.md`](api/go/README.md) and [`api/go/_bdd/README.md`](api/go/_bdd/README.md).

## License

NovusPack is licensed under the **Apache License 2.0**.

This license allows:

- ✅ **Commercial use** of the library in production applications (no restrictions).
- ✅ **Modification** and creation of derivative works.
- ✅ **Distribution** in source or binary form.
- ✅ **Contributions** from anyone without restriction.
- ✅ **Patent grant** - explicit protection from patent litigation by contributors.

This license provides:

- 📜 **Explicit patent grant** - important for cryptography projects.
- 🛡️ **Defensive termination** - deters frivolous patent litigation.
- 🌍 **Industry standard** - widely recognized and trusted.

For the complete license text, see [LICENSE](LICENSE.txt).

## Security

NovusPack implements quantum-safe cryptography to protect against future quantum computing threats while maintaining compatibility with traditional cryptographic systems.
For detailed security information, see the [Security Architecture](docs/tech_specs/security.md) documentation.

[badge-docs-check]: https://github.com/novus-engine/novuspack/actions/workflows/docs-check.yml/badge.svg?branch=main
[badge-go-ci]: https://github.com/novus-engine/novuspack/actions/workflows/go-ci.yml/badge.svg?branch=main
[badge-go-bdd]: https://github.com/novus-engine/novuspack/actions/workflows/go-bdd.yml/badge.svg?branch=main
[badge-python-lint]: https://github.com/novus-engine/novuspack/actions/workflows/python-lint.yml/badge.svg?branch=main
[badge-license]: https://img.shields.io/badge/license-Apache%202.0-blue

[workflow-docs-check]: https://github.com/novus-engine/novuspack/actions/workflows/docs-check.yml
[workflow-go-ci]: https://github.com/novus-engine/novuspack/actions/workflows/go-ci.yml
[workflow-go-bdd]: https://github.com/novus-engine/novuspack/actions/workflows/go-bdd.yml
[workflow-python-lint]: https://github.com/novus-engine/novuspack/actions/workflows/python-lint.yml
[license-file]: LICENSE.txt
