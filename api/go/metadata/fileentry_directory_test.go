package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
)

type fileEntryTableCase struct {
	name  string
	setup func() *FileEntry
	want  interface{}
}

func runFileEntryTableTest(t *testing.T, tests []fileEntryTableCase, getter func(*FileEntry) interface{}, format string) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := getter(fe)
			if got != tt.want {
				t.Errorf(format, got, tt.want)
			}
		})
	}
}

// fileEntryWithPathAndParent returns a FileEntry with one path and PathMetadataEntry with given parent path (depth 1).
func fileEntryWithPathAndParent(filePath, parentPath string) *FileEntry {
	fe := NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{Path: filePath, PathLength: uint16(len(filePath))},
	}
	pme := &PathMetadataEntry{
		Path: generics.PathEntry{Path: filePath, PathLength: uint16(len(filePath))},
	}
	parent := &PathMetadataEntry{
		Path: generics.PathEntry{Path: parentPath, PathLength: uint16(len(parentPath))},
	}
	pme.ParentPath = parent
	if fe.PathMetadataEntries == nil {
		fe.PathMetadataEntries = make(map[string]*PathMetadataEntry)
	}
	fe.PathMetadataEntries[filePath] = pme
	return fe
}

// TestGetParentPath tests GetParentPath method
func TestGetParentPath(t *testing.T) {
	tests := []fileEntryTableCase{
		{name: "no path metadata entries", setup: NewFileEntry, want: ""},
		{
			name: "with path metadata entry",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file.txt", PathLength: 14},
				}
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{Path: "/test/file.txt", PathLength: 14},
				}
				if fe.PathMetadataEntries == nil {
					fe.PathMetadataEntries = make(map[string]*PathMetadataEntry)
				}
				fe.PathMetadataEntries["/test/file.txt"] = pme
				return fe
			},
			want: "/test",
		},
	}
	runFileEntryTableTest(t, tests, func(fe *FileEntry) interface{} { return fe.GetParentPath() }, "GetParentPath() = %q, want %q")
}

// TestGetDirectoryDepth tests GetDirectoryDepth method
func TestGetDirectoryDepth(t *testing.T) {
	tests := []fileEntryTableCase{
		{name: "no path metadata entries", setup: NewFileEntry, want: 0},
		{name: "one level deep", setup: func() *FileEntry { return fileEntryWithPathAndParent("/test/file.txt", "/test") }, want: 1},
	}
	runFileEntryTableTest(t, tests, func(fe *FileEntry) interface{} { return fe.GetDirectoryDepth() }, "GetDirectoryDepth() = %d, want %d")
}

// TestIsRootRelative tests IsRootRelative method
func TestIsRootRelative(t *testing.T) {
	tests := []fileEntryTableCase{
		{name: "root relative", setup: NewFileEntry, want: true},
		{name: "not root relative", setup: func() *FileEntry { return fileEntryWithPathAndParent("/test/file.txt", "/test") }, want: false},
	}
	runFileEntryTableTest(t, tests, func(fe *FileEntry) interface{} { return fe.IsRootRelative() }, "IsRootRelative() = %v, want %v")
}
