# NovusPack

A modern, quantum-safe package archive library that provides comprehensive alternatives to traditional packaging formats like ZIP, TAR, and 7Z.

## Overview

NovusPack is a multi-language library designed for creating, managing, and manipulating modern package archives with support for compression, encryption, digital signatures, and streaming capabilities.
It provides a comprehensive solution for secure file packaging with quantum-safe cryptography.

The repository supports multiple language implementations (Go, Rust, Zig, and future languages) while maintaining shared specifications, feature files, and documentation.

## Key Features

- **Quantum-Safe Encryption**: ML-KEM (CRYSTALS-Kyber) for key exchange and file encryption
- **Multiple Digital Signatures**: Support for ML-DSA, SLH-DSA, PGP, and X.509 signatures
- **Advanced Compression**: Multiple algorithms (Zstd, LZ4, LZMA) with per-file selection
- **Streaming Interface**: Memory-efficient file streaming for large files
- **Unified Format**: Single .npk format supporting both encrypted and unencrypted files
- **Buffer Management**: Intelligent buffer pooling optimized for encrypted content

## File Format

NovusPack uses the `.npk` (Novus Package) format, which provides:

- **Encryption**: Quantum-safe ML-KEM with AES-256-GCM compatibility
- **Signatures**: Multiple signature support for package integrity verification
- **Compression**: Per-file compression algorithm selection
- **Key Sizes**: ML-KEM keys ranging from 800-1,568 bytes based on security level

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

- **Shared Resources at Root**: Feature files, documentation, and specifications are shared across all implementations to ensure consistency.
- **Language-Specific Code in `api/`**: Each language implementation is self-contained in versioned directories (e.g., `api/go/v1/`, `api/rust/v1/`).
- **Feature Parity**: All implementations at the same version number must have identical feature sets and pass the same tests.
- **Independent Development**: Each language can be developed, tested, and released independently.

### Versioning

NovusPack uses a two-tier versioning system:

- **API Version Tags**: Unified tags (e.g., `v1.0.2`) indicate all implementations have the same feature set.
- **Language-Specific Tags**: Implementation-specific tags (e.g., `go/v1.0.2`) when code versions differ.

For complete versioning policy, see [Versioning Documentation](docs/VERSIONING.md).

## Documentation

Comprehensive technical specifications are available in the `docs/tech_specs/` directory.

- **[Technical Specifications Main Index](docs/tech_specs/_main.md)**
- [System Overview](docs/tech_specs/_overview.md)
- [Package File Format](docs/tech_specs/package_file_format.md)
- [Security Architecture](docs/tech_specs/security.md)

## Quick Start

The repository structure and initial Go BDD test framework are now established.

To run BDD tests:

```bash
cd api/go/v1
make bdd
```

## Project Status

This repository contains:

- **Complete initial specifications** for NovusPack
- **Go BDD test framework** with 898 feature files covering all domains
- **Repository structure** established for multi-language support
- **Requirements documentation** organized by domain

The Go implementation is in progress, with the BDD framework ready for functional implementation.

## Contributing

This project has completed the initial specification phase and is moving into implementation.

For Go development:

- Feature files: [`features/`](features/) at repository root
- Go implementation: [`api/go/v1/`](api/go/v1/)
- BDD tests: [`api/go/v1/bdd/`](api/go/v1/bdd/)
- Requirements: [`docs/requirements/`](docs/requirements/)

## License

NovusPack is licensed under the **Business Source License 1.1 (BSL 1.1)**.

This license allows:

- ✅ **Commercial use** of the library in production applications (no restrictions)
- ✅ **Contributions** from anyone without restriction
- ✅ **Non-commercial forks** and modifications
- ✅ **Internal use** and modifications within organizations

This license restricts:

- ❌ **Commercial forks** that create competing products or services substantially similar to NovusPack

**Change Date:** January 1, 2030

On the Change Date, the license will automatically convert to the **Apache License 2.0**, making the code fully open source.

For the complete license text, see [LICENSE](LICENSE).

## Security

NovusPack implements quantum-safe cryptography to protect against future quantum computing threats while maintaining compatibility with traditional cryptographic systems.
For detailed security information, see the [Security Architecture](docs/tech_specs/security.md) documentation.
