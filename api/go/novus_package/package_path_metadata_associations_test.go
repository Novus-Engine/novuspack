// This file contains unit tests for file-path association operations.
// It tests AssociateFileWithPath, DisassociateFileFromPath, UpdateFilePathAssociations,
// and GetFilePathAssociations methods from package_path_metadata_associations.go.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"context"
	"testing"

	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// =============================================================================
// TEST: AssociateFileWithPath
// =============================================================================

// TestPackage_AssociateFileWithPath_Success tests AssociateFileWithPath successful association.
// Expected: Should link FileEntry to PathMetadataEntry
func TestPackage_AssociateFileWithPath_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true

	// Create FileEntry
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "test.txt", PathLength: 8},
	}
	fpkg.FileEntries = []*metadata.FileEntry{fe}

	// Create PathMetadataEntry
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "test.txt", PathLength: 8},
		Type: metadata.PathMetadataTypeFile,
	}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pme}

	// AssociateFileWithPath should succeed
	err = fpkg.AssociateFileWithPath(ctx, "test.txt", "test.txt")
	if err != nil {
		t.Fatalf("AssociateFileWithPath failed: %v", err)
	}

	// Verify association
	if fe.PathMetadataEntries == nil || fe.PathMetadataEntries["test.txt"] != pme {
		t.Error("FileEntry should be associated with PathMetadataEntry")
	}
	// Verify bidirectional
	if len(pme.AssociatedFileEntries) == 0 || pme.AssociatedFileEntries[0] != fe {
		t.Error("PathMetadataEntry should be associated with FileEntry")
	}
}

// TestPackage_AssociateFileWithPath_FileNotFound tests AssociateFileWithPath when file doesn't exist.
// Expected: Should return validation error
func TestPackage_AssociateFileWithPath_FileNotFound(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.FileEntries = []*metadata.FileEntry{} // Empty - no files

	// Create PathMetadataEntry
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "test.txt", PathLength: 8},
		Type: metadata.PathMetadataTypeFile,
	}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pme}

	// AssociateFileWithPath should fail - file not found
	err = fpkg.AssociateFileWithPath(ctx, "nonexistent.txt", "test.txt")
	if err == nil {
		t.Fatal("AssociateFileWithPath should fail when file not found")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want ErrTypeValidation", pkgErr.Type)
	}
}

// TestPackage_AssociateFileWithPath_PathNotFound tests AssociateFileWithPath when path metadata doesn't exist.
// Expected: Should return validation error
func TestPackage_AssociateFileWithPath_PathNotFound(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true

	// Create FileEntry
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "test.txt", PathLength: 8},
	}
	fpkg.FileEntries = []*metadata.FileEntry{fe}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{} // Empty - no path metadata

	// AssociateFileWithPath should fail - path metadata not found
	err = fpkg.AssociateFileWithPath(ctx, "test.txt", "nonexistent.txt")
	if err == nil {
		t.Fatal("AssociateFileWithPath should fail when path metadata not found")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want ErrTypeValidation", pkgErr.Type)
	}
}

// TestPackage_AssociateFileWithPath_CancelledContext tests AssociateFileWithPath with cancelled context.
// Expected: Should return context error
func TestPackage_AssociateFileWithPath_CancelledContext(t *testing.T) {
	cancelledCtx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// AssociateFileWithPath should fail with cancelled context
	err = fpkg.AssociateFileWithPath(cancelledCtx, "file.txt", "file.txt")
	if err == nil {
		t.Fatal("AssociateFileWithPath should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want ErrTypeContext", pkgErr.Type)
	}
}

// =============================================================================
// TEST: DisassociateFileFromPath
// =============================================================================

// TestPackage_DisassociateFileFromPath_Success tests DisassociateFileFromPath removing association.
// Expected: Should remove bidirectional association
func TestPackage_DisassociateFileFromPath_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true

	// Create FileEntry with association
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "test.txt", PathLength: 8},
	}
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "test.txt", PathLength: 8},
		Type: metadata.PathMetadataTypeFile,
	}

	// Set up bidirectional association manually
	fe.PathMetadataEntries = make(map[string]*metadata.PathMetadataEntry)
	fe.PathMetadataEntries["test.txt"] = pme
	pme.AssociatedFileEntries = []*metadata.FileEntry{fe}

	fpkg.FileEntries = []*metadata.FileEntry{fe}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pme}

	// DisassociateFileFromPath should remove association
	err = fpkg.DisassociateFileFromPath(ctx, "test.txt")
	if err != nil {
		t.Fatalf("DisassociateFileFromPath failed: %v", err)
	}

	// Verify association was removed from FileEntry
	if fe.PathMetadataEntries != nil && fe.PathMetadataEntries["test.txt"] != nil {
		t.Error("Association should be removed from FileEntry.PathMetadataEntries")
	}
	// Verify association was removed from PathMetadataEntry
	if len(pme.AssociatedFileEntries) != 0 {
		t.Error("Association should be removed from PathMetadataEntry.AssociatedFileEntries")
	}
}

// TestPackage_DisassociateFileFromPath_FileNotFound tests DisassociateFileFromPath when file doesn't exist.
// Expected: Should return validation error
func TestPackage_DisassociateFileFromPath_FileNotFound(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.FileEntries = []*metadata.FileEntry{} // Empty

	// DisassociateFileFromPath should fail - file not found
	err = fpkg.DisassociateFileFromPath(ctx, "nonexistent.txt")
	if err == nil {
		t.Fatal("DisassociateFileFromPath should fail when file not found")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want ErrTypeValidation", pkgErr.Type)
	}
}

// TestPackage_DisassociateFileFromPath_CancelledContext tests DisassociateFileFromPath with cancelled context.
// Expected: Should return context error
func TestPackage_DisassociateFileFromPath_CancelledContext(t *testing.T) {
	cancelledCtx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// DisassociateFileFromPath should fail with cancelled context
	err = fpkg.DisassociateFileFromPath(cancelledCtx, "file.txt")
	if err == nil {
		t.Fatal("DisassociateFileFromPath should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want ErrTypeContext", pkgErr.Type)
	}
}

// =============================================================================
// TEST: UpdateFilePathAssociations
// =============================================================================

// TestPackage_UpdateFilePathAssociations_BasicMatching tests UpdateFilePathAssociations path matching.
// Expected: Should match FileEntry.Paths to PathMetadataEntry.Path
func TestPackage_UpdateFilePathAssociations_BasicMatching(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true

	// Create FileEntry
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "file1.txt", PathLength: 9},
	}
	fpkg.FileEntries = []*metadata.FileEntry{fe}

	// Create PathMetadataEntry
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "file1.txt", PathLength: 9},
		Type: metadata.PathMetadataTypeFile,
	}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pme}

	// UpdateFilePathAssociations should match and associate
	err = fpkg.UpdateFilePathAssociations(ctx)
	if err != nil {
		t.Fatalf("UpdateFilePathAssociations failed: %v", err)
	}

	// Verify association was created
	if fe.PathMetadataEntries == nil {
		t.Fatal("FileEntry.PathMetadataEntries should not be nil")
	}
	if fe.PathMetadataEntries["file1.txt"] != pme {
		t.Error("FileEntry should be associated with PathMetadataEntry")
	}
	// Verify bidirectional association
	if len(pme.AssociatedFileEntries) == 0 {
		t.Error("PathMetadataEntry.AssociatedFileEntries should contain FileEntry")
	}
	if pme.AssociatedFileEntries[0] != fe {
		t.Error("PathMetadataEntry should be associated with FileEntry")
	}
}

// TestPackage_UpdateFilePathAssociations_ParentChain tests UpdateFilePathAssociations parent chain building.
// Expected: Should build parent path hierarchy
func TestPackage_UpdateFilePathAssociations_ParentChain(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.FileEntries = []*metadata.FileEntry{}

	// Create hierarchical PathMetadataEntries
	pmeRoot := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "dir/", PathLength: 4},
		Type: metadata.PathMetadataTypeDirectory,
	}
	pmeSubdir := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "dir/subdir/", PathLength: 11},
		Type: metadata.PathMetadataTypeDirectory,
	}
	pmeFile := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "dir/subdir/file.txt", PathLength: 19},
		Type: metadata.PathMetadataTypeFile,
	}

	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pmeRoot, pmeSubdir, pmeFile}

	// UpdateFilePathAssociations should build parent chain
	err = fpkg.UpdateFilePathAssociations(ctx)
	if err != nil {
		t.Fatalf("UpdateFilePathAssociations failed: %v", err)
	}

	// Verify parent path chain
	if pmeFile.ParentPath != pmeSubdir {
		t.Error("pmeFile.ParentPath should point to pmeSubdir")
	}
	if pmeSubdir.ParentPath != pmeRoot {
		t.Error("pmeSubdir.ParentPath should point to pmeRoot")
	}
	if pmeRoot.ParentPath != nil {
		t.Error("pmeRoot.ParentPath should be nil (root level)")
	}
}

// TestPackage_UpdateFilePathAssociations_MultiplePaths tests UpdateFilePathAssociations with files having multiple paths.
// Expected: Should associate all paths correctly
func TestPackage_UpdateFilePathAssociations_MultiplePaths(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true

	// Create FileEntry with multiple paths
	fe := metadata.NewFileEntry()
	fe.FileID = 1
	fe.Paths = []generics.PathEntry{
		{Path: "path1/file.txt", PathLength: 14},
		{Path: "path2/file.txt", PathLength: 14},
	}
	fpkg.FileEntries = []*metadata.FileEntry{fe}

	// Create PathMetadataEntries for both paths
	pme1 := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "path1/file.txt", PathLength: 14},
		Type: metadata.PathMetadataTypeFile,
	}
	pme2 := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "path2/file.txt", PathLength: 14},
		Type: metadata.PathMetadataTypeFile,
	}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pme1, pme2}

	// UpdateFilePathAssociations should associate all paths
	err = fpkg.UpdateFilePathAssociations(ctx)
	if err != nil {
		t.Fatalf("UpdateFilePathAssociations failed: %v", err)
	}

	// Verify both associations
	if fe.PathMetadataEntries["path1/file.txt"] != pme1 {
		t.Error("First path should be associated with pme1")
	}
	if fe.PathMetadataEntries["path2/file.txt"] != pme2 {
		t.Error("Second path should be associated with pme2")
	}
}

// TestPackage_UpdateFilePathAssociations_CancelledContext tests UpdateFilePathAssociations with cancelled context.
// Expected: Should return context error
func TestPackage_UpdateFilePathAssociations_CancelledContext(t *testing.T) {
	runContextCancelledTest(t, func(fpkg *filePackage, ctx context.Context) error {
		return fpkg.UpdateFilePathAssociations(ctx)
	})
}

// =============================================================================
// TEST: GetFilePathAssociations
// =============================================================================

// TestPackage_GetFilePathAssociations_Basic tests basic GetFilePathAssociations operation.
func TestPackage_GetFilePathAssociations_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetFilePathAssociations should succeed even with empty package
	fpkg := pkg.(*filePackage)
	associations, err := fpkg.GetFilePathAssociations(ctx)
	if err != nil {
		t.Errorf("GetFilePathAssociations failed: %v", err)
	}
	if associations == nil {
		t.Error("GetFilePathAssociations should return non-nil map")
	}
	if len(associations) != 0 {
		t.Errorf("GetFilePathAssociations should return empty map for new package, got %d entries", len(associations))
	}
}

// TestPackage_GetFilePathAssociations_WithAssociations tests GetFilePathAssociations with file associations.
func TestPackage_GetFilePathAssociations_WithAssociations(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Create FileEntry with path
	fe1 := metadata.NewFileEntry()
	fe1.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "file1.txt"},
	}

	// Create PathMetadataEntry
	pme1 := createValidPathMetadataEntry("file1.txt", metadata.PathMetadataTypeFile)

	// Associate them
	fe1.PathMetadataEntries = make(map[string]*metadata.PathMetadataEntry)
	fe1.PathMetadataEntries["file1.txt"] = pme1

	// Create another FileEntry with multiple paths
	fe2 := metadata.NewFileEntry()
	fe2.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "file2.txt"},
		{PathLength: 9, Path: "file3.txt"},
	}

	pme2 := createValidPathMetadataEntry("file2.txt", metadata.PathMetadataTypeFile)
	pme3 := createValidPathMetadataEntry("file3.txt", metadata.PathMetadataTypeFile)

	fe2.PathMetadataEntries = make(map[string]*metadata.PathMetadataEntry)
	fe2.PathMetadataEntries["file2.txt"] = pme2
	fe2.PathMetadataEntries["file3.txt"] = pme3

	// Set FileEntries
	fpkg.FileEntries = []*metadata.FileEntry{fe1, fe2}

	associations, err := fpkg.GetFilePathAssociations(ctx)
	if err != nil {
		t.Errorf("GetFilePathAssociations should succeed, got error: %v", err)
	}
	if len(associations) != 3 {
		t.Errorf("GetFilePathAssociations should return 3 associations, got %d", len(associations))
	}
	if associations["file1.txt"] != pme1 {
		t.Error("GetFilePathAssociations should return correct association for file1.txt")
	}
	if associations["file2.txt"] != pme2 {
		t.Error("GetFilePathAssociations should return correct association for file2.txt")
	}
	if associations["file3.txt"] != pme3 {
		t.Error("GetFilePathAssociations should return correct association for file3.txt")
	}
}

// TestPackage_GetFilePathAssociations_WithNilEntries tests GetFilePathAssociations with nil entries.
func TestPackage_GetFilePathAssociations_WithNilEntries(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Create FileEntry with nil PathMetadataEntry in map
	fe := metadata.NewFileEntry()
	fe.Paths = []generics.PathEntry{
		{PathLength: 8, Path: "file.txt"},
	}
	fe.PathMetadataEntries = make(map[string]*metadata.PathMetadataEntry)
	fe.PathMetadataEntries["file.txt"] = nil // nil entry should be skipped

	fpkg.FileEntries = []*metadata.FileEntry{fe}

	associations, err := fpkg.GetFilePathAssociations(ctx)
	if err != nil {
		t.Errorf("GetFilePathAssociations should succeed, got error: %v", err)
	}
	if len(associations) != 0 {
		t.Errorf("GetFilePathAssociations should skip nil entries, got %d", len(associations))
	}
}

// TestPackage_GetFilePathAssociations_WithContext tests GetFilePathAssociations with cancelled context.
func TestPackage_GetFilePathAssociations_WithContext(t *testing.T) {
	runWithCancelledContext(t, func(fpkg *filePackage, ctx context.Context) (interface{}, error) {
		return fpkg.GetFilePathAssociations(ctx)
	}, "GetFilePathAssociations")
}
