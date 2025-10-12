# NovusPack

A modern, quantum-safe package archive library that provides comprehensive alternatives to traditional packaging formats like ZIP, TAR, and 7Z.

## Overview

NovusPack is a library designed for creating, managing, and manipulating modern package archives with support for compression, encryption, digital signatures, and streaming capabilities.
It provides a comprehensive solution for secure file packaging with quantum-safe cryptography.

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
