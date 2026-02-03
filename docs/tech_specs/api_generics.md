# NovusPack Technical Specifications - Generic Types and Patterns

- [0. Overview](#0-overview)
  - [0.1 Cross-References Overview](#01-cross-references-overview)
- [1. Core Generic Types](#1-core-generic-types)
  - [1.1 Option Type](#11-option-type)
    - [1.1.1 Option Type Definition](#111-option-struct)
    - [1.1.2 Option.Set Method](#112-optiontset-method)
    - [1.1.3 Option.Get Method](#113-optiontget-method)
    - [1.1.4 Option.GetOrDefault Method](#114-optiontgetordefault-method)
    - [1.1.5 Option.IsSet Method](#115-optiontisset-method)
    - [1.1.6 Option.Clear Method](#116-optiontclear-method)
    - [1.1.7 Option Type Usage Examples](#117-option-type-usage-examples)
  - [1.2 Result Type](#12-result-type)
    - [1.2.1 Result Type Definition](#121-result-struct)
    - [1.2.2 Ok Function](#122-ok-function)
    - [1.2.3 Err Function](#123-err-function)
    - [1.2.4 Result.Unwrap Method](#124-resulttunwrap-method)
    - [1.2.5 Result.IsOk Method](#125-resulttisok-method)
    - [1.2.6 Result.IsErr Method](#126-resulttiserr-method)
    - [1.2.7 Result Type Usage Examples](#127-result-type-usage-examples)
  - [1.3 PathEntry Type](#13-pathentry-type)
    - [1.3.1 PathEntry Structure](#131-pathentry-structure)
    - [1.3.2 Binary Format Specification](#132-binary-format-specification)
    - [1.3.3 Path Storage Rules](#133-path-storage-rules)
    - [1.3.4 Validation Rules](#134-validation-rules)
    - [1.3.5 Usage Context](#135-usage-context)
  - [1.4 Collection Operations](#14-collection-operations)
  - [1.5 Data Structure Operations](#15-data-structure-operations)
    - [1.5.1 Map Operations](#151-map-operations)
    - [1.5.2 Set Operations](#152-set-operations)
    - [1.5.3 Aggregation Operations](#153-aggregation-operations)
  - [1.6 Strategy Interface](#16-strategy-interface)
  - [1.7 Validator Interface](#17-validator-interface)
    - [1.7.1 Validator Interface Definition](#171-validator-interface-definition)
    - [1.7.2 ValidationRule Structure](#172-validationrule-structure)
  - [1.8 Generic Concurrency Patterns](#18-generic-concurrency-patterns)
    - [1.8.1 WorkerPool Structure](#181-workerpool-structure)
    - [1.8.2 Worker Structure](#182-worker-structure)
    - [1.8.3 Job Structure](#183-job-structure)
    - [1.8.4 ConcurrencyConfig Structure](#184-concurrencyconfig-structure)
    - [1.8.5 ThreadSafetyMode Type](#185-threadsafetymode-type)
  - [1.9 Generic Concurrency Methods](#19-generic-concurrency-methods)
    - [1.9.1 WorkerPool.Start Method](#191-workerpooltstart-method)
    - [1.9.2 WorkerPool.Stop Method](#192-workerpooltstop-method)
    - [1.9.3 WorkerPool SubmitJob Method](#193-workerpooltsubmitjob-method)
    - [1.9.4 WorkerPool.GetWorkerStats Method](#194-workerpooltgetworkerstats-method)
    - [1.9.5 ProcessConcurrently Function](#195-processconcurrently-function)
  - [1.10 Generic Configuration Patterns](#110-generic-configuration-patterns)
    - [1.10.1 Config Structure](#1101-config-structure)
    - [1.10.2 ConfigBuilder Structure](#1102-configbuilder-structure)
- [2. Generic Function Patterns](#2-generic-function-patterns)
  - [2.1 Collection Operations (Generic Function Patterns)](#21-collection-operations-generic-function-patterns)
    - [2.1.1 Common samber/lo Functions](#211-common-samberlo-functions)
    - [2.1.2 Example Usage](#212-example-usage)
  - [2.2 Validation Functions](#22-validation-functions)
    - [2.2.1 ValidateWith Function](#221-validatewith-function)
    - [2.2.2 ValidateAll Function](#222-validateall-function)
    - [2.2.3 ComposeValidators Function](#223-composevalidators-function)
  - [2.3 Factory Functions](#23-factory-functions)
    - [2.3.1 NewTag Function](#231-newtag-function)
    - [2.3.2 NewBufferPool Function](#232-newbufferpool-function)
    - [2.3.3 NewConfigBuilder Function (Streaming Configuration)](#233-newconfigbuilder-function-streaming-configuration)
    - [2.3.4 Cross-References (Generic Function Patterns)](#234-cross-references-generic-function-patterns)
- [3. Best Practices](#3-best-practices)
  - [3.1 Naming Conventions](#31-naming-conventions)
  - [3.2 Type Parameter Constraints](#32-type-parameter-constraints)
    - [3.2.1 When to Use any](#321-when-to-use-any)
    - [3.2.2 When to Use comparable](#322-when-to-use-comparable)
    - [3.2.3 When to Use Interface Constraints](#323-when-to-use-interface-constraints)
    - [3.2.4 Documenting Constraints](#324-documenting-constraints)
    - [3.2.5 When to Avoid any as a Concrete Type](#325-when-to-avoid-any-as-a-concrete-type)
  - [3.3 Error Handling](#33-error-handling)
    - [3.3.1 Generic Error Context Helpers](#331-generic-error-context-helpers)
    - [3.3.2 Error Context in Validation Functions](#332-error-context-in-validation-functions)
    - [3.3.3 Error Transformation](#333-error-transformation)
  - [3.4 Documentation](#34-documentation)
  - [3.5 Testing](#35-testing)
  - [3.6 Generic Function Usage](#36-generic-function-usage)
    - [3.6.1 Always Use Generics](#361-always-use-generics)
    - [3.6.2 Go Language Limitation: Generic Methods on Non-Generic Types](#362-go-language-limitation-generic-methods-on-non-generic-types)
    - [3.6.3 Example: Tag Operations](#363-example-tag-operations)

---

## 0. Overview

This document defines the generic types and patterns used throughout the NovusPack API.
It provides the technical specifications for type-safe generic implementations across all API modules.

**Note:** For collection operations (filtering, mapping, searching, aggregation), NovusPack uses `samber/lo` instead of custom generic functions.
See [samber/lo Usage Standards](../implementations/go/samber_lo_usage.md) for detailed guidelines on when and how to use `samber/lo` functions.

### 0.1 Cross-References Overview

- [samber/lo Usage Standards](../implementations/go/samber_lo_usage.md) - Guidelines for using samber/lo collection operations
- [Core Package Interface](api_core.md) - Generic types and interfaces
- [File Management API](api_file_mgmt_index.md) - Generic tag operations and collections
- [Package Compression API](api_package_compression.md) - Generic strategy patterns and compression concurrency
- [Streaming and Buffer Management](api_streaming.md) - Generic buffer pools and streaming concurrency
- [Security Validation API](api_security.md) - Generic encryption patterns and type-safe security

## 1. Core Generic Types

This section describes the core generic types used throughout the API.

### 1.1 Option Type

The Option type provides type-safe handling of optional values.

#### 1.1.1 Option Struct

```go
// Option provides type-safe optional configuration values
type Option[T any] struct {
    value T
    set   bool
}
```

#### 1.1.2 Option[T].Set Method

```go
// Set sets the option value.
func (o *Option[T]) Set(value T)
```

#### 1.1.3 Option[T].Get Method

```go
// Get returns the value and a boolean indicating if the value is set.
func (o *Option[T]) Get() (T, bool)
```

#### 1.1.4 Option[T].GetOrDefault Method

```go
// GetOrDefault returns the value if set, otherwise returns the default value.
func (o *Option[T]) GetOrDefault(defaultValue T) T
```

#### 1.1.5 Option[T].IsSet Method

```go
// IsSet returns true if the option has a value set.
func (o *Option[T]) IsSet() bool
```

#### 1.1.6 Option[T].Clear Method

```go
// Clear clears the option value.
func (o *Option[T]) Clear()
```

#### 1.1.7 Option Type Usage Examples

```go
// Create an option
var opt Option[string]
opt.Set("hello")

// Get value with default
value := opt.GetOrDefault("default")

// Check if set
if opt.IsSet() {
    val, _ := opt.Get()
    fmt.Println(val)
}
```

### 1.2 Result Type

The Result type provides type-safe handling of operations that may fail.

#### 1.2.1 Result Struct

```go
// Result represents a value that may be an error.
type Result[T any] struct {
    value T
    err   error
}
```

#### 1.2.2 Ok Function

```go
// Ok creates a successful Result with the given value.
func Ok[T any](value T) Result[T]
```

#### 1.2.3 Err Function

```go
// Err creates a failed Result with the given error.
func Err[T any](err error) Result[T]
```

#### 1.2.4 Result[T].Unwrap Method

```go
// Unwrap returns the value and error from the Result.
func (r Result[T]) Unwrap() (T, error)
```

#### 1.2.5 Result[T].IsOk Method

```go
// IsOk returns true if the Result contains a value (no error).
func (r Result[T]) IsOk() bool
```

#### 1.2.6 Result[T].IsErr Method

```go
// IsErr returns true if the Result contains an error.
func (r Result[T]) IsErr() bool
```

#### 1.2.7 Result Type Usage Examples

```go
// Success case
result := Ok("success")
if result.IsOk() {
    value, _ := result.Unwrap()
    fmt.Println(value)
}

// Error case
result := Err[string](errors.New("failed"))
if result.IsErr() {
    _, err := result.Unwrap()
    fmt.Println(err)
}
```

### 1.3 PathEntry Type

PathEntry represents a minimal file or directory path.

PathEntry is a shared type used by both FileEntry and PathMetadataEntry to represent path information.
It contains only the path string itself - no metadata, permissions, or symlink information.

Path metadata (permissions, timestamps, ownership, tags, symlinks) is stored separately in metadata structures
(see [Package Metadata API - Path Metadata System](api_metadata.md#8-pathmetadata-system)) to allow
the same file content to have different permissions at different paths.

#### 1.3.1 PathEntry Structure

This section describes the PathEntry structure for path management.

##### 1.3.1.1 PathEntry Struct

```go
// PathEntry represents a minimal file or directory path
type PathEntry struct {
    // PathLength is the length of the path in bytes (UTF-8)
    // Must match the actual length of the Path field
    PathLength uint16

    // Path is the UTF-8 encoded file or directory path (not null-terminated)
    // Must be valid UTF-8 and non-empty (after trimming whitespace)
    // All paths MUST be stored with a leading "/" to ensure full path references (see Package Path Semantics in api_core.md)
    // The root path "/" represents the package root
    // All paths use forward slashes (/) as separators, regardless of source platform
    Path string
}
```

##### 1.3.1.2 PathEntry.Validate Method

```go
// Validate performs validation checks on the PathEntry
// Returns error if PathLength doesn't match Path length, or if Path is empty/invalid
func (p *PathEntry) Validate() error
```

##### 1.3.1.3 PathEntry.Size Method

```go
// Size returns the total size of the PathEntry in bytes
// Formula: 2 (PathLength) + PathLength (Path)
func (p *PathEntry) Size() int
```

##### 1.3.1.4 PathEntry.ReadFrom Method

```go
// ReadFrom reads a PathEntry from the provided io.Reader
// Implements io.ReaderFrom interface
// Returns number of bytes read and any error encountered
func (p *PathEntry) ReadFrom(r io.Reader) (int64, error)
```

##### 1.3.1.5 PathEntry.WriteTo Method

```go
// WriteTo writes a PathEntry to the provided io.Writer
// Implements io.WriterTo interface
// Returns number of bytes written and any error encountered
func (p *PathEntry) WriteTo(w io.Writer) (int64, error)
```

##### 1.3.1.6 PathEntry.GetPath Method

```go
// GetPath returns the path string as stored (Unix-style with forward slashes)
func (p *PathEntry) GetPath() string
```

##### 1.3.1.7 PathEntry.GetPathForPlatform Method

```go
// GetPathForPlatform returns the path string converted for the specified platform
// On Windows, converts forward slashes to backslashes
// On Unix/Linux, returns the path as stored (with forward slashes)
func (p *PathEntry) GetPathForPlatform(isWindows bool) string
```

#### 1.3.2 Binary Format Specification

The binary format for PathEntry is a minimal variable-length structure containing only the path information.

**Field Layout** (in order):

| Field      | Size     | Type         | Description                                  |
| ---------- | -------- | ------------ | -------------------------------------------- |
| PathLength | 2 bytes  | uint16       | Length of path in bytes (UTF-8 encoded)      |
| Path       | variable | UTF-8 string | File or directory path (not null-terminated) |

**Total Size**: 2 + PathLength bytes

**Byte Order**: All multi-byte integers are stored in little-endian format.

##### Encoding Details

- **PathLength**: 2-byte unsigned integer (little-endian) specifying the number of bytes in the Path field
- **Path**: UTF-8 encoded string, exactly PathLength bytes, not null-terminated

**Path Metadata**: Permissions, timestamps, ownership, tags, and symlink information for paths are stored separately
in metadata structures (see [Package Metadata API - Path Metadata System](api_metadata.md#8-pathmetadata-system)). This design allows
the same file content to have different permissions and metadata at different paths.

#### 1.3.3 Path Storage Rules

For path storage rules, normalization, display, and extraction behavior, see [Package Path Semantics](api_core.md#2-package-path-semantics) in the Core Package Interface API.

The path storage rules cover:

- Path normalization rules (separator normalization, leading slash requirement, dot segment canonicalization, Unicode normalization, path length limits)
- Path normalization on storage
- Path display and extraction behavior
- Case sensitivity rules
- Path length validation utilities

#### 1.3.4 Validation Rules

The `Validate()` method performs the following checks:

1. **Path Length Consistency**: `PathLength` must exactly match the byte length of the `Path` field
2. **Path Non-Empty**: `Path` must not be empty or contain only whitespace (after trimming)
3. **UTF-8 Validity**: `Path` must be valid UTF-8 (enforced by Go's string type)
4. **Path Format**: `Path` must begin with `/` (leading slash is mandatory for all stored paths)
   - Exception: The root path `/` itself is valid (as a directory)
5. **Path Separators**: `Path` must use forward slashes (`/`) as separators, not backslashes (`\`)
6. **Unicode Normalization**: `Path` must be normalized to NFC form (see [Unicode Normalization](api_core.md#214-unicode-normalization))
7. **No Null Bytes**: `Path` must not contain null bytes (`\x00`)
8. **No Trailing Slash for Files**: File paths must not end with `/` (directories must end with `/`)

If any validation check fails, `Validate()` returns a `PackageError` with `ErrTypeValidation`.

#### 1.3.5 Usage Context

This section describes the usage context for generic types.

##### 1.3.5.1 In FileEntry

- `FileEntry.Paths` is a slice of `PathEntry` structures
- The first path entry (index 0) is considered the **primary path**
- Additional paths (index 1+) are **secondary paths** that point to the same content
- Multiple paths enable efficient storage of hard links and symbolic links
- Each path can have its own metadata (permissions, timestamps, ownership) stored in `PathMetadataEntry`
- All paths MUST be stored with a leading `/` (see [Package Path Semantics](api_core.md#2-package-path-semantics))

##### 1.3.5.2 In PathMetadataEntry

- `PathMetadataEntry.Path` is a string path
- Directory paths must end with `/` (forward slash)
- The path represents the path location in the package hierarchy

**Cross-Reference**: For information about how PathEntry is used within FileEntry structures in the package file format, see [Package File Format - Path Entries](package_file_format.md#4142-path-entries).

### 1.4 Collection Operations

**Use `samber/lo` for all collection operations.**

NovusPack uses `samber/lo` for collection operations including filtering, mapping, searching, aggregation, and duplicate detection.
No custom collection interface is needed.

For collection operation patterns and examples, see [samber/lo Usage Standards](../implementations/go/samber_lo_usage.md).

### 1.5 Data Structure Operations

**Use native Go data structures with `samber/lo` for operations.**

NovusPack uses native Go data structures (slices, maps) with `samber/lo` functions for operations.
No custom wrapper types are needed.

#### 1.5.1 Map Operations

Use native `map[K]V` with `samber/lo` functions:

```go
import "github.com/samber/lo"

// Extract keys
keys := lo.Keys(myMap)

// Extract values
values := lo.Values(myMap)

// Extract entries (key-value pairs)
entries := lo.Entries(myMap)

// Transform map entries
transformed := lo.MapEntries(myMap, func(k K, v V) (K2, V2) {
    return transformKey(k), transformValue(v)
})
```

#### 1.5.2 Set Operations

Use native `map[T]bool` or slices with `samber/lo` functions:

```go
import "github.com/samber/lo"

// Remove duplicates from slice
unique := lo.Uniq(items)

// Remove duplicates by key
unique := lo.UniqBy(items, func(item T) K { return item.Key() })

// Check if slice contains value
contains := lo.Contains(items, value)
```

#### 1.5.3 Aggregation Operations

Use `samber/lo` aggregation functions:

```go
import "github.com/samber/lo"

// Sum values
total := lo.SumBy(items, func(item T) int { return item.Size() })

// Count matching items
count := lo.CountBy(items, func(item T) bool { return item.IsValid() })

// Reduce/accumulate
result := lo.Reduce(items, func(acc U, item T, _ int) U {
    return accumulate(acc, item)
}, initialValue)
```

See [samber/lo Usage Standards](../implementations/go/samber_lo_usage.md) for more patterns and examples.

### 1.6 Strategy Interface

Generic strategy pattern for processing different data types.

```go
// Strategy defines a generic strategy pattern for processing different data types.
type Strategy[T any, U any] interface {
    Process(ctx context.Context, input T) (U, error)
    Name() string
    Type() string
}
```

### 1.7 Validator Interface

Generic validation system for type-safe validation.

#### 1.7.1 Validator Interface Definition

```go
// Validator defines a generic validation interface for type-safe validation.
type Validator[T any] interface {
    Validate(value T) error
}
```

#### 1.7.2 ValidationRule Structure

This section describes the ValidationRule structure for validation operations.

##### 1.7.2.1 ValidationRule Struct

```go
// ValidationRule represents a single validation rule
type ValidationRule[T any] struct {
    Name      string
    Predicate func(T) bool
    Message   string
}
```

##### 1.7.2.2 ValidationRule[T].Validate Method

```go
// Returns *PackageError on failure
func (r *ValidationRule[T]) Validate(value T) error
```

### 1.8 Generic Concurrency Patterns

Generic concurrency and thread safety patterns for type-safe concurrent operations.

#### 1.8.1 WorkerPool Structure

```go
// WorkerPool manages concurrent workers for any data type
type WorkerPool[T any] struct {
    mu       sync.RWMutex
    workers  []*Worker[T]
    workChan chan Job[T]
    done     chan struct{}
    config   *ConcurrencyConfig
}
```

#### 1.8.2 Worker Structure

```go
// Worker represents a single worker in the pool
type Worker[T any] struct {
    mu       sync.RWMutex
    id       int
    workChan chan Job[T]
    done     chan struct{}
    strategy Strategy[T, T]
}
```

#### 1.8.3 Job Structure

```go
// Job represents a unit of work for any data type
type Job[T any] struct {
    ID       string
    Data     T
    Result   chan Result[T]
    Context  context.Context
    Priority int
}
```

#### 1.8.4 ConcurrencyConfig Structure

```go
// ConcurrencyConfig defines thread safety and worker management settings
type ConcurrencyConfig struct {
    // Worker management
    MaxWorkers        int           // Maximum parallel workers (0 = auto-detect)
    WorkerTimeout     time.Duration // Worker timeout for graceful shutdown
    WorkerBufferSize  int           // Worker channel buffer size (0 = auto-calculate)

    // Thread safety
    UseMutex          bool          // Use mutex for shared state protection
    UseRWMutex        bool          // Use read-write mutex for better read performance
    LockTimeout       time.Duration // Lock acquisition timeout

    // Resource management
    MaxConcurrentOps  int           // Maximum concurrent operations (0 = no limit)
    ResourcePoolSize  int           // Resource pool size for workers
}
```

#### 1.8.5 ThreadSafetyMode Type

```go
// ThreadSafetyMode defines the level of thread safety guarantees
type ThreadSafetyMode int

const (
    ThreadSafetyNone       ThreadSafetyMode = iota // No thread safety guarantees
    ThreadSafetyReadOnly                          // Read-only operations are safe
    ThreadSafetyConcurrent                        // Concurrent read/write operations
    ThreadSafetyFull                              // Full thread safety with synchronization
)
```

### 1.9 Generic Concurrency Methods

This section describes generic concurrency methods and patterns.

#### 1.9.1 WorkerPool[T].Start Method

```go
// Start initializes and starts the worker pool
// Returns *PackageError on failure
func (p *WorkerPool[T]) Start(ctx context.Context) error
```

#### 1.9.2 WorkerPool[T].Stop Method

```go
// Stop gracefully shuts down the worker pool
// Returns *PackageError on failure
func (p *WorkerPool[T]) Stop(ctx context.Context) error
```

#### 1.9.3 WorkerPool[T].SubmitJob Method

```go
// SubmitJob submits a job to the worker pool
// Returns *PackageError on failure
func (p *WorkerPool[T]) SubmitJob(ctx context.Context, job Job[T]) error
```

#### 1.9.4 WorkerPool[T].GetWorkerStats Method

```go
// GetWorkerStats returns current worker pool statistics
func (p *WorkerPool[T]) GetWorkerStats() WorkerStats
```

#### 1.9.5 ProcessConcurrently Function

```go
// ProcessConcurrently processes multiple items concurrently using the worker pool
func ProcessConcurrently[T any](ctx context.Context, items []T, processor Strategy[T, T], config *ConcurrencyConfig) ([]Result[T], error)
```

### 1.10 Generic Configuration Patterns

Generic configuration patterns for type-safe configuration management.

#### 1.10.1 Config Structure

```go
// Config provides type-safe configuration for any data type
type Config[T any] struct {
    ChunkSize        Option[int64]
    MaxMemoryUsage   Option[int64]
    CompressionLevel Option[int]
    Strategy         Option[Strategy[T, T]]
    Validator        Option[Validator[T]]
}
```

#### 1.10.2 ConfigBuilder Structure

This section describes the ConfigBuilder structure for building configurations.

##### 1.10.2.1 ConfigBuilder Struct

```go
// ConfigBuilder provides fluent configuration building
type ConfigBuilder[T any] struct {
    config *Config[T]
}
```

##### 1.10.2.2 NewConfigBuilder Function

See [NewConfigBuilder Function (Streaming Configuration)](api_generics.md#233-newconfigbuilder-function-streaming-configuration) for the complete function definition.

##### 1.10.2.3 ConfigBuilder[T].WithChunkSize Method

```go
// WithChunkSize sets the chunk size for the configuration.
func (b *ConfigBuilder[T]) WithChunkSize(size int64) *ConfigBuilder[T]
```

##### 1.10.2.4 ConfigBuilder[T].WithMemoryUsage Method

```go
// WithMemoryUsage sets the memory usage limit for the configuration.
func (b *ConfigBuilder[T]) WithMemoryUsage(usage int64) *ConfigBuilder[T]
```

##### 1.10.2.5 ConfigBuilder[T].WithCompressionLevel Method

```go
// WithCompressionLevel sets the compression level for the configuration.
func (b *ConfigBuilder[T]) WithCompressionLevel(level int) *ConfigBuilder[T]
```

##### 1.10.2.6 ConfigBuilder[T].WithStrategy Method

```go
// WithStrategy sets the processing strategy for the configuration.
func (b *ConfigBuilder[T]) WithStrategy(strategy Strategy[T, T]) *ConfigBuilder[T]
```

##### 1.10.2.7 ConfigBuilder[T].WithValidator Method

```go
// WithValidator sets the validator for the configuration.
func (b *ConfigBuilder[T]) WithValidator(validator Validator[T]) *ConfigBuilder[T]
```

##### 1.10.2.8 ConfigBuilder[T].Build Method

```go
// Build constructs and returns the final configuration.
func (b *ConfigBuilder[T]) Build() *Config[T]
```

## 2. Generic Function Patterns

This section describes generic function patterns used throughout the API.

### 2.1 Collection Operations (Generic Function Patterns)

**Use `samber/lo` for all collection operations.**

NovusPack uses `samber/lo` for all collection operations including filtering, mapping, searching, and aggregation.
This provides better performance, wider adoption, and reduces maintenance burden.

See also [Collection Operations](#14-collection-operations) and [Data Structure Operations](#15-data-structure-operations) for additional patterns.

#### 2.1.1. Common Samber/lo Functions

Use these `samber/lo` functions for collection operations:

- **Filtering:** `lo.Filter()` or `lo.Reject()`
- **Mapping:** `lo.Map()`
- **Searching:** `lo.Find()`, `lo.Contains()`, or `lo.IndexOf()`
- **Aggregation:** `lo.Reduce()`, `lo.SumBy()`, or `lo.CountBy()`
- **Duplicate Detection:** `lo.UniqBy()` or `lo.FindDuplicates()`
- **Validation:** `lo.EveryBy()` or `lo.SomeBy()`

#### 2.1.2 Example Usage

```go
import "github.com/samber/lo"

// Filter items
filtered := lo.Filter(items, func(item T, _ int) bool {
    return item.IsValid()
})

// Map transformation
mapped := lo.Map(items, func(item T, _ int) U {
    return transform(item)
})

// Find element
found, ok := lo.Find(items, func(item T) bool {
    return item.Matches(criteria)
})

// Sum aggregation
total := lo.SumBy(items, func(item T) int {
    return item.Size()
})

// Reduce operation
result := lo.Reduce(items, func(acc U, item T, _ int) U {
    return accumulate(acc, item)
}, initialValue)
```

See [samber/lo Usage Standards](../implementations/go/samber_lo_usage.md) for detailed guidelines and patterns.

### 2.2 Validation Functions

Generic functions for type-safe validation operations.

#### 2.2.1 ValidateWith Function

```go
// ValidateWith validates a single value using a validator
// Returns *PackageError on failure
func ValidateWith[T any](ctx context.Context, value T, validator Validator[T]) error
```

#### 2.2.2 ValidateAll Function

```go
// ValidateAll validates multiple values using a validator
func ValidateAll[T any](ctx context.Context, values []T, validator Validator[T]) []error
```

#### 2.2.3 ComposeValidators Function

```go
// ComposeValidators creates a validator that runs multiple validators
func ComposeValidators[T any](validators ...Validator[T]) Validator[T]
```

**Error Handling**: Validation functions return errors using the structured error system.
Errors use `NewPackageError` or `WrapErrorWithContext` with typed error context structures.
See [Error Handling](#33-error-handling) for details on generic error context helpers.

### 2.3 Factory Functions

Generic factory functions for creating type-safe instances.

#### 2.3.1 NewTag Function

See [NewTag Function](api_file_mgmt_file_entry.md#192-newtag-function) for the complete function definition.

#### 2.3.2 NewBufferPool Function

```go
// NewBufferPool creates a new buffer pool for the specified type
func NewBufferPool[T any](config *BufferConfig) *BufferPool[T]
```

#### 2.3.3 NewConfigBuilder Function (Streaming Configuration)

```go
// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder[T any]() *ConfigBuilder[T]
```

#### 2.3.4 Cross-References (Generic Function Patterns)

- `NewTag`: Used in [FileEntry API - Tag Type](api_file_mgmt_file_entry.md#32-tag-type)
- `NewBufferPool`: Used in [Streaming and Buffer Management](api_streaming.md#2-buffer-management-system)
- `NewConfigBuilder`: Used throughout the API for configuration management

## 3. Best Practices

This section describes best practices for using generic types and patterns.

### 3.1 Naming Conventions

- Use descriptive type parameter names: `T`, `U`, `V` for simple cases
- Use meaningful names for complex cases: `Key`, `Value`, `Element`
- Use simple, descriptive names without "Generic" prefix: `Option[T]`, `Config[T]`, `WorkerPool[T]`
- Let the generic syntax `[T any]` indicate genericity, not the name itself

### 3.2 Type Parameter Constraints

- Use the most restrictive constraint that works
- Prefer `comparable` over `any` when possible
- Use interface constraints for behavior requirements

#### 3.2.1. When to Use Any

Use `any` (or no constraint) when defining a type parameter constraint and:

- The type parameter doesn't need any specific operations
- The type is only used for storage or passing through
- No comparison, ordering, or specific method calls are needed
- Examples: `Option[T any]`, `Result[T any]`, `Config[T any]`

Avoid using `any` as a concrete parameter type or return type in exported APIs.
Prefer typed structs, dedicated sum types, or generic functions over `any` inputs.

#### 3.2.2. When to Use Comparable

Use `comparable` when:

- The type needs to be used in map keys or compared with `==` or `!=`
- The type is used in set operations or duplicate detection
- Examples: Map key types, set element types, types used in comparison operations

**Note**: `GetTag[T]` uses `any` constraint, not `comparable`, because:

- The key is a `string` parameter, not type `T`
- The implementation only performs type assertion, not comparison
- Using `any` allows `GetTag[any]("key")` when the tag type is unknown

#### 3.2.3 When to Use Interface Constraints

Use interface constraints when:

- The type needs to implement specific methods
- Behavior requirements must be enforced at compile time
- Examples: `Strategy[T, U]` where T and U might need specific interfaces

##### 3.2.3.1 Serializable Interface

Example interface showing how interface constraints can be used:

```go
// Example interface: demonstrates interface constraint pattern
// This is NOT an actual API type - it's shown for illustration only
type Serializable interface {
    Serialize() ([]byte, error)
}
```

##### 3.2.3.2 Strategy Interface with Constraint

See [Strategy Interface](api_generics.md#16-strategy-interface) for the full definition.

The actual `Strategy` interface uses `T any` (unconstrained).

This example shows how you could add interface constraints if needed:

```go
// Example: A hypothetical constrained version of Strategy
// The actual Strategy interface uses T any (see line 558)
type ConstrainedStrategy[T Serializable, U any] interface {
    Process(ctx context.Context, input T) (U, error)
}
```

The canonical `Strategy` definition uses `T any` to allow maximum flexibility.

#### 3.2.4 Documenting Constraints

Always document why a specific constraint is chosen:

- Explain the operations that require the constraint
- Provide examples of valid and invalid type instantiations
- Note any performance or type safety implications

#### 3.2.5. When to Avoid Any As a Concrete Type

Avoid `any` as a concrete type in exported APIs unless the API is explicitly an untyped boundary.
Prefer compile-time type safety by representing variability with named types.

Common alternatives include:

- Use a typed struct for known fields and options.
- Use a dedicated sum type (tagged union) for "one of N" values.
- Use generic functions to keep callers typed and avoid runtime type assertions.
- Use `[]byte` or `json.RawMessage` for opaque payloads when the format is external and untyped.

### 3.3 Error Handling

- Use `Result[T]` for operations that may fail
- Use `Option[T]` for optional values
- Always use generic versions for type safety
- When using `samber/lo` functions, ensure error messages provide sufficient context (use explicit loops when index information is needed for debugging)

#### 3.3.1 Generic Error Context Helpers

When returning errors from generic functions, use the generic error context helpers from the structured error system:

- **`AddErrorContext[T any](err error, key string, value T) error`**: Add type-safe context to errors
- **`GetErrorContext[T any](err error, key string) (T, bool)`**: Retrieve type-safe context from errors
- **`NewPackageError[T any](errType ErrorType, message string, cause error, context T) *PackageError`**: Create structured errors with typed context
- **`WrapErrorWithContext[T any](err error, errType ErrorType, message string, context T) *PackageError`**: Wrap errors with typed context

For complete documentation, see [Structured Error System](api_core.md#10-structured-error-system).

#### 3.3.2 Error Context in Validation Functions

Validation functions should use typed error context when returning validation errors:

See [ValidationErrorContext Structure](api_signatures.md#534-validationerrorcontext-structure) for the complete structure definition.

Example usage:

```go
// Example: Validation error with typed context
// Note: ValidationErrorContext is defined in api_signatures.md
type ExampleValidationErrorContext struct {
    Field    string
    Value    interface{}
    Expected string
}

err := NewPackageError(ErrTypeValidation, "validation failed", nil, ExampleValidationErrorContext{
    Field:    "path",
    Value:    path,
    Expected: "non-empty string",
})
```

#### 3.3.3 Error Transformation

Use `MapError[T, U]` to transform error contexts when needed.
See [`MapError` Function](./api_core.md#1056-maperror-function).

Example context type for demonstrating error transformation:

```go
// Example context type: demonstrates error transformation pattern
// This is NOT an actual API type - it's shown for illustration only
type OldContext struct {
    Path string
}

// Example context type: demonstrates error transformation pattern
// This is NOT an actual API type - it's shown for illustration only
type NewContext struct {
    FilePath string
}

// Example: transform from OldContext to NewContext
transformedErr := MapError(err, func(old OldContext) NewContext {
    return NewContext{FilePath: old.Path}
})
```

### 3.4 Documentation

- Document type parameter constraints
- Provide usage examples
- Document generic function usage patterns and type safety benefits
- Reference `samber/lo` usage guidelines when documenting collection operations

### 3.5 Testing

- Test with multiple type instantiations
- Use type-specific test cases
- Verify compile-time type safety

### 3.6 Generic Function Usage

**Type safety is always a priority in the NovusPack API.**
All functions that can benefit from generics should use generic type parameters.
Since this is v1 of the API, there are no backward compatibility concerns.

#### 3.6.1 Always Use Generics

Use generic versions for:

- Type safety with compile-time checks
- Working with multiple types that share the same operations
- Avoiding type assertions or `interface{}` conversions
- Performance benefits from avoiding runtime type checks
- Examples: `GetFileEntryTag[T]`, `SetFileEntryTag[T]`, `EncryptFile[T]`, `GetFileEntryTagsByType[T]`

Even for single well-known types, prefer generics when they provide type safety benefits.
The performance overhead of generics in Go is minimal and the type safety benefits are significant.

#### 3.6.2 Go Language Limitation: Generic Methods on Non-Generic Types

**Important**: Go 1.25 and earlier versions do not support generic methods on non-generic receiver types.
This means that while you can define generic methods on generic types (e.g., `func (p *WorkerPool[T]) Start(...)`),
you cannot define generic methods on non-generic types (e.g., `func (fe *FileEntry) GetTag[T any](...)`).

##### 3.6.2.1 Impact on NovusPack API

- Tag operations on `FileEntry` and `PathMetadataEntry` are implemented as standalone generic functions rather than methods.
- Function names are prefixed with the type name to avoid confusion: `GetFileEntryTag[T]`, `GetPathMetaTag[T]`, etc.
- This is a known limitation documented in the API specifications and does not affect functionality.

##### 3.6.2.2 Not Supported Example

```go
// ❌ Not supported in Go 1.25:
func (fe *FileEntry) GetTag[T any](key string) (*Tag[T], error)
```

##### 3.6.2.3 Supported Example

✅ **Supported** - standalone generic function pattern:

See [GetFileEntryTag Function](api_file_mgmt_file_entry.md#3123-getfileentrytag-function) for the complete function definition.

The function signature follows the pattern: `func GetFileEntryTag[T any](fe *FileEntry, key string) (*Tag[T], error)`

##### 3.6.2.4 Future Compatibility

If future Go versions add support for generic methods on non-generic types, the NovusPack API may be updated to use methods.
However, the current function-based approach will remain supported for backward compatibility.

#### 3.6.3 Example: Tag Operations

All tag operations use typed tags exclusively.
The tag system provides both collection-level and individual tag operations.

##### 3.6.3.1 GetFileEntryTags Function

See [GetFileEntryTags Function](api_file_mgmt_file_entry.md#3121-getfileentrytags-function) for the complete function definition.

##### 3.6.3.2 GetFileEntryTagsByType Function

See [GetFileEntryTagsByType Function](api_file_mgmt_file_entry.md#3122-getfileentrytagsbytype-function) for the complete function definition.

##### 3.6.3.3 AddFileEntryTags Function

See [AddFileEntryTags Function](api_file_mgmt_file_entry.md#3124-addfileentrytags-function) for the complete function definition.

##### 3.6.3.4 SetFileEntryTags Function

See [SetFileEntryTags Function](api_file_mgmt_file_entry.md#3125-setfileentrytags-function) for the complete function definition.

##### 3.6.3.5 GetFileEntryTag Function

See [GetFileEntryTag Function](api_file_mgmt_file_entry.md#3123-getfileentrytag-function) for the complete function definition.

##### 3.6.3.6 AddFileEntryTag Function

See [AddFileEntryTag Function](api_file_mgmt_file_entry.md#3126-addfileentrytag-function) for the complete function definition.

##### 3.6.3.7 SetFileEntryTag Function

See [SetFileEntryTag Function](api_file_mgmt_file_entry.md#3127-setfileentrytag-function) for the complete function definition.

##### 3.6.3.8 GetPathMetaTags Function

See [GetPathMetaTags Function](api_metadata.md#8171-getpathmetatags-function) for the complete function definition.

##### 3.6.3.9 GetPathMetaTagsByType Function

See [GetPathMetaTagsByType Function](api_metadata.md#8172-getpathmetatagsbytype-function) for the complete function definition.

##### 3.6.3.10 GetPathMetaTag Function

See [GetPathMetaTag Function](api_metadata.md#8175-getpathmetatag-function) for the complete function definition.

##### 3.6.3.11 AddPathMetaTag Function

See [AddPathMetaTag Function](api_metadata.md#8176-addpathmetatag-function) for the complete function definition.

##### 3.6.3.12 SetPathMetaTag Function

See [SetPathMetaTag Function](api_metadata.md#8177-setpathmetatag-function) for the complete function definition.

##### 3.6.3.13 Usage Guidelines

- Use `GetFileEntryTags()` when you need all tags or are iterating over tags.
  Returns `([]*Tag[any], error)`, returns `*PackageError` on failure (corruption, I/O).
- Use `GetFileEntryTagsByType[T]()` when you need all tags of a specific type.
  Returns `([]*Tag[T], error)`, returns `*PackageError` on failure (corruption, I/O).
- Use `GetFileEntryTag[T]()` when you need type-safe access to a specific tag.
  Returns `(*Tag[T], error)` where `(nil, nil)` means the tag was not found (normal case) and `(nil, error)` means an underlying error occurred (corruption, I/O).
  If you don't know the tag type, use `GetFileEntryTag[any](fe, "key")` to retrieve the tag and inspect its `Type` field to determine the actual type.
- Use `AddFileEntryTag[T]()` when adding a new tag (will error if key already exists).
- Use `SetFileEntryTag[T]()` when updating an existing tag (will error if key does not exist).
- Use `AddFileEntryTags()` when adding multiple new tags at once (will error if any key already exists).
- Use `SetFileEntryTags()` when updating multiple existing tags at once (will error if any key does not exist).
