# NovusPack Technical Specifications - Testing Requirements

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Dual Encryption Testing Requirements](#1-dual-encryption-testing-requirements)
  - [1.1 ML-KEM Encryption Testing](#11-ml-kem-encryption-testing)
  - [1.2 AES-256-GCM Encryption Testing](#12-aes-256-gcm-encryption-testing)
  - [1.3 Dual Encryption Integration Testing](#13-dual-encryption-integration-testing)
- [2. File Validation Testing Requirements](#2-file-validation-testing-requirements)
  - [2.1 Empty File Testing](#21-empty-file-testing)
  - [2.2 Path Normalization Testing](#22-path-normalization-testing)
  - [2.3 Compression Error Handling Testing](#23-compression-error-handling-testing)
  - [2.4 Hash-based Deduplication Testing](#24-hash-based-deduplication-testing)

---

## 0. Overview

This document defines the comprehensive testing requirements for the NovusPack system, including encryption testing, file validation testing, and performance validation.

### 0.1 Cross-References

- [Main Index](_main.md) - Central navigation for all NovusPack specifications
- [System Overview](_overview.md) - System overview and core components
- [Package File Format](package_file_format.md) - .nvpk format and FileEntry structure
- [File Types System](file_type_system.md) - Comprehensive file type system
- [Metadata System](metadata.md) - Package metadata and tags system
- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Security and Encryption](security.md) - Comprehensive security architecture, encryption implementation, and digital signature requirements
- [File Validation](file_validation.md) - File validation and transparency requirements

---

## 1. Dual Encryption Testing Requirements

This section defines testing requirements for dual encryption functionality, covering both ML-KEM and AES-256-GCM encryption methods.

### 1.1 ML-KEM Encryption Testing

- **Key generation:** Test ML-KEM key generation for all security levels
- **Encryption/decryption:** Test file encryption and decryption with ML-KEM
- **Performance testing:** Benchmark ML-KEM operations for various file sizes
- **Security validation:** Verify ML-KEM implementation meets NIST standards
- **Cross-platform testing:** Test ML-KEM on different operating systems

### 1.2 AES-256-GCM Encryption Testing

- **Key generation:** Test AES key generation and management
- **Encryption/decryption:** Test file encryption and decryption with AES-256-GCM
- **Performance testing:** Benchmark AES operations for various file sizes
- **Security validation:** Verify AES-256-GCM implementation meets industry standards
- **Cross-platform testing:** Test AES on different operating systems

### 1.3 Dual Encryption Integration Testing

- **Mixed encryption packages:** Test packages containing files with different encryption types
- **Default behavior:** Test that ML-KEM is used when no encryption type is specified
- **User selection:** Test that users can choose encryption type per file
- **Package operations:** Test all package operations work with mixed encryption types
- **Key management:** Test that appropriate keys are used for each encryption type

## 2. File Validation Testing Requirements

This section defines testing requirements for file validation functionality.

### 2.1 Empty File Testing

- **Empty file acceptance:** Test that files with zero bytes are successfully added to packages
- **Empty file retrieval:** Test that empty files can be extracted and retrieved correctly
- **Empty file integrity:** Test that empty files maintain their integrity during package operations
- **Empty file metadata:** Test that empty files have correct metadata (name, size, timestamps)

### 2.2 Path Normalization Testing

- **Tar-like path normalization:** Test that paths like "dir//file.txt" become "dir/file.txt"
- **Relative reference resolution:** Test that "dir/./file.txt" becomes "dir/file.txt"
- **Parent directory resolution:** Test that "dir/../file.txt" becomes "file.txt"
- **Multiple normalization:** Test that "dir/./../subdir//file.txt" becomes "subdir/file.txt"

### 2.3 Compression Error Handling Testing

- **Compression failure testing:** Test that compression algorithm failures return appropriate errors
- **Memory exhaustion:** Test that insufficient memory during compression returns errors
- **Invalid data handling:** Test that data that cannot be compressed returns errors
- **No fallback behavior:** Test that failed compression does not result in storing uncompressed data

### 2.4 Hash-based Deduplication Testing

- **Processing order validation:** Test that deduplication occurs AFTER compression/encryption, not on raw file content
- **Processed content hash calculation:** Test that content hashes are correctly calculated on processed content (compressed/encrypted)
- **Duplicate detection accuracy:** Test that files with identical processed content are properly identified
- **Single storage:** Test that duplicate processed content is only stored once in the package
- **Path preservation:** Test that all paths to duplicate content are preserved and accessible
