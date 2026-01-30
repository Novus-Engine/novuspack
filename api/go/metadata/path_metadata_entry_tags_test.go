package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestGetPathMetaTags tests the GetPathMetaTags function.
func TestGetPathMetaTags(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *PathMetadataEntry
		wantErr bool
		wantLen int
	}{
		{
			name: "nil PathMetadataEntry",
			setup: func() *PathMetadataEntry {
				return nil
			},
			wantErr: true,
			wantLen: 0,
		},
		{
			name: "empty tags",
			setup: func() *PathMetadataEntry {
				return &PathMetadataEntry{
					Properties: []*generics.Tag[any]{},
				}
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name: "with tags",
			setup: func() *PathMetadataEntry {
				return &PathMetadataEntry{
					Properties: []*generics.Tag[any]{
						{Key: "tag1", Value: "value1", Type: generics.TagValueTypeString},
						{Key: "tag2", Value: "value2", Type: generics.TagValueTypeString},
					},
				}
			},
			wantErr: false,
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pme := tt.setup()
			tags, err := GetPathMetaTags(pme)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPathMetaTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(tags) != tt.wantLen {
				t.Errorf("GetPathMetaTags() returned %d tags, want %d", len(tags), tt.wantLen)
			}
		})
	}
}

// TestGetPathMetaTagsByType tests the GetPathMetaTagsByType function.
func TestGetPathMetaTagsByType(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			{Key: "str-tag", Value: "string", Type: generics.TagValueTypeString},
			{Key: "int-tag", Value: int64(42), Type: generics.TagValueTypeInteger},
			{Key: "str-tag2", Value: "another", Type: generics.TagValueTypeString},
		},
	}

	// Test string tags
	strTags, err := GetPathMetaTagsByType[string](pme)
	if err != nil {
		t.Fatalf("GetPathMetaTagsByType[string]() error = %v", err)
	}

	if len(strTags) != 2 {
		t.Errorf("GetPathMetaTagsByType[string]() returned %d tags, want 2", len(strTags))
	}

	// Test int tags
	intTags, err := GetPathMetaTagsByType[int64](pme)
	if err != nil {
		t.Fatalf("GetPathMetaTagsByType[int64]() error = %v", err)
	}

	if len(intTags) != 1 {
		t.Errorf("GetPathMetaTagsByType[int64]() returned %d tags, want 1", len(intTags))
	}
}

// TestGetPathMetaTag tests the GetPathMetaTag function.
func TestGetPathMetaTag(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			{Key: "test-tag", Value: "test-value", Type: generics.TagValueTypeString},
		},
	}

	// Test found tag
	tag, err := GetPathMetaTag[string](pme, "test-tag")
	if err != nil {
		t.Fatalf("GetPathMetaTag() error = %v", err)
	}

	if tag == nil {
		t.Fatal("GetPathMetaTag() returned nil for existing tag")
	}

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so tag is not nil after check
	if tag.Key != "test-tag" {
		t.Errorf("GetPathMetaTag() tag.Key = %q, want %q", tag.Key, "test-tag")
	}

	// Test not found tag
	tag, err = GetPathMetaTag[string](pme, "nonexistent")
	if err != nil {
		t.Fatalf("GetPathMetaTag() for nonexistent error = %v", err)
	}

	if tag != nil {
		t.Error("GetPathMetaTag() for nonexistent should return nil")
	}

	// Test nil PathMetadataEntry
	tag, err = GetPathMetaTag[string](nil, "test-tag")
	if err == nil {
		t.Error("GetPathMetaTag() with nil PathMetadataEntry should return error")
	}

	if tag != nil {
		t.Error("GetPathMetaTag() with nil PathMetadataEntry should return nil tag")
	}
}

// TestAddPathMetaTag tests the AddPathMetaTag function.
func TestAddPathMetaTag(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{},
	}

	// Test add new tag
	err := AddPathMetaTag(pme, "new-tag", "new-value", generics.TagValueTypeString)
	if err != nil {
		t.Fatalf("AddPathMetaTag() error = %v", err)
	}

	if len(pme.Properties) != 1 {
		t.Errorf("AddPathMetaTag() Properties length = %d, want 1", len(pme.Properties))
	}

	// Test duplicate key
	err = AddPathMetaTag(pme, "new-tag", "another-value", generics.TagValueTypeString)
	if err == nil {
		t.Error("AddPathMetaTag() with duplicate key should return error")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Error("AddPathMetaTag() should return PackageError")
	}
}

// TestSetPathMetaTag tests the SetPathMetaTag function.
func TestSetPathMetaTag(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			{Key: "existing-tag", Value: "old-value", Type: generics.TagValueTypeString},
		},
	}

	// Test update existing tag
	err := SetPathMetaTag(pme, "existing-tag", "new-value", generics.TagValueTypeString)
	if err != nil {
		t.Fatalf("SetPathMetaTag() error = %v", err)
	}

	tag, _ := GetPathMetaTag[string](pme, "existing-tag")
	if tag == nil {
		t.Fatal("SetPathMetaTag() tag not found after update")
	}

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so tag is not nil after check
	if tag.Value != "new-value" {
		t.Errorf("SetPathMetaTag() tag.Value = %q, want %q", tag.Value, "new-value")
	}

	// Test non-existent key
	err = SetPathMetaTag(pme, "nonexistent", "value", generics.TagValueTypeString)
	if err == nil {
		t.Error("SetPathMetaTag() with non-existent key should return error")
	}
}

// TestAddPathMetaTags tests the AddPathMetaTags function.
func TestAddPathMetaTags(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{},
	}

	tags := []*generics.Tag[any]{
		{Key: "tag1", Value: "value1", Type: generics.TagValueTypeString},
		{Key: "tag2", Value: "value2", Type: generics.TagValueTypeString},
	}

	// Test add multiple tags
	err := AddPathMetaTags(pme, tags)
	if err != nil {
		t.Fatalf("AddPathMetaTags() error = %v", err)
	}

	if len(pme.Properties) != 2 {
		t.Errorf("AddPathMetaTags() Properties length = %d, want 2", len(pme.Properties))
	}

	// Test duplicate key
	err = AddPathMetaTags(pme, tags)
	if err == nil {
		t.Error("AddPathMetaTags() with duplicate keys should return error")
	}
}

// TestSetPathMetaTags tests the SetPathMetaTags function.
func TestSetPathMetaTags(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			{Key: "tag1", Value: "old1", Type: generics.TagValueTypeString},
			{Key: "tag2", Value: "old2", Type: generics.TagValueTypeString},
		},
	}

	tags := []*generics.Tag[any]{
		{Key: "tag1", Value: "new1", Type: generics.TagValueTypeString},
		{Key: "tag2", Value: "new2", Type: generics.TagValueTypeString},
	}

	// Test update multiple tags
	err := SetPathMetaTags(pme, tags)
	if err != nil {
		t.Fatalf("SetPathMetaTags() error = %v", err)
	}

	tag1, _ := GetPathMetaTag[string](pme, "tag1")
	//nolint:staticcheck // SA5011: false positive - checking nil before accessing Value
	if tag1 == nil || tag1.Value != "new1" {
		t.Error("SetPathMetaTags() did not update tag1 correctly")
	}

	// Test non-existent key
	tags2 := []*generics.Tag[any]{
		{Key: "nonexistent", Value: "value", Type: generics.TagValueTypeString},
	}

	err = SetPathMetaTags(pme, tags2)
	if err == nil {
		t.Error("SetPathMetaTags() with non-existent key should return error")
	}
}

// TestRemovePathMetaTag tests the RemovePathMetaTag function.
func TestRemovePathMetaTag(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			{Key: "tag1", Value: "value1", Type: generics.TagValueTypeString},
			{Key: "tag2", Value: "value2", Type: generics.TagValueTypeString},
		},
	}

	// Test remove existing tag
	err := RemovePathMetaTag(pme, "tag1")
	if err != nil {
		t.Fatalf("RemovePathMetaTag() error = %v", err)
	}

	if len(pme.Properties) != 1 {
		t.Errorf("RemovePathMetaTag() Properties length = %d, want 1", len(pme.Properties))
	}

	// Test remove non-existent tag (returns nil, no error)
	err = RemovePathMetaTag(pme, "nonexistent")
	if err != nil {
		t.Errorf("RemovePathMetaTag() with non-existent key error = %v, want nil", err)
	}
}

// TestHasPathMetaTag tests the HasPathMetaTag function.
func TestHasPathMetaTag(t *testing.T) {
	pme := &PathMetadataEntry{
		Properties: []*generics.Tag[any]{
			{Key: "test-tag", Value: "value", Type: generics.TagValueTypeString},
		},
	}

	// Test existing tag
	if !HasPathMetaTag(pme, "test-tag") {
		t.Error("HasPathMetaTag() for existing tag = false, want true")
	}

	// Test non-existent tag
	if HasPathMetaTag(pme, "nonexistent") {
		t.Error("HasPathMetaTag() for non-existent tag = true, want false")
	}

	// Test nil PathMetadataEntry
	if HasPathMetaTag(nil, "test-tag") {
		t.Error("HasPathMetaTag() with nil PathMetadataEntry = true, want false")
	}
}
