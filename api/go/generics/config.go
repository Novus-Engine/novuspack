// This file implements generic configuration patterns: Config[T] and
// ConfigBuilder[T]. It contains type-safe configuration management with
// Option-based fields and fluent builder interface. This file should contain
// all code related to generic configuration as specified in api_generics.md
// Section 1.10.
//
// Specification: api_generics.md: 1.10 Generic Configuration Patterns

// Package generics provides generic types and patterns for the NovusPack API.
package generics

// Config provides type-safe configuration for any data type.
//
// Config[T] uses Option fields for all configuration values, allowing
// optional configuration parameters with type safety.
//
// Type Parameters:
//   - T: The data type this configuration applies to
//
// Example:
//
//	config := NewConfigBuilder[string]().
//	    WithChunkSize(1024).
//	    WithMemoryUsage(1024 * 1024).
//	    WithCompressionLevel(5).
//	    Build()
type Config[T any] struct {
	ChunkSize        Option[int64]
	MaxMemoryUsage   Option[int64]
	CompressionLevel Option[int]
	Strategy         Option[Strategy[T, T]]
	Validator        Option[Validator[T]]
}

// ConfigBuilder provides fluent configuration building.
//
// ConfigBuilder[T] allows building Config[T] instances using a fluent API.
// All methods return the builder itself for method chaining.
//
// Type Parameters:
//   - T: The data type this configuration applies to
//
// Example:
//
//	builder := NewConfigBuilder[int]()
//	config := builder.
//	    WithChunkSize(2048).
//	    WithMemoryUsage(2 * 1024 * 1024).
//	    Build()
type ConfigBuilder[T any] struct {
	config *Config[T]
}

// NewConfigBuilder creates a new ConfigBuilder for the given type.
func NewConfigBuilder[T any]() *ConfigBuilder[T] {
	return &ConfigBuilder[T]{
		config: &Config[T]{},
	}
}

// WithChunkSize sets the chunk size configuration option.
func (b *ConfigBuilder[T]) WithChunkSize(size int64) *ConfigBuilder[T] {
	b.config.ChunkSize.Set(size)
	return b
}

// WithMemoryUsage sets the maximum memory usage configuration option.
func (b *ConfigBuilder[T]) WithMemoryUsage(usage int64) *ConfigBuilder[T] {
	b.config.MaxMemoryUsage.Set(usage)
	return b
}

// WithCompressionLevel sets the compression level configuration option.
func (b *ConfigBuilder[T]) WithCompressionLevel(level int) *ConfigBuilder[T] {
	b.config.CompressionLevel.Set(level)
	return b
}

// WithStrategy sets the strategy configuration option.
func (b *ConfigBuilder[T]) WithStrategy(strategy Strategy[T, T]) *ConfigBuilder[T] {
	b.config.Strategy.Set(strategy)
	return b
}

// WithValidator sets the validator configuration option.
func (b *ConfigBuilder[T]) WithValidator(validator Validator[T]) *ConfigBuilder[T] {
	b.config.Validator.Set(validator)
	return b
}

// Build creates and returns the configured Config instance.
func (b *ConfigBuilder[T]) Build() *Config[T] {
	return b.config
}
