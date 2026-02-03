package generics

import (
	"testing"
)

// TestConfigBuilder tests ConfigBuilder creation
func TestConfigBuilder(t *testing.T) {
	builder := NewConfigBuilder[string]()
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so builder is not nil after check
	if builder == nil {
		t.Fatal("NewConfigBuilder should not return nil")
	}
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so builder is not nil after check
	if builder.config == nil {
		t.Fatal("ConfigBuilder should have config")
	}
}

func assertConfigOptionSet[T comparable](t *testing.T, opt Option[T], expected T, fieldName string) {
	t.Helper()
	if !opt.IsSet() {
		t.Errorf("%s should be set", fieldName)
	}
	val, ok := opt.Get()
	if !ok {
		t.Errorf("%s should be retrievable", fieldName)
	}
	if val != expected {
		t.Errorf("%s should be %v, got %v", fieldName, expected, val)
	}
}

// TestConfigBuilder_WithChunkSize tests WithChunkSize
func TestConfigBuilder_WithChunkSize(t *testing.T) {
	builder := NewConfigBuilder[string]()
	builder.WithChunkSize(1024)
	assertConfigOptionSet(t, builder.Build().ChunkSize, 1024, "ChunkSize")
}

// TestConfigBuilder_WithMemoryUsage tests WithMemoryUsage
func TestConfigBuilder_WithMemoryUsage(t *testing.T) {
	builder := NewConfigBuilder[int]()
	builder.WithMemoryUsage(2 * 1024 * 1024)
	assertConfigOptionSet(t, builder.Build().MaxMemoryUsage, 2*1024*1024, "MaxMemoryUsage")
}

// TestConfigBuilder_WithCompressionLevel tests WithCompressionLevel
func TestConfigBuilder_WithCompressionLevel(t *testing.T) {
	builder := NewConfigBuilder[float64]()
	builder.WithCompressionLevel(5)
	assertConfigOptionSet(t, builder.Build().CompressionLevel, 5, "CompressionLevel")
}

// TestConfigBuilder_WithStrategy tests WithStrategy
func TestConfigBuilder_WithStrategy(t *testing.T) {
	strategy := &testStrategy{
		name:         "test-strategy",
		strategyType: "test",
	}

	builder := NewConfigBuilder[string]()
	builder.WithStrategy(strategy)

	config := builder.Build()
	if !config.Strategy.IsSet() {
		t.Error("Strategy should be set")
	}
	strategyVal, ok := config.Strategy.Get()
	if !ok {
		t.Error("Strategy should be retrievable")
	}
	if strategyVal != strategy {
		t.Error("Strategy should match the set strategy")
	}
}

// TestConfigBuilder_WithValidator tests WithValidator
func TestConfigBuilder_WithValidator(t *testing.T) {
	validator := &testValidator{shouldFail: false}

	builder := NewConfigBuilder[string]()
	builder.WithValidator(validator)

	config := builder.Build()
	if !config.Validator.IsSet() {
		t.Error("Validator should be set")
	}
	validatorVal, ok := config.Validator.Get()
	if !ok {
		t.Error("Validator should be retrievable")
	}
	if validatorVal != validator {
		t.Error("Validator should match the set validator")
	}
}

// TestConfigBuilder_FluentAPI tests method chaining
func TestConfigBuilder_FluentAPI(t *testing.T) {
	strategy := &testStrategy{
		name:         "test",
		strategyType: "test",
	}
	validator := &testValidator{shouldFail: false}

	config := NewConfigBuilder[string]().
		WithChunkSize(2048).
		WithMemoryUsage(4 * 1024 * 1024).
		WithCompressionLevel(7).
		WithStrategy(strategy).
		WithValidator(validator).
		Build()

	if !config.ChunkSize.IsSet() {
		t.Error("ChunkSize should be set")
	}
	if !config.MaxMemoryUsage.IsSet() {
		t.Error("MaxMemoryUsage should be set")
	}
	if !config.CompressionLevel.IsSet() {
		t.Error("CompressionLevel should be set")
	}
	if !config.Strategy.IsSet() {
		t.Error("Strategy should be set")
	}
	if !config.Validator.IsSet() {
		t.Error("Validator should be set")
	}

	chunkSize, _ := config.ChunkSize.Get()
	if chunkSize != 2048 {
		t.Errorf("ChunkSize should be 2048, got %d", chunkSize)
	}

	memoryUsage, _ := config.MaxMemoryUsage.Get()
	if memoryUsage != 4*1024*1024 {
		t.Errorf("MaxMemoryUsage should be 4MB, got %d", memoryUsage)
	}

	level, _ := config.CompressionLevel.Get()
	if level != 7 {
		t.Errorf("CompressionLevel should be 7, got %d", level)
	}
}

// TestConfigBuilder_Build tests Build method
func TestConfigBuilder_Build(t *testing.T) {
	builder := NewConfigBuilder[int]()
	config1 := builder.Build()
	config2 := builder.Build()

	// Build should return the same config instance
	if config1 != config2 {
		t.Error("Build should return the same config instance")
	}
}

// TestConfig_OptionalFields tests that Config fields are optional
func TestConfig_OptionalFields(t *testing.T) {
	config := NewConfigBuilder[string]().Build()

	// All fields should be unset initially
	if config.ChunkSize.IsSet() {
		t.Error("ChunkSize should not be set initially")
	}
	if config.MaxMemoryUsage.IsSet() {
		t.Error("MaxMemoryUsage should not be set initially")
	}
	if config.CompressionLevel.IsSet() {
		t.Error("CompressionLevel should not be set initially")
	}
	if config.Strategy.IsSet() {
		t.Error("Strategy should not be set initially")
	}
	if config.Validator.IsSet() {
		t.Error("Validator should not be set initially")
	}
}

// TestConfig_TypeParameter tests Config with different type parameters
func TestConfig_TypeParameter(t *testing.T) {
	// Test with string
	configStr := NewConfigBuilder[string]().
		WithChunkSize(1024).
		Build()
	if !configStr.ChunkSize.IsSet() {
		t.Error("Config[string] should work")
	}

	// Test with int
	configInt := NewConfigBuilder[int]().
		WithChunkSize(2048).
		Build()
	if !configInt.ChunkSize.IsSet() {
		t.Error("Config[int] should work")
	}

	// Test with custom type
	configCustom := NewConfigBuilder[CustomType]().
		WithChunkSize(4096).
		Build()
	if !configCustom.ChunkSize.IsSet() {
		t.Error("Config[CustomType] should work")
	}
}
