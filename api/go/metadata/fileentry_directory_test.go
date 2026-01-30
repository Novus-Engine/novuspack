package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
)

// TestGetParentPath tests GetParentPath method
func TestGetParentPath(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FileEntry
		want  string
	}{
		{
			name: "no path metadata entries",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			want: "",
		},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := fe.GetParentPath()

			if got != tt.want {
				t.Errorf("GetParentPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestGetDirectoryDepth tests GetDirectoryDepth method
func TestGetDirectoryDepth(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FileEntry
		want  int
	}{
		{
			name: "no path metadata entries",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			want: 0,
		},
		{
			name: "one level deep",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file.txt", PathLength: 14},
				}
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{Path: "/test/file.txt", PathLength: 14},
				}
				// Set depth to 1
				parent := &PathMetadataEntry{
					Path: generics.PathEntry{Path: "/test", PathLength: 5},
				}
				pme.ParentPath = parent
				if fe.PathMetadataEntries == nil {
					fe.PathMetadataEntries = make(map[string]*PathMetadataEntry)
				}
				fe.PathMetadataEntries["/test/file.txt"] = pme
				return fe
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := fe.GetDirectoryDepth()

			if got != tt.want {
				t.Errorf("GetDirectoryDepth() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestIsRootRelative tests IsRootRelative method
func TestIsRootRelative(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *FileEntry
		want  bool
	}{
		{
			name: "root relative",
			setup: func() *FileEntry {
				return NewFileEntry()
			},
			want: true,
		},
		{
			name: "not root relative",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.Paths = []generics.PathEntry{
					{Path: "/test/file.txt", PathLength: 14},
				}
				pme := &PathMetadataEntry{
					Path: generics.PathEntry{Path: "/test/file.txt", PathLength: 14},
				}
				parent := &PathMetadataEntry{
					Path: generics.PathEntry{Path: "/test", PathLength: 5},
				}
				pme.ParentPath = parent
				if fe.PathMetadataEntries == nil {
					fe.PathMetadataEntries = make(map[string]*PathMetadataEntry)
				}
				fe.PathMetadataEntries["/test/file.txt"] = pme
				return fe
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			got := fe.IsRootRelative()

			if got != tt.want {
				t.Errorf("IsRootRelative() = %v, want %v", got, tt.want)
			}
		})
	}
}
