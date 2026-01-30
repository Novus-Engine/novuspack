package generics

import (
	"testing"
)

// TestNewTag tests NewTag factory function
func TestNewTag(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		tagType TagValueType
	}{
		{"string tag", "key1", "value1", TagValueTypeString},
		{"integer tag", "key2", int64(42), TagValueTypeInteger},
		{"float tag", "key3", 3.14, TagValueTypeFloat},
		{"boolean tag", "key4", true, TagValueTypeBoolean},
		{"UUID tag", "key5", "550e8400-e29b-41d4-a716-446655440000", TagValueTypeUUID},
		{"hash tag", "key6", "abc123", TagValueTypeHash},
		{"version tag", "key7", "1.0.0", TagValueTypeVersion},
		{"timestamp tag", "key8", "2024-01-01T00:00:00Z", TagValueTypeTimestamp},
		{"URL tag", "key9", "https://example.com", TagValueTypeURL},
		{"email tag", "key10", "test@example.com", TagValueTypeEmail},
		{"path tag", "key11", "/path/to/file", TagValueTypePath},
		{"mime type tag", "key12", "text/plain", TagValueTypeMimeType},
		{"language tag", "key13", "en", TagValueTypeLanguage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tag interface{}

			switch v := tt.value.(type) {
			case string:
				tag = NewTag[string](tt.key, v, tt.tagType)
			case int64:
				tag = NewTag[int64](tt.key, v, tt.tagType)
			case float64:
				tag = NewTag[float64](tt.key, v, tt.tagType)
			case bool:
				tag = NewTag[bool](tt.key, v, tt.tagType)
			}

			if tag == nil {
				t.Fatal("NewTag() returned nil")
			}
			// Use type assertion to check tag properties
			switch tagVal := tag.(type) {
			case *Tag[string]:
				if tagVal.Key != tt.key {
					t.Errorf("NewTag() tag.Key = %q, want %q", tagVal.Key, tt.key)
				}
				if tagVal.Type != tt.tagType {
					t.Errorf("NewTag() tag.Type = %v, want %v", tagVal.Type, tt.tagType)
				}
			case *Tag[int64]:
				if tagVal.Key != tt.key {
					t.Errorf("NewTag() tag.Key = %q, want %q", tagVal.Key, tt.key)
				}
				if tagVal.Type != tt.tagType {
					t.Errorf("NewTag() tag.Type = %v, want %v", tagVal.Type, tt.tagType)
				}
			case *Tag[float64]:
				if tagVal.Key != tt.key {
					t.Errorf("NewTag() tag.Key = %q, want %q", tagVal.Key, tt.key)
				}
				if tagVal.Type != tt.tagType {
					t.Errorf("NewTag() tag.Type = %v, want %v", tagVal.Type, tt.tagType)
				}
			case *Tag[bool]:
				if tagVal.Key != tt.key {
					t.Errorf("NewTag() tag.Key = %q, want %q", tagVal.Key, tt.key)
				}
				if tagVal.Type != tt.tagType {
					t.Errorf("NewTag() tag.Type = %v, want %v", tagVal.Type, tt.tagType)
				}
			default:
				t.Errorf("NewTag() returned unexpected type: %T", tag)
			}
		})
	}
}

// TestTag_GetValue tests Tag.GetValue method
func TestTag_GetValue(t *testing.T) {
	tag := NewTag[string]("test_key", "test_value", TagValueTypeString)

	value := tag.GetValue()
	if value != "test_value" {
		t.Errorf("Tag.GetValue() = %q, want test_value", value)
	}
}

// TestTag_SetValue tests Tag.SetValue method
func TestTag_SetValue(t *testing.T) {
	tag := NewTag[string]("test_key", "old_value", TagValueTypeString)

	tag.SetValue("new_value")
	if tag.Value != "new_value" {
		t.Errorf("Tag.SetValue() tag.Value = %q, want new_value", tag.Value)
	}
}

// TestTag_TypeSafety tests type safety of Tag[T]
func TestTag_TypeSafety(t *testing.T) {
	// String tag
	strTag := NewTag[string]("str_key", "string_value", TagValueTypeString)
	if strTag.Value != "string_value" {
		t.Errorf("String tag value = %q, want string_value", strTag.Value)
	}

	// Integer tag
	intTag := NewTag[int64]("int_key", int64(42), TagValueTypeInteger)
	if intTag.Value != int64(42) {
		t.Errorf("Integer tag value = %d, want 42", intTag.Value)
	}

	// Boolean tag
	boolTag := NewTag[bool]("bool_key", true, TagValueTypeBoolean)
	if boolTag.Value != true {
		t.Errorf("Boolean tag value = %v, want true", boolTag.Value)
	}
}

// TestTagValueType_Constants tests all TagValueType constants
func TestTagValueType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		constant TagValueType
		expected TagValueType
	}{
		{"String", TagValueTypeString, 0x00},
		{"Integer", TagValueTypeInteger, 0x01},
		{"Float", TagValueTypeFloat, 0x02},
		{"Boolean", TagValueTypeBoolean, 0x03},
		{"JSON", TagValueTypeJSON, 0x04},
		{"YAML", TagValueTypeYAML, 0x05},
		{"StringList", TagValueTypeStringList, 0x06},
		{"UUID", TagValueTypeUUID, 0x07},
		{"Hash", TagValueTypeHash, 0x08},
		{"Version", TagValueTypeVersion, 0x09},
		{"Timestamp", TagValueTypeTimestamp, 0x0A},
		{"URL", TagValueTypeURL, 0x0B},
		{"Email", TagValueTypeEmail, 0x0C},
		{"Path", TagValueTypePath, 0x0D},
		{"MimeType", TagValueTypeMimeType, 0x0E},
		{"Language", TagValueTypeLanguage, 0x0F},
		{"NovusPackMetadata", TagValueTypeNovusPackMetadata, 0x10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("TagValueType%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
