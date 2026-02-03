package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
)

// fileEntryWithTwoPathsFirstSymlink returns a FileEntry with two paths; the first is a symlink with the given target.
func fileEntryWithTwoPathsFirstSymlink(path1, linkTarget, path2 string) *FileEntry {
	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{Path: path1, IsSymlink: true, LinkTarget: linkTarget},
		{Path: path2},
	}
	return fe
}

// TestHasSymlinks tests HasSymlinks method
func TestHasSymlinks(t *testing.T) {
	tests := []fileEntryTableCase{
		{
			name: "no symlinks",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file1"},
					{Path: "/test/file2"},
				}
				return fe
			},
			want: false,
		},
		{
			name: "has symlinks",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file1", IsSymlink: true, LinkTarget: "/target"},
					{Path: "/test/file2"},
				}
				return fe
			},
			want: true,
		},
		{name: "no paths", setup: NewFileEntry, want: false},
	}
	runFileEntryTableTest(t, tests, func(fe *FileEntry) interface{} { return fe.HasSymlinks() }, "HasSymlinks() = %v, want %v")
}

// TestGetSymlinkPaths tests GetSymlinkPaths method
func TestGetSymlinkPaths(t *testing.T) {
	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{Path: "/test/file1", IsSymlink: true, LinkTarget: "/target1"},
		{Path: "/test/file2"},
		{Path: "/test/file3", IsSymlink: true, LinkTarget: "/target2"},
	}

	symlinks := fe.GetSymlinkPaths()

	if len(symlinks) != 2 {
		t.Errorf("GetSymlinkPaths() returned %d symlinks, want 2", len(symlinks))
	}

	if symlinks[0].Path != "/test/file1" {
		t.Errorf("GetSymlinkPaths() first symlink path = %q, want %q", symlinks[0].Path, "/test/file1")
	}

	if symlinks[1].Path != "/test/file3" {
		t.Errorf("GetSymlinkPaths() second symlink path = %q, want %q", symlinks[1].Path, "/test/file3")
	}
}

// TestGetPrimaryPath tests GetPrimaryPath method
func TestGetPrimaryPath(t *testing.T) {
	tests := []fileEntryTableCase{
		{
			name: "has paths",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file1"},
					{Path: "/test/file2"},
				}
				return fe
			},
			want: "test/file1",
		},
		{name: "no paths", setup: NewFileEntry, want: ""},
	}
	runFileEntryTableTest(t, tests, func(fe *FileEntry) interface{} { return fe.GetPrimaryPath() }, "GetPrimaryPath() = %q, want %q")
}

// TestResolveAllSymlinks tests ResolveAllSymlinks method
func TestResolveAllSymlinks(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FileEntry
		want  []string
	}{
		{
			name: "absolute symlink target",
			setup: func() *FileEntry {
				return fileEntryWithTwoPathsFirstSymlink("/test/file1", "/absolute/target", "/test/file2")
			},
			want: []string{"/absolute/target", "/test/file2"},
		},
		{
			name: "relative symlink target",
			setup: func() *FileEntry {
				return fileEntryWithTwoPathsFirstSymlink("/test/file1", "relative/target", "/test/file2")
			},
			want: []string{"/test/relative/target", "/test/file2"},
		},
		{
			name: "no symlinks",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file1"},
					{Path: "/test/file2"},
				}
				return fe
			},
			want: []string{"/test/file1", "/test/file2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := fe.ResolveAllSymlinks()

			if len(got) != len(tt.want) {
				t.Errorf("ResolveAllSymlinks() returned %d paths, want %d", len(got), len(tt.want))
				return
			}

			for i, wantPath := range tt.want {
				if got[i] != wantPath {
					t.Errorf("ResolveAllSymlinks() [%d] = %q, want %q", i, got[i], wantPath)
				}
			}
		})
	}
}

// TestAssociateWithPathMetadata tests AssociateWithPathMetadata method
func TestAssociateWithPathMetadata(t *testing.T) {
	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "dir/file"},
	}

	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}

	// Test successful association
	err := fe.AssociateWithPathMetadata(pme)
	if err != nil {
		t.Fatalf("AssociateWithPathMetadata() error = %v", err)
	}

	if fe.PathMetadataEntries == nil {
		t.Fatal("AssociateWithPathMetadata() did not initialize PathMetadataEntries map")
	}

	if fe.PathMetadataEntries["dir/file"] != pme {
		t.Error("AssociateWithPathMetadata() did not add PathMetadataEntry correctly")
	}

	// Test nil PathMetadataEntry
	err = fe.AssociateWithPathMetadata(nil)
	if err == nil {
		t.Error("AssociateWithPathMetadata() with nil should return error")
	}

	// Test path mismatch
	pme2 := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 7, Path: "no/match"},
		Type: PathMetadataTypeFile,
	}

	err = fe.AssociateWithPathMetadata(pme2)
	if err == nil {
		t.Error("AssociateWithPathMetadata() with path mismatch should return error")
	}
}

// TestGetPathMetadataForPath tests GetPathMetadataForPath method
func TestGetPathMetadataForPath(t *testing.T) {
	fe := NewFileEntry()
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{PathLength: 8, Path: "dir/file"},
		Type: PathMetadataTypeFile,
	}

	// Test nil map
	if got := fe.GetPathMetadataForPath("dir/file"); got != nil {
		t.Errorf("GetPathMetadataForPath() with nil map = %v, want nil", got)
	}

	// Test with entry
	fe.PathMetadataEntries = make(map[string]*PathMetadataEntry)
	fe.PathMetadataEntries["dir/file"] = pme

	if got := fe.GetPathMetadataForPath("dir/file"); got != pme {
		t.Errorf("GetPathMetadataForPath() = %v, want %v", got, pme)
	}

	// Test non-existent path
	if got := fe.GetPathMetadataForPath("nonexistent"); got != nil {
		t.Errorf("GetPathMetadataForPath() for nonexistent = %v, want nil", got)
	}
}

// TestGetPaths tests GetPaths method
func TestGetPaths(t *testing.T) {
	fe := NewFileEntry()
	expectedPaths := []generics.PathEntry{
		{PathLength: 4, Path: "file"},
		{PathLength: 8, Path: "dir/file"},
	}
	fe.Paths = expectedPaths

	got := fe.GetPaths()
	if len(got) != len(expectedPaths) {
		t.Errorf("GetPaths() returned %d paths, want %d", len(got), len(expectedPaths))
	}

	for i, expected := range expectedPaths {
		if got[i].Path != expected.Path {
			t.Errorf("GetPaths() [%d] = %q, want %q", i, got[i].Path, expected.Path)
		}
	}
}

// TestGetFileID tests GetFileID method
func TestGetFileID(t *testing.T) {
	fe := NewFileEntry()
	fe.FileID = 12345

	if got := fe.GetFileID(); got != 12345 {
		t.Errorf("GetFileID() = %d, want %d", got, 12345)
	}
}
