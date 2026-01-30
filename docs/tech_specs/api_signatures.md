# NovusPack Technical Specifications - Digital Signature API

Status: Deferred to v2.

This document is provided as future work specification for signature management, signing, and signature validation.
V1 only enforces signed package immutability based on signature presence, and does not validate signature contents.

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Signature Management](#1-signature-management)
  - [1.1 Multiple Signature Management (Incremental Signing)](#11-multiple-signature-management-incremental-signing)
    - [1.1.1 Package AddSignature Method](#111-packageaddsignature-method)
    - [1.1.2 Package RemoveSignature Method](#112-packageremovesignature-method)
    - [1.1.3 Package GetSignatureCount Method](#113-packagegetsignaturecount-method)
    - [1.1.4 Package GetSignature Method](#114-packagegetsignature-method)
    - [1.1.5 Package GetAllSignatures Method](#115-packagegetallsignatures-method)
    - [1.1.6 Package ClearAllSignatures Method](#116-packageclearallsignatures-method)
    - [1.1.7 Implementation Requirements](#117-implementation-requirements)
  - [1.2 Incremental Signing Implementation](#12-incremental-signing-implementation)
    - [1.2.1 Adding Subsequent Signatures](#121-adding-subsequent-signatures)
    - [1.2.2 Signature Validation Process](#122-signature-validation-process)
    - [1.2.3 Key Implementation Points](#123-key-implementation-points)
  - [1.3 Immutability Check](#13-immutability-check)
- [2. Signature Types](#2-signature-types)
  - [2.1 Signature Type Constants](#21-signature-type-constants)
    - [2.1.1 Signature Type Usage](#211-signature-type-usage)
  - [2.2 Signature Information Structure](#22-signature-information-structure)
    - [2.2.1 SignatureInfo struct](#221-signatureinfo-struct)
    - [2.2.2 SignatureValidationResult struct](#222-signaturevalidationresult-struct)
  - [2.3 ML-DSA (CRYSTALS-Dilithium) Implementation](#23-ml-dsa-crystals-dilithium-implementation)
  - [2.4 SLH-DSA (SPHINCS+) Implementation](#24-slh-dsa-sphincs-implementation)
  - [2.5 PGP (OpenPGP) Implementation](#25-pgp-openpgp-implementation)
  - [2.6 X.509/PKCS#7 Implementation](#26-x509pkcs7-implementation)
  - [2.7 Signature Validation](#27-signature-validation)
    - [2.7.1 General Signature Validation Methods](#271-general-signature-validation-methods)
    - [2.7.2 PGP-Specific Validation Methods](#272-pgp-specific-validation-methods)
    - [2.7.3 X.509-Specific Validation Methods](#273-x509-specific-validation-methods)
  - [2.8 Existing Package Signing](#28-existing-package-signing)
    - [2.8.1 General Signing Methods](#281-general-signing-methods)
    - [2.8.2 PGP-Specific Signing Methods](#282-pgp-specific-signing-methods)
    - [2.8.3 X.509-Specific Signing Methods](#283-x509-specific-signing-methods)
    - [2.8.4 Package UpdateSignature Method](#284-packageupdatesignature-method)
    - [2.8.5 Implementation Requirements (SignPackage)](#285-implementation-requirements-signpackage)
    - [2.8.6 Function Usage Guide](#286-function-usage-guide)
- [3. Comparison with Other Signed File Implementations](#3-comparison-with-other-signed-file-implementations)
  - [3.1 Industry Standard Comparison](#31-industry-standard-comparison)
  - [3.2 NovusPack Advantages](#32-novuspack-advantages)
  - [3.3 Industry Standard Compliance](#33-industry-standard-compliance)
  - [3.4 Signature Size Comparison](#34-signature-size-comparison)
  - [3.5 Verification Performance](#35-verification-performance)
- [4. Generic Signature Patterns](#4-generic-signature-patterns)
  - [4.1 SignatureStrategy Interface](#41-signaturestrategy-interface)
    - [4.1.1 SignatureStrategy Interface Definition](#411-signaturestrategy-interface-definition)
    - [4.1.2 ByteSignatureStrategy Interface](#412-bytesignaturestrategy-interface)
    - [4.1.3 SigningKey Structure](#413-signingkey-structure)
    - [4.1.4 Signature Structure](#414-signature-structure)
    - [4.1.5 GetKey and SetKey Behavior](#415-getkey-and-setkey-behavior)
    - [4.1.6 SignatureStrategy Secure SigningKey Operations with runtime/secret](#416-signaturestrategy-secure-signingkey-operations-with-runtimesecret)
  - [4.2 IsValid and IsExpired Semantics](#42-isvalid-and-isexpired-semantics)
    - [4.2.1 IsValid() Requirements](#421-isvalid-requirements)
    - [4.2.2 IsExpired() Requirements](#422-isexpired-requirements)
    - [4.2.3 Expiration Semantics](#423-expiration-semantics)
    - [4.2.4 Validation Order](#424-validation-order)
  - [4.3 Generic Signature Configuration](#43-generic-signature-configuration)
    - [4.3.1 SignatureConfig Structure](#431-signatureconfig-structure)
    - [4.3.2 SignatureConfigBuilder Structure](#432-signatureconfigbuilder-structure)
  - [4.4 Generic Signature Validation](#44-generic-signature-validation)
    - [4.4.1 SignatureValidator Structure](#441-signaturevalidator-struct)
    - [4.4.2 SignatureValidator AddSignatureRule Method](#443-signaturevalidatortaddsignaturerule-method)
    - [4.4.3 SignatureValidator ValidateSignatureData Method](#444-signaturevalidatortvalidatesignaturedata-method)
    - [4.4.4 SignatureValidator ValidateSignatureKey Method](#445-signaturevalidatortvalidatesignaturekey-method)
    - [4.4.5 SignatureValidator ValidateSignatureFormat Method](#446-signaturevalidatortvalidatesignatureformat-method)
- [5. Error Handling](#5-error-handling)
  - [5.1 Structured Error System](#51-structured-error-system)
  - [5.2 Signature Error Messages](#52-signature-error-messages)
  - [5.3 Signature-Specific Error Context Types](#53-signature-specific-error-context-types)
    - [5.3.1 SignatureErrorContext Structure](#531-signatureerrorcontext-structure)
    - [5.3.2 UnsupportedErrorContext Structure](#532-unsupportederrorcontext-structure)
    - [5.3.3 SecurityErrorContext Structure](#533-securityerrorcontext-structure)
    - [5.3.4 ValidationErrorContext Structure](#534-validationerrorcontext-structure)

---

## 0. Overview

This document defines the digital signature API for the NovusPack system, including signature management, types, validation, and comparison with industry standards.

**Note:** This is the authoritative source for all signature implementation details.
The [Package File Format](package_file_format.md) document contains only the binary format specifications.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Package Writing Operations](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Security Validation API](api_security.md) - Generic encryption patterns and type-safe security

## 1. Signature Management

**Note:** This API implements the incremental signing approach as defined in [Package File Format - Digital Signatures Section](package_file_format.md#8-digital-signatures-section-optional).

### 1.1 Multiple Signature Management (Incremental Signing)

**Function Hierarchy**: High-level signing functions (`SignPackage*`) internally call `AddSignature` after generating the signature data.

#### 1.1.1 Package.AddSignature Method

```go
// AddSignature adds a new digital signature (appends incrementally)
// LOW-LEVEL: Use when you have pre-computed signature data
// Automatically sets the "Has signatures" bit (Flags Bit 0) and SignatureOffset if this is the first signature
// Ensures signature integrity by validating header state before signing
// Returns *PackageError on failure
func (p *Package) AddSignature(ctx context.Context, signatureType uint32, signatureData []byte, flags uint32) error
```

#### 1.1.2 Package.RemoveSignature Method

```go
// RemoveSignature removes signature by index and all later signatures
// Returns *PackageError on failure
func (p *Package) RemoveSignature(ctx context.Context, signatureIndex int) error
```

#### 1.1.3 Package.GetSignatureCount Method

```go
// GetSignatureCount gets total number of signatures
func (p *Package) GetSignatureCount() int
```

#### 1.1.4 Package.GetSignature Method

```go
// GetSignature gets signature by index
// Returns *PackageError on failure
func (p *Package) GetSignature(index int) (SignatureInfo, error)
```

#### 1.1.5 Package.GetAllSignatures Method

```go
// GetAllSignatures gets all signatures
func (p *Package) GetAllSignatures() []SignatureInfo
```

#### 1.1.6 Package.ClearAllSignatures Method

```go
// ClearAllSignatures removes all signatures
// Returns *PackageError on failure
func (p *Package) ClearAllSignatures(ctx context.Context) error
```

#### 1.1.7 Implementation Requirements

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

This section describes signature types supported by the API.

### 2.1 Signature Type Constants

- **SignatureTypeNone (0x00)**: No signature
- **SignatureTypeMLDSA (0x01)**: ML-DSA (CRYSTALS-Dilithium)
- **SignatureTypeSLHDSA (0x02)**: SLH-DSA (SPHINCS+)
- **SignatureTypePGP (0x03)**: PGP (OpenPGP)
- **SignatureTypeX509 (0x04)**: X.509/PKCS#7
- **0x05-0xFF**: Reserved for future signature types

#### 2.1.1 Signature Type Usage

`SignatureType` is used throughout the API to identify signature algorithms.
The constants (0x01, 0x02, etc.) are used both for on-disk representation and in-memory type identification.

### 2.2 Signature Information Structure

This section describes the signature information structure.

#### 2.2.1. SignatureInfo Struct

- **Type (uint32)**: Signature type identifier
- **Size (uint32)**: Size of signature data in bytes
- **Offset (uint64)**: Offset to signature data from start of file
- **Flags (uint32)**: Signature-specific flags
- **Timestamp (uint32)**: Unix timestamp when signature was created
- **Data ([]byte)**: Raw signature data
- **Algorithm (string)**: Algorithm name/description
- **SecurityLevel (int)**: Algorithm security level (signature algorithm specific)
- **Valid (bool)**: Whether signature is valid
- **Error (string)**: Error message if validation failed

#### 2.2.2. SignatureValidationResult Struct

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

This section describes signature validation methods.

#### 2.7.1 General Signature Validation Methods

This section describes general signature validation methods.

##### 2.7.1.1 Package.ValidateSignature Method

```go
// General signature validation
// Returns *PackageError on failure
func (p *Package) ValidateSignature(ctx context.Context, publicKey []byte) error
```

##### 2.7.1.2 Package.ValidateSignatureWithKey Method

```go
// Returns *PackageError on failure
func (p *Package) ValidateSignatureWithKey(ctx context.Context, keyFile string) error
```

##### 2.7.1.3 Package.GetSignatureStatus Method

```go
// GetSignatureStatus returns the current signature status of the package.
func (p *Package) GetSignatureStatus() SignatureStatus
```

#### 2.7.2 PGP-Specific Validation Methods

This section describes PGP-specific signature validation methods.

##### 2.7.2.1 Package.ValidatePGPSignature Method

```go
// PGP-specific validation
// Returns *PackageError on failure
func (p *Package) ValidatePGPSignature(ctx context.Context, keyringFile string) error
```

##### 2.7.2.2 Package.ValidatePGPSignatureWithKey Method

```go
// Returns *PackageError on failure
func (p *Package) ValidatePGPSignatureWithKey(ctx context.Context, keyFile string) error
```

#### 2.7.3 X.509-Specific Validation Methods

This section describes X.509-specific signature validation methods.

##### 2.7.3.1 Package.ValidateX509Signature Method

```go
// X.509-specific validation
// Returns *PackageError on failure
func (p *Package) ValidateX509Signature(ctx context.Context, certFile string, caFile string) error
```

##### 2.7.3.2 Package.ValidateX509SignatureWithChain Method

```go
// Returns *PackageError on failure
func (p *Package) ValidateX509SignatureWithChain(ctx context.Context, certChain []*x509.Certificate) error
```

### 2.8 Existing Package Signing

**Function Hierarchy**: All `SignPackage*` functions internally call `AddSignature` after generating the signature data.

#### 2.8.1 General Signing Methods

This section describes general package signing methods.

##### 2.8.1.1 Package.SignPackage Method

```go
// HIGH-LEVEL: Use when you have a private key and want to generate + add signature
// Internally calls AddSignature after generating signature data
// Returns *PackageError on failure
func (p *Package) SignPackage(ctx context.Context, privateKey []byte, signatureType uint32) error
```

##### 2.8.1.2 Package.SignPackageWithKeyFile Method

```go
// Returns *PackageError on failure
func (p *Package) SignPackageWithKeyFile(ctx context.Context, keyFile string, signatureType uint32) error
```

#### 2.8.2 PGP-Specific Signing Methods

This section describes PGP-specific package signing methods.

##### 2.8.2.1 Package.SignPackageWithPGP Method

```go
// PGP-specific signing (internally calls AddSignature)
// Returns *PackageError on failure
func (p *Package) SignPackageWithPGP(ctx context.Context, keyFile string, passphrase string) error
```

##### 2.8.2.2 Package.SignPackageWithPGPKeyring Method

```go
// Returns *PackageError on failure
func (p *Package) SignPackageWithPGPKeyring(ctx context.Context, keyringFile string, keyID string, passphrase string) error
```

#### 2.8.3 X.509-Specific Signing Methods

This section describes X.509-specific package signing methods.

##### 2.8.3.1 Package.SignPackageWithX509 Method

```go
// X.509-specific signing (internally calls AddSignature)
// Returns *PackageError on failure
func (p *Package) SignPackageWithX509(ctx context.Context, certFile string, keyFile string, passphrase string) error
```

##### 2.8.3.2 Package.SignPackageWithX509Chain Method

```go
// Returns *PackageError on failure
func (p *Package) SignPackageWithX509Chain(ctx context.Context, certChain []*x509.Certificate, privateKey []byte) error
```

#### 2.8.4 Package.UpdateSignature Method

```go
// UpdateSignature replaces the most recent signature with new signature data
// Internally calls AddSignature after removing the previous signature
// Returns *PackageError on failure
func (p *Package) UpdateSignature(ctx context.Context, newSignature []byte, signatureType uint32) error
```

#### 2.8.5 Implementation Requirements (SignPackage)

The `SignPackage` functions follow the same implementation requirements as `AddSignature` (see [Implementation Requirements](#117-implementation-requirements)) with one difference:

- **Step 3**: Generate signature using the provided private key, then append signature data
- All other steps (signature bit setting, header validation, immutability) are identical

##### 2.8.5.1. Secure Signing Operations with Runtime/secret

**MUST Requirements**: All signing operations that use keys MUST execute within Go's `runtime/secret.Do` function to protect sensitive cryptographic material in memory.

- Signature generation operations that access key material (including all `SignPackage*` functions) MUST wrap the key access and signature computation within the secret execution context
- This ensures that keys and intermediate cryptographic values are promptly erased from memory registers and stack frames after use
- The `SignPackageWithKeyFile` and related functions MUST use `runtime/secret.Do` when loading keys from files to prevent key material from persisting in memory
- Implementations MUST ensure that all operations involving key material are executed within the secret execution context to maximize protection against memory analysis attacks

#### 2.8.6 Function Usage Guide

This section provides guidance on when to use different signing functions.

##### 2.8.6.1. When to Use AddSignature (Low-Level)

- You have pre-computed signature data
- You want direct control over the signature addition process
- You're implementing custom signature generation logic
- You need to add signatures from external signature services

##### 2.8.6.2. When to Use SignPackage\* Functions (High-Level)

- You have a private key and want to generate + add a signature
- You want the convenience of automatic signature generation
- You're using standard signature types (ML-DSA, SLH-DSA, PGP, X.509)
- You want automatic key management and signature generation

##### 2.8.6.3 Implementation Pattern

**Implementation Pattern**: High-level functions internally follow this pattern:

1. Generate signature data using the provided private key
2. Call `AddSignature` with the generated signature data
3. Handle any errors from the signature generation or addition process

## 3. Comparison with Other Signed File Implementations

This section compares the NovusPack signature implementation with other signed file formats.

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

The signatures API provides generic signature patterns that extend the generic configuration patterns defined in [api_generics.md](api_generics.md#1-core-generic-types) for type-safe signature operations.

### 4.1 SignatureStrategy Interface

The `SignatureStrategy[T]` interface extends the generic [Core Generic Types](api_generics.md#1-core-generic-types) pattern for signature-specific operations.
`SignatureStrategy[T]` embeds `Strategy[T, Signature[T]]` where the input is the data type T and the output is `Signature[T]`.
The `Process` method from `Strategy[T, Signature[T]]` represents the Sign operation, while `Sign` and `Verify` provide more specific signature operations with key management.

#### 4.1.1 SignatureStrategy Interface Definition

```go
// SignatureStrategy extends Strategy[T, Signature[T]] for signature operations
// The Process method from Strategy represents the Sign operation
// The Strategy.Type() method returns "signature" as the category
type SignatureStrategy[T any] interface {
    Strategy[T, Signature[T]]  // Extends the generic Strategy interface

    Sign(ctx context.Context, data T, key SigningKey[T]) (Signature[T], error)
    Verify(ctx context.Context, data T, signature Signature[T], key SigningKey[T]) error
    SignatureType() SignatureType  // Returns the specific signature algorithm type
    Name() string
    KeySize() int
    ValidateKey(ctx context.Context, key SigningKey[T]) error
}
```

#### 4.1.2 ByteSignatureStrategy Interface

```go
// ByteSignatureStrategy is the concrete implementation for []byte data
type ByteSignatureStrategy interface {
    SignatureStrategy[[]byte]
}
```

#### 4.1.3 SigningKey Structure

This section describes the SigningKey structure for managing signing keys.

##### 4.1.3.1 SigningKey Struct

```go
// SigningKey provides type-safe key management for signatures
// Stores private key material only - public keys are handled separately for verification
// Uses Option[T] internally for type-safe key storage
// All private key material must be handled within runtime/secret.Do for security
type SigningKey[T any] struct {
    *Option[T]  // See [Option Type](api_generics.md#11-option-type) for details
    KeyType    SignatureType
    KeyID      string
    CreatedAt  time.Time
    ExpiresAt  *time.Time
    Data       []byte  // Private key material (for signing operations)
}
```

##### 4.1.3.2 NewSigningKey Function

```go
// NewSigningKey creates a new signing key with the specified type, ID, and key material.
func NewSigningKey[T any](keyType SignatureType, keyID string, key T) *SigningKey[T]
```

##### 4.1.3.3 SigningKey[T].GetKey Method

```go
// GetKey returns the signing key material.
func (k *SigningKey[T]) GetKey() (T, error)
```

##### 4.1.3.4 SigningKey[T].SetKey Method

```go
// SetKey sets the signing key material.
func (k *SigningKey[T]) SetKey(key T)
```

##### 4.1.3.5 SigningKey[T].IsValid Method

```go
// IsValid returns true if the signing key is valid.
func (k *SigningKey[T]) IsValid() bool
```

##### 4.1.3.6 SigningKey[T].IsExpired Method

```go
// IsExpired returns true if the signing key has expired.
func (k *SigningKey[T]) IsExpired() bool
```

#### 4.1.4 Signature Structure

This section describes the Signature structure for managing signatures.

##### 4.1.4.1 Signature Struct

```go
// Signature provides type-safe signature data
// Uses Option[T] internally for type-safe signature data storage
type Signature[T any] struct {
    *Option[T]  // See [Option Type](api_generics.md#11-option-type) for details
    SignatureType SignatureType
    CreatedAt     time.Time
    Data          []byte
}
```

##### 4.1.4.2 NewSignature Function

```go
// Note: Current implementation is simplified for v1 (signatures deferred to v2)
// Future v2 implementation: func NewSignature[T any](sigType SignatureType, data []byte) *Signature[T]
func NewSignature() *Signature
```

##### 4.1.4.3 Signature[T].GetData Method

```go
// GetData returns the signature data.
func (s *Signature[T]) GetData() (T, error)
```

##### 4.1.4.4 Signature[T].SetData Method

```go
// SetData sets the signature data.
func (s *Signature[T]) SetData(data T)
```

##### 4.1.4.5 Signature[T].IsValid Method

```go
// IsValid returns true if the signature is valid.
func (s *Signature[T]) IsValid() bool
```

##### 4.1.4.6 Signature[T].GetSignatureType Method

```go
// GetSignatureType returns the type of the signature.
func (s *Signature[T]) GetSignatureType() SignatureType
```

#### 4.1.5 GetKey and SetKey Behavior

This section describes the behavior of GetKey and SetKey methods.

##### 4.1.5.1 GetKey Behavior

- `GetKey()` returns a **copy** of the key value, not a reference to the original.
- This ensures that modifications to the returned value do not affect the stored key.
- For slice types (`[]byte`, etc.), a deep copy is performed to prevent accidental mutation.
- Returns an error if the key is missing, invalid, or expired.
- The returned copy should be used within `runtime/secret.Do` for secure key material handling.

##### 4.1.5.2 SetKey Behavior

- `SetKey()` overwrites any existing key value.
- The key value is stored as a copy (for slice types, a deep copy is made).
- `SetKey()` does not update timestamps (`CreatedAt` remains unchanged).
- To update timestamps, create a new `SigningKey` instance or manually update the timestamp fields.
- **Validation Requirements**: `SetKey()` must validate that:
  - The key is valid for the signature type (`KeyType` matches the expected signature algorithm).
  - The key is not expired (if `ExpiresAt` is set, it must be in the future).
  - Returns `*PackageError` with `ErrTypeSignature` if validation fails.

##### 4.1.5.3 Error Conditions

- `GetKey()` returns `*PackageError` with `ErrTypeSignature` if:
  - The key is not set (Option is None) - error message should indicate "key not set"
  - The key is invalid (`IsValid()` returns false) - error message should indicate "key invalid"
  - The key is expired (`IsExpired()` returns true) - error message should indicate "key expired"
  - The key type does not match the expected type - error message should indicate "key type mismatch"

##### 4.1.5.4 Operation Requirements

All operations that use `SigningKey` must:

1. **Check key is set**: Before attempting to use the key, verify that it is set (Option is not None).
2. **Re-validate before use**: Re-validate that the key is:
   - Valid for the signature type (matches `KeyType`)
   - Not expired (if `ExpiresAt` is set, check it's in the future)
   - Passes `IsValid()` check
3. **Return PackageError on failure**: If any validation fails, return `*PackageError` with:
   - `ErrTypeSignature` as the error type
   - A descriptive error message indicating the specific reason for failure
   - Appropriate context in the error structure

This applies to:

- `Sign()` and `Verify()` methods of `SignatureStrategy` implementations
- Any other operations that access or use the key material

#### 4.1.6. SignatureStrategy Secure SigningKey Operations with Runtime/secret

Operations on `SigningKey` that handle private key material should use Go's `runtime/secret` package to protect sensitive cryptographic data in memory.
The `GetKey` method should use `runtime/secret.Do` when retrieving private signing keys to ensure that key material is handled securely and promptly erased from memory after use.
The `SetKey` method should wrap key assignment operations within the secret execution context to protect sensitive key data during storage operations.
All signing operations that use `SigningKey` instances containing private keys should execute within `runtime/secret.Do` to ensure comprehensive memory protection.
This applies to the `Sign` method of `SignatureStrategy` implementations that access private key material during signature generation.

### 4.2 IsValid and IsExpired Semantics

This section describes the semantics of IsValid and IsExpired methods.

#### 4.2.1 IsValid() Requirements

`IsValid()` returns `true` only when all of the following conditions are met:

1. **Key is set**: The key value must be set (Option is not None).
2. **KeyID is non-empty**: `KeyID` must be a non-empty string (length > 0).
3. **CreatedAt is non-zero**: `CreatedAt` must be a non-zero `time.Time` value (not the zero value).
4. **KeyType is valid**: `KeyType` must be a valid signature type (passes signature type validation).
5. **ExpiresAt is valid if set**: If `ExpiresAt` is not `nil`, it must be after `CreatedAt` (expiration time must be in the future relative to creation time).

#### 4.2.2 IsExpired() Requirements

`IsExpired()` returns `true` when:

1. **ExpiresAt is set**: `ExpiresAt` is not `nil`.
2. **Current time is past expiration**: The current time (`time.Now()`) is equal to or after `ExpiresAt`.

#### 4.2.3 Expiration Semantics

- **`ExpiresAt == nil` means never expires**: If `ExpiresAt` is `nil`, the key never expires and `IsExpired()` always returns `false`.
- **`ExpiresAt` before `CreatedAt` is invalid**: If `ExpiresAt` is set and is before or equal to `CreatedAt`, the key is considered invalid (`IsValid()` returns `false`).
- **Zero time handling**: If `CreatedAt` is the zero value, `IsValid()` returns `false` regardless of `ExpiresAt`.

#### 4.2.4 Validation Order

When checking key validity, operations should:

1. First check `IsValid()` to ensure all required fields are present and valid.
2. Then check `IsExpired()` to ensure the key is still within its validity period.
3. Both checks must pass for the key to be usable.

### 4.3 Generic Signature Configuration

The `SignatureConfig[T]` type extends the generic [Core Generic Types](api_generics.md#1-core-generic-types) for signature-specific settings.

#### 4.3.1 SignatureConfig Structure

```go
// SignatureConfig provides type-safe signature configuration
type SignatureConfig[T any] struct {
    *Config[T]  // See [Core Generic Types](api_generics.md#1-core-generic-types) for base configuration

    // Signature-specific settings
    SignatureType     Option[SignatureType]     // Signature algorithm type
    KeySize           Option[int]               // Key size in bits
    UseTimestamp      Option[bool]              // Include timestamp in signature
    IncludeMetadata   Option[bool]              // Include metadata in signature
    CompressionLevel  Option[int]               // Compression level for signed data
}
```

#### 4.3.2 SignatureConfigBuilder Structure

This section describes the SignatureConfigBuilder structure for building signature configurations.

##### 4.3.2.1 SignatureConfigBuilder Struct

```go
// SignatureConfigBuilder provides fluent configuration building for signatures
type SignatureConfigBuilder[T any] struct {
    config *SignatureConfig[T]
}
```

##### 4.3.2.2 NewSignatureConfigBuilder Function

```go
// NewSignatureConfigBuilder creates a new signature configuration builder.
func NewSignatureConfigBuilder[T any]() *SignatureConfigBuilder[T]
```

##### 4.3.2.3 SignatureConfigBuilder[T].WithSignatureType Method

```go
// WithSignatureType sets the signature type for the configuration.
func (b *SignatureConfigBuilder[T]) WithSignatureType(sigType SignatureType) *SignatureConfigBuilder[T]
```

##### 4.3.2.4 SignatureConfigBuilder[T].WithKeySize Method

```go
// WithKeySize sets the key size for the configuration.
func (b *SignatureConfigBuilder[T]) WithKeySize(size int) *SignatureConfigBuilder[T]
```

##### 4.3.2.5 SignatureConfigBuilder[T].WithTimestamp Method

```go
// WithTimestamp enables or disables timestamp inclusion for the configuration.
func (b *SignatureConfigBuilder[T]) WithTimestamp(useTimestamp bool) *SignatureConfigBuilder[T]
```

##### 4.3.2.6 SignatureConfigBuilder[T].WithMetadata Method

```go
// WithMetadata enables or disables metadata inclusion for the configuration.
func (b *SignatureConfigBuilder[T]) WithMetadata(includeMetadata bool) *SignatureConfigBuilder[T]
```

##### 4.3.2.7 SignatureConfigBuilder[T].Build Method

```go
// Build constructs and returns the final signature configuration.
func (b *SignatureConfigBuilder[T]) Build() *SignatureConfig[T]
```

### 4.4 Generic Signature Validation

The `SignatureValidator[T]` type extends the generic [Core Generic Types](api_generics.md#1-core-generic-types) for signature-specific validation.

#### 4.4.1 SignatureValidator Struct

```go
// SignatureValidator provides type-safe signature validation
type SignatureValidator[T any] struct {
    *Validator[T]  // See [Core Generic Types](api_generics.md#1-core-generic-types) for base validation
    signatureRules []SignatureValidationRule[T]
}
```

#### 4.4.2 SignatureValidationRule Alias

```go
// SignatureValidationRule is an alias for the generic ValidationRule
type SignatureValidationRule[T any] = ValidationRule[T]
```

#### 4.4.3 SignatureValidator[T].AddSignatureRule Method

```go
// AddSignatureRule adds a signature validation rule to the validator.
func (v *SignatureValidator[T]) AddSignatureRule(rule SignatureValidationRule[T])
```

#### 4.4.4 SignatureValidator[T].ValidateSignatureData Method

```go
// Returns *PackageError on failure
func (v *SignatureValidator[T]) ValidateSignatureData(ctx context.Context, data T) error
```

#### 4.4.5 SignatureValidator[T].ValidateSignatureKey Method

```go
// Returns *PackageError on failure
func (v *SignatureValidator[T]) ValidateSignatureKey(key SigningKey[T]) error
```

#### 4.4.6 SignatureValidator[T].ValidateSignatureFormat Method

```go
// Returns *PackageError on failure
func (v *SignatureValidator[T]) ValidateSignatureFormat(signature Signature[T]) error
```

## 5. Error Handling

This section describes error handling for signature operations.

### 5.1 Structured Error System

The NovusPack digital signature API uses the comprehensive structured error system defined in [api_core.md](api_core.md#10-structured-error-system).

**Error Types Used**: This API uses the following error types from the core error system:

- **ErrTypeSignature**: Digital signature validation, signing failures, signature format errors
- **ErrTypeValidation**: Input validation errors, invalid parameters, format errors
- **ErrTypeSecurity**: Security-related errors, access denied, authentication failures
- **ErrTypeUnsupported**: Unsupported features, versions, or operations
- **ErrTypeCorruption**: Data corruption, checksum failures, integrity violations

For complete error system documentation, see [Structured Error System](api_core.md#10-structured-error-system).

### 5.2 Signature Error Messages

Signature errors use descriptive error messages with type-safe `ErrorType` constants and typed context structures.
Error messages should be descriptive strings that clearly explain what went wrong.

Common signature error message patterns:

- Invalid signature format or data: `"invalid signature"`
- Signature not found at index: `"signature not found"`
- Signature validation failed: `"signature validation failed"`
- Unsupported signature algorithm: `"unsupported signature type"`
- Invalid signature data format: `"invalid signature data"`
- Signature exceeds maximum size: `"signature too large"`
- Signing key not found: `"signing key not found"`
- Invalid signing key format: `"invalid signing key"`
- Signature generation failed: `"signature generation failed"`
- Signature data corrupted: `"signature data corrupted"`

All signature errors use `NewPackageError[T]` or `WrapErrorWithContext[T]` with appropriate `ErrorType` constants (`ErrTypeSignature`, `ErrTypeValidation`, `ErrTypeSecurity`, `ErrTypeUnsupported`, `ErrTypeCorruption`) and typed context structures.

### 5.3 Signature-Specific Error Context Types

The signature API uses the following typed context structures for type-safe error handling:

#### 5.3.1 SignatureErrorContext Structure

```go
// Signature-specific error context types
type SignatureErrorContext struct {
    SignatureIndex int
    Algorithm      string
    Operation      string
}
```

#### 5.3.2 UnsupportedErrorContext Structure

```go
// UnsupportedErrorContext provides error context for unsupported signature type errors.
type UnsupportedErrorContext struct {
    SignatureType   uint32
    SupportedTypes  []uint32
    Operation       string
}
```

#### 5.3.3 SecurityErrorContext Structure

See [SecurityErrorContext Structure](api_basic_operations.md#203-securityerrorcontext-structure) for the complete structure definition.

#### 5.3.4 ValidationErrorContext Structure

```go
// ValidationErrorContext provides error context for signature validation errors.
type ValidationErrorContext struct {
    Field    string
    Value    string
    Expected string
}
```

All signature error-returning functions use `NewPackageError[T]` or `WrapErrorWithContext[T]` with these typed context structures.

For error handling patterns, examples, and helper functions, see:

- [Error Helper Functions](api_core.md#105-error-helper-functions) - Generic error context helpers
- [Error Handling Patterns](api_core.md#10-structured-error-system) - Error creation, inspection, and propagation patterns
- [Structured Error System](api_core.md#10-structured-error-system) - Complete error system documentation
