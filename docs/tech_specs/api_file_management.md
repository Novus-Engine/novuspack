# NovusPack Technical Specifications - File Management API

## Table of Contents

- [Table of Contents](#table-of-contents)
- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Core Data Structures](#1-core-data-structures)
  - [1.1 FileEntry Structure](#11-fileentry-structure)
  - [1.2.1 Tag Management](#121-tag-management)
  - [1.2.2 Data Management](#122-data-management)
  - [1.2.3 Path and Directory Management](#123-path-and-directory-management)
  - [1.2.4 Serialization](#124-serialization)
  - [1.2.5 FileEntry Directory Association Methods](#125-fileentry-directory-association-methods)
- [2. Add File Operations](#2-add-file-operations)
  - [2.1 Add File](#21-add-file)
  - [2.2 Add File Pattern](#22-add-file-pattern)
  - [2.3 FileSource Interface](#23-filesource-interface)
  - [2.4 FileSource Implementations](#24-filesource-implementations)
  - [2.5 AddFileOptions Configuration](#25-addfileoptions-configuration)
  - [2.6 Usage Notes](#26-usage-notes)
- [3. File Addition Implementation Flow](#3-file-addition-implementation-flow)
  - [3.1 Processing Order Requirements](#31-processing-order-requirements)
- [4. Remove File Operations](#4-remove-file-operations)
  - [4.1 Remove File](#41-remove-file)
  - [4.2 Remove File Pattern](#42-remove-file-pattern)
- [5. Extract File Operations](#5-extract-file-operations)
  - [5.1 Extract File](#51-extract-file)
- [6. Update File Operations](#6-update-file-operations)
  - [6.1 Update File](#61-update-file)
  - [6.2 Update File Pattern](#62-update-file-pattern)
  - [6.3 Update File Metadata](#63-update-file-metadata)
  - [6.4 Add File Path](#64-add-file-path)
  - [6.5 Remove File Path](#65-remove-file-path)
  - [6.6 Add File Hash](#66-add-file-hash)
- [7. File Compression Operations](#7-file-compression-operations)
  - [7.1 Purpose](#71-purpose)
  - [7.1.1 CompressFile Parameters](#711-compressfile-parameters)
- [8. File Encryption Operations](#8-file-encryption-operations)
- [9. Deduplication Operations](#9-deduplication-operations)
  - [9.1 File Deduplication](#91-file-deduplication)
- [10. File Information and Queries](#10-file-information-and-queries)
  - [10.1 File Existence and Properties](#101-file-existence-and-properties)
  - [10.2 File Lookup by Metadata](#102-file-lookup-by-metadata)
- [11. FileEntry Methods](#11-fileentry-methods)
  - [11.1 File Entry Properties](#111-file-entry-properties)
  - [11.2 File Entry Encryption](#112-file-entry-encryption)
  - [11.3 File Entry Data Management](#113-file-entry-data-management)
- [12. Error Handling](#12-error-handling)
  - [12.1 Structured Error System](#121-structured-error-system)
  - [12.2 Common Error Types](#122-common-error-types)
  - [12.3 Structured Error Examples](#123-structured-error-examples)
  - [12.4 Error Handling Best Practices](#124-error-handling-best-practices)
  - [12.5 Generic FileEntry Operations](#125-generic-fileentry-operations)
- [13. Best Practices](#13-best-practices)
  - [13.1 File Path Management](#131-file-path-management)
  - [13.2 Encryption Management](#132-encryption-management)
  - [13.3 Performance Considerations](#133-performance-considerations)

## 0. Overview

This document defines the file management API for the NovusPack system, covering basic file operations, encryption-aware file management, and cryptographic key management.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Basic Operations API](api_basic_operations.md) - Package creation, opening, closing, and lifecycle management
- [Package Metadata API](api_metadata.md) - Comment management, AppID/VendorID, and metadata operations
- [Security Validation API](api_security.md) - Package validation and security status structures
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Package Compression API](api_package_compression.md) - Generic strategy patterns and compression concurrency

## 1. Core Data Structures

### 1.1 FileEntry Structure

```go
// FileType represents the file type identifier from the file type system
type FileType uint16

// FileEntry represents a file entry in the package with complete metadata
type FileEntry struct {
    // Static fields (64 bytes total)
    FileID             uint64    // Unique file identifier (8 bytes)
    OriginalSize       uint64    // Original file size before processing (8 bytes)
    StoredSize         uint64    // Final file size after compression/encryption (8 bytes)
    RawChecksum        uint32    // CRC32 of raw file content (4 bytes)
    StoredChecksum     uint32    // CRC32 of processed file content (4 bytes)
    FileVersion        uint32    // File data version (4 bytes)
    MetadataVersion    uint32    // File metadata version (4 bytes)
    PathCount          uint16    // Number of paths (2 bytes)
    Type               FileType  // File type identifier (2 bytes)
    CompressionType    uint8     // Compression algorithm identifier (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    CompressionLevel   uint8     // Compression level (1 byte)
    EncryptionType     uint8     // Encryption algorithm identifier (1 byte)
    HashCount          uint8     // Number of hash entries (1 byte)
    HashDataOffset     uint32    // Offset to hash data (4 bytes)
    HashDataLen        uint16    // Length of hash data (2 bytes)
    OptionalDataLen    uint16    // Length of optional data (2 bytes)
    OptionalDataOffset uint32    // Offset to optional data (4 bytes)
    Reserved           uint32    // Reserved for future use (4 bytes)

    // Variable-length data (populated on demand)
    Paths              []PathEntry    // File paths with metadata
    Hashes             []HashEntry    // Content hashes
    OptionalData       OptionalData   // Structured optional data

    // Convenience properties (computed from OptionalData)
    Tags               *[]Tag         // Direct access to OptionalData.Tags

    // Data management (runtime only, not stored in file)
    Data               []byte          // File content in memory (only for small files being processed)
    SourceFile         *os.File        // Source file handle for streaming
    SourceOffset       int64           // Offset in source file
    SourceSize         int64           // Size of data to read from source
    TempFilePath       string          // Path to temp file for large files being processed
    IsDataLoaded       bool            // Whether data is currently loaded in memory
    IsTempFile         bool            // Whether file is stored as temp file during processing
    ProcessingState    ProcessingState // Current processing state of the file

    // Directory association (runtime only, not stored in file)
    ParentDirectory    *DirectoryEntry // Pointer to parent directory (nil for root-relative)
    InheritedTags      *[]Tag          // Pointer to cached inherited tags from directory hierarchy
}

// PathEntry represents a file path with metadata
type PathEntry struct {
    Path       string    // File path (UTF-8)
    Mode       uint32    // File permissions and type (Unix-style)
    UserID     uint32    // User ID (Unix-style)
    GroupID    uint32    // Group ID (Unix-style)
    ModTime    time.Time // Last modification time
    CreateTime time.Time // File creation time
    AccessTime time.Time // Last access time

    // Symbolic link support
    IsSymlink  bool      // Whether this path is a symbolic link
    LinkTarget string    // Target path for symbolic links (empty if not a symlink)
}

// HashEntry represents a hash with type and purpose
type HashEntry struct {
    Type    HashType    // Hash algorithm type
    Purpose HashPurpose // Hash purpose
    Data    []byte      // Hash data
}

// HashType represents hash algorithm types
type HashType uint8
const (
    HashTypeSHA256   HashType = 0x00  // SHA-256 (32 bytes) - Standard cryptographic hash
    HashTypeSHA512   HashType = 0x01  // SHA-512 (64 bytes) - Stronger cryptographic hash
    HashTypeBLAKE3   HashType = 0x02  // BLAKE3 (32 bytes) - Fast cryptographic hash
    HashTypeXXH3     HashType = 0x03  // XXH3 (8 bytes) - Ultra-fast non-cryptographic hash
    HashTypeBLAKE2b  HashType = 0x04  // BLAKE2b (64 bytes) - Cryptographic hash with configurable output
    HashTypeBLAKE2s  HashType = 0x05  // BLAKE2s (32 bytes) - Cryptographic hash optimized for 32-bit systems
    HashTypeSHA3_256 HashType = 0x06  // SHA-3-256 (32 bytes) - SHA-3 family hash
    HashTypeSHA3_512 HashType = 0x07  // SHA-3-512 (64 bytes) - SHA-3 family hash
    HashTypeCRC32    HashType = 0x08  // CRC32 (4 bytes) - Fast checksum for error detection
    HashTypeCRC64    HashType = 0x09  // CRC64 (8 bytes) - Stronger checksum for error detection
)

// HashPurpose represents hash purposes
type HashPurpose uint8
const (
    HashPurposeContentVerification HashPurpose = 0x00  // Content verification - Verify file content integrity
    HashPurposeDeduplication       HashPurpose = 0x01  // Deduplication - Identify duplicate content
    HashPurposeIntegrityCheck      HashPurpose = 0x02  // Integrity check - General integrity verification
    HashPurposeFastLookup          HashPurpose = 0x03  // Fast lookup - Quick content identification
    HashPurposeErrorDetection      HashPurpose = 0x04  // Error detection - Detect data corruption
)

// Tag represents a key-value pair with typed value
type Tag struct {
    Key       string    // Tag key (UTF-8 string)
    ValueType TagValueType // Value type
    Value     string    // Value as UTF-8 string (type-specific encoding)
}

// TagValueType represents the type of a tag value
type TagValueType uint8
const (
    // Basic Types
    TagValueTypeString      TagValueType = 0x00  // String value
    TagValueTypeInteger     TagValueType = 0x01  // 64-bit signed integer
    TagValueTypeFloat       TagValueType = 0x02  // 64-bit floating point number
    TagValueTypeBoolean     TagValueType = 0x03  // Boolean value

    // Structured Data
    TagValueTypeJSON        TagValueType = 0x04  // JSON-encoded object or array
    TagValueTypeYAML        TagValueType = 0x05  // YAML-encoded data
    TagValueTypeStringList  TagValueType = 0x06  // Comma-separated list of strings

    // Identifiers
    TagValueTypeUUID        TagValueType = 0x07  // UUID string
    TagValueTypeHash        TagValueType = 0x08  // Hash/checksum string
    TagValueTypeVersion     TagValueType = 0x09  // Semantic version string

    // Time
    TagValueTypeTimestamp   TagValueType = 0x0A  // ISO8601 timestamp

    // Network/Communication
    TagValueTypeURL         TagValueType = 0x0B  // URL string
    TagValueTypeEmail       TagValueType = 0x0C  // Email address

    // File System
    TagValueTypePath        TagValueType = 0x0D  // File system path
    TagValueTypeMimeType    TagValueType = 0x0E  // MIME type string

    // Localization
    TagValueTypeLanguage    TagValueType = 0x0F  // Language code (ISO 639-1)

    // NovusPack Special Files
    TagValueTypeNovusPackMetadata TagValueType = 0x10  // NovusPack special metadata file reference
)

// ProcessingState defines the current state of file processing
type ProcessingState uint8

const (
    ProcessingStateIdle ProcessingState = iota // File not being processed
    ProcessingStateLoading                     // Loading data from source
    ProcessingStateProcessing                  // Processing data (compression, encryption, etc.)
    ProcessingStateWriting                     // Writing to package
    ProcessingStateComplete                    // Processing complete, ready for cleanup
    ProcessingStateError                       // Error occurred during processing
)

// OptionalData represents structured optional data for a file entry
type OptionalData struct {
    Tags                []Tag               // Per-file tags data (DataType 0x00)
    PathEncoding        *uint8              // Path encoding type (DataType 0x01)
    PathFlags           *uint8              // Path handling flags (DataType 0x02)
    CompressionDictID   *uint32             // Dictionary ID for solid compression (DataType 0x03)
    SolidGroupID        *uint32             // Solid compression group ID (DataType 0x04)
    FileSystemFlags     *uint16             // File system specific flags (DataType 0x05)
    WindowsAttributes   *uint32             // Windows file attributes (DataType 0x06)
    ExtendedAttributes  map[string]string   // Unix extended attributes (DataType 0x07)
    ACLData             []byte              // Access Control List data (DataType 0x08)
    CustomData          map[uint8][]byte    // Custom data for reserved types (0x09-0xFF)
}

// OptionalDataType represents the type of optional data
type OptionalDataType uint8
const (
    OptionalDataTypeTags                OptionalDataType = 0x00  // TagsData
    OptionalDataTypePathEncoding        OptionalDataType = 0x01  // PathEncoding
    OptionalDataTypePathFlags           OptionalDataType = 0x02  // PathFlags
    OptionalDataTypeCompressionDictID   OptionalDataType = 0x03  // CompressionDictionaryID
    OptionalDataTypeSolidGroupID        OptionalDataType = 0x04  // SolidGroupID
    OptionalDataTypeFileSystemFlags     OptionalDataType = 0x05  // FileSystemFlags
    OptionalDataTypeWindowsAttributes   OptionalDataType = 0x06  // WindowsAttributes
    OptionalDataTypeExtendedAttributes  OptionalDataType = 0x07  // ExtendedAttributes
    OptionalDataTypeACLData             OptionalDataType = 0x08  // ACLData
)

// Helper functions for creating pointers to optional data
func uint8Ptr(v uint8) *uint8
func uint16Ptr(v uint16) *uint16
func uint32Ptr(v uint32) *uint32

// Note: These methods assume the following imports:
// import (
//     "context"
//     "encoding/json"
//     "fmt"
//     "io"
//     "os"
//     "strconv"
//     "strings"
//     "time"
//     "gopkg.in/yaml.v3"
//     "github.com/google/uuid"
// )
// DirectoryEntry, DirectoryInheritance, DirectoryMetadata, DirectoryFileSystem, ACLEntry from api_metadata

// FileEntry Methods
```

### 1.2.1 Tag Management

```go
func (fe *FileEntry) GetTags(ctx context.Context) []Tag
func (fe *FileEntry) SetTags(ctx context.Context, tags []Tag)
func (fe *FileEntry) SetTag(ctx context.Context, key string, valueType TagValueType, value string)
func (fe *FileEntry) SetStringTag(ctx context.Context, key, value string)
func (fe *FileEntry) SetIntegerTag(ctx context.Context, key string, value int64)
func (fe *FileEntry) SetBooleanTag(ctx context.Context, key string, value bool)
func (fe *FileEntry) SetJSONTag(ctx context.Context, key string, value interface{}) error
func (fe *FileEntry) SetStringListTag(ctx context.Context, key string, values []string)
func (fe *FileEntry) GetStringListTag(ctx context.Context, key string) ([]string, bool)
func (fe *FileEntry) SetYAMLTag(ctx context.Context, key string, value interface{}) error
func (fe *FileEntry) GetYAMLTag(ctx context.Context, key string, dest interface{}) error
func (fe *FileEntry) SetFloatTag(ctx context.Context, key string, value float64)
func (fe *FileEntry) GetFloatTag(ctx context.Context, key string) (float64, bool)
func (fe *FileEntry) SetTimestampTag(ctx context.Context, key string, value time.Time)
func (fe *FileEntry) GetTimestampTag(ctx context.Context, key string) (time.Time, bool)
func (fe *FileEntry) SetUUIDTag(ctx context.Context, key string, value uuid.UUID)
func (fe *FileEntry) GetUUIDTag(ctx context.Context, key string) (uuid.UUID, bool)
func (fe *FileEntry) SetURLTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetURLTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetEmailTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetEmailTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetVersionTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetVersionTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetHashTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetHashTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetPathTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetPathTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetMimeTypeTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetMimeTypeTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetLanguageTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetLanguageTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) SetNovusPackMetadataTag(ctx context.Context, key, value string)
func (fe *FileEntry) GetNovusPackMetadataTag(ctx context.Context, key string) (string, bool)
func (fe *FileEntry) GetInheritedTag(ctx context.Context, key string) (Tag, bool)
func (fe *FileEntry) RemoveTag(ctx context.Context, key string)
func (fe *FileEntry) HasTag(ctx context.Context, key string) bool
func (fe *FileEntry) GetTag(ctx context.Context, key string) (Tag, bool)
func (fe *FileEntry) HasTags(ctx context.Context) bool
func (fe *FileEntry) SyncTags(ctx context.Context)
func (fe *FileEntry) GetEffectiveTags(ctx context.Context) []Tag
```

#### 1.2.1.1 Generic Tag Types

Type-safe tag operations using generics for improved type safety and code reuse.

```go
// TypedTag represents a type-safe tag with a specific value type
type TypedTag[T any] struct {
    Key   string
    Value T
    Type  TagValueType
}

func NewTypedTag[T any](key string, value T, tagType TagValueType) *TypedTag[T]
func (t *TypedTag[T]) GetValue() T
func (t *TypedTag[T]) SetValue(value T)
```

#### 1.2.1.2 Generic Tag Operations

Type-safe tag operations for FileEntry objects.

```go
// GetTypedTag retrieves a type-safe tag value
func (fe *FileEntry) GetTypedTag[T any](ctx context.Context, key string) (T, bool)

// SetTypedTag sets a type-safe tag value
func (fe *FileEntry) SetTypedTag[T any](ctx context.Context, key string, value T, tagType TagValueType) error

// GetTagAs retrieves a tag and converts it to the specified type
func (fe *FileEntry) GetTagAs[T any](ctx context.Context, key string, converter func(interface{}) (T, error)) (T, error)
```

### 1.2.2 Data Management

```go
func (fe *FileEntry) LoadData(ctx context.Context) error
func (fe *FileEntry) UnloadData(ctx context.Context)
func (fe *FileEntry) GetData(ctx context.Context) ([]byte, error)
func (fe *FileEntry) SetData(ctx context.Context, data []byte)
func (fe *FileEntry) CreateTempFile(ctx context.Context) error
func (fe *FileEntry) StreamToTempFile(ctx context.Context) error
func (fe *FileEntry) WriteToTempFile(ctx context.Context, data []byte) error
func (fe *FileEntry) ReadFromTempFile(ctx context.Context, offset int64, size int64) ([]byte, error)
func (fe *FileEntry) CleanupTempFile(ctx context.Context) error
func (fe *FileEntry) GetProcessingState(ctx context.Context) ProcessingState
func (fe *FileEntry) SetProcessingState(ctx context.Context, state ProcessingState)
func (fe *FileEntry) SetSourceFile(ctx context.Context, file *os.File, offset, size int64)
func (fe *FileEntry) GetSourceFile(ctx context.Context) (*os.File, int64, int64)
func (fe *FileEntry) SetTempPath(ctx context.Context, path string)
func (fe *FileEntry) GetTempPath(ctx context.Context) string
```

### 1.2.3 Path and Directory Management

```go
func (fe *FileEntry) HasSymlinks(ctx context.Context) bool
func (fe *FileEntry) GetSymlinkPaths(ctx context.Context) []PathEntry
func (fe *FileEntry) GetPrimaryPath(ctx context.Context) string
func (fe *FileEntry) ResolveAllSymlinks(ctx context.Context) []string
func (fe *FileEntry) SetParentDirectory(ctx context.Context, parent *DirectoryEntry)
func (fe *FileEntry) GetParentDirectory(ctx context.Context) *DirectoryEntry
func (fe *FileEntry) GetParentPath(ctx context.Context) string
func (fe *FileEntry) GetDirectoryDepth(ctx context.Context) int
func (fe *FileEntry) IsRootRelative(ctx context.Context) bool
func (fe *FileEntry) GetAncestorDirectories(ctx context.Context) []*DirectoryEntry
func (fe *FileEntry) GetParentDirectoryPath(ctx context.Context) string
func (fe *FileEntry) SetInheritedTags(ctx context.Context, tags []Tag)
func (fe *FileEntry) GetInheritedTags(ctx context.Context) []Tag
func (fe *FileEntry) UpdateInheritedTags(ctx context.Context, tags []Tag)
func (fe *FileEntry) ClearDirectoryAssociations(ctx context.Context)
```

### 1.2.4 Serialization

```go
func ParseFileEntry(data []byte) (*FileEntry, error)
```

### 1.2.5 FileEntry Directory Association Methods

Cross-Reference: For directory metadata structures and package-level directory management, see [Package Metadata API - Directory Metadata System](api_metadata.md#8-directory-metadata-system).

```go
// FileEntry directory association methods
// NewFileEntry creates a new FileEntry with proper tag synchronization
func NewFileEntry() *FileEntry

// LoadFileEntry creates a FileEntry from package data with proper tag synchronization
func LoadFileEntry(data []byte) (*FileEntry, error)
```

## 2. Add File Operations

### 2.1 Add File

```go
// AddFile adds a file to the package from various sources
func (p *Package) AddFile(ctx context.Context, path string, source FileSource, options *AddFileOptions) (*FileEntry, error)
```

#### 2.1.1 Purpose

Adds a file to the package from any data source using a unified interface and returns the created FileEntry.

#### 2.1.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Package-internal path for the file
- `source`: FileSource providing the file data
- `options`: Optional configuration for file processing (can be nil for defaults)

#### 2.1.3 Returns

- `*FileEntry`: The created file entry with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file addition

#### 2.1.4 Behavior

- Reads file content from the provided FileSource
- Uses streaming for large files to manage memory efficiently (when supported by source)
- Automatically determines file type based on extension and content analysis
- Creates a new file entry in the package with complete metadata

#### 2.1.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrInvalidPath`: Invalid or malformed file path
- `ErrFileExists`: File already exists at the specified path
- `ErrContentTooLarge`: File content exceeds size limits
- `ErrIOError`: I/O error during file operations
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

### 2.2 Add File Pattern

```go
// AddFilePattern adds files matching a pattern with options
func (p *Package) AddFilePattern(ctx context.Context, pattern string, options *AddFileOptions) ([]*FileEntry, error)
```

#### 2.2.1 Purpose

Adds multiple files to the package based on a file system pattern and returns the created FileEntry objects.

#### 2.2.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `pattern`: File system pattern (e.g., "_.txt", "documents/\*\*/_.pdf")
- `options`: Configuration options for file processing (can be nil for defaults)

#### 2.2.3 Returns

- `[]*FileEntry`: Slice of created file entries with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file addition (if error occurs, some files may have been added successfully)

#### 2.2.4 Behavior

- Scans file system for files matching the pattern
- Applies pattern-specific filters (exclude patterns, max file size)
- Adds each matching file to the package with specified options
- Preserves directory structure if requested

#### 2.2.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrInvalidPattern`: Invalid or malformed pattern
- `ErrNoFilesFound`: No files match the pattern
- `ErrIOError`: I/O error during file operations

#### 2.2.6 Usage Notes

AddFilePattern accepts AddFileOptions to configure encryption, compression, and filtering behavior. See AddFileOptions Configuration for available fields.

### 2.3 FileSource Interface

```go
// FileSource represents the source of file data
type FileSource interface {
    // Read reads file data into the provided buffer
    Read(ctx context.Context, p []byte) (n int, err error)
    // Size returns the total size of the file data
    Size(ctx context.Context) (int64, error)
    // Close closes the source
    Close(ctx context.Context) error
}
```

#### 2.3.1 Purpose

Defines the interface for providing file data from various sources.

#### 2.3.2 Methods

- `Read(ctx, p)`: Reads file data into the provided buffer
- `Size(ctx)`: Returns the total size of the file data
- `Close(ctx)`: Closes the source and releases resources

### 2.4 FileSource Implementations

```go
// FilePathSource implements FileSource for filesystem files
type FilePathSource struct {
    Path string
    file *os.File
}

// MemorySource implements FileSource for in-memory data
type MemorySource struct {
    Data []byte
    pos  int
}

// NewFilePathSource creates a FileSource from a filesystem path
func NewFilePathSource(ctx context.Context, path string) (FileSource, error)

// NewMemorySource creates a FileSource from byte data
func NewMemorySource(ctx context.Context, data []byte) FileSource
```

#### 2.4.1 Purpose

Built-in implementations of the FileSource interface.

#### 2.4.2 FilePathSource

- Reads from filesystem files with streaming support
- Automatically determines file type from extension and content
- Handles large files efficiently with memory management

#### 2.4.3 MemorySource

- Reads from in-memory byte data
- Suitable for small files or generated content
- No filesystem overhead

### 2.5 AddFileOptions Configuration

```go
// AddFileOptions configures file addition behavior for both individual files and patterns
type AddFileOptions struct {
    // File processing options
    Compress        Option[bool]            // Whether to compress the file
    CompressionType Option[uint8]           // Compression algorithm (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    CompressionLevel Option[int]            // Compression level (1-9, 0 = default)
    FileType        Option[uint16]          // File type identifier
    Tags            Option[map[string]interface{}] // Per-file tags (key-value pairs)

    // Encryption options
    Encrypt         Option[bool]            // Whether to encrypt the file
    EncryptionType  Option[EncryptionType]  // Encryption algorithm type
    EncryptionKey   Option[*EncryptionKey]  // Specific encryption key (overrides EncryptionType)

    // Pattern-specific options (only used for pattern operations)
    ExcludePatterns Option[[]string]        // Patterns to exclude from processing
    MaxFileSize     Option[int64]           // Maximum file size to include (0 = no limit)
    PreservePaths   Option[bool]            // Whether to preserve directory structure
}
```

#### 2.5.1 Purpose

Unified configuration options for all file addition operations, supporting both individual files and pattern-based operations.

#### 2.5.2 Fields

##### 2.5.2.1 File Processing Options

- `Compress`: Whether to compress the file (default: false)
- `CompressionType`: Compression algorithm identifier (default: 0 = none, 1=Zstd, 2=LZ4, 3=LZMA)
- `CompressionLevel`: Compression level 1-9 (default: 0 = default)
- `FileType`: File type identifier (default: 0 = regular file)
- `Tags`: Per-file tags as key-value pairs (default: nil)

##### 2.5.2.2 Encryption Options

- `Encrypt`: Whether to encrypt the file (default: false)
- `EncryptionType`: Encryption algorithm type (default: EncryptionNone)
- `EncryptionKey`: Specific encryption key (overrides EncryptionType)

##### 2.5.2.3 Pattern-Specific Options

- `ExcludePatterns`: Patterns to exclude from processing (default: nil)
- `MaxFileSize`: Maximum file size to include, 0 = no limit (default: 0)
- `PreservePaths`: Whether to preserve directory structure (default: false)

### 2.6 Usage Notes

AddFile supports various FileSource implementations including filesystem files, memory data, and custom sources. Use AddFileOptions to configure compression, encryption, and metadata.

## 3. File Addition Implementation Flow

### 3.1 Processing Order Requirements

The file addition process must follow a specific sequence to ensure proper compression, encryption, and deduplication:

#### 3.1.1 Required Processing Sequence

1. File Validation

   - Check file exists and is not a directory
   - Validate file name and path format
   - Verify file size limits and permissions

2. Compression (if requested)

   - Apply compression algorithm if requested
   - Compression may fail and must return error
   - No fallback to uncompressed storage

3. Encryption (if requested)

   - Apply encryption algorithm if requested
   - Encryption may fail and must return error
   - No fallback to unencrypted storage

4. Deduplication Check

   - Check for existing processed content using:
     - Processed file size (after compression/encryption)
     - CRC32 checksum of processed content
     - SHA-256 hash of processed content (if size and CRC32 match)
   - Use processed content, not raw file content

5. Storage Decision
   - If unique: Store file with new FileEntry
   - If duplicate: Add path reference to existing FileEntry

#### 3.1.2 Error Handling Requirements

- **Compression failures**: Must prevent file addition and return appropriate error
- **Encryption failures**: Must prevent file addition and return appropriate error
- **Resource cleanup**: Failed operations must properly clean up allocated resources
- **User feedback**: Provide clear error messages explaining failures and recovery options

#### 3.1.3 Performance Requirements

- **Deduplication efficiency**: Use processed size and CRC32 as early elimination filters
- **SHA-256 optimization**: Only compute expensive SHA-256 when size and CRC32 match
- **Memory management**: Handle large files efficiently with streaming when needed

#### 3.1.4 Implementation Function Signature

```go
func (p *Package) AddFile(ctx context.Context, path string, source FileSource, options *AddFileOptions) (*FileEntry, error)
```

## 4. Remove File Operations

### 4.1 Remove File

```go
// RemoveFile removes a file from the package by FileEntry reference
func (p *Package) RemoveFile(ctx context.Context, entry *FileEntry) error

// RemoveFileByPath removes a file from the package by path
func (p *Package) RemoveFileByPath(ctx context.Context, path string) error
```

#### 3.1.1 Purpose

Removes a file from the package using either a FileEntry reference or a virtual path.

#### 3.1.2 Parameters

#### 3.1.2.1 RemoveFile

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file to remove

#### 3.1.2.2 RemoveFileByPath

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to remove

#### 3.1.3 Behavior

- Removes file entry from the package index
- Marks file data as deleted (space reclaimed during defragmentation)
- Updates package metadata and file count
- Preserves package integrity and signatures

#### 3.1.4 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrFileNotFound`: File does not exist at the specified path
- `ErrInvalidPath`: Invalid or malformed file path
- `ErrPackageReadOnly`: Package is in read-only mode
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 3.1.5 Usage Notes

RemoveFile accepts either a FileEntry reference or a path string. Use GetFileByPath to obtain FileEntry references when needed.

### 4.2 Remove File Pattern

```go
// RemoveFilePattern removes files matching a pattern from the package
func (p *Package) RemoveFilePattern(ctx context.Context, pattern string) ([]*FileEntry, error)
```

#### 3.2.1 Purpose

Removes multiple files from the package based on a file system pattern and returns the removed FileEntry objects.

#### 3.2.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `pattern`: File system pattern (e.g., "\*.txt", "documents/\*\*/\*.pdf")

#### 3.2.3 Returns

- `[]*FileEntry`: Slice of removed file entries
- `error`: Any error that occurred during file removal (if error occurs, some files may have been removed successfully)

#### 3.2.4 Behavior

- Scans package for files matching the pattern
- Removes each matching file from the package
- Marks file data as deleted (space reclaimed during defragmentation)
- Updates package metadata and file count
- Preserves package integrity and signatures

#### 3.2.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrInvalidPattern`: Invalid or malformed pattern
- `ErrNoFilesFound`: No files match the pattern
- `ErrPackageReadOnly`: Package is in read-only mode
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 3.2.6 Usage Notes

RemoveFilePattern returns a slice of removed FileEntry objects for further processing or logging.

## 5. Extract File Operations

### 5.1 Extract File

```go
// ExtractFile extracts file content from the package
func (p *Package) ExtractFile(ctx context.Context, path string) ([]byte, error)
```

#### 4.1.1 Purpose

Extracts and returns the content of a file from the package.

#### 4.1.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to extract

#### 4.1.3 Returns

File content as byte slice

#### 4.1.4 Behavior

- Locates file entry in the package index
- Reads file content from the data section
- Decompresses content if necessary
- Decrypts content if encrypted
- Returns raw file content

#### 4.1.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrFileNotFound`: File does not exist at the specified path
- `ErrInvalidPath`: Invalid or malformed file path
- `ErrDecryptionFailed`: Failed to decrypt encrypted file
- `ErrDecompressionFailed`: Failed to decompress file content
- `ErrIOError`: I/O error during file extraction
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 4.1.6 Usage Notes

ExtractFile returns the raw file content as a byte slice, automatically handling decompression and decryption.

## 6. Update File Operations

### 6.1 Update File

```go
// UpdateFile updates file content and metadata in the package
func (p *Package) UpdateFile(ctx context.Context, entry *FileEntry, source FileSource, options *AddFileOptions) (*FileEntry, error)
```

#### 5.1.1 Purpose

Updates an existing file's content and metadata in the package, returning the updated FileEntry.

#### 5.1.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file to update
- `source`: FileSource interface providing new file data
- `options`: Configuration options for file update (can be nil for defaults)

#### 5.1.3 Returns

- `*FileEntry`: The updated file entry with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file update

#### 5.1.4 Behavior

- Reads new file content from the provided FileSource
- Updates file entry metadata (size, checksums, timestamps)
- Reapplies compression and encryption settings from options
- Preserves file path and basic metadata unless explicitly changed
- Updates package metadata and file count
- Validates content size limits
- Performs deduplication checks
- Automatically closes the FileSource when done

#### 5.1.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrFileNotFound`: File entry does not exist or is invalid
- `ErrContentTooLarge`: File content exceeds size limits
- `ErrUnsupportedEncryption`: Unsupported encryption type
- `ErrEncryptionFailed`: Failed to encrypt file content
- `ErrIOError`: I/O error during file update
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 5.1.6 Usage Notes

UpdateFile updates file content and metadata, supporting various FileSource implementations and AddFileOptions for configuration.

### 6.2 Update File Pattern

```go
// UpdateFilePattern updates files matching a pattern in the package
func (p *Package) UpdateFilePattern(ctx context.Context, pattern string, sourceDir string, options *AddFileOptions) ([]*FileEntry, error)
```

#### 5.2.1 Purpose

Updates multiple files in the package based on a file system pattern and returns the updated FileEntry objects.

#### 5.2.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `pattern`: File system pattern (e.g., "\*.txt", "documents/\*\*/\*.pdf")
- `sourceDir`: Base directory to search for updated files
- `options`: Configuration options for file processing (can be nil for defaults)

#### 5.2.3 Returns

- `[]*FileEntry`: Slice of updated file entries with all metadata, compression status, encryption details, and checksums
- `error`: Any error that occurred during file update (if error occurs, some files may have been updated successfully)

#### 5.2.4 Behavior

- Scans file system for files matching the pattern
- Finds corresponding files in the package
- Updates each matching file with new content from filesystem
- Applies pattern-specific filters (exclude patterns, max file size)
- Preserves directory structure if requested
- Reports progress for large file sets
- Performs deduplication checks

#### 5.2.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrInvalidPattern`: Invalid or malformed pattern
- `ErrNoFilesFound`: No files match the pattern
- `ErrIOError`: I/O error during file operations
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 5.2.6 Usage Notes

UpdateFilePattern updates files matching a pattern, allowing specification of a source directory and AddFileOptions.

### 6.3 Update File Metadata

```go
// UpdateFileMetadata updates file metadata without changing content
func (p *Package) UpdateFileMetadata(ctx context.Context, entry *FileEntry, metadata *FileMetadataUpdate) (*FileEntry, error)
```

#### 5.3.1 Purpose

Updates file metadata (tags, attributes, compression settings, etc.) without modifying file content.

#### 5.3.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file to update
- `metadata`: FileMetadataUpdate structure containing new metadata

#### 5.3.3 FileMetadataUpdate Structure

```go
type FileMetadataUpdate struct {
    // Basic metadata
    Tags            []string        // File tags
    CompressionType uint8           // New compression type
    CompressionLevel uint8          // New compression level
    EncryptionType  uint8           // New encryption type

    // Path management
    AddPaths        []PathEntry     // Additional paths to add
    RemovePaths     []string        // Paths to remove (by path string)
    UpdatePaths     []PathEntry     // Paths to update

    // Hash management
    AddHashes       []HashEntry     // Additional hashes to add
    RemoveHashes    []HashType      // Hash types to remove
    UpdateHashes    []HashEntry     // Hashes to update

    // Optional data
    OptionalData    OptionalData   // Structured optional data
}
```

#### 5.3.4 Returns

- `*FileEntry`: The updated file entry with new metadata
- `error`: Any error that occurred during metadata update

#### 5.3.5 Behavior

- Updates file metadata fields without changing content
- Increments MetadataVersion field
- Preserves file content and data integrity
- Updates package metadata if needed
- Validates new compression/encryption settings
- Applies new compression/encryption to existing content if settings changed
- Manages multiple paths:
  - Adds new paths with per-path metadata (permissions, ownership, timestamps)
  - Removes specified paths
  - Updates existing path metadata
- Manages hash data:
  - Adds new hash entries with specified type and purpose
  - Removes hash entries by type
  - Updates existing hash entries
- Manages optional data:
  - Updates tags data (DataType 0x00)
  - Sets path encoding and flags (DataTypes 0x01-0x02)
  - Updates compression dictionary and solid group IDs (DataTypes 0x03-0x04)
  - Sets file system flags and Windows attributes (DataTypes 0x05-0x06)
  - Updates extended attributes and ACL data (DataTypes 0x07-0x08)

#### 5.3.6 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrFileNotFound`: File entry does not exist or is invalid
- `ErrUnsupportedCompression`: Unsupported compression type
- `ErrUnsupportedEncryption`: Unsupported encryption type
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 5.3.7 Usage Notes

UpdateFileMetadata accepts a FileMetadataUpdate structure to modify tags, attributes, compression settings, and optional data without changing file content.

### 6.4 Add File Path

```go
// AddFilePath adds an additional path to an existing file entry
func (p *Package) AddFilePath(ctx context.Context, entry *FileEntry, path PathEntry) (*FileEntry, error)
```

#### 5.4.1 Purpose

Adds an additional path to an existing file entry, enabling multiple paths to point to the same content.

#### 5.4.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file
- `path`: PathEntry with path and metadata

#### 5.4.3 Returns

- `*FileEntry`: The updated file entry with additional path
- `error`: Any error that occurred during path addition

### 6.5 Remove File Path

```go
// RemoveFilePath removes a path from an existing file entry
func (p *Package) RemoveFilePath(ctx context.Context, entry *FileEntry, path string) (*FileEntry, error)
```

#### 5.5.1 Purpose

Removes a specific path from an existing file entry.

#### 5.5.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file
- `path`: Path string to remove

#### 5.5.3 Returns

- `*FileEntry`: The updated file entry with path removed
- `error`: Any error that occurred during path removal

### 6.6 Add File Hash

```go
// AddFileHash adds a hash entry to an existing file entry
func (p *Package) AddFileHash(ctx context.Context, entry *FileEntry, hash HashEntry) (*FileEntry, error)
```

#### 5.6.1 Purpose

Adds a hash entry to an existing file entry for content verification, deduplication, or integrity checking.

#### 5.6.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `entry`: FileEntry reference to the file
- `hash`: HashEntry with type, purpose, and data

#### 5.6.3 Returns

- `*FileEntry`: The updated file entry with additional hash
- `error`: Any error that occurred during hash addition

## 7. File Compression Operations

```go
// CompressFile compresses an existing file in the package
func (p *Package) CompressFile(ctx context.Context, path string, compressionType uint8) error

// DecompressFile decompresses an existing file in the package
func (p *Package) DecompressFile(ctx context.Context, path string) error

// GetFileCompressionInfo gets compression information for a file
func (p *Package) GetFileCompressionInfo(ctx context.Context, path string) (*FileCompressionInfo, error)

// FileCompressionInfo contains file compression details
type FileCompressionInfo struct {
    IsCompressed     bool
    CompressionType  uint8
    OriginalSize     int64
    CompressedSize   int64
    CompressionRatio float64
}
```

### 7.1 Purpose

Manages file-level compression operations on existing files in the package.

### 7.1.1 CompressFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to compress
- `compressionType`: Compression algorithm to use

#### 6.1.2 DecompressFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to decompress

#### 6.1.3 GetFileCompressionInfo Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Virtual path of the file to inspect

#### 6.1.4 Behavior

- Compresses/decompresses file content in-place
- Updates file entry metadata
- Preserves file integrity and signatures
- Maintains package structure

#### 6.1.5 Error Conditions

- `ErrPackageNotOpen`: Package is not currently open
- `ErrFileNotFound`: File does not exist at the specified path
- `ErrInvalidPath`: Invalid or malformed file path
- `ErrFileAlreadyCompressed`: File is already compressed
- `ErrFileNotCompressed`: File is not compressed
- `ErrCompressionFailed`: Failed to compress file content
- `ErrDecompressionFailed`: Failed to decompress file content
- `ErrContextCancelled`: Context was cancelled
- `ErrContextTimeout`: Context timeout exceeded

#### 6.1.6 Usage Notes

CompressFile compresses the content of an existing file within the package using the specified compression type.

## 8. File Encryption Operations

File encryption operations are defined in the [Security Validation API](api_security.md#45-package-file-encryption-operations). The file management API provides access to these operations through the Package interface.

## 9. Deduplication Operations

### 9.1 File Deduplication

```go
// FindExistingEntryByCRC32 finds existing entry by CRC32 checksum
func (p *Package) FindExistingEntryByCRC32(rawChecksum uint32) *FileEntry

// FindExistingEntryMultiLayer performs multi-layer deduplication
func (p *Package) FindExistingEntryMultiLayer(originalSize int64, rawChecksum uint32, content []byte) (*FileEntry, []byte)

// AddPathToExistingEntry adds a path to an existing entry
func (p *Package) AddPathToExistingEntry(existingEntry *FileEntry, newPath string)
```

#### 9.1.1 Purpose

Provides deduplication functionality to avoid storing duplicate content.

#### 9.1.2 FindExistingEntryByCRC32 Parameters

- `rawChecksum`: CRC32 checksum to search for

#### 9.1.3 FindExistingEntryMultiLayer Parameters

- `originalSize`: Original file size
- `rawChecksum`: CRC32 checksum
- `content`: File content for verification

#### 9.1.4 AddPathToExistingEntry Parameters

- `existingEntry`: Existing file entry
- `newPath`: New path to add

#### 9.1.5 Behavior

- Searches for existing files with matching checksums
- Performs multi-layer verification (CRC32 + content hash)
- Adds new paths to existing entries when duplicates found
- Reduces storage space by sharing content

#### 9.1.6 Usage Notes

Deduplication functions support both simple CRC32-based lookup and multi-layer verification for accurate duplicate detection.

## 10. File Information and Queries

### 10.1 File Existence and Properties

```go
// FileExists checks if a file with the given path exists in the package
func (p *Package) FileExists(ctx context.Context, path string) (bool, error)

// ListFiles returns all file entries in the package
func (p *Package) ListFiles(ctx context.Context) ([]*FileEntry, error)


// GetFileByPath gets a file entry by path
func (p *Package) GetFileByPath(ctx context.Context, path string) (*FileEntry, bool)

// GetFileByOffset gets a file entry by offset
func (p *Package) GetFileByOffset(ctx context.Context, offset int64) (*FileEntry, bool)

// GetFileByFileID gets a file entry by its unique FileID
func (p *Package) GetFileByFileID(ctx context.Context, fileID uint64) (*FileEntry, bool)

// GetFileByHash gets a file entry by content hash
func (p *Package) GetFileByHash(ctx context.Context, hashType HashType, hashData []byte) (*FileEntry, bool)

// GetFileByChecksum gets a file entry by CRC32 checksum
func (p *Package) GetFileByChecksum(ctx context.Context, checksum uint32) (*FileEntry, bool)

// FindEntriesByTag finds all file entries with a specific tag
func (p *Package) FindEntriesByTag(ctx context.Context, tag string) ([]*FileEntry, error)

// FindEntriesByType finds all file entries of a specific type
func (p *Package) FindEntriesByType(ctx context.Context, fileType uint16) ([]*FileEntry, error)

// GetFileCount returns the total number of files in the package
func (p *Package) GetFileCount(ctx context.Context) (int, error)

// GetPatterns gets files matching patterns from the package
func (p *Package) FindEntriesByPathPatterns(ctx context.Context, patterns []string) ([]*FileEntry, error)

// ListCompressedFiles returns all compressed file entries
func (p *Package) ListCompressedFiles(ctx context.Context) ([]*FileEntry, error)

// ListEncryptedFiles returns all encrypted file entries
func (p *Package) ListEncryptedFiles(ctx context.Context) ([]*FileEntry, error)
```

#### 10.1.1 Purpose

Provides file information and querying capabilities.

#### 10.1.2 FileEntry Access

All file query functions return `*FileEntry` objects or `[]*FileEntry` arrays, which provide comprehensive file information including all metadata, compression status, encryption details, checksums, and timestamps. The `GetFileByPath()` function replaces the need for a separate `GetFileInfo()` function by returning the complete `FileEntry` directly.

#### 10.1.3 Usage Notes

File query functions provide comprehensive FileEntry objects. GetFileByPath replaces GetFileInfo.

### 10.2 File Lookup by Metadata

#### 9.2.1 GetFileByFileID

```go
// GetFileByFileID gets a file entry by its unique FileID
func (p *Package) GetFileByFileID(ctx context.Context, fileID uint64) (*FileEntry, bool)
```

#### 9.2.1 Purpose

Finds a file entry by its unique 64-bit FileID.

#### 9.2.1.1 Parameters

- `ctx`: Context for cancellation and timeout handling
- `fileID`: Unique 64-bit file identifier

#### 9.2.1.2 Returns

- `*FileEntry`: The file entry with the specified FileID
- `bool`: True if found, false otherwise

#### 9.2.1.3 Use Cases

- Stable file references across package modifications
- Database-style lookups by primary key
- File tracking and management systems

#### 9.2.2 GetFileByHash

```go
// GetFileByHash gets a file entry by content hash
func (p *Package) GetFileByHash(ctx context.Context, hashType HashType, hashData []byte) (*FileEntry, bool)
```

#### 9.2.2.1 Purpose

Finds a file entry by its content hash (SHA-256, SHA-512, BLAKE3, or XXH3).

#### 9.2.2.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `hashType`: Type of hash algorithm used
- `hashData`: Hash data to search for

#### 9.2.2.3 Returns

- `*FileEntry`: The file entry with matching hash
- `bool`: True if found, false otherwise

#### 9.2.2.4 Use Cases

- Content deduplication
- Integrity verification
- Finding files by content rather than path

#### 9.2.3 GetFileByChecksum

```go
// GetFileByChecksum gets a file entry by CRC32 checksum
func (p *Package) GetFileByChecksum(ctx context.Context, checksum uint32) (*FileEntry, bool)
```

#### 9.2.3.1 Purpose

Finds a file entry by its CRC32 checksum (fast lookup).

#### 9.2.3.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `checksum`: CRC32 checksum value to search for

#### 9.2.3.3 Returns

- `*FileEntry`: The file entry with matching checksum
- `bool`: True if found, false otherwise

#### 9.2.3.4 Use Cases

- Fast content identification
- Quick duplicate detection
- Lightweight file matching

#### 9.2.4 FindEntriesByTag

```go
// FindEntriesByTag finds all file entries with a specific tag
func (p *Package) FindEntriesByTag(ctx context.Context, tag string) ([]*FileEntry, error)
```

#### 9.2.4.1 Purpose

Finds all file entries that have a specific tag.

#### 9.2.4.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `tag`: Tag string to search for

#### 9.2.4.3 Returns

- `[]*FileEntry`: All file entries with the specified tag
- `error`: Any error that occurred during the search

#### 9.2.4.4 Use Cases

- Finding all files with a specific label
- Organizing files by category
- Tag-based file management

#### 9.2.5 FindEntriesByType

```go
// FindEntriesByType finds all file entries of a specific type
func (p *Package) FindEntriesByType(ctx context.Context, fileType uint16) ([]*FileEntry, error)
```

#### 9.2.5.1 Purpose

Finds all file entries of a specific file type.

#### 9.2.5.2 Parameters

- `ctx`: Context for cancellation and timeout handling
- `fileType`: File type identifier (0-65535)

#### 9.2.5.3 Returns

- `[]*FileEntry`: All file entries of the specified type
- `error`: Any error that occurred during the search

#### 9.2.5.4 Use Cases

- Finding all files of a specific format
- Type-based file processing
- File organization by category

#### 9.2.6 GetFileCount

```go
// GetFileCount returns the total number of files in the package
func (p *Package) GetFileCount(ctx context.Context) (int, error)
```

#### 9.2.6.1 Purpose

Returns the total number of files in the package.

#### 9.2.6.2 Parameters

- `ctx`: Context for cancellation and timeout handling

#### 9.2.6.3 Returns

- `int`: Total number of files in the package
- `error`: Any error that occurred

#### 9.2.6.4 Use Cases

- Package statistics
- Progress tracking
- Validation and bounds checking

## 11. FileEntry Methods

### 11.1 File Entry Properties

```go
// IsCompressed checks if the file is compressed
func (entry *FileEntry) IsCompressed() bool

// HasEncryptionKey checks if the file has an encryption key set
func (entry *FileEntry) HasEncryptionKey() bool

// GetEncryptionType returns the encryption type used for this file
func (entry *FileEntry) GetEncryptionType() EncryptionType

// IsEncrypted checks if the file is encrypted
func (entry *FileEntry) IsEncrypted() bool

// ToBinaryFormat converts the file entry to binary format
func (entry *FileEntry) ToBinaryFormat() (*FileEntryBinaryFormat, error)
```

#### 11.1.1 Purpose

Provides access to file entry properties and serialization.

#### 11.1.2 IsCompressed Returns

`bool` indicating if the file is compressed

#### 11.1.3 HasEncryptionKey Returns

`bool` indicating if the file has an encryption key

#### 11.1.4 GetEncryptionType Returns

`EncryptionType` value indicating the encryption algorithm used (e.g., EncryptionNone, EncryptionAES256GCM, EncryptionMLKEM)

#### 11.1.5 IsEncrypted Returns

`bool` indicating if the file is encrypted (equivalent to `GetEncryptionType() != EncryptionNone`)

#### 11.1.6 ToBinaryFormat Returns

`FileEntryBinaryFormat` structure and error

#### 11.1.7 Usage Notes

ToBinaryFormat converts the FileEntry to a binary format for storage or transmission.

### 11.2 File Entry Encryption

```go
// SetEncryptionKey sets the encryption key for the file
func (entry *FileEntry) SetEncryptionKey(key *EncryptionKey) error

// Encrypt encrypts data using the file's encryption key
func (entry *FileEntry) Encrypt(data []byte) ([]byte, error)

// Decrypt decrypts data using the file's encryption key
func (entry *FileEntry) Decrypt(data []byte) ([]byte, error)

// UnsetEncryptionKey removes the encryption key from the file
func (entry *FileEntry) UnsetEncryptionKey()
```

#### 11.2.1 Purpose

Manages encryption for individual file entries.

#### 11.2.2 SetEncryptionKey Parameters

- `key`: Encryption key to set

#### 11.2.3 Encrypt/Decrypt Parameters

- `data`: Data to encrypt or decrypt

#### 11.2.4 Error Conditions

- `ErrInvalidKey`: Invalid encryption key
- `ErrEncryptionFailed`: Encryption operation failed
- `ErrDecryptionFailed`: Decryption operation failed

#### 11.2.5 Usage Notes

SetEncryptionKey and UnsetEncryptionKey manage the encryption key for the file entry.

### 11.3 File Entry Data Management

```go
// LoadData loads the file data into memory
func (entry *FileEntry) LoadData() error

// ProcessData processes the file data (compression, encryption, etc.)
func (entry *FileEntry) ProcessData() error
```

#### 11.3.1 Purpose

Manages file data loading and processing.

#### 11.3.2 LoadData Behavior

- Loads file content from package into memory
- Prepares data for access and processing
- May trigger decompression or decryption

#### 11.3.3 ProcessData Behavior

- Applies compression and encryption to file data
- Updates file entry metadata
- Prepares data for storage

#### 11.3.4 Error Conditions

- `ErrIOError`: I/O error during data operations
- `ErrDecryptionFailed`: Failed to decrypt data
- `ErrDecompressionFailed`: Failed to decompress data

#### 11.3.5 Usage Notes

LoadData and UnloadData manage the file's content in memory, while ProcessData applies compression/encryption.

## 12. Error Handling

### 12.1 Structured Error System

The NovusPack file management API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities. For complete error system documentation, see [Structured Error System](api_core.md#10-structured-error-system).

### 12.2 Common Error Types

#### 12.2.1 Sentinel Errors (Legacy Support)

```go
var (
    ErrFileNotFound        = errors.New("file not found")
    ErrFileExists          = errors.New("file already exists")
    ErrInvalidPath         = errors.New("invalid file path")
    ErrInvalidPattern      = errors.New("invalid file pattern")
    ErrContentTooLarge     = errors.New("file content too large")
    ErrNoFilesFound        = errors.New("no files found matching pattern")
    ErrUnsupportedEncryption = errors.New("unsupported encryption type")
    ErrEncryptionFailed    = errors.New("encryption failed")
    ErrDecryptionFailed    = errors.New("decryption failed")
    ErrDecompressionFailed = errors.New("decompression failed")
    ErrInvalidSecurityLevel = errors.New("invalid security level")
    ErrKeyGenerationFailed = errors.New("key generation failed")
    ErrInvalidKey          = errors.New("invalid key")
    ErrPackageNotOpen      = errors.New("package is not open")
    ErrPackageReadOnly     = errors.New("package is read-only")
    ErrIOError            = errors.New("I/O error")
    ErrContextCancelled   = errors.New("context cancelled")
    ErrContextTimeout     = errors.New("context timeout")
)
```

#### 12.2.2 Error Type Mapping

| Sentinel Error           | ErrorType          | Description                     |
| ------------------------ | ------------------ | ------------------------------- |
| ErrFileNotFound          | ErrTypeValidation  | File not found in package       |
| ErrFileExists            | ErrTypeValidation  | File already exists             |
| ErrInvalidPath           | ErrTypeValidation  | Invalid file path               |
| ErrInvalidPattern        | ErrTypeValidation  | Invalid file pattern            |
| ErrContentTooLarge       | ErrTypeValidation  | File content too large          |
| ErrNoFilesFound          | ErrTypeValidation  | No files found matching pattern |
| ErrUnsupportedEncryption | ErrTypeEncryption  | Unsupported encryption type     |
| ErrEncryptionFailed      | ErrTypeEncryption  | Encryption operation failed     |
| ErrDecryptionFailed      | ErrTypeEncryption  | Decryption operation failed     |
| ErrDecompressionFailed   | ErrTypeCompression | Decompression operation failed  |
| ErrInvalidSecurityLevel  | ErrTypeSecurity    | Invalid security level          |
| ErrKeyGenerationFailed   | ErrTypeEncryption  | Key generation failed           |
| ErrInvalidKey            | ErrTypeEncryption  | Invalid encryption key          |
| ErrPackageNotOpen        | ErrTypeValidation  | Package is not open             |
| ErrPackageReadOnly       | ErrTypeValidation  | Package is read-only            |
| ErrIOError               | ErrTypeIO          | I/O error                       |
| ErrContextCancelled      | ErrTypeContext     | Context cancelled               |
| ErrContextTimeout        | ErrTypeContext     | Context timeout                 |

### 12.3 Structured Error Examples

#### 12.2.3 Creating File Management Errors

```go
// File not found with context
err := NewPackageError(ErrTypeValidation, "file not found", ErrFileNotFound).
    WithContext("path", "/path/to/file").
    WithContext("operation", "ExtractFile")

// Encryption failure with context
err := NewPackageError(ErrTypeEncryption, "encryption failed", ErrEncryptionFailed).
    WithContext("algorithm", "AES-256-GCM").
    WithContext("keySize", 256).
    WithContext("file", "/path/to/file")

// Pattern matching error with context
err := NewPackageError(ErrTypeValidation, "no files found", ErrNoFilesFound).
    WithContext("pattern", "*.txt").
    WithContext("directory", "/src").
    WithContext("operation", "AddFilePattern")
```

### 12.4 Error Handling Best Practices

#### 12.4.1 Always check for errors

```go
err := package.AddFile(ctx, path, data)
if err != nil {
    return fmt.Errorf("failed to add file %s: %w", path, err)
}
```

#### 12.4.2 Use structured errors for better debugging

Use structured errors to provide context about operations and parameters for better debugging.

#### 12.4.3 Handle specific error types with structured errors

Handle different error types appropriately using structured error inspection and context.

#### 12.4.4 Use context for cancellation with structured errors

Use context for cancellation and timeout handling with appropriate error context.

#### 12.4.5 Handle encryption errors with context

Handle encryption errors with appropriate context and logging for debugging.

### 12.5 Generic FileEntry Operations

Generic collection operations for FileEntry objects using type-safe predicates and mappers.

```go
// FindFileEntry returns the first FileEntry matching the predicate
func FindFileEntry(entries []*FileEntry, predicate FilterFunc[*FileEntry]) (*FileEntry, bool)

// FilterFileEntries returns FileEntries matching the predicate
func FilterFileEntries(entries []*FileEntry, predicate FilterFunc[*FileEntry]) []*FileEntry

// MapFileEntries transforms FileEntries using the mapper function
func MapFileEntries[U any](entries []*FileEntry, mapper MapFunc[*FileEntry, U]) []U
```

## 13. Best Practices

### 13.1 File Path Management

#### 13.1.1 Use consistent path formats

```go
// Good: Use forward slashes, no leading slash
path := "documents/subfolder/file.txt"

// Bad: Mixed separators or leading slash
path := "/documents\\subfolder/file.txt"
```

#### 13.1.2 Validate paths before use

```go
if !isValidPath(path) {
    return fmt.Errorf("invalid file path: %s", path)
}
```

### 13.2 Encryption Management

#### 13.2.1 Choose appropriate encryption types

```go
// For sensitive data
err := package.AddFileWithEncryption(ctx, path, data, EncryptionAES256GCM)

// For post-quantum security
err := package.AddFileWithEncryption(ctx, path, data, EncryptionMLKEM)
```

#### 13.2.2 Secure key management

```go
key, err := GenerateMLKEMKey(ctx, 3)
if err != nil {
    return err
}
defer key.Clear() // Always clear sensitive data
```

### 13.3 Performance Considerations

#### 13.3.1 Use patterns for bulk operations

Use AddFilePattern for bulk operations instead of individual AddFile calls for better performance.

#### 13.3.2 Handle large files with streaming

For very large files, consider using the streaming API (see [Streaming and Buffer Management](api_streaming.md)).

#### 13.3.3 Use appropriate context timeouts

Use context.WithTimeout for operations that may take a long time to prevent indefinite blocking.
