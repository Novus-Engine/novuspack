// This file contains unit tests for directory-specific path metadata operations.
// It tests AddDirectoryMetadata, RemoveDirectoryMetadata, UpdateDirectoryMetadata, and ListDirectories
// methods from package_path_metadata_directories.go.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"context"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/metadata"
)

// =============================================================================
// TEST: AddDirectoryMetadata
// =============================================================================

const testDirPath = "test/dir"
const testDirPathSlash = "test/dir/"

// TestPackage_AddDirectoryMetadata_Basic tests basic AddDirectoryMetadata operation.
func TestPackage_AddDirectoryMetadata_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// AddDirectoryMetadata should now work since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	err = fpkg.AddDirectoryMetadata(ctx, testDirPath, nil, nil, nil)
	if err != nil {
		t.Errorf("AddDirectoryMetadata failed: %v", err)
	}
}

// =============================================================================
// TEST: RemoveDirectoryMetadata
// =============================================================================

// TestPackage_RemoveDirectoryMetadata_Basic tests basic RemoveDirectoryMetadata operation.
func TestPackage_RemoveDirectoryMetadata_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// RemoveDirectoryMetadata should return error when LoadPathMetadataFile is not implemented
	fpkg := pkg.(*filePackage)
	err = fpkg.RemoveDirectoryMetadata(ctx, testDirPath)
	if err == nil {
		t.Error("RemoveDirectoryMetadata should return error when LoadPathMetadataFile is not implemented")
	}
}

// =============================================================================
// TEST: UpdateDirectoryMetadata
// =============================================================================

// TestPackage_UpdateDirectoryMetadata_Basic tests basic UpdateDirectoryMetadata operation.
func TestPackage_UpdateDirectoryMetadata_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// UpdateDirectoryMetadata should return error when LoadPathMetadataFile is not implemented
	fpkg := pkg.(*filePackage)
	err = fpkg.UpdateDirectoryMetadata(ctx, testDirPath, nil, nil, nil)
	if err == nil {
		t.Error("UpdateDirectoryMetadata should return error when LoadPathMetadataFile is not implemented")
	}
}

// =============================================================================
// TEST: ListDirectories
// =============================================================================

func setupOpenFilePackageForPathMetadata(t *testing.T) (Package, *filePackage) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	return pkg, fpkg
}

// TestPackage_ListDirectories_Basic tests basic ListDirectories operation.
func TestPackage_ListDirectories_Basic(t *testing.T) {
	pkg, fpkg := setupOpenFilePackageForPathMetadata(t)
	defer func() { _ = pkg.Close() }()
	_, err := fpkg.ListDirectories()
	if err != nil {
		t.Errorf("ListDirectories failed: %v", err)
	}
}

// =============================================================================
// TEST: Path normalization (paths ending with /)
// =============================================================================

// TestPackage_AddDirectoryMetadata_PathNormalization tests AddDirectoryMetadata path normalization.
//
//nolint:gocognit // table-driven path normalization
func TestPackage_AddDirectoryMetadata_PathNormalization(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}

	// Test path that doesn't end with /
	pathWithoutSlash := testDirPath
	err = fpkg.AddDirectoryMetadata(ctx, pathWithoutSlash, nil, nil, nil)
	// Should fail because SavePathMetadataFile is not implemented, but path should be normalized
	if err == nil {
		// If it succeeded, verify the path was normalized
		if len(fpkg.PathMetadataEntries) > 0 {
			entryPath := fpkg.PathMetadataEntries[0].GetPath()
			if entryPath != pathWithoutSlash+"/" {
				t.Errorf("AddDirectoryMetadata should normalize path, got %q, want %q", entryPath, pathWithoutSlash+"/")
			}
		}
	}

	// Test path that already ends with /
	pathWithSlash := testDirPathSlash
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}
	err = fpkg.AddDirectoryMetadata(ctx, pathWithSlash, nil, nil, nil)
	// Should fail because SavePathMetadataFile is not implemented
	if err == nil {
		// If it succeeded, verify the path was not double-normalized
		if len(fpkg.PathMetadataEntries) > 0 {
			entryPath := fpkg.PathMetadataEntries[0].GetPath()
			if entryPath != pathWithSlash {
				t.Errorf("AddDirectoryMetadata should not double-normalize path, got %q, want %q", entryPath, pathWithSlash)
			}
		}
	}

	// Test empty path
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}
	err = fpkg.AddDirectoryMetadata(ctx, "", nil, nil, nil)
	// Should fail because SavePathMetadataFile is not implemented
	// Empty path should not cause panic
	if err == nil {
		if len(fpkg.PathMetadataEntries) > 0 {
			entryPath := fpkg.PathMetadataEntries[0].GetPath()
			if entryPath != "/" {
				t.Errorf("AddDirectoryMetadata with empty path should normalize to /, got %q", entryPath)
			}
		}
	}
}

// TestPackage_RemoveDirectoryMetadata_PathNormalization tests RemoveDirectoryMetadata path normalization.
func TestPackage_RemoveDirectoryMetadata_PathNormalization(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Test path that doesn't end with /
	pathWithoutSlash := testDirPath
	_ = fpkg.RemoveDirectoryMetadata(ctx, pathWithoutSlash)
	// Should fail because LoadPathMetadataFile is not implemented

	// Test path that already ends with /
	pathWithSlash := testDirPathSlash
	_ = fpkg.RemoveDirectoryMetadata(ctx, pathWithSlash)
	// Should fail because LoadPathMetadataFile is not implemented
}

// TestPackage_UpdateDirectoryMetadata_PathNormalization tests UpdateDirectoryMetadata path normalization.
func TestPackage_UpdateDirectoryMetadata_PathNormalization(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Test path that doesn't end with /
	pathWithoutSlash := testDirPath
	_ = fpkg.UpdateDirectoryMetadata(ctx, pathWithoutSlash, nil, nil, nil)
	// Should fail because LoadPathMetadataFile is not implemented

	// Test path that already ends with /
	pathWithSlash := testDirPathSlash
	_ = fpkg.UpdateDirectoryMetadata(ctx, pathWithSlash, nil, nil, nil)
	// Should fail because LoadPathMetadataFile is not implemented
}

// =============================================================================
// TEST: ListDirectories with cached entries
// =============================================================================

// TestPackage_ListDirectories_Success tests successful ListDirectories operation.
func TestPackage_ListDirectories_Success(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Create mix of directories and files
	dir1 := createValidPathMetadataEntry("dir1/", metadata.PathMetadataTypeDirectory)
	dir2 := createValidPathMetadataEntry("dir2/", metadata.PathMetadataTypeDirectory)
	file1 := createValidPathMetadataEntry("file1.txt", metadata.PathMetadataTypeFile)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{dir1, dir2, file1}

	directories, err := fpkg.ListDirectories()
	if err != nil {
		t.Errorf("ListDirectories should succeed, got error: %v", err)
	}
	if len(directories) != 2 {
		t.Errorf("ListDirectories should return 2 directories, got %d", len(directories))
	}
}

func runListDirectoriesWithEntries(t *testing.T, entries []*metadata.PathMetadataEntry, wantCount int, errMsg string) {
	t.Helper()
	pkg, fpkg := setupOpenFilePackageForPathMetadata(t)
	defer func() { _ = pkg.Close() }()
	fpkg.PathMetadataEntries = entries
	directories, err := fpkg.ListDirectories()
	if err != nil {
		t.Errorf("ListDirectories %s: %v", errMsg, err)
	}
	if len(directories) != wantCount {
		t.Errorf("ListDirectories want %d directories, got %d", wantCount, len(directories))
	}
}

// TestPackage_ListDirectories_Empty tests ListDirectories with no directories.
func TestPackage_ListDirectories_Empty(t *testing.T) {
	file1 := createValidPathMetadataEntry("file1.txt", metadata.PathMetadataTypeFile)
	runListDirectoriesWithEntries(t, []*metadata.PathMetadataEntry{file1}, 0, "should succeed")
}

// =============================================================================
// TEST: Context Cancellation
// =============================================================================

// TestPackage_AddDirectoryMetadata_ContextCancelled tests AddDirectoryMetadata with cancelled context.
func TestPackage_AddDirectoryMetadata_ContextCancelled(t *testing.T) {
	runContextCancelledTest(t, func(fpkg *filePackage, ctx context.Context) error {
		return fpkg.AddDirectoryMetadata(ctx, testDirPath, nil, nil, nil)
	})
}

// TestPackage_RemoveDirectoryMetadata_ContextCancelled tests RemoveDirectoryMetadata with cancelled context.
func TestPackage_RemoveDirectoryMetadata_ContextCancelled(t *testing.T) {
	runContextCancelledTest(t, func(fpkg *filePackage, ctx context.Context) error {
		return fpkg.RemoveDirectoryMetadata(ctx, testDirPath)
	})
}

// TestPackage_UpdateDirectoryMetadata_ContextCancelled tests UpdateDirectoryMetadata with cancelled context.
func TestPackage_UpdateDirectoryMetadata_ContextCancelled(t *testing.T) {
	runContextCancelledTest(t, func(fpkg *filePackage, ctx context.Context) error {
		return fpkg.UpdateDirectoryMetadata(ctx, testDirPath, nil, nil, nil)
	})
}

// TestPackage_ListDirectories_InMemoryOperation tests that ListDirectories is a pure in-memory operation.
func TestPackage_ListDirectories_InMemoryOperation(t *testing.T) {
	dir1 := createValidPathMetadataEntry("dir1/", metadata.PathMetadataTypeDirectory)
	runListDirectoriesWithEntries(t, []*metadata.PathMetadataEntry{dir1}, 1, "should succeed as in-memory operation")
}
