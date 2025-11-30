# NovusPack Technical Specifications - Streaming and Buffer Management API

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. File Streaming Interface](#1-file-streaming-interface)
  - [1.1 Purpose](#11-purpose)
  - [1.2 Core Types](#12-core-types)
  - [1.3 Key Methods](#13-key-methods)
  - [1.4 Features](#14-features)
  - [1.5 Additional Methods](#15-additional-methods)
- [2. Buffer Management System](#2-buffer-management-system)
  - [2.1 Purpose](#21-purpose)
  - [2.2 Core Types](#22-core-types)
  - [2.3 Key Methods](#23-key-methods)
  - [2.4 Features](#24-features)
  - [2.5 Additional Methods](#25-additional-methods)
  - [2.6 Default Configuration](#26-default-configuration)
- [3. Streaming Concurrency Patterns](#3-streaming-concurrency-patterns)
  - [3.1 Purpose](#31-purpose)
  - [3.2 Core Types](#32-core-types)
  - [3.3 Key Methods](#33-key-methods)
  - [3.4 Features](#34-features)
- [4. Streaming Configuration Patterns](#4-streaming-configuration-patterns)
  - [4.1 Purpose](#41-purpose)
  - [4.2 Core Types](#42-core-types)
  - [4.3 Key Methods](#43-key-methods)

---

## 0. Overview

This document defines the streaming and buffer management API for the NovusPack system, providing memory-efficient handling of large files with compression and encryption support.

### 0.1 Cross-References

- [API Signatures Index](api_func_signatures_index.md) - Complete index of all functions, types, and structures
- [Core Package Interface](api_core.md) - Package operations and compression
- [Multi-Layer Deduplication](api_deduplication.md) - Content deduplication strategies and processing levels
- [File Format Specifications](package_file_format.md) - .npk format structure and signature implementation
- [Generic Types and Patterns](api_generics.md) - Generic concurrency patterns and type-safe configuration
- [Package Compression API](api_package_compression.md) - Compression-specific concurrency patterns

## 1. File Streaming Interface

### 1.1 Purpose

Provides memory-efficient streaming for large files with compression and encryption support.

### 1.2 Core Types

#### 1.2.1 FileStream struct

```go
type FileStream struct {
    reader    io.Reader     // Underlying reader interface
    file      *os.File      // File handle for direct file access
    bufReader *bufio.Reader // Buffered reader for efficient reading
    fileSize  int64         // Total size of the file
    position  int64         // Current read position
    closed    bool          // Whether the stream is closed
    chunkSize int           // Size of read chunks
    config    *StreamConfig // Stream configuration
    fileEntry *FileEntry    // Associated file entry
}
```

#### 1.2.2 StreamConfig struct

```go
type StreamConfig struct {
    BufferSize    int   // Size of read buffer (0 = default)
    ChunkSize     int   // Size of each read chunk (0 = calculated)
    MemoryLimit   int64 // Maximum memory for buffering (0 = no limit)
    UseBufferPool bool  // Whether to use global buffer pool
    IsCompressed  bool  // Whether file is compressed
    IsEncrypted   bool  // Whether file is encrypted
}
```

### 1.3 Key Methods

```go
func NewFileStream(ctx context.Context, reader io.Reader, config *StreamConfig, entry ...*FileEntry) *FileStream
func (s *FileStream) ReadChunk(ctx context.Context) ([]byte, error)
func (s *FileStream) Seek(ctx context.Context, offset int64) error
func (s *FileStream) Close() error
func (s *FileStream) GetStats() StreamStats

// Additional FileStream methods
func (s *FileStream) Size() int64
func (s *FileStream) Position() int64
func (s *FileStream) IsClosed() bool
func (s *FileStream) Progress() (bytesRead int64, totalBytes int64, readSpeed int64, elapsed time.Duration)
func (s *FileStream) EstimatedTimeRemaining() time.Duration
func (s *FileStream) Read(p []byte) (int, error)
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

#### 1.5.1 Stream Information

```go
// Size returns the total size of the stream
func (s *FileStream) Size() int64

// Position returns the current position in the stream
func (s *FileStream) Position() int64

// IsClosed checks if the stream is closed
func (s *FileStream) IsClosed() bool
```

#### 1.5.1.1 Purpose

Provides basic information about the stream state.

#### 1.5.1.2 Size Returns

Total size of the stream in bytes

#### 1.5.1.3 Position Returns

Current read position in the stream

#### 1.5.1.4 IsClosed Returns

Boolean indicating if stream is closed

#### 1.5.1.5 Example Usage

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

```go
// Progress returns detailed progress information
func (s *FileStream) Progress() (bytesRead int64, totalBytes int64, readSpeed int64, elapsed time.Duration)

// EstimatedTimeRemaining estimates time remaining for completion
func (s *FileStream) EstimatedTimeRemaining() time.Duration
```

#### 1.5.2.1 Purpose

Provides progress monitoring and time estimation.

#### 1.5.2.2 Progress Returns

- `bytesRead`: Number of bytes read so far
- `totalBytes`: Total bytes to read
- `readSpeed`: Current read speed in bytes per second
- `elapsed`: Time elapsed since stream creation

#### 1.5.2.3 EstimatedTimeRemaining Returns

Estimated time remaining for completion

#### 1.5.2.4 Example Usage

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

```go
// Read implements io.Reader interface
func (s *FileStream) Read(p []byte) (int, error)

// ReadAt implements io.ReaderAt interface
func (s *FileStream) ReadAt(p []byte, off int64) (int, error)
```

#### 1.5.3.1 Purpose

Implements standard Go I/O interfaces for compatibility.

#### 1.5.3.2 Read Parameters

- `p`: Buffer to read data into

#### 1.5.3.3 Read Returns

Number of bytes read and error

#### 1.5.3.4 ReadAt Parameters

- `p`: Buffer to read data into
- `off`: Offset to read from

#### 1.5.3.5 ReadAt Returns

Number of bytes read and error

#### 1.5.3.6 Example Usage

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

### 2.1 Purpose

Provides intelligent buffer pooling with eviction policies and memory management.

### 2.2 Core Types

#### 2.2.1 BufferPool struct

```go
type BufferPool struct {
    buffers     map[int][]byte        // Map of buffer IDs to buffer data
    lastUsed    map[int]time.Time     // Last access time for each buffer
    accessCount map[int]int64         // Access count for each buffer
}
```

#### 2.2.2 BufferConfig struct

```go
type BufferConfig struct {
    MaxTotalSize     int64         // Maximum total size of all buffers
    MaxBufferSize    int           // Maximum size of a single buffer
    EvictionPolicy   string        // "lru" or "time" eviction policy
    EvictionTimeout  time.Duration // Time after which unused buffers are evicted
}
```

#### 2.2.3 BufferPool struct

Buffer pool for type-safe buffer management of any type.

```go
// BufferPool manages buffers of any type
type BufferPool[T any] struct {
    buffers     map[int][]T
    lastUsed    map[int]time.Time
    accessCount map[int]int64
}

func NewBufferPool[T any](config *BufferConfig) *BufferPool[T]
func (bp *BufferPool[T]) Get(size int) []T
func (bp *BufferPool[T]) Put(buf []T)
```

### 2.3 Key Methods

```go
func NewBufferPool(config *BufferConfig) *BufferPool
func (bp *BufferPool) Get(size int) []byte
func (bp *BufferPool) Put(buf []byte)
func (bp *BufferPool) GetStats() BufferPoolStats

// Additional BufferPool methods
func (bp *BufferPool) TotalSize() int64
func (bp *BufferPool) SetMaxTotalSize(maxSize int64)
```

### 2.4 Features

- **Size-Based Pools**: Separate pools for different buffer sizes
- **LRU Eviction**: Least Recently Used eviction policy
- **Time-Based Eviction**: Automatic cleanup of unused buffers
- **Memory Limits**: Configurable total memory usage limits
- **Access Tracking**: Statistics on buffer usage patterns
- **Thread Safety**: Concurrent access with proper synchronization

### 2.5 Additional Methods

```go
// TotalSize returns the total size of all buffers in the pool
func (bp *BufferPool) TotalSize() int64

// SetMaxTotalSize sets the maximum total size for the buffer pool
func (bp *BufferPool) SetMaxTotalSize(maxSize int64)
```

#### 2.5.1 Purpose

Provides additional buffer pool management and monitoring capabilities.

#### 2.5.2 TotalSize Returns

Total size of all buffers currently in the pool

#### 2.5.3 SetMaxTotalSize Parameters

- `maxSize`: Maximum total size in bytes

#### 2.5.4 Behavior

- `TotalSize()` returns the current memory usage of the pool
- `SetMaxTotalSize()` dynamically adjusts the memory limit
- When limit is exceeded, eviction policies are triggered

#### 2.5.5 Example Usage

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

#### 2.6.1 DefaultBufferConfig() *BufferConfig

```go
func DefaultBufferConfig() *BufferConfig {
    return &BufferConfig{
        MaxTotalSize:     1 << 30,           // 1GB maximum total buffer size
        MaxBufferSize:    1 << 20,           // 1MB maximum single buffer size
        EvictionPolicy:   "lru",             // Least Recently Used
        EvictionTimeout:  5 * time.Minute,   // 5 minutes for unused buffer cleanup
    }
}
```

## 3. Streaming Concurrency Patterns

### 3.1 Purpose

Provides concurrent streaming operations with thread safety and resource management for large file processing.

### 3.2 Core Types

#### 3.2.1 StreamingWorkerPool struct

```go
// StreamingWorkerPool manages concurrent streaming workers
type StreamingWorkerPool struct {
    mu       sync.RWMutex
    workers  []*StreamingWorker
    workChan chan StreamingJob
    done     chan struct{}
    config   *StreamingConcurrencyConfig
}

// StreamingWorker represents a single streaming worker
type StreamingWorker struct {
    mu       sync.RWMutex
    id       int
    workChan chan StreamingJob
    done     chan struct{}
    stream   *FileStream
}

// StreamingJob represents a unit of streaming work
type StreamingJob struct {
    ID       string
    Stream   *FileStream
    Result   chan StreamingResult
    Context  context.Context
    Priority int
}

// StreamingConcurrencyConfig defines streaming-specific concurrency settings
type StreamingConcurrencyConfig struct {
    // Inherits from generic ConcurrencyConfig
    *ConcurrencyConfig

    // Streaming-specific settings
    MaxStreamsPerWorker int           // Maximum streams per worker (0 = no limit)
    StreamBufferSize    int           // Buffer size for stream operations
    ChunkProcessingMode ChunkMode     // How to process chunks concurrently
}

// ChunkMode defines how chunks are processed concurrently
type ChunkMode int

const (
    ChunkModeSequential ChunkMode = iota // Process chunks sequentially
    ChunkModeParallel                    // Process chunks in parallel
    ChunkModeAdaptive                    // Adapt based on system load
)
```

### 3.3 Key Methods

```go
// Start initializes and starts the streaming worker pool
func (p *StreamingWorkerPool) Start(ctx context.Context) error

// Stop gracefully shuts down the streaming worker pool
func (p *StreamingWorkerPool) Stop(ctx context.Context) error

// SubmitStreamingJob submits a streaming job to the worker pool
func (p *StreamingWorkerPool) SubmitStreamingJob(ctx context.Context, job StreamingJob) error

// ProcessStreamsConcurrently processes multiple streams concurrently
func ProcessStreamsConcurrently(ctx context.Context, streams []*FileStream, processor func(*FileStream) error, config *StreamingConcurrencyConfig) []error

// GetStreamingStats returns current streaming worker pool statistics
func (p *StreamingWorkerPool) GetStreamingStats() StreamingStats
```

### 3.4 Features

- **Concurrent Stream Processing**: Multiple streams processed simultaneously
- **Thread-Safe Operations**: All streaming operations are thread-safe
- **Resource Management**: Intelligent resource allocation and cleanup
- **Adaptive Chunking**: Dynamic chunk size based on system load
- **Progress Tracking**: Real-time progress monitoring across workers

## 4. Streaming Configuration Patterns

### 4.1 Purpose

Provides streaming-specific configuration patterns that extend the generic configuration patterns for streaming operations.

### 4.2 Core Types

#### 4.2.1 StreamingConfig struct

```go
// StreamingConfig extends Config for streaming-specific settings
type StreamingConfig struct {
    *Config[*FileStream]

    // Streaming-specific settings
    StreamBufferSize    Option[int]           // Buffer size for stream operations
    ChunkProcessingMode Option[ChunkMode]     // How to process chunks concurrently
    MaxStreamsPerWorker Option[int]           // Maximum streams per worker
    StreamTimeout       Option[time.Duration] // Timeout for stream operations
}

// StreamingConfigBuilder provides fluent configuration building for streaming
type StreamingConfigBuilder struct {
    config *StreamingConfig
}

func NewStreamingConfigBuilder() *StreamingConfigBuilder
func (b *StreamingConfigBuilder) WithStreamBufferSize(size int) *StreamingConfigBuilder
func (b *StreamingConfigBuilder) WithChunkProcessingMode(mode ChunkMode) *StreamingConfigBuilder
func (b *StreamingConfigBuilder) WithMaxStreamsPerWorker(max int) *StreamingConfigBuilder
func (b *StreamingConfigBuilder) WithStreamTimeout(timeout time.Duration) *StreamingConfigBuilder
func (b *StreamingConfigBuilder) Build() *StreamingConfig
```

### 4.3 Key Methods

```go
// CreateStreamingConfig creates a streaming configuration with intelligent defaults
func CreateStreamingConfig() *StreamingConfig

// ValidateStreamingConfig validates streaming configuration settings
func ValidateStreamingConfig(config *StreamingConfig) error

// GetStreamingConfigDefaults returns default streaming configuration values
func GetStreamingConfigDefaults() map[string]interface{}
```
