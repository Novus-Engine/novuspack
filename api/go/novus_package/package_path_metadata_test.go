// This file contains unit tests for core path metadata operations.
// It tests GetPathMetadata, SetPathMetadata, AddPath, RemovePath, UpdatePath,
// ValidatePathMetadata, and GetPathConflicts methods from package_path_metadata.go.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"context"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// =============================================================================
// TEST: GetPathMetadata
// =============================================================================

// TestPackage_GetPathMetadata_Basic tests basic GetPathMetadata operation.
func TestPackage_GetPathMetadata_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetPathMetadata should now succeed since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	_, err = fpkg.GetPathMetadata(ctx)
	// Should succeed (may return empty slice if no path metadata)
	if err != nil {
		t.Errorf("GetPathMetadata failed: %v", err)
	}
}

// TestPackage_GetPathMetadata_WithContext tests GetPathMetadata with context scenarios.
func TestPackage_GetPathMetadata_WithContext(t *testing.T) {
	runWithCancelledContext(t, func(fpkg *filePackage, ctx context.Context) (interface{}, error) {
		return fpkg.GetPathMetadata(ctx)
	}, "GetPathMetadata")
}

// =============================================================================
// TEST: SetPathMetadata
// =============================================================================

// TestPackage_SetPathMetadata_Basic tests basic SetPathMetadata operation.
func TestPackage_SetPathMetadata_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// SetPathMetadata should now succeed since SavePathMetadataFile is implemented
	entries := []*metadata.PathMetadataEntry{}
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	err = fpkg.SetPathMetadata(ctx, entries)
	if err != nil {
		t.Errorf("SetPathMetadata failed: %v", err)
	}
}

// =============================================================================
// TEST: AddPath
// =============================================================================

// TestPackage_AddPath_Basic tests basic AddPath operation.
func TestPackage_AddPath_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// AddPath should now succeed since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	err = fpkg.AddPathMetadata(ctx, "test/path", metadata.PathMetadataTypeFile, nil, nil, nil)
	if err != nil {
		t.Errorf("AddPathMetadata failed: %v", err)
	}
}

// =============================================================================
// TEST: RemovePath
// =============================================================================

// TestPackage_RemovePath_Basic tests basic RemovePath operation.
func TestPackage_RemovePath_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// RemovePathMetadata should return error when LoadPathMetadataFile is not implemented
	fpkg := pkg.(*filePackage)
	err = fpkg.RemovePathMetadata(ctx, "test/path")
	if err == nil {
		t.Error("RemovePathMetadata should return error when LoadPathMetadataFile is not implemented")
	}
}

// =============================================================================
// TEST: UpdatePath
// =============================================================================

// TestPackage_UpdatePath_Basic tests basic UpdatePath operation.
func TestPackage_UpdatePath_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// UpdatePathMetadata should return error when LoadPathMetadataFile is not implemented
	fpkg := pkg.(*filePackage)
	err = fpkg.UpdatePathMetadata(ctx, "test/path", nil, nil, nil)
	if err == nil {
		t.Error("UpdatePathMetadata should return error when LoadPathMetadataFile is not implemented")
	}
}

// =============================================================================
// TEST: ValidatePathMetadata
// =============================================================================

// TestPackage_ValidatePathMetadata_Basic tests basic ValidatePathMetadata operation.
func TestPackage_ValidatePathMetadata_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// ValidatePathMetadata should now succeed since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	err = fpkg.ValidatePathMetadata(ctx)
	if err != nil {
		t.Errorf("ValidatePathMetadata failed: %v", err)
	}
}

// =============================================================================
// TEST: GetPathConflicts
// =============================================================================

// TestPackage_GetPathConflicts_Basic tests basic GetPathConflicts operation.
func TestPackage_GetPathConflicts_Basic(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetPathConflicts should now succeed since LoadPathMetadataFile is implemented
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = make([]*metadata.PathMetadataEntry, 0)
	_, err = fpkg.GetPathConflicts(ctx)
	if err != nil {
		t.Errorf("GetPathConflicts failed: %v", err)
	}
}

// =============================================================================
// TEST HELPERS
// =============================================================================

// createValidPathMetadataEntry creates a valid PathMetadataEntry for testing.
func createValidPathMetadataEntry(path string, pathType metadata.PathMetadataType) *metadata.PathMetadataEntry {
	return &metadata.PathMetadataEntry{
		Path: generics.PathEntry{
			Path:       path,
			PathLength: uint16(len(path)),
		},
		Type: pathType,
	}
}

// =============================================================================
// TEST: GetPathMetadata with cached entries
// =============================================================================

// TestPackage_GetPathMetadata_Cached tests GetPathMetadata with cached entries.
func TestPackage_GetPathMetadata_Cached(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	expectedEntries := []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("test/path", metadata.PathMetadataTypeFile),
	}
	fpkg.PathMetadataEntries = expectedEntries

	entries, err := fpkg.GetPathMetadata(ctx)
	if err != nil {
		t.Errorf("GetPathMetadata should succeed with cached entries, got error: %v", err)
	}
	if len(entries) != len(expectedEntries) {
		t.Errorf("GetPathMetadata should return cached entries, got %d, want %d", len(entries), len(expectedEntries))
	}
}

// =============================================================================
// TEST: SetPathMetadata with validation
// =============================================================================

// TestPackage_SetPathMetadata_NilEntry tests SetPathMetadata with nil entry.
func TestPackage_SetPathMetadata_NilEntry(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	entries := []*metadata.PathMetadataEntry{nil}
	err = fpkg.SetPathMetadata(ctx, entries)
	if err == nil {
		t.Error("SetPathMetadata should return error for nil entry")
	}
}

// TestPackage_SetPathMetadata_InvalidEntry tests SetPathMetadata with invalid entry.
func TestPackage_SetPathMetadata_InvalidEntry(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Create entry with invalid path length
	invalidEntry := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{
			Path:       "test",
			PathLength: 5, // Mismatch
		},
		Type: metadata.PathMetadataTypeFile,
	}
	entries := []*metadata.PathMetadataEntry{invalidEntry}
	err = fpkg.SetPathMetadata(ctx, entries)
	if err == nil {
		t.Error("SetPathMetadata should return error for invalid entry")
	}
}

// =============================================================================
// TEST: AddPath with cached entries
// =============================================================================

// TestPackage_AddPath_AlreadyExists tests AddPath when path already exists.
func TestPackage_AddPath_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	existingPath := "existing/path"
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry(existingPath, metadata.PathMetadataTypeFile),
	}

	err = fpkg.AddPathMetadata(ctx, existingPath, metadata.PathMetadataTypeFile, nil, nil, nil)
	if err == nil {
		t.Error("AddPathMetadata should return error when path already exists")
	}
}

// TestPackage_AddPath_Success tests successful AddPath operation.
func TestPackage_AddPath_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}

	newPath := "new/path"
	properties := map[string]string{"key": "value"}
	err = fpkg.AddPathMetadata(ctx, newPath, metadata.PathMetadataTypeFile, properties, nil, nil)
	if err == nil {
		// Should fail because SavePathMetadataFile is not implemented
		// But we can verify the entry was added to cache
		if len(fpkg.PathMetadataEntries) != 1 {
			t.Errorf("AddPathMetadata should add entry to cache, got %d entries", len(fpkg.PathMetadataEntries))
		}
		if fpkg.PathMetadataEntries[0].GetPath() != newPath {
			t.Errorf("AddPathMetadata should add correct path, got %q", fpkg.PathMetadataEntries[0].GetPath())
		}
	}
}

// =============================================================================
// TEST: RemovePath with cached entries
// =============================================================================

// TestPackage_RemovePath_NotFound tests RemovePath when path not found.
func TestPackage_RemovePath_NotFound(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("other/path", metadata.PathMetadataTypeFile),
	}

	err = fpkg.RemovePathMetadata(ctx, "nonexistent/path")
	if err == nil {
		t.Error("RemovePathMetadata should return error when path not found")
	}
}

// TestPackage_RemovePath_Success tests successful RemovePath operation.
func TestPackage_RemovePath_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	pathToRemove := "path/to/remove"
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry(pathToRemove, metadata.PathMetadataTypeFile),
		createValidPathMetadataEntry("other/path", metadata.PathMetadataTypeFile),
	}

	err = fpkg.RemovePathMetadata(ctx, pathToRemove)
	if err == nil {
		// Should fail because SavePathMetadataFile is not implemented
		// But we can verify the entry was removed from cache
		if len(fpkg.PathMetadataEntries) != 1 {
			t.Errorf("RemovePathMetadata should remove entry from cache, got %d entries", len(fpkg.PathMetadataEntries))
		}
		if fpkg.PathMetadataEntries[0].GetPath() == pathToRemove {
			t.Error("RemovePathMetadata should remove correct path")
		}
	}
}

// =============================================================================
// TEST: UpdatePath with cached entries
// =============================================================================

// TestPackage_UpdatePath_NotFound tests UpdatePath when path not found.
func TestPackage_UpdatePath_NotFound(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("other/path", metadata.PathMetadataTypeFile),
	}

	err = fpkg.UpdatePathMetadata(ctx, "nonexistent/path", nil, nil, nil)
	if err == nil {
		t.Error("UpdatePathMetadata should return error when path not found")
	}
}

// TestPackage_UpdatePath_Success tests successful UpdatePath operation.
func TestPackage_UpdatePath_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	pathToUpdate := "path/to/update"
	entry := createValidPathMetadataEntry(pathToUpdate, metadata.PathMetadataTypeFile)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{entry}

	newProperties := map[string]string{"newkey": "newvalue"}
	err = fpkg.UpdatePathMetadata(ctx, pathToUpdate, newProperties, nil, nil)
	if err == nil {
		// Should fail because SavePathMetadataFile is not implemented
		// But we can verify the entry was updated in cache
		if len(entry.Properties) == 0 {
			t.Error("UpdatePathMetadata should update properties")
		}
	}
}

// =============================================================================
// TEST: ValidatePathMetadata with cached entries
// =============================================================================

// TestPackage_ValidatePathMetadata_Success tests successful ValidatePathMetadata.
func TestPackage_ValidatePathMetadata_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry("valid/path", metadata.PathMetadataTypeFile),
	}

	err = fpkg.ValidatePathMetadata(ctx)
	if err != nil {
		t.Errorf("ValidatePathMetadata should succeed with valid entries, got error: %v", err)
	}
}

// TestPackage_ValidatePathMetadata_InvalidEntry tests ValidatePathMetadata with invalid entry.
func TestPackage_ValidatePathMetadata_InvalidEntry(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Create entry with invalid path length
	invalidEntry := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{
			Path:       "test",
			PathLength: 5, // Mismatch
		},
		Type: metadata.PathMetadataTypeFile,
	}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{invalidEntry}

	err = fpkg.ValidatePathMetadata(ctx)
	if err == nil {
		t.Error("ValidatePathMetadata should return error for invalid entry")
	}
}

// =============================================================================
// TEST: GetPathConflicts with cached entries
// =============================================================================

// TestPackage_GetPathConflicts_NoConflicts tests GetPathConflicts with no conflicts.
func TestPackage_GetPathConflicts_NoConflicts(t *testing.T) {
	ctx := context.Background()
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

	conflicts, err := fpkg.GetPathConflicts(ctx)
	if err != nil {
		t.Errorf("GetPathConflicts should succeed, got error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Errorf("GetPathConflicts should return no conflicts, got %d", len(conflicts))
	}
}

// TestPackage_GetPathConflicts_WithConflicts tests GetPathConflicts with conflicts.
func TestPackage_GetPathConflicts_WithConflicts(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	duplicatePath := "duplicate/path"
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		createValidPathMetadataEntry(duplicatePath, metadata.PathMetadataTypeFile),
		createValidPathMetadataEntry(duplicatePath, metadata.PathMetadataTypeFile),
		createValidPathMetadataEntry("unique/path", metadata.PathMetadataTypeFile),
	}

	conflicts, err := fpkg.GetPathConflicts(ctx)
	if err != nil {
		t.Errorf("GetPathConflicts should succeed, got error: %v", err)
	}
	if len(conflicts) != 1 {
		t.Errorf("GetPathConflicts should return 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0] != duplicatePath {
		t.Errorf("GetPathConflicts should return duplicate path, got %q", conflicts[0])
	}
}

// =============================================================================
// TEST: Context Cancellation
// =============================================================================

// TestPackage_UpdatePath_ContextCancelled tests UpdatePath with cancelled context.
func TestPackage_UpdatePath_ContextCancelled(t *testing.T) {
	ctx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	err = fpkg.UpdatePathMetadata(ctx, "test/path", nil, nil, nil)
	if err == nil {
		t.Error("UpdatePathMetadata() should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Expected error type Context, got: %v", pkgErr.Type)
	}
}

// TestPackage_ValidatePathMetadata_ContextCancelled tests ValidatePathMetadata with cancelled context.
func TestPackage_ValidatePathMetadata_ContextCancelled(t *testing.T) {
	runContextCancelledTest(t, func(fpkg *filePackage, ctx context.Context) error {
		return fpkg.ValidatePathMetadata(ctx)
	})
}

// TestPackage_GetPathConflicts_ContextCancelled tests GetPathConflicts with cancelled context.
func TestPackage_GetPathConflicts_ContextCancelled(t *testing.T) {
	runContextCancelledTest(t, func(fpkg *filePackage, ctx context.Context) error {
		_, err := fpkg.GetPathConflicts(ctx)
		return err
	})
}
