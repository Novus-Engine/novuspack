# NovusPack Technical Specifications - Basic Operations API

## Table of Contents

- [Table of Contents](#table-of-contents)
- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
  - [0.2 Context Integration](#02-context-integration)
- [1. Package Structure and Loading](#1-package-structure-and-loading)
  - [1.1 Package Structure](#11-package-structure)
  - [1.2 Package Loading Process](#12-package-loading-process)
- [2. Package Constants](#2-package-constants)
  - [2.1 Package Format Constants](#21-package-format-constants)
- [3. Package Lifecycle Operations](#3-package-lifecycle-operations)
- [4. Package Creation](#4-package-creation)
  - [4.1 Package Constructor](#41-package-constructor)
  - [4.2 Create Method](#42-create-method)
  - [4.3 Create with Options](#43-create-with-options)
  - [4.4 Package Builder Pattern](#44-package-builder-pattern)
- [5. Package Opening](#5-package-opening)
  - [5.1 Open Method](#51-open-method)
  - [5.2 Open with Validation](#52-open-with-validation)
- [6. Package Closing](#6-package-closing)
  - [6.1 Close Method](#61-close-method)
  - [6.2 Close with Cleanup](#62-close-with-cleanup)
- [7. Package Operations](#7-package-operations)
  - [7.1 Package Validation](#71-package-validation)
  - [7.2 Package Defragmentation](#72-package-defragmentation)
  - [7.3 Package Information](#73-package-information)
  - [7.4 Header Inspection](#74-header-inspection)
  - [7.5 Check Package State](#75-check-package-state)
- [8. Error Handling](#8-error-handling)
  - [8.1 Structured Error System](#81-structured-error-system)
  - [8.2 Common Error Types](#82-common-error-types)
  - [8.3 Structured Error Examples](#83-structured-error-examples)
  - [8.4 Error Handling Best Practices](#84-error-handling-best-practices)
- [9. Best Practices](#9-best-practices)
  - [9.1 Package Lifecycle Management](#91-package-lifecycle-management)
  - [9.2 Error Handling](#92-error-handling)
  - [9.3 Resource Management](#93-resource-management)

## 0. Overview

This document defines the basic package operations for the NovusPack system, covering the fundamental lifecycle operations of creating, opening, and closing packages.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Package Writing Operations](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [Package Metadata API](api_metadata.md) - Comment management, AppID/VendorID, and metadata operations
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation

### 0.2 Context Integration

All public methods in the NovusPack Basic Operations API accept `context.Context` as the first parameter to support:

- Request cancellation and timeout handling
- Request-scoped values and configuration
- Graceful shutdown and resource cleanup
- Integration with Go's standard context patterns

This follows 2025 Go best practices and ensures the API is compatible with modern Go applications and frameworks.

These methods assume the following imports:

```go
import (
    "context"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)
// FileEntry, DirectoryEntry, PackageInfo, PackageHeader, PackageIndex from other API modules
```

## 1. Package Structure and Loading

### 1.1 Package Structure

```go
// Package represents a NovusPack package with all its metadata and files
type Package struct {
    // Package metadata
    Info *PackageInfo

    // File entries (loaded on demand)
    FileEntries []*FileEntry

    // Directory metadata (loaded from special files)
    DirectoryEntries []*DirectoryEntry

    // Special metadata files
    SpecialFiles map[uint16]*FileEntry // fileType -> FileEntry

    // Package state
    IsOpen bool
    FilePath string

    // Internal state
    header *PackageHeader
    index *PackageIndex
    fileHandle *os.File
}
```

### 1.2 Package Loading Process

When a package is opened, the following initialization sequence occurs:

1. **Load package header** - Validates magic number and version
2. **Load package info** - Retrieves metadata (comment, VendorID, AppID)
3. **Load file entries** - Reads file index and entry structures
4. **Load special metadata files** - Processes special file types (65000-65535)
5. **Load directory metadata** - Parses directory structure from YAML
6. **Update file-directory associations** - Links files to their parent directories

```go
// OpenPackage loads a package from file and initializes all metadata
func OpenPackage(ctx context.Context, filePath string) (*Package, error)

// loadSpecialMetadataFiles loads all special metadata files
func (p *Package) loadSpecialMetadataFiles(ctx context.Context) error

// loadDirectoryMetadata loads directory metadata from special files
func (p *Package) loadDirectoryMetadata(ctx context.Context) error

// updateFileDirectoryAssociations links files to their directory metadata
func (p *Package) updateFileDirectoryAssociations(ctx context.Context) error
```

## 2. Package Constants

### 2.1 Package Format Constants

```go
const (
    // NPKMagic is the magic number for .npk files
    NPKMagic = 0x4E56504B // "NVPK" in hex

    // NPKVersion is the current version of the .npk format
    NPKVersion = 1

    // HeaderSize is the fixed size of the package header
    // See: Package File Format - Package Header for authoritative definition
    HeaderSize int64 = 112
)
```

These constants define fundamental values for the NovusPack format.

#### 2.1.1 Constants

- `NPKMagic`: Package identifier (0x4E56504B "NVPK")
- `NPKVersion`: Current format version (1)
- `HeaderSize`: Fixed header size in bytes (see [Package File Format - Package Header](package_file_format.md#2-package-header))

**Usage**: Validate package header magic number and version before processing.

## 3. Package Lifecycle Operations

The NovusPack system follows a simple lifecycle pattern:

1. **Create** - Create a new package
2. **Open** - Open an existing package
3. **Operations** - Perform various operations (add files, metadata, etc.)
4. **Close** - Close the package and release resources

## 4. Package Creation

### 4.1 Package Constructor

```go
// NewPackage creates a new empty package
func NewPackage(ctx context.Context) *Package
```

This function creates a new, empty NovusPack package in memory with default values.
The package exists only in memory until written to disk using one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`).

Returns a new `Package` instance with:

- Default header values (magic number, version, timestamps)
- Empty file index
- Empty package comment
- Closed state set to false

#### 4.1.1 NewPackage Behavior

- Creates package structure in memory only (no file I/O operations performed)
- Initializes package with standard NovusPack header
- Sets creation timestamp to current time
- Prepares package for file operations
- Package must be written to disk using one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`) before it can be persisted

#### 4.1.2 NewPackage Example Usage

```go
package := NewPackage(ctx)
defer package.Close()
```

### 4.2 Create Method

```go
// Create configures a package for creation at the specified path
func (p *Package) Create(ctx context.Context, path string) error
```

This function configures an existing package (typically created with `NewPackage`) for writing to disk at the specified path.
**This function does not write to disk** - it only prepares the package structure and sets the target path.
The package file is only written to disk when one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`) is called.

**Note**: This function internally uses `NewPackage` if called on a nil package, ensuring the package exists in memory before configuration.

**Path Validation**: This function validates that the provided path is valid and the target directory is writable, even though it doesn't write to disk. This ensures early detection of path-related issues.

#### 4.2.1 Create Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: File system path where the package will be written (when Write is called)

#### 4.2.2 Create Behavior

- Validates that the provided path is valid and well-formed
- Validates that the target directory exists and is writable (fails if directory does not exist)
- Configures package with the standard NovusPack header structure (in memory)
- Initializes all package metadata to default values
- Sets up the basic package structure (file entries, data sections)
- Stores the target path for later writing operations
- The package remains in an unsigned, uncompressed state until written
- No file I/O operations are performed on the target file - package remains in memory

#### 4.2.3 Create Error Conditions

- **Validation Errors**:
  - Invalid or malformed file path
  - Target directory does not exist (parent directories are not created)
  - Target directory is not writable
  - Insufficient permissions to create file in target directory
- **Security Errors**: Insufficient permissions to access target directory
- **Context Errors**: Context cancellation or timeout exceeded

**Note**: While the target file is not created during `Create`, the path and directory are validated to ensure they exist and are writable. This enables early error detection before file operations begin. The parent directory must already exist - `Create` will not create missing parent directories.

#### 4.2.4 Create Example Usage

```go
// Create package in memory
package := NewPackage(ctx)

// Configure for writing at specific path (still in memory)
err := package.Create(ctx, "/path/to/new-package.npk")
if err != nil {
    return err
}

// Add files, metadata, etc...
// ...

// Write to disk
err = package.Write(ctx, "/path/to/new-package.npk", 0, false)
if err != nil {
    return err
}
```

### 4.3 Create with Options

```go
// CreateWithOptions configures a package with specified options for later writing
func (p *Package) CreateWithOptions(ctx context.Context, path string, options CreateOptions) error

type CreateOptions struct {
    Comment     string    // Initial package comment
    VendorID    uint32    // Vendor identifier
    AppID       uint64    // Application identifier
    Permissions os.FileMode // File permissions (default: 0644)
}
```

This function configures an existing package (typically created with `NewPackage`) with initial metadata and target path.
**This function does not write to disk** - it only configures the package structure in memory.
The package file is only written to disk when one of the Write functions (`Write`, `SafeWrite`, or `FastWrite`) is called.

**Note**: This function internally uses `NewPackage` if called on a nil package, ensuring the package exists in memory before configuration.

**Implementation**: This function calls `Create` internally to validate the path and set up the basic package structure, then applies the provided options.

#### 4.3.1 Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: File system path where the package will be written (when Write is called)
- `options`: Initial package configuration and metadata

#### 4.3.2 Behavior

- Calls `Create` internally to validate the path and set up the basic package structure
- Applies the provided options to configure the package:
  - Sets initial package comment, VendorID, and AppID if provided
  - Stores file permissions for later use during Write operations
  - Initializes package with provided metadata
- Stores the target path for later writing operations
- No file I/O operations are performed on the target file - package remains in memory

**Path Validation**: Path validation is performed by the underlying `Create` method. Any path validation errors from `Create` will be propagated.

#### 4.3.3 CreateWithOptions Error Conditions

This function inherits all error conditions from `Create` since it calls `Create` internally:

- **Validation Errors**:
  - Invalid or malformed file path
  - Target directory does not exist and cannot be created
  - Target directory is not writable
  - Insufficient permissions to create file in target directory
- **Security Errors**: Insufficient permissions to access or create target directory
- **Context Errors**: Context cancellation or timeout exceeded
- **Additional Validation Errors**: Invalid option values (e.g., invalid VendorID format)

#### 4.3.4 CreateWithOptions Example Usage

```go
// Create package in memory
package := NewPackage(ctx)

// Configure with options (still in memory)
options := CreateOptions{
    Comment:  "My Game Package",
    VendorID: 0x00000001, // Steam
    AppID:    0x00000000000002DA, // CS:GO
    Permissions: 0644,
}

err := package.CreateWithOptions(ctx, "/path/to/game-package.npk", options)
if err != nil {
    return err
}

// Add files, metadata, etc...
// ...

// Write to disk
err = package.Write(ctx, "/path/to/game-package.npk", 0, false)
if err != nil {
    return err
}
```

### 4.4 Package Builder Pattern

The builder pattern provides a fluent interface for creating packages with complex configurations.

```go
type PackageBuilder interface {
    WithCompression(comp CompressionType) PackageBuilder
    WithEncryption(enc EncryptionType) PackageBuilder
    WithMetadata(metadata map[string]string) PackageBuilder
    WithComment(comment string) PackageBuilder
    WithVendorID(vendorID uint32) PackageBuilder
    WithAppID(appID uint64) PackageBuilder
    Build(ctx context.Context) (Package, error)
}

// NewBuilder creates a new package builder
func NewBuilder() PackageBuilder
```

#### 4.4.1 Purpose

Provides a fluent interface for creating packages with complex configurations, improving code readability and reducing parameter complexity.

#### 4.4.2 Example Usage

```go
package, err := NewBuilder().
    WithCompression(CompressionZstandard).
    WithEncryption(EncryptionAES256GCM).
    WithComment("My application package").
    WithVendorID(0x12345678).
    WithAppID(0x87654321).
    Build(ctx)
```

## 5. Package Opening

### 5.1 Open Method

```go
// Open opens an existing package from the specified path
func (p *Package) Open(ctx context.Context, path string) error
```

This function opens an existing NovusPack package file for reading and modification.

#### 5.1.1 Open Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: File system path to the existing package file

#### 5.1.2 Open Behavior

- Opens the package file for reading and writing
- Validates the package header and basic structure
- Loads package metadata (comment, VendorID, AppID, etc.)
- Reads file entries and prepares for file operations
- Loads signature information if present
- Sets up package state for subsequent operations

#### 5.1.3 Open Error Conditions

- **Validation Errors**: Package file not found, corrupted or invalid format
- **Unsupported Errors**: Package version not supported
- **Security Errors**: Insufficient permissions to open file
- **I/O Errors**: File system errors during opening
- **Context Errors**: Context cancellation or timeout exceeded

#### 5.1.4 Open Example Usage

```go
err := package.Open(ctx, "/path/to/existing-package.npk")
if err != nil {
    return err
}
```

### 5.2 Open with Validation

```go
// OpenWithValidation opens a package and performs full validation
func (p *Package) OpenWithValidation(ctx context.Context, path string) error
```

This function opens a package and performs comprehensive validation of its integrity.

#### 5.2.1 OpenWithValidation Behavior

- Opens the package file
- Performs full package validation (structure, checksums, signatures)
- Returns detailed error information if validation fails
- Ensures package integrity before allowing operations

#### 5.2.2 OpenWithValidation Error Conditions

- All errors from `Open` method
- **Validation Errors**: Package validation failed, invalid signatures
- **Corruption Errors**: File checksums don't match, data integrity issues

## 6. Package Closing

### 6.1 Close Method

```go
// Close closes the package and releases resources
func (p *Package) Close() error
```

This function closes the package file and releases all associated resources.

#### 6.1.1 Close Behavior

- Closes the package file handle
- Releases memory buffers and caches
- Clears package state and metadata
- Performs any necessary cleanup operations
- Does not modify the package file (use Write methods to save changes)

#### 6.1.2 Close Error Conditions

- **I/O Errors**: File system errors during closing
- **Validation Errors**: Package is not currently open

#### 6.1.3 Close Example Usage

```go
err := package.Close()
if err != nil {
    return err
}
```

### 6.2 Close with Cleanup

```go
// CloseWithCleanup closes the package and performs cleanup operations
func (p *Package) CloseWithCleanup(ctx context.Context) error
```

This function closes the package and performs additional cleanup operations.

#### 6.2.1 CloseWithCleanup Behavior

- Closes the package file
- Performs cleanup operations (defragmentation, optimization)
- Releases all resources
- May take longer than standard Close due to cleanup operations

## 7. Package Operations

### 7.1 Package Validation

```go
// Validate validates package format, structure, and integrity
func (p *Package) Validate(ctx context.Context) error
```

This function performs comprehensive validation of the package format, structure, and integrity.

#### 7.1.1 Validate Behavior

- Validates package header format and version
- Checks file entry structure and consistency
- Verifies data section integrity and checksums
- Validates digital signatures if present
- Ensures package follows NovusPack specifications
- Returns detailed error information for any issues found

#### 7.1.2 Validate Error Conditions

- **Validation Errors**: Package not open, invalid format, validation failed
- **Corruption Errors**: Invalid signatures, checksum mismatches
- **Context Errors**: Context cancellation or timeout exceeded

#### 7.1.3 Validate Example Usage

```go
err := package.Validate(ctx)
if err != nil {
    return err
}
```

### 7.2 Package Defragmentation

```go
// Defragment optimizes package structure and removes unused space
func (p *Package) Defragment(ctx context.Context) error
```

This function optimizes package structure by removing unused space and reorganizing data for better performance.

#### 7.2.1 Defragment Behavior

- Removes unused space from deleted files
- Reorganizes file entries for optimal access
- Compacts data sections to reduce file size
- Updates internal indexes and references
- Preserves all package metadata and signatures
- May take significant time for large packages

#### 7.2.2 Defragment Error Conditions

- **Validation Errors**: Package not open, read-only mode
- **I/O Errors**: File system errors during defragmentation
- **Context Errors**: Context cancellation or timeout exceeded

#### 7.2.3 Defragment Example Usage

```go
err := package.Defragment(ctx)
if err != nil {
    return err
}
```

### 7.3 Package Information

```go
// GetInfo gets comprehensive package information
func (p *Package) GetInfo(ctx context.Context) PackageInfo
```

This function retrieves comprehensive information about the current package.

Returns a `PackageInfo` structure containing:

- Basic package information (file count, sizes)
- Package identity (VendorID, AppID)
- Package comment and metadata
- Digital signature information
- Security and compression status
- Timestamps and feature flags

#### 7.3.1 GetInfo Error Conditions

- **Validation Errors**: Package not currently open
- **Context Errors**: Context cancellation or timeout exceeded

#### 7.3.2 GetInfo Example Usage

```go
// Get comprehensive package information
info := package.GetInfo(ctx)
fmt.Printf("Package has %d files\n", info.FileCount)
fmt.Printf("Package version: %d\n", info.Version)
```

### 7.4 Header Inspection

```go
// ReadHeader reads the package header from the reader
func ReadHeader(ctx context.Context, reader io.Reader) (*Header, error)
```

This is a low-level function for header-only inspection without opening the full package.

#### 7.4.1 ReadHeader Use Cases

- Validate .npk file format without loading package data
- Inspect package metadata before deciding to open
- Debugging corrupted or partially readable packages
- Stream processing where only header information is needed

#### 7.4.2 ReadHeader Parameters

- `ctx`: Context for cancellation and timeout handling
- `reader`: Input stream to read header from

Returns a `Header` structure and error.

#### 7.4.3 ReadHeader Error Conditions

- **Validation Errors**: Invalid package header format
- **Unsupported Errors**: Unsupported package version
- **Context Errors**: Context cancellation or timeout exceeded

#### 7.4.4 ReadHeader Example Usage

```go
header, err := ReadHeader(ctx, file)
if err != nil {
    return err
}
```

### 7.5 Check Package State

```go
// IsOpen checks if the package is currently open
func (p *Package) IsOpen() bool

// IsReadOnly checks if the package is in read-only mode
func (p *Package) IsReadOnly() bool

// GetPath returns the current package file path
func (p *Package) GetPath() string
```

These functions provide information about the current package state.

## 8. Error Handling

### 8.1 Structured Error System

The NovusPack API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities.
For complete error system documentation, see [Structured Error System](api_core.md#10-structured-error-system).

### 8.2 Common Error Types

#### 8.2.1 Error Types Used

The NovusPack Basic Operations API uses the following error types from the structured error system:

- `ErrTypeValidation`: Input validation errors, invalid parameters, format errors
- `ErrTypeIO`: I/O errors, file system operations, network issues
- `ErrTypeSecurity`: Security-related errors, access denied, authentication failures
- `ErrTypeUnsupported`: Unsupported features, versions, or operations
- `ErrTypeContext`: Context cancellation, timeout, and lifecycle errors
- `ErrTypeCorruption`: Data corruption, checksum failures, integrity violations

### 8.3 Structured Error Examples

#### 8.3.1 Creating Structured Errors

```go
// Create a validation error with context
err := NewPackageError(ErrTypeValidation, "package file not found", nil).
    WithContext("path", "/path/to/package.npk").
    WithContext("operation", "Open")

// Wrap an existing error
err := WrapError(io.ErrUnexpectedEOF, ErrTypeIO, "unexpected end of file").
    WithContext("file", "package.npk").
    WithContext("offset", 1024)

// Create a security error
err := NewPackageError(ErrTypeSecurity, "permission denied", nil).
    WithContext("path", "/path/to/package.npk").
    WithContext("user", "anonymous")

// Create an I/O error
err := NewPackageError(ErrTypeIO, "failed to read package file", nil).
    WithContext("path", "/path/to/package.npk").
    WithContext("operation", "Open")
```

#### 8.3.2 Error Inspection

**Usage**: Check error types and handle them with appropriate logging and context extraction.

#### 8.3.3 Common Error Scenarios

**Usage**: Create structured errors with rich context for different error scenarios.

### 8.4 Error Handling Best Practices

#### 8.4.1 Always check for errors

Always check for errors after calling package operations and handle them appropriately.
Never ignore error return values as they indicate critical failures that must be addressed.

#### 8.4.2 Use structured errors for better debugging

Use the structured error system to provide rich context for debugging.
Wrap errors with additional context information to help identify the source of problems and provide better error messages to users.

#### 8.4.3 Use context for cancellation

Use context timeouts and cancellation to prevent operations from hanging indefinitely.
Set appropriate timeouts for long-running operations and handle context cancellation gracefully.

#### 8.4.4 Handle different error types appropriately

Handle different error types with appropriate responses.
Provide user-friendly messages for validation errors, log security errors, and implement retry logic for I/O errors.
Use the structured error system to determine the appropriate handling strategy.

#### 8.4.5 Clean up resources

Always clean up resources properly using defer statements.
Ensure packages are closed even when errors occur, and handle cleanup errors gracefully to prevent resource leaks.

## 9. Best Practices

### 9.1 Package Lifecycle Management

#### 9.1.1 Always use defer for cleanup

Always use defer statements to ensure resources are properly cleaned up, even when errors occur.
This prevents resource leaks and ensures consistent cleanup behavior.

#### 9.1.2 Check package state before operations

Always verify that a package is in the correct state before performing operations.
Check if the package is open, not read-only, and in a valid state for the intended operation.

#### 9.1.3 Use appropriate context timeouts

Use appropriate context timeouts for long-running operations to prevent indefinite blocking.
Set timeouts based on the expected operation duration and handle timeout errors gracefully.

### 9.2 Error Handling

#### 9.2.1 Wrap errors with context

Wrap errors with additional context information to provide better debugging information.
Include relevant details such as file paths, operation names, and parameter values in error messages.

#### 9.2.2 Handle specific error types

Handle different error types with appropriate responses based on the error category.
Use the structured error system to determine the appropriate handling strategy for validation, security, I/O, and other error types.

### 9.3 Resource Management

#### 9.3.1 Use context for resource management

Use context for resource management and cancellation.
Pass context to long-running operations and handle context cancellation to ensure proper resource cleanup and operation termination.

#### 9.3.2 Handle cleanup errors gracefully

Handle cleanup errors gracefully by logging warnings rather than failing.
Use defer functions to ensure cleanup occurs even when errors happen, and log cleanup failures as warnings rather than errors.
