# NovusPack Technical Specifications - Security Validation API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Package Validation](#1-package-validation)
  - [1.1 Multiple Signature Validation (Incremental)](#11-multiple-signature-validation-incremental)
    - [1.1.1 Package ValidateIntegrity Method](#111-packagevalidateintegrity-method)
    - [1.1.2 Package GetSecurityStatus Method](#112-packagegetsecuritystatus-method)
    - [1.1.3 Deferred to v2](#113-deferred-to-v2)
  - [1.2 Incremental Validation Process](#12-incremental-validation-process)
- [2. SecurityStatus Structure](#2-securitystatus-structure)
  - [2.1 SecurityValidationResult struct](#21-securityvalidationresult-struct)
  - [2.2 SecurityStatus struct](#22-securitystatus-struct)
  - [2.3 SignatureValidationResult struct](#23-signaturevalidationresult-struct)
- [3. EncryptionType System](#3-encryptiontype-system)
  - [3.1 EncryptionType Definition](#31-encryptiontype-definition)
    - [3.1.1 EncryptionAlgorithm Type](#311-encryptionalgorithm-type)
    - [3.1.2 EncryptionType Alias](#312-encryptiontype-alias)
    - [3.1.3 EncryptionType Purpose](#313-encryptiontype-purpose)
    - [3.1.4 EncryptionType Values](#314-encryptiontype-values)
  - [3.2 EncryptionType Validation](#32-encryptiontype-validation)
    - [3.2.1 IsValidEncryptionType Function](#321-isvalidencryptiontype-function)
    - [3.2.2 GetEncryptionTypeName Function](#322-getencryptiontypename-function)
    - [3.2.3 Purpose (GetEncryptionTypeName)](#323-purpose-getencryptiontypename)
    - [3.2.4 Example Usage](#324-example-usage)
  - [3.3 On-Disk Mapping](#33-on-disk-mapping)
  - [3.4 Algorithm Selection Guidelines](#34-algorithm-selection-guidelines)
- [4. Generic Encryption Patterns](#4-generic-encryption-patterns)
  - [4.1 EncryptionStrategy Interface](#41-encryptionstrategy-interface)
    - [4.1.1 EncryptionStrategy Interface Definition](#411-encryptionstrategy-interface-definition)
    - [4.1.2 ByteEncryptionStrategy Interface](#412-byteencryptionstrategy-interface)
    - [4.1.3 EncryptionKey Structure](#413-encryptionkey-struct)
    - [4.1.4 GetKey and SetKey Behavior](#414-getkey-and-setkey-behavior)
    - [4.1.5 Secure EncryptionKey Operations with runtime/secret](#415-secure-encryptionkey-operations-with-runtimesecret)
  - [4.2 IsValid and IsExpired Semantics](#42-isvalid-and-isexpired-semantics)
    - [4.2.1 IsValid() Requirements](#421-isvalid-requirements)
    - [4.2.2 IsExpired() Requirements](#422-isexpired-requirements)
    - [4.2.3 Expiration Semantics](#423-expiration-semantics)
    - [4.2.4 Validation Order](#424-validation-order)
  - [4.3 Generic Encryption Configuration](#43-generic-encryption-configuration)
    - [4.3.1 EncryptionConfig Structure](#431-encryptionconfig-structure)
    - [4.3.2 EncryptionConfigBuilder Structure](#432-encryptionconfigbuilder-structure)
  - [4.4 Generic Encryption Validation](#44-generic-encryption-validation)
    - [4.4.1 EncryptionValidator Structure](#441-encryptionvalidator-struct)
    - [4.4.2 EncryptionValidator AddEncryptionRule Method](#443-encryptionvalidatortaddencryptionrule-method)
    - [4.4.3 EncryptionValidator ValidateEncryptionData Method](#444-encryptionvalidatortvalidateencryptiondata-method)
    - [4.4.4 EncryptionValidator ValidateDecryptionData Method](#445-encryptionvalidatortvalidatedecryptiondata-method)
    - [4.4.5 EncryptionValidator ValidateEncryptionKey Method](#446-encryptionvalidatortvalidateencryptionkey-method)
  - [4.5 File Encryption Operations](#45-file-encryption-operations)
    - [4.5.1 FileEncryptionHandler Interface](#451-fileencryptionhandler-interface)
    - [4.5.2 Built-in File Encryption Handlers](#452-built-in-file-encryption-handlers)
  - [4.6 Package File Encryption Operations](#46-package-file-encryption-operations)
    - [4.6.1 Package EncryptFile Method](#461-packageencryptfile-method)
    - [4.6.2 Package DecryptFile Method](#462-packagedecryptfile-method)
    - [4.6.3 Package ValidateFileEncryption Method](#463-packagevalidatefileencryption-method)
    - [4.6.4 Package GetFileEncryptionInfo Method](#464-packagegetfileencryptioninfo-method)
    - [4.6.5 Secure File Encryption Operations with runtime/secret](#465-secure-file-encryption-operations-with-runtimesecret)
  - [4.7 EncryptionErrorContext Structure](#47-encryptionerrorcontext-structure)
    - [4.7.1 Structure Definition](#471-encryptionerrorcontext-struct)
    - [4.7.2 Usage Examples](#472-usage-examples)
    - [4.7.3 Retrieving Context Example](#473-retrieving-context-example)
- [5. ML-KEM Key Structure and Operations](#5-ml-kem-key-structure-and-operations)
  - [5.1 ML-KEM Key Structure](#51-mlkemkey-struct)
    - [5.1.1 Purpose (MLKEMKey)](#511-purpose-mlkemkey)
    - [5.1.2 MLKEMKey Fields](#512-mlkemkey-fields)
  - [5.2 ML-KEM Encryption Operations](#52-ml-kem-encryption-operations)
    - [5.2.1 MLKEMKey Encrypt Method](#521-mlkemkeyencrypt-method)
    - [5.2.2 MLKEMKey Decrypt Method](#522-mlkemkeydecrypt-method)
    - [5.2.3 ML-KEM Parameters](#523-ml-kem-parameters)
    - [5.2.4 ML-KEM Returns](#524-ml-kem-returns)
    - [5.2.5 ML-KEM Error Conditions](#525-ml-kem-error-conditions)
    - [5.2.6 ML-KEM Example Usage](#526-ml-kem-example-usage)
    - [5.2.7 Secure Encryption Operations with runtime/secret](#527-secure-encryption-operations-with-runtimesecret)
  - [5.3 ML-KEM Key Information](#53-ml-kem-key-information)
    - [5.3.1 MLKEMKey GetPublicKey Method](#531-mlkemkeygetpublickey-method)
    - [5.3.2 MLKEMKey GetLevel Method](#532-mlkemkeygetlevel-method)
    - [5.3.3 MLKEMKey Clear Method](#533-mlkemkeyclear-method)
    - [5.3.4 On-Disk Storage](#534-on-disk-storage)
    - [5.3.5 Key Management Scope](#535-key-management-scope)
    - [5.3.6 GetPublicKey Returns](#536-getpublickey-returns)
    - [5.3.7 GetLevel Returns](#537-getlevel-returns)
    - [5.3.8 Clear Behavior](#538-clear-behavior)
    - [5.3.9 Secure Key Clearing with runtime/secret](#539-secure-key-clearing-with-runtimesecret)

---

## 0. Overview

This document defines the security validation API for the NovusPack system, including package validation, encryption, and security status reporting.

Signature verification and signature validation are deferred to v2.
V1 only enforces signed package immutability based on header fields that indicate signature presence (for example, Flags bit 0 "Has signatures" and `SignatureOffset > 0`).
There are no other immutability-related flags in the v1 header.
V1 does not validate signature contents.

There is no package-wide security level.
Security levels apply only to algorithms that define security levels (for example, ML-KEM level 1, 3, or 5) and are selected per encrypted FileEntry and key material.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Package Compression API](api_package_compression.md) - Generic strategy patterns and compression concurrency
- [File Management API](api_file_mgmt_index.md) - References this API for encryption patterns

## 1. Package Validation

**Note:** V1 focuses on format and integrity validation.
Signature validation is deferred to v2.

### 1.1 Multiple Signature Validation (Incremental)

Package validation methods:

- [Package.Validate](api_basic_operations.md#9-package-configuration) - Validates package format, structure, and integrity

#### 1.1.1 Package.ValidateIntegrity Method

```go
// ValidateIntegrity validates package integrity using checksums
// Returns *PackageError on failure
func (p *Package) ValidateIntegrity() error
```

#### 1.1.2 Package.GetSecurityStatus Method

```go
// GetSecurityStatus gets comprehensive security status
func (p *Package) GetSecurityStatus() SecurityStatus
```

#### 1.1.3. Deferred to V2

This functionality is deferred to version 2 of the API.

##### 1.1.3.1 Package.ValidateAllSignatures Method

```go
// ValidateAllSignatures validates all signatures in order
func (p *Package) ValidateAllSignatures() []SignatureValidationResult
```

##### 1.1.3.2 Package.ValidateSignatureType Method

```go
// ValidateSignatureType validates signatures of specific type
func (p *Package) ValidateSignatureType(signatureType uint32) []SignatureValidationResult
```

##### 1.1.3.3 Package.ValidateSignatureIndex Method

```go
// ValidateSignatureIndex validates signature by index
// Returns *PackageError on failure
func (p *Package) ValidateSignatureIndex(index int) (SignatureValidationResult, error)
```

##### 1.1.3.4 Package.ValidateSignatureChain Method

```go
// ValidateSignatureChain validates signature chain integrity
func (p *Package) ValidateSignatureChain() []SignatureValidationResult
```

### 1.2 Incremental Validation Process

- Signature validation is deferred to v2.
- Each signature validates complete package state up to its position
- Chain validation ensures no signatures were removed or modified
- Results provide detailed validation status for each signature

## 2. SecurityStatus Structure

This section describes the SecurityStatus structure used for security validation results.

### 2.1. SecurityValidationResult Struct

**Note:** Signature-related fields are deferred to v2.
V1 returns signature-related fields as zero values.
There is no package-wide security level in v1.

- **SignatureCount (int)**: Number of signatures in package
- **ValidSignatures (int)**: Number of valid signatures
- **TrustedSignatures (int)**: Number of trusted signatures
- **SignatureResults ([]SignatureValidationResult)**: Individual results
- **HasChecksums (bool)**: Checksums present
- **ChecksumsValid (bool)**: Checksums valid
- **ValidationErrors ([]string)**: Validation errors

### 2.2. SecurityStatus Struct

See [SecurityStatus Structure](api_metadata.md#73-securitystatus-structure) for the complete structure definition.

### 2.3. SignatureValidationResult Struct

**Note:** This structure is v2-only.

```go
// SignatureValidationResult provides information about individual signature validation results.
type SignatureValidationResult struct {
    Index       int       // Signature index in the package
    Type        uint32    // Signature type identifier
    Valid       bool      // Whether signature is valid
    Trusted     bool      // Whether signature is trusted
    Error       string    // Error message if validation failed
    Timestamp   uint32    // When signature was created
    PublicKey   []byte    // Public key used for validation (if available)
}
```

## 3. EncryptionType System

This section describes the EncryptionType system for specifying encryption algorithms.

### 3.1 EncryptionType Definition

This section describes the EncryptionType definition and structure.

#### 3.1.1 EncryptionAlgorithm Type

```go
// EncryptionAlgorithm enumeration
type EncryptionAlgorithm int

const (
    EncryptionAlgorithmNone EncryptionAlgorithm = iota
    EncryptionAlgorithmAES256GCM
    EncryptionAlgorithmChaCha20Poly1305
    EncryptionAlgorithmMLKEM512
    EncryptionAlgorithmMLKEM768
    EncryptionAlgorithmMLKEM1024
)
```

#### 3.1.2 EncryptionType Alias

```go
// EncryptionType is a v1 alias of EncryptionAlgorithm.
// This exists for compatibility with existing APIs.
type EncryptionType = EncryptionAlgorithm

const (
    EncryptionNone EncryptionType = EncryptionAlgorithmNone
    EncryptionAES256GCM EncryptionType = EncryptionAlgorithmAES256GCM
    EncryptionChaCha20Poly1305 EncryptionType = EncryptionAlgorithmChaCha20Poly1305
    EncryptionMLKEM512 EncryptionType = EncryptionAlgorithmMLKEM512
    EncryptionMLKEM768 EncryptionType = EncryptionAlgorithmMLKEM768
    EncryptionMLKEM1024 EncryptionType = EncryptionAlgorithmMLKEM1024
)
```

#### 3.1.3 EncryptionType Purpose

Defines the available encryption algorithms for file encryption.

#### 3.1.4 EncryptionType Values

- `EncryptionAlgorithmNone`: No encryption (default)
- `EncryptionAlgorithmAES256GCM`: AES-256-GCM symmetric encryption
- `EncryptionAlgorithmChaCha20Poly1305`: ChaCha20-Poly1305 symmetric encryption
- `EncryptionAlgorithmMLKEM512`: ML-KEM-512 post-quantum encryption (security level 1)
- `EncryptionAlgorithmMLKEM768`: ML-KEM-768 post-quantum encryption (security level 3)
- `EncryptionAlgorithmMLKEM1024`: ML-KEM-1024 post-quantum encryption (security level 5)

### 3.2 EncryptionType Validation

This section describes validation functions for EncryptionType values.

#### 3.2.1 IsValidEncryptionType Function

```go
// IsValidEncryptionType checks if the encryption type is valid
func IsValidEncryptionType(encType EncryptionType) bool
```

#### 3.2.2 GetEncryptionTypeName Function

```go
// GetEncryptionTypeName returns the human-readable name of the encryption type
func GetEncryptionTypeName(encType EncryptionType) string
```

#### 3.2.3 Purpose (GetEncryptionTypeName)

Provides validation and naming utilities for encryption types.

#### 3.2.4 Example Usage

```go
if !IsValidEncryptionType(encType) {
    return fmt.Errorf("invalid encryption type: %d", encType)
}

name := GetEncryptionTypeName(encType)
fmt.Printf("Using encryption: %s\n", name)
```

### 3.3 On-Disk Mapping

The on-disk representation of encryption is defined by the NovusPack file format.
See [Compression and Encryption Types](package_file_format.md#4113-compression-and-encryption-types) for the `EncryptionType` field and its encoded values.
See [Encrypted File Data Framing](package_file_format.md#4114-encrypted-file-data-framing) for the per-file ciphertext encoding (nonces, encapsulation, and ciphertext layout).

In v1:

- `EncryptionAES256GCM` maps to the on-disk AES-256-GCM value.
- `EncryptionChaCha20Poly1305` maps to the on-disk ChaCha20-Poly1305 value.
- `EncryptionMLKEM512`, `EncryptionMLKEM768`, and `EncryptionMLKEM1024` map to the on-disk quantum-safe encryption value.

ML-KEM variants are represented in memory by the selected algorithm constant.
The file format stores quantum-safe encryption as a single on-disk value.
In v1, the specific ML-KEM variant is derived from the key material provided to the API (for example, the `Level` carried by `MLKEMKey`) and is not encoded as a distinct `EncryptionType` value.

### 3.4 Algorithm Selection Guidelines

This section provides non-normative guidance for selecting encryption algorithms.

- Prefer `EncryptionAES256GCM` for most desktop and server environments.
- Prefer `EncryptionChaCha20Poly1305` for mobile, embedded, and ARM environments.
- Prefer `EncryptionMLKEM768` or `EncryptionMLKEM1024` when post-quantum security is required.

## 4. Generic Encryption Patterns

The security API provides generic encryption patterns that extend the generic configuration patterns defined in [api_generics.md](api_generics.md#1-core-generic-types) for type-safe encryption operations.

### 4.1 EncryptionStrategy Interface

The `EncryptionStrategy[T]` interface extends the generic [Core Generic Types](api_generics.md#1-core-generic-types) pattern for encryption-specific operations.
`EncryptionStrategy[T]` embeds `Strategy[T, T]` where both input and output are the same type.
The `Process` method from `Strategy[T, T]` can be used for encryption operations, while `Encrypt` and `Decrypt` provide more specific encryption/decryption methods with key management.

#### 4.1.1 EncryptionStrategy Interface Definition

```go
// EncryptionStrategy extends Strategy[T, T] for encryption operations
// Both input and output are the same type T
// The Strategy.Type() method returns "encryption" as the category
type EncryptionStrategy[T any] interface {
    Strategy[T, T]  // Extends the generic Strategy interface

    Encrypt(ctx context.Context, data T, key EncryptionKey[T]) (T, error)
    Decrypt(ctx context.Context, data T, key EncryptionKey[T]) (T, error)
    EncryptionType() EncryptionType  // Returns the specific encryption algorithm type
    Name() string
    KeySize() int
    ValidateKey(ctx context.Context, key EncryptionKey[T]) error
}
```

#### 4.1.2 ByteEncryptionStrategy Interface

```go
// ByteEncryptionStrategy is the concrete implementation for []byte data
type ByteEncryptionStrategy interface {
    EncryptionStrategy[[]byte]
}
```

#### 4.1.3 EncryptionKey Struct

This section describes the EncryptionKey structure for managing encryption keys.

##### 4.1.3.1 EncryptionKey Struct Type Definition

```go
// EncryptionKey provides type-safe key management
// Uses Option[T] internally for type-safe key storage
type EncryptionKey[T any] struct {
    *Option[T]  // See [Option Type](api_generics.md#11-option-type) for details
    KeyType    EncryptionType
    KeyID      string
    CreatedAt  time.Time
    ExpiresAt  *time.Time
}
```

##### 4.1.3.2 NewEncryptionKey Function

```go
// NewEncryptionKey creates a new encryption key with the specified type, ID, and key material.
func NewEncryptionKey[T any](keyType EncryptionType, keyID string, key T) *EncryptionKey[T]
```

##### 4.1.3.3 EncryptionKey[T].GetKey Method

```go
// GetKey returns the encryption key material.
func (k *EncryptionKey[T]) GetKey() (T, error)
```

##### 4.1.3.4 EncryptionKey[T].SetKey Method

```go
// SetKey sets the encryption key material.
func (k *EncryptionKey[T]) SetKey(key T)
```

##### 4.1.3.5 EncryptionKey[T].IsValid Method

```go
// IsValid returns true if the encryption key is valid.
func (k *EncryptionKey[T]) IsValid() bool
```

##### 4.1.3.6 EncryptionKey[T].IsExpired Method

```go
// IsExpired returns true if the encryption key has expired.
func (k *EncryptionKey[T]) IsExpired() bool
```

#### 4.1.4 GetKey and SetKey Behavior

This section describes the behavior of GetKey and SetKey methods.

##### 4.1.4.1 GetKey Behavior

- `GetKey()` returns a **copy** of the key value, not a reference to the original.
- This ensures that modifications to the returned value do not affect the stored key.
- For slice types (`[]byte`, etc.), a deep copy is performed to prevent accidental mutation.
- Returns an error if the key is missing, invalid, or expired.
- The returned copy SHOULD be used within `runtime/secret.Do` for secure key material handling.
- The API cannot enforce `runtime/secret.Do` usage outside the NovusPack API boundary.

##### 4.1.4.2 SetKey Behavior

- `SetKey()` overwrites any existing key value.
- The key value is stored as a copy (for slice types, a deep copy is made).
- `SetKey()` does not update timestamps (`CreatedAt` remains unchanged).
- To update timestamps, create a new `EncryptionKey` instance or manually update the timestamp fields.
- **Validation Requirements**: `SetKey()` must validate that:
  - The key is valid for the encryption type (`KeyType` matches the expected encryption algorithm).
  - The key is not expired (if `ExpiresAt` is set, it must be in the future).
  - Returns `*PackageError` with `ErrTypeEncryption` if validation fails.

##### 4.1.4.3 EncryptionKey Error Conditions

- `GetKey()` returns `*PackageError` with `ErrTypeEncryption` if:
  - The key is not set (Option is None) - error message should indicate "key not set"
  - The key is invalid (`IsValid()` returns false) - error message should indicate "key invalid"
  - The key is expired (`IsExpired()` returns true) - error message should indicate "key expired"
  - The key type does not match the expected type - error message should indicate "key type mismatch"

##### 4.1.4.4 Operation Requirements

All operations that use `EncryptionKey` must:

1. **Check key is set**: Before attempting to use the key, verify that it is set (Option is not None).
2. **Re-validate before use**: Re-validate that the key is:
   - Valid for the encryption type (matches `KeyType`)
   - Not expired (if `ExpiresAt` is set, check it's in the future)
   - Passes `IsValid()` check
3. **Return PackageError on failure**: If any validation fails, return `*PackageError` with:
   - `ErrTypeEncryption` as the error type
   - A descriptive error message indicating the specific reason for failure
   - Appropriate context in the error structure

This applies to:

- `Encrypt()` and `Decrypt()` methods of `EncryptionStrategy` implementations
- Any other operations that access or use the key material

#### 4.1.5. Secure EncryptionKey Operations with Runtime/secret

This section describes secure operations for EncryptionKey using runtime/secret.

##### 4.1.5.1 Private Key Material Policy

For API simplicity and safety, ALL keys stored in `EncryptionKey[T]` MUST be treated as private key material and handled with `runtime/secret.Do`.
This approach:

- **Simplifies code paths**: No need to distinguish between public and private keys at runtime
- **Provides defense in depth**: Even if a key is technically "public", protecting it doesn't hurt
- **Makes the API contract clear**: All keys receive the same security guarantees

ALL keys (regardless of `T` type, `KeyType`, or `EncryptionType`) MUST be handled with `runtime/secret.Do`.
This is a blanket policy with no exceptions for public keys or non-sensitive data.
This ensures consistent security handling across all key operations.

##### 4.1.5.2 Supported Key Types

The following `T` types are supported for storage in `EncryptionKey[T]`:

- `[]byte` - Raw byte slices for symmetric keys
- `*rsa.PrivateKey` - RSA private keys
- `*ecdsa.PrivateKey` - ECDSA private keys
- `*ed25519.PrivateKey` - Ed25519 private keys
- Other cryptographic key types as specified by the encryption strategy implementation

##### 4.1.5.3 Method Requirements

Since ALL keys are treated as private key material, the following requirements apply.

##### 4.1.5.4 Internal API Requirements (MUST)

- `GetKey()` MUST execute key retrieval within `runtime/secret.Do` for all keys.
- `SetKey()` MUST execute key assignment within `runtime/secret.Do` for all keys.
- `Encrypt()` and `Decrypt()` MUST execute within `runtime/secret.Do` when accessing keys.
- `Clear()` MUST execute key clearing within `runtime/secret.Do` for all keys.

##### 4.1.5.5 External User Code Recommendations (SHOULD)

- Users SHOULD use `runtime/secret.Do` when handling values returned by `GetKey()`.
- Users SHOULD NOT store key material outside the secret execution context.

##### 4.1.5.6 Key Material Lifetime

Key material has the following lifetime guarantees:

- The returned copy from `GetKey()` is valid only within the `runtime/secret.Do` context where it was retrieved.
- Key material is automatically erased from memory after the `runtime/secret.Do` function completes.
- Callers SHOULD use the returned key within `runtime/secret.Do` for any operations that access the key material.
- Callers SHOULD NOT store the key copy outside the secret execution context.

### 4.2 IsValid and IsExpired Semantics

This section describes the semantics of IsValid and IsExpired methods.

#### 4.2.1 IsValid() Requirements

`IsValid()` returns `true` only when all of the following conditions are met:

1. **Key is set**: The key value must be set (Option is not None).
2. **KeyID is non-empty**: `KeyID` must be a non-empty string (length > 0).
3. **CreatedAt is non-zero**: `CreatedAt` must be a non-zero `time.Time` value (not the zero value).
4. **KeyType is valid**: `KeyType` must be a valid encryption type (passes `IsValidEncryptionType()`).
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

### 4.3 Generic Encryption Configuration

The `EncryptionConfig[T]` type extends the generic [Configuration Patterns](api_generics.md#110-generic-configuration-patterns) for encryption-specific settings.

#### 4.3.1 EncryptionConfig Structure

```go
// EncryptionConfig provides type-safe encryption configuration
type EncryptionConfig[T any] struct {
    *Config[T]  // See [Core Generic Types](api_generics.md#1-core-generic-types) for base configuration

    // Encryption-specific settings
    EncryptionType     Option[EncryptionType]     // Encryption algorithm type
    KeySize           Option[int]                 // Key size in bits
    UseRandomIV       Option[bool]                // Use random initialization vector
    AuthenticationTag Option[bool]                // Include authentication tag
    CompressionLevel  Option[int]                 // Compression level for encrypted data
}
```

#### 4.3.2 EncryptionConfigBuilder Structure

This section describes the EncryptionConfigBuilder structure for building encryption configurations.

##### 4.3.2.1 EncryptionConfigBuilder Struct

```go
// EncryptionConfigBuilder provides fluent configuration building for encryption
type EncryptionConfigBuilder[T any] struct {
    config *EncryptionConfig[T]
}
```

##### 4.3.2.2 NewEncryptionConfigBuilder Function

```go
// NewEncryptionConfigBuilder creates a new encryption configuration builder.
func NewEncryptionConfigBuilder[T any]() *EncryptionConfigBuilder[T]
```

##### 4.3.2.3 EncryptionConfigBuilder[T].WithEncryptionType Method

```go
// WithEncryptionType sets the encryption type for the configuration.
func (b *EncryptionConfigBuilder[T]) WithEncryptionType(encType EncryptionType) *EncryptionConfigBuilder[T]
```

##### 4.3.2.4 EncryptionConfigBuilder[T].WithKeySize Method

```go
// WithKeySize sets the key size for the configuration.
func (b *EncryptionConfigBuilder[T]) WithKeySize(size int) *EncryptionConfigBuilder[T]
```

##### 4.3.2.5 EncryptionConfigBuilder[T].WithRandomIV Method

```go
// WithRandomIV enables or disables random IV generation for the configuration.
func (b *EncryptionConfigBuilder[T]) WithRandomIV(useRandom bool) *EncryptionConfigBuilder[T]
```

##### 4.3.2.6 EncryptionConfigBuilder[T].WithAuthenticationTag Method

```go
// WithAuthenticationTag enables or disables authentication tag generation for the configuration.
func (b *EncryptionConfigBuilder[T]) WithAuthenticationTag(useAuth bool) *EncryptionConfigBuilder[T]
```

##### 4.3.2.7 EncryptionConfigBuilder[T].Build Method

```go
// Build constructs and returns the final encryption configuration.
func (b *EncryptionConfigBuilder[T]) Build() *EncryptionConfig[T]
```

### 4.4 Generic Encryption Validation

The `EncryptionValidator[T]` type extends the generic [Core Generic Types](api_generics.md#1-core-generic-types) for encryption-specific validation.

#### 4.4.1 EncryptionValidator Struct

```go
// EncryptionValidator provides type-safe encryption validation
type EncryptionValidator[T any] struct {
    *Validator[T]  // See [Core Generic Types](api_generics.md#1-core-generic-types) for base validation
    encryptionRules []EncryptionValidationRule[T]
}
```

#### 4.4.2 EncryptionValidationRule Alias

```go
// EncryptionValidationRule is an alias for the generic ValidationRule
type EncryptionValidationRule[T any] = ValidationRule[T]
```

#### 4.4.3 EncryptionValidator[T].AddEncryptionRule Method

```go
// AddEncryptionRule adds an encryption validation rule to the validator.
func (v *EncryptionValidator[T]) AddEncryptionRule(rule EncryptionValidationRule[T])
```

#### 4.4.4 EncryptionValidator[T].ValidateEncryptionData Method

```go
// Returns *PackageError on failure
func (v *EncryptionValidator[T]) ValidateEncryptionData(ctx context.Context, data T) error
```

#### 4.4.5 EncryptionValidator[T].ValidateDecryptionData Method

```go
// Returns *PackageError on failure
func (v *EncryptionValidator[T]) ValidateDecryptionData(ctx context.Context, data T) error
```

#### 4.4.6 EncryptionValidator[T].ValidateEncryptionKey Method

```go
// Returns *PackageError on failure
func (v *EncryptionValidator[T]) ValidateEncryptionKey(key EncryptionKey[T]) error
```

### 4.5 File Encryption Operations

This section describes file-level encryption operations.

#### 4.5.1 FileEncryptionHandler Interface

```go
// FileEncryptionHandler provides file-specific encryption operations
type FileEncryptionHandler[T any] interface {
    EncryptFile(ctx context.Context, filePath string, data T, key EncryptionKey[T]) error
    DecryptFile(ctx context.Context, filePath string, key EncryptionKey[T]) (T, error)
    ValidateFileEncryption(ctx context.Context, filePath string) error
}
```

#### 4.5.2 Built-in File Encryption Handlers

This section describes built-in file encryption handler implementations.

##### 4.5.2.1 AES256GCMFileHandler Structure

```go
// Built-in file encryption handlers
type AES256GCMFileHandler struct { ... }
```

##### 4.5.2.2 ChaCha20Poly1305FileHandler Structure

```go
// ChaCha20Poly1305FileHandler provides file encryption using ChaCha20-Poly1305 algorithm.
type ChaCha20Poly1305FileHandler struct { ... }
```

##### 4.5.2.3 MLKEMFileHandler Structure

```go
// MLKEMFileHandler provides file encryption using ML-KEM (post-quantum) algorithm.
type MLKEMFileHandler struct { ... }
```

### 4.6 Package File Encryption Operations

This section describes package-level file encryption operations.

#### 4.6.1 Package.EncryptFile Method

```go
// EncryptFile encrypts a file using the security API's file encryption patterns
// Returns *PackageError on failure
func (p *Package) EncryptFile[T any](ctx context.Context, path string, data T, handler FileEncryptionHandler[T], key EncryptionKey[T]) error
```

#### 4.6.2 Package.DecryptFile Method

```go
// DecryptFile decrypts a file using the security API's file encryption patterns
// Returns *PackageError on failure
func (p *Package) DecryptFile[T any](ctx context.Context, path string, handler FileEncryptionHandler[T], key EncryptionKey[T]) (T, error)
```

#### 4.6.3 Package.ValidateFileEncryption Method

```go
// ValidateFileEncryption validates file encryption using the security API's file encryption patterns
// Returns *PackageError on failure
func (p *Package) ValidateFileEncryption[T any](ctx context.Context, path string, handler FileEncryptionHandler[T]) error
```

#### 4.6.4 Package.GetFileEncryptionInfo Method

```go
// GetFileEncryptionInfo gets encryption information for a file using the security API's patterns
func (p *Package) GetFileEncryptionInfo[T any](path string) (*EncryptionConfig[T], error)
```

**Type Constraints**: The type parameter `T` in these functions is typically `[]byte` for file data operations, but can be any type that the `FileEncryptionHandler[T]` supports.
For most use cases, `T` should be `[]byte` to work with file content directly.

**Error Handling**: All encryption operations return errors using `NewPackageError` or `WrapErrorWithContext` with `EncryptionErrorContext` for type-safe error handling.
See [EncryptionErrorContext Structure](#47-encryptionerrorcontext-structure) for details.

#### 4.6.5. Secure File Encryption Operations with Runtime/secret

**MUST Requirements**: File encryption and decryption operations that use keys MUST execute within Go's `runtime/secret.Do` function to protect sensitive cryptographic material.

- `DecryptFile` method MUST wrap decryption operations within the secret execution context to ensure that keys and decrypted plaintext are handled securely
- `EncryptFile` method MUST wrap encryption operations within the secret execution context when accessing keys
- `FileEncryptionHandler` implementations MUST use `runtime/secret.Do` when accessing key material during encryption and decryption operations
- This ensures that sensitive key data and decrypted content are promptly erased from memory registers and stack frames, reducing exposure to memory analysis attacks

### 4.7 EncryptionErrorContext Structure

The security API defines a typed error context structure for type-safe error handling:

#### 4.7.1 EncryptionErrorContext Struct

```go
// EncryptionErrorContext provides type-safe error context for encryption operations
type EncryptionErrorContext struct {
    Path          string         // File path that caused the error
    Operation     string         // Operation name (e.g., "EncryptFile", "DecryptFile", "ValidateFileEncryption")
    EncryptionType EncryptionType // Encryption algorithm type
    KeyID         string         // Encryption key ID (if applicable)
    KeySize       int            // Key size in bits (if applicable)
    ErrorStage    string         // Stage where error occurred (e.g., "key_generation", "encryption", "decryption")
}
```

**Usage**: All encryption error-returning functions use `EncryptionErrorContext` with `NewPackageError` or `WrapErrorWithContext`:

#### 4.7.2 Usage Examples

```go
// Example: Encryption failure with typed context
err := NewPackageError(ErrTypeEncryption, "encryption failed", nil, EncryptionErrorContext{
    Path:          "/path/to/file",
    Operation:     "EncryptFile",
    EncryptionType: EncryptionAES256GCM,
    KeyID:         "key-12345",
    KeySize:       256,
    ErrorStage:    "encryption",
})

// Example: Key validation error with typed context
err := WrapErrorWithContext(validationErr, ErrTypeSecurity, "invalid encryption key", EncryptionErrorContext{
    Path:          "/path/to/file",
    Operation:     "EncryptFile",
    EncryptionType: EncryptionMLKEM768,
    KeyID:         "key-67890",
    KeySize:       0,
    ErrorStage:    "key_validation",
})
```

**Retrieving Context**: Use `GetErrorContext` to retrieve typed error context:

#### 4.7.3 Retrieving Context Example

```go
if ctx, ok := GetErrorContext[EncryptionErrorContext](err, "_typed_context"); ok {
    log.Printf("Encryption error: Path=%s, Operation=%s, Type=%v, KeyID=%s",
        ctx.Path, ctx.Operation, ctx.EncryptionType, ctx.KeyID)
}
```

## 5. ML-KEM Key Structure and Operations

**Note**: NovusPack does not provide key generation functions.
Users are responsible for generating cryptographic keys using appropriate external libraries (e.g., `crypto/aes`, `crypto/rand`, Cloudflare CIRCL for ML-KEM).
NovusPack accepts pre-generated keys wrapped in the `EncryptionKey[T]` structure.

### 5.1 MLKEMKey Struct

```go
// ML-KEM Key Structure
type MLKEMKey struct {
    PublicKey  []byte  // ML-KEM public key data
    PrivateKey []byte  // ML-KEM private key data
    Level      int     // Security level (1, 3, or 5)
}
```

#### 5.1.1 Purpose (MLKEMKey)

Represents an ML-KEM key pair with associated security level.

#### 5.1.2 MLKEMKey Fields

- `PublicKey`: Public key data for encryption
- `PrivateKey`: Private key data for decryption
- `Level`: Security level (1, 3, or 5) mapped to ML-KEM-512, ML-KEM-768, and ML-KEM-1024

### 5.2 ML-KEM Encryption Operations

Performs encryption and decryption operations using ML-KEM keys.

#### 5.2.1 MLKEMKey.Encrypt Method

```go
// Encrypt encrypts plaintext using ML-KEM key
func (k *MLKEMKey) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error)
```

#### 5.2.2 MLKEMKey.Decrypt Method

```go
// Decrypt decrypts ciphertext using ML-KEM key
func (k *MLKEMKey) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error)
```

#### 5.2.3 ML-KEM Parameters

- `ctx`: Context for cancellation and timeout handling
- `plaintext`: Data to encrypt (for Encrypt method)
- `ciphertext`: Encrypted data to decrypt (for Decrypt method)

#### 5.2.4 ML-KEM Returns

Encrypted data (for Encrypt) or decrypted data (for Decrypt)

#### 5.2.5 ML-KEM Error Conditions

- `ErrEncryptionFailed`: Failed to encrypt data
- `ErrDecryptionFailed`: Failed to decrypt data
- `ErrInvalidKey`: Key is invalid or corrupted
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 5.2.6 ML-KEM Example Usage

```go
// Encrypt data
ciphertext, err := key.Encrypt(ctx, []byte("Sensitive data"))
if err != nil {
    return fmt.Errorf("encryption failed: %w", err)
}

// Decrypt data
plaintext, err := key.Decrypt(ctx, ciphertext)
if err != nil {
    return fmt.Errorf("decryption failed: %w", err)
}
```

#### 5.2.7. Secure Encryption Operations with Runtime/secret

**MUST Requirements**: Encryption and decryption operations that access key material MUST use Go's `runtime/secret` package to protect sensitive cryptographic data in memory.

- `Decrypt` method MUST wrap operations that use keys within `runtime/secret.Do` to ensure that key data and decrypted plaintext are handled securely
- `Encrypt` method MUST wrap operations that use keys within `runtime/secret.Do` to ensure that key data is handled securely
- This prevents sensitive data from persisting in memory registers, stack frames, or heap allocations longer than necessary
- Implementations MUST ensure that all operations involving key material are executed within the secret execution context provided by `runtime/secret.Do`

### 5.3 ML-KEM Key Information

This section describes ML-KEM key information and operations.

#### 5.3.1 MLKEMKey.GetPublicKey Method

```go
// GetPublicKey returns the public key data
func (k *MLKEMKey) GetPublicKey() []byte
```

#### 5.3.2 MLKEMKey.GetLevel Method

```go
// GetLevel returns the security level of the key
func (k *MLKEMKey) GetLevel() int
```

#### 5.3.3 MLKEMKey.Clear Method

```go
// Clear clears sensitive key data from memory
func (k *MLKEMKey) Clear()
```

#### 5.3.4 On-Disk Storage

The ML-KEM public key is not stored in the on-disk package.
The package stores only per-file `EncryptionType` values and encrypted file data.
Key management is out of scope for the package file format.
Callers must provide the appropriate key material when writing and reading encrypted files.

If you need to persist key identity (for example, a key fingerprint or a key ID), that must be handled out-of-band in v1 (for example, in application metadata).

#### 5.3.5 Key Management Scope

Provides access to key information and secure key cleanup.

#### 5.3.6 GetPublicKey Returns

Copy of the public key data (safe to share)

#### 5.3.7 GetLevel Returns

Security level (1, 3, or 5)

#### 5.3.8 Clear Behavior

- Securely zeroes out all sensitive key material in memory
- Should be called using `defer key.Clear()` immediately after key creation
- Uses `runtime/secret.Do` to ensure secure cleanup

#### 5.3.9. Secure Key Clearing with Runtime/secret

**MUST Requirements**: Key clearing operations MUST use Go's `runtime/secret` package to ensure that sensitive data is securely erased from memory.

- `Clear()` method MUST wrap key clearing operations within `runtime/secret.Do` to ensure that key data is securely zeroed
- This provides defense-in-depth by ensuring that even the cleanup operations are protected from memory analysis attacks
- Implementations MUST wrap key clearing operations within the secret execution context to maximize protection of sensitive cryptographic material
