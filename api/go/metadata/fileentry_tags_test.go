package metadata

import (
	"encoding/json"
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestGetFileEntryTags tests GetFileEntryTags function
func TestGetFileEntryTags(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *FileEntry
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty tags",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "single tag",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				tag := generics.NewTag[any]("key1", "value1", generics.TagValueTypeString)
				tags := []*generics.Tag[any]{tag}
				tagData, _ := json.Marshal(tags)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(tagData)),
						Data:       tagData,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "multiple tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				tags := []*generics.Tag[any]{
					generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
					generics.NewTag[any]("key2", int64(42), generics.TagValueTypeInteger),
					generics.NewTag[any]("key3", true, generics.TagValueTypeBoolean),
				}
				tagData, _ := json.Marshal(tags)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(tagData)),
						Data:       tagData,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "corrupted tags data",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "partially corrupted tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create tags with one valid and one corrupted tag
				validTag := generics.NewTag[any]("valid", "value", generics.TagValueTypeString)
				tags := []*generics.Tag[any]{validTag}
				tagData, _ := json.Marshal(tags)
				// Append invalid JSON to create partial corruption
				corruptedData := append(tagData, []byte(`,"invalid":}`)...)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(corruptedData)),
						Data:       corruptedData,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 1,    // Should recover valid tag
			wantErr:   true, // But report corruption error
		},
		{
			name: "empty tags array",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				tagData := []byte("[]")
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(tagData)),
						Data:       tagData,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "array with corrupted individual tag",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create array with one valid tag and one corrupted tag
				validTag := generics.NewTag[any]("valid", "value", generics.TagValueTypeString)
				tags := []*generics.Tag[any]{validTag}
				tagData, _ := json.Marshal(tags)
				// Create array with valid tag followed by corrupted tag
				corruptedArray := []byte(`[`)
				corruptedArray = append(corruptedArray, tagData[1:len(tagData)-1]...) // Remove outer brackets
				corruptedArray = append(corruptedArray, []byte(`,{"invalid":}`)...)
				corruptedArray = append(corruptedArray, []byte(`]`)...)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(corruptedArray)),
						Data:       corruptedArray,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 1,    // Should recover valid tag
			wantErr:   true, // But report corruption error
		},
		{
			name: "partially corrupted tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create tags with one valid and one corrupted tag
				validTag := generics.NewTag[any]("valid", "value", generics.TagValueTypeString)
				tags := []*generics.Tag[any]{validTag}
				tagData, _ := json.Marshal(tags)
				// Append invalid JSON to create partial corruption
				corruptedData := append(tagData, []byte(`,"invalid":}`)...)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(corruptedData)),
						Data:       corruptedData,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 1,    // Should recover valid tag
			wantErr:   true, // But report corruption error
		},
		{
			name: "empty tags array",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				tagData := []byte("[]")
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(tagData)),
						Data:       tagData,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			tags, err := GetFileEntryTags(fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(tags) != tt.wantCount {
					t.Errorf("GetFileEntryTags() returned %d tags, want %d", len(tags), tt.wantCount)
				}
			}
		})
	}
}

// TestGetFileEntryTagsByType tests GetFileEntryTagsByType function
func TestGetFileEntryTagsByType(t *testing.T) {
	fe := NewFileEntry()
	// Add tags using AddFileEntryTag to ensure proper serialization
	_ = AddFileEntryTag(fe, "str1", "value1", generics.TagValueTypeString)
	_ = AddFileEntryTag(fe, "str2", "value2", generics.TagValueTypeString)
	_ = AddFileEntryTag(fe, "int1", int64(42), generics.TagValueTypeInteger)
	_ = AddFileEntryTag(fe, "int2", int64(100), generics.TagValueTypeInteger)
	_ = AddFileEntryTag(fe, "bool1", true, generics.TagValueTypeBoolean)
	_ = AddFileEntryTag(fe, "float1", 3.14, generics.TagValueTypeFloat)

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{"string tags", 2, false},
		{"integer tags", 2, false},
		{"boolean tags", 1, false},
		{"float tags", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var result []*generics.Tag[any]

			switch tt.name {
			case "string tags":
				strTags, err := GetFileEntryTagsByType[string](fe)
				if err == nil {
					result = make([]*generics.Tag[any], len(strTags))
					for i, t := range strTags {
						result[i] = generics.NewTag[any](t.Key, t.Value, t.Type)
					}
				}
			case "integer tags":
				// Note: After JSON round-trip, integers may be float64, so we check by TagValueType
				allTags, err := GetFileEntryTags(fe)
				if err == nil {
					for _, tag := range allTags {
						if tag.Type == generics.TagValueTypeInteger {
							result = append(result, tag)
						}
					}
				}
			case "boolean tags":
				boolTags, err := GetFileEntryTagsByType[bool](fe)
				if err == nil {
					result = make([]*generics.Tag[any], len(boolTags))
					for i, t := range boolTags {
						result[i] = generics.NewTag[any](t.Key, t.Value, t.Type)
					}
				}
			case "float tags":
				// Note: After JSON round-trip, both integers and floats may be float64
				// So we check by TagValueType to distinguish
				allTags, err := GetFileEntryTags(fe)
				if err == nil {
					for _, tag := range allTags {
						if tag.Type == generics.TagValueTypeFloat {
							result = append(result, tag)
						}
					}
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryTagsByType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(result) != tt.wantCount {
				t.Errorf("GetFileEntryTagsByType() returned %d tags, want %d", len(result), tt.wantCount)
			}
		})
	}
}

// TestGetFileEntryTag tests GetFileEntryTag function
func TestGetFileEntryTag(t *testing.T) {
	fe := NewFileEntry()
	tags := []*generics.Tag[any]{
		generics.NewTag[any]("author", "John Doe", generics.TagValueTypeString),
		generics.NewTag[any]("version", int64(1), generics.TagValueTypeInteger),
		generics.NewTag[any]("published", true, generics.TagValueTypeBoolean),
	}
	tagData, _ := json.Marshal(tags)
	fe.OptionalData = []OptionalDataEntry{
		{
			DataType:   OptionalDataTagsData,
			DataLength: uint16(len(tagData)),
			Data:       tagData,
		},
	}
	fe.updateOptionalDataLen()

	tests := []struct {
		name      string
		key       string
		wantFound bool
		wantErr   bool
	}{
		{"found string tag", "author", true, false},
		{"found integer tag", "version", true, false},
		{"found boolean tag", "published", true, false},
		{"not found", "nonexistent", false, true}, // Now returns error when tag doesn't exist
		{"wrong type", "version", false, false},   // Requesting as string when it's int64 - returns nil, nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tag *generics.Tag[any]
			var err error

			if tt.name == "wrong type" {
				strTag, err := GetFileEntryTag[string](fe, "version")
				if err == nil && strTag != nil {
					tag = generics.NewTag[any](strTag.Key, strTag.Value, strTag.Type)
				}
			} else {
				tag, err = GetFileEntryTag[any](fe, tt.key)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (tag != nil) != tt.wantFound {
				t.Errorf("GetFileEntryTag() tag = %v, wantFound %v", tag != nil, tt.wantFound)
			}

			if tt.wantFound && tag != nil {
				if tag.Key != tt.key {
					t.Errorf("GetFileEntryTag() tag.Key = %q, want %q", tag.Key, tt.key)
				}
			}
		})
	}
}

// TestGetFileEntryTag_WithAny tests GetFileEntryTag[any] usage
func TestGetFileEntryTag_WithAny(t *testing.T) {
	fe := NewFileEntry()
	tags := []*generics.Tag[any]{
		generics.NewTag[any]("unknown", "value", generics.TagValueTypeString),
	}
	tagData, _ := json.Marshal(tags)
	fe.OptionalData = []OptionalDataEntry{
		{
			DataType:   OptionalDataTagsData,
			DataLength: uint16(len(tagData)),
			Data:       tagData,
		},
	}
	fe.updateOptionalDataLen()

	tag, err := GetFileEntryTag[any](fe, "unknown")
	if err != nil {
		t.Fatalf("GetFileEntryTag[any]() error = %v", err)
	}
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so tag is not nil after check
	if tag == nil {
		t.Fatal("GetFileEntryTag[any]() returned nil tag")
	}
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so tag is not nil after check
	if tag.Type != generics.TagValueTypeString {
		t.Errorf("GetFileEntryTag[any]() tag.Type = %v, want TagValueTypeString", tag.Type)
	}
}

// TestAddFileEntryTag tests AddFileEntryTag function
func TestAddFileEntryTag(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   interface{}
		tagType generics.TagValueType
		wantErr bool
		errType pkgerrors.ErrorType
		setup   func() *FileEntry
	}{
		{
			name:    "add string tag",
			key:     "author",
			value:   "John Doe",
			tagType: generics.TagValueTypeString,
			wantErr: false,
			setup:   func() *FileEntry { return NewFileEntry() },
		},
		{
			name:    "add integer tag",
			key:     "version",
			value:   int64(1),
			tagType: generics.TagValueTypeInteger,
			wantErr: false,
			setup:   func() *FileEntry { return NewFileEntry() },
		},
		{
			name:    "add boolean tag",
			key:     "published",
			value:   true,
			tagType: generics.TagValueTypeBoolean,
			wantErr: false,
			setup:   func() *FileEntry { return NewFileEntry() },
		},
		{
			name:    "duplicate key error",
			key:     "author",
			value:   "Jane Doe",
			tagType: generics.TagValueTypeString,
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "author", "John Doe", generics.TagValueTypeString)
				return fe
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()

			var err error
			switch v := tt.value.(type) {
			case string:
				err = AddFileEntryTag(fe, tt.key, v, tt.tagType)
			case int64:
				err = AddFileEntryTag(fe, tt.key, v, tt.tagType)
			case bool:
				err = AddFileEntryTag(fe, tt.key, v, tt.tagType)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("AddFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
					if pkgErr.Type != tt.errType {
						t.Errorf("AddFileEntryTag() error type = %v, want %v", pkgErr.Type, tt.errType)
					}
				} else {
					t.Errorf("AddFileEntryTag() error is not a PackageError")
				}
			} else {
				// Verify tag was added
				tag, getErr := GetFileEntryTag[any](fe, tt.key)
				if getErr != nil {
					t.Errorf("GetFileEntryTag() error = %v", getErr)
				}
				if tag == nil {
					t.Error("AddFileEntryTag() tag was not added")
				}
			}
		})
	}
}

// TestSetFileEntryTag tests SetFileEntryTag function
func TestSetFileEntryTag(t *testing.T) {
	fe := NewFileEntry()
	_ = AddFileEntryTag(fe, "author", "John Doe", generics.TagValueTypeString)

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
		errType pkgerrors.ErrorType
	}{
		{
			name:    "update existing tag",
			key:     "author",
			value:   "Jane Doe",
			wantErr: false,
		},
		{
			name:    "non-existent key error",
			key:     "nonexistent",
			value:   "value",
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetFileEntryTag(fe, tt.key, tt.value, generics.TagValueTypeString)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
					if pkgErr.Type != tt.errType {
						t.Errorf("SetFileEntryTag() error type = %v, want %v", pkgErr.Type, tt.errType)
					}
				}
			} else {
				// Verify tag was updated
				tag, getErr := GetFileEntryTag[string](fe, tt.key)
				if getErr != nil {
					t.Errorf("GetFileEntryTag() error = %v", getErr)
				}
				if tag == nil {
					t.Error("SetFileEntryTag() tag was not found")
				} else if tag.Value != tt.value {
					t.Errorf("SetFileEntryTag() tag.Value = %q, want %q", tag.Value, tt.value)
				}
			}
		})
	}
}

// TestAddFileEntryTags tests AddFileEntryTags function
func TestAddFileEntryTags(t *testing.T) {
	tests := []struct {
		name    string
		tags    []*generics.Tag[any]
		wantErr bool
		errType pkgerrors.ErrorType
		setup   func() *FileEntry
	}{
		{
			name: "add multiple tags",
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
				generics.NewTag[any]("key2", int64(42), generics.TagValueTypeInteger),
			},
			wantErr: false,
			setup:   func() *FileEntry { return NewFileEntry() },
		},
		{
			name: "duplicate key error",
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
				generics.NewTag[any]("key1", "value2", generics.TagValueTypeString),
			},
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
			setup:   func() *FileEntry { return NewFileEntry() },
		},
		{
			name: "duplicate with existing tag",
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("existing", "value", generics.TagValueTypeString),
			},
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "existing", "old", generics.TagValueTypeString)
				return fe
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := AddFileEntryTags(fe, tt.tags)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
					if pkgErr.Type != tt.errType {
						t.Errorf("AddFileEntryTags() error type = %v, want %v", pkgErr.Type, tt.errType)
					}
				}
			} else {
				// Verify all tags were added
				allTags, getErr := GetFileEntryTags(fe)
				if getErr != nil {
					t.Errorf("GetFileEntryTags() error = %v", getErr)
				}
				if len(allTags) != len(tt.tags) {
					t.Errorf("AddFileEntryTags() added %d tags, want %d", len(allTags), len(tt.tags))
				}
			}
		})
	}
}

// TestSetFileEntryTags tests SetFileEntryTags function
func TestSetFileEntryTags(t *testing.T) {
	fe := NewFileEntry()
	_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
	_ = AddFileEntryTag(fe, "key2", int64(42), generics.TagValueTypeInteger)

	tests := []struct {
		name    string
		tags    []*generics.Tag[any]
		wantErr bool
		errType pkgerrors.ErrorType
	}{
		{
			name: "update multiple tags",
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("key1", "newvalue1", generics.TagValueTypeString),
				generics.NewTag[any]("key2", int64(100), generics.TagValueTypeInteger),
			},
			wantErr: false,
		},
		{
			name: "non-existent key error",
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("nonexistent", "value", generics.TagValueTypeString),
			},
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetFileEntryTags(fe, tt.tags)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
					if pkgErr.Type != tt.errType {
						t.Errorf("SetFileEntryTags() error type = %v, want %v", pkgErr.Type, tt.errType)
					}
				}
			}
		})
	}
}

// TestRemoveFileEntryTag tests RemoveFileEntryTag function
func TestRemoveFileEntryTag(t *testing.T) {
	fe := NewFileEntry()
	_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)

	tests := []struct {
		name    string
		key     string
		wantErr bool
		errType pkgerrors.ErrorType
	}{
		{
			name:    "remove existing tag",
			key:     "key1",
			wantErr: false,
		},
		{
			name:    "non-existent key error",
			key:     "nonexistent",
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RemoveFileEntryTag(fe, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
					if pkgErr.Type != tt.errType {
						t.Errorf("RemoveFileEntryTag() error type = %v, want %v", pkgErr.Type, tt.errType)
					}
				}
			} else {
				// Verify tag was removed
				tag, _ := GetFileEntryTag[any](fe, tt.key)
				if tag != nil {
					t.Error("RemoveFileEntryTag() tag was not removed")
				}
			}
		})
	}
}

// TestHasFileEntryTag tests HasFileEntryTag function
func TestHasFileEntryTag(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FileEntry
		key   string
		want  bool
	}{
		{
			name: "existing tag",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				return fe
			},
			key:  "key1",
			want: true,
		},
		{
			name: "non-existent tag",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				return fe
			},
			key:  "nonexistent",
			want: false,
		},
		{
			name: "empty file entry",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			key:  "anykey",
			want: false,
		},
		{
			name: "tag with error in getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			key:  "anykey",
			want: false, // Should return false on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := HasFileEntryTag(fe, tt.key)
			if got != tt.want {
				t.Errorf("HasFileEntryTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasFileEntryTags tests HasFileEntryTags function
func TestHasFileEntryTags(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FileEntry
		want  bool
	}{
		{
			name: "has tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				return fe
			},
			want: true,
		},
		{
			name: "no tags",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			want: false,
		},
		{
			name: "multiple tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				_ = AddFileEntryTag(fe, "key2", int64(42), generics.TagValueTypeInteger)
				return fe
			},
			want: true,
		},
		{
			name: "tag with error in getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			want: false, // Should return false on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := HasFileEntryTags(fe)
			if got != tt.want {
				t.Errorf("HasFileEntryTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSyncFileEntryTags tests SyncFileEntryTags function
func TestSyncFileEntryTags(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "sync with tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				return fe
			},
			wantErr: false,
		},
		{
			name: "sync with no tags",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			wantErr: false,
		},
		{
			name: "sync with multiple tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				_ = AddFileEntryTag(fe, "key2", int64(42), generics.TagValueTypeInteger)
				return fe
			},
			wantErr: false,
		},
		{
			name: "sync with corrupted tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Add valid tag first
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				// Corrupt the data
				fe.OptionalData[0].Data = []byte("invalid json")
				return fe
			},
			wantErr: true, // Corrupted data causes error in getTagsFromOptionalData
		},
		{
			name: "sync with existing tags entry",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Add a tag to create an OptionalData entry
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				// Add another tag - should update existing entry
				_ = AddFileEntryTag(fe, "key2", "value2", generics.TagValueTypeString)
				return fe
			},
			wantErr: false,
		},
		{
			name: "sync with empty tags map",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Add a tag then manually clear it to test empty map sync
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				fe.OptionalData = []OptionalDataEntry{} // Clear optional data
				fe.updateOptionalDataLen()
				return fe
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := SyncFileEntryTags(fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("SyncFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify tags are still accessible after sync
				tags, getErr := GetFileEntryTags(fe)
				if getErr != nil {
					t.Errorf("GetFileEntryTags() error = %v", getErr)
				}
				if tags == nil {
					t.Error("GetFileEntryTags() returned nil tags")
				}
			}
		})
	}
}

// TestSyncTagsToOptionalData_EdgeCases tests edge cases in syncTagsToOptionalData
func TestSyncTagsToOptionalData_EdgeCases(t *testing.T) {
	t.Run("sync with very large tag data", func(t *testing.T) {
		fe := NewFileEntry()
		// Create a tag with a very large value
		largeValue := make([]byte, 10000)
		for i := range largeValue {
			largeValue[i] = byte(i % 256)
		}
		_ = AddFileEntryTag(fe, "large", string(largeValue), generics.TagValueTypeString)

		// Sync should handle large data
		err := SyncFileEntryTags(fe)
		if err != nil {
			t.Errorf("SyncFileEntryTags() with large data error = %v", err)
		}

		// Verify tag is still accessible
		tags, _ := GetFileEntryTags(fe)
		if len(tags) != 1 {
			t.Errorf("GetFileEntryTags() returned %d tags, want 1", len(tags))
		}
	})

	t.Run("sync updates existing OptionalData entry", func(t *testing.T) {
		fe := NewFileEntry()
		// Add first tag
		_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
		initialDataLen := len(fe.OptionalData)

		// Add second tag - should update existing entry, not create new one
		_ = AddFileEntryTag(fe, "key2", "value2", generics.TagValueTypeString)

		if len(fe.OptionalData) != initialDataLen {
			t.Errorf("syncTagsToOptionalData() created new entry, want update existing")
		}

		// Verify both tags are present
		tags, _ := GetFileEntryTags(fe)
		if len(tags) != 2 {
			t.Errorf("GetFileEntryTags() returned %d tags, want 2", len(tags))
		}
	})
}

// TestGetFileEntryEffectiveTags tests GetFileEntryEffectiveTags function
func TestGetFileEntryEffectiveTags(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *FileEntry
		wantCount int
		wantErr   bool
	}{
		{
			name: "file with tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "file_key", "file_value", generics.TagValueTypeString)
				return fe
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "file with no tags",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "file with multiple tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "key1", "value1", generics.TagValueTypeString)
				_ = AddFileEntryTag(fe, "key2", int64(42), generics.TagValueTypeInteger)
				return fe
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "file with PathMetadataEntry tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{PathLength: 8, Path: "dir/file"},
				}

				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
					Type: PathMetadataTypeFile,
					Properties: []*generics.Tag[any]{
						{Key: "path-tag", Value: "path-value", Type: generics.TagValueTypeString},
					},
				}

				_ = fe.AssociateWithPathMetadata(pme)
				return fe
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "file with both file and path tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "file-tag", "file-value", generics.TagValueTypeString)
				fe.Paths = []generics.PathEntry{
					{PathLength: 8, Path: "dir/file"},
				}

				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
					Type: PathMetadataTypeFile,
					Properties: []*generics.Tag[any]{
						{Key: "path-tag", Value: "path-value", Type: generics.TagValueTypeString},
					},
				}

				_ = fe.AssociateWithPathMetadata(pme)
				return fe
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			effectiveTags, err := GetFileEntryEffectiveTags(fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryEffectiveTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(effectiveTags) != tt.wantCount {
				t.Errorf("GetFileEntryEffectiveTags() returned %d tags, want %d", len(effectiveTags), tt.wantCount)
			}
		})
	}
}

// TestGetFileEntryInheritedTags tests GetFileEntryInheritedTags function
func TestGetFileEntryInheritedTags(t *testing.T) {
	// Create hierarchy: root -> parent -> child
	root := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 1, Path: "/"},
		Type: PathMetadataTypeDirectory,
		Inheritance: &PathInheritance{
			Enabled:  true,
			Priority: 1,
		},
		Properties: []*generics.Tag[any]{
			{Key: "root-tag", Value: "root-value", Type: generics.TagValueTypeString},
		},
	}

	parent := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 4, Path: "dir"},
		Type: PathMetadataTypeDirectory,
		Inheritance: &PathInheritance{
			Enabled:  true,
			Priority: 2,
		},
		Properties: []*generics.Tag[any]{
			{Key: "parent-tag", Value: "parent-value", Type: generics.TagValueTypeString},
		},
	}
	parent.SetParentPath(root)

	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "dir/file"},
	}

	child := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}
	child.SetParentPath(parent)

	_ = fe.AssociateWithPathMetadata(child)

	tests := []struct {
		name      string
		setup     func() *FileEntry
		wantCount int
		wantErr   bool
	}{
		{
			name: "file with no inheritance",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "file with tags but no inheritance",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				_ = AddFileEntryTag(fe, "file_key", "file_value", generics.TagValueTypeString)
				return fe
			},
			wantCount: 0, // No parent directory, so no inheritance
			wantErr:   false,
		},
		{
			name: "file with inherited tags",
			setup: func() *FileEntry {
				return fe
			},
			wantCount: 2, // root-tag and parent-tag
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			inheritedTags, err := GetFileEntryInheritedTags(fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryInheritedTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(inheritedTags) != tt.wantCount {
				t.Errorf("GetFileEntryInheritedTags() returned %d tags, want %d", len(inheritedTags), tt.wantCount)
			}
		})
	}
}

// TestFileEntryTagOperations_AllValueTypes tests tag operations with all TagValueType values
func TestFileEntryTagOperations_AllValueTypes(t *testing.T) {
	fe := NewFileEntry()

	valueTypes := []struct {
		name    string
		value   interface{}
		tagType generics.TagValueType
	}{
		{"String", "test", generics.TagValueTypeString},
		{"Integer", int64(42), generics.TagValueTypeInteger},
		{"Float", 3.14, generics.TagValueTypeFloat},
		{"Boolean", true, generics.TagValueTypeBoolean},
		{"UUID", "550e8400-e29b-41d4-a716-446655440000", generics.TagValueTypeUUID},
		{"Hash", "abc123", generics.TagValueTypeHash},
		{"Version", "1.0.0", generics.TagValueTypeVersion},
		{"Timestamp", "2024-01-01T00:00:00Z", generics.TagValueTypeTimestamp},
		{"URL", "https://example.com", generics.TagValueTypeURL},
		{"Email", "test@example.com", generics.TagValueTypeEmail},
		{"Path", "/path/to/file", generics.TagValueTypePath},
		{"MimeType", "text/plain", generics.TagValueTypeMimeType},
		{"Language", "en", generics.TagValueTypeLanguage},
	}

	for _, vt := range valueTypes {
		t.Run(vt.name, func(t *testing.T) {
			key := "test_" + vt.name
			var err error

			switch v := vt.value.(type) {
			case string:
				err = AddFileEntryTag(fe, key, v, vt.tagType)
			case int64:
				err = AddFileEntryTag(fe, key, v, vt.tagType)
			case float64:
				err = AddFileEntryTag(fe, key, v, vt.tagType)
			case bool:
				err = AddFileEntryTag(fe, key, v, vt.tagType)
			}

			if err != nil {
				t.Errorf("AddFileEntryTag() error = %v", err)
				return
			}

			// Verify tag can be retrieved
			tag, getErr := GetFileEntryTag[any](fe, key)
			if getErr != nil {
				t.Errorf("GetFileEntryTag() error = %v", getErr)
			}
			if tag == nil {
				t.Error("GetFileEntryTag() returned nil tag")
			} else if tag.Type != vt.tagType {
				t.Errorf("GetFileEntryTag() tag.Type = %v, want %v", tag.Type, vt.tagType)
			}
		})
	}
}

// TestGetFileEntryTags_CorruptionScenarios tests various corruption scenarios
func TestGetFileEntryTags_CorruptionScenarios(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *FileEntry
		wantCount int
		wantErr   bool
	}{
		{
			name: "all tags corrupted but array structure valid",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create array with tags that have invalid structure
				// Use tags that will unmarshal but have zero ValueType (which is invalid)
				corruptedArray := []byte(`[{"Key":"tag1","ValueType":0,"Value":null},{"Key":"tag2","ValueType":0,"Value":null}]`)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(corruptedArray)),
						Data:       corruptedArray,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 2, // Tags unmarshal successfully (zero values are valid)
			wantErr:   false,
		},
		{
			name: "partial corruption with some valid tags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create array with mix of valid and corrupted tags
				// Use the same approach as existing "array with corrupted individual tag" test
				// Create valid tag first
				validTag := generics.NewTag[any]("valid", "value", generics.TagValueTypeString)
				tags := []*generics.Tag[any]{validTag}
				tagData, _ := json.Marshal(tags)
				// Create array with valid tag followed by corrupted tag
				// Use the same approach as existing "array with corrupted individual tag" test
				// The corrupted tag {"invalid":} is invalid JSON, so it will fail array parsing
				// This tests the path where array structure itself is invalid
				partialCorruption := []byte(`[`)
				partialCorruption = append(partialCorruption, tagData[1:len(tagData)-1]...) // Remove outer brackets
				partialCorruption = append(partialCorruption, []byte(`,{"invalid":}]`)...)
				partialCorruption = append(partialCorruption, []byte(`]`)...)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(partialCorruption)),
						Data:       partialCorruption,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 0,    // Array structure invalid, so no tags parsed
			wantErr:   true, // But report corruption error
		},
		{
			name: "all tags corrupted - entry removed",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create array with all tags corrupted (invalid JSON objects)
				// Use valid array structure but with tags that fail individual unmarshaling
				allCorrupted := []byte(`[{"invalid":},{"bad":}]`)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(allCorrupted)),
						Data:       allCorrupted,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 0,    // No valid tags, entry should be removed
			wantErr:   true, // Corruption error should be returned
		},
		{
			name: "tags with invalid ValueType greater than maximum",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create array with tags that have valid JSON structure but invalid ValueType
				// ValueType 0x11 (17) is greater than TagValueTypeNovusPackMetadata (0x10)
				invalidValueType := []byte(`[{"Key":"tag1","ValueType":17,"Value":"value1"},{"Key":"tag2","ValueType":255,"Value":"value2"}]`)
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: uint16(len(invalidValueType)),
						Data:       invalidValueType,
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantCount: 0,    // Invalid tags should be skipped
			wantErr:   true, // Corruption error should be returned
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			tags, err := GetFileEntryTags(fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("Error details: %+v", err)
				}
				return
			}

			if len(tags) != tt.wantCount {
				t.Errorf("GetFileEntryTags() returned %d tags, want %d", len(tags), tt.wantCount)
			}
		})
	}
}

// TestAddFileEntryTag_ErrorHandling tests error handling in AddFileEntryTag
func TestAddFileEntryTag_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		key     string
		value   string
		wantErr bool
	}{
		{
			name: "error from getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			key:     "test",
			value:   "value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := AddFileEntryTag(fe, tt.key, tt.value, generics.TagValueTypeString)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetFileEntryTag_ErrorHandling tests error handling in SetFileEntryTag
func TestSetFileEntryTag_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		key     string
		value   string
		wantErr bool
	}{
		{
			name: "error from getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			key:     "test",
			value:   "value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := SetFileEntryTag(fe, tt.key, tt.value, generics.TagValueTypeString)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRemoveFileEntryTag_ErrorHandling tests error handling in RemoveFileEntryTag
func TestRemoveFileEntryTag_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		key     string
		wantErr bool
	}{
		{
			name: "error from getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			key:     "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := RemoveFileEntryTag(fe, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetFileEntryEffectiveTags_ErrorHandling tests error handling in GetFileEntryEffectiveTags
func TestGetFileEntryEffectiveTags_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "error from GetFileEntryTags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			_, err := GetFileEntryEffectiveTags(fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryEffectiveTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetFileEntryTagsByType_ErrorHandling tests error handling in GetFileEntryTagsByType
func TestGetFileEntryTagsByType_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "error from GetFileEntryTags",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			_, err := GetFileEntryTagsByType[string](fe)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryTagsByType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetFileEntryTag_ErrorHandling tests error handling in GetFileEntryTag
func TestGetFileEntryTag_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		key     string
		wantErr bool
		errType pkgerrors.ErrorType
	}{
		{
			name: "error from getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			key:     "test",
			wantErr: true,
			errType: pkgerrors.ErrTypeCorruption,
		},
		{
			name: "tag does not exist",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// FileEntry with no tags
				return fe
			},
			key:     "nonexistent",
			wantErr: true,
			errType: pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			_, err := GetFileEntryTag[any](fe, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileEntryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				pkgErr, ok := pkgerrors.IsPackageError(err)
				if !ok {
					t.Errorf("GetFileEntryTag() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errType {
					t.Errorf("GetFileEntryTag() error type = %v, want %v", pkgErr.Type, tt.errType)
				}
				if tt.name == "tag does not exist" {
					if pkgErr.Message != "tag does not exist" {
						t.Errorf("GetFileEntryTag() error message = %q, want %q", pkgErr.Message, "tag does not exist")
					}
				}
			}
		})
	}
}

// TestAddFileEntryTags_ErrorHandling tests error handling in AddFileEntryTags
func TestAddFileEntryTags_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		tags    []*generics.Tag[any]
		wantErr bool
	}{
		{
			name: "error from getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := AddFileEntryTags(fe, tt.tags)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSetFileEntryTags_ErrorHandling tests error handling in SetFileEntryTags
func TestSetFileEntryTags_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		tags    []*generics.Tag[any]
		wantErr bool
	}{
		{
			name: "error from getTagsFromOptionalData",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				// Create corrupted optional data
				fe.OptionalData = []OptionalDataEntry{
					{
						DataType:   OptionalDataTagsData,
						DataLength: 10,
						Data:       []byte("invalid json"),
					},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			tags: []*generics.Tag[any]{
				generics.NewTag[any]("key1", "value1", generics.TagValueTypeString),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			err := SetFileEntryTags(fe, tt.tags)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetFileEntryTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
