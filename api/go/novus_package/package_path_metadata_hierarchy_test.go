// This file contains unit tests for path hierarchy operations.
// It tests GetPathInfo, ListPaths, and GetPathHierarchy methods from
// package_path_metadata_hierarchy.go.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/metadata"
)

// =============================================================================
// TEST: GetPathInfo
// =============================================================================

// TestPackage_GetPathInfo_Basic tests basic GetPathInfo operation.
func TestPackage_GetPathInfo_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetPathInfo should now work since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	_, err = fpkg.GetPathInfo("test/path")
	// Expected to return "not found" error since path doesn't exist
	if err == nil {
		t.Error("GetPathInfo should return error when path not found")
	}
}

// =============================================================================
// TEST: ListPaths
// =============================================================================

// TestPackage_ListPaths_Basic tests basic ListPaths operation.
func TestPackage_ListPaths_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// ListPaths should now succeed since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	_, err = fpkg.ListPaths()
	if err != nil {
		t.Errorf("ListPaths failed: %v", err)
	}
}

// =============================================================================
// TEST: GetPathHierarchy
// =============================================================================

// TestPackage_GetPathHierarchy_Basic tests basic GetPathHierarchy operation.
func TestPackage_GetPathHierarchy_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetPathHierarchy should now succeed since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	_, err = fpkg.GetPathHierarchy()
	if err != nil {
		t.Errorf("GetPathHierarchy failed: %v", err)
	}
}

// =============================================================================
// TEST: GetPathInfo with cached entries
// =============================================================================

// TestPackage_GetPathInfo_NotFound tests GetPathInfo when path not found.
func TestPackage_GetPathInfo_NotFound(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("other/path", metadata.PathMetadataTypeFile),
	}

	_, err = fpkg.GetPathInfo("nonexistent/path")
	if err == nil {
		t.Error("GetPathInfo should return error when path not found")
	}
}

// TestPackage_GetPathInfo_Success tests successful GetPathInfo operation.
func TestPackage_GetPathInfo_Success(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	targetPath := "test/path"
	entry := createValidPathMetadataEntry(targetPath, metadata.PathMetadataTypeFile)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{entry}

	pathInfo, err := fpkg.GetPathInfo(targetPath)
	if err != nil {
		t.Errorf("GetPathInfo should succeed, got error: %v", err)
	}
	if pathInfo == nil {
		t.Fatal("GetPathInfo should return non-nil PathInfo")
	}
	if pathInfo.Entry.GetPath() != targetPath {
		t.Errorf("GetPathInfo should return correct entry, got path %q", pathInfo.Entry.GetPath())
	}
	if pathInfo.FileCount != 0 {
		t.Errorf("GetPathInfo should return correct file count, got %d", pathInfo.FileCount)
	}
}

// TestPackage_GetPathInfo_WithParent tests GetPathInfo with parent path.
func TestPackage_GetPathInfo_WithParent(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	parentPath := "parent"
	childPath := "parent/child"
	parent := createValidPathMetadataEntry(parentPath, metadata.PathMetadataTypeDirectory)
	child := createValidPathMetadataEntry(childPath, metadata.PathMetadataTypeFile)
	child.ParentPath = parent
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{parent, child}

	pathInfo, err := fpkg.GetPathInfo(childPath)
	if err != nil {
		t.Errorf("GetPathInfo should succeed, got error: %v", err)
	}
	if pathInfo.ParentPath != parentPath {
		t.Errorf("GetPathInfo should return correct parent path, got %q", pathInfo.ParentPath)
	}
}

// TestPackage_GetPathInfo_WithSubDirs tests GetPathInfo with subdirectories.
func TestPackage_GetPathInfo_WithSubDirs(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	parentPath := "parent"
	subDir1 := "parent/subdir1"
	subDir2 := "parent/subdir2"
	parent := createValidPathMetadataEntry(parentPath, metadata.PathMetadataTypeDirectory)
	sub1 := createValidPathMetadataEntry(subDir1, metadata.PathMetadataTypeDirectory)
	sub2 := createValidPathMetadataEntry(subDir2, metadata.PathMetadataTypeDirectory)
	sub1.ParentPath = parent
	sub2.ParentPath = parent
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{parent, sub1, sub2}

	pathInfo, err := fpkg.GetPathInfo(parentPath)
	if err != nil {
		t.Errorf("GetPathInfo should succeed, got error: %v", err)
	}
	if len(pathInfo.SubDirs) != 2 {
		t.Errorf("GetPathInfo should return 2 subdirectories, got %d", len(pathInfo.SubDirs))
	}
}

// TestPackage_GetPathInfo_WithAssociatedFiles tests GetPathInfo with associated files.
func TestPackage_GetPathInfo_WithAssociatedFiles(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	targetPath := "test/path"
	entry := createValidPathMetadataEntry(targetPath, metadata.PathMetadataTypeFile)
	fileEntry := metadata.NewFileEntry()
	entry.AssociatedFileEntries = []*metadata.FileEntry{fileEntry}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{entry}

	pathInfo, err := fpkg.GetPathInfo(targetPath)
	if err != nil {
		t.Errorf("GetPathInfo should succeed, got error: %v", err)
	}
	if pathInfo.FileCount != 1 {
		t.Errorf("GetPathInfo should return correct file count, got %d, want 1", pathInfo.FileCount)
	}
}

// =============================================================================
// TEST: ListPaths with cached entries
// =============================================================================

// TestPackage_ListPaths_Success tests successful ListPaths operation.
func TestPackage_ListPaths_Success(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("path1", metadata.PathMetadataTypeFile),
		createValidPathMetadataEntry("path2", metadata.PathMetadataTypeFile),
	}

	pathInfos, err := fpkg.ListPaths()
	if err != nil {
		t.Errorf("ListPaths should succeed, got error: %v", err)
	}
	if len(pathInfos) != 2 {
		t.Errorf("ListPaths should return 2 paths, got %d", len(pathInfos))
	}
}

// TestPackage_ListPaths_Empty tests ListPaths with empty entries.
func TestPackage_ListPaths_Empty(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}

	pathInfos, err := fpkg.ListPaths()
	if err != nil {
		t.Errorf("ListPaths should succeed with empty entries, got error: %v", err)
	}
	if len(pathInfos) != 0 {
		t.Errorf("ListPaths should return empty slice, got %d", len(pathInfos))
	}
}

// =============================================================================
// TEST: GetPathHierarchy with cached entries
// =============================================================================

// TestPackage_GetPathHierarchy_Success tests successful GetPathHierarchy operation.
func TestPackage_GetPathHierarchy_Success(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	parentPath := "parent"
	child1 := "parent/child1"
	child2 := "parent/child2"
	rootPath := "root"
	parent := createValidPathMetadataEntry(parentPath, metadata.PathMetadataTypeDirectory)
	ch1 := createValidPathMetadataEntry(child1, metadata.PathMetadataTypeFile)
	ch2 := createValidPathMetadataEntry(child2, metadata.PathMetadataTypeFile)
	root := createValidPathMetadataEntry(rootPath, metadata.PathMetadataTypeFile)
	ch1.ParentPath = parent
	ch2.ParentPath = parent
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{parent, ch1, ch2, root}

	hierarchy, err := fpkg.GetPathHierarchy()
	if err != nil {
		t.Errorf("GetPathHierarchy should succeed, got error: %v", err)
	}
	if len(hierarchy[parentPath]) != 2 {
		t.Errorf("GetPathHierarchy should return 2 children for parent, got %d", len(hierarchy[parentPath]))
	}
	// Parent is also a root path (no parent), and root is also a root path
	if len(hierarchy[""]) < 2 {
		t.Errorf("GetPathHierarchy should return at least 2 root paths (parent and root), got %d", len(hierarchy[""]))
	}
}

// TestPackage_GetPathHierarchy_Empty tests GetPathHierarchy with empty entries.
func TestPackage_GetPathHierarchy_Empty(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}

	hierarchy, err := fpkg.GetPathHierarchy()
	if err != nil {
		t.Errorf("GetPathHierarchy should succeed with empty entries, got error: %v", err)
	}
	if len(hierarchy) != 0 {
		t.Errorf("GetPathHierarchy should return empty map, got %d entries", len(hierarchy))
	}
}

// TestPackage_GetPathHierarchy_AllRoot tests GetPathHierarchy with all root paths.
func TestPackage_GetPathHierarchy_AllRoot(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("root1", metadata.PathMetadataTypeFile),
		createValidPathMetadataEntry("root2", metadata.PathMetadataTypeFile),
	}

	hierarchy, err := fpkg.GetPathHierarchy()
	if err != nil {
		t.Errorf("GetPathHierarchy should succeed, got error: %v", err)
	}
	if len(hierarchy[""]) != 2 {
		t.Errorf("GetPathHierarchy should return 2 root paths, got %d", len(hierarchy[""]))
	}
}

// =============================================================================
// NOTE: Context cancellation tests removed
// =============================================================================
// GetPathInfo, ListPaths, ListDirectories, and GetPathHierarchy are now
// pure in-memory operations and do not require context parameters.
// Per spec: api_metadata.md Section 8.2 lines 1092-1100
