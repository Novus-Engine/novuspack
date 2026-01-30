# NovusPack Technical Specifications - Package Compression API

## Table of Contents

- [Table of Contents](#table-of-contents)
- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
  - [0.2 Context Integration](#02-context-integration)
- [1. Package Compression Overview](#1-package-compression-overview)
  - [1.1 Compression Scope](#11-compression-scope)
    - [1.1.1 Compressed Content](#111-compressed-content)
    - [1.1.2 Uncompressed Content](#112-uncompressed-content)
  - [1.2 Compression Types](#12-compression-types)
  - [1.3 Compression Information Structure](#13-packagecompressioninfo-struct)
  - [1.4 Compression Constraints](#14-compression-constraints)
- [2. Strategy Pattern Interfaces](#2-strategy-pattern-interfaces)
  - [2.1 CompressionStrategy Interface](#21-compressionstrategy-interface)
    - [2.1.1 CompressionStrategy Interface Definition](#211-compressionstrategy-interface-definition)
    - [2.1.2 ByteCompressionStrategy Interface](#212-bytecompressionstrategy-interface)
    - [2.1.3 AdvancedCompressionStrategy Interface](#213-advancedcompressionstrategy-interface)
    - [2.1.4 StreamConfig Structure](#214-streamconfig-structure)
    - [2.1.5 MemoryStrategy Type](#215-memorystrategy-type)
  - [2.2 Built-in Compression Strategies](#22-built-in-compression-strategies)
    - [2.2.1 ZstandardStrategy Structure](#221-zstandardstrategy-structure)
    - [2.2.2 LZ4Strategy Structure](#222-lz4strategy-structure)
    - [2.2.3 LZMAStrategy Structure](#223-lzmastrategy-structure)
    - [2.2.4 CompressionJob Structure](#224-compressionjob-structure)
- [3. Interface Granularity and Composition](#3-interface-granularity-and-composition)
  - [3.1 Compression Information Interface](#31-compressioninfo-interface)
  - [3.2 CompressionOperations Interface](#32-compressionoperations-interface)
  - [3.3 CompressionStreaming Interface](#33-compressionstreaming-interface)
  - [3.4 CompressionFileOperations Interface](#34-compressionfileoperations-interface)
  - [3.5 Generic Compression Interface](#35-generic-compression-interface)
- [4. In-Memory Compression Methods](#4-in-memory-compression-methods)
  - [4.1 CompressPackage](#41-packagecompresspackage-method)
    - [4.1.1 CompressPackage Purpose](#411-compresspackage-purpose)
    - [4.1.2 CompressPackage Parameters](#412-compresspackage-parameters)
    - [4.1.3 CompressPackage Behavior](#413-compresspackage-behavior)
    - [4.1.4 CompressPackage Error Conditions](#414-compresspackage-error-conditions)
  - [4.2 DecompressPackage](#42-packagedecompresspackage-method)
    - [4.2.1 DecompressPackage Purpose](#421-decompresspackage-purpose)
    - [4.2.2 DecompressPackage Parameters](#422-decompresspackage-parameters)
    - [4.2.3 DecompressPackage Behavior](#423-decompresspackage-behavior)
    - [4.2.4 DecompressPackage Error Conditions](#424-decompresspackage-error-conditions)
- [5. Streaming Compression Methods](#5-streaming-compression-methods)
  - [5.1 CompressPackageStream](#51-packagecompresspackagestream-method)
    - [5.1.1 CompressPackageStream Purpose](#511-compresspackagestream-purpose)
    - [5.1.2 CompressPackageStream Parameters](#512-compresspackagestream-parameters)
    - [5.1.3 CompressPackageStream Behavior](#513-compresspackagestream-behavior)
    - [5.1.4 CompressPackageStream Error Conditions](#514-compresspackagestream-error-conditions)
    - [5.1.5 Configuration Usage Patterns](#515-configuration-usage-patterns)
  - [5.2 DecompressPackageStream](#52-packagedecompresspackagestream-method)
    - [5.2.1 DecompressPackageStream Purpose](#521-decompresspackagestream-purpose)
    - [5.2.2 DecompressPackageStream Parameters](#522-decompresspackagestream-parameters)
    - [5.2.3 DecompressPackageStream Behavior](#523-decompresspackagestream-behavior)
    - [5.2.4 DecompressPackageStream Error Conditions](#524-decompresspackagestream-error-conditions)
- [6. File-Based Compression Methods](#6-file-based-compression-methods)
  - [6.1 CompressPackageFile](#61-packagecompresspackagefile-method)
    - [6.1.1 CompressPackageFile Purpose](#611-compresspackagefile-purpose)
    - [6.1.2 CompressPackageFile Parameters](#612-compresspackagefile-parameters)
    - [6.1.3 CompressPackageFile Behavior](#613-compresspackagefile-behavior)
    - [6.1.4 CompressPackageFile Error Conditions](#614-compresspackagefile-error-conditions)
  - [6.2 DecompressPackageFile](#62-packagedecompresspackagefile-method)
    - [6.2.1 DecompressPackageFile Purpose](#621-decompresspackagefile-purpose)
    - [6.2.2 DecompressPackageFile Parameters](#622-decompresspackagefile-parameters)
    - [6.2.3 DecompressPackageFile Behavior](#623-decompresspackagefile-behavior)
    - [6.2.4 DecompressPackageFile Error Conditions](#624-decompresspackagefile-error-conditions)
- [7. Compression Information and Status](#7-compression-information-and-status)
  - [7.1 Compression Information Structure Reference](#71-compression-information-structure-reference)
  - [7.2 Compression Status Methods](#72-compression-status-methods)
    - [7.2.1 Package Compression Query Methods](#721-package-compression-query-methods)
    - [7.2.2 Package GetPackageCompressionInfo Method](#722-packagegetpackagecompressioninfo-method)
    - [7.2.3 Package IsPackageCompressed Method](#723-packageispackagecompressed-method)
    - [7.2.4 Package GetPackageCompressionType Method](#724-packagegetpackagecompressiontype-method)
    - [7.2.5 Package GetPackageCompressionRatio Method](#725-packagegetpackagecompressionratio-method)
    - [7.2.6 Package GetPackageOriginalSize Method](#726-packagegetpackageoriginalsize-method)
    - [7.2.7 Package GetPackageCompressedSize Method](#727-packagegetpackagecompressedsize-method)
    - [7.2.8 Package Compression Control Methods](#728-package-compression-control-methods)
    - [7.2.9 Metadata Index Methods](#729-metadata-index-methods)
    - [7.2.10 Generic Compression Methods](#7210-generic-compression-methods)
  - [7.3 Internal Compression Methods](#73-internal-compression-methods)
    - [7.3.1 Package compressPackageContent Method](#731-packagecompresspackagecontent-method)
    - [7.3.2 Package decompressPackageContent Method](#732-packagedecompresspackagecontent-method)
- [8. Concurrency Patterns and Thread Safety](#8-concurrency-patterns-and-thread-safety)
  - [8.1 Thread Safety Guarantees](#81-thread-safety-guarantees)
    - [8.1.1 Thread Safety None Mode](#811-thread-safety-none-mode)
    - [8.1.2 Thread Safety Read-Only Mode](#812-thread-safety-read-only-mode)
    - [8.1.3 Thread Safety Concurrent Mode](#813-thread-safety-concurrent-mode)
    - [8.1.4 Thread Safety Full Mode](#814-thread-safety-full-mode)
  - [8.2 Worker Pool Management](#82-worker-pool-management)
    - [8.2.1 CompressionWorkerPool Structure](#821-compressionworkerpool-structure)
    - [8.2.2 CompressionWorkerPool CompressConcurrently Method](#822-compressionworkerpooltcompressconcurrently-method)
    - [8.2.3 CompressionWorkerPool DecompressConcurrently Method](#823-compressionworkerpooltdecompressconcurrently-method)
    - [8.2.4 CompressionWorkerPool GetCompressionStats Method](#824-compressionworkerpooltgetcompressionstats-method)
  - [8.3 Concurrent Compression Methods](#83-concurrent-compression-methods)
    - [8.3.1 Package CompressPackageConcurrent Method](#831-packagecompresspackageconcurrent-method)
    - [8.3.2 Package DecompressPackageConcurrent Method](#832-packagedecompresspackageconcurrent-method)
  - [8.4 Resource Management](#84-resource-management)
    - [8.4.1 CompressionResourcePool Structure](#841-compressionresourcepool-structure)
    - [8.4.2 CompressionResource Structure](#842-compressionresource-structure)
- [9. Compression Configuration Patterns](#9-compression-configuration-patterns)
  - [9.1 Compression-Specific Configuration](#91-compression-specific-configuration)
    - [9.1.1 CompressionConfig Structure](#911-compressionconfig-structure)
    - [9.1.2 CompressionConfigBuilder Structure](#912-compressionconfigbuilder-structure)
  - [9.2 Compression Validation Patterns](#92-compression-validation-patterns)
    - [9.2.1 CompressionValidator Structure](#921-compressionvalidator-structure)
    - [9.2.2 CompressionValidationRule Structure](#922-compressionvalidationrule-structure)
- [10. Compression and Signing Relationship](#10-compression-and-signing-relationship)
  - [10.1 Signing Compressed Packages](#101-signing-compressed-packages)
    - [10.1.1 Supported Operation](#1011-supported-operation)
    - [10.1.2 Signing Compressed Packages Process](#1012-signing-compressed-packages-process)
    - [10.1.3 Signing Compressed Packages Benefits](#1013-signing-compressed-packages-benefits)
  - [10.2 Compressing Signed Packages](#102-compressing-signed-packages)
    - [10.2.1 Not Supported](#1021-not-supported)
    - [10.2.2 Compressing Signed Packages Reasoning](#1022-compressing-signed-packages-reasoning)
    - [10.2.3 Compressing Signed Packages Workflow](#1023-compressing-signed-packages-workflow)
- [11. CompressionStrategy Selection](#11-compressionstrategy-selection)
  - [11.1 Compression Type Selection](#111-compression-type-selection)
    - [11.1.1 Zstandard Compression (Type 1)](#1111-zstandard-compression-type-1)
    - [11.1.2 LZ4 Compression (Type 2)](#1112-lz4-compression-type-2)
    - [11.1.3 LZMA Compression (Type 3)](#1113-lzma-compression-type-3)
  - [11.2 Compression Decision Matrix](#112-compression-decision-matrix)
    - [11.2.1 User Guidance Matrix](#1121-user-guidance-matrix)
    - [11.2.2 Automatic Compression Type Selection](#1122-automatic-compression-type-selection)
  - [11.3 Compression Workflow Options](#113-compression-workflow-options)
    - [11.3.1 Option 1: Compress Before Writing](#1131-option-1-compress-before-writing)
    - [11.3.2 Option 2: Compress and Write in One Step](#1132-option-2-compress-and-write-in-one-step)
    - [11.3.3 Option 3: Write with Compression](#1133-option-3-write-with-compression)
    - [11.3.4 Option 4: Stream Compression for Large Packages](#1134-option-4-stream-compression-for-large-packages)
    - [11.3.5 Option 5: Advanced Streaming Compression](#1135-option-5-advanced-streaming-compression)
    - [11.3.6 Option 6: Custom Memory Management](#1136-option-6-custom-memory-management)
- [12. Error Handling](#12-error-handling)
  - [12.1 Common Error Conditions](#121-common-error-conditions)
    - [12.1.1 Common Error Conditions Compression Errors](#1211-common-error-conditions-compression-errors)
    - [12.1.2 Common Error Conditions Decompression Errors](#1212-common-error-conditions-decompression-errors)
    - [12.1.3 Common Error Conditions File Operation Errors](#1213-common-error-conditions-file-operation-errors)
  - [12.2 Error Recovery](#122-error-recovery)
    - [12.2.1 Error Recovery Compression Failure](#1221-error-recovery-compression-failure)
    - [12.2.2 Error Recovery Decompression Failure](#1222-error-recovery-decompression-failure)
- [13. Modern Best Practices](#13-modern-best-practices)
  - [13.1 Industry Standards Alignment](#131-industry-standards-alignment)
    - [13.1.1 Streaming Compression (Universal Standard)](#1311-streaming-compression-universal-standard)
    - [13.1.2 Parallel Processing (Performance Critical)](#1312-parallel-processing-performance-critical)
    - [13.1.3 Chunked Processing (Industry Standard)](#1313-chunked-processing-industry-standard)
    - [13.1.4 Memory Management (Resource Critical)](#1314-memory-management-resource-critical)
  - [13.2 Intelligent Defaults and Memory Management](#132-intelligent-defaults-and-memory-management)
    - [13.2.1 MemoryStrategy Defaults](#1321-memorystrategy-defaults)
    - [13.2.2 Auto-Detection Logic](#1322-auto-detection-logic)
    - [13.2.3 Adaptive Memory Management](#1323-adaptive-memory-management)
  - [13.3 Performance Considerations (CompressionOperations)](#133-performance-considerations-compressionoperations)
  - [13.4 Memory Usage](#134-memory-usage)
    - [13.4.1 Compression Memory Usage](#1341-compression-memory-usage)
    - [13.4.2 Decompression Memory Usage](#1342-decompression-memory-usage)
  - [13.5 CPU Usage](#135-cpu-usage)
    - [13.5.1 CPU Usage Compression](#1351-cpu-usage-compression)
    - [13.5.2 CPU Usage Decompression](#1352-cpu-usage-decompression)
  - [13.6 I/O Considerations](#136-io-considerations)
    - [13.6.1 I/O Considerations File-Based Operations](#1361-io-considerations-file-based-operations)
    - [13.6.2 ~~I/O Considerations Network Operations~~ (REMOVED)](#1362-io-considerations-network-operations-removed)
- [14. Structured Error System](#14-structured-error-system)
  - [14.1 Structured Error System (Compression API)](#141-structured-error-system-compression-api)
  - [14.2 Common Compression Error Types](#142-common-compression-error-types)
    - [14.2.1 Compression Error Types](#1421-compression-error-types)
  - [14.3 Structured Error Examples](#143-structured-error-examples)
    - [14.3.1 Creating Compression Errors](#1431-creating-compression-errors)
    - [14.3.2 Error Handling Patterns](#1432-error-handling-patterns)

## 0. Overview

This document defines the NovusPack package compression API, providing methods for compressing and decompressing package content while maintaining package integrity and signature compatibility.

### 0.1 Cross-References

- [Core Package Interface API](api_core.md) - Package operations and compression
- [Package Writing API](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation
- [File Compression API](api_file_mgmt_compression.md) - Individual file compression operations (FileEntry.Compress, Package.CompressFile, etc.)
- [Security and Encryption](security.md) - Comprehensive security architecture and encryption implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Streaming and Buffer Management](api_streaming.md) - Streaming concurrency patterns and buffer management

### 0.2 Context Integration

All public methods in the NovusPack Compression API accept `context.Context` as the first parameter to support:

- Request cancellation and timeout handling
- Request-scoped values and configuration
- Graceful shutdown and resource cleanup
- Integration with Go's standard context patterns

This follows 2025 Go best practices and ensures the API is compatible with modern Go applications and frameworks.

## 1. Package Compression Overview

Package compression in NovusPack compresses package content using separate compression for metadata and data blocks, while preserving the header, metadata index, package comment, and signatures in an uncompressed state for direct access.
This enables selective decompression of metadata without requiring full package decompression.

### 1.1 Compression Scope

This section describes the scope of compression operations.

#### 1.1.1 Compressed Content

When package compression is enabled (header flags bits 15-8 != 0), the following content is compressed:

- **FileEntry metadata**: Each FileEntry (64 bytes + variable data) is compressed individually using LZ4 compression
- **File data**: Each file's data is compressed individually using the package compression type (Zstd, LZ4, or LZMA)
- **File index**: The regular file index is compressed as a single block using LZ4 compression

**Note**: Package compression compresses all file data as part of the package structure.
For compressing individual files within a package (without compressing the entire package), see [File Compression API](api_file_mgmt_compression.md).

Special metadata files (types 65000-65535) are handled as regular FileEntry objects:

- FileEntry metadata compressed with LZ4 (same as all FileEntry metadata)
- File data (YAML content) compressed with LZ4 for fast access
  - Note that this is only a requirement for fully-compressed packages; special metadata files can also be stored uncompressed when NOT implementing full package compression

#### 1.1.2 Uncompressed Content

The following content remains uncompressed for direct access:

- Package header (see [Package File Format - Package Header](package_file_format.md#2-package-header))
- Metadata index (see [Package File Format - Metadata Index Section](package_file_format.md#5-metadata-index-section)) - enables fast access to compressed blocks
- Package comment
- Digital signatures

### 1.2 Compression Types

```go
const (
    CompressionNone = 0  // No compression
    CompressionZstd = 1  // Zstd compression
    CompressionLZ4  = 2  // LZ4 compression
    CompressionLZMA = 3  // LZMA compression
)
```

### 1.3 PackageCompressionInfo Struct

```go
// PackageCompressionInfo contains package compression details
type PackageCompressionInfo struct {
    Type         uint8   // Compression type (0-3)
    IsCompressed bool    // Whether package is compressed
    OriginalSize int64   // Original package size before compression
    CompressedSize int64 // Compressed package size
    Ratio        float64 // Compression ratio (0.0-1.0)
}
```

### 1.4 Compression Constraints

- **Signed Package Restriction**: Packages with signatures cannot be compressed
- **Compression Before Signing**: Packages must be compressed before signing
- **Header Immutability**: Once compressed, the header becomes immutable
- **Metadata Index Location**: Metadata index is located at fixed offset 112 bytes (immediately after header) when compression is enabled

## 2. Strategy Pattern Interfaces

The compression API supports pluggable compression algorithms through the strategy pattern.

### 2.1 CompressionStrategy Interface

This section describes the CompressionStrategy interface.

#### 2.1.1 CompressionStrategy Interface Definition

```go
// CompressionStrategy extends Strategy[T, T] for compression operations
// Both input and output are the same type T
// The Strategy.Type() method returns "compression" as the category
type CompressionStrategy[T any] interface {
    Strategy[T, T]  // Extends the generic Strategy interface

    Compress(ctx context.Context, data T) (T, error)
    Decompress(ctx context.Context, data T) (T, error)
    CompressionType() CompressionType  // Returns the specific compression algorithm type
    Name() string
}
```

#### 2.1.2 ByteCompressionStrategy Interface

```go
// ByteCompressionStrategy is the concrete implementation for []byte data
type ByteCompressionStrategy interface {
    CompressionStrategy[[]byte]
}
```

#### 2.1.3 AdvancedCompressionStrategy Interface

```go
// AdvancedCompressionStrategy for compression with additional validation and metrics
type AdvancedCompressionStrategy[T any] interface {
    CompressionStrategy[T]
    ValidateInput(ctx context.Context, data T) error
    GetCompressionRatio(ctx context.Context, original T, compressed T) float64
}
```

#### 2.1.4 StreamConfig Structure

```go
// StreamConfig handles streaming compression for files of any size
type StreamConfig struct {
    // Basic settings
    ChunkSize        int64  // Size of processing chunks (0 = auto-calculate)
    TempDir         string // Directory for temporary files ("" = system temp)
    MaxMemoryUsage  int64  // Maximum memory usage (0 = auto-detect, -1 = no limit)
    UseDiskBuffering bool  // Use disk for intermediate buffering
    CleanupTempFiles bool  // Clean up temporary files after completion
    ProgressCallback func(bytesProcessed int64, totalBytes int64) // Progress reporting

    // Advanced settings (optional - use nil for defaults)
    UseParallelProcessing bool  // Enable multi-core processing (default: true)
    MaxWorkers          int   // Maximum parallel workers (0 = auto-detect)
    CompressionLevel    int   // Compression level (0 = auto-select, 1-22 for zstd, 1-9 for others)
    UseSolidCompression bool  // Use solid compression for better ratios
    ResumeFromOffset    int64 // Resume from specific offset (0 = start)

    // Memory management enhancements
    MemoryStrategy      MemoryStrategy // Memory management strategy
    AdaptiveChunking    bool           // Enable adaptive chunk sizing based on memory
    BufferPoolSize      int            // Buffer pool size (0 = auto-calculate)
    MaxTempFileSize     int64          // Maximum temp file size before rotation (0 = no limit)

    // Concurrency and thread safety (see api_generics.md for ConcurrencyConfig and ThreadSafetyMode)
    ConcurrencyConfig   *ConcurrencyConfig // Thread safety and worker management
    ThreadSafetyMode    ThreadSafetyMode   // Thread safety guarantees
}
```

#### 2.1.5 MemoryStrategy Type

```go
// MemoryStrategy defines memory management approach
type MemoryStrategy int

const (
    MemoryStrategyConservative MemoryStrategy = iota // Use 25% of available RAM
    MemoryStrategyBalanced                          // Use 50% of available RAM (default)
    MemoryStrategyAggressive                        // Use 75% of available RAM
    MemoryStrategyCustom                            // Use MaxMemoryUsage value
)
```

### 2.2 Built-in Compression Strategies

This section describes built-in compression strategy implementations.

#### 2.2.1 ZstandardStrategy Structure

```go
// Zstandard compression strategy with generic support
type ZstandardStrategy[T any] struct {
    level    int
    strategy CompressionStrategy[T]
}
```

#### 2.2.2 LZ4Strategy Structure

```go
// LZ4 compression strategy with generic support
type LZ4Strategy[T any] struct {
    level    int
    strategy CompressionStrategy[T]
}
```

#### 2.2.3 LZMAStrategy Structure

```go
// LZMA compression strategy with generic support
type LZMAStrategy[T any] struct {
    level    int
    strategy CompressionStrategy[T]
}
```

#### 2.2.4 CompressionJob Structure

```go
// CompressionJob represents a unit of work for compression (extends Job)
type CompressionJob[T any] struct {
    *Job[T]
    CompressionType uint8
    CompressionLevel int
}
```

## 3. Interface Granularity and Composition

The compression API uses focused interfaces to provide clear separation of concerns and enable flexible composition.

### 3.1 CompressionInfo Interface

```go
// CompressionInfo provides read-only access to compression information
type CompressionInfo interface {
    GetCompressionInfo(ctx context.Context) PackageCompressionInfo
    IsCompressed() bool
    GetCompressionType() (uint8, error) // Returns compression type, or error if not compressed
    GetCompressionRatio() (float64, error) // Returns compression ratio, or error if not compressed
    CanCompress() bool
}
```

### 3.2 CompressionOperations Interface

```go
// CompressionOperations provides basic compression/decompression operations
type CompressionOperations interface {
    CompressPackage(ctx context.Context, compressionType uint8) error
    DecompressPackage(ctx context.Context) error
    SetCompressionType(ctx context.Context, compressionType uint8) error
}
```

### 3.3 CompressionStreaming Interface

```go
// CompressionStreaming provides streaming compression for large packages
type CompressionStreaming interface {
    CompressPackageStream(ctx context.Context, compressionType uint8, config *StreamConfig) error
    DecompressPackageStream(ctx context.Context, config *StreamConfig) error
}
```

### 3.4 CompressionFileOperations Interface

```go
// CompressionFileOperations provides file-based compression operations
type CompressionFileOperations interface {
    CompressPackageFile(ctx context.Context, path string, compressionType uint8, overwrite bool) error
    DecompressPackageFile(ctx context.Context, path string, overwrite bool) error
}
```

**Note**: These methods compress or decompress the entire package structure.
For compressing individual files within a package (without compressing the entire package), see [File Compression API](api_file_mgmt_compression.md).

### 3.5 Generic Compression Interface

The `CompressionStrategy[T]` interface extends the generic [Core Generic Types](api_generics.md#1-core-generic-types) pattern for compression-specific operations.
`CompressionStrategy[T]` embeds `Strategy[T, T]` where both input and output are the same type.
The `Process` method from `Strategy[T, T]` can be used for compression operations, while `Compress` and `Decompress` provide more specific compression/decompression methods.

```go
// Compression provides type-safe compression for any data type
type Compression[T any] interface {
    CompressGeneric(ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
    DecompressGeneric(ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
    ValidateCompressionData(ctx context.Context, data T) error
}
```

**Cross-Reference**: For the base strategy pattern, see [Core Generic Types](api_generics.md#1-core-generic-types).

## 4. In-Memory Compression Methods

These methods operate on packages in memory without writing to disk.

**Note**: For large packages, consider using [Streaming Compression Methods](#5-streaming-compression-methods) to avoid memory limitations.

### 4.1 Package.CompressPackage Method

```go
// CompressPackage compresses package content in memory
// Compresses file entries and data separately using LZ4 for metadata and specified type for data
// Compresses file index with LZ4
// Creates metadata index for fast access (NOT header, metadata index, comment, or signatures)
// Signed packages cannot be compressed
// Returns *PackageError on failure
func (p *Package) CompressPackage(ctx context.Context, compressionType uint8) error
```

#### 4.1.1 CompressPackage Purpose

Handle compression/decompression of in-memory packages with separate metadata and data compression.

#### 4.1.2 CompressPackage Parameters

- `ctx`: Context for cancellation and timeout handling
- `compressionType`: Compression algorithm to use for file data (1-3), metadata always uses LZ4

#### 4.1.3 CompressPackage Behavior

- Compresses FileEntry metadata individually using LZ4
- Compresses file data individually using specified compression type (Zstd, LZ4, or LZMA)
- Compresses special metadata files (types 65000-65535) with LZ4 for fast access
- Compresses file index with LZ4 as a single block
- Creates metadata index for fast access to compressed blocks
- Updates package compression state in memory
- Returns error if package is signed
- Updates package header compression flags (bits 15-8)
- Writes metadata index at fixed offset 112 bytes (immediately after header)

#### 4.1.4 CompressPackage Error Conditions

- Package is already signed (cannot compress signed packages)
- Invalid compression type (must be 1-3)
- Package is already compressed with different type
- Context cancellation
- Metadata index creation failure

### 4.2 Package.DecompressPackage Method

```go
// DecompressPackage decompresses the package in memory
// Decompresses all compressed content
// Returns *PackageError on failure
func (p *Package) DecompressPackage(ctx context.Context) error
```

#### 4.2.1 DecompressPackage Purpose

Decompress package content in memory

#### 4.2.2 DecompressPackage Parameters

- `ctx`: Context for cancellation and timeout handling

#### 4.2.3 DecompressPackage Behavior

- Decompresses all compressed content (metadata blocks, data blocks, and file index)
- Updates package compression state in memory
- Clears package header compression flags (bits 15-8)
- Removes metadata index (no longer needed when uncompressed)
- Preserves all other package data

#### 4.2.4 DecompressPackage Error Conditions

- Package is not compressed
- Decompression failure
- Context cancellation

## 5. Streaming Compression Methods

These methods handle compression/decompression of large packages using streaming to avoid memory limitations.

**For Large Files**: These methods use temporary files and chunked processing to handle files that exceed available RAM, with adaptive strategies based on configuration.

### 5.1 Package.CompressPackageStream Method

```go
// CompressPackageStream compresses large package content using streaming
// Uses temporary files and chunked processing to handle files of any size
// Configuration determines the level of optimization and memory management
// Returns *PackageError on failure
func (p *Package) CompressPackageStream(ctx context.Context, compressionType uint8, config *StreamConfig) error
```

#### 5.1.1 CompressPackageStream Purpose

Handle compression of large packages using streaming, temporary files, and configurable optimization strategies for files of any size

#### 5.1.2 CompressPackageStream Parameters

- `ctx`: Context for cancellation and timeout handling
- `compressionType`: Compression algorithm to use (1-3)
- `config`: Unified streaming configuration for memory management and optimization

#### 5.1.3 CompressPackageStream Behavior

- Uses streaming for large package content
- Creates temporary files when needed for memory management
- Compresses file entries + data + index (NOT header, comment, or signatures)
- Returns error if package is signed
- Updates package header compression flags
- **Adaptive Processing**: Automatically adjusts strategy based on file size and configuration
- **Memory Management**: Respects `StreamConfig.MaxMemoryUsage` to prevent OOM
- **Progress Reporting**: Provides progress updates for long-running operations
- **Parallel Processing**: Uses multiple CPU cores when enabled in configuration
- **Chunked Processing**: Processes files in configurable chunks (auto-calculated if 0)

#### 5.1.4 CompressPackageStream Error Conditions

- **Security Errors**: Package is already signed (cannot compress signed packages)
- **Validation Errors**: Invalid compression type, invalid stream configuration
- **I/O Errors**: Temporary file creation failed, insufficient disk space
- **Context Errors**: Context cancellation or timeout exceeded

#### 5.1.5 Configuration Usage Patterns

The unified `StreamConfig` supports different usage patterns based on requirements:

##### 5.1.5.1 Simple Usage

**Simple Usage** (basic settings only):

```go
config := &StreamConfig{
    ChunkSize: 0,        // Auto-calculate
    MaxMemoryUsage: 0,   // Auto-detect
    TempDir: "",         // System temp
}
```

##### 5.1.5.2 Advanced Usage

**Advanced Usage** (full configuration):

```go
config := &StreamConfig{
    ChunkSize: 1024 * 1024 * 1024,  // 1GB chunks
    MaxMemoryUsage: 8 * 1024 * 1024 * 1024,  // 8GB limit
    UseParallelProcessing: true,
    MaxWorkers: 0,  // Auto-detect
    CompressionLevel: 0,  // Auto-select
    UseSolidCompression: true,
    MemoryStrategy: MemoryStrategyBalanced,
    AdaptiveChunking: true,
}
```

### 5.2 Package.DecompressPackageStream Method

```go
// DecompressPackageStream decompresses large package content using streaming
// Uses streaming to manage memory efficiently for large packages
// Returns *PackageError on failure
func (p *Package) DecompressPackageStream(ctx context.Context, config *StreamConfig) error
```

#### 5.2.1 DecompressPackageStream Purpose

Decompress large package content using streaming

#### 5.2.2 DecompressPackageStream Parameters

- `ctx`: Context for cancellation and timeout handling
- `config`: Streaming configuration for memory management

#### 5.2.3 DecompressPackageStream Behavior

- Uses streaming for large package content
- Decompresses all compressed content
- Updates package compression state in memory
- Clears package header compression flags
- Preserves all other package data

#### 5.2.4 DecompressPackageStream Error Conditions

- **Validation Errors**: Package is not compressed, invalid stream configuration
- **Compression Errors**: Decompression operation failed, algorithm-specific failures
- **I/O Errors**: Streaming operation failed, insufficient disk space
- **Context Errors**: Context cancellation or timeout exceeded

## 6. File-Based Compression Methods

These methods handle both compression/decompression and writing to a file.

**Note**: These methods compress or decompress the entire package structure.
For compressing individual files within a package (FileEntry.Compress, Package.CompressFile), see [File Compression API](api_file_mgmt_compression.md).

### 6.1 Package.CompressPackageFile Method

```go
// CompressPackageFile compresses package content and writes to specified path
// Compresses file entries + data + index (NOT header, comment, or signatures)
// Signed packages cannot be compressed
// Returns *PackageError on failure
func (p *Package) CompressPackageFile(ctx context.Context, path string, compressionType uint8, overwrite bool) error
```

#### 6.1.1 CompressPackageFile Purpose

Handle compression AND write to file

#### 6.1.2 CompressPackageFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Target file path for compressed package
- `compressionType`: Compression algorithm to use (1-3)
- `overwrite`: Whether to overwrite existing file

#### 6.1.3 CompressPackageFile Behavior

- Compresses package content in memory
- Writes compressed package to specified path
- Creates new file by default, overwrites if `overwrite=true`
- Compresses file entries + data + index (NOT header, comment, or signatures)
- Returns error if package is signed

#### 6.1.4 CompressPackageFile Error Conditions

- Package is already signed
- Invalid compression type
- File already exists and `overwrite=false`
- I/O errors
- Context cancellation

### 6.2 Package.DecompressPackageFile Method

```go
// DecompressPackageFile decompresses the package and writes to specified path
// Decompresses all compressed content and writes uncompressed package
// Returns *PackageError on failure
func (p *Package) DecompressPackageFile(ctx context.Context, path string, overwrite bool) error
```

#### 6.2.1 DecompressPackageFile Purpose

Decompress package and write to file

#### 6.2.2 DecompressPackageFile Parameters

- `ctx`: Context for cancellation and timeout handling
- `path`: Target file path for uncompressed package
- `overwrite`: Whether to overwrite existing file

#### 6.2.3 DecompressPackageFile Behavior

- Decompresses package content in memory
- Writes uncompressed package to specified path
- Creates new file by default, overwrites if `overwrite=true`
- Decompresses all compressed content

#### 6.2.4 DecompressPackageFile Error Conditions

- Package is not compressed
- File already exists and `overwrite=false`
- I/O errors
- Context cancellation

## 7. Compression Information and Status

This section describes compression information and status operations.

### 7.1 Compression Information Structure Reference

See [1.3 PackageCompressionInfo Struct](#13-packagecompressioninfo-struct) for the complete structure definition.

### 7.2 Compression Status Methods

This section describes compression status methods.

#### 7.2.1 Package Compression Query Methods

This section describes package compression query methods.

#### 7.2.2 Package.GetPackageCompressionInfo Method

```go
// GetPackageCompressionInfo returns package compression information
func (p *Package) GetPackageCompressionInfo() PackageCompressionInfo
```

#### 7.2.3 Package.IsPackageCompressed Method

```go
// IsPackageCompressed checks if the package is compressed
// Checks header flags bits 15-8 for compression type
func (p *Package) IsPackageCompressed() bool
```

#### 7.2.4 Package.GetPackageCompressionType Method

```go
// GetPackageCompressionType returns the package compression type
// Returns compression type from header flags bits 15-8
// Returns *PackageError if package is not compressed
func (p *Package) GetPackageCompressionType() (uint8, error)
```

#### 7.2.5 Package.GetPackageCompressionRatio Method

```go
// GetPackageCompressionRatio returns the compression ratio
// Returns *PackageError if package is not compressed
func (p *Package) GetPackageCompressionRatio() (float64, error)
```

#### 7.2.6 Package.GetPackageOriginalSize Method

```go
// GetPackageOriginalSize returns the original size before compression
// Returns *PackageError if package is not compressed
func (p *Package) GetPackageOriginalSize() (int64, error)
```

#### 7.2.7 Package.GetPackageCompressedSize Method

```go
// GetPackageCompressedSize returns the compressed size
// Returns *PackageError if package is not compressed
func (p *Package) GetPackageCompressedSize() (int64, error)
```

#### 7.2.8 Package Compression Control Methods

This section describes package compression control methods.

##### 7.2.8.1 Package.SetPackageCompressionType Method

```go
// SetPackageCompressionType sets the package compression type (without compressing)
// Returns *PackageError on failure
func (p *Package) SetPackageCompressionType(compressionType uint8) error
```

##### 7.2.8.2 Package.CanCompressPackage Method

```go
// CanCompressPackage checks if package can be compressed (not signed)
func (p *Package) CanCompressPackage() bool
```

#### 7.2.9 Metadata Index Methods

This section describes metadata index methods.

##### 7.2.9.1 Package.HasMetadataIndex Method

```go
// HasMetadataIndex checks if package has metadata index (compression enabled)
// Returns true if header flags bits 15-8 != 0
func (p *Package) HasMetadataIndex() bool
```

##### 7.2.9.2 Package.GetMetadataIndexOffset Method

```go
// GetMetadataIndexOffset returns the offset to metadata index
// Returns fixed offset 112 bytes (PackageHeaderSize) when compression enabled
// Returns *PackageError if package is not compressed (no metadata index)
func (p *Package) GetMetadataIndexOffset() (int64, error)
```

#### 7.2.10 Generic Compression Methods

This section describes generic compression methods.

##### 7.2.10.1 Package.CompressGeneric Method

```go
// Generic compression methods for type-safe operations
// CompressionStrategy[T] embeds Strategy[T, T] from the generics package
// See [Core Generic Types](api_generics.md#1-core-generic-types) for base strategy pattern
func (p *Package) CompressGeneric[T any](ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
```

##### 7.2.10.2 Package.DecompressGeneric Method

```go
// DecompressGeneric decompresses data using a generic compression strategy.
func (p *Package) DecompressGeneric[T any](ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
```

##### 7.2.10.3 Package.ValidateCompressionData Method

```go
// Returns *PackageError on failure
func (p *Package) ValidateCompressionData[T any](ctx context.Context, data T) error
```

**Type Constraints**: The type parameter `T` in `CompressGeneric` and `DecompressGeneric` is typically `[]byte` for compression operations, but can be any type that the `CompressionStrategy[T]` supports.
For most use cases, `T` should be `[]byte` to work with data directly.
The constraint `any` is used because compression strategies may work with different data types (e.g., `[]byte`, custom serializable types).

**Error Handling**: All compression operations return errors using `NewPackageError` or `WrapErrorWithContext` with typed error context for type-safe error handling.
See [Error Handling](#12-error-handling) for details.

### 7.3 Internal Compression Methods

This section describes internal compression methods.

#### 7.3.1 Package.compressPackageContent Method

```go
// Internal compression methods (used by CompressPackage and Write)
// Returns *PackageError on failure
func (p *Package) compressPackageContent(ctx context.Context, compressionType uint8) error
```

#### 7.3.2 Package.decompressPackageContent Method

```go
// Returns *PackageError on failure
func (p *Package) decompressPackageContent(ctx context.Context) error
```

## 8. Concurrency Patterns and Thread Safety

The compression API provides explicit concurrency patterns and thread safety guarantees for safe concurrent usage.

### 8.1 Thread Safety Guarantees

The compression API provides different levels of thread safety based on the `ThreadSafetyMode` configuration:

#### 8.1.1 Thread Safety None Mode

No thread safety guarantees.
Operations should not be called concurrently.

#### 8.1.2 Thread Safety Read-Only Mode

Read-only operations are safe for concurrent access.
Multiple goroutines can safely call read methods simultaneously.

#### 8.1.3 Thread Safety Concurrent Mode

Concurrent read/write operations are supported.
Uses read-write mutex for optimal read performance.

#### 8.1.4 Thread Safety Full Mode

Full thread safety with complete synchronization.
All operations are protected by appropriate locking mechanisms.

### 8.2 Worker Pool Management

The compression API uses the generic worker pool patterns defined in [api_generics.md](api_generics.md#2-generic-function-patterns) with compression-specific extensions.

#### 8.2.1 CompressionWorkerPool Structure

```go
// CompressionWorkerPool extends WorkerPool for compression operations
type CompressionWorkerPool[T any] struct {
    *WorkerPool[T]
    compressionStrategy CompressionStrategy[T]
}
```

#### 8.2.2 CompressionWorkerPool[T].CompressConcurrently Method

```go
// Compression-specific methods
func (p *CompressionWorkerPool[T]) CompressConcurrently(ctx context.Context, data []T, strategy CompressionStrategy[T]) ([]T, error)
```

#### 8.2.3 CompressionWorkerPool[T].DecompressConcurrently Method

```go
// DecompressConcurrently decompresses multiple data items concurrently using a worker pool.
func (p *CompressionWorkerPool[T]) DecompressConcurrently(ctx context.Context, data []T, strategy CompressionStrategy[T]) ([]T, error)
```

#### 8.2.4 CompressionWorkerPool[T].GetCompressionStats Method

```go
// GetCompressionStats returns statistics about compression operations performed by the worker pool.
func (p *CompressionWorkerPool[T]) GetCompressionStats() CompressionStats
```

### 8.3 Concurrent Compression Methods

This section describes concurrent compression methods.

#### 8.3.1 Package.CompressPackageConcurrent Method

```go
// CompressPackageConcurrent compresses package content using worker pool
// Returns *PackageError on failure
func (p *Package) CompressPackageConcurrent(ctx context.Context, compressionType uint8, config *StreamConfig) error
```

#### 8.3.2 Package.DecompressPackageConcurrent Method

```go
// DecompressPackageConcurrent decompresses package content using worker pool
// Returns *PackageError on failure
func (p *Package) DecompressPackageConcurrent(ctx context.Context, config *StreamConfig) error
```

### 8.4 Resource Management

The compression API uses the generic resource management patterns defined in [api_generics.md](api_generics.md#2-generic-function-patterns) with compression-specific resource types.

#### 8.4.1 CompressionResourcePool Structure

This section describes the CompressionResourcePool structure.

##### 8.4.1.1 CompressionResourcePool Struct

```go
// CompressionResourcePool manages compression-specific resources
type CompressionResourcePool struct {
    *ResourcePool[CompressionResource]
    compressionConfig *CompressionConfig
}
```

##### 8.4.1.2 CompressionResourcePool.AcquireCompressionResource Method

```go
// Compression-specific resource management methods
func (p *CompressionResourcePool) AcquireCompressionResource(ctx context.Context, strategyType uint8) (*CompressionResource, error)
```

##### 8.4.1.3 CompressionResourcePool.ReleaseCompressionResource Method

```go
// Returns *PackageError on failure
func (p *CompressionResourcePool) ReleaseCompressionResource(resource *CompressionResource) error
```

##### 8.4.1.4 CompressionResourcePool.GetCompressionResourceStats Method

```go
// GetCompressionResourceStats returns statistics about compression resource usage.
func (p *CompressionResourcePool) GetCompressionResourceStats() CompressionResourceStats
```

#### 8.4.2 CompressionResource Structure

```go
// CompressionResource represents a compression-specific resource
type CompressionResource struct {
    ID           string
    Strategy     CompressionStrategy[[]byte]
    Buffer       []byte
    LastUsed     time.Time
    AccessCount  int64
}
```

## 9. Compression Configuration Patterns

The compression API provides compression-specific configuration patterns that extend the generic configuration patterns defined in [api_generics.md](api_generics.md#110-generic-configuration-patterns).

### 9.1 Compression-Specific Configuration

This section describes compression-specific configuration options.

#### 9.1.1 CompressionConfig Structure

```go
// CompressionConfig extends Config for compression-specific settings
type CompressionConfig struct {
    *Config[[]byte]

    // Compression-specific settings
    CompressionType     Option[uint8]           // Compression algorithm type
    CompressionLevel    Option[int]             // Compression level (1-22 for zstd, 1-9 for others)
    UseSolidCompression Option[bool]            // Use solid compression for better ratios
    ResumeFromOffset    Option[int64]           // Resume from specific offset
    MemoryStrategy      Option[MemoryStrategy]  // Memory management strategy
}
```

#### 9.1.2 CompressionConfigBuilder Structure

This section describes the CompressionConfigBuilder structure.

##### 9.1.2.1 CompressionConfigBuilder Struct

```go
// CompressionConfigBuilder provides fluent configuration building for compression
type CompressionConfigBuilder struct {
    config *CompressionConfig
}
```

##### 9.1.2.2 NewCompressionConfigBuilder Function

```go
// NewCompressionConfigBuilder creates a new compression configuration builder.
func NewCompressionConfigBuilder() *CompressionConfigBuilder
```

##### 9.1.2.3 CompressionConfigBuilder.WithCompressionType Method

```go
// WithCompressionType sets the compression type for the configuration.
func (b *CompressionConfigBuilder) WithCompressionType(compType uint8) *CompressionConfigBuilder
```

##### 9.1.2.4 CompressionConfigBuilder.WithCompressionLevel Method

```go
// WithCompressionLevel sets the compression level for the configuration.
func (b *CompressionConfigBuilder) WithCompressionLevel(level int) *CompressionConfigBuilder
```

##### 9.1.2.5 CompressionConfigBuilder.WithSolidCompression Method

```go
// WithSolidCompression enables or disables solid compression for the configuration.
func (b *CompressionConfigBuilder) WithSolidCompression(useSolid bool) *CompressionConfigBuilder
```

##### 9.1.2.6 CompressionConfigBuilder.WithMemoryStrategy Method

```go
// WithMemoryStrategy sets the memory strategy for the configuration.
func (b *CompressionConfigBuilder) WithMemoryStrategy(strategy MemoryStrategy) *CompressionConfigBuilder
```

##### 9.1.2.7 CompressionConfigBuilder.Build Method

```go
// Build constructs and returns the final compression configuration.
func (b *CompressionConfigBuilder) Build() *CompressionConfig
```

### 9.2 Compression Validation Patterns

This section describes compression validation patterns.

#### 9.2.1 CompressionValidator Structure

This section describes the CompressionValidator structure.

##### 9.2.1.1 CompressionValidator Struct

```go
// CompressionValidator provides compression-specific validation
type CompressionValidator struct {
    *Validator[[]byte]
    compressionRules []CompressionValidationRule
}
```

##### 9.2.1.2 CompressionValidator.AddCompressionRule Method

```go
// AddCompressionRule adds a compression validation rule to the validator.
func (v *CompressionValidator) AddCompressionRule(rule CompressionValidationRule)
```

##### 9.2.1.3 CompressionValidator.ValidateCompressionData Method

```go
// Returns *PackageError on failure
func (v *CompressionValidator) ValidateCompressionData(ctx context.Context, data []byte) error
```

##### 9.2.1.4 CompressionValidator.ValidateDecompressionData Method

```go
// Returns *PackageError on failure
func (v *CompressionValidator) ValidateDecompressionData(ctx context.Context, data []byte) error
```

#### 9.2.2 CompressionValidationRule Structure

```go
// CompressionValidationRule represents a compression-specific validation rule
type CompressionValidationRule struct {
    Name      string
    Predicate func([]byte) bool
    Message   string
}
```

## 10. Compression and Signing Relationship

This section describes the relationship between compression and signing operations.

### 10.1 Signing Compressed Packages

This section describes signing compressed packages.

#### 10.1.1 Supported Operation

Compressed packages can be signed

#### 10.1.2 Signing Compressed Packages Process

1. Compress package content using `CompressPackage` or `CompressPackageFile`
2. Sign the compressed package using signature methods
3. Signatures validate the compressed content

#### 10.1.3 Signing Compressed Packages Benefits

- Faster signature validation (less data to hash during validation)
- Reduced overall package storage requirements (compressed content reduces total package size)

### 10.2 Compressing Signed Packages

This section describes compressing signed packages.

#### 10.2.1 Not Supported

Signed packages cannot be compressed

#### 10.2.2 Compressing Signed Packages Reasoning

- Signatures validate specific content
- Compression would change the content being validated
- Would invalidate existing signatures

#### 10.2.3 Compressing Signed Packages Workflow

1. If package is signed, decompress first
2. Make changes to package
3. Recompress if desired
4. Re-sign the package

## 11. CompressionStrategy Selection

This section describes how to select compression strategies.

### 11.1 Compression Type Selection

This section describes compression type selection criteria.

#### 11.1.1 Zstandard Compression (Type 1)

- Best compression ratio
- Moderate CPU usage
- Good for archival storage

#### 11.1.2 LZ4 Compression (Type 2)

- Fastest compression/decompression
- Lower compression ratio
- Good for real-time applications

#### 11.1.3 LZMA Compression (Type 3)

- Highest compression ratio
- Highest CPU usage
- Best for long-term storage

### 11.2 Compression Decision Matrix

This section provides a decision matrix for selecting compression strategies.

#### 11.2.1 User Guidance Matrix

The following table provides guidance for manual compression type selection based on intended use case:

| Use Case             | Recommended Type | Reason               |
| -------------------- | ---------------- | -------------------- |
| Real-time processing | LZ4              | Speed priority       |
| Archival storage     | Zstandard        | Balanced performance |
| Maximum compression  | LZMA             | Size priority        |
| Network transfer     | Zstandard        | Good balance         |

#### 11.2.2 Automatic Compression Type Selection

When compression is requested but compression type is not explicitly specified (compressionType = 0 or omitted), the API automatically selects the optimal compression algorithm based on package characteristics.

##### 11.2.2.1 Selection Algorithm

The automatic selection analyzes the following package properties:

- **Total Package Size**: Uncompressed size of all file entries, data, and index
- **File Count**: Number of files in the package
- **File Type Distribution**: Classification of files by type (text, binary, already-compressed)
- **Average File Size**: Total size divided by file count
- **Content Compressibility**: Estimated compression potential based on file types

##### 11.2.2.2 Selection Rules

The algorithm applies the following rules in order of priority:

1. **Already-Compressed Content Detection**:
    - If >50% of package content consists of already-compressed formats (JPEG, PNG, GIF, MP3, MP4, OGG, FLAC), select **LZ4** (fast, minimal benefit from heavy compression)

2. **Small Package Optimization**:
    - If total package size < 10MB, select **LZ4** (speed over compression ratio for small packages)
    - Rationale: Compression overhead outweighs benefits for small packages

3. **Many Small Files**:
    - If file count > 100 AND average file size < 10KB, select **LZ4** (fast compression for many small files)
    - Rationale: Package structure overhead makes compression ratio less important than speed

4. **Large Package with Text-Heavy Content**:
    - If total package size > 100MB AND text-based files (text, scripts, configs) represent >60% of content, select **LZMA** (maximum compression for compressible content)
    - Rationale: Text compresses well, large size justifies CPU cost

5. **Large Package with Mixed Content**:
    - If total package size > 100MB AND text-based files represent 30-60% of content, select **Zstandard** (balanced compression)
    - Rationale: Mixed content benefits from balanced approach

6. **Large Package with Binary-Heavy Content**:
    - If total package size > 100MB AND binary files represent >60% of content, select **Zstandard** (good compression for binary with reasonable speed)
    - Rationale: Binary doesn't compress as well, balanced approach optimal

7. **Medium Package Default**:
    - If total package size 10MB - 100MB, select **Zstandard** (balanced performance for medium packages)
    - Rationale: Default balanced approach for moderate sizes

8. **Fallback Default**:
    - If no specific rules apply, select **Zstandard** (safe balanced default)
    - Rationale: Zstandard provides good balance of speed and compression for most scenarios

##### 11.2.2.3 File Type Classification

For selection algorithm purposes, files are classified as:

- **Text-based**: Text files, scripts, configuration files, source code, JSON, XML, CSV
- **Binary**: Executables, compiled binaries, databases, proprietary formats
- **Already-compressed**: JPEG, PNG, GIF, MP3, MP4, OGG, FLAC, ZIP, GZIP, etc.
- **Media**: Images, audio, video (excludes already-compressed formats)

The algorithm uses `SelectCompressionType` logic for individual file classification where applicable.

##### 11.2.2.4 Implementation Behavior

When automatic selection is triggered:

- Selection occurs during `CompressPackage`, `CompressPackageFile`, `CompressPackageStream`, or `Write` operations
- Selected compression type is logged for debugging/monitoring
- User can override by explicitly specifying compressionType (1-3)
- Selection is consistent: same package properties always yield same selection

##### 11.2.2.5 Performance Considerations

Automatic selection has minimal overhead:

- Analysis uses existing package metadata (file entries, sizes, types)
- No content scanning required beyond metadata lookup
- Selection decision is O(n) where n is file count
- Selection adds <1ms overhead for typical packages

### 11.3 Compression Workflow Options

This section describes different compression workflow options.

#### 11.3.1 Option 1: Compress Before Writing

Compress the package content in memory first, then write the compressed package to disk.

##### 11.3.1.1 Process (Option 1)

Call `CompressPackage` with the desired compression type to compress the package content in memory.

Call `Write` with `CompressionNone` to write the already-compressed package to the output file without additional compression.

#### 11.3.2 Option 2: Compress and Write In One Step

Compress the package content and write it to disk in a single operation.

##### 11.3.2.1 Option 2 Process

Call `CompressPackageFile` with the target file path, compression type, and overwrite flag to compress and write the package in one step.

#### 11.3.3 Option 3: Write with Compression

Write the package to disk with compression applied during the write operation.

##### 11.3.3.1 Process (Option 3)

Call `Write` with the target file path, compression type, and overwrite flag to write the package with compression applied during the write process.

#### 11.3.4 Option 4: Stream Compression for Large Packages

Use streaming compression for large packages that may exceed available memory.

##### 11.3.4.1 Stream Configuration

Create a `StreamConfig` with appropriate chunk size settings.

Set `ChunkSize` to a reasonable size such as 1MB for processing chunks.

Enable `UseTempFiles` to use temporary files for large packages that exceed memory limits.

##### 11.3.4.2 Option 4 Process

Call `CompressPackageStream` with the compression type and stream configuration to compress the package using streaming.

Call `Write` with `CompressionNone` to write the compressed package to the output file.

#### 11.3.5 Option 5: Advanced Streaming Compression

For extremely large packages or when maximum performance is required, use advanced streaming compression with full configuration options that align with modern best practices from 7zip, zstd, and tar.

##### 11.3.5.1 Configuration Setup

Create a `StreamConfig` with intelligent defaults that allow the system to auto-detect optimal values.

Set `ChunkSize` to 0 for automatic calculation based on available memory.

Use an empty string for `TempDir` to utilize the system's temporary directory.

Set `MaxMemoryUsage` to 0 for automatic detection based on system RAM.

Select `MemoryStrategyBalanced` to use 50% of available RAM for optimal performance.

Enable `AdaptiveChunking` to allow the system to adjust chunk size based on memory pressure.

##### 11.3.5.2 Performance Configuration

Enable `UseDiskBuffering` for intermediate buffering when memory limits are reached.

Set `CleanupTempFiles` to true for automatic cleanup of temporary files.

Enable `UseParallelProcessing` for multi-core processing.

Set `MaxWorkers` to 0 for automatic CPU core detection.

Set `CompressionLevel` to 0 for automatic selection of the optimal compression level.

##### 11.3.5.3 Advanced Features

Enable `UseSolidCompression` for better compression ratios by treating multiple files as a single stream.

Set `ResumeFromOffset` to 0 to start from the beginning.

Set `BufferPoolSize` to 0 for automatic calculation of buffer pool size.

Set `MaxTempFileSize` to 0 for no limit on temporary file size.

Configure a `ProgressCallback` function to receive real-time progress updates during compression.

##### 11.3.5.4 Execution (Option 5)

Call `CompressPackageStream` with the ZSTD compression type and the configured settings.

Write the compressed package to the output file using `Write` with no additional compression.

#### 11.3.6 Option 6: Custom Memory Management

For specific memory constraints or performance requirements, configure custom memory management settings.

##### 11.3.6.1 Custom Configuration

Set `ChunkSize` to a specific value such as 512MB for controlled chunk processing.

Specify a custom `TempDir` path for temporary file storage.

Set `MaxMemoryUsage` to a specific limit such as 1GB for strict memory control.

Use `MemoryStrategyCustom` to utilize the explicit `MaxMemoryUsage` value.

Disable `AdaptiveChunking` to prevent automatic chunk size adjustments.

Set `BufferPoolSize` to a specific number of buffers for predictable memory usage.

Configure `MaxTempFileSize` to limit individual temporary file sizes.

##### 11.3.6.2 Performance Settings

Enable `UseParallelProcessing` for multi-core utilization.

Set `MaxWorkers` to a specific number to limit concurrent workers.

Specify a particular `CompressionLevel` for consistent compression behavior.

##### 11.3.6.3 Option 6 Execution

Call `CompressPackageStream` with the ZSTD compression type and the custom configuration.

## 12. Error Handling

The compression API uses the comprehensive structured error system defined in [api_core.md](api_core.md#10-structured-error-system).

**Generic Error Context**: All compression error-returning functions use `WrapErrorWithContext` or `NewPackageError` with typed error context structures for type-safe error handling.

Functions like `CompressPackage`, `DecompressPackage`, `CompressGeneric`, and other compression operations return errors that use the generic error context helpers:

- `WrapErrorWithContext[T]`: Wraps errors with typed context structures
- `NewPackageError[T]`: Creates structured errors with typed context
- `GetErrorContext[T]`: Retrieves type-safe context from errors

See [Generic Error Context Helpers](api_core.md#105-error-helper-functions) for complete documentation.

### 12.1 Common Error Conditions

This section describes common error conditions in compression operations.

#### 12.1.1 Common Error Conditions Compression Errors

- **Security Errors**: Package is already signed (cannot compress signed packages)
- **Validation Errors**: Invalid compression algorithm, package already compressed with different type
- **Compression Errors**: Compression operation failed, algorithm-specific failures

#### 12.1.2 Common Error Conditions Decompression Errors

- **Validation Errors**: Package is not compressed, invalid compressed data format
- **Compression Errors**: Decompression operation failed, algorithm-specific failures
- **Corruption Errors**: Compressed data is corrupted, checksum validation failed

#### 12.1.3 Common Error Conditions File Operation Errors

- **Validation Errors**: Target file exists and overwrite=false, invalid file path
- **I/O Errors**: I/O operation failed, disk space insufficient
- **Security Errors**: Insufficient permissions, access denied

### 12.2 Error Recovery

This section describes error recovery strategies for compression operations.

#### 12.2.1 Error Recovery Compression Failure

- Package remains in original state
- No partial compression state
- Can retry with different compression type

#### 12.2.2 Error Recovery Decompression Failure

- Package remains compressed
- Original compressed data preserved
- Can attempt recovery or use backup

## 13. Modern Best Practices

This section describes modern best practices for compression operations.

### 13.1 Industry Standards Alignment

Our compression API aligns with modern best practices used by leading compression systems:

#### 13.1.1 Streaming Compression (Universal Standard)

- **ZSTD Streaming**: Uses `ZSTD_compressStream2` and `ZSTD_decompressStream2` for large files
- **Memory Efficiency**: Constant memory usage regardless of file size
- **Real-time Processing**: Enables compression of files larger than available RAM
- **Progress Reporting**: Industry-standard progress callbacks for user feedback

#### 13.1.2 Parallel Processing (Performance Critical)

- **Multi-core Utilization**: Automatically detects and uses available CPU cores
- **Worker Pool Management**: Configurable worker count for optimal performance
- **Load Balancing**: Distributes chunks across workers for maximum throughput
- **Memory Isolation**: Each worker operates within memory limits

#### 13.1.3 Chunked Processing (Industry Standard)

- **Configurable Chunk Size**: Default 1GB chunks, adjustable based on system resources
- **Adaptive Sizing**: Automatically adjusts chunk size based on available memory
- **Resumable Operations**: Can resume from any chunk boundary
- **Progress Tracking**: Real-time progress reporting per chunk

#### 13.1.4 Memory Management (Resource Critical)

- **Strict Limits**: Enforces maximum memory usage to prevent OOM
- **Disk Fallback**: Automatic fallback to disk buffering when memory limits hit
- **Temporary File Management**: Intelligent temp file cleanup and management
- **Buffer Pooling**: Reuses buffers to minimize allocation overhead
- **Intelligent Defaults**: Auto-detects system capabilities and sets optimal values
- **Adaptive Sizing**: Automatically adjusts memory usage based on available RAM
- **Memory Strategies**: Conservative, Balanced, Aggressive, or Custom approaches

### 13.2 Intelligent Defaults and Memory Management

This section describes intelligent defaults and memory management for compression.

#### 13.2.1 MemoryStrategy Defaults

The system provides intelligent defaults based on system capabilities:

##### 13.2.1.1 Conservative Strategy (25% RAM)

- Use when system has limited RAM or other processes need memory
- Default for systems with <4GB RAM
- Ensures system stability during compression

##### 13.2.1.2 Balanced Strategy (50% RAM) - DEFAULT

- Optimal balance between performance and system stability
- Default for systems with 4-16GB RAM
- Provides good compression speed while leaving system responsive

##### 13.2.1.3 Aggressive Strategy (75% RAM)

- Maximum performance for dedicated compression systems
- Default for systems with >16GB RAM
- Use when system is dedicated to compression tasks

##### 13.2.1.4 Custom Strategy

- Use explicit `MaxMemoryUsage` value
- Override automatic detection
- Useful for specific memory constraints

#### 13.2.2 Auto-Detection Logic

The system automatically detects optimal memory settings based on available system resources.

##### 13.2.2.1 Memory Detection Process

The system queries available system RAM and calculates appropriate memory limits based on the selected strategy.

For systems with less than 4GB RAM, the Conservative strategy is automatically selected, allocating 25% of total RAM for compression operations.

Systems with 4-16GB RAM use the Balanced strategy by default, utilizing 50% of available RAM for optimal performance while maintaining system responsiveness.

Systems with more than 16GB RAM automatically select the Aggressive strategy, using 75% of available RAM for maximum compression performance.

##### 13.2.2.2 Chunk Size Calculation

When chunk size is not explicitly specified, the system calculates an optimal chunk size as 25% of the allocated memory limit.

This ensures that each processing chunk fits comfortably within the memory constraints while allowing for multiple concurrent operations.

##### 13.2.2.3 Worker Count Detection

The system automatically detects the number of available CPU cores and sets the worker count accordingly.

This enables optimal parallel processing without overloading the system with excessive worker threads.

#### 13.2.3 Adaptive Memory Management

- **Memory Monitoring**: Continuously monitors available memory during compression
- **Dynamic Adjustment**: Reduces chunk size if memory pressure detected
- **Disk Fallback**: Automatically switches to disk buffering when memory limits hit
- **Buffer Pooling**: Reuses buffers to minimize allocation overhead
- **Temp File Rotation**: Rotates temp files when they exceed `MaxTempFileSize`

### 13.3 Performance Considerations (CompressionOperations)

This section describes performance considerations for compression operations.

### 13.4 Memory Usage

This section describes memory usage considerations for compression.

#### 13.4.1 Compression Memory Usage

- Requires additional memory for compression buffers
- Memory usage scales with package size
- Consider streaming for very large packages
- **Large Files**: Use `CompressPackageStream` with appropriate memory limits and advanced configuration
- **Memory Management**: Automatic fallback to disk buffering when memory limits exceeded

#### 13.4.2 Decompression Memory Usage

- Requires memory for decompressed content
- May need temporary storage for large packages
- Use streaming for memory-constrained environments
- **Large Files**: Uses chunked decompression with temp file management
- **Memory Limits**: Enforces `MaxMemoryUsage` to prevent system OOM

### 13.5 CPU Usage

This section describes CPU usage considerations for compression.

#### 13.5.1 CPU Usage Compression

- LZ4: Lowest CPU usage
- Zstandard: Moderate CPU usage
- LZMA: Highest CPU usage

#### 13.5.2 CPU Usage Decompression

- Generally faster than compression
- LZ4: Fastest decompression
- Zstandard: Moderate decompression speed
- LZMA: Slowest decompression

### 13.6 I/O Considerations

This section describes I/O considerations for compression operations.

#### 13.6.1 I/O Considerations File-Based Operations

- Use streaming for large packages
- Consider disk space requirements
- Monitor I/O performance impact

#### 13.6.2 ~~I/O Considerations Network Operations~~ (REMOVED)

UPDATE: Removed: Out of scope.

- ~~Compressed packages transfer faster~~
- ~~Consider compression overhead vs. transfer time~~
- ~~Use appropriate compression type for network speed~~

## 14. Structured Error System

This section describes the structured error system for compression operations.

### 14.1 Structured Error System (Compression API)

The NovusPack package compression API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities.
For complete error system documentation, see [Structured Error System](api_core.md#10-structured-error-system).

### 14.2 Common Compression Error Types

This section describes common compression error types.

#### 14.2.1 Compression Error Types

The NovusPack compression API uses the structured error system with the following error types:

- `ErrTypeCompression`: Compression and decompression operation failures
- `ErrTypeValidation`: Invalid compression parameters and data validation errors
- `ErrTypeIO`: I/O errors during compression operations
- `ErrTypeContext`: Context cancellation and timeout errors
- `ErrTypeCorruption`: Corrupted compressed data errors
- `ErrTypeUnsupported`: Unsupported compression algorithms and features

### 14.3 Structured Error Examples

This section provides examples of structured errors in compression operations.

#### 14.3.1 Creating Compression Errors

This section describes how to create compression errors.

##### 14.3.1.1 Error Context Type Definitions

This section defines error context type definitions used in compression errors.

##### 14.3.1.2 CompressionErrorContext Structure

```go
// Define error context types
type CompressionErrorContext struct {
    Algorithm string
    Level     int
    InputSize int64
    Operation string
}
```

##### 14.3.1.3 UnsupportedCompressionErrorContext Structure

```go
// UnsupportedCompressionErrorContext provides error context for unsupported compression type errors.
type UnsupportedCompressionErrorContext struct {
    CompressionType uint8
    SupportedTypes  []uint8
    Operation       string
}
```

##### 14.3.1.4 MemoryErrorContext Structure

```go
// MemoryErrorContext provides error context for memory-related compression errors.
type MemoryErrorContext struct {
    RequiredMemory  string
    AvailableMemory string
    Algorithm       string
    Operation       string
}
```

##### 14.3.1.5 Creating Errors with Context

```go
// Compression failure with typed context
err := NewPackageError(ErrTypeCompression, "compression failed", nil, CompressionErrorContext{
    Algorithm: "Zstd",
    Level:     6,
    InputSize: 1024 * 1024,
    Operation: "CompressPackage",
})

// Unsupported compression type with typed context
err := NewPackageError(ErrTypeUnsupported, "unsupported compression type", nil, UnsupportedCompressionErrorContext{
    CompressionType: 99,
    SupportedTypes:  []uint8{0, 1, 2, 3},
    Operation:       "SetCompressionType",
})

// Memory error with typed context
err := NewPackageError(ErrTypeIO, "insufficient memory", nil, MemoryErrorContext{
    RequiredMemory:  "512MB",
    AvailableMemory: "256MB",
    Algorithm:       "LZMA",
    Operation:       "CompressPackage",
})
```

#### 14.3.2 Error Handling Patterns

Use the structured error system to handle compression errors appropriately.

Check error types and extract context information for proper error handling and logging.

Handle different error categories (compression, I/O, context) with appropriate responses.
