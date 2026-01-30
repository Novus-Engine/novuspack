# NovusPack Technical Specifications - System Overview

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Core Components](#1-core-components)
- [2. File Formats](#2-file-formats)
- [3. Multiple Signatures Support](#3-multiple-signatures-support)
- [4. Use Cases and Applications](#4-use-cases-and-applications)

---

## 0. Overview

NovusPack is a multi-language library that serves as the core component for creating, managing, and manipulating modern package archives with support for compression, encryption, digital signatures, and streaming capabilities.
It provides a comprehensive alternative to traditional packaging formats like ZIP, TAR, and 7Z.

The library is implemented in multiple languages (Go, Rust, Zig, and future languages) with shared specifications ensuring consistency across all implementations.

> Note: This overview is superseded by the unified index in `_main.md`. For the latest entry point to the specifications, see: [_main.md](_main.md).

### 0.1 Cross-References

See the consolidated navigation in [_main.md](_main.md).

## 1. Core Components

- **Package Management:** Creation, reading, writing, and manipulation of unified .nvpk files (supports both encrypted and unencrypted files within the same package)
- **FileEntry System:** Individual file metadata and content management with per-file encryption selection
- **Compression Engine:** Multiple compression algorithms (Zstd, LZ4, LZMA) for optimal file size reduction
- **Quantum-Safe Encryption System:** ML-KEM (CRYSTALS-Kyber) for key exchange and file encryption
- **Digital Signature System:** Multiple signature support with ML-DSA (CRYSTALS-Dilithium), SLH-DSA (SPHINCS+), PGP, and X.509 for package integrity verification
- **Streaming Interface:** Memory-efficient file streaming for large files with encryption support
- **Buffer Management:** Intelligent buffer pooling and memory management optimized for encrypted content

## 2. File Formats

- **.nvpk:** Unified Novus Package format (supports both encrypted and unencrypted files within the same package)
  - **Encryption:** Quantum-safe ML-KEM (CRYSTALS-Kyber) with AES-256-GCM for compatibility
  - **Signatures:** Multiple signature support (ML-DSA, SLH-DSA, PGP, X.509) for package integrity verification
  - **Compression:** Multiple algorithms (Zstd, LZ4, LZMA) with per-file selection
  - **Key Sizes:** ML-KEM keys range from 800-1,568 bytes depending on security level

## 3. Multiple Signatures Support

NovusPack supports multiple digital signatures per package, bringing it in line with industry standards while maintaining unique quantum-safe signature advantages:

- **Multiple Signature Types**: Support for ML-DSA, SLH-DSA, PGP, and X.509 signatures in the same package
- **Incremental Signatures**: Signatures are appended incrementally; no separate signature index is used
- **Enhanced Security**: Quantum-safe signatures alongside traditional signatures for maximum security
- **Industry Compliance**: Brings NovusPack in line with PGP, X.509, Windows Authenticode, and macOS Code Signing
- **Backward Compatibility**: Future format versions will maintain backward compatibility with v1 packages

For detailed implementation information, see:

- [File Format](package_file_format.md) - Header structure and signature index specification
- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation

## 4. Use Cases and Applications

NovusPack is designed as a general-purpose archive format suitable for a wide range of applications:

- **Software Distribution**: Secure distribution of applications, libraries, and updates
- **Data Archival**: Long-term storage of sensitive data with quantum-safe encryption
- **Content Delivery**: Efficient delivery of large content packages with integrity verification
- **Backup Systems**: Secure backup solutions with multiple signature verification
- **Enterprise Applications**: Corporate data packaging with X.509 certificate integration
- **Open Source Projects**: Community-driven content distribution with PGP signatures
- **Media Distribution**: Large media file packages with selective encryption
- **Document Management**: Secure document archives with access control and audit trails

The format's comprehensive feature set makes it suitable for any application requiring secure, efficient, and verifiable file packaging.
