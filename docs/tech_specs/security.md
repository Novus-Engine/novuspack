# NovusPack Technical Specifications - Security and Encryption

This document provides a comprehensive overview of NovusPack's security architecture, including package signing, encryption implementation, validation, and integrity verification systems.

---

## 0. Overview

NovusPack implements a multi-layered security architecture designed to provide comprehensive protection for modern package archives.
The security system combines quantum-safe cryptography, multiple signature support, and robust validation mechanisms to ensure package integrity, authenticity, and confidentiality across all use cases.

### 0.1 Cross-References

- [Main Index](_main.md) - Central navigation for all NovusPack specifications
- [System Overview](_overview.md) - High-level system architecture
- [Package File Format](package_file_format.md) - .npk format structure and incremental signatures
- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation
- [File Validation](file_validation.md) - File validation and transparency requirements

---

## 1. Security Architecture Overview

### 1.1 Security Layers

NovusPack implements a comprehensive security architecture with multiple layers of protection:

1. **Package Integrity**: Digital signatures for package authenticity and integrity
2. **Content Encryption**: Quantum-safe encryption for sensitive asset protection
3. **File Validation**: Comprehensive validation of file content and metadata
4. **Access Control**: Per-file encryption selection and security metadata
5. **Transparency**: Antivirus-friendly design with inspectable structure

### 1.2 Security Principles

- **Defense in Depth**: Multiple security layers provide comprehensive protection
- **Quantum-Safe**: Future-proof cryptography using NIST PQC standards
- **Industry Standard**: Compatible with existing security infrastructure
- **Transparent**: Package structure is inspectable and antivirus-friendly
- **Flexible**: Granular control over security features per file and package

---

## 2. Package Signing System

**Cross-Reference**: For complete signature management, validation, and implementation details, see [Digital Signature API](api_signatures.md).

### 2.1 Security Overview

NovusPack supports multiple digital signatures per package, bringing it in line with industry standards while maintaining unique quantum-safe signature advantages.

#### 2.1.1 Key Security Features

- **Quantum-Safe Signatures**: ML-DSA and SLH-DSA support for future-proof security
- **Incremental Signing**: Signatures are appended sequentially without requiring separate index
- **Immutability Protection**: Signed packages are protected from unauthorized modifications
- **Multiple Algorithm Support**: Traditional (PGP, X.509) and quantum-safe (ML-DSA, SLH-DSA) algorithms

### 2.3 Signature Types and Algorithms

#### 2.3.1 ML-DSA (CRYSTALS-Dilithium)

- **Algorithm**: NIST PQC Standard ML-DSA
- **Security Levels**: Support for all three security levels
  - Level 2: ~2,420-byte signatures (128-bit security)
  - Level 3: ~3,293-byte signatures (192-bit security)
  - Level 5: ~4,595-byte signatures (256-bit security)
- **Performance**: Optimized for archive signing
- **Key Management**: Secure key generation and storage

#### 2.3.2 SLH-DSA (SPHINCS+)

- **Algorithm**: NIST PQC Standard SLH-DSA
- **Security Levels**: Support for all three security levels
  - Level 1: ~7,856-byte signatures (128-bit security)
  - Level 3: ~12,272-byte signatures (192-bit security)
  - Level 5: ~17,088-byte signatures (256-bit security)
- **Performance**: Stateless hash-based signatures
- **Key Management**: Single-use key generation

#### 2.3.3 PGP (OpenPGP)

- **Algorithm**: Traditional PGP signatures
- **Compatibility**: Full OpenPGP standard compliance
- **Key Management**: PGP keyring integration
- **Use Cases**: Developer signatures, community verification

#### 2.3.4 X.509/PKCS#7

- **Algorithm**: Certificate-based signatures
- **Compatibility**: Enterprise PKI integration
- **Key Management**: X.509 certificate chains
- **Use Cases**: Corporate signing, code signing certificates

---

## 3. Encryption System

### 3.1 Quantum-Safe Encryption

NovusPack implements a dual encryption strategy combining quantum-safe and traditional encryption:

- **Primary**: ML-KEM (CRYSTALS-Kyber) for quantum-safe key exchange
- **Secondary**: AES-256-GCM for compatibility and performance
- **Per-File Selection**: Users can choose encryption type per file

### 3.2 Encryption Algorithms

#### 3.2.1 ML-KEM (CRYSTALS-Kyber)

- **Algorithm**: NIST PQC Standard ML-KEM
- **Security Levels**: Support for all three security levels
  - Level 1: 800-byte public keys (128-bit security)
  - Level 3: 1,184-byte public keys (192-bit security)
  - Level 5: 1,568-byte public keys (256-bit security)
- **Performance**: Optimized for file encryption
- **Key Management**: Secure key generation, storage, and distribution

#### 3.2.2 AES-256-GCM

- **Algorithm**: Advanced Encryption Standard with Galois/Counter Mode
- **Key Size**: 256-bit keys
- **Performance**: High-speed encryption for large files
- **Compatibility**: Industry standard for maximum compatibility

### 3.3 Per-File Encryption

Selective file encryption allows users to choose which files to encrypt:

- **User Selection**: Choose encryption per file during package creation
- **Mixed Packages**: Packages can contain both encrypted and unencrypted files
- **Performance Optimization**: Unencrypted files accessed without decryption overhead
- **Security Granularity**: Fine-grained control over content protection

### 3.4 Encryption Implementation Details

#### 3.4.1 ML-KEM Key Management

The ML-KEM implementation provides quantum-safe key generation and management:

- **Key Generation**: Generate ML-KEM keys at specified security levels (1-5)
- **Encryption/Decryption**: Encrypt and decrypt data using ML-KEM keys
- **Key Access**: Retrieve public keys and security level information
- **Key Structure**: Secure storage of public and private key data

For complete function signatures and implementation details, see [File Management API](api_file_management.md) and [API Signatures Index](api_func_signatures_index.md).

#### 3.4.2 Per-File Encryption Operations

The per-file encryption system provides granular control over file protection:

- **AddFileWithEncryption**: Add files with specific encryption types
- **Encryption Types**: Support for None, AES-256-GCM, and ML-KEM encryption
- **GetFileEncryptionType**: Retrieve encryption type for specific files
- **GetEncryptedFiles**: List all encrypted files in the package

For complete function signatures and implementation details, see [File Management API](api_file_management.md) and [API Signatures Index](api_func_signatures_index.md).

#### 3.4.3 Dual Encryption Strategy

- **Default Encryption**: ML-KEM is the default encryption method for new packages
- **ML-KEM Benefits**: Full quantum resistance with optimized performance for file archives
- **AES Support**: AES-256-GCM maintained for compatibility and user preference
- **Per-File Selection**: Users can choose encryption type per file (ML-KEM, AES, or none)
- **Backward Compatibility**: Existing AES-encrypted packages continue to work
- **Hybrid Approach**: ML-KEM for key exchange + AES-256-GCM for data encryption

#### 3.4.4 Implementation Considerations

- **CIRCL Library**: Use Cloudflare's CIRCL for Go implementation of quantum-safe algorithms
- **Standard Library**: Use Go's crypto/aes and crypto/cipher for AES implementation
- **Key Storage**: Secure storage for both quantum-safe and AES keys
- **Performance**: Optimize both encryption methods for large archive packages
- **Compatibility**: Maintain backward compatibility with existing packages while supporting dual encryption

---

## 4. Security Validation and Integrity

### 4.1 Package Validation

#### 4.1.1 Comprehensive Validation

The package validation system provides comprehensive security validation:

- **Validate**: Validate package format, structure, and integrity
- **ValidateAllSignatures**: Validate all signatures in the package
- **ValidateIntegrity**: Validate package integrity using checksums
- **GetSecurityStatus**: Get comprehensive security status information

For complete function signatures and implementation details, see [Security Validation API](api_security.md) and [API Signatures Index](api_func_signatures_index.md).

#### 4.1.2 Security Status Information

The security validation system provides comprehensive status information:

- **Signature Information**: Count of signatures, valid signatures, and trusted signatures
- **Validation Results**: Individual signature validation results
- **Checksum Status**: Information about checksums and their validity
- **Security Level**: Overall security assessment of the package
- **Error Reporting**: Detailed validation errors and issues

For complete structure definitions and implementation details, see [Security Validation API](api_security.md) and [API Signatures Index](api_func_signatures_index.md).

### 4.2 File Validation

#### 4.2.1 Content Validation

- **Empty Files**: Supported and valid
- **Nil Data**: Prohibited and rejected
- **Content Integrity**: Validated before package addition
- **Path Validation**: Normalized and validated file paths

#### 4.2.2 Transparency Requirements

- **No Obfuscation**: Package format is transparent and inspectable
- **Antivirus-Friendly**: Headers and indexes designed for easy scanning
- **Standard Operations**: Uses standard file system operations
- **Clear Structure**: Well-documented and readable package structure

---

## 5. Security Metadata and Access Control

### 5.1 Per-File Security Metadata

#### 5.1.1 Security Classification

- **Script Validation**: Mark files requiring Lua script validation
- **Resource Limits**: Memory and CPU limits for file processing
- **Network Access**: Mark files that should not access network
- **Security Level**: Security classification for file content

#### 5.1.2 Access Control

- **Encryption Selection**: Per-file encryption type selection
- **Compression Selection**: Per-file compression algorithm selection
- **Security Flags**: File-specific security options and restrictions
- **Metadata Protection**: Secure storage of sensitive metadata

### 5.2 Package-Level Security

#### 5.2.1 Security Flags

Package-level security flags provide comprehensive security configuration:

- **Bit 15-8**: Signature features
  - Bit 7: Multiple signatures enabled
  - Bit 6: Has quantum-safe signatures
  - Bit 5: Has traditional signatures
  - Bit 4: Has timestamps
  - Bit 3: Has metadata
  - Bit 2: Has chain validation
  - Bit 1: Has revocation
  - Bit 0: Has expiration

#### 5.2.2 Vendor and Application Identification

- **VendorID**: Storefront/platform identifier for trusted sources
- **AppID**: Application/game identifier for package association
- **LocaleID**: Locale identifier for path encoding
- **CreatorID**: Creator identifier for package attribution

---

## 6. Industry Standard Compliance

### 6.1 Comparison with Industry Standards

| Feature                 | NovusPack | PGP Files | X.509/PKCS#7 | Windows Authenticode | macOS Code Signing |
| ----------------------- | --------- | --------- | ------------ | -------------------- | ------------------ |
| **Multiple Signatures** | Yes       | Yes       | Yes          | Yes                  | Yes                |
| **Signature Types**     | 4 types   | 1 type    | 1 type       | 1 type               | 1 type             |
| **Quantum-Safe**        | Yes       | No        | No           | No                   | No                 |
| **Cross-Platform**      | Yes       | Yes       | Yes          | No                   | No                 |
| **Key Management**      | Multiple  | PGP       | X.509        | Windows cert store   | macOS keychain     |

### 6.2 NovusPack Security Advantages

1. **Quantum-Safe Signatures**: First package format with ML-DSA/SLH-DSA support
2. **Unified Format**: Single format supporting multiple signature types
3. **Cross-Platform**: Works on all major operating systems
4. **Flexible Encryption**: Per-file encryption selection with quantum-safe options
5. **Transparent Design**: Antivirus-friendly and easily inspectable
6. **Industry Compliance**: Compatible with existing security infrastructure

---

## 7. Implementation Considerations

### 7.1 Security Best Practices

#### 7.1.1 Key Management

- **Secure Generation**: Use cryptographically secure random number generators
- **Key Storage**: Implement secure key storage mechanisms
- **Key Rotation**: Support for key rotation and renewal
- **Access Control**: Implement proper access controls for private keys

#### 7.1.2 Signature Validation

- **Multiple Validation**: Validate all signatures, not just one
- **Trust Verification**: Implement proper trust chain validation
- **Timestamp Verification**: Validate signature timestamps
- **Revocation Checking**: Check for revoked certificates and keys

### 7.2 Performance Considerations

#### 7.2.1 Signature Performance

- **Batch Validation**: Validate multiple signatures efficiently
- **Caching**: Cache validation results for performance
- **Parallel Processing**: Use parallel processing for large packages
- **Memory Management**: Optimize memory usage for signature operations

#### 7.2.2 Encryption Performance

- **Selective Encryption**: Only encrypt sensitive assets
- **Hardware Acceleration**: Use hardware acceleration when available
- **Streaming**: Support streaming for large encrypted files
- **Compression**: Combine encryption with compression efficiently

---

## 8. Comment Security and Injection Prevention

### 8.1 Comment Security Architecture

NovusPack implements comprehensive security measures to prevent code execution and malicious injection attacks through comment sections, including package comments and signature comments.

#### 8.1.1 Security Principles

- **No Code Execution**: Comments are treated as pure text data with no executable content
- **Input Sanitization**: All comment data is sanitized and validated before storage
- **Encoding Validation**: Strict UTF-8 validation prevents encoding-based attacks
- **Length Limits**: Enforced maximum lengths prevent buffer overflow attacks
- **Character Filtering**: Dangerous characters and sequences are filtered or escaped

#### 8.1.2 Package Comment Security

Package comments are subject to strict security validation:

- **UTF-8 Validation**: All comment data must be valid UTF-8 encoding
- **Length Limits**: Maximum 65,535 bytes per package comment
- **Character Filtering**: Control characters (0x00-0x1F, 0x7F-0x9F) are filtered
- **No Null Bytes**: Null bytes (0x00) are prohibited in comments
- **No Script Tags**: HTML/XML script tags are escaped or filtered
- **No Command Injection**: Shell command characters are escaped

#### 8.1.3 Signature Comment Security

Signature comments have additional security requirements:

- **UTF-8 Validation**: All signature comment data must be valid UTF-8 encoding
- **Length Limits**: Maximum 4,095 bytes per signature comment
- **Character Filtering**: Same filtering as package comments plus additional restrictions
- **No Executable Content**: No executable code, scripts, or commands allowed
- **Signature Validation**: Comments are included in signature validation to prevent tampering
- **Audit Trail**: All signature comments are logged for security auditing

### 8.2 Input Validation and Sanitization

#### 8.2.1 Comment Validation Process

All comment data undergoes comprehensive validation:

1. **Encoding Check**: Verify valid UTF-8 encoding
2. **Length Validation**: Check against maximum length limits
3. **Character Filtering**: Remove or escape dangerous characters
4. **Content Scanning**: Scan for potential injection patterns
5. **Sanitization**: Apply appropriate sanitization based on content type
6. **Final Validation**: Final check before storage

#### 8.2.2 Dangerous Pattern Detection

The system detects and prevents common injection patterns:

- **Script Injection**: `<script>`, `javascript:`, `vbscript:` patterns
- **Command Injection**: Shell metacharacters (`;`, `|`, `&`, `$`, `` ` ``, `\`)
- **SQL Injection**: SQL keywords and special characters
- **Path Traversal**: `../`, `..\\`, absolute path patterns
- **Unicode Attacks**: Unicode normalization attacks and homograph attacks
- **Control Characters**: All control characters and escape sequences

#### 8.2.3 Sanitization Methods

Different sanitization methods are applied based on content:

- **HTML Escaping**: Convert `<`, `>`, `&`, `"`, `'` to HTML entities
- **URL Encoding**: Encode special characters in URL contexts
- **Character Removal**: Remove dangerous characters entirely
- **Character Replacement**: Replace dangerous characters with safe alternatives
- **Length Truncation**: Truncate overly long content to safe limits

### 8.3 Security Implementation

#### 8.3.1 Comment Storage Security

- **Immutable Storage**: Comments are stored in immutable sections after signing
- **Integrity Protection**: Comments are protected by digital signatures
- **Access Control**: Comments are read-only after package creation
- **Audit Logging**: All comment modifications are logged for security auditing

#### 8.3.2 Runtime Security

- **Safe Display**: Comments are safely displayed without code execution
- **Context Isolation**: Comments are isolated from executable contexts
- **Memory Protection**: Comments are stored in protected memory regions
- **Buffer Overflow Prevention**: Strict bounds checking prevents buffer overflows

### 8.4 Security Testing Requirements

#### 8.4.1 Comment Security Testing

- **Injection Testing**: Test resistance to all types of injection attacks
- **Encoding Testing**: Test with various character encodings and Unicode attacks
- **Length Testing**: Test with maximum and oversized comment lengths
- **Pattern Testing**: Test with known malicious patterns and sequences
- **Sanitization Testing**: Verify proper sanitization of dangerous content

#### 8.4.2 Signature Comment Testing

- **Tamper Testing**: Test resistance to signature comment tampering
- **Validation Testing**: Test signature validation with malicious comments
- **Audit Testing**: Test security audit logging for comment modifications
- **Performance Testing**: Test security measures don't impact performance

---

## 9. Security Testing and Validation

### 9.1 Testing Requirements

#### 9.1.1 Signature Testing

- **Multiple Signature Creation**: Test creating packages with multiple signatures
- **Signature Validation**: Test validation of all signature types
- **Invalid Signature Handling**: Test handling of invalid or corrupted signatures
- **Performance Testing**: Test signature performance with large numbers of signatures

#### 9.1.2 Encryption Testing

- **Encryption/Decryption**: Test all encryption algorithms
- **Key Management**: Test key generation, storage, and retrieval
- **Performance Testing**: Test encryption performance with various file sizes
- **Compatibility Testing**: Test compatibility with existing encrypted packages

### 9.2 Security Validation

#### 9.2.1 Penetration Testing

- **Signature Bypass**: Attempt to bypass signature validation
- **Encryption Bypass**: Attempt to bypass encryption
- **Metadata Manipulation**: Test resistance to metadata manipulation
- **Format Attacks**: Test resistance to malformed package attacks

#### 9.2.2 Compliance Testing

- **Industry Standards**: Verify compliance with industry standards
- **Interoperability**: Test interoperability with other security systems
- **Cross-Platform**: Test security features across all supported platforms
- **Performance**: Verify security features don't significantly impact performance
