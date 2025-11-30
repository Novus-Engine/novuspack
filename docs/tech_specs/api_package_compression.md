# NovusPack Technical Specifications - Package Compression API

## Table of Contents

- [Table of Contents](#table-of-contents)
- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
  - [0.2 Context Integration](#02-context-integration)
- [1. Package Compression Overview](#1-package-compression-overview)
  - [1.1 Compression Scope](#11-compression-scope)
  - [1.2 Compression Types](#12-compression-types)
  - [1.3 Compression Information Structure](#13-compression-information-structure)
  - [1.4 Compression Constraints](#14-compression-constraints)
- [2. Strategy Pattern Interfaces](#2-strategy-pattern-interfaces)
  - [2.1 Compression Strategy Interface](#21-compression-strategy-interface)
  - [2.2 Built-in Compression Strategies](#22-built-in-compression-strategies)
- [3. Interface Granularity and Composition](#3-interface-granularity-and-composition)
  - [3.1 Compression Information Interface](#31-compression-information-interface)
  - [3.2 Compression Operations Interface](#32-compression-operations-interface)
  - [3.3 Compression Streaming Interface](#33-compression-streaming-interface)
  - [3.4 Compression File Operations Interface](#34-compression-file-operations-interface)
  - [3.5 Generic Compression Interface](#35-generic-compression-interface)
- [4. In-Memory Compression Methods](#4-in-memory-compression-methods)
  - [4.1 CompressPackage](#41-compresspackage)
  - [4.2 DecompressPackage](#42-decompresspackage)
- [5. Streaming Compression Methods](#5-streaming-compression-methods)
  - [5.1 CompressPackageStream](#51-compresspackagestream)
  - [5.2 DecompressPackageStream](#52-decompresspackagestream)
- [6. File-Based Compression Methods](#6-file-based-compression-methods)
  - [6.1 CompressPackageFile](#61-compresspackagefile)
  - [6.2 DecompressPackageFile](#62-decompresspackagefile)
- [7. Compression Information and Status](#7-compression-information-and-status)
  - [7.1 Compression Information Structure](#71-compression-information-structure)
  - [7.2 Compression Status Methods](#72-compression-status-methods)
  - [7.3 Internal Compression Methods](#73-internal-compression-methods)
- [8. Concurrency Patterns and Thread Safety](#8-concurrency-patterns-and-thread-safety)
  - [8.1 Thread Safety Guarantees](#81-thread-safety-guarantees)
  - [8.2 Worker Pool Management](#82-worker-pool-management)
  - [8.3 Concurrent Compression Methods](#83-concurrent-compression-methods)
  - [8.4 Resource Management](#84-resource-management)
- [9. Compression Configuration Patterns](#9-compression-configuration-patterns)
  - [9.1 Compression-Specific Configuration](#91-compression-specific-configuration)
  - [9.2 Compression Validation Patterns](#92-compression-validation-patterns)
- [10. Compression and Signing Relationship](#10-compression-and-signing-relationship)
  - [10.1 Signing Compressed Packages](#101-signing-compressed-packages)
  - [10.2 Compressing Signed Packages](#102-compressing-signed-packages)
- [11. Compression Strategy Selection](#11-compression-strategy-selection)
  - [11.1 Compression Type Selection](#111-compression-type-selection)
  - [11.2 Compression Decision Matrix](#112-compression-decision-matrix)
  - [11.3 Compression Workflow Options](#113-compression-workflow-options)
- [12. Error Handling](#12-error-handling)
  - [12.1 Common Error Conditions](#121-common-error-conditions)
  - [12.2 Error Recovery](#122-error-recovery)
- [13. Modern Best Practices](#13-modern-best-practices)
  - [13.1 Industry Standards Alignment](#131-industry-standards-alignment)
  - [13.2 Intelligent Defaults and Memory Management](#132-intelligent-defaults-and-memory-management)
  - [13.3 Performance Considerations](#133-performance-considerations)
  - [13.4 Memory Usage](#134-memory-usage)
  - [13.5 CPU Usage](#135-cpu-usage)
  - [13.6 I/O Considerations](#136-io-considerations)
- [14. Structured Error System](#14-structured-error-system)
  - [14.1 Structured Error System](#141-structured-error-system)
  - [14.2 Common Compression Error Types](#142-common-compression-error-types)
  - [14.3 Structured Error Examples](#143-structured-error-examples)

## 0. Overview

This document defines the NovusPack package compression API, providing methods for compressing and decompressing package content while maintaining package integrity and signature compatibility.

### 0.1 Cross-References

- [Core Package Interface API](api_core.md) - Package operations and compression
- [Package Writing API](api_writing.md) - SafeWrite, FastWrite, and write strategy selection
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
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

Package compression in NovusPack compresses the package content while preserving the header, package comment, and signatures in an uncompressed state for direct access.

### 1.1 Compression Scope

#### 1.1.1 Compressed Content

- File entries (directory structure)
- File data (actual file contents)
- Package index

#### 1.1.2 Uncompressed Content

- Package header (see [Package File Format - Package Header](package_file_format.md#2-package-header))
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

### 1.3 Compression Information Structure

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

## 2. Strategy Pattern Interfaces

The compression API supports pluggable compression algorithms through the strategy pattern.

### 2.1 Compression Strategy Interface

```go
// CompressionStrategy implements Strategy[[]byte, []byte] with generic support
type CompressionStrategy[T any] interface {
    Compress(ctx context.Context, data T) (T, error)
    Decompress(ctx context.Context, data T) (T, error)
    Type() CompressionType
    Name() string
}

// ByteCompressionStrategy is the concrete implementation for []byte data
type ByteCompressionStrategy interface {
    CompressionStrategy[[]byte]
}

// AdvancedCompressionStrategy for compression with additional validation and metrics
type AdvancedCompressionStrategy[T any] interface {
    CompressionStrategy[T]
    ValidateInput(ctx context.Context, data T) error
    GetCompressionRatio(ctx context.Context, original T, compressed T) float64
}

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

```go
// Zstandard compression strategy with generic support
type ZstandardStrategy[T any] struct {
    level    int
    strategy CompressionStrategy[T]
}

// LZ4 compression strategy with generic support
type LZ4Strategy[T any] struct {
    level    int
    strategy CompressionStrategy[T]
}

// LZMA compression strategy with generic support
type LZMAStrategy[T any] struct {
    level    int
    strategy CompressionStrategy[T]
}

// CompressionJob represents a unit of work for compression (extends Job)
type CompressionJob[T any] struct {
    *Job[T]
    CompressionType uint8
    CompressionLevel int
}
```

## 3. Interface Granularity and Composition

The compression API uses focused interfaces to provide clear separation of concerns and enable flexible composition.

### 3.1 Compression Information Interface

```go
// CompressionInfo provides read-only access to compression information
type CompressionInfo interface {
    GetCompressionInfo(ctx context.Context) PackageCompressionInfo
    IsCompressed() bool
    GetCompressionType() (uint8, bool) // Returns (type, isSet)
    GetCompressionRatio() (float64, bool) // Returns (ratio, isSet)
    CanCompress() bool
}
```

### 3.2 Compression Operations Interface

```go
// CompressionOperations provides basic compression/decompression operations
type CompressionOperations interface {
    CompressPackage(ctx context.Context, compressionType uint8) error
    DecompressPackage(ctx context.Context) error
    SetCompressionType(ctx context.Context, compressionType uint8) error
}
```

### 3.3 Compression Streaming Interface

```go
// CompressionStreaming provides streaming compression for large packages
type CompressionStreaming interface {
    CompressPackageStream(ctx context.Context, compressionType uint8, config *StreamConfig) error
    DecompressPackageStream(ctx context.Context, config *StreamConfig) error
}
```

### 3.4 Compression File Operations Interface

```go
// CompressionFileOperations provides file-based compression operations
type CompressionFileOperations interface {
    CompressPackageFile(ctx context.Context, path string, compressionType uint8, overwrite bool) error
    DecompressPackageFile(ctx context.Context, path string, overwrite bool) error
}
```

### 3.5 Generic Compression Interface

```go
// Compression provides type-safe compression for any data type
type Compression[T any] interface {
    CompressGeneric(ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
    DecompressGeneric(ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
    ValidateCompressionData(ctx context.Context, data T) error
}
```

## 4. In-Memory Compression Methods

These methods operate on packages in memory without writing to disk.

**Note**: For large packages, consider using [Streaming Compression Methods](#5-streaming-compression-methods) to avoid memory limitations.

### 4.1 CompressPackage

```go
// CompressPackage compresses package content in memory
// Compresses file entries + data + index (NOT header, comment, or signatures)
// Signed packages cannot be compressed
func (p *Package) CompressPackage(ctx context.Context, compressionType uint8) error
```

#### 4.1.1 CompressPackage Purpose

Handle compression/decompression of in-memory packages

#### 4.1.2 CompressPackage Parameters

- `ctx`: Context for cancellation and timeout handling
- `compressionType`: Compression algorithm to use (1-3)

#### 4.1.3 CompressPackage Behavior

- Compresses file entries + data + index (NOT header, comment, or signatures)
- Updates package compression state in memory
- Returns error if package is signed
- Updates package header compression flags

#### 4.1.4 CompressPackage Error Conditions

- Package is already signed
- Invalid compression type
- Package is already compressed with different type
- Context cancellation

### 4.2 DecompressPackage

```go
// DecompressPackage decompresses the package in memory
// Decompresses all compressed content
func (p *Package) DecompressPackage(ctx context.Context) error
```

#### 4.2.1 DecompressPackage Purpose

Decompress package content in memory

#### 4.2.2 DecompressPackage Parameters

- `ctx`: Context for cancellation and timeout handling

#### 4.2.3 DecompressPackage Behavior

- Decompresses all compressed content
- Updates package compression state in memory
- Clears package header compression flags
- Preserves all other package data

#### 4.2.4 DecompressPackage Error Conditions

- Package is not compressed
- Decompression failure
- Context cancellation

## 5. Streaming Compression Methods

These methods handle compression/decompression of large packages using streaming to avoid memory limitations.

**For Large Files**: These methods use temporary files and chunked processing to handle files that exceed available RAM, with adaptive strategies based on configuration.

### 5.1 CompressPackageStream

```go
// CompressPackageStream compresses large package content using streaming
// Uses temporary files and chunked processing to handle files of any size
// Configuration determines the level of optimization and memory management
func (p *Package) CompressPackageStream(ctx context.Context, compressionType uint8, config *StreamConfig) error
```

#### 5.1.1 Purpose

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

**Simple Usage** (basic settings only):

```go
config := &StreamConfig{
    ChunkSize: 0,        // Auto-calculate
    MaxMemoryUsage: 0,   // Auto-detect
    TempDir: "",         // System temp
}
```

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

### 5.2 DecompressPackageStream

```go
// DecompressPackageStream decompresses large package content using streaming
// Uses streaming to manage memory efficiently for large packages
func (p *Package) DecompressPackageStream(ctx context.Context, config *StreamConfig) error
```

#### 5.2.1 Purpose

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

### 6.1 CompressPackageFile

```go
// CompressPackageFile compresses package content and writes to specified path
// Compresses file entries + data + index (NOT header, comment, or signatures)
// Signed packages cannot be compressed
func (p *Package) CompressPackageFile(ctx context.Context, path string, compressionType uint8, overwrite bool) error
```

#### 6.1.1 Purpose

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

### 6.2 DecompressPackageFile

```go
// DecompressPackageFile decompresses the package and writes to specified path
// Decompresses all compressed content and writes uncompressed package
func (p *Package) DecompressPackageFile(ctx context.Context, path string, overwrite bool) error
```

#### 6.2.1 Purpose

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

### 7.1 Compression Information Structure

See [Compression Information Structure](#71-compression-information-structure) for the complete structure definition.

### 7.2 Compression Status Methods

```go
// GetPackageCompressionInfo returns package compression information
func (p *Package) GetPackageCompressionInfo(ctx context.Context) PackageCompressionInfo

// IsPackageCompressed checks if the package is compressed
func (p *Package) IsPackageCompressed() bool

// GetPackageCompressionType returns the package compression type with optional pattern
func (p *Package) GetPackageCompressionType() (uint8, bool)

// GetPackageCompressionRatio returns the compression ratio with optional pattern
func (p *Package) GetPackageCompressionRatio() (float64, bool)

// GetPackageOriginalSize returns the original size before compression with optional pattern
func (p *Package) GetPackageOriginalSize() (int64, bool)

// GetPackageCompressedSize returns the compressed size with optional pattern
func (p *Package) GetPackageCompressedSize() (int64, bool)

// SetPackageCompressionType sets the package compression type (without compressing)
func (p *Package) SetPackageCompressionType(ctx context.Context, compressionType uint8) error

// CanCompressPackage checks if package can be compressed (not signed)
func (p *Package) CanCompressPackage() bool

// Generic compression methods for type-safe operations
func (p *Package) CompressGeneric[T any](ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
func (p *Package) DecompressGeneric[T any](ctx context.Context, data T, strategy CompressionStrategy[T]) (T, error)
func (p *Package) ValidateCompressionData[T any](ctx context.Context, data T) error
```

### 7.3 Internal Compression Methods

```go
// Internal compression methods (used by CompressPackage and Write)
func (p *Package) compressPackageContent(ctx context.Context, compressionType uint8) error
func (p *Package) decompressPackageContent(ctx context.Context) error
```

## 8. Concurrency Patterns and Thread Safety

The compression API provides explicit concurrency patterns and thread safety guarantees for safe concurrent usage.

### 8.1 Thread Safety Guarantees

The compression API provides different levels of thread safety based on the `ThreadSafetyMode` configuration:

#### 8.1.1 ThreadSafetyNone

No thread safety guarantees.
Operations should not be called concurrently.

#### 8.1.2 ThreadSafetyReadOnly

Read-only operations are safe for concurrent access.
Multiple goroutines can safely call read methods simultaneously.

#### 8.1.3 ThreadSafetyConcurrent

Concurrent read/write operations are supported.
Uses read-write mutex for optimal read performance.

#### 8.1.4 ThreadSafetyFull

Full thread safety with complete synchronization.
All operations are protected by appropriate locking mechanisms.

### 8.2 Worker Pool Management

The compression API uses the generic worker pool patterns defined in [api_generics.md](api_generics.md#26-generic-concurrency-patterns) with compression-specific extensions.

```go
// CompressionWorkerPool extends WorkerPool for compression operations
type CompressionWorkerPool[T any] struct {
    *WorkerPool[T]
    compressionStrategy CompressionStrategy[T]
}

// Compression-specific methods
func (p *CompressionWorkerPool[T]) CompressConcurrently(ctx context.Context, data []T, strategy CompressionStrategy[T]) ([]T, error)
func (p *CompressionWorkerPool[T]) DecompressConcurrently(ctx context.Context, data []T, strategy CompressionStrategy[T]) ([]T, error)
func (p *CompressionWorkerPool[T]) GetCompressionStats() CompressionStats
```

### 8.3 Concurrent Compression Methods

```go
// CompressPackageConcurrent compresses package content using worker pool
func (p *Package) CompressPackageConcurrent(ctx context.Context, compressionType uint8, config *StreamConfig) error

// DecompressPackageConcurrent decompresses package content using worker pool
func (p *Package) DecompressPackageConcurrent(ctx context.Context, config *StreamConfig) error

// CompressMultiplePackages compresses multiple packages concurrently
func CompressMultiplePackages[T any](ctx context.Context, packages []*Package, compressionType uint8, config *StreamConfig) []error
```

### 8.4 Resource Management

The compression API uses the generic resource management patterns defined in [api_generics.md](api_generics.md#26-generic-concurrency-patterns) with compression-specific resource types.

```go
// CompressionResourcePool manages compression-specific resources
type CompressionResourcePool struct {
    *ResourcePool[CompressionResource]
    compressionConfig *CompressionConfig
}

// CompressionResource represents a compression-specific resource
type CompressionResource struct {
    ID           string
    Strategy     CompressionStrategy[[]byte]
    Buffer       []byte
    LastUsed     time.Time
    AccessCount  int64
}

// Compression-specific resource management methods
func (p *CompressionResourcePool) AcquireCompressionResource(ctx context.Context, strategyType uint8) (*CompressionResource, error)
func (p *CompressionResourcePool) ReleaseCompressionResource(resource *CompressionResource) error
func (p *CompressionResourcePool) GetCompressionResourceStats() CompressionResourceStats
```

## 9. Compression Configuration Patterns

The compression API provides compression-specific configuration patterns that extend the generic configuration patterns defined in [api_generics.md](api_generics.md#28-generic-configuration-patterns).

### 9.1 Compression-Specific Configuration

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

// CompressionConfigBuilder provides fluent configuration building for compression
type CompressionConfigBuilder struct {
    config *CompressionConfig
}

func NewCompressionConfigBuilder() *CompressionConfigBuilder
func (b *CompressionConfigBuilder) WithCompressionType(compType uint8) *CompressionConfigBuilder
func (b *CompressionConfigBuilder) WithCompressionLevel(level int) *CompressionConfigBuilder
func (b *CompressionConfigBuilder) WithSolidCompression(useSolid bool) *CompressionConfigBuilder
func (b *CompressionConfigBuilder) WithMemoryStrategy(strategy MemoryStrategy) *CompressionConfigBuilder
func (b *CompressionConfigBuilder) Build() *CompressionConfig
```

### 9.2 Compression Validation Patterns

```go
// CompressionValidator provides compression-specific validation
type CompressionValidator struct {
    *Validator[[]byte]
    compressionRules []CompressionValidationRule
}

// CompressionValidationRule represents a compression-specific validation rule
type CompressionValidationRule struct {
    Name      string
    Predicate func([]byte) bool
    Message   string
}

func (v *CompressionValidator) AddCompressionRule(rule CompressionValidationRule)
func (v *CompressionValidator) ValidateCompressionData(ctx context.Context, data []byte) error
func (v *CompressionValidator) ValidateDecompressionData(ctx context.Context, data []byte) error
```

## 10. Compression and Signing Relationship

### 10.1 Signing Compressed Packages

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

## 11. Compression Strategy Selection

### 11.1 Compression Type Selection

#### 11.1.1 Zstandard Compression (1)

- Best compression ratio
- Moderate CPU usage
- Good for archival storage

#### 11.1.2 LZ4 Compression (2)

- Fastest compression/decompression
- Lower compression ratio
- Good for real-time applications

#### 11.1.3 LZMA Compression (3)

- Highest compression ratio
- Highest CPU usage
- Best for long-term storage

### 11.2 Compression Decision Matrix

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

#### 11.3.1 Option 1: Compress Before Writing

Compress the package content in memory first, then write the compressed package to disk.

#### 11.3.1.1 Process

Call `CompressPackage` with the desired compression type to compress the package content in memory.

Call `Write` with `CompressionNone` to write the already-compressed package to the output file without additional compression.

#### 11.3.2 Option 2: Compress and Write in One Step

Compress the package content and write it to disk in a single operation.

#### 11.3.2.1 Process

Call `CompressPackageFile` with the target file path, compression type, and overwrite flag to compress and write the package in one step.

#### 11.3.3 Option 3: Write with Compression

Write the package to disk with compression applied during the write operation.

#### 11.3.3.1 Process

Call `Write` with the target file path, compression type, and overwrite flag to write the package with compression applied during the write process.

#### 11.3.4 Option 4: Stream Compression for Large Packages

Use streaming compression for large packages that may exceed available memory.

#### 11.3.4.1 Configuration

Create a `StreamConfig` with appropriate chunk size settings.

Set `ChunkSize` to a reasonable size such as 1MB for processing chunks.

Enable `UseTempFiles` to use temporary files for large packages that exceed memory limits.

#### 11.3.4.2 Process

Call `CompressPackageStream` with the compression type and stream configuration to compress the package using streaming.

Call `Write` with `CompressionNone` to write the compressed package to the output file.

#### 11.3.5 Option 5: Advanced Streaming Compression

For extremely large packages or when maximum performance is required, use advanced streaming compression with full configuration options that align with modern best practices from 7zip, zstd, and tar.

#### 11.3.5.1 Configuration Setup

Create a `StreamConfig` with intelligent defaults that allow the system to auto-detect optimal values.

Set `ChunkSize` to 0 for automatic calculation based on available memory.

Use an empty string for `TempDir` to utilize the system's temporary directory.

Set `MaxMemoryUsage` to 0 for automatic detection based on system RAM.

Select `MemoryStrategyBalanced` to use 50% of available RAM for optimal performance.

Enable `AdaptiveChunking` to allow the system to adjust chunk size based on memory pressure.

#### 11.3.5.2 Performance Configuration

Enable `UseDiskBuffering` for intermediate buffering when memory limits are reached.

Set `CleanupTempFiles` to true for automatic cleanup of temporary files.

Enable `UseParallelProcessing` for multi-core processing.

Set `MaxWorkers` to 0 for automatic CPU core detection.

Set `CompressionLevel` to 0 for automatic selection of the optimal compression level.

#### 11.3.5.3 Advanced Features

Enable `UseSolidCompression` for better compression ratios by treating multiple files as a single stream.

Set `ResumeFromOffset` to 0 to start from the beginning.

Set `BufferPoolSize` to 0 for automatic calculation of buffer pool size.

Set `MaxTempFileSize` to 0 for no limit on temporary file size.

Configure a `ProgressCallback` function to receive real-time progress updates during compression.

#### 11.3.5.4 Execution

Call `CompressPackageStream` with the ZSTD compression type and the configured settings.

Write the compressed package to the output file using `Write` with no additional compression.

#### 11.3.6 Option 6: Custom Memory Management

For specific memory constraints or performance requirements, configure custom memory management settings.

#### 11.3.6.1 Custom Configuration

Set `ChunkSize` to a specific value such as 512MB for controlled chunk processing.

Specify a custom `TempDir` path for temporary file storage.

Set `MaxMemoryUsage` to a specific limit such as 1GB for strict memory control.

Use `MemoryStrategyCustom` to utilize the explicit `MaxMemoryUsage` value.

Disable `AdaptiveChunking` to prevent automatic chunk size adjustments.

Set `BufferPoolSize` to a specific number of buffers for predictable memory usage.

Configure `MaxTempFileSize` to limit individual temporary file sizes.

#### 11.3.6.2 Performance Settings

Enable `UseParallelProcessing` for multi-core utilization.

Set `MaxWorkers` to a specific number to limit concurrent workers.

Specify a particular `CompressionLevel` for consistent compression behavior.

#### 11.3.6.3 Execution

Call `CompressPackageStream` with the ZSTD compression type and the custom configuration.

## 12. Error Handling

### 12.1 Common Error Conditions

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

#### 12.2.1 Error Recovery Compression Failure

- Package remains in original state
- No partial compression state
- Can retry with different compression type

#### 12.2.2 Error Recovery Decompression Failure

- Package remains compressed
- Original compressed data preserved
- Can attempt recovery or use backup

## 13. Modern Best Practices

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

#### 13.2.1 Memory Strategy Defaults

The system provides intelligent defaults based on system capabilities:

#### 13.2.1.1 Conservative Strategy (25% RAM)

- Use when system has limited RAM or other processes need memory
- Default for systems with <4GB RAM
- Ensures system stability during compression

#### 13.2.1.2 Balanced Strategy (50% RAM) - DEFAULT

- Optimal balance between performance and system stability
- Default for systems with 4-16GB RAM
- Provides good compression speed while leaving system responsive

#### 13.2.1.3 Aggressive Strategy (75% RAM)

- Maximum performance for dedicated compression systems
- Default for systems with >16GB RAM
- Use when system is dedicated to compression tasks

#### 13.2.1.4 Custom Strategy

- Use explicit `MaxMemoryUsage` value
- Override automatic detection
- Useful for specific memory constraints

#### 13.2.2 Auto-Detection Logic

The system automatically detects optimal memory settings based on available system resources.

#### 13.2.2.1 Memory Detection Process

The system queries available system RAM and calculates appropriate memory limits based on the selected strategy.

For systems with less than 4GB RAM, the Conservative strategy is automatically selected, allocating 25% of total RAM for compression operations.

Systems with 4-16GB RAM use the Balanced strategy by default, utilizing 50% of available RAM for optimal performance while maintaining system responsiveness.

Systems with more than 16GB RAM automatically select the Aggressive strategy, using 75% of available RAM for maximum compression performance.

#### 13.2.2.2 Chunk Size Calculation

When chunk size is not explicitly specified, the system calculates an optimal chunk size as 25% of the allocated memory limit.

This ensures that each processing chunk fits comfortably within the memory constraints while allowing for multiple concurrent operations.

#### 13.2.2.3 Worker Count Detection

The system automatically detects the number of available CPU cores and sets the worker count accordingly.

This enables optimal parallel processing without overloading the system with excessive worker threads.

#### 13.2.3 Adaptive Memory Management

- **Memory Monitoring**: Continuously monitors available memory during compression
- **Dynamic Adjustment**: Reduces chunk size if memory pressure detected
- **Disk Fallback**: Automatically switches to disk buffering when memory limits hit
- **Buffer Pooling**: Reuses buffers to minimize allocation overhead
- **Temp File Rotation**: Rotates temp files when they exceed `MaxTempFileSize`

### 13.3 Performance Considerations

### 13.4 Memory Usage

#### 13.4.1 Compression

- Requires additional memory for compression buffers
- Memory usage scales with package size
- Consider streaming for very large packages
- **Large Files**: Use `CompressPackageStream` with appropriate memory limits and advanced configuration
- **Memory Management**: Automatic fallback to disk buffering when memory limits exceeded

#### 13.4.2 Decompression

- Requires memory for decompressed content
- May need temporary storage for large packages
- Use streaming for memory-constrained environments
- **Large Files**: Uses chunked decompression with temp file management
- **Memory Limits**: Enforces `MaxMemoryUsage` to prevent system OOM

### 13.5 CPU Usage

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

### 14.1 Structured Error System

The NovusPack package compression API uses a comprehensive structured error system that provides better error categorization, context, and debugging capabilities.
For complete error system documentation, see [Structured Error System](api_core.md#11-structured-error-system).

### 14.2 Common Compression Error Types

#### 14.2.1 Compression Error Types

The NovusPack compression API uses the structured error system with the following error types:

- `ErrTypeCompression`: Compression and decompression operation failures
- `ErrTypeValidation`: Invalid compression parameters and data validation errors
- `ErrTypeIO`: I/O errors during compression operations
- `ErrTypeContext`: Context cancellation and timeout errors
- `ErrTypeCorruption`: Corrupted compressed data errors
- `ErrTypeUnsupported`: Unsupported compression algorithms and features

### 14.3 Structured Error Examples

#### 14.3.1 Creating Compression Errors

```go
// Compression failure with context
err := NewPackageError(ErrTypeCompression, "compression failed", nil).
    WithContext("algorithm", "Zstd").
    WithContext("level", 6).
    WithContext("inputSize", 1024*1024).
    WithContext("operation", "CompressPackage")

// Unsupported compression type with context
err := NewPackageError(ErrTypeUnsupported, "unsupported compression type", nil).
    WithContext("compressionType", 99).
    WithContext("supportedTypes", []uint8{0, 1, 2, 3}).
    WithContext("operation", "SetCompressionType")

// Memory error with context
err := NewPackageError(ErrTypeIO, "insufficient memory", nil).
    WithContext("requiredMemory", "512MB").
    WithContext("availableMemory", "256MB").
    WithContext("algorithm", "LZMA").
    WithContext("operation", "CompressPackage")
```

#### 14.3.2 Error Handling Patterns

Use the structured error system to handle compression errors appropriately.

Check error types and extract context information for proper error handling and logging.

Handle different error categories (compression, I/O, context) with appropriate responses.
