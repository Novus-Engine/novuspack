# NovusPack Technical Specifications - Security Validation API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Package Validation](#1-package-validation)
  - [1.1 Multiple Signature Validation (Incremental)](#11-multiple-signature-validation-incremental)
  - [1.2 Incremental Validation Process](#12-incremental-validation-process)
- [2. Security Status Structure](#2-security-status-structure)
  - [2.1 SecurityValidationResult struct](#21-securityvalidationresult-struct)
  - [2.2 SecurityStatus struct](#22-securitystatus-struct)
  - [2.3 SignatureValidationResult struct](#23-signaturevalidationresult-struct)
  - [2.4 SecurityLevel enum](#24-securitylevel-enum)
- [3. Encryption Type System](#3-encryption-type-system)
  - [3.1 Encryption Type Definition](#31-encryption-type-definition)
  - [3.2 Encryption Type Validation](#32-encryption-type-validation)
- [4. Generic Encryption Patterns](#4-generic-encryption-patterns)
  - [4.1 Generic Encryption Strategy Interface](#41-generic-encryption-strategy-interface)
  - [4.2 Generic Encryption Configuration](#42-generic-encryption-configuration)
  - [4.3 Generic Encryption Validation](#43-generic-encryption-validation)
  - [4.4 File Encryption Operations](#44-file-encryption-operations)
  - [4.5 Package File Encryption Operations](#45-package-file-encryption-operations)
- [5. ML-KEM Key Management](#5-ml-kem-key-management)
  - [5.1 ML-KEM Key Structure](#51-ml-kem-key-structure)
  - [5.2 ML-KEM Key Generation](#52-ml-kem-key-generation)
  - [5.3 ML-KEM Encryption Operations](#53-ml-kem-encryption-operations)
  - [5.4 ML-KEM Key Information](#54-ml-kem-key-information)

---

## 0. Overview

This document defines the security validation API for the NovusPack system, including package validation, signature verification, and security status reporting.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Package Compression API](api_package_compression.md) - Generic strategy patterns and compression concurrency
- [File Management API](api_file_management.md) - References this API for encryption patterns

## 1. Package Validation

**Note:** This API implements the incremental signature validation approach as defined in [Package File Format - Digital Signatures](package_file_format.md#7-digital-signatures-section-optional).

### 1.1 Multiple Signature Validation (Incremental)

```go
// Validate validates package format, structure, and integrity
func (p *Package) Validate() error

// ValidateAllSignatures validates all signatures in order
func (p *Package) ValidateAllSignatures() []SignatureValidationResult

// ValidateSignatureType validates signatures of specific type
func (p *Package) ValidateSignatureType(signatureType uint32) []SignatureValidationResult

// ValidateSignatureIndex validates signature by index
func (p *Package) ValidateSignatureIndex(index int) (SignatureValidationResult, error)

// ValidateSignatureChain validates signature chain integrity
func (p *Package) ValidateSignatureChain() []SignatureValidationResult

// ValidateIntegrity validates package integrity using checksums
func (p *Package) ValidateIntegrity() error

// GetSecurityStatus gets comprehensive security status
func (p *Package) GetSecurityStatus() SecurityStatus
```

### 1.2 Incremental Validation Process

- Each signature validates complete package state up to its position
- Chain validation ensures no signatures were removed or modified
- Results provide detailed validation status for each signature

## 2. Security Status Structure

### 2.1 SecurityValidationResult struct

- **SignatureCount (int)**: Number of signatures in package
- **ValidSignatures (int)**: Number of valid signatures
- **TrustedSignatures (int)**: Number of trusted signatures
- **SignatureResults ([]SignatureValidationResult)**: Individual results
- **HasChecksums (bool)**: Checksums present
- **ChecksumsValid (bool)**: Checksums valid
- **SecurityLevel (SecurityLevel)**: Overall security level
- **ValidationErrors ([]string)**: Validation errors

### 2.2 SecurityStatus struct

```go
type SecurityStatus struct {
    SignatureCount      int                           // Number of signatures
    ValidSignatures     int                           // Number of valid signatures
    TrustedSignatures   int                           // Number of trusted signatures
    SignatureResults    []SignatureValidationResult   // Individual results
    HasChecksums        bool                          // Checksums present
    ChecksumsValid      bool                          // Checksums valid
    SecurityLevel       SecurityLevel                 // Overall security level
    ValidationErrors    []string                      // Validation errors
}
```

### 2.3 SignatureValidationResult struct

```go
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

### 2.4 SecurityLevel enum

```go
type SecurityLevel int

const (
    SecurityLevelNone SecurityLevel = iota
    SecurityLevelLow
    SecurityLevelMedium
    SecurityLevelHigh
    SecurityLevelMaximum
)
```

## 3. Encryption Type System

### 3.1 Encryption Type Definition

```go
// EncryptionType enumeration
type EncryptionType int

const (
    EncryptionNone EncryptionType = iota
    EncryptionAES256GCM
    EncryptionMLKEM
)
```

#### 3.1.1 Purpose

Defines the available encryption algorithms for file encryption.

#### 3.1.2 Values

- `EncryptionNone`: No encryption (default)
- `EncryptionAES256GCM`: AES-256-GCM symmetric encryption
- `EncryptionMLKEM`: ML-KEM post-quantum encryption

### 3.2 Encryption Type Validation

```go
// IsValidEncryptionType checks if the encryption type is valid
func IsValidEncryptionType(encType EncryptionType) bool

// GetEncryptionTypeName returns the human-readable name of the encryption type
func GetEncryptionTypeName(encType EncryptionType) string
```

#### 3.2.1 Purpose

Provides validation and naming utilities for encryption types.

#### 3.2.2 Example Usage

```go
if !IsValidEncryptionType(encType) {
    return fmt.Errorf("invalid encryption type: %d", encType)
}

name := GetEncryptionTypeName(encType)
fmt.Printf("Using encryption: %s\n", name)
```

## 4. Generic Encryption Patterns

The security API provides generic encryption patterns that extend the generic configuration patterns defined in [api_generics.md](api_generics.md#28-generic-configuration-patterns) for type-safe encryption operations.

### 4.1 Generic Encryption Strategy Interface

```go
// EncryptionStrategy provides type-safe encryption for any data type
type EncryptionStrategy[T any] interface {
    Encrypt(ctx context.Context, data T, key EncryptionKey[T]) (T, error)
    Decrypt(ctx context.Context, data T, key EncryptionKey[T]) (T, error)
    Type() EncryptionType
    Name() string
    KeySize() int
    ValidateKey(ctx context.Context, key EncryptionKey[T]) error
}

// ByteEncryptionStrategy is the concrete implementation for []byte data
type ByteEncryptionStrategy interface {
    EncryptionStrategy[[]byte]
}

// EncryptionKey provides type-safe key management
type EncryptionKey[T any] struct {
    *Option[T]
    KeyType    EncryptionType
    KeyID      string
    CreatedAt  time.Time
    ExpiresAt  *time.Time
}

func NewEncryptionKey[T any](keyType EncryptionType, keyID string, key T) *EncryptionKey[T]
func (k *EncryptionKey[T]) GetKey() (T, bool)
func (k *EncryptionKey[T]) SetKey(key T)
func (k *EncryptionKey[T]) IsValid() bool
func (k *EncryptionKey[T]) IsExpired() bool
```

### 4.2 Generic Encryption Configuration

```go
// EncryptionConfig provides type-safe encryption configuration
type EncryptionConfig[T any] struct {
    *Config[T]

    // Encryption-specific settings
    EncryptionType     Option[EncryptionType]     // Encryption algorithm type
    KeySize           Option[int]                 // Key size in bits
    UseRandomIV       Option[bool]                // Use random initialization vector
    AuthenticationTag Option[bool]                // Include authentication tag
    CompressionLevel  Option[int]                 // Compression level for encrypted data
}

// EncryptionConfigBuilder provides fluent configuration building for encryption
type EncryptionConfigBuilder[T any] struct {
    config *EncryptionConfig[T]
}

func NewEncryptionConfigBuilder[T any]() *EncryptionConfigBuilder[T]
func (b *EncryptionConfigBuilder[T]) WithEncryptionType(encType EncryptionType) *EncryptionConfigBuilder[T]
func (b *EncryptionConfigBuilder[T]) WithKeySize(size int) *EncryptionConfigBuilder[T]
func (b *EncryptionConfigBuilder[T]) WithRandomIV(useRandom bool) *EncryptionConfigBuilder[T]
func (b *EncryptionConfigBuilder[T]) WithAuthenticationTag(useAuth bool) *EncryptionConfigBuilder[T]
func (b *EncryptionConfigBuilder[T]) Build() *EncryptionConfig[T]
```

### 4.3 Generic Encryption Validation

```go
// EncryptionValidator provides type-safe encryption validation
type EncryptionValidator[T any] struct {
    *Validator[T]
    encryptionRules []EncryptionValidationRule[T]
}

// EncryptionValidationRule is an alias for the generic ValidationRule
type EncryptionValidationRule[T any] = ValidationRule[T]

func (v *EncryptionValidator[T]) AddEncryptionRule(rule EncryptionValidationRule[T])
func (v *EncryptionValidator[T]) ValidateEncryptionData(ctx context.Context, data T) error
func (v *EncryptionValidator[T]) ValidateDecryptionData(ctx context.Context, data T) error
func (v *EncryptionValidator[T]) ValidateEncryptionKey(ctx context.Context, key EncryptionKey[T]) error
```

### 4.4 File Encryption Operations

```go
// FileEncryptionHandler provides file-specific encryption operations
type FileEncryptionHandler[T any] interface {
    EncryptFile(ctx context.Context, filePath string, data T, key EncryptionKey[T]) error
    DecryptFile(ctx context.Context, filePath string, key EncryptionKey[T]) (T, error)
    ValidateFileEncryption(ctx context.Context, filePath string) error
}

// Built-in file encryption handlers
type AES256GCMFileHandler struct { ... }
type ChaCha20Poly1305FileHandler struct { ... }
type MLKEMFileHandler struct { ... }
```

### 4.5 Package File Encryption Operations

```go
// EncryptFile encrypts a file using the security API's file encryption patterns
func (p *Package) EncryptFile[T any](ctx context.Context, path string, data T, handler FileEncryptionHandler[T], key EncryptionKey[T]) error

// DecryptFile decrypts a file using the security API's file encryption patterns
func (p *Package) DecryptFile[T any](ctx context.Context, path string, handler FileEncryptionHandler[T], key EncryptionKey[T]) (T, error)

// ValidateFileEncryption validates file encryption using the security API's file encryption patterns
func (p *Package) ValidateFileEncryption[T any](ctx context.Context, path string, handler FileEncryptionHandler[T]) error

// GetFileEncryptionInfo gets encryption information for a file using the security API's patterns
func (p *Package) GetFileEncryptionInfo[T any](ctx context.Context, path string) (*EncryptionConfig[T], error)
```

## 5. ML-KEM Key Management

### 5.1 ML-KEM Key Structure

```go
// ML-KEM Key Structure
type MLKEMKey struct {
    PublicKey  []byte  // ML-KEM public key data
    PrivateKey []byte  // ML-KEM private key data
    Level      int     // Security level (1-5)
}
```

#### 5.1.1 Purpose

Represents an ML-KEM key pair with associated security level.

#### 5.1.2 Fields

- `PublicKey`: Public key data for encryption
- `PrivateKey`: Private key data for decryption
- `Level`: Security level (1-5, higher is more secure)

### 5.2 ML-KEM Key Generation

```go
// GenerateMLKEMKey generates a new ML-KEM key at specified security level
func GenerateMLKEMKey(ctx context.Context, level int) (*MLKEMKey, error)
```

#### 5.2.1 Purpose

Generates a new ML-KEM key pair at the specified security level.

#### 5.2.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `level`: Security level (1-5)

#### 5.2.3 Returns

New `MLKEMKey` instance

#### 5.2.4 Error Conditions

- `ErrInvalidSecurityLevel`: Security level must be 1-5
- `ErrKeyGenerationFailed`: Failed to generate key pair
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 5.2.5 Example Usage

```go
key, err := GenerateMLKEMKey(ctx, 3)
if err != nil {
    return fmt.Errorf("failed to generate ML-KEM key: %w", err)
}
defer key.Clear() // Clear sensitive data
```

### 5.3 ML-KEM Encryption Operations

```go
// Encrypt encrypts plaintext using ML-KEM key
func (k *MLKEMKey) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error)

// Decrypt decrypts ciphertext using ML-KEM key
func (k *MLKEMKey) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error)
```

#### 5.3.1 Purpose

Performs encryption and decryption operations using ML-KEM keys.

#### 5.3.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `plaintext`: Data to encrypt (for Encrypt method)
- `ciphertext`: Encrypted data to decrypt (for Decrypt method)

#### 5.3.3 Returns

Encrypted data (for Encrypt) or decrypted data (for Decrypt)

#### 5.3.4 Error Conditions

- `ErrEncryptionFailed`: Failed to encrypt data
- `ErrDecryptionFailed`: Failed to decrypt data
- `ErrInvalidKey`: Key is invalid or corrupted
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 5.3.5 Example Usage

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

### 5.4 ML-KEM Key Information

```go
// GetPublicKey returns the public key data
func (k *MLKEMKey) GetPublicKey() []byte

// GetLevel returns the security level of the key
func (k *MLKEMKey) GetLevel() int

// Clear clears sensitive key data from memory
func (k *MLKEMKey) Clear()
```

#### 5.4.1 Purpose

Provides access to key information and secure cleanup.

#### 5.4.2 Example Usage

```go
publicKey := key.GetPublicKey()
level := key.GetLevel()
fmt.Printf("ML-KEM key level %d, public key size: %d bytes\n", level, len(publicKey))

// Clear sensitive data when done
key.Clear()
```
