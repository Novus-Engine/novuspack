# NovusPack Technical Specifications - Digital Signature API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Signature Management](#1-signature-management)
  - [1.1 Multiple Signature Management (Incremental Signing)](#11-multiple-signature-management-incremental-signing)
  - [1.2 Incremental Signing Implementation](#12-incremental-signing-implementation)
  - [1.3 Immutability Check](#13-immutability-check)
- [2. Signature Types](#2-signature-types)
  - [2.1 Signature Type Constants](#21-signature-type-constants)
  - [2.2 Signature Information Structure](#22-signature-information-structure)
  - [2.3 ML-DSA (CRYSTALS-Dilithium) Implementation](#23-ml-dsa-crystals-dilithium-implementation)
  - [2.4 SLH-DSA (SPHINCS+) Implementation](#24-slh-dsa-sphincs-implementation)
  - [2.5 PGP (OpenPGP) Implementation](#25-pgp-openpgp-implementation)
  - [2.6 X.509/PKCS#7 Implementation](#26-x509pkcs7-implementation)
  - [2.7 Signature Validation](#27-signature-validation)
  - [2.8 Existing Package Signing](#28-existing-package-signing)
  - [2.9 Signing Key Management](#29-signing-key-management)
- [3. Comparison with Other Signed File Implementations](#3-comparison-with-other-signed-file-implementations)
  - [3.1 Industry Standard Comparison](#31-industry-standard-comparison)
  - [3.2 NovusPack Advantages](#32-novuspack-advantages)
  - [3.3 Industry Standard Compliance](#33-industry-standard-compliance)
  - [3.4 Signature Size Comparison](#34-signature-size-comparison)
  - [3.5 Verification Performance](#35-verification-performance)
- [4. Generic Signature Patterns](#4-generic-signature-patterns)
  - [4.1 Generic Signature Strategy Interface](#41-generic-signature-strategy-interface)
  - [4.2 Generic Signature Configuration](#42-generic-signature-configuration)
  - [4.3 Generic Signature Validation](#43-generic-signature-validation)
- [5. Error Handling](#5-error-handling)
  - [5.1 Structured Error System](#51-structured-error-system)
  - [5.2 Common Signature Error Types](#52-common-signature-error-types)
  - [5.3 Structured Error Examples](#53-structured-error-examples)

---

## 0. Overview

This document defines the digital signature API for the NovusPack system, including signature management, types, validation, and comparison with industry standards.

**Note:** This is the authoritative source for all signature implementation details.
The [Package File Format](package_file_format.md) document contains only the binary format specifications.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Package Writing Operations](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Security Validation API](api_security.md) - Generic encryption patterns and type-safe security

## 1. Signature Management

**Note:** This API implements the incremental signing approach as defined in [Package File Format - Digital Signatures](package_file_format.md#7-digital-signatures-section-optional).

### 1.1 Multiple Signature Management (Incremental Signing)

**Function Hierarchy**: High-level signing functions (`SignPackage*`) internally call `AddSignature` after generating the signature data.

```go
// AddSignature adds a new digital signature (appends incrementally)
// LOW-LEVEL: Use when you have pre-computed signature data
// Automatically sets the "Has signatures" bit (Flags Bit 0) and SignatureOffset if this is the first signature
// Ensures signature integrity by validating header state before signing
func (p *Package) AddSignature(ctx context.Context, signatureType uint32, signatureData []byte, flags uint32) error

// RemoveSignature removes signature by index and all later signatures
func (p *Package) RemoveSignature(ctx context.Context, signatureIndex int) error

// GetSignatureCount gets total number of signatures
func (p *Package) GetSignatureCount(ctx context.Context) int

// GetSignature gets signature by index
func (p *Package) GetSignature(ctx context.Context, index int) (SignatureInfo, error)

// GetAllSignatures gets all signatures
func (p *Package) GetAllSignatures(ctx context.Context) []SignatureInfo

// ClearAllSignatures removes all signatures
func (p *Package) ClearAllSignatures(ctx context.Context) error
```

#### 1.1.1 Implementation Requirements

The `AddSignature` function must:

1. **Check if this is the first signature**:

   - If `SignatureOffset == 0`, this is the first signature
   - Set the "Has signatures" bit (Bit 0) in package header flags to 1
   - Set `SignatureOffset` to point to the new signature location

2. **Validate header state**:

   - Ensure the package header is in a valid state for signing
   - Verify that no content modifications have occurred since package creation
   - Check that the package is not already signed if this is not the first signature

3. **Append signature**:

   - Add signature metadata header (18 bytes)
   - Add signature comment (if provided)
   - Add signature data
   - Update file size and any necessary offsets

4. **Maintain immutability** (see [Immutability Check](#13-immutability-check)):
   - After the first signature, the entire package becomes immutable
   - Only signature addition operations are allowed on signed packages

**Implementation Constraint**: The signature bit must be set before the first signature is added because:

- The signature validates the entire header, including the flags field
- Once signed, the header becomes immutable and cannot be modified
- Setting the bit after signing would invalidate the signature
- This ensures signature integrity and prevents tampering

### 1.2 Incremental Signing Implementation

- First signature signs all content up to its own metadata and signature comment and is appended at the end of the file
- Subsequent signatures sign all content up to that point (including previous signatures), its metadata and signature comment and are appended
- Each signature validates all content up to that point, including its own metadata and comment
- Signature removal removes the signature and all later signatures
- No separate signature index - signatures are stored directly

#### 1.2.1 Adding Subsequent Signatures

1. New signature signs all content up to that point (including previous signatures), its metadata and signature comment
2. New signature metadata header (18 bytes) + signature comment + signature data is appended to the end of the file
3. All previous signatures remain valid and unchanged

#### 1.2.2 Signature Validation Process

- Each signature validates all content up to its creation point, including its own metadata
- First signature validates all content from header through its own metadata and signature comment
- Subsequent signatures validate all content up to that point (including previous signatures), plus their own metadata and signature comment
- This ensures incremental validation without invalidating existing signatures
- Signature metadata and comments are protected by the signature itself, preventing tampering

#### 1.2.3 Key Implementation Points

- The `SignatureOffset` in the header points directly to the first signature
- Additional signatures follow immediately after the first signature
- No separate signature index is needed - signatures are read sequentially
- Each signature is self-contained with its own metadata header

### 1.3 Immutability Check

This is the authoritative source for immutability requirements. All signature functions must follow these rules:

- **Pre-signature check**: All write operations must check `SignatureOffset > 0` before proceeding
- **Post-signature restrictions**: If signed, only signature addition operations are allowed
- **Header protection**: Header modifications are prohibited on signed packages
- **Content protection**: File entries, file data, file index, and package comment cannot be modified after first signature
- **Signature integrity**: The signature bit and SignatureOffset cannot be changed after first signature without invalidating the signature

## 2. Signature Types

### 2.1 Signature Type Constants

- **SignatureTypeNone (0x00)**: No signature
- **SignatureTypeMLDSA (0x01)**: ML-DSA (CRYSTALS-Dilithium)
- **SignatureTypeSLHDSA (0x02)**: SLH-DSA (SPHINCS+)
- **SignatureTypePGP (0x03)**: PGP (OpenPGP)
- **SignatureTypeX509 (0x04)**: X.509/PKCS#7
- **0x05-0xFF**: Reserved for future signature types

### 2.2 Signature Information Structure

#### 2.2.1 SignatureInfo struct

- **Type (uint32)**: Signature type identifier
- **Size (uint32)**: Size of signature data in bytes
- **Offset (uint64)**: Offset to signature data from start of file
- **Flags (uint32)**: Signature-specific flags
- **Timestamp (uint32)**: Unix timestamp when signature was created
- **Data ([]byte)**: Raw signature data
- **Algorithm (string)**: Algorithm name/description
- **SecurityLevel (int)**: Security level (1-5)
- **Valid (bool)**: Whether signature is valid
- **Error (string)**: Error message if validation failed

#### 2.2.2 SignatureValidationResult struct

- **Index (int)**: Signature index in the package
- **Type (uint32)**: Signature type identifier
- **Valid (bool)**: Whether signature is valid
- **Trusted (bool)**: Whether signature is trusted
- **Error (string)**: Error message if validation failed
- **Timestamp (uint32)**: When signature was created
- **PublicKey ([]byte)**: Public key used for validation (if available)

### 2.3 ML-DSA (CRYSTALS-Dilithium) Implementation

- **Algorithm**: NIST PQC Standard ML-DSA
- **Security Levels**: Support for all three security levels (see [Signature Size Comparison](#34-signature-size-comparison) for detailed sizes)
- **Performance**: Optimized for package signing
- **Key Management**: Secure key generation and storage

### 2.4 SLH-DSA (SPHINCS+) Implementation

- **Algorithm**: NIST PQC Standard SLH-DSA
- **Security Levels**: Support for all three security levels (see [Signature Size Comparison](#34-signature-size-comparison) for detailed sizes)
- **Performance**: Stateless hash-based signatures
- **Key Management**: Single-use key generation

### 2.5 PGP (OpenPGP) Implementation

- **Algorithm**: OpenPGP standard (RFC 4880)
- **Key Types**: RSA, DSA, ECDSA, EdDSA
- **Key Sizes**: Variable based on algorithm (RSA: 2048-4096 bits, ECDSA: P-256/P-384/P-521)
- **Performance**: Fast verification, moderate signing speed
- **Key Management**: PGP keyring support with passphrase protection

### 2.6 X.509/PKCS#7 Implementation

- **Algorithm**: X.509 certificates with PKCS#7 signatures
- **Key Types**: RSA, ECDSA, EdDSA
- **Certificate Chains**: Full certificate chain validation
- **Performance**: Fast verification with certificate chain validation
- **Key Management**: X.509 certificate and private key files

### 2.7 Signature Validation

```go
// General signature validation
func (p *Package) ValidateSignature(ctx context.Context, publicKey []byte) error
func (p *Package) ValidateSignatureWithKey(ctx context.Context, keyFile string) error
func (p *Package) GetSignatureStatus() SignatureStatus

// PGP-specific validation
func (p *Package) ValidatePGPSignature(ctx context.Context, keyringFile string) error
func (p *Package) ValidatePGPSignatureWithKey(ctx context.Context, keyFile string) error

// X.509-specific validation
func (p *Package) ValidateX509Signature(ctx context.Context, certFile string, caFile string) error
func (p *Package) ValidateX509SignatureWithChain(ctx context.Context, certChain []*x509.Certificate) error
```

### 2.8 Existing Package Signing

**Function Hierarchy**: All `SignPackage*` functions internally call `AddSignature` after generating the signature data.

```go
// HIGH-LEVEL: Use when you have a private key and want to generate + add signature
// Internally calls AddSignature after generating signature data
func (p *Package) SignPackage(ctx context.Context, privateKey []byte, signatureType uint32) error
func (p *Package) SignPackageWithKeyFile(ctx context.Context, keyFile string, signatureType uint32) error

// Sign package with key generation
func (p *Package) SignPackageWithNewKey(ctx context.Context, signatureType uint32, securityLevel int) (*SigningKey, error)

// PGP-specific signing (internally calls AddSignature)
func (p *Package) SignPackageWithPGP(ctx context.Context, keyFile string, passphrase string) error
func (p *Package) SignPackageWithPGPKeyring(ctx context.Context, keyringFile string, keyID string, passphrase string) error

// X.509-specific signing (internally calls AddSignature)
func (p *Package) SignPackageWithX509(ctx context.Context, certFile string, keyFile string, passphrase string) error
func (p *Package) SignPackageWithX509Chain(ctx context.Context, certChain []*x509.Certificate, privateKey []byte) error

// UpdateSignature replaces the most recent signature with new signature data
// Internally calls AddSignature after removing the previous signature
func (p *Package) UpdateSignature(ctx context.Context, newSignature []byte, signatureType uint32) error
```

#### 2.8.1 Implementation Requirements

The `SignPackage` functions follow the same implementation requirements as `AddSignature` (see [Implementation Requirements](#111-implementation-requirements)) with one difference:

- **Step 3**: Generate signature using the provided private key, then append signature data
- All other steps (signature bit setting, header validation, immutability) are identical

#### 2.8.2 Function Usage Guide

##### 2.8.2.1 When to use AddSignature (Low-Level)

- You have pre-computed signature data
- You want direct control over the signature addition process
- You're implementing custom signature generation logic
- You need to add signatures from external signature services

##### 2.8.2.2 When to use SignPackage\* functions (High-Level)

- You have a private key and want to generate + add a signature
- You want the convenience of automatic signature generation
- You're using standard signature types (ML-DSA, SLH-DSA, PGP, X.509)
- You want automatic key management and signature generation

##### 2.8.2.3 Implementation Pattern

**Implementation Pattern**: High-level functions internally follow this pattern:

1. Generate signature data using the provided private key
2. Call `AddSignature` with the generated signature data
3. Handle any errors from the signature generation or addition process

### 2.9 Signing Key Management

```go
type SigningKey struct {
    PrivateKey []byte
    PublicKey  []byte
    Type       uint32
    Level      int
}

func (p *Package) GenerateSigningKey(ctx context.Context, signatureType uint32, securityLevel int) (*SigningKey, error)
func (p *Package) SaveSigningKey(ctx context.Context, key *SigningKey, keyFile string) error
func (p *Package) LoadSigningKey(ctx context.Context, keyFile string) (*SigningKey, error)
```

## 3. Comparison with Other Signed File Implementations

### 3.1 Industry Standard Comparison

| Feature                 | NovusPack                             | PGP Files                | X.509/PKCS#7             | Windows Authenticode     | macOS Code Signing       |
| ----------------------- | ------------------------------------- | ------------------------ | ------------------------ | ------------------------ | ------------------------ |
| **Signature Location**  | End of file                           | Detached/Inline          | End of file              | End of file              | End of file              |
| **Header Metadata**     | [OK] Extended header                  | [NO] Separate files      | [OK] PKCS#7 structure    | [OK] PE structure        | [OK] Mach-O structure    |
| **Multiple Signatures** | [OK] Multiple signatures              | [OK] Multiple signatures | [OK] Multiple signatures | [OK] Multiple signatures | [OK] Multiple signatures |
| **Signature Types**     | 4 types (ML-DSA, SLH-DSA, PGP, X.509) | 1 type (PGP)             | 1 type (X.509)           | 1 type (Authenticode)    | 1 type (Code Signing)    |
| **Quantum-Safe**        | [OK] ML-DSA/SLH-DSA                   | [NO] No                  | [NO] No                  | [NO] No                  | [NO] No                  |
| **Cross-Platform**      | [OK] All platforms                    | [OK] All platforms       | [OK] All platforms       | [NO] Windows only        | [NO] macOS only          |
| **Key Management**      | [OK] Multiple types                   | [OK] PGP keyrings        | [OK] X.509 certificates  | [OK] Windows cert store  | [OK] macOS keychain      |
| **Signature Size**      | Variable (100-17K bytes)              | Variable (100-1K bytes)  | Variable (200-2K bytes)  | Variable (200-2K bytes)  | Variable (200-2K bytes)  |
| **Verification Speed**  | Fast                                  | Fast                     | Fast                     | Fast                     | Fast                     |
| **Industry Adoption**   | New                                   | High                     | High                     | High (Windows)           | High (macOS)             |

### 3.2 NovusPack Advantages

1. **Quantum-Safe Signatures**: First package format with ML-DSA/SLH-DSA support
2. **Unified Format**: Single format supporting multiple signature types
3. **Cross-Platform**: Works on all platforms unlike platform-specific solutions
4. **Future-Proof**: Extensible header design for new signature types
5. **General-Purpose**: Designed for general-purpose archive applications

### 3.3 Industry Standard Compliance

- **PGP Compatibility**: Follows OpenPGP standard (RFC 4880)
- **X.509 Compliance**: Follows PKCS#7 standard (RFC 2315)
- **Signature Placement**: Follows industry standard (end of file)
- **Hash Algorithm**: Uses SHA-256 (industry standard)
- **Key Management**: Supports standard key formats

### 3.4 Signature Size Comparison

- **NovusPack ML-DSA**: ~2,420-4,595 bytes (quantum-safe)
- **NovusPack SLH-DSA**: ~7,856-17,088 bytes (quantum-safe)
- **NovusPack PGP**: ~100-1,000 bytes (traditional)
- **NovusPack X.509**: ~200-2,000 bytes (traditional)
- **PGP Files**: ~100-1,000 bytes (traditional)
- **X.509/PKCS#7**: ~200-2,000 bytes (traditional)
- **Windows Authenticode**: ~200-2,000 bytes (traditional)
- **macOS Code Signing**: ~200-2,000 bytes (traditional)

### 3.5 Verification Performance

- **NovusPack**: Fast verification with optimized hash calculation
- **PGP**: Fast verification with established algorithms
- **X.509**: Fast verification with certificate chain validation
- **Platform-Specific**: Fast verification with OS integration

## 4. Generic Signature Patterns

The signatures API provides generic signature patterns that extend the generic configuration patterns defined in [api_generics.md](api_generics.md#28-generic-configuration-patterns) for type-safe signature operations.

### 4.1 Generic Signature Strategy Interface

```go
// SignatureStrategy provides type-safe signing for any data type
type SignatureStrategy[T any] interface {
    Sign(ctx context.Context, data T, key SigningKey[T]) (Signature[T], error)
    Verify(ctx context.Context, data T, signature Signature[T], key SigningKey[T]) error
    Type() SignatureType
    Name() string
    KeySize() int
    ValidateKey(ctx context.Context, key SigningKey[T]) error
}

// ByteSignatureStrategy is the concrete implementation for []byte data
type ByteSignatureStrategy interface {
    SignatureStrategy[[]byte]
}

// SigningKey provides type-safe key management for signatures
type SigningKey[T any] struct {
    *Option[T]
    KeyType    SignatureType
    KeyID      string
    CreatedAt  time.Time
    ExpiresAt  *time.Time
    Algorithm  string
}

func NewSigningKey[T any](keyType SignatureType, keyID string, key T) *SigningKey[T]
func (k *SigningKey[T]) GetKey() (T, bool)
func (k *SigningKey[T]) SetKey(key T)
func (k *SigningKey[T]) IsValid() bool
func (k *SigningKey[T]) IsExpired() bool

// Signature provides type-safe signature data
type Signature[T any] struct {
    *Option[T]
    SignatureType SignatureType
    Algorithm     string
    CreatedAt     time.Time
    Data          []byte
}

func NewSignature[T any](sigType SignatureType, data []byte) *Signature[T]
func (s *Signature[T]) GetData() (T, bool)
func (s *Signature[T]) SetData(data T)
func (s *Signature[T]) IsValid() bool
func (s *Signature[T]) GetSignatureType() SignatureType
```

### 4.2 Generic Signature Configuration

```go
// SignatureConfig provides type-safe signature configuration
type SignatureConfig[T any] struct {
    *Config[T]

    // Signature-specific settings
    SignatureType     Option[SignatureType]     // Signature algorithm type
    KeySize          Option[int]                // Key size in bits
    UseTimestamp     Option[bool]              // Include timestamp in signature
    IncludeMetadata  Option[bool]              // Include metadata in signature
    CompressionLevel Option[int]               // Compression level for signed data
}

// SignatureConfigBuilder provides fluent configuration building for signatures
type SignatureConfigBuilder[T any] struct {
    config *SignatureConfig[T]
}

func NewSignatureConfigBuilder[T any]() *SignatureConfigBuilder[T]
func (b *SignatureConfigBuilder[T]) WithSignatureType(sigType SignatureType) *SignatureConfigBuilder[T]
func (b *SignatureConfigBuilder[T]) WithKeySize(size int) *SignatureConfigBuilder[T]
func (b *SignatureConfigBuilder[T]) WithTimestamp(useTimestamp bool) *SignatureConfigBuilder[T]
func (b *SignatureConfigBuilder[T]) WithMetadata(includeMetadata bool) *SignatureConfigBuilder[T]
func (b *SignatureConfigBuilder[T]) Build() *SignatureConfig[T]
```

### 4.3 Generic Signature Validation

```go
// SignatureValidator provides type-safe signature validation
type SignatureValidator[T any] struct {
    *Validator[T]
    signatureRules []SignatureValidationRule[T]
}

// SignatureValidationRule is an alias for the generic ValidationRule
type SignatureValidationRule[T any] = ValidationRule[T]

func (v *SignatureValidator[T]) AddSignatureRule(rule SignatureValidationRule[T])
func (v *SignatureValidator[T]) ValidateSignatureData(ctx context.Context, data T) error
func (v *SignatureValidator[T]) ValidateSignatureKey(ctx context.Context, key SigningKey[T]) error
func (v *SignatureValidator[T]) ValidateSignatureFormat(ctx context.Context, signature Signature[T]) error
```

## 5. Error Handling

### 5.1 Structured Error System

The NovusPack digital signature API uses the comprehensive structured error system defined in [api_core.md](api_core.md#11-structured-error-system).

**Error Types Used**: This API uses the following error types from the core error system:

- **ErrTypeSignature**: Digital signature validation, signing failures, signature format errors
- **ErrTypeValidation**: Input validation errors, invalid parameters, format errors
- **ErrTypeSecurity**: Security-related errors, access denied, authentication failures
- **ErrTypeUnsupported**: Unsupported features, versions, or operations
- **ErrTypeCorruption**: Data corruption, checksum failures, integrity violations

For complete error system documentation, see [Structured Error System](api_core.md#11-structured-error-system).

### 5.2 Common Signature Error Types

The signature API uses structured errors with specific error types for granular error handling. Each error type provides detailed context and categorization.

#### 5.2.1 Specific Signature Error Types

```go
// Signature-specific error types for granular error handling
const (
    ErrSignatureInvalid      = "invalid signature"
    ErrSignatureNotFound     = "signature not found"
    ErrSignatureValidationFailed = "signature validation failed"
    ErrUnsupportedSignatureType = "unsupported signature type"
    ErrInvalidSignatureData  = "invalid signature data"
    ErrSignatureTooLarge     = "signature too large"
    ErrKeyNotFound          = "signing key not found"
    ErrInvalidKey           = "invalid signing key"
    ErrSignatureGenerationFailed = "signature generation failed"
    ErrSignatureCorrupted   = "signature data corrupted"
)

// Typed context structures for modern error handling
type SignatureErrorContext struct {
    SignatureIndex int
    Algorithm      string
    Operation      string
}

type UnsupportedErrorContext struct {
    SignatureType   uint32
    SupportedTypes  []uint32
    Operation       string
}

type SecurityErrorContext struct {
    KeyID     string
    KeyType   string
    Operation string
}

type ValidationErrorContext struct {
    Field    string
    Value    string
    Expected string
}
```

#### 5.2.2 Error Type Mapping

| Error Message                | ErrorType          | Description                      |
| ---------------------------- | ------------------ | -------------------------------- |
| ErrSignatureInvalid          | ErrTypeSignature   | Invalid signature format or data |
| ErrSignatureNotFound         | ErrTypeValidation  | Signature not found at index     |
| ErrSignatureValidationFailed | ErrTypeSignature   | Signature validation failed      |
| ErrUnsupportedSignatureType  | ErrTypeUnsupported | Unsupported signature algorithm  |
| ErrInvalidSignatureData      | ErrTypeValidation  | Invalid signature data format    |
| ErrSignatureTooLarge         | ErrTypeValidation  | Signature exceeds maximum size   |
| ErrKeyNotFound               | ErrTypeSecurity    | Signing key not found            |
| ErrInvalidKey                | ErrTypeSecurity    | Invalid signing key format       |
| ErrSignatureGenerationFailed | ErrTypeSignature   | Signature generation failed      |
| ErrSignatureCorrupted        | ErrTypeCorruption  | Signature data corrupted         |

### 5.3 Structured Error Examples

#### 5.3.1 Creating Signature Errors

```go
// Modern approach using generic error helpers
// Signature validation failure with typed context
err := NewTypedPackageError(ErrTypeSignature, ErrSignatureValidationFailed, nil, SignatureErrorContext{
    SignatureIndex: 0,
    Algorithm:     "ML-DSA",
    Operation:     "ValidateSignature",
})

// Unsupported signature type with typed context
err := NewTypedPackageError(ErrTypeUnsupported, ErrUnsupportedSignatureType, nil, UnsupportedErrorContext{
    SignatureType:   999,
    SupportedTypes: []uint32{1, 2, 3, 4},
    Operation:       "AddSignature",
})

// Key not found with typed context
err := NewTypedPackageError(ErrTypeSecurity, ErrKeyNotFound, nil, SecurityErrorContext{
    KeyID:      "0x12345678",
    KeyType:    "ML-DSA",
    Operation:  "SignPackage",
})
```

#### 5.3.2 Error Handling Patterns

The signature API uses modern error handling patterns with typed context for better debugging and error categorization.

**Error Inspection Pattern**: Use `IsPackageError(err)` to check for structured errors, then switch on `pkgErr.Type` to handle different error categories appropriately.

**Type-Safe Context Access**: Use `GetTypedContext[SignatureErrorContext](err, "context")` to access typed error context for debugging information.

**Error Creation Pattern**: Use `NewTypedPackageError` with appropriate context structures for creating structured errors with type-safe context information.

#### 5.3.3 Function Signatures

The signature API provides several key functions for error handling and validation.
These functions follow the structured error system and use typed context for enhanced debugging capabilities.

#### 5.3.3.1 Core Error Handling Functions

```go
// Internal signature validation with structured error context
func (p *Package) validateSignatureInternal(ctx context.Context, sigData []byte) error
```

**Cross-Reference**: For generic error helper functions like `GetTypedContext[T any]`, `NewTypedPackageError[T any]`, `WithTypedContext[T any]`, `WrapWithContext[T any]`, and `MapError[T, U]`, see [Structured Error System](api_core.md#11-structured-error-system).
