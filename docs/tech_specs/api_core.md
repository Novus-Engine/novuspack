# NovusPack Technical Specifications - Core Package Interface API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
  - [0.2 Context Integration](#02-context-integration)
- [1. Core Interfaces](#1-core-interfaces)
  - [1.1 Package Reader Interface](#11-package-reader-interface)
  - [1.2 Package Writer Interface](#12-package-writer-interface)
  - [1.3 Package Interface](#13-package-interface)
- [2. Basic Operations](#2-basic-operations)
- [3. Package Writing Operations](#3-package-writing-operations)
- [4. File Management](#4-file-management)
- [5. Encryption Management](#5-encryption-management)
- [6. Package Compression Operations](#6-package-compression-operations)
  - [6.1 Package Compression Types](#61-package-compression-types)
  - [6.2 Package Compression Functions](#62-package-compression-functions)
  - [6.3 Package Compression Behavior](#63-package-compression-behavior)
- [7. Digital Signatures and Security](#7-digital-signatures-and-security)
  - [7.1 Core Integration Points](#71-core-integration-points)
  - [7.4 Write Protection and Immutability Enforcement](#74-write-protection-and-immutability-enforcement)
- [8. Per-File Tags Management](#8-per-file-tags-management)
- [9. Package Metadata Management](#9-package-metadata-management)
  - [9.1 General Metadata Operations](#91-general-metadata-operations)
  - [9.2 AppID/VendorID Management](#92-appidvendorid-management)
  - [9.3 Package Information Structures](#93-package-information-structures)
- [10. File Validation Requirements](#10-file-validation-requirements)
- [11. Structured Error System](#11-structured-error-system)
  - [11.1 Error Types and Categories](#111-error-types-and-categories)
  - [Error Type Categories](#error-type-categories)
  - [11.2 PackageError Structure](#112-packageerror-structure)
  - [11.3 Error Helper Functions](#113-error-helper-functions)
  - [11.4 Error Handling Patterns](#114-error-handling-patterns)
  - [11.5 Migration from Sentinel Errors](#115-migration-from-sentinel-errors)
  - [Benefits of Structured Errors](#benefits-of-structured-errors)
- [12. Generic Types](#12-generic-types)

---

## 0. Overview

This document defines the core package interface for the NovusPack system, including basic operations, file management, encryption, digital signatures, and metadata handling.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Basic Operations API](api_basic_operations.md) - Package creation, opening, closing, and lifecycle management
- [File Management API](api_file_management.md) - File operations and encryption
- [Package Writing Operations](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [Streaming and Buffer Management](api_streaming.md) - File streaming interface and buffer management system
- [Multi-Layer Deduplication](api_deduplication.md) - Content deduplication strategies and processing levels
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation
- [Package Metadata API](api_metadata.md) - Comment management, AppID/VendorID, and metadata operations
- [Security Validation API](api_security.md) - Package validation and security status structures
- [Generic Types and Patterns](api_generics.md) - Generic type system, patterns, and best practices
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation

### 0.2 Context Integration

All public methods in the NovusPack API accept `context.Context` as the first parameter to support:

- Request cancellation and timeout handling
- Request-scoped values and configuration
- Graceful shutdown and resource cleanup
- Integration with Go's standard context patterns

This follows 2025 Go best practices and ensures the API is compatible with modern Go applications and frameworks.

## 1. Core Interfaces

The NovusPack API is designed around core interfaces that provide clear separation of concerns and enable better testability.

### 1.1 Package Reader Interface

```go
type PackageReader interface {
    ReadFile(ctx context.Context, path string) ([]byte, error)
    ListFiles(ctx context.Context) ([]FileInfo, error)
    GetMetadata(ctx context.Context) (*PackageInfo, error)
    Validate(ctx context.Context) error
    GetInfo(ctx context.Context) PackageInfo
}
```

### 1.2 Package Writer Interface

```go
type PackageWriter interface {
    WriteFile(ctx context.Context, path string, data []byte, opts *AddFileOptions) error
    RemoveFile(ctx context.Context, path string) error
    Write(ctx context.Context, path string, compression CompressionType, sign bool) error
    SafeWrite(ctx context.Context, path string, compression CompressionType, sign bool) error
    FastWrite(ctx context.Context, path string, compression CompressionType, sign bool) error
}
```

### 1.3 Package Interface

```go
type Package interface {
    PackageReader
    PackageWriter
    Close() error
    IsOpen() bool
    Defragment(ctx context.Context) error
}
```

## 2. Basic Operations

See [Basic Operations API](api_basic_operations.md) for detailed method signatures and implementation details.

- **Create**: Creates a new package at the specified path
- **Open**: Opens an existing package from the specified path
- **Close**: Closes the package and releases resources
- **Write**: General write method with compression handling using SafeWrite or FastWrite methods
  - See [Package Writing API](api_writing.md) for detailed method signatures and implementation details.
- **Defragment**: Optimizes package structure and removes unused space
- **Validate**: Validates package format, structure, and integrity
- **GetInfo**: Gets comprehensive package information

## 3. Package Writing Operations

See [Package Writing API](api_writing.md) for detailed method signatures and implementation details.

- **SafeWrite**: Atomic write with temp file strategy for data safety
- **FastWrite**: In-place updates for existing packages for performance
- **Write**: General write method with compression handling using SafeWrite or FastWrite methods

## 4. File Management

See [File Management API](api_file_management.md) for detailed method signatures and implementation details.

- **Basic File Operations**: Add, remove, and extract files from packages
- **Encryption-Aware File Management**: Add files with specific encryption types
- **Encryption Type System**: Define and validate encryption algorithms
- **ML-KEM Key Management**: Generate and manage post-quantum encryption keys
- **File Pattern Operations**: Add multiple files using patterns
- **File Information and Queries**: Get file information and search capabilities

## 5. Encryption Management

See [File Management API - ML-KEM Key Management](api_file_management.md#4-ml-kem-key-management) for detailed method signatures and implementation details.

- **ML-KEM Key Management**: Generate and manage post-quantum encryption keys
- **Key Operations**: Encrypt, decrypt, and manage key lifecycle
- **Security Levels**: Support for multiple security levels (1-5)

## 6. Package Compression Operations

### 6.1 Package Compression Types

See [Package Compression API - Compression Types](api_package_compression.md#12-compression-types) for compression type constants and [Compression Information Structure](api_package_compression.md#13-compression-information-structure) for the PackageCompressionInfo structure.

### 6.2 Package Compression Functions

See [Package Compression API](api_package_compression.md) for detailed method signatures and implementation details.

- **CompressPackage**: Compresses package content in memory
- **DecompressPackage**: Decompresses the package in memory
- **CompressPackageFile**: Compresses package content and writes to specified path
- **DecompressPackageFile**: Decompresses the package and writes to specified path
- **GetPackageCompressionInfo**: Returns package compression information
- **IsPackageCompressed**: Checks if the package is compressed
- **GetPackageCompressionType**: Returns the package compression type
- **SetPackageCompressionType**: Sets the package compression type (without compressing)
- **CanCompressPackage**: Checks if package can be compressed (not signed)

### 6.3 Package Compression Behavior

- **Compression Scope**: Compresses package content (file entries + data + index) but NOT header, comment, or signatures
- **Header Exclusion**: Package header must remain uncompressed for compression type detection
- **Comment Exclusion**: Package comment must remain uncompressed for easy reading without decompression
- **Signature Exclusion**: Digital signatures must remain uncompressed for validation
- **Relationship to Per-File Compression**: Package compression is applied after per-file compression
- **Decompression Order**: Package decompression must occur before per-file decompression
- **Signing Compatibility**: Compressed packages can be signed, but signed packages cannot be compressed
- **Use Cases**:
  - Small packages where package-level compression is more efficient
  - Packages with many small files that benefit from global compression
  - Archival scenarios where maximum compression is desired

#### 5.2.4 Signing and Compression Relationship

#### 5.2.4.1 Supported Operations

- **Signing Compressed Packages**: Compressed packages can be signed
  - Process: Compress package first, then add signatures
  - Signatures validate the compressed content
  - Header, comment, and signatures remain uncompressed for access

#### 5.2.4.2 Unsupported Operations

- **Compressing Signed Packages**: Signed packages cannot be compressed
  - Reason: Would require decompression to access signatures for validation
  - Error: `CompressSignedPackageError` returned if attempted
  - Workflow: Must clear signatures first, then compress, then re-sign

#### 5.2.4.3 Error Handling

- **CompressSignedPackageError**: Returned when attempting to compress signed package
- **Validation**: All compression functions check for existing signatures
- **Clear Workflow**: Must clear signatures => compress => re-sign

## 7. Digital Signatures and Security

**Cross-Reference**: For complete signature management, validation, and implementation details, see [Digital Signature API](api_signatures.md).

### 7.1 Core Integration Points

The core package interface integrates with the signature system through:

- **Immutability Enforcement**: All write operations check `SignatureOffset > 0` before proceeding
- **Write Protection**: Signed packages are protected from write operations by default
- **Context Integration**: All signature operations accept `ctx context.Context` as first parameter
- **Error Handling**: Signature operations use the structured error system defined in [Structured Error System](#11-structured-error-system)

### 7.4 Write Protection and Immutability Enforcement

- **Signed File Detection**: All write operations must check if `SignatureOffset > 0` before proceeding
- **Write Protection**: Signed packages are protected from write operations by default
- **Clear-Signatures Flag**: Write operations are refused unless `clearSignatures` flag is passed
- **Allowed Operations**: If signed, only signature addition operations are allowed
- **Prohibited Operations**: Header modifications and content changes are prohibited on signed packages
- **Signature Removal**: Removing signatures is not recommended as it invalidates all subsequent signatures
- **Detailed Behavior**: See [Package Writing Operations](api_writing.md#signed-file-write-operations) for complete implementation details
- **Purpose**: This prevents accidental signature invalidation and maintains package integrity

## 8. Per-File Tags Management

```go
// SetFileTags sets tags for a specific file
func (p *Package) SetFileTags(ctx context.Context, path string, tags map[string]interface{}) error

// GetFileTags retrieves tags for a specific file
func (p *Package) GetFileTags(ctx context.Context, path string) (map[string]interface{}, error)

// UpdateFileTags updates existing tags for a file
func (p *Package) UpdateFileTags(ctx context.Context, path string, updates map[string]interface{}) error

// RemoveFileTags removes specific tag keys from a file
func (p *Package) RemoveFileTags(ctx context.Context, path string, keys []string) error

// ClearFileTags removes all tags from a file
func (p *Package) ClearFileTags(ctx context.Context, path string) error

// GetFilesByTag searches for files by tag key-value pair
func (p *Package) GetFilesByTag(ctx context.Context, key string, value interface{}) ([]string, error)

// GetInheritedTags retrieves tags including inheritance from parent directories
func (p *Package) GetInheritedTags(ctx context.Context, path string) (map[string]interface{}, error)
```

## 9. Package Metadata Management

See [Package Metadata API](metadata.md) for detailed method signatures and implementation details.

### 9.1 General Metadata Operations

- **SetMetadata**: Sets package metadata
- **GetMetadata**: Retrieves package metadata
- **UpdateMetadata**: Updates package metadata
- **ValidateMetadata**: Validates metadata structure and content
- **HasMetadata**: Checks if package has metadata
- **AddMetadataFile**: Adds metadata as special file
- **GetMetadataFile**: Retrieves metadata from special file
- **UpdateMetadataFile**: Updates metadata file
- **RemoveMetadataFile**: Removes metadata special file

### 9.2 AppID/VendorID Management

- **SetAppID**: Sets the application identifier
- **GetAppID**: Gets the current application identifier
- **ClearAppID**: Clears the application identifier
- **HasAppID**: Checks if application identifier is set
- **SetVendorID**: Sets the vendor/platform identifier
- **GetVendorID**: Gets the current vendor identifier
- **ClearVendorID**: Clears the vendor identifier
- **HasVendorID**: Checks if vendor identifier is set

### 9.3 Package Information Structures

- **PackageInfo**: Comprehensive package information and metadata
- **SignatureInfo**: Detailed signature information
- **SecurityStatus**: Current security status and validation results

## 10. File Validation Requirements

Requirements for files that are added to packages.

- Path must not be empty or contain only whitespace
- Path is normalized like tar files (removes redundant separators, resolves relative references)
- Data must not be nil (empty files with len = 0 are allowed)
- Invalid files will be rejected with appropriate error messages

For complete file validation specifications, see [File Validation Requirements](file_validation.md).

## 11. Structured Error System

The NovusPack API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities while maintaining compatibility with Go's standard error handling patterns.

### 11.1 Error Types and Categories

```go
type ErrorType int

const (
    ErrTypeValidation ErrorType = iota
    ErrTypeIO
    ErrTypeSecurity
    ErrTypeCorruption
    ErrTypeUnsupported
    ErrTypeContext
    ErrTypeEncryption
    ErrTypeCompression
    ErrTypeSignature
)
```

### Error Type Categories

- **ErrTypeValidation**: Input validation errors, invalid parameters, format errors
- **ErrTypeIO**: File I/O errors, permission errors, disk space issues
- **ErrTypeSecurity**: Security-related errors, access denied, authentication failures
- **ErrTypeCorruption**: Data corruption, checksum failures, integrity violations
- **ErrTypeUnsupported**: Unsupported features, versions, or operations
- **ErrTypeContext**: Context cancellation, timeout, or context-related errors
- **ErrTypeEncryption**: Encryption/decryption failures, key management errors
- **ErrTypeCompression**: Compression/decompression failures, algorithm errors
- **ErrTypeSignature**: Digital signature validation, signing failures

### 11.2 PackageError Structure

```go
type PackageError struct {
    Type    ErrorType              // Error category
    Message string                 // Human-readable error message
    Cause   error                  // Underlying error (wrapped)
    Context map[string]interface{} // Additional error context
}

// Error implements the error interface
func (e *PackageError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *PackageError) Unwrap() error {
    return e.Cause
}

// Is implements error matching for sentinel errors
func (e *PackageError) Is(target error) bool {
    if e.Cause != nil {
        return errors.Is(e.Cause, target)
    }
    return false
}

// WithContext adds additional context to the error
func (e *PackageError) WithContext(key string, value interface{}) *PackageError {
    if e.Context == nil {
        e.Context = make(map[string]interface{})
    }
    e.Context[key] = value
    return e
}
```

### 11.3 Error Helper Functions

```go
// NewPackageError creates a new structured error
func NewPackageError(errType ErrorType, message string, cause error) *PackageError {
    return &PackageError{
        Type:    errType,
        Message: message,
        Cause:   cause,
        Context: make(map[string]interface{}),
    }
}

// WrapError wraps an existing error with structured information
func WrapError(err error, errType ErrorType, message string) *PackageError {
    return NewPackageError(errType, message, err)
}

// IsPackageError checks if an error is a PackageError
func IsPackageError(err error) (*PackageError, bool) {
    var pkgErr *PackageError
    if errors.As(err, &pkgErr) {
        return pkgErr, true
    }
    return nil, false
}

// GetErrorType returns the error type if the error is a PackageError
func GetErrorType(err error) (ErrorType, bool) {
    if pkgErr, ok := IsPackageError(err); ok {
        return pkgErr.Type, true
    }
    return 0, false
}

// WithTypedContext adds type-safe context to errors
func WithTypedContext[T any](err error, key string, value T) error

// GetTypedContext retrieves type-safe context from errors
func GetTypedContext[T any](err error, key string) (T, bool)

// NewTypedPackageError creates a structured error with type-safe context
func NewTypedPackageError[T any](errType ErrorType, message string, cause error, context T) *PackageError

// WrapWithContext wraps an error with type-safe context
func WrapWithContext[T any](err error, errType ErrorType, message string, context T) *PackageError

// MapError transforms an error with a generic mapper function
func MapError[T any, U any](err error, mapper func(T) U) error
```

### 11.4 Error Handling Patterns

#### 11.4.1 Creating Structured Errors

```go
// Create a new validation error
err := NewPackageError(ErrTypeValidation, "file not found", ErrFileNotFound).
    WithContext("path", "/path/to/file").
    WithContext("operation", "AddFile")

// Wrap an existing error
err := WrapError(io.ErrUnexpectedEOF, ErrTypeIO, "unexpected end of file").
    WithContext("file", "package.npk").
    WithContext("offset", 1024)
```

#### 11.4.2 Error Inspection and Handling

```go
// Check if error is a PackageError
if pkgErr, ok := IsPackageError(err); ok {
    switch pkgErr.Type {
    case ErrTypeValidation:
        // Handle validation errors
        log.Printf("Validation error: %s", pkgErr.Message)
        if path, exists := pkgErr.Context["path"]; exists {
            log.Printf("Failed path: %v", path)
        }
    case ErrTypeIO:
        // Handle I/O errors
        log.Printf("I/O error: %s", pkgErr.Message)
    case ErrTypeSecurity:
        // Handle security errors
        log.Printf("Security error: %s", pkgErr.Message)
    default:
        // Handle other error types
        log.Printf("Error (%d): %s", pkgErr.Type, pkgErr.Message)
    }
}

// Check specific error types
if errType, ok := GetErrorType(err); ok && errType == ErrTypeValidation {
    // Handle validation errors specifically
}
```

#### 11.4.3 Error Propagation

```go
func (p *Package) AddFile(ctx context.Context, path string, data []byte) error {
    if err := p.validateFile(path, data); err != nil {
        return WrapError(err, ErrTypeValidation, "file validation failed").
            WithContext("path", path).
            WithContext("size", len(data))
    }

    if err := p.writeFileData(path, data); err != nil {
        return WrapError(err, ErrTypeIO, "failed to write file data").
            WithContext("path", path)
    }

    return nil
}
```

#### 11.4.4 Sentinel Error Compatibility

```go
// Sentinel errors are still supported and can be wrapped
var (
    ErrFileNotFound = errors.New("file not found")
    ErrInvalidPath  = errors.New("invalid file path")
    ErrIOError      = errors.New("I/O error")
)

// Convert sentinel errors to structured errors
func wrapSentinelError(err error, message string) *PackageError {
    switch err {
    case ErrFileNotFound:
        return NewPackageError(ErrTypeValidation, message, err)
    case ErrInvalidPath:
        return NewPackageError(ErrTypeValidation, message, err)
    case ErrIOError:
        return NewPackageError(ErrTypeIO, message, err)
    default:
        return NewPackageError(ErrTypeUnsupported, message, err)
    }
}
```

#### 11.4.5 Error Logging and Debugging

```go
// Log structured errors with full context
func logError(err error, operation string) {
    if pkgErr, ok := IsPackageError(err); ok {
        log.Printf("Error in %s: Type=%d, Message=%s, Context=%+v",
            operation, pkgErr.Type, pkgErr.Message, pkgErr.Context)
        if pkgErr.Cause != nil {
            log.Printf("Caused by: %v", pkgErr.Cause)
        }
    } else {
        log.Printf("Error in %s: %v", operation, err)
    }
}
```

### 11.5 Migration from Sentinel Errors

The structured error system is designed to work alongside existing sentinel errors.
Existing code using sentinel errors will continue to work, while new code can take advantage of structured errors for better error handling and debugging.

### Benefits of Structured Errors

- **Better Error Categorization**: Errors grouped by type for easier handling
- **Rich Error Context**: Additional context fields for debugging
- **Type Safety**: Structured errors can be inspected with type assertions
- **Backward Compatible**: Coexists with sentinel errors
- **Better Logging**: Structured errors provide more information for logs
- **Testing**: Easier to test error conditions with typed errors

## 12. Generic Types

For comprehensive generic type definitions, usage examples, and best practices, see [Generic Types and Patterns](api_generics.md).

The NovusPack API provides generic types for improved type safety and code reuse across different data types. All generic type definitions, interfaces, and usage patterns are documented in the dedicated [Generic Types and Patterns](api_generics.md) specification.
