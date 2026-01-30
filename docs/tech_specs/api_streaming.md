# NovusPack Technical Specifications - Streaming and Buffer Management API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. File Streaming Interface](#1-file-streaming-interface)
  - [1.1 Purpose](#11-purpose)
  - [1.2 Core Types](#12-core-types)
    - [1.2.1 FileStream struct](#121-filestream-struct)
    - [1.2.2 StreamConfig struct](#122-streamconfig-struct)
  - [1.3 Key Methods](#13-key-methods)
    - [1.3.1 Factory Function](#131-newfilestream-function)
    - [1.3.2 Core Operations](#132-core-operations)
    - [1.3.3 Status and Query Methods](#133-status-and-query-methods)
    - [1.3.4 Standard Go Interface Methods](#134-standard-go-interface-methods)
  - [1.4 Features](#14-features)
  - [1.5 Additional Methods](#15-additional-methods)
    - [1.5.1 Stream Information](#151-stream-information)
    - [1.5.2 Progress Monitoring](#152-progress-monitoring)
    - [1.5.3 Standard Go Interfaces](#153-standard-go-interfaces)
- [2. Buffer Management System](#2-buffer-management-system)
  - [2.1 Buffer Management Purpose](#21-buffer-management-purpose)
  - [2.2 Buffer Management Core Types](#22-buffer-management-core-types)
    - [2.2.1 BufferPool struct](#221-bufferpool-struct)
    - [2.2.2 BufferConfig struct](#222-bufferconfig-struct)
  - [2.3 Buffer Management Key Methods](#23-buffer-management-key-methods)
    - [2.3.1 BufferPool Creation and Core Operations](#231-bufferpool-creation-and-core-operations)
    - [2.3.2 BufferPool Management Methods](#232-bufferpool-management-methods)
  - [2.4 Buffer Management Features](#24-buffer-management-features)
  - [2.5 Buffer Management Additional Methods](#25-buffer-management-additional-methods)
    - [2.5.1 TotalSize Method Details](#251-totalsize-method-details)
    - [2.5.2 SetMaxTotalSize Method Details](#252-setmaxtotalsize-method-details)
    - [2.5.3 BufferPool Management Purpose](#253-bufferpool-management-purpose)
    - [2.5.4 TotalSize Returns](#254-totalsize-returns)
    - [2.5.5 SetMaxTotalSize Parameters](#255-setmaxtotalsize-parameters)
    - [2.5.6 SetMaxTotalSize Behavior](#256-setmaxtotalsize-behavior)
    - [2.5.7 BufferPool Management Example Usage](#257-bufferpool-management-example-usage)
  - [2.6 Default Configuration](#26-default-configuration)
    - [2.6.1 DefaultBufferConfig Function](#261-defaultbufferconfig-function)
- [3. Streaming Concurrency Patterns](#3-streaming-concurrency-patterns)
  - [3.1 Streaming Concurrency Purpose](#31-streaming-concurrency-purpose)
  - [3.2 Streaming Concurrency Core Types](#32-streaming-concurrency-core-types)
    - [3.2.1 StreamingWorkerPool struct](#321-streamingworkerpool-struct)
  - [3.3 Streaming Concurrency Key Methods](#33-streaming-concurrency-key-methods)
    - [3.3.1 Worker Pool Lifecycle Methods](#331-worker-pool-lifecycle-methods)
    - [3.3.2 Job Submission and Processing Methods](#332-job-submission-and-processing-methods)
    - [3.3.3 Statistics Method](#333-streamingworkerpoolgetstreamingstats-method)
  - [3.4 Streaming Concurrency Features](#34-streaming-concurrency-features)
- [4. Streaming Configuration Patterns](#4-streaming-configuration-patterns)
  - [4.1 Streaming Configuration Purpose](#41-streaming-configuration-purpose)
  - [4.2 Streaming Configuration Core Types](#42-streaming-configuration-core-types)
    - [4.2.1 StreamingConfig struct](#421-streamingconfig-struct)
  - [4.3 Streaming Configuration Key Methods](#43-streaming-configuration-key-methods)
    - [4.3.1 StreamingConfigDefaults Structure](#431-streamingconfigdefaults-structure)
    - [4.3.2 CreateStreamingConfig Function](#432-createstreamingconfig-function)
    - [4.3.3 ValidateStreamingConfig Function](#433-validatestreamingconfig-function)
    - [4.3.4 GetStreamingConfigDefaults Function](#434-getstreamingconfigdefaults-function)

---

## 0. Overview

This document defines the streaming and buffer management API for the NovusPack system, providing memory-efficient handling of large files with compression and encryption support.

### 0.1 Cross-References

- [Go API Definitions Index](api_go_defs_index.md) - Complete index of all Go API functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Multi-Layer Deduplication](api_deduplication.md) - Content deduplication strategies and processing levels
- [File Format Specifications](package_file_format.md) - .nvpk format structure and signature implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Package Compression API](api_package_compression.md) - Compression-specific concurrency patterns

## 1. File Streaming Interface

This section describes the file streaming interface for efficient file processing.

### 1.1 Purpose

Provides memory-efficient streaming for large files with compression and encryption support.

### 1.2 Core Types

This section describes core types used in the file streaming interface.

#### 1.2.1. FileStream Struct

```go
// FileStream provides streaming access to file data with buffering and progress tracking.
type FileStream struct {
    reader    io.Reader     // Underlying reader interface
    file      *os.File      // File handle for direct file access
    bufReader *bufio.Reader // Buffered reader for efficient reading
    fileSize  int64         // Total size of the file
    position  int64         // Current read position
    closed    bool          // Whether the stream is closed
    chunkSize int           // Size of read chunks
    config    *StreamConfig // Stream configuration
    fileEntry *FileEntry    // Associated FileEntry
}
```

#### 1.2.2. StreamConfig Struct

See [StreamConfig Structure](api_package_compression.md#214-streamconfig-structure) for the complete structure definition.

### 1.3 Key Methods

This section describes key methods of the file streaming interface.

#### 1.3.1 NewFileStream Function

```go
// NewFileStream creates a new file stream with the specified configuration.
func NewFileStream(ctx context.Context, reader io.Reader, config *StreamConfig, entry ...*FileEntry) *FileStream
```

#### 1.3.2 Core Operations

This section describes core streaming operations.

##### 1.3.2.1 FileStream.ReadChunk Method

```go
// Returns *PackageError on failure
func (s *FileStream) ReadChunk(ctx context.Context) ([]byte, error)
```

##### 1.3.2.2 FileStream.Seek Method

```go
// Returns *PackageError on failure
func (s *FileStream) Seek(ctx context.Context, offset int64) error
```

##### 1.3.2.3 FileStream.Close Method

```go
// Returns *PackageError on failure
func (s *FileStream) Close() error
```

#### 1.3.3 Status and Query Methods

This section describes status and query methods for file streams.

##### 1.3.3.1 FileStream.GetStats Method

```go
// GetStats returns statistics about the stream's read operations.
func (s *FileStream) GetStats() StreamStats
```

##### 1.3.3.2 FileStream.Size Method

```go
// Size returns the total size of the stream in bytes.
func (s *FileStream) Size() int64
```

##### 1.3.3.3 FileStream.Position Method

```go
// Position returns the current read position in the stream.
func (s *FileStream) Position() int64
```

##### 1.3.3.4 FileStream.IsClosed Method

```go
// IsClosed returns true if the stream has been closed.
func (s *FileStream) IsClosed() bool
```

##### 1.3.3.5 FileStream.Progress Method

```go
// Progress returns progress information about the stream read operation.
func (s *FileStream) Progress() (bytesRead int64, totalBytes int64, readSpeed int64, elapsed time.Duration)
```

##### 1.3.3.6 FileStream.EstimatedTimeRemaining Method

```go
// EstimatedTimeRemaining returns an estimate of the time remaining to complete the stream read.
func (s *FileStream) EstimatedTimeRemaining() time.Duration
```

#### 1.3.4 Standard Go Interface Methods

This section describes standard Go interface implementations for file streams.

##### 1.3.4.1 FileStream.Read Method

```go
// Read reads data from the stream into the provided buffer.
func (s *FileStream) Read(p []byte) (int, error)
```

##### 1.3.4.2 FileStream.ReadAt Method

```go
// ReadAt reads data from the stream at the specified offset.
func (s *FileStream) ReadAt(p []byte, off int64) (int, error)
```

### 1.4 Features

- **Chunked Reading**: Configurable chunk sizes for optimal performance
- **Buffer Pool Integration**: Reuses buffers to reduce memory allocation
- **Compression Support**: Handles compressed data transparently
- **Encryption Support**: Handles encrypted data transparently
- **Memory Management**: Configurable memory limits and pressure handling
- **Performance Monitoring**: Built-in read speed and statistics tracking

### 1.5 Additional Methods

This section describes additional methods for file streaming.

#### 1.5.1 Stream Information

This section describes stream information methods.

##### 1.5.1.1 Size Method Details

See [FileStream Size Method](#1332-filestreamsize-method) for the method signature.

##### 1.5.1.2 Position Method Details

See [FileStream Position Method](#1333-filestreamposition-method) for the method signature.

##### 1.5.1.3 IsClosed Method Details

See [FileStream IsClosed Method](#1334-filestreamisclosed-method) for the method signature.

##### 1.5.1.4 FileStream Information Purpose

Provides basic information about the stream state.

##### 1.5.1.5 Size Returns

Total size of the stream in bytes

##### 1.5.1.6 Position Returns

Current read position in the stream

##### 1.5.1.7 IsClosed Returns

Boolean indicating if stream is closed

##### 1.5.1.8 FileStream Information Example Usage

```go
stream := NewFileStream(ctx, reader, config)
defer stream.Close()

fmt.Printf("Stream size: %d bytes\n", stream.Size())
fmt.Printf("Current position: %d\n", stream.Position())

if stream.IsClosed() {
    fmt.Println("Stream is closed")
}
```

#### 1.5.2 Progress Monitoring

This section describes progress monitoring methods for file streams.

##### 1.5.2.1 Progress Method Details

See [FileStream Progress Method](#1335-filestreamprogress-method) for the method signature.

##### 1.5.2.2 EstimatedTimeRemaining Details

This section describes the EstimatedTimeRemaining method details.

##### 1.5.2.3 FileStream Progress Purpose

Provides progress monitoring and time estimation.

##### 1.5.2.4 Progress Returns

- `bytesRead`: Number of bytes read so far
- `totalBytes`: Total bytes to read
- `readSpeed`: Current read speed in bytes per second
- `elapsed`: Time elapsed since stream creation

##### 1.5.2.5 EstimatedTimeRemaining Returns

Estimated time remaining for completion

##### 1.5.2.6 FileStream Progress Example Usage

```go
for {
    data, err := stream.ReadChunk(ctx)
    if err != nil {
        break
    }

    bytesRead, totalBytes, speed, elapsed := stream.Progress()
    remaining := stream.EstimatedTimeRemaining()

    fmt.Printf("Progress: %d/%d bytes (%.1f%%) - Speed: %d B/s - Remaining: %v\n",
        bytesRead, totalBytes, float64(bytesRead)/float64(totalBytes)*100, speed, remaining)
}
```

#### 1.5.3 Standard Go Interfaces

This section describes standard Go interface implementations for file streams.

##### 1.5.3.1 Read Method Details

See [FileStream Read Method](#1341-filestreamread-method) for the method signature.

##### 1.5.3.2 ReadAt Method Details

See [FileStream ReadAt Method](#1342-filestreamreadat-method) for the method signature.

##### 1.5.3.3 FileStream I/O Purpose

Implements standard Go I/O interfaces for compatibility.

##### 1.5.3.4 Read Parameters

- `p`: Buffer to read data into

##### 1.5.3.5 Read Returns

Number of bytes read and error

##### 1.5.3.6 ReadAt Parameters

- `p`: Buffer to read data into
- `off`: Offset to read from

##### 1.5.3.7 ReadAt Returns

Number of bytes read and error

##### 1.5.3.8 FileStream I/O Example Usage

```go
// Use as io.Reader
data := make([]byte, 1024)
n, err := stream.Read(data)
if err != nil {
    return fmt.Errorf("read failed: %w", err)
}

// Use as io.ReaderAt
data = make([]byte, 512)
n, err = stream.ReadAt(data, 1024) // Read 512 bytes starting at offset 1024
if err != nil {
    return fmt.Errorf("read at failed: %w", err)
}
```

## 2. Buffer Management System

This section describes the buffer management system for streaming operations.

### 2.1 Buffer Management Purpose

Provides intelligent buffer pooling with eviction policies and memory management.

### 2.2 Buffer Management Core Types

This section describes core types used in buffer management.

#### 2.2.1. BufferPool Struct

Buffer pool for type-safe buffer management of any type.

The generic `BufferPool[T]` uses the factory function pattern from [Generic Factory Functions](api_generics.md#23-factory-functions).

##### 2.2.1.1 BufferPool Struct Type Definition

```go
// BufferPool manages buffers of any type
type BufferPool[T any] struct {
    buffers     map[int][]T
    lastUsed    map[int]time.Time
    accessCount map[int]int64
}
```

##### 2.2.1.2 NewBufferPool Function

See [NewBufferPool Function](api_generics.md#232-newbufferpool-function) for the complete function definition.

##### 2.2.1.3 BufferPool[T].Get Method

```go
// Get retrieves a buffer of the specified size from the pool.
func (bp *BufferPool[T]) Get(size int) []T
```

##### 2.2.1.4 BufferPool[T].Put Method

```go
// Put returns a buffer to the pool for reuse.
func (bp *BufferPool[T]) Put(buf []T)
```

#### 2.2.2. BufferConfig Struct

```go
// BufferConfig configures buffer pool behavior and limits.
type BufferConfig struct {
    MaxTotalSize     int64         // Maximum total size of all buffers
    MaxBufferSize    int           // Maximum size of a single buffer
    EvictionPolicy   string        // "lru" or "time" eviction policy
    EvictionTimeout  time.Duration // Time after which unused buffers are evicted
}
```

### 2.3 Buffer Management Key Methods

This section describes key methods for buffer management.

#### 2.3.1 BufferPool Creation and Core Operations

This section describes buffer pool creation and core operations.

##### 2.3.1.1 NewBufferPool Function (BufferPool Creation)

See [NewBufferPool Function](api_generics.md#232-newbufferpool-function) for the complete function definition.

##### 2.3.1.2 Get Method Details

See [BufferPool Get Method](#2213-bufferpooltget-method) for the method signature.

##### 2.3.1.3 Put Method Details

See [BufferPool Put Method](#2214-bufferpooltput-method) for the method signature.

##### 2.3.1.4 BufferPool[T].GetStats Method

```go
// GetStats returns statistics about buffer pool usage.
func (bp *BufferPool[T]) GetStats() BufferPoolStats
```

#### 2.3.2 BufferPool Management Methods

This section describes buffer pool management methods.

##### 2.3.2.1 BufferPool[T].TotalSize Method

```go
// Additional BufferPool methods
func (bp *BufferPool[T]) TotalSize() int64
```

##### 2.3.2.2 BufferPool[T].SetMaxTotalSize Method

```go
// SetMaxTotalSize sets the maximum total size for all buffers in the pool.
func (bp *BufferPool[T]) SetMaxTotalSize(maxSize int64)
```

### 2.4 Buffer Management Features

- **Size-Based Pools**: Separate pools for different buffer sizes
- **LRU Eviction**: Least Recently Used eviction policy
- **Time-Based Eviction**: Automatic cleanup of unused buffers
- **Memory Limits**: Configurable total memory usage limits
- **Access Tracking**: Statistics on buffer usage patterns
- **Thread Safety**: Concurrent access with proper synchronization

### 2.5 Buffer Management Additional Methods

This section describes additional buffer management methods.

#### 2.5.1 TotalSize Method Details

See [BufferPool TotalSize Method](#2321-bufferpoolttotalsize-method) for the method signature.

#### 2.5.2 SetMaxTotalSize Method Details

See [BufferPool SetMaxTotalSize Method](#2322-bufferpooltsetmaxtotalsize-method) for the method signature.

#### 2.5.3 BufferPool Management Purpose

Provides additional buffer pool management and monitoring capabilities.

#### 2.5.4 TotalSize Returns

Total size of all buffers currently in the pool

#### 2.5.5 SetMaxTotalSize Parameters

- `maxSize`: Maximum total size in bytes

#### 2.5.6 SetMaxTotalSize Behavior

- `TotalSize()` returns the current memory usage of the pool
- `SetMaxTotalSize()` dynamically adjusts the memory limit
- When limit is exceeded, eviction policies are triggered

#### 2.5.7 BufferPool Management Example Usage

```go
pool := NewBufferPool(config)

// Check current memory usage
currentSize := pool.TotalSize()
fmt.Printf("Current pool size: %d bytes\n", currentSize)

// Adjust memory limit
pool.SetMaxTotalSize(2 << 30) // Set to 2GB

// Monitor memory usage
for {
    size := pool.TotalSize()
    if size > 1.5<<30 { // 1.5GB
        fmt.Println("High memory usage detected")
    }
    time.Sleep(time.Second)
}
```

### 2.6 Default Configuration

This section describes default configuration for buffer management.

#### 2.6.1 DefaultBufferConfig Function

```go
// DefaultBufferConfig returns a buffer configuration with default values.
func DefaultBufferConfig() *BufferConfig
```

**DefaultBufferConfig** returns a BufferConfig with default values suitable for most use cases.

The function creates and returns a new BufferConfig instance with the following default settings:

- **MaxTotalSize**: 1 GB (1 << 30 bytes) - Maximum total size of all buffers in the pool
- **MaxBufferSize**: 1 MB (1 << 20 bytes) - Maximum size of a single buffer
- **EvictionPolicy**: "lru" (Least Recently Used) - Buffer eviction strategy for managing memory when limits are reached
- **EvictionTimeout**: 5 minutes - Time after which unused buffers are automatically evicted from the pool

These defaults provide a balance between memory efficiency and performance for typical streaming operations. Applications can create custom configurations using the BufferConfig structure directly if different values are needed.

## 3. Streaming Concurrency Patterns

This section describes concurrency patterns for streaming operations.

### 3.1 Streaming Concurrency Purpose

Provides concurrent streaming operations with thread safety and resource management for large file processing.

### 3.2 Streaming Concurrency Core Types

This section describes core types used in streaming concurrency.

#### 3.2.1. StreamingWorkerPool Struct

This section describes the StreamingWorkerPool structure.

##### 3.2.1.1 StreamingWorkerPool Structure

```go
// StreamingWorkerPool manages concurrent streaming workers
type StreamingWorkerPool struct {
    mu       sync.RWMutex
    workers  []*StreamingWorker
    workChan chan StreamingJob
    done     chan struct{}
    config   *StreamingConcurrencyConfig
}
```

##### 3.2.1.2 StreamingWorker Structure

```go
// StreamingWorker represents a single streaming worker
type StreamingWorker struct {
    mu       sync.RWMutex
    id       int
    workChan chan StreamingJob
    done     chan struct{}
    stream   *FileStream
}
```

##### 3.2.1.3 StreamingJob Structure

```go
// StreamingJob represents a unit of streaming work
type StreamingJob struct {
    ID       string
    Stream   *FileStream
    Result   chan StreamingResult
    Context  context.Context
    Priority int
}
```

##### 3.2.1.4 StreamingConcurrencyConfig Structure

```go
// StreamingConcurrencyConfig defines streaming-specific concurrency settings
type StreamingConcurrencyConfig struct {
    // Inherits from generic ConcurrencyConfig
    // See [Generic Function Patterns](api_generics.md#2-generic-function-patterns) for base configuration
    *ConcurrencyConfig

    // Streaming-specific settings
    MaxStreamsPerWorker int           // Maximum streams per worker (0 = no limit)
    StreamBufferSize    int           // Buffer size for stream operations
    ChunkProcessingMode ChunkMode     // How to process chunks concurrently
}
```

##### 3.2.1.5 ChunkMode Type

```go
// ChunkMode defines how chunks are processed concurrently
type ChunkMode int

const (
    ChunkModeSequential ChunkMode = iota // Process chunks sequentially
    ChunkModeParallel                    // Process chunks in parallel
    ChunkModeAdaptive                    // Adapt based on system load
)
```

### 3.3 Streaming Concurrency Key Methods

This section describes key methods for streaming concurrency.

#### 3.3.1 Worker Pool Lifecycle Methods

This section describes worker pool lifecycle methods.

##### 3.3.1.1 StreamingWorkerPool.Start Method

```go
// Start initializes and starts the streaming worker pool
// Returns *PackageError on failure
func (p *StreamingWorkerPool) Start(ctx context.Context) error
```

##### 3.3.1.2 StreamingWorkerPool.Stop Method

```go
// Stop gracefully shuts down the streaming worker pool
// Returns *PackageError on failure
func (p *StreamingWorkerPool) Stop(ctx context.Context) error
```

#### 3.3.2 Job Submission and Processing Methods

This section describes job submission and processing methods.

##### 3.3.2.1 StreamingWorkerPool.SubmitStreamingJob Method

```go
// SubmitStreamingJob submits a streaming job to the worker pool
// Returns *PackageError on failure
func (p *StreamingWorkerPool) SubmitStreamingJob(ctx context.Context, job StreamingJob) error
```

##### 3.3.2.2 ProcessStreamsConcurrently Function

```go
// ProcessStreamsConcurrently processes multiple streams concurrently
func ProcessStreamsConcurrently(ctx context.Context, streams []*FileStream, processor func(*FileStream) error, config *StreamingConcurrencyConfig) []error
```

#### 3.3.3 StreamingWorkerPool.GetStreamingStats Method

```go
// GetStreamingStats returns current streaming worker pool statistics
func (p *StreamingWorkerPool) GetStreamingStats() StreamingStats
```

### 3.4 Streaming Concurrency Features

- **Concurrent Stream Processing**: Multiple streams processed simultaneously
- **Thread-Safe Operations**: All streaming operations are thread-safe
- **Resource Management**: Intelligent resource allocation and cleanup
- **Adaptive Chunking**: Dynamic chunk size based on system load
- **Progress Tracking**: Real-time progress monitoring across workers

## 4. Streaming Configuration Patterns

This section describes configuration patterns for streaming operations.

### 4.1 Streaming Configuration Purpose

Provides streaming-specific configuration patterns that extend the generic configuration patterns for streaming operations.

### 4.2 Streaming Configuration Core Types

This section describes core types used in streaming configuration.

#### 4.2.1. StreamingConfig Struct

This section describes the StreamingConfig structure.

##### 4.2.1.1 StreamingConfig Structure

```go
// StreamingConfig extends Config for streaming-specific settings
type StreamingConfig struct {
    *Config[*FileStream]  // See [Core Generic Types](api_generics.md#1-core-generic-types) for base configuration

    // Streaming-specific settings
    StreamBufferSize    Option[int]           // Buffer size for stream operations (see [Option Type](api_generics.md#11-option-type))
    ChunkProcessingMode Option[ChunkMode]     // How to process chunks concurrently (see [Option Type](api_generics.md#11-option-type))
    MaxStreamsPerWorker Option[int]           // Maximum streams per worker (see [Option Type](api_generics.md#11-option-type))
    StreamTimeout       Option[time.Duration] // Timeout for stream operations (see [Option Type](api_generics.md#11-option-type))
}
```

##### 4.2.1.2 StreamingConfigBuilder Struct

```go
// StreamingConfigBuilder provides fluent configuration building for streaming
type StreamingConfigBuilder struct {
    config *StreamingConfig
}
```

##### 4.2.1.3 NewStreamingConfigBuilder Function

```go
// NewStreamingConfigBuilder creates a new streaming configuration builder.
func NewStreamingConfigBuilder() *StreamingConfigBuilder
```

##### 4.2.1.4 StreamingConfigBuilder.WithStreamBufferSize Method

```go
// WithStreamBufferSize sets the stream buffer size for the configuration.
func (b *StreamingConfigBuilder) WithStreamBufferSize(size int) *StreamingConfigBuilder
```

##### 4.2.1.5 StreamingConfigBuilder.WithChunkProcessingMode Method

```go
// WithChunkProcessingMode sets the chunk processing mode for the configuration.
func (b *StreamingConfigBuilder) WithChunkProcessingMode(mode ChunkMode) *StreamingConfigBuilder
```

##### 4.2.1.6 StreamingConfigBuilder.WithMaxStreamsPerWorker Method

```go
// WithMaxStreamsPerWorker sets the maximum number of streams per worker for the configuration.
func (b *StreamingConfigBuilder) WithMaxStreamsPerWorker(max int) *StreamingConfigBuilder
```

##### 4.2.1.7 StreamingConfigBuilder.WithStreamTimeout Method

```go
// WithStreamTimeout sets the stream timeout for the configuration.
func (b *StreamingConfigBuilder) WithStreamTimeout(timeout time.Duration) *StreamingConfigBuilder
```

##### 4.2.1.8 StreamingConfigBuilder.Build Method

```go
// Build constructs and returns the final streaming configuration.
func (b *StreamingConfigBuilder) Build() *StreamingConfig
```

### 4.3 Streaming Configuration Key Methods

This section describes key methods for streaming configuration.

#### 4.3.1 StreamingConfigDefaults Structure

```go
// StreamingConfigDefaults represents default streaming configuration values.
type StreamingConfigDefaults struct {
    StreamBufferSize    int
    ChunkProcessingMode ChunkMode
    MaxStreamsPerWorker int
    StreamTimeout       time.Duration
}
```

#### 4.3.2 CreateStreamingConfig Function

```go
// CreateStreamingConfig creates a streaming configuration with intelligent defaults
func CreateStreamingConfig() *StreamingConfig
```

#### 4.3.3 ValidateStreamingConfig Function

```go
// ValidateStreamingConfig validates streaming configuration settings
// Returns *PackageError on failure
func ValidateStreamingConfig(config *StreamingConfig) error
```

#### 4.3.4 GetStreamingConfigDefaults Function

```go
// GetStreamingConfigDefaults returns default streaming configuration values
func GetStreamingConfigDefaults() StreamingConfigDefaults
```
