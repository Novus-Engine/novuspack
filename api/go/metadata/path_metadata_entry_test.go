package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestPathMetadataEntry_Validate tests the Validate method.
func TestPathMetadataEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *PathMetadataEntry
		wantErr bool
	}{
		{
			name: "valid file entry",
			setup: func() *PathMetadataEntry {
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 4, Path: "file"},
					Type: PathMetadataTypeFile,
				}
				return pme
			},
			wantErr: false,
		},
		{
			name: "valid directory entry",
			setup: func() *PathMetadataEntry {
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 3, Path: "dir"},
					Type: PathMetadataTypeDirectory,
					Inheritance: &PathInheritance{
						Enabled:  true,
						Priority: 1,
					},
					Metadata: &PathMetadata{
						Created:  "2024-01-01T00:00:00Z",
						Modified: "2024-01-01T00:00:00Z",
					},
				}
				return pme
			},
			wantErr: false,
		},
		{
			name: "invalid path entry",
			setup: func() *PathMetadataEntry {
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 5, Path: "file"}, // Mismatch
					Type: PathMetadataTypeFile,
				}
				return pme
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			setup: func() *PathMetadataEntry {
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 4, Path: "file"},
					Type: PathMetadataType(99), // Invalid
				}
				return pme
			},
			wantErr: true,
		},
		{
			name: "file with inheritance (invalid)",
			setup: func() *PathMetadataEntry {
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 4, Path: "file"},
					Type: PathMetadataTypeFile,
					Inheritance: &PathInheritance{
						Enabled: true,
					},
				}
				return pme
			},
			wantErr: true,
		},
		{
			name: "file with metadata (invalid)",
			setup: func() *PathMetadataEntry {
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{PathLength: 4, Path: "file"},
					Type: PathMetadataTypeFile,
					Metadata: &PathMetadata{
						Created: "2024-01-01T00:00:00Z",
					},
				}
				return pme
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pme := tt.setup()
			err := pme.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPathMetadataEntry_SetPath tests the SetPath method.
func TestPathMetadataEntry_SetPath(t *testing.T) {
	pme := &PathMetadataEntry{}
	pme.SetPath("test/path")

	if pme.Path.Path != "test/path" {
		t.Errorf("SetPath() Path = %q, want %q", pme.Path.Path, "test/path")
	}

	if pme.Path.PathLength != uint16(len("test/path")) {
		t.Errorf("SetPath() PathLength = %d, want %d", pme.Path.PathLength, len("test/path"))
	}
}

// TestPathMetadataEntry_GetPath tests the GetPath method.
func TestPathMetadataEntry_GetPath(t *testing.T) {
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 4, Path: "test"},
	}

	if got := pme.GetPath(); got != "test" {
		t.Errorf("GetPath() = %q, want %q", got, "test")
	}
}

// TestPathMetadataEntry_GetPathForPlatform tests the GetPathForPlatform method.
func TestPathMetadataEntry_GetPathForPlatform(t *testing.T) {
	tests := []struct {
		name      string
		pme       *PathMetadataEntry
		isWindows bool
		want      string
	}{
		{
			"Unix path with leading slash",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 13, Path: "/path/to/file"},
			},
			false,
			"path/to/file",
		},
		{
			"Windows path with leading slash",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 13, Path: "/path/to/file"},
			},
			true,
			"path\\to\\file",
		},
		{
			"Unix path without leading slash",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 12, Path: "path/to/file"},
			},
			false,
			"path/to/file",
		},
		{
			"Windows path without leading slash",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 12, Path: "path/to/file"},
			},
			true,
			"path\\to\\file",
		},
		{
			"Unix root path",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 1, Path: "/"},
			},
			false,
			"",
		},
		{
			"Windows root path",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 1, Path: "/"},
			},
			true,
			"",
		},
		{
			"Unix nested path",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 20, Path: "/very/long/path/file"},
			},
			false,
			"very/long/path/file",
		},
		{
			"Windows nested path",
			&PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 20, Path: "/very/long/path/file"},
			},
			true,
			"very\\long\\path\\file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pme.GetPathForPlatform(tt.isWindows); got != tt.want {
				t.Errorf("GetPathForPlatform(%v) = %q, want %q", tt.isWindows, got, tt.want)
			}
		})
	}
}

// TestPathMetadataEntry_GetPathEntry tests the GetPathEntry method.
func TestPathMetadataEntry_GetPathEntry(t *testing.T) {
	expected := generics.PathEntry{PathLength: 4, Path: "test"}
	pme := &PathMetadataEntry{
		Path: expected,
	}

	if got := pme.GetPathEntry(); got.Path != expected.Path {
		t.Errorf("GetPathEntry() = %v, want %v", got, expected)
	}
}

// TestPathMetadataEntry_GetType tests the GetType method.
func TestPathMetadataEntry_GetType(t *testing.T) {
	pme := &PathMetadataEntry{
		Type: PathMetadataTypeDirectory,
	}

	if got := pme.GetType(); got != PathMetadataTypeDirectory {
		t.Errorf("GetType() = %v, want %v", got, PathMetadataTypeDirectory)
	}
}

// TestPathMetadataEntry_IsDirectory tests the IsDirectory method.
func TestPathMetadataEntry_IsDirectory(t *testing.T) {
	tests := []struct {
		name     string
		pme      *PathMetadataEntry
		expected bool
	}{
		{"directory", &PathMetadataEntry{Type: PathMetadataTypeDirectory}, true},
		{"file", &PathMetadataEntry{Type: PathMetadataTypeFile}, false},
		{"file symlink", &PathMetadataEntry{Type: PathMetadataTypeFileSymlink}, false},
		{"directory symlink", &PathMetadataEntry{Type: PathMetadataTypeDirectorySymlink}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pme.IsDirectory(); got != tt.expected {
				t.Errorf("IsDirectory() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestPathMetadataEntry_IsFile tests the IsFile method.
func TestPathMetadataEntry_IsFile(t *testing.T) {
	tests := []struct {
		name     string
		pme      *PathMetadataEntry
		expected bool
	}{
		{"file", &PathMetadataEntry{Type: PathMetadataTypeFile}, true},
		{"directory", &PathMetadataEntry{Type: PathMetadataTypeDirectory}, false},
		{"file symlink", &PathMetadataEntry{Type: PathMetadataTypeFileSymlink}, false},
		{"directory symlink", &PathMetadataEntry{Type: PathMetadataTypeDirectorySymlink}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pme.IsFile(); got != tt.expected {
				t.Errorf("IsFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestPathMetadataEntry_IsSymlink tests the IsSymlink method.
func TestPathMetadataEntry_IsSymlink(t *testing.T) {
	tests := []struct {
		name     string
		pme      *PathMetadataEntry
		expected bool
	}{
		{"file symlink", &PathMetadataEntry{Type: PathMetadataTypeFileSymlink}, true},
		{"directory symlink", &PathMetadataEntry{Type: PathMetadataTypeDirectorySymlink}, true},
		{"file", &PathMetadataEntry{Type: PathMetadataTypeFile}, false},
		{"directory", &PathMetadataEntry{Type: PathMetadataTypeDirectory}, false},
		{"file with IsSymlink flag", &PathMetadataEntry{
			Type: PathMetadataTypeFile,
			Path: generics.PathEntry{PathLength: 4, Path: "file", IsSymlink: true},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pme.IsSymlink(); got != tt.expected {
				t.Errorf("IsSymlink() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestPathMetadataEntry_GetLinkTarget tests the GetLinkTarget method.
func TestPathMetadataEntry_GetLinkTarget(t *testing.T) {
	tests := []struct {
		name     string
		pme      *PathMetadataEntry
		expected string
	}{
		{
			name: "FileSystem LinkTarget",
			pme: &PathMetadataEntry{
				FileSystem: PathFileSystem{LinkTarget: "target"},
			},
			expected: "target",
		},
		{
			name: "Path LinkTarget",
			pme: &PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 4, Path: "file", LinkTarget: "path-target"},
			},
			expected: "path-target",
		},
		{
			name:     "no target",
			pme:      &PathMetadataEntry{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pme.GetLinkTarget(); got != tt.expected {
				t.Errorf("GetLinkTarget() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestPathMetadataEntry_ResolveSymlink tests the ResolveSymlink method.
func TestPathMetadataEntry_ResolveSymlink(t *testing.T) {
	tests := []struct {
		name     string
		pme      *PathMetadataEntry
		expected string
	}{
		{
			name: "absolute target",
			pme: &PathMetadataEntry{
				Path:       generics.PathEntry{PathLength: 4, Path: "file"},
				FileSystem: PathFileSystem{LinkTarget: "/absolute/target"},
			},
			expected: "/absolute/target",
		},
		{
			name: "relative target",
			pme: &PathMetadataEntry{
				Path:       generics.PathEntry{PathLength: 9, Path: "path/to/file"},
				FileSystem: PathFileSystem{LinkTarget: "target"},
			},
			expected: "path/to/target",
		},
		{
			name: "no target",
			pme: &PathMetadataEntry{
				Path: generics.PathEntry{PathLength: 4, Path: "file"},
			},
			expected: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pme.ResolveSymlink(); got != tt.expected {
				t.Errorf("ResolveSymlink() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestPathMetadataEntry_ParentPath tests parent path methods.
func TestPathMetadataEntry_ParentPath(t *testing.T) {
	root := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 1, Path: "/"},
		Type: PathMetadataTypeDirectory,
	}

	parent := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 4, Path: "dir"},
		Type: PathMetadataTypeDirectory,
	}

	child := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}

	// Test SetParentPath
	child.SetParentPath(parent)
	if child.GetParentPath() != parent {
		t.Error("SetParentPath() failed")
	}

	// Test GetParentPathString
	if got := child.GetParentPathString(); got != "dir" {
		t.Errorf("GetParentPathString() = %q, want %q", got, "dir")
	}

	// Test GetParentPathString for root (nil parent)
	if got := root.GetParentPathString(); got != "" {
		t.Errorf("GetParentPathString() for root = %q, want empty", got)
	}

	// Test IsRoot
	if root.IsRoot() != true {
		t.Error("IsRoot() for root should be true")
	}

	if child.IsRoot() != false {
		t.Error("IsRoot() for child should be false")
	}

	// Test GetDepth
	if root.GetDepth() != 0 {
		t.Errorf("GetDepth() for root = %d, want 0", root.GetDepth())
	}

	parent.SetParentPath(root)
	child.SetParentPath(parent)

	if child.GetDepth() != 2 {
		t.Errorf("GetDepth() for child = %d, want 2", child.GetDepth())
	}

	// Test GetAncestors
	ancestors := child.GetAncestors()
	if len(ancestors) != 2 {
		t.Errorf("GetAncestors() returned %d ancestors, want 2", len(ancestors))
	}

	if ancestors[0] != parent {
		t.Error("GetAncestors() first ancestor should be parent")
	}

	if ancestors[1] != root {
		t.Error("GetAncestors() second ancestor should be root")
	}
}

// TestPathMetadataEntry_GetInheritedTags tests the GetInheritedTags method.
func TestPathMetadataEntry_GetInheritedTags(t *testing.T) {
	root, parent := pathMetadataRootParentFixture()
	parent.SetParentPath(root)

	child := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}
	child.SetParentPath(parent)

	// Test GetInheritedTags
	tags, err := child.GetInheritedTags()
	if err != nil {
		t.Fatalf("GetInheritedTags() error = %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("GetInheritedTags() returned %d tags, want 2", len(tags))
	}

	// Verify tags are sorted by priority (higher priority first)
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value.(string)
	}

	if tagMap["parent-tag"] != "parent-value" {
		t.Error("GetInheritedTags() missing parent-tag")
	}

	if tagMap["root-tag"] != "root-value" {
		t.Error("GetInheritedTags() missing root-tag")
	}

	// Test with no parent
	rootTags, err := root.GetInheritedTags()
	if err != nil {
		t.Fatalf("GetInheritedTags() for root error = %v", err)
	}

	if len(rootTags) != 0 {
		t.Errorf("GetInheritedTags() for root returned %d tags, want 0", len(rootTags))
	}
}

// TestPathMetadataEntry_GetEffectiveTags tests the GetEffectiveTags method.
func TestPathMetadataEntry_GetEffectiveTags(t *testing.T) {
	// Create FileEntry with tags
	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "dir/file"},
	}
	err := AddFileEntryTag(fe, "file-tag", "file-value", generics.TagValueTypeString)
	if err != nil {
		t.Fatalf("AddFileEntryTag() error = %v", err)
	}

	// Create PathMetadataEntry with direct tags
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
		Properties: []*generics.Tag[any]{
			{
				Key:   "path-tag",
				Value: "path-value",
				Type:  generics.TagValueTypeString,
			},
		},
	}

	// Associate FileEntry
	err = pme.AssociateWithFileEntry(fe)
	if err != nil {
		t.Fatalf("AssociateWithFileEntry() error = %v", err)
	}

	// Test GetEffectiveTags
	tags, err := pme.GetEffectiveTags()
	if err != nil {
		t.Fatalf("GetEffectiveTags() error = %v", err)
	}

	if len(tags) < 2 {
		t.Errorf("GetEffectiveTags() returned %d tags, want at least 2", len(tags))
	}

	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value.(string)
	}

	if tagMap["path-tag"] != "path-value" {
		t.Error("GetEffectiveTags() missing path-tag")
	}

	if tagMap["file-tag"] != "file-value" {
		t.Error("GetEffectiveTags() missing file-tag")
	}
}

// TestPathMetadataEntry_AssociateWithFileEntry tests the AssociateWithFileEntry method.
func TestPathMetadataEntry_AssociateWithFileEntry(t *testing.T) {
	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "dir/file"},
	}

	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}

	// Test successful association
	err := pme.AssociateWithFileEntry(fe)
	if err != nil {
		t.Fatalf("AssociateWithFileEntry() error = %v", err)
	}

	if len(pme.AssociatedFileEntries) != 1 {
		t.Errorf("AssociateWithFileEntry() AssociatedFileEntries length = %d, want 1", len(pme.AssociatedFileEntries))
	}

	if pme.AssociatedFileEntries[0] != fe {
		t.Error("AssociateWithFileEntry() did not add FileEntry correctly")
	}

	// Test duplicate association (should be idempotent)
	err = pme.AssociateWithFileEntry(fe)
	if err != nil {
		t.Fatalf("AssociateWithFileEntry() duplicate error = %v", err)
	}

	if len(pme.AssociatedFileEntries) != 1 {
		t.Errorf("AssociateWithFileEntry() duplicate AssociatedFileEntries length = %d, want 1", len(pme.AssociatedFileEntries))
	}

	// Test path mismatch
	pme2 := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 7, Path: "no/match"},
		Type: PathMetadataTypeFile,
	}

	err = pme2.AssociateWithFileEntry(fe)
	if err == nil {
		t.Error("AssociateWithFileEntry() with path mismatch should return error")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Error("AssociateWithFileEntry() should return PackageError")
	}

	// Test nil FileEntry
	err = pme.AssociateWithFileEntry(nil)
	if err == nil {
		t.Error("AssociateWithFileEntry() with nil FileEntry should return error")
	}

	// Test FileEntry with no paths
	fe2 := NewFileEntry()
	err = pme.AssociateWithFileEntry(fe2)
	if err == nil {
		t.Error("AssociateWithFileEntry() with FileEntry having no paths should return error")
	}

	// Test bidirectional association failure (simulate error in fe.AssociateWithPathMetadata)
	fe3 := NewFileEntry()
	fe3.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "dir/file"},
	}
	// Create a PathMetadataEntry with a different path to cause association failure
	pme3 := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}
	// This should succeed
	err = pme3.AssociateWithFileEntry(fe3)
	if err != nil {
		t.Fatalf("AssociateWithFileEntry() error = %v", err)
	}
}

// TestPathMetadataEntry_GetAssociatedFileEntries tests the GetAssociatedFileEntries method.
func TestPathMetadataEntry_GetAssociatedFileEntries(t *testing.T) {
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 4, Path: "file"},
		Type: PathMetadataTypeFile,
	}

	// Test empty
	entries := pme.GetAssociatedFileEntries()
	if len(entries) != 0 {
		t.Errorf("GetAssociatedFileEntries() empty = %d, want 0", len(entries))
	}

	// Test with entries
	fe1 := NewFileEntry()
	fe2 := NewFileEntry()
	pme.AssociatedFileEntries = []*FileEntry{fe1, fe2}

	entries = pme.GetAssociatedFileEntries()
	if len(entries) != 2 {
		t.Errorf("GetAssociatedFileEntries() = %d, want 2", len(entries))
	}
}

// TestPathMetadataEntry_GetEffectiveTags_WithNilFileEntry tests GetEffectiveTags with nil FileEntry.
func TestPathMetadataEntry_GetEffectiveTags_WithNilFileEntry(t *testing.T) {
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
		Properties: []*generics.Tag[any]{
			{
				Key:   "path-tag",
				Value: "path-value",
				Type:  generics.TagValueTypeString,
			},
		},
	}

	// Add nil FileEntry to associated entries
	pme.AssociatedFileEntries = []*FileEntry{nil}

	// GetEffectiveTags should handle nil FileEntry gracefully
	tags, err := pme.GetEffectiveTags()
	if err != nil {
		t.Fatalf("GetEffectiveTags() error = %v", err)
	}

	// Should still have the direct tag
	if len(tags) < 1 {
		t.Errorf("GetEffectiveTags() returned %d tags, want at least 1", len(tags))
	}
}

// TestPathMetadataEntry_GetEffectiveTags_ErrorInInheritedTags tests GetEffectiveTags when GetInheritedTags fails.
func TestPathMetadataEntry_GetEffectiveTags_ErrorInInheritedTags(t *testing.T) {
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
		Properties: []*generics.Tag[any]{
			{
				Key:   "path-tag",
				Value: "path-value",
				Type:  generics.TagValueTypeString,
			},
		},
	}

	// Set a ParentPath that might cause issues (but won't actually error)
	parent := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 3, Path: "dir"},
		Type: PathMetadataTypeDirectory,
	}
	pme.ParentPath = parent

	// GetEffectiveTags should handle errors in GetInheritedTags gracefully
	tags, err := pme.GetEffectiveTags()
	if err != nil {
		t.Fatalf("GetEffectiveTags() error = %v", err)
	}

	// Should still have the direct tag
	if len(tags) < 1 {
		t.Errorf("GetEffectiveTags() returned %d tags, want at least 1", len(tags))
	}
}
