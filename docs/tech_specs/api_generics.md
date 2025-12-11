# NovusPack Technical Specifications - Generic Types and Patterns

- [0. Overview](#0-overview)
  - [0.1 Cross-References](#01-cross-references)
- [1. Core Generic Types](#1-core-generic-types)
  - [1.1 Option Type](#11-option-type)
  - [1.2 Result Type](#12-result-type)
  - [1.3 Collection Interface](#13-collection-interface)
  - [1.4 Basic Data Structures](#14-basic-data-structures)
  - [1.5 Strategy Interface](#15-strategy-interface)
  - [1.6 Validator Interface](#16-validator-interface)
  - [1.7 Generic Concurrency Patterns](#17-generic-concurrency-patterns)
  - [1.8 Generic Concurrency Methods](#18-generic-concurrency-methods)
  - [1.9 Generic Configuration Patterns](#19-generic-configuration-patterns)
- [2. Generic Function Patterns](#2-generic-function-patterns)
  - [2.1 Collection Operations](#21-collection-operations)
  - [2.2 Validation Functions](#22-validation-functions)
  - [2.3 Factory Functions](#23-factory-functions)
- [3. Best Practices](#3-best-practices)
  - [3.1 Naming Conventions](#31-naming-conventions)
  - [3.2 Type Parameter Constraints](#32-type-parameter-constraints)
  - [3.3 Error Handling](#33-error-handling)
  - [3.4 Documentation](#34-documentation)
  - [3.5 Testing](#35-testing)

---

## 0. Overview

This document defines the generic types and patterns used throughout the NovusPack API.
It provides the technical specifications for type-safe generic implementations across all API modules.

### 0.1 Cross-References

- [Core Package Interface](api_core.md) - Generic types and interfaces
- [File Management API](api_file_management.md) - Generic tag operations and collections
- [Package Compression API](api_package_compression.md) - Generic strategy patterns and compression concurrency
- [Streaming and Buffer Management](api_streaming.md) - Generic buffer pools and streaming concurrency
- [Security Validation API](api_security.md) - Generic encryption patterns and type-safe security

## 1. Core Generic Types

### 1.1 Option Type

The Option type provides type-safe handling of optional values.

```go
// Option provides type-safe optional configuration values
type Option[T any] struct {
    value T
    set   bool
}

func (o *Option[T]) Set(value T)
func (o *Option[T]) Get() (T, bool)
func (o *Option[T]) GetOrDefault(defaultValue T) T
func (o *Option[T]) IsSet() bool
func (o *Option[T]) Clear()
```

#### 1.1.1 Option Type Usage Examples

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

```go
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](value T) Result[T]
func Err[T any](err error) Result[T]
func (r Result[T]) Unwrap() (T, error)
func (r Result[T]) IsOk() bool
func (r Result[T]) IsErr() bool
```

#### 1.2.1 Result Type Usage Examples

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

### 1.3 Collection Interface

Generic collection interface for type-safe collection operations.

```go
type Collection[T comparable] interface {
    Add(item T) error
    Remove(item T) error
    Contains(item T) bool
    Size() int
    Clear()
    ToSlice() []T
}

// FilterFunc is a predicate function for filtering
type FilterFunc[T any] func(T) bool

// MapFunc transforms items from one type to another
type MapFunc[T any, U any] func(T) U
```

### 1.4 Basic Data Structures

Core generic data structures for type-safe operations.

```go
// Map provides type-safe key-value storage
type Map[K comparable, V any] struct {
    data map[K]V
}

func (m *Map[K, V]) Set(key K, value V)
func (m *Map[K, V]) Get(key K) (V, bool)
func (m *Map[K, V]) Delete(key K)
func (m *Map[K, V]) Keys() []K
func (m *Map[K, V]) Values() []V
func (m *Map[K, V]) Size() int

// Set provides type-safe set operations
type Set[T comparable] struct {
    data map[T]bool
}

func (s *Set[T]) Add(item T)
func (s *Set[T]) Remove(item T)
func (s *Set[T]) Contains(item T) bool
func (s *Set[T]) Size() int
func (s *Set[T]) ToSlice() []T

// Writer provides type-safe writer operations
type Writer[T io.Writer] struct {
    writer T
}

func (w *Writer[T]) Write(data []byte) (int, error)
func (w *Writer[T]) WriteString(s string) (int, error)
func (w *Writer[T]) Flush() error

// Numeric constraint for numeric types
type Numeric interface {
    int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Sum[T Numeric](items []T) T
func Average[T Numeric](items []T) T
```

### 1.5 Strategy Interface

Generic strategy pattern for processing different data types.

```go
type Strategy[T any, U any] interface {
    Process(ctx context.Context, input T) (U, error)
    Name() string
    Type() string
}
```

### 1.6 Validator Interface

Generic validation system for type-safe validation.

```go
type Validator[T any] interface {
    Validate(ctx context.Context, value T) error
}

// ValidationRule represents a single validation rule
type ValidationRule[T any] struct {
    Name      string
    Predicate func(T) bool
    Message   string
}

func (r *ValidationRule[T]) Validate(value T) error
```

### 1.7 Generic Concurrency Patterns

Generic concurrency and thread safety patterns for type-safe concurrent operations.

```go
// WorkerPool manages concurrent workers for any data type
type WorkerPool[T any] struct {
    mu       sync.RWMutex
    workers  []*Worker[T]
    workChan chan Job[T]
    done     chan struct{}
    config   *ConcurrencyConfig
}

// Worker represents a single worker in the pool
type Worker[T any] struct {
    mu       sync.RWMutex
    id       int
    workChan chan Job[T]
    done     chan struct{}
    strategy Strategy[T, T]
}

// Job represents a unit of work for any data type
type Job[T any] struct {
    ID       string
    Data     T
    Result   chan Result[T]
    Context  context.Context
    Priority int
}

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

// ThreadSafetyMode defines the level of thread safety guarantees
type ThreadSafetyMode int

const (
    ThreadSafetyNone       ThreadSafetyMode = iota // No thread safety guarantees
    ThreadSafetyReadOnly                          // Read-only operations are safe
    ThreadSafetyConcurrent                        // Concurrent read/write operations
    ThreadSafetyFull                              // Full thread safety with synchronization
)
```

### 1.8 Generic Concurrency Methods

```go
// Start initializes and starts the worker pool
func (p *WorkerPool[T]) Start(ctx context.Context) error

// Stop gracefully shuts down the worker pool
func (p *WorkerPool[T]) Stop(ctx context.Context) error

// SubmitJob submits a job to the worker pool
func (p *WorkerPool[T]) SubmitJob(ctx context.Context, job Job[T]) error

// GetWorkerStats returns current worker pool statistics
func (p *WorkerPool[T]) GetWorkerStats() WorkerStats

// ProcessConcurrently processes multiple items concurrently using the worker pool
func ProcessConcurrently[T any](ctx context.Context, items []T, processor Strategy[T, T], config *ConcurrencyConfig) ([]Result[T], error)
```

### 1.9 Generic Configuration Patterns

Generic configuration patterns for type-safe configuration management.

```go
// Config provides type-safe configuration for any data type
type Config[T any] struct {
    ChunkSize        Option[int64]
    MaxMemoryUsage   Option[int64]
    CompressionLevel Option[int]
    Strategy         Option[Strategy[T, T]]
    Validator        Option[Validator[T]]
}

// ConfigBuilder provides fluent configuration building
type ConfigBuilder[T any] struct {
    config *Config[T]
}

func NewConfigBuilder[T any]() *ConfigBuilder[T]
func (b *ConfigBuilder[T]) WithChunkSize(size int64) *ConfigBuilder[T]
func (b *ConfigBuilder[T]) WithMemoryUsage(usage int64) *ConfigBuilder[T]
func (b *ConfigBuilder[T]) WithCompressionLevel(level int) *ConfigBuilder[T]
func (b *ConfigBuilder[T]) WithStrategy(strategy Strategy[T, T]) *ConfigBuilder[T]
func (b *ConfigBuilder[T]) Build() *Config[T]
```

## 2. Generic Function Patterns

### 2.1 Collection Operations

Generic functions for type-safe collection operations.

```go
// Filter returns items that match the predicate
func Filter[T any](items []T, predicate FilterFunc[T]) []T

// Map transforms items from one type to another
func Map[T any, U any](items []T, mapper MapFunc[T, U]) []U

// Find returns the first item that matches the predicate
func Find[T any](items []T, predicate FilterFunc[T]) (T, bool)

// Reduce applies a reducer function to accumulate a result
func Reduce[T any, U any](items []T, initial U, reducer func(U, T) U) U
```

### 2.2 Validation Functions

Generic functions for type-safe validation operations.

```go
// ValidateWith validates a single value using a validator
func ValidateWith[T any](ctx context.Context, value T, validator Validator[T]) error

// ValidateAll validates multiple values using a validator
func ValidateAll[T any](ctx context.Context, values []T, validator Validator[T]) []error

// ComposeValidators creates a validator that runs multiple validators
func ComposeValidators[T any](validators ...Validator[T]) Validator[T]
```

### 2.3 Factory Functions

Generic factory functions for creating type-safe instances.

```go
// NewTypedTag creates a new typed tag with the specified type
func NewTypedTag[T any](key string, value T, tagType TagValueType) *TypedTag[T]

// NewBufferPool creates a new buffer pool for the specified type
func NewBufferPool[T any](config *BufferConfig) *BufferPool[T]

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder[T any]() *ConfigBuilder[T]
```

## 3. Best Practices

### 3.1 Naming Conventions

- Use descriptive type parameter names: `T`, `U`, `V` for simple cases
- Use meaningful names for complex cases: `Key`, `Value`, `Element`
- Use simple, descriptive names without "Generic" prefix: `Option[T]`, `Config[T]`, `WorkerPool[T]`
- Let the generic syntax `[T any]` indicate genericity, not the name itself

### 3.2 Type Parameter Constraints

- Use the most restrictive constraint that works
- Prefer `comparable` over `any` when possible
- Use interface constraints for behavior requirements

### 3.3 Error Handling

- Use `Result[T]` for operations that may fail
- Use `Option[T]` for optional values
- Provide both generic and non-generic versions when needed

### 3.4 Documentation

- Document type parameter constraints
- Provide usage examples
- Explain when to use generic vs non-generic versions

### 3.5 Testing

- Test with multiple type instantiations
- Use type-specific test cases
- Verify compile-time type safety
