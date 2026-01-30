// This file implements the Tag[T] type representing type-safe tags with
// specific value types. It contains the Tag type definition, TagValueType
// constants, and NewTag constructor. This file should contain all code
// related to typed tags as specified in api_generics.md and api_file_management.md.
//
// Specification: api_generics.md: 1. Core Generic Types

package generics

// TagValueType represents the type of a tag value.
//
// Specification: api_file_mgmt_file_entry.md: 13. TagValueType Type
type TagValueType uint8

const (
	// Basic Types
	TagValueTypeString  TagValueType = 0x00 // String value
	TagValueTypeInteger TagValueType = 0x01 // 64-bit signed integer
	TagValueTypeFloat   TagValueType = 0x02 // 64-bit floating point number
	TagValueTypeBoolean TagValueType = 0x03 // Boolean value

	// Structured Data
	TagValueTypeJSON       TagValueType = 0x04 // JSON-encoded object or array
	TagValueTypeYAML       TagValueType = 0x05 // YAML-encoded data
	TagValueTypeStringList TagValueType = 0x06 // Comma-separated list of strings

	// Identifiers
	TagValueTypeUUID    TagValueType = 0x07 // UUID string
	TagValueTypeHash    TagValueType = 0x08 // Hash/checksum string
	TagValueTypeVersion TagValueType = 0x09 // Semantic version string

	// Time
	TagValueTypeTimestamp TagValueType = 0x0A // ISO8601 timestamp

	// Network/Communication
	TagValueTypeURL   TagValueType = 0x0B // URL string
	TagValueTypeEmail TagValueType = 0x0C // Email address

	// File System
	TagValueTypePath     TagValueType = 0x0D // File system path
	TagValueTypeMimeType TagValueType = 0x0E // MIME type string

	// Localization
	TagValueTypeLanguage TagValueType = 0x0F // Language code (ISO 639-1)

	// NovusPack Special Files
	TagValueTypeNovusPackMetadata TagValueType = 0x10 // NovusPack special metadata file reference
)

// Tag represents a type-safe tag with a specific value type.
//
// Tag[T] provides type-safe tag operations where T is the type of the tag value.
// All tags are stored and accessed as typed tags for compile-time type safety.
//
// Specification: api_file_mgmt_file_entry.md: 3.2 Tag Type
type Tag[T any] struct {
	// Key is the tag key (UTF-8 string)
	Key string

	// Value is the typed tag value
	Value T

	// Type is the tag value type identifier
	Type TagValueType
}

// NewTag creates a new typed tag with the specified key, value, and type.
//
// NewTag[T] provides a factory function for creating type-safe tags.
// The tag type must match the value type T (e.g., TagValueTypeString for string,
// TagValueTypeInteger for int64, etc.).
//
// Type Parameters:
//   - T: The type of the tag value
//
// Parameters:
//   - key: The tag key (UTF-8 string)
//   - value: The typed tag value
//   - tagType: The tag value type identifier
//
// Returns:
//   - *Tag[T]: A new typed tag instance
//
// Specification: api_file_mgmt_file_entry.md: 3.2 Tag Type
func NewTag[T any](key string, value T, tagType TagValueType) *Tag[T] {
	return &Tag[T]{
		Key:   key,
		Value: value,
		Type:  tagType,
	}
}

// GetValue returns the typed tag value.
//
// Returns:
//   - T: The typed tag value
func (t *Tag[T]) GetValue() T {
	return t.Value
}

// SetValue sets the typed tag value.
//
// Parameters:
//   - value: The new typed tag value
func (t *Tag[T]) SetValue(value T) {
	t.Value = value
}
