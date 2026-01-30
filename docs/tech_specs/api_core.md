# NovusPack Technical Specifications - Core Package Interface API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
  - [0.2 Context Integration](#02-context-integration)
- [1. Core Interfaces](#1-core-interfaces)
  - [1.1 Package Interface](#11-package-interface)
    - [1.1.1 filePackage Struct](#111-filepackage-struct)
    - [1.1.2 Basic Operations](#112-basic-operations)
    - [1.1.3 File Management Operations](#113-file-management-operations)
  - [1.2 PackageReader Interface](#12-packagereader-interface)
    - [1.2.1 PackageReader Contract](#121-packagereader-contract)
    - [1.2.2 PackageReader.ReadFile Method](#122-packagereaderreadfile-method)
    - [1.2.3 PackageReader.ListFiles Method](#123-packagereaderlistfiles-method)
    - [1.2.4 FileInfo Structure](#124-fileinfo-structure)
    - [1.2.5 PackageReader.GetInfo Method](#125-packagereadergetinfo-method)
    - [1.2.6 PackageReader.GetMetadata Method](#126-packagereadergetmetadata-method)
    - [1.2.7 PackageReader.Validate Method](#127-packagereadervalidate-method)
    - [1.2.8 PackageReader Common Error Mapping Table](#128-packagereader-common-error-mapping-table)
  - [1.3 PackageWriter Interface](#13-packagewriter-interface)
    - [1.3.1 Memory Versus Disk Side Effects](#131-memory-versus-disk-side-effects)
    - [1.3.2 Common Writer Error Mapping Table](#132-common-writer-error-mapping-table)
- [2. Package Path Semantics](#2-package-path-semantics)
  - [2.1 Path Normalization Rules](#21-path-normalization-rules)
    - [2.1.1 Separator Normalization](#211-separator-normalization)
    - [2.1.2 Leading Slash Requirement](#212-leading-slash-requirement)
    - [2.1.3 Dot Segment Canonicalization](#213-dot-segment-canonicalization)
    - [2.1.4 Unicode Normalization](#214-unicode-normalization)
    - [2.1.5 Path Length Limits](#215-path-length-limits)
    - [2.1.6 Path Normalization on Storage](#216-path-normalization-on-storage)
  - [2.2 Path Rules](#22-path-rules)
    - [2.2.1 Case Sensitivity](#221-case-sensitivity)
  - [2.3 Path Display and Extraction](#23-path-display-and-extraction)
- [3. Package Writing Operations](#3-package-writing-operations)
- [4. File Management](#4-file-management)
- [5. Encryption Management](#5-encryption-management)
- [6. Package Compression Operations](#6-package-compression-operations)
- [7. Digital Signatures and Security](#7-digital-signatures-and-security)
  - [7.1 Core Integration Points](#71-core-integration-points)
  - [7.2 Write Protection and Immutability Enforcement](#72-write-protection-and-immutability-enforcement)
- [8. Package Metadata Management](#8-package-metadata-management)
  - [8.1 General Metadata Operations](#81-general-metadata-operations)
  - [8.2 AppID/VendorID Management](#82-appidvendorid-management)
  - [8.3 Package Information Structures](#83-package-information-structures)
- [9. File Validation Requirements](#9-file-validation-requirements)
- [10. Structured Error System](#10-structured-error-system)
  - [10.1 Benefits of Structured Errors](#101-benefits-of-structured-errors)
  - [10.2 ErrorType Types and Categories](#102-errortype-types-and-categories)
  - [10.3 ErrorType Categories](#103-errortype-categories)
  - [10.4 PackageError Structure](#104-packageerror-structure)
    - [10.4.1 PackageError.Error Method](#1041-packageerrorerror-method)
    - [10.4.2 PackageError.Unwrap Method](#1042-packageerrorunwrap-method)
    - [10.4.3 PackageError.Is Method](#1043-packageerroris-method)
  - [10.5 Error Helper Functions](#105-error-helper-functions)
    - [10.5.1 NewPackageError Function](#1051-newpackageerror-function)
    - [10.5.2 WrapErrorWithContext Function](#1052-wraperrorwithcontext-function)
    - [10.5.3 AsPackageError Function](#1053-aspackageerror-function)
    - [10.5.4 GetErrorContext Function](#1054-geterrorcontext-function)
    - [10.5.5 AddErrorContext Function](#1055-adderrorcontext-function)
    - [10.5.6 MapError Function](#1056-maperror-function)
    - [10.5.7 Example - Creating Errors with Context](#1057-example---creating-errors-with-context)
    - [10.5.8 Example - Error Inspection and Handling](#1058-example---error-inspection-and-handling)
    - [10.5.9 Example - Error Propagation](#1059-example---error-propagation)
    - [10.5.10 Error Example Pattern](#10510-error-example-pattern)
    - [10.5.11 Error Wrapping Patterns](#10511-error-wrapping-patterns)
    - [10.5.12 Error Logging and Debugging](#10512-error-logging-and-debugging)
- [11. Generic Types](#11-generic-types)
- [12. Utility Functions](#12-utility-functions)
  - [12.1 NormalizePackagePath Function](#121-normalizepackagepath-function)
    - [12.1.1 NormalizePackagePath Error Handling](#1211-normalizepackagepath-error-handling)
    - [12.1.2 NormalizePackagePath Return Value](#1212-normalizepackagepath-return-value)
  - [12.2 ToDisplayPath Function](#122-todisplaypath-function)
    - [12.2.1 ToDisplayPath Behavior](#1221-todisplaypath-behavior)
    - [12.2.2 ToDisplayPath Usage](#1222-todisplaypath-usage)
  - [12.3 ValidatePackagePath Function](#123-validatepackagepath-function)
  - [12.4 ValidatePathLength Function](#124-validatepathlength-function)

---

## 0. Overview

This document defines the core package interface for the NovusPack system, including basic operations, file management, encryption, and metadata handling.

Signature management and signature validation are deferred to v2.
V1 enforces signed package immutability based on signature presence and does not validate signature contents.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Basic Operations API](api_basic_operations.md) - Package creation, opening, closing, and lifecycle management
- [File Management API](api_file_mgmt_index.md) - File operations and encryption
- [Package Writing Operations](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [Package Compression API](api_package_compression.md) - Package compression and decompression operations
- [Streaming and Buffer Management](api_streaming.md) - File streaming interface and buffer management system
- [Multi-Layer Deduplication](api_deduplication.md) - Content deduplication strategies and processing levels
- [Digital Signature API](api_signatures.md) - Signature management, types, and validation (deferred to v2)
- [Package Metadata API](api_metadata.md) - Comment management, AppID/VendorID, and metadata operations
- [Security Validation API](api_security.md) - Package validation and security status structures
- [Generic Types and Patterns](api_generics.md) - Generic type system, patterns, and best practices
- [File Format Specifications](package_file_format.md) - .nvpk format structure and optional signature binary format
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation

### 0.2 Context Integration

The NovusPack API follows 2025 Go best practices for context integration. Methods that perform I/O operations, network calls, or long-running operations accept `context.Context` as the first parameter to support:

- Request cancellation and timeout handling
- Request-scoped values and configuration
- Graceful shutdown and resource cleanup
- Integration with Go's standard context patterns

Pure data structure methods that perform only in-memory operations do not require `context.Context`. These are synchronous, deterministic operations with no I/O, network calls, or long-running behavior. Examples include:

- `Validate()` methods on data structures (PackageHeader, FileEntry, etc.)
- `Size()` calculation methods
- `Get*()`, `Set*()`, `Has*()`, `Clear*()` accessor methods
- Flag manipulation methods (`SetFlag`, `ClearFlag`, `HasFlag`)

These methods operate exclusively on in-memory data structures and do not perform I/O operations, making context unnecessary.

## 1. Core Interfaces

The NovusPack API is designed around core interfaces that provide clear separation of concerns and enable better testability.

### 1.1 Package Interface

The `Package` interface provides a unified API that combines:

- **PackageReader** methods (embedded) - Read-only operations on opened packages
- **PackageWriter** methods (embedded) - Write operations to persist package changes
- **Lifecycle operations** - Package creation, opening, closing, and state management
- **File management operations** - Adding, removing, and extracting files (see [File Management API](api_file_mgmt_index.md))
- **Metadata operations** - Comment, AppID, and VendorID management (see [Package Metadata API](api_metadata.md))
- **Compression operations** - Package-level compression and decompression (see [Package Compression API](api_package_compression.md))
- **Session base management** - Automatic path derivation for file operations

**Note on Embedded Interface Methods in Go**: When `Package` embeds `PackageReader` and `PackageWriter`, the methods from those interfaces become part of the `Package` interface. A single implementation satisfies both the embedded interface and the `Package` interface.
For example, `Package.Validate` and `PackageReader.Validate` are the same method - there is no delegation or wrapper.
The concrete `filePackage` type implements one `Validate` method that satisfies both interfaces.

```go
// Package defines the main interface for NovusPack package operations.
// Package combines PackageReader and PackageWriter interfaces, providing
// complete package lifecycle management including opening, closing, and
// defragmentation operations.
type Package interface {
    // Embedded interfaces
    PackageReader
    PackageWriter

    // Lifecycle operations
    SetTargetPath(ctx context.Context, path string) error
    Close() error
    CloseWithCleanup(ctx context.Context) error
    IsOpen() bool
    IsReadOnly() bool
    GetPath() string
    Defragment(ctx context.Context) error

    // Session base path management
    // See [File Addition API - Session Base](api_file_mgmt_addition.md#264-session-base---package-level-automatic-basepath)
    SetSessionBase(basePath string) error
    GetSessionBase() string
    ClearSessionBase()
    HasSessionBase() bool

    // File management operations
    // See [File Management API](api_file_mgmt_index.md) for detailed specifications
    AddFile(ctx context.Context, sourcePath string, opts *AddFileOptions) (*FileEntry, error)
    AddFileFromMemory(ctx context.Context, path string, data []byte, opts *AddFileOptions) (*FileEntry, error)
    AddFilePattern(ctx context.Context, pattern string, opts *AddFileOptions) ([]*FileEntry, error)
    AddDirectory(ctx context.Context, sourcePath string, opts *AddFileOptions) ([]*FileEntry, error)
    RemoveFile(ctx context.Context, path string) error
    RemoveFilePattern(ctx context.Context, pattern string) ([]string, error)
    ExtractPath(ctx context.Context, storedPath string, isWindows bool, opts *ExtractPathOptions) error


    // Package metadata management
    // See [Package Metadata API](api_metadata.md) for detailed specifications
    SetComment(comment string) error
    GetComment() string
    ClearComment() error
    HasComment() bool
    SetAppID(appID uint64) error
    GetAppID() uint64
    ClearAppID() error
    HasAppID() bool
    SetVendorID(vendorID uint32) error
    GetVendorID() uint32
    ClearVendorID() error
    HasVendorID() bool
    SetPackageIdentity(vendorID uint32, appID uint64) error
    GetPackageIdentity() (uint32, uint64)
    ClearPackageIdentity() error

    // Package compression operations
    // See [Package Compression API](api_package_compression.md) for detailed specifications
    CompressPackage(ctx context.Context) error
    DecompressPackage(ctx context.Context) error
    CompressPackageFile(ctx context.Context, outputPath string) error
    DecompressPackageFile(ctx context.Context, outputPath string) error
    GetPackageCompressionInfo() (*PackageCompressionInfo, error)
    IsPackageCompressed() bool
    GetPackageCompressionType() uint8
    SetPackageCompressionType(compressionType uint8) error
    CanCompressPackage() bool
}
```

**Implementation**: The `Package` interface is implemented by the concrete `filePackage` struct.
See [FilePackage Implementation](#111-filepackage-struct) for the complete struct definition and field descriptions.

#### 1.1.1 filePackage Struct

The NovusPack API uses an interface-based design where the `Package` interface is implemented by the concrete `filePackage` struct.

**Source of Truth**: The following `filePackage` struct definition is the authoritative specification for the internal package implementation:

```go
// filePackage is the concrete implementation of the Package interface
type filePackage struct {
    // Exported fields (accessible via interface methods)
    Info *metadata.PackageInfo

    FileEntries []*metadata.FileEntry

    PathMetadataEntries []*metadata.PathMetadataEntry

    SpecialFiles map[uint16]*metadata.FileEntry

    FilePath string

    // Internal state (unexported fields)
    header      *fileformat.PackageHeader
    index       *fileformat.FileIndex
    fileHandle  *os.File
    isOpen      bool
    sessionBase string  // Package-level session base path (runtime-only, not persisted)
}
```

##### 1.1.1.1 FilePackage Field Descriptions

- `Info`: Package metadata including file counts, timestamps, and security information.
  - Populated during package opening or creation.
  - Updated as package state changes.

- `FileEntries`: Slice of all file entries in the package.
  - Loaded on demand to minimize memory usage.
  - Initially empty when package is opened.
  - Populated when file operations require entry data.

- `PathMetadataEntries`: Slice of path metadata entries.
  - Loaded from special metadata files (type 65001).
  - Contains path-specific metadata including permissions, ownership, and tags.
  - Used for tag inheritance and filesystem property management.

- `SpecialFiles`: Map of special file type IDs to their FileEntry instances.
  - Keys are file type identifiers (65000-65535).
  - Values are FileEntry instances for metadata, signatures, and other special files.
  - Loaded during package opening.

- `FilePath`: Path to the package file on disk.
  - Set during NewPackageWithOptions() or OpenPackage() operations.
  - Used for file I/O operations.

- `header`: Binary package header structure.
  - Contains magic number, version, flags, and offsets.
  - Loaded immediately when package is opened.
  - Used for package format validation and navigation.

- `index`: File index structure.
  - Contains entry count and list of IndexEntry structures.
  - Each IndexEntry maps FileID to FileEntry offset.
  - Loaded immediately when package is opened.
  - Used to locate FileEntry objects in the package file.

- `fileHandle`: Open file handle for reading package data.
  - Opened during OpenPackage() operation.
  - Closed during Close() operation.
  - Nil when package is not open.

- `isOpen`: Boolean flag indicating if package is open for reading.
  - Set to true during OpenPackage() operation.
  - Set to false during Close() operation.
  - Used to validate operation state - if false, package is considered closed.
  - All read operations (GetInfo, GetMetadata, ListFiles, ReadFile, Validate) require isOpen to be true.
  - Close() checks this flag to make the operation idempotent (multiple Close() calls are safe).

- `sessionBase`: Package-level session base path for automatic path derivation.
  - Runtime-only field, not persisted to package file.
  - Used for automatic BasePath determination when adding files.
    - See [Session Base Management](api_basic_operations.md#19-package-session-base-management) for details.
  - Can be set explicitly via SetSessionBase() or established automatically from first absolute path.
  - Cleared via ClearSessionBase().

#### 1.1.2 Basic Operations

See [Basic Operations API](api_basic_operations.md) for detailed method signatures and implementation details.

- **NewPackage**: Creates a new empty package
- **NewPackageWithOptions**: Creates a new package with specified configuration options
- **SetTargetPath**: Changes the package's target write path
- **OpenPackage**: Opens an existing package from the specified path
- **OpenPackageReadOnly**: Opens an existing package from the specified path in read-only mode
- **OpenBrokenPackage**: Opens a broken package for repair workflows
- **Close**: Closes the package and releases resources
- **CloseWithCleanup**: Closes the package and performs cleanup operations
- **Validate**: Validates package format, structure, and integrity
- **Defragment**: Optimizes package structure and removes unused space
- **GetInfo**: Gets comprehensive package information
- **Write**: General write method with compression handling using SafeWrite or FastWrite methods
  - See [Package Writing API](api_writing.md) for detailed method signatures and implementation details.

#### 1.1.3 File Management Operations

File management operations (adding, removing files):

- `AddFile` - Adds a file from filesystem
- `AddFileFromMemory` - Adds a file from memory
- `AddFilePattern` - Adds files matching a pattern
- `AddDirectory` - Recursively adds directory contents
- `RemoveFile` - Removes a file
- `RemoveFilePattern` - Removes files matching a pattern

For complete documentation, see [File Management API](api_file_mgmt_index.md).

### 1.2 PackageReader Interface

```go
// PackageReader defines the interface for reading operations on a package.
// PackageReader provides methods for reading files, listing files, retrieving
// metadata, and validating package contents.
type PackageReader interface {
    ReadFile(ctx context.Context, path string) ([]byte, error)
    ListFiles() ([]FileInfo, error)
    GetMetadata() (*PackageMetadata, error)
    Validate(ctx context.Context) error
    GetInfo() (*PackageInfo, error)
}
```

#### 1.2.1 PackageReader Contract

`PackageReader` is the read-only interface for opened packages.
It is embedded in the larger `Package` interface and describes read-only operations available on an opened package instance.

##### 1.2.1.1 Reader Contract Scope

`PackageReader` methods assume the package has already been opened (via `OpenPackage` or equivalent).
`PackageReader` is not intended to represent header-only or lightweight on-disk inspection.
Header-only inspection is handled by separate functions (for example, `ReadHeader` in [Basic Operations API](api_basic_operations.md)) rather than by a special `PackageReader` implementation.

##### 1.2.1.2 OpenPackage Eager Metadata Load

`OpenPackage` MUST read into memory all package metadata required for `PackageReader` operations, including:

- Package header
- File index
- All FileEntry metadata
- Package comment metadata
- Signature metadata
- Special metadata files (type 65000-65535)
- Path metadata YAML parsed into `PathMetadataEntries` (fully eager)
- File-path associations updated during open

##### 1.2.1.3 Context Usage

Pure in-memory operations do not require `context.Context`:

- `ListFiles()` - Pure in-memory operation that returns file information from already-loaded metadata
- `GetInfo()` - Pure in-memory operation that returns lightweight package information
- `GetMetadata()` - Pure in-memory operation that returns comprehensive metadata from already-loaded data

Operations that perform I/O retain `context.Context`:

- `ReadFile(ctx, path)` - Reads file data from disk, may decrypt/decompress, and should be cancellable
- `Validate(ctx)` - Validation can be non-trivial and should be cancellable

##### 1.2.1.4 Unsupported Operations

`PackageReader` should not include methods that are only meaningful for un-opened, on-disk inspection.

##### 1.2.1.5 Code Reuse Requirement

Operations that share underlying functionality (for example, header reading and validation) must use shared helper functions to avoid duplication.
For example, `internal.ReadAndValidateHeader` is shared by both `ReadHeader` (standalone function) and `OpenPackage` (full package load).

##### 1.2.1.6 ReadFile Cross-Reference

See [ReadHeader](api_basic_operations.md) function as an example of a lightweight operation for header-only inspection.

#### 1.2.2 PackageReader.ReadFile Method

Reads file content from the package, applying decryption and decompression as needed.

The canonical signature for `ReadFile` is defined in the [PackageReader Interface](#12-packagereader-interface).

##### 1.2.2.1 PackageReader.ReadFile Parameters

- `ctx context.Context` - Context for cancellation and timeout handling
- `path string` - Package-internal path to the file (see [Package Path Semantics](#2-package-path-semantics))

##### 1.2.2.2 PackageReader.ReadFile Returns

- `[]byte` - File content (decrypted and decompressed)
- `error` - Returns `*PackageError` on failure

##### 1.2.2.3 PackageReader.ReadFile Behavior

- Reads file data from disk
- Locates the FileEntry by normalized package path.
- Uses FileEntry.EntryOffset when available to locate file data efficiently (file data starts at EntryOffset + TotalSize()).
- If FileEntry.EntryOffset is not available, uses the file index to locate file data.
- Applies decryption if the file is encrypted
- Applies decompression if the file is compressed
- Returns decrypted and decompressed content

##### 1.2.2.4 PackageReader.ReadFile Error Conditions

See [Common Error Mapping Table](#128-packagereader-common-error-mapping-table).

##### 1.2.2.5 PackageReader.ReadFile Concurrency

Safe for concurrent reads from different goroutines.

#### 1.2.3 PackageReader.ListFiles Method

Returns information about all files in the package.

The canonical signature for `ListFiles` is defined in the [PackageReader Interface](#12-packagereader-interface).

##### 1.2.3.1 PackageReader.ListFiles Parameters

None (pure in-memory operation).

##### 1.2.3.2 PackageReader.ListFiles Returns

- `[]FileInfo` - Slice of file information, sorted by PrimaryPath alphabetically
- `error` - Returns `*PackageError` on failure

##### 1.2.3.3 PackageReader.ListFiles Behavior

- Results MUST be sorted by PrimaryPath (normalized package path), alphabetically
- Results MUST be stable across calls when the in-memory package state has not changed
- Files that have been removed from the in-memory Package MUST NOT be listed by `ListFiles()`
- In-memory mutations (for example, via `AddFile` or `RemoveFile`) affect `ListFiles` results immediately, even before a write operation
- For files with multiple paths, PrimaryPath is the first path (lexicographically) and all paths appear in the Paths array

##### 1.2.3.4 PackageReader.ListFiles Error Conditions

See [Common Error Mapping Table](#128-packagereader-common-error-mapping-table).

##### 1.2.3.5 PackageReader.ListFiles Concurrency

Safe for concurrent calls from different goroutines.

#### 1.2.4 FileInfo Structure

`FileInfo` provides lightweight file information for package contents, optimized for listing operations and filtering without requiring full `FileEntry` access.

```go
// FileInfo provides lightweight file information for package contents.
// FileInfo is optimized for listing operations and filtering without requiring full FileEntry access.
type FileInfo struct {
    // Basic Identification
    PrimaryPath  string   // Primary display path (leading '/' removed)
    Paths        []string // All paths for this file (aliases/hard links, leading '/' removed)
    FileID       uint64   // Unique file identifier
    FileType     uint16   // File type identifier (0-64999: content, 65000-65535: special metadata)
    FileTypeName string   // Human-readable file type name (e.g., "Texture", "Audio", "Unknown")

    // Size Information
    Size       int64 // Original file size in bytes (before compression/encryption)
    StoredSize int64 // Actual stored size in bytes (after compression/encryption)

    // Processing Status
    IsCompressed   bool  // Whether file is compressed
    IsEncrypted    bool  // Whether file is encrypted
    CompressionType uint8 // Compression algorithm (0=none, 1=Zstd, 2=LZ4, 3=LZMA)

    // Content Verification
    RawChecksum    uint32 // CRC32 checksum of original content
    StoredChecksum uint32 // CRC32 checksum of stored content (after compression/encryption)

    // Multi-Path Support
    PathCount uint16 // Number of paths (aliases) for this FileEntry

    // Version Tracking
    FileVersion     uint32 // File content version
    MetadataVersion uint32 // File metadata version

    // Metadata Indicators
    HasTags bool // Whether file has custom tags/metadata
}
```

##### 1.2.4.1 FileInfo Field Descriptions

Basic Identification fields:

- `PrimaryPath`: Primary display path with leading `/` stripped (first path from Paths array)
- `Paths`: All paths referencing this file content (aliases/hard links), with leading `/` stripped
- `FileID`: Unique 64-bit identifier, stable across package operations
- `FileType`: File type from type system (0-64999 for content, 65000-65535 for special files)
- `FileTypeName`: Human-readable file type name for display (derived from FileType via type system lookup)

Size Information fields:

- `Size`: Original uncompressed, unencrypted size
- `StoredSize`: Actual bytes consumed in package file

Processing Status fields:

- `IsCompressed`: Quick check for compression without examining type
- `IsEncrypted`: Quick check for encryption status
- `CompressionType`: Specific algorithm used (0=none, 1=Zstd, 2=LZ4, 3=LZMA)

Content Verification fields:

- `RawChecksum`: CRC32 checksum of original unprocessed content (for deduplication and verification)
- `StoredChecksum`: CRC32 checksum of stored content after compression/encryption (for integrity verification)

Multi-Path Support fields:

- `PathCount`: Number of paths referencing this content (1=single path, 2+=hard links/aliases)

Version Tracking fields:

- `FileVersion`: Increments when file content changes
- `MetadataVersion`: Increments when file metadata changes

Metadata Indicators fields:

- `HasTags`: Indicates if file has custom tags (requires `GetFileByPath` for full tag data)

##### 1.2.4.2 FileInfo Design Rationale

`FileInfo` balances comprehensiveness with performance:

- **Lightweight**: All fields from `FileEntry` static section (no variable-length data)
- **Filtering-friendly**: Supports common filtering operations (type, compression, encryption)
- **Drill-down enabled**: Provides enough information to decide if full `FileEntry` is needed
- **Query alignment**: Matches existing query functions (`ListCompressedFiles`, `FindEntriesByType`)

##### 1.2.4.3 Example - FileInfo Usage - Basic Listing

```go
// Example: Basic listing
files, err := pkg.ListFiles()
for _, f := range files {
    fmt.Printf("%s (%s): %d bytes\n", f.PrimaryPath, f.FileTypeName, f.Size)
    if len(f.Paths) > 1 {
        fmt.Printf("  Aliases: %v\n", f.Paths[1:])
    }
}
```

##### 1.2.4.4 Example - FileInfo Usage - Filter by Type

```go
// Example: Filter by numeric type range
textures := lo.Filter(files, func(f FileInfo, _ int) bool {
    return f.FileType >= 1000 && f.FileType < 2000 // Texture types
})

// Example: Filter by type name
audioFiles := lo.Filter(files, func(f FileInfo, _ int) bool {
    return f.FileTypeName == "Audio"
})
```

##### 1.2.4.5 Example - FileInfo Usage - Find Compressed Files

```go
// Example: Find compressed files
compressed := lo.Filter(files, func(f FileInfo, _ int) bool {
    return f.IsCompressed
})
```

##### 1.2.4.6 Example - FileInfo Usage - Calculate Compression Ratios

```go
// Example: Calculate compression ratios
for _, f := range files {
    if f.IsCompressed {
        ratio := float64(f.StoredSize) / float64(f.Size)
        fmt.Printf("%s: %.2f%% of original\n", f.Path, ratio*100)
    }
}
```

##### 1.2.4.7 Example - FileInfo Usage - Check for Duplicates

```go
// Example: Deduplicate by original content
unique := lo.UniqBy(files, func(f FileInfo) uint32 {
    return f.RawChecksum
})
```

##### 1.2.4.8 Example - FileInfo Usage - Verify Content Integrity

```go
// Example: Check if stored content matches expected checksum
for _, f := range files {
    if f.IsCompressed || f.IsEncrypted {
        // StoredChecksum verifies compressed/encrypted data integrity
        fmt.Printf("%s: stored checksum 0x%08x\n", f.PrimaryPath, f.StoredChecksum)
    } else {
        // For unprocessed files, both checksums should match
        if f.RawChecksum != f.StoredChecksum {
            fmt.Printf("Warning: checksum mismatch for %s\n", f.PrimaryPath)
        }
    }
}
```

##### 1.2.4.9 Example - FileInfo Usage - Find Files with Multiple Paths

```go
// Example: Find files with multiple paths
aliased := lo.Filter(files, func(f FileInfo, _ int) bool {
    return len(f.Paths) > 1
})
for _, f := range aliased {
    fmt.Printf("File %d has %d paths: %v\n", f.FileID, len(f.Paths), f.Paths)
}
```

##### 1.2.4.10 FileInfo Usage - Collection Example Patterns

When working with collections of results (for example `[]FileInfo` from `ListFiles`, or `[]*FileEntry` from other APIs), use `samber/lo` for common operations such as find, filter, map, and de-duplication.

See [samber/lo Usage Standards](../implementations/go/samber_lo_usage.md) for detailed guidance.

```go
import "github.com/samber/lo"

// Example: Find first entry matching predicate.
targetID := uint64(123)
entry, found := lo.Find(files, func(f FileInfo) bool {
    return f.FileID == targetID
})
_ = entry
_ = found

// Example: Filter entries matching predicate.
filtered := lo.Filter(files, func(f FileInfo, _ int) bool {
    return f.IsCompressed
})
_ = filtered

// Example: Map entries to another type.
paths := lo.Map(files, func(f FileInfo, _ int) string {
    return f.PrimaryPath
})
_ = paths

// Example: Check for duplicates by raw checksum.
unique := lo.UniqBy(files, func(f FileInfo) uint32 {
    return f.RawChecksum
})
_ = unique
```

#### 1.2.5 PackageReader.GetInfo Method

Returns lightweight package information derived from header and computed package-level statistics.

The canonical signature for `GetInfo` is defined in the [PackageReader Interface](#12-packagereader-interface).

##### 1.2.5.1 PackageReader.GetInfo Parameters

None (pure in-memory operation).

##### 1.2.5.2 PackageReader.GetInfo Returns

- `*PackageInfo` - Lightweight package information (see [PackageInfo Structure](api_metadata.md#71-packageinfo-structure))
- `error` - Returns `*PackageError` on failure

##### 1.2.5.3 PackageReader.GetInfo Scope

Lightweight view over already-loaded package state.
This method MUST NOT perform additional disk I/O.
This method MUST NOT perform additional parsing beyond what `OpenPackage` already loaded.

##### 1.2.5.4 PackageReader.GetInfo Error Conditions

See [Common Error Mapping Table](#128-packagereader-common-error-mapping-table).

##### 1.2.5.5 PackageReader.GetInfo Concurrency

- Go: Safe for concurrent calls from different goroutines.

#### 1.2.6 PackageReader.GetMetadata Method

Returns comprehensive metadata including all package information plus detailed file and metadata file contents.

The canonical signature for `GetMetadata` is defined in the [PackageReader Interface](#12-packagereader-interface).

##### 1.2.6.1 PackageReader.GetMetadata Parameters

None (pure in-memory operation).

##### 1.2.6.2 PackageReader.GetMetadata Returns

- `*PackageMetadata` - Comprehensive package metadata; see [Package Metadata API - PackageInfo Structure](api_metadata.md#71-packageinfo-structure)
- `error` - Returns `*PackageError` on failure

##### 1.2.6.3 PackageReader.GetMetadata Scope

Full metadata view over already-loaded package state.
This method MUST NOT perform additional disk I/O.
This method MUST NOT perform additional parsing beyond what `OpenPackage` already loaded.

##### 1.2.6.4 PackageReader.GetMetadata Serialization

See [Package Metadata API - Package Information Methods](api_metadata.md#74-package-information-methods).

##### 1.2.6.5 PackageReader.GetMetadata Error Conditions

`GetMetadata()` MUST still return an error for internal consistency failures (for example, metadata not loaded due to an invariant violation) even though `OpenPackage` is required to eagerly load metadata.
`GetMetadata()` requires a fully opened package instance (it is not applicable to header-only inspection).

##### 1.2.6.6 PackageReader.GetMetadata Concurrency

Safe for concurrent calls from different goroutines.

#### 1.2.7 PackageReader.Validate Method

Validates package format, structure, and integrity.

The canonical signature for `Validate` is defined in the [PackageReader Interface](#12-packagereader-interface).

##### 1.2.7.1 PackageReader.Validate Parameters

- `ctx context.Context` - Context for cancellation and timeout handling

##### 1.2.7.2 PackageReader.Validate Returns

- `error` - Returns `*PackageError` on failure

##### 1.2.7.3 PackageReader.Validate Behavior

- Performs comprehensive package validation
- Can be non-trivial and should be cancellable
- Validates format, structure, checksums, and integrity

##### 1.2.7.4 PackageReader.Validate Error Conditions

See [Common Error Mapping Table](#128-packagereader-common-error-mapping-table).

##### 1.2.7.5 PackageReader.Validate Concurrency

- Go: Safe for concurrent calls from different goroutines.

#### 1.2.8 PackageReader Common Error Mapping Table

The following error mapping applies to all `PackageReader` methods:

| Condition                              | Error Type          | Notes                                                                                           |
| -------------------------------------- | ------------------- | ----------------------------------------------------------------------------------------------- |
| Package is not open / is closed        | `ErrTypeValidation` | Must be consistent across all reader methods.                                                   |
| Invalid package-internal path          | `ErrTypeValidation` | Includes empty path, malformed path, or path normalization that would escape package root.      |
| File not found at path                 | `ErrTypeValidation` | Use a consistent message pattern.                                                               |
| Package integrity check failed         | `ErrTypeCorruption` | Integrity failures are corruption.                                                              |
| Context cancelled or deadline exceeded | `ErrTypeContext`    | Applies only to methods that accept `context.Context` (for example, `ReadFile` and `Validate`). |

### 1.3 PackageWriter Interface

```go
// PackageWriter defines the interface for writing packages to disk.
type PackageWriter interface {
    Write(ctx context.Context) error
    SafeWrite(ctx context.Context, overwrite bool) error
    FastWrite(ctx context.Context) error
}
```

#### 1.3.1 Memory Versus Disk Side Effects

The `PackageWriter` interface provides methods for writing the in-memory package state to disk.

File management operations (add, remove) are part of the `Package` interface, not `PackageWriter`.
See [File Management API](api_file_mgmt_index.md) for file operations.

##### 1.3.1.1 Write Operations

- `Write` - General write method that selects appropriate strategy (SafeWrite or FastWrite)
- `SafeWrite` - Atomic write with temp file strategy (see [Package Writing API](api_writing.md#1-safewrite---atomic-package-writing))
- `FastWrite` - In-place updates for existing packages (see [Package Writing API](api_writing.md#2-fastwrite---in-place-package-updates))

##### 1.3.1.2 Write Durability

- Changes are NOT written to disk until `Write`, `SafeWrite`, or `FastWrite` is called
- File management operations (`AddFile`, `RemoveFile`, etc.) mutate in-memory state
- The in-memory package state reflects all pending changes

##### 1.3.1.3 Writing - Target Path Configuration

The package's target path is configured in memory using:

- [`NewPackageWithOptions`](api_basic_operations.md#7-newpackagewithoptions-function) - For initial package creation with configuration options
- [`SetTargetPath`](api_basic_operations.md#8-packagesettargetpath-method) - To change the path on an existing package

`SafeWrite` and `FastWrite` write to the Package's configured target path.
Writing to a new path requires changing the Package's configured target path (via `SetTargetPath`) prior to calling `SafeWrite` or `FastWrite`.

For detailed information about allowed target paths, overwrite behavior, path restrictions, and signed package handling, see the [Package Writing API](api_writing.md).

#### 1.3.2 Common Writer Error Mapping Table

The following error mapping applies to all `PackageWriter` methods:

| Condition                                                                         | Error Type        | Notes                                                                       |
| --------------------------------------------------------------------------------- | ----------------- | --------------------------------------------------------------------------- |
| Package is not open / is closed                                                   | ErrTypeValidation | Must be consistent across all writer methods.                               |
| Read-only package                                                                 | ErrTypeSecurity   | Mutations and writes must be rejected with a structured security error.     |
| Invalid package-internal path                                                     | ErrTypeValidation | Applies to file path parameters (for example, `RemoveFile`, `ExtractPath`). |
| Overwrite disallowed (target exists, overwrite == false)                          | ErrTypeValidation | Applies to `SafeWrite`.                                                     |
| Attempt to overwrite a signed package                                             | ErrTypeSecurity   | Overwriting signed packages is prohibited.                                  |
| Attempt to write signed package content to a new path without clearing signatures | ErrTypeSecurity   | Changing the target path MUST clear signature information before writing.   |
| Permission denied on target directory or file                                     | ErrTypeSecurity   | Use `ErrTypeSecurity` for permission failures.                              |
| I/O failure during write                                                          | ErrTypeIO         | Includes short writes, fsync failures, rename failures, etc.                |
| FastWrite interrupted / corruption detected                                       | ErrTypeCorruption | Corruption is expected risk of FastWrite interruption.                      |

## 2. Package Path Semantics

Package-internal paths are treated like archive entry names (tar/zip style), not OS filesystem paths.

For path utility functions that implement these rules, see [Utility Functions](#12-utility-functions).

### 2.1 Path Normalization Rules

Package paths are normalized according to the following rules when stored in the package format:

#### 2.1.1 Separator Normalization

All path separators are converted to forward slashes (`/`), regardless of the source platform.

- Windows-style backslashes (`\`) are converted to forward slashes (`/`)
- Example: `Users\file.txt` becomes `/Users/file.txt`
- Example: `/home/user/file.txt` remains `/home/user/file.txt`

All paths are stored using Unix-style forward slashes (`/`) as separators, regardless of the source platform.
This provides cross-platform compatibility and consistent path representation within the package format.

#### 2.1.2 Leading Slash Requirement

All stored paths MUST have a leading `/` (added if missing during normalization).

- The leading `/` indicates the package root, not the OS filesystem root
- Root path is represented as `/`
- Example: `path/to/file.txt` becomes `/path/to/file.txt`
- Example: `/path/to/file.txt` remains `/path/to/file.txt`

All paths MUST be stored with a leading `/` to ensure full path references within the package.
The leading `/` indicates the package root, not the OS filesystem root.
This provides unambiguous path representation and consistent path semantics within the package format.

All stored paths MUST begin with `/`:

1. **Root path**: `/` (single forward slash) - represents the package root directory
2. **All file paths**: Must have leading `/` - e.g., `/file.txt`, `/path/to/file.txt`
3. **All directory paths**: Must have leading `/` - e.g., `/assets/`, `/path/to/dir/`
4. **Normalization**: Input paths without leading slashes are automatically prefixed with `/` during normalization

Valid stored paths (all with leading `/`):

- `/file.txt` - file at package root
- `/path/to/file.txt` - file in nested directory
- `/assets/textures/ui/button.png` - deeply nested file
- `/usr/bin/myapp` - Unix-style path structure
- `/` - the package root directory itself

#### 2.1.3 Dot Segment Canonicalization

Dot segments (`.` and `..`) are resolved to canonical paths.

- `.` segments are removed (current directory references)
- `..` segments are resolved relative to the package root
- Dot segments that would escape the package root are rejected with an error
- Example: `/path/../file.txt` becomes `/file.txt`
- Example: `../../../etc/passwd` is rejected (would escape package root)

This ensures all paths are canonical and prevents path traversal attacks.

#### 2.1.4 Unicode Normalization

Path strings are normalized to NFC (Normalization Form Composed) before storage.

All path strings MUST be normalized to NFC (Normalization Form Composed) before storage.

The rationale for this normalization:

- Ensures consistent lookups across platforms
- Prevents duplicate entries for visually identical paths
- Resolves macOS (NFD) vs Windows/Linux (NFC) differences

Example:

- Input: `café` (can be represented as NFC or NFD in different filesystems)
- Stored: Always NFC form regardless of input
- Benefit: Lookups work consistently, deduplication works correctly

Implementation:

- Go: Use Go's `golang.org/x/text/unicode/norm` package

Cross-platform behavior:

- macOS typically uses NFD (decomposed)
- Windows/Linux typically use NFC (composed)
- NovusPack normalizes to NFC on storage
- Extraction preserves NFC form (may need conversion on macOS for proper display)

For complete Unicode normalization details, including rationale, examples, implementation, and cross-platform behavior, see [`NormalizePackagePath` Function - Unicode Normalization](#214-unicode-normalization).

#### 2.1.5 Path Length Limits

The path length is validated to ensure it is within format limits.

- **Format Limit**: The `PathEntry.PathLength` field is `uint16`, allowing paths up to **65,535 bytes**
- Paths exceeding this limit are rejected with an error
- Paths within the limit are accepted without artificial restrictions

**Storage Policy**: Accept any path up to the format limit (65,535 bytes) without artificial restrictions.

**Warning Thresholds**: The implementation MUST emit warnings during file addition operations when paths exceed platform-specific limits:

1. **> 260 bytes** (Info): "Path exceeds Windows default limit (260 bytes). Extended paths will be used automatically on Windows extraction."

   - Not an error - Windows extraction handles this automatically
   - Informational only

2. **> 1,024 bytes** (Warning): "Path exceeds macOS limit (1,024 bytes). Extraction may fail on macOS."

   - macOS PATH_MAX limit
   - Operation proceeds, but extraction portability affected

3. **> 4,096 bytes** (Warning): "Path exceeds Linux limit (4,096 bytes). Extraction may fail on most filesystems."

   - Linux PATH_MAX limit
   - Operation proceeds, but extraction portability affected

4. **> 32,767 bytes** (Warning): "Path exceeds Windows extended path limit (32,767 bytes). Extraction will fail on Windows."
   - Windows extended path maximum
   - Operation proceeds, but Windows extraction will fail

Warnings are non-fatal - they inform the user about potential portability issues but do not prevent the operation.

Extraction Behavior:

- **Windows**:

  - Paths ≤ 260 bytes: Use regular paths (better compatibility)
  - Paths > 260 bytes: Automatically use extended paths with `\\?\` prefix
  - Seamless support up to ~32,767 bytes
  - No user intervention or registry modifications required
  - Paths > 32,767 bytes: Fail with error

- **Linux**:

  - Validate against PATH_MAX (4,096 bytes)
  - Paths > 4,096 bytes: Fail with error message

- **macOS**:
  - Validate against PATH_MAX (1,024 bytes)
  - Paths > 1,024 bytes: Fail with error message

Benefits:

- Maximum flexibility during package creation
- Windows: Automatic extended path handling (no manual configuration)
- Early warnings about cross-platform portability issues
- Validation only when needed (at extraction time)

**Comparison**: ZIP/TAR require manual extended path enabling; NovusPack handles automatically.

#### 2.1.6 Path Normalization on Storage

When adding files to a package, paths are normalized using the `NormalizePackagePath` function.

For complete normalization details, see [`NormalizePackagePath` Function](#121-normalizepackagepath-function).

The normalization process includes:

1. **Separator normalization**: All path separators are converted to forward slashes (`/`)
2. **Leading slash handling**: Paths without leading `/` are automatically prefixed with `/`
3. **Dot segment canonicalization**: Dot segments (`.` and `..`) are resolved to canonical paths
4. **Unicode normalization**: Path strings are normalized to NFC (Normalization Form Composed)
5. **Path length validation**: Path length is validated against format limits

Filesystem input mapping for file addition is specified by the File Management API.
See [File Addition API - AddFile Operations](api_file_mgmt_addition.md#2-addfile-operations).

### 2.2 Path Rules

- **Stored paths are package-internal**: Paths that refer to package content are stored paths, not OS filesystem paths
  - Examples: `ReadFile(ctx, path)`, `ExtractPath(ctx, path, isWindows, opts)`, `RemoveFile(ctx, path)`
  - Stored path format is specified in [Path Normalization Rules](#21-path-normalization-rules) and [Path Rules](#22-path-rules)
  - **File-add operations accept filesystem inputs**: `AddFile` and `AddFilePattern` accept filesystem-style inputs and derive stored paths
    - See [File Addition API - AddFile Operations](api_file_mgmt_addition.md#2-addfile-operations) for filesystem input mapping rules
- **Leading slash**: All stored paths MUST have a leading slash to ensure full path references
  - `/` by itself refers to the package root, not the OS filesystem root
  - All file and directory paths are stored with a leading `/` (e.g., `/path/to/file.txt`)
  - Input paths without leading slashes are automatically prefixed with `/` during normalization
  - **Display paths**: When displaying paths to end users, the leading `/` MUST be stripped (e.g., stored `/path/to/file.txt` is displayed as `path/to/file.txt`)
- **Dot segments**: `.` and `..` segments are NOT permitted in stored paths
  - Files which are added with dot segments must have those segments converted to canonical paths WITHIN the established package root
  - Canonical paths that would resolve to a path outside of the package root MUST be rejected with an invalid path error
- **Separator normalization**: Separators are normalized to `/`
- **Case sensitivity**: Paths are case sensitive
- **Trailing slashes**: Trailing slashes are meaningful (distinguish files from directories)
- **Empty paths**: Empty paths are not allowed

#### 2.2.1 Case Sensitivity

Storage Policy: Store all paths case-sensitively to preserve exact names.

Examples:

- `file.txt` and `FILE.txt` are stored as distinct paths
- `Config.json` preserves capital 'C'

Cross-Platform Issues:

- **Case-sensitive filesystems** (Linux, some macOS volumes): No issues
- **Case-insensitive filesystems** (Windows, default macOS): Potential conflicts

Extraction Behavior:

On case-insensitive filesystems, if the package contains paths that differ only in case:

1. **Default**: Error with message

   ```text
   Path 'FILE.txt' conflicts with existing path 'file.txt' (case-insensitive filesystem)
   ```

2. **Future Option** (not yet implemented): `AllowCaseConflicts: auto-rename`

   - First file: `file.txt`
   - Second file: `FILE(1).txt` (auto-renamed with suffix)

Portability Warning: Packages containing case-conflicting paths are not portable to case-insensitive filesystems.
Consider this during package design.

Comparison: ZIP/TAR have the same issue - 7-Zip errors, other tools use last-write-wins.
NovusPack provides explicit error with clear message.

### 2.3 Path Display and Extraction

When extracting files or displaying paths to end users:

1. **Platform detection**: The system detects the target platform (Windows vs Unix/Linux)
2. **Strip leading slash**: The leading `/` MUST be removed for display and extraction (it represents package root, not filesystem root)
   - Stored path: `/path/to/file.txt` => Display/Extract as: `path/to/file.txt` (relative to extraction directory)
   - **CRITICAL**: End users should NEVER see the leading `/` in path displays, listings, or extraction operations
3. **Windows conversion**: On Windows systems, forward slashes are converted to backslashes for display and file system operations
   - Stored path: `/path/to/file.txt` => Windows display: `path\to\file.txt`
   - Stored path: `/path/to/file.txt` => Windows extraction: `path\to\file.txt`
4. **Unix/Linux**: On Unix/Linux systems, paths remain with forward slashes (without leading `/`)
   - Stored path: `/path/to/file.txt` => Unix display: `path/to/file.txt`
   - Stored path: `/path/to/file.txt` => Unix extraction: `path/to/file.txt`
5. **API methods**: Methods that return paths for display (e.g., `ListFiles`, `GetPrimaryPath`) MUST strip the leading `/` before returning to callers

**Go API**: `novuspack.ToDisplayPath(storedPath string) string` converts a stored path to display format (strips leading `/`).
See [`ToDisplayPath` Function](#122-todisplaypath-function) for complete path conversion rules.

## 3. Package Writing Operations

See [Package Writing API](api_writing.md) for detailed method signatures and implementation details.

- **SafeWrite**: Atomic write with temp file strategy for data safety
- **FastWrite**: In-place updates for existing packages for performance
- **Write**: General write method with compression handling using SafeWrite or FastWrite methods

## 4. File Management

See [File Management API](api_file_mgmt_index.md) for detailed method signatures and implementation details.

- **Basic File Operations**: Add, remove, and extract files from packages
- **Encryption-Aware File Management**: Add files with specific encryption types
- **Encryption Type System**: Define and validate encryption algorithms
- **ML-KEM Key Management**: Generate and manage post-quantum encryption keys
- **File Pattern Operations**: Add multiple files using patterns
- **File Information and Queries**: Get file information and search capabilities

## 5. Encryption Management

See [Security Validation API - ML-KEM Key Structure and Operations](api_security.md#5-ml-kem-key-structure-and-operations) for detailed method signatures and implementation details.

- **ML-KEM Key Management**: Generate and manage post-quantum encryption keys
- **Key Operations**: Encrypt, decrypt, and manage key lifecycle
- **Security Levels**: Support for multiple security levels (1-5)

## 6. Package Compression Operations

- **CompressPackage**: Compresses package content in memory
- **DecompressPackage**: Decompresses the package in memory
- **CompressPackageFile**: Compresses package content and writes to specified path
- **DecompressPackageFile**: Decompresses the package and writes to specified path
- **GetPackageCompressionInfo**: Returns package compression information
- **IsPackageCompressed**: Checks if the package is compressed
- **GetPackageCompressionType**: Returns the package compression type
- **SetPackageCompressionType**: Sets the package compression type (without compressing)
- **CanCompressPackage**: Checks if package can be compressed (not signed)

See [Package Compression API - Compression Types](api_package_compression.md#12-compression-types) for compression type constants and [Compression Information Structure](api_package_compression.md#71-compression-information-structure-reference) for the PackageCompressionInfo structure.

See [Package Compression API - Compression Scope](api_package_compression.md#11-compression-scope) and [Compression and Signing Relationship](api_package_compression.md#10-compression-and-signing-relationship) for detailed behavior specifications.

See [Package Compression API](api_package_compression.md) for detailed method signatures and implementation details.

## 7. Digital Signatures and Security

**Cross-Reference**: Signature management, signing, and signature validation are deferred to v2.
See [Digital Signature API](api_signatures.md) for future work.

### 7.1 Core Integration Points

The core package interface integrates with signature presence through:

- **Immutability Enforcement**: All write operations check `SignatureOffset > 0` before proceeding
- **Write Protection**: Signed packages are protected from write operations by default
- **Context Integration**: All signature operations accept `ctx context.Context` as first parameter
- **Error Handling**: Signature operations use the structured error system defined in [Structured Error System](#10-structured-error-system)

### 7.2 Write Protection and Immutability Enforcement

- **Signed File Detection**: All write operations must check if `SignatureOffset > 0` before proceeding
- **Write Protection**: Signed packages are protected from write operations by default
- **Writing Signed Package Content**: Writing signed package content is allowed only after reconfiguring to a new target path, which MUST clear signature information from the in-memory Package
- **Allowed Operations**: If signed, only read operations are allowed unless signatures are cleared by writing to a new path
- **Prohibited Operations**: Header modifications and content changes are prohibited on signed packages
- **Signature Removal**: Clearing signatures is allowed only as part of an explicit unsigned copy workflow
- **Detailed Behavior**: See [Package Writing Operations - Signed File Write Operations](api_writing.md#4-signed-file-write-operations) for complete implementation details
- **Purpose**: This prevents accidental signature invalidation and maintains package integrity

## 8. Package Metadata Management

See [Package Metadata API](metadata.md) for detailed method signatures and implementation details.

### 8.1 General Metadata Operations

- **SetMetadata**: Sets package metadata
- **GetMetadata**: Retrieves package metadata
- **UpdateMetadata**: Updates package metadata
- **ValidateMetadata**: Validates metadata structure and content
- **HasMetadata**: Checks if package has metadata
- **AddMetadataFile**: Adds metadata as special file
- **GetMetadataFile**: Retrieves metadata from special file
- **UpdateMetadataFile**: Updates metadata file
- **RemoveMetadataFile**: Removes metadata special file

### 8.2 AppID/VendorID Management

- **SetAppID**: Sets the application identifier
- **GetAppID**: Gets the current application identifier
- **ClearAppID**: Clears the application identifier
- **HasAppID**: Checks if application identifier is set
- **SetVendorID**: Sets the vendor/platform identifier
- **GetVendorID**: Gets the current vendor identifier
- **ClearVendorID**: Clears the vendor identifier
- **HasVendorID**: Checks if vendor identifier is set

### 8.3 Package Information Structures

- **PackageInfo**: Comprehensive package information and metadata (see [PackageInfo Structure](api_metadata.md#71-packageinfo-structure))
- **SignatureInfo**: Detailed signature information
- **SecurityStatus**: Current security status and validation results

## 9. File Validation Requirements

Requirements for files that are added to packages.

- **Path validation**: Paths must conform to [Package Path Semantics](#2-package-path-semantics)
  - Path must not be empty (see [Path Rules](#22-path-rules))
  - Path must be validated according to [Path Validation `ValidatePackagePath`](#123-validatepackagepath-function) rules
- Data must not be nil (empty files with len = 0 are allowed)
- Invalid files will be rejected with appropriate error messages

For complete file validation specifications, see [File Validation Requirements](file_validation.md).

## 10. Structured Error System

The NovusPack API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities while maintaining compatibility with Go's standard error handling patterns.

All errors are returned as `PackageError` instances with appropriate error types, messages, and contextual information.

### 10.1 Benefits of Structured Errors

- **Better Error Categorization**: Errors grouped by type for easier handling
- **Rich Error Context**: Additional context fields for debugging
- **Type Safety**: Structured errors can be inspected with type assertions
- **Consistent API**: All errors follow the same structured pattern
- **Better Logging**: Structured errors provide more information for logs
- **Testing**: Easier to test error conditions with typed errors

### 10.2 ErrorType Types and Categories

```go
// ErrorType categorizes errors for programmatic handling.
// ErrorType provides structured error categorization for better error handling and diagnostics.
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
    ErrTypeState
    ErrTypeFormat
    ErrTypeVersion
    ErrTypeResource
    ErrTypeMetadata
)
```

### 10.3 ErrorType Categories

- **ErrTypeValidation**: Input validation errors, invalid parameters, invalid file paths, invalid patterns
- **ErrTypeIO**: File I/O errors, permission errors, disk space issues
- **ErrTypeSecurity**: Security-related errors, access denied, authentication failures
- **ErrTypeCorruption**: Data corruption, checksum failures, integrity violations
- **ErrTypeUnsupported**: Unsupported features or operations (not version-related)
- **ErrTypeContext**: Context cancellation, timeout, or context-related errors
- **ErrTypeEncryption**: Encryption/decryption failures, key management errors
- **ErrTypeCompression**: Compression/decompression failures, algorithm errors
- **ErrTypeSignature**: Digital signature validation, signing failures
- **ErrTypeState**: Package state errors (package not open, package read-only, package already open, package already closed)
- **ErrTypeFormat**: Package format/parsing errors (invalid package format, invalid header format, malformed package structure)
- **ErrTypeVersion**: Version incompatibility errors (unsupported package version, version mismatch)
- **ErrTypeResource**: Resource exhaustion errors (memory insufficient, resource limits exceeded)
- **ErrTypeMetadata**: Metadata-specific errors (metadata validation failures, metadata corruption, invalid metadata format)

### 10.4 PackageError Structure

```go
// PackageError represents a structured error in package operations.
type PackageError struct {
    Type    ErrorType              // Error category
    Message string                 // Human-readable error message
    Cause   error                  // Underlying error (wrapped)
    Context map[string]interface{} // Additional error context
}
```

`PackageError` is a structured error type that provides categorized error information with context support and compatibility with Go's standard error handling patterns.

#### 10.4.1 PackageError.Error Method

```go
// Error implements the error interface
func (e *PackageError) Error() string
```

`Error` implements the `error` interface, returning a formatted error message.

- If a cause error exists, returns a formatted string combining the message and cause: `"{Message}: {Cause}"`
- If no cause exists, returns only the message

#### 10.4.2 PackageError.Unwrap Method

```go
// Unwrap returns the underlying error for error unwrapping
func (e *PackageError) Unwrap() error
```

`Unwrap` returns the underlying error for error unwrapping, enabling compatibility with Go's `errors.Unwrap` function and error chain traversal.

- Returns the `Cause` field directly
- Returns `nil` if no cause error is set

#### 10.4.3 PackageError.Is Method

```go
// Is implements error matching for error comparison
func (e *PackageError) Is(target error) bool
```

`Is` implements error matching for error comparison, enabling compatibility with Go's `errors.Is` function.

- If a cause error exists, delegates to `errors.Is(e.Cause, target)` to check if the cause matches the target error
- If no cause exists, returns `false`
- This allows `PackageError` to participate in Go's standard error matching patterns

### 10.5 Error Helper Functions

Helper functions for creating and managing structured errors.

#### 10.5.1 NewPackageError Function

```go
// NewPackageError creates a structured error with type-safe context
// All errors must include typed context for type safety
func NewPackageError[T any](errType ErrorType, message string, cause error, context T) *PackageError
```

`NewPackageError` creates a structured error with type-safe context, combining error creation with type-safe context in a single operation.

- Creates a new `PackageError` struct with the provided error type, message, and cause
- Stores the typed context in the error's Context map with the special key `"_typed_context"` to preserve type information
- Returns a pointer to the newly created `PackageError` with typed context
- All errors must include typed context for compile-time type safety

#### 10.5.2 WrapErrorWithContext Function

```go
// WrapErrorWithContext wraps an error with type-safe context
func WrapErrorWithContext[T any](err error, errType ErrorType, message string, context T) *PackageError
```

`WrapErrorWithContext` wraps an error with type-safe context, providing a convenient way to wrap errors with typed contextual information.

#### 10.5.3 AsPackageError Function

```go
// AsPackageError checks if an error is a PackageError and returns it if found
func AsPackageError(err error) (*PackageError, bool)
```

`AsPackageError` checks if an error is a `PackageError` and returns the `PackageError` pointer if found, enabling type-safe error inspection and access.

- Uses Go's `errors.As` function to attempt type assertion and extraction
- If the error is a `PackageError` (or wraps one), returns the `PackageError` pointer and `true`
- If the error is not a `PackageError`, returns `nil` and `false`
- This enables safe error type checking and access: `if pkgErr, ok := AsPackageError(err); ok { ... }`
- Follows the same pattern as Go's standard `errors.As` function

#### 10.5.4 GetErrorContext Function

```go
// GetErrorContext retrieves type-safe context from errors
func GetErrorContext[T any](err error, key string) (T, bool)
```

`GetErrorContext` retrieves type-safe context from errors using generics, enabling type-safe access to error context values.

#### 10.5.5 AddErrorContext Function

```go
// AddErrorContext adds type-safe context to errors
func AddErrorContext[T any](err error, key string, value T) error
```

`AddErrorContext` adds type-safe context to errors using generics, providing compile-time type safety for error context values.

#### 10.5.6 MapError Function

```go
// MapError transforms an error with a generic mapper function
func MapError[T any, U any](err error, mapper func(T) U) error
```

`MapError` transforms an error with a generic mapper function, enabling error transformation patterns with type safety.

- Extracts typed context from the error using `GetErrorContext[T]`
- Applies the mapper function to transform the context from type `T` to type `U`
- Creates a new error with the transformed context
- Returns the original error unchanged if no typed context is found
- This enables error transformation patterns: `MapError(err, func(old OldContext) NewContext { ... })`

#### 10.5.7 Example - Creating Errors with Context

Define error context types for type safety, then use them with `NewPackageError` or `WrapErrorWithContext`:

See [FileErrorContext Structure](api_core.md#104-packageerror-structure) for the complete structure definition.

Example usage:

```go
// Example: Define error context types for type safety
// Note: FileErrorContext is defined in api_core.md
type ExampleFileErrorContext struct {
    Path      string
    Operation string
    Size      int64
}

// Example: Create a new validation error with typed context
err := NewPackageError(ErrTypeValidation, "file not found", nil, ExampleFileErrorContext{
    Path:      "/path/to/file",
    Operation: "AddFile",
    Size:      0,
})

// Example: Create error with context
err := NewPackageError(ErrTypeValidation, "file validation failed", nil, FileErrorContext{
    Path:      "/path/to/file",
    Operation: "Validate",
    Size:      0,
})

// Example: Wrap error with typed context
err := WrapErrorWithContext(io.ErrUnexpectedEOF, ErrTypeIO, "failed to read file", FileErrorContext{
    Path:      "/path/to/file",
    Operation: "ReadFile",
    Size:      0,
})

// Example: Wrap existing error with typed context
err := WrapErrorWithContext(io.ErrUnexpectedEOF, ErrTypeIO, "unexpected end of file", FileErrorContext{
    Path:      "/path/to/file",
    Operation: "ReadFile",
    Size:      0,
})
```

#### 10.5.8 Example - Error Inspection and Handling

```go
// Example: Check if error is a PackageError and extract it
if pkgErr, ok := AsPackageError(err); ok {
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
if pkgErr, ok := AsPackageError(err); ok && pkgErr.Type == ErrTypeValidation {
    // Handle validation errors specifically
}
```

#### 10.5.9 Example - Error Propagation

```go
// Example: file operation that returns structured errors.
// This is an example only and does not define the canonical signature.
//
// For AddFile semantics and signature, see:
// - docs/tech_specs/api_file_mgmt_addition.md#21-addfile
//
// Example implementation pattern (not canonical):
func exampleReadFileOperation(ctx context.Context, pkg PackageReader, filePath string) ([]byte, error) {
    // Implementation would wrap errors with structured context
    return nil, nil
}
```

Error propagation in the NovusPack API follows a consistent pattern that wraps errors with structured information and contextual details.

When an error occurs during an operation:

1. **Wrap the error** using `WrapErrorWithContext[T]` with an appropriate error type, descriptive message, and typed context
2. **Use typed context** with `NewPackageError[T]` or `WrapErrorWithContext[T]` to provide type-safe contextual information about the operation
3. **Return the structured error** to allow callers to inspect and handle errors appropriately

**Generic Error Context**: For type-safe error context, use `WrapErrorWithContext[T]` or `NewPackageError[T]` with typed context structures.
This provides compile-time type safety and better error inspection capabilities.

#### 10.5.10 Error Example Pattern

For validation errors:

- Create the error with `NewPackageError` and typed context
- Include relevant information in the context structure
- Example: Define a context type and use it: `type FileErrorContext struct { Path string; Size int64 }` then `NewPackageError(ErrTypeValidation, "file validation failed", nil, FileErrorContext{Path: path, Size: len(data)})`

For I/O errors:

- Create the error with `NewPackageError` and typed context
- Include relevant information in the context structure
- Example: `NewPackageError(ErrTypeIO, "failed to write file data", err, FileErrorContext{Path: path})`

This pattern ensures that errors carry rich contextual information while maintaining compatibility with Go's standard error handling patterns.

#### 10.5.11 Error Wrapping Patterns

When wrapping external errors or standard library errors, use `WrapErrorWithContext[T]` to convert them to structured errors with appropriate error types and typed context.

**Pattern**: `WrapErrorWithContext(externalErr, ErrTypeX, "descriptive message", ErrorContext{...})`

- Wrap standard library errors (e.g., `io.EOF`, `os.ErrNotExist`) with appropriate error types and typed context
- Wrap third-party library errors with appropriate error types and typed context
- Always provide descriptive messages that explain the context of the error
- Always include typed context structures for type safety

##### 10.5.11.1 Example - Error Wrapping Patterns

```go
// Example: Wrap a standard library error with typed context
// Note: FileErrorContext is defined in api_core.md
type ExampleFileErrorContext struct {
    Path string
}

if err := os.Open(filePath); err != nil {
    return WrapErrorWithContext(err, ErrTypeIO, "failed to open file", ExampleFileErrorContext{
        Path: filePath,
    })
}
```

This ensures all errors in the NovusPack API are structured and provide rich contextual information for debugging and error handling.

#### 10.5.12 Error Logging and Debugging

Error logging should leverage the structured error system to provide comprehensive debugging information.

When logging errors, check if the error is a `PackageError` using `AsPackageError`:

1. **For PackageError instances**:

   - Log the operation name, error type (as integer), error message, and full context map
   - If a cause error exists, log it separately with a "Caused by:" prefix
   - Format: `"Error in {operation}: Type={type}, Message={message}, Context={context}"`
   - Format (with cause): `"Caused by: {cause}"`

2. **For non-PackageError instances**:
   - Log the operation name and error value
   - Format: `"Error in {operation}: {error}"`

This pattern ensures that structured errors provide rich debugging information while maintaining compatibility with standard Go errors.

##### 10.5.12.1 Example - Error Logging Pattern

```go
// Example:
if pkgErr, ok := AsPackageError(err); ok {
    log.Printf("Error in %s: Type=%d, Message=%s, Context=%+v",
        operation, pkgErr.Type, pkgErr.Message, pkgErr.Context)
    if pkgErr.Cause != nil {
        log.Printf("Caused by: %v", pkgErr.Cause)
    }
} else {
    log.Printf("Error in %s: %v", operation, err)
}
```

##### 10.5.12.2 Example - Type-Safe Error Context Retrieval

```go
// Example: Retrieve typed context from error
// Note: FileErrorContext is defined in api_core.md
type ExampleFileErrorContext struct {
    Path      string
    Operation string
    Size      int64
}

if ctx, ok := GetErrorContext[ExampleFileErrorContext](err, "_typed_context"); ok {
    log.Printf("Error context: Path=%s, Operation=%s, Size=%d",
        ctx.Path, ctx.Operation, ctx.Size)
} else {
    // Fallback to untyped context
    if path, exists := pkgErr.Context["path"]; exists {
        log.Printf("Path: %v", path)
    }
}
```

##### 10.5.12.3 Example - Error Transformation

This example demonstrates how to transform errors between different context types.

```go
// Example: Transform error context using MapError
type OldContext struct {
    FilePath string
    Size     int64
}

//Example struct
type NewContext struct {
    Path     string
    FileSize int64
}

transformedErr := MapError(err, func(old OldContext) NewContext {
    return NewContext{
        Path:     old.FilePath,
        FileSize: old.Size,
    }
})

// Use transformed error
if newCtx, ok := GetErrorContext[NewContext](transformedErr, "_typed_context"); ok {
    log.Printf("Transformed context: Path=%s, FileSize=%d", newCtx.Path, newCtx.FileSize)
}
```

## 11. Generic Types

For comprehensive generic type definitions, usage examples, and best practices, see [Generic Types and Patterns](api_generics.md).

The NovusPack API provides generic types for improved type safety and code reuse across different data types.
The following generic types are used throughout the NovusPack API:

- **[Option Type](api_generics.md#11-option-type)**: Type-safe optional values (`Option[T]`)
- **[Result Type](api_generics.md#12-result-type)**: Type-safe error handling (`Result[T]`)
- **[Core Generic Types](api_generics.md#1-core-generic-types)**: Generic strategy pattern (`Strategy[T, U]`)
- **[Core Generic Types](api_generics.md#1-core-generic-types)**: Generic validation (`Validator[T]`)
- **[Core Generic Types](api_generics.md#1-core-generic-types)**: Type-safe configuration (`Config[T]`, `ConfigBuilder[T]`)

All generic type definitions, interfaces, and usage patterns are documented in the dedicated [Generic Types and Patterns](api_generics.md) specification.

## 12. Utility Functions

Utility functions for working with package paths.

### 12.1 NormalizePackagePath Function

```go
// NormalizePackagePath normalizes a package-internal path for consistent
// comparison and storage. Applies separator normalization, dot-segment
// canonicalization, leading slash, NFC, and path-length checks.
// Returns the normalized path with leading "/" or an error if the path is
// invalid or would escape the package root.
func NormalizePackagePath(path string) (string, error)
```

`NormalizePackagePath` normalizes a package-internal path according to NovusPack path storage rules.

This function applies the normalization rules defined in [Separator Normalization](#211-separator-normalization) in the following order:

1. **Separator normalization**: Converts all path separators to forward slashes (`/`)
2. **Leading slash handling**: Ensures the path has a leading `/` (adds it if missing)
3. **Dot segment canonicalization**: Resolves `.` and `..` segments to canonical paths
4. **Unicode normalization**: Normalizes path strings to NFC (Normalization Form Composed)
5. **Path length validation**: Validates that the path length is within format limits (up to 65,535 bytes)

#### 12.1.1 NormalizePackagePath Error Handling

- Returns an error if dot segments would escape the package root
- Returns an error if the path length exceeds the format limit (65,535 bytes)
- Returns an error if the path is empty or malformed

#### 12.1.2 NormalizePackagePath Return Value

The normalized path is suitable for storage in the package format and for use in package operations such as `ReadFile`, `RemoveFile`, and `ExtractPath`.

For complete path normalization rules and details, see [Path Normalization Rules](#21-path-normalization-rules).
For complete path storage format rules, see [Path Rules](#22-path-rules).

### 12.2 ToDisplayPath Function

```go
// ToDisplayPath converts a stored package path to display format by stripping the leading slash.
func ToDisplayPath(storedPath string) string
```

`ToDisplayPath` converts a stored package path to display format by stripping the leading slash.

This function implements the display path conversion rule defined in [Path Rules](#22-path-rules).

#### 12.2.1 ToDisplayPath Behavior

- **Input**: A stored path with leading `/` (e.g., `/path/to/file.txt`)
- **Output**: Display path without leading `/` (e.g., `path/to/file.txt`)
- **Root path**: The root path `/` is converted to an empty string `""`
- **No validation**: This function does not validate the input path; it only performs the leading slash removal

#### 12.2.2 ToDisplayPath Usage

This function is used when displaying paths to end users in:

- File listings (e.g., `ListFiles()` returns display paths)
- Path metadata displays
- Error messages and logging
- User-facing API responses

**Important**: The leading `/` in stored paths represents the _package root_, not the OS filesystem root.
When displaying paths to users, this leading slash MUST be stripped to avoid confusion.

For platform-specific filesystem display (e.g., converting to Windows backslashes), additional conversion may be needed after calling `ToDisplayPath`.

For complete path conversion rules, see [Path Display and Extraction](#23-path-display-and-extraction).

### 12.3 ValidatePackagePath Function

```go
// ValidatePackagePath validates a package path according to package path semantics.
// Validates path format, rejects empty or whitespace-only paths, normalizes separators,
// canonicalizes dot segments, and ensures the path does not escape the package root.
func ValidatePackagePath(path string) error
```

Invalid paths include:

- Empty path
- Malformed path
- Paths containing dot segments (`.` or `..`) that, when converted to canonical paths, would escape the package root (for example, `../../../etc/passwd` resolves outside the package root and must be rejected)

All reader and writer methods that accept path parameters must validate paths according to these rules and return `*PackageError` with `ErrTypeValidation` for invalid paths.

### 12.4 ValidatePathLength Function

```go
// ValidatePathLength validates platform path length portability constraints.
// Returns:
// - warnings: Non-fatal portability warnings.
// - error: ErrTypeValidation when the path exceeds the hard limit.
func ValidatePathLength(path string) ([]string, error)
```

`ValidatePathLength` validates path length constraints and returns portability warnings.

This function is intended for use in user-facing validation and portability checks.
It does not change the path.

Behavior:

- **Warnings (`[]string` return value)**: The function returns non-fatal portability warnings in the `[]string` return value when paths exceed platform-specific limits:
  - Paths > 260 bytes: Info warning about Windows default limit (extended paths will be used automatically)
  - Paths > 1,024 bytes: Warning about macOS PATH_MAX limit (extraction may fail on macOS)
  - Paths > 4,096 bytes: Warning about Linux PATH_MAX limit (extraction may fail on most filesystems)
  - Paths > 32,767 bytes: Warning about Windows extended path limit (extraction will fail on Windows)
  - Multiple warnings may be returned if a path exceeds multiple thresholds
  - Warnings are informational only and do not prevent the operation from proceeding
- **Error return**: The function returns an error (`ErrTypeValidation`) when the path exceeds the hard maximum limit (Windows extended-length 32,767 bytes), indicating that extraction will fail on Windows systems
- **No warnings or errors**: If the path is within all platform limits, the function returns an empty `[]string` and `nil` error

For complete details on path length limits and warning thresholds, see [Path Length Limits](#215-path-length-limits).

**Go API**: `novuspack.NormalizePackagePath`, `novuspack.ToDisplayPath`, `novuspack.ValidatePackagePath`, and `novuspack.ValidatePathLength` export path normalization, stored-to-display conversion, path validation, and path length validation.

See [`NormalizePackagePath` Function](#121-normalizepackagepath-function), [`ToDisplayPath` Function](#122-todisplaypath-function), [`ValidatePackagePath` Function](#123-validatepackagepath-function), and [`ValidatePathLength` Function](#124-validatepathlength-function) for complete function documentation.
