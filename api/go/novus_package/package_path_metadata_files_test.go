// This file contains unit tests for special metadata file operations.
// It tests SavePathMetadataFile, LoadPathMetadataFile, and UpdateSpecialMetadataFlags
// methods from package_path_metadata_files.go.
//
// Specification: api_metadata.md: 8.3 Special Metadata File Management

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
// TEST: SavePathMetadataFile
// =============================================================================

// TestPackage_SavePathMetadataFile_CreateNew tests SavePathMetadataFile creating a new special file.
// Expected: Should create FileEntry type 65001 with YAML data
func TestPackage_SavePathMetadataFile_CreateNew(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.FileEntries = make([]*metadata.FileEntry, 0)

	// Add PathMetadataEntries
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		{
			Path: generics.PathEntry{Path: "assets/", PathLength: 7},
			Type: metadata.PathMetadataTypeDirectory,
		},
	}

	// SavePathMetadataFile should create special file
	err = fpkg.SavePathMetadataFile(ctx)
	if err != nil {
		t.Fatalf("SavePathMetadataFile failed: %v", err)
	}

	// Verify special file was created
	specialFile, exists := fpkg.SpecialFiles[65001]
	if !exists {
		t.Fatal("Special file type 65001 should be created")
	}
	if specialFile.Type != 65001 {
		t.Errorf("Special file Type = %d, want 65001", specialFile.Type)
	}
	// Verify file path
	if len(specialFile.Paths) == 0 {
		t.Fatal("Special file should have path")
	}
	if specialFile.Paths[0].Path != "__NVPK_PATH_65001__.nvpkpath" {
		t.Errorf("Special file path = %s, want '__NVPK_PATH_65001__.nvpkpath'", specialFile.Paths[0].Path)
	}
	// Verify compression type is uncompressed (per spec)
	if specialFile.CompressionType != 0 {
		t.Errorf("CompressionType = %d, want 0 (uncompressed)", specialFile.CompressionType)
	}
	// Verify FileID is sequential (not 65001)
	if specialFile.FileID == 65001 {
		t.Error("FileID should be sequential, not 65001")
	}
	// Verify FileID is non-zero
	if specialFile.FileID == 0 {
		t.Error("FileID should be non-zero")
	}
	// Verify encryption type is none
	if specialFile.EncryptionType != 0 {
		t.Errorf("EncryptionType = %d, want 0 (none)", specialFile.EncryptionType)
	}
}

// TestPackage_SavePathMetadataFile_UpdateExisting tests SavePathMetadataFile updating existing special file.
// Expected: Should update existing FileEntry
func TestPackage_SavePathMetadataFile_UpdateExisting(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.FileEntries = make([]*metadata.FileEntry, 0)

	// Create existing special file with sequential FileID
	existingFile := metadata.NewFileEntry()
	existingFile.FileID = 1 // Sequential FileID
	existingFile.Type = 65001
	existingFile.Paths = []generics.PathEntry{
		{Path: "__NVPK_PATH_65001__.nvpkpath", PathLength: 27},
	}
	existingFile.Data = []byte("old data")
	fpkg.SpecialFiles[65001] = existingFile
	fpkg.FileEntries = append(fpkg.FileEntries, existingFile)

	// Add PathMetadataEntries
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		{
			Path: generics.PathEntry{Path: "newfile.txt", PathLength: 11},
			Type: metadata.PathMetadataTypeFile,
		},
	}

	// SavePathMetadataFile should update existing file
	err = fpkg.SavePathMetadataFile(ctx)
	if err != nil {
		t.Fatalf("SavePathMetadataFile failed: %v", err)
	}

	// Verify file was updated (data changed)
	if string(fpkg.SpecialFiles[65001].Data) == "old data" {
		t.Error("Special file data should be updated")
	}
}

// TestPackage_SavePathMetadataFile_EmptyEntries tests SavePathMetadataFile with empty PathMetadataEntries.
// Expected: Should remove special file if it exists
func TestPackage_SavePathMetadataFile_EmptyEntries(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.FileEntries = make([]*metadata.FileEntry, 0)

	// Create existing special file
	existingFile := metadata.NewFileEntry()
	existingFile.Type = 65001
	fpkg.SpecialFiles[65001] = existingFile
	fpkg.FileEntries = append(fpkg.FileEntries, existingFile)

	// Set PathMetadataEntries to empty
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}

	// SavePathMetadataFile should remove special file
	err = fpkg.SavePathMetadataFile(ctx)
	if err != nil {
		t.Fatalf("SavePathMetadataFile failed: %v", err)
	}

	// Verify special file was removed
	_, exists := fpkg.SpecialFiles[65001]
	if exists {
		t.Error("Special file should be removed when PathMetadataEntries is empty")
	}
}

// TestPackage_SavePathMetadataFile_YAMLStructure tests SavePathMetadataFile YAML structure.
// Expected: YAML should have correct structure
func TestPackage_SavePathMetadataFile_YAMLStructure(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.FileEntries = make([]*metadata.FileEntry, 0)

	// Add PathMetadataEntry
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{
		{
			Path: generics.PathEntry{Path: "test.txt", PathLength: 8},
			Type: metadata.PathMetadataTypeFile,
		},
	}

	// SavePathMetadataFile should create file with correct YAML structure
	err = fpkg.SavePathMetadataFile(ctx)
	if err != nil {
		t.Fatalf("SavePathMetadataFile failed: %v", err)
	}

	// Verify YAML data contains "paths:" key
	specialFile := fpkg.SpecialFiles[65001]
	if specialFile == nil {
		t.Fatal("Special file should exist")
	}
	if len(specialFile.Data) == 0 {
		t.Fatal("Special file should have data")
	}
	yamlStr := string(specialFile.Data)
	if !testhelpers.Contains(yamlStr, "paths:") {
		t.Error("YAML should contain 'paths:' key")
	}
}

// TestPackage_SavePathMetadataFile_CancelledContext tests SavePathMetadataFile with cancelled context.
// Expected: Should return context error
func TestPackage_SavePathMetadataFile_CancelledContext(t *testing.T) {
	cancelledCtx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()

	// SavePathMetadataFile should fail with cancelled context
	err = fpkg.SavePathMetadataFile(cancelledCtx)
	if err == nil {
		t.Fatal("SavePathMetadataFile should fail with cancelled context")
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
// TEST: LoadPathMetadataFile
// =============================================================================

// TestPackage_LoadPathMetadataFile_MissingFile tests LoadPathMetadataFile when no special file exists.
// Expected: Should succeed without error (path metadata is optional)
func TestPackage_LoadPathMetadataFile_MissingFile(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Initialize package state
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = nil

	// LoadPathMetadataFile should succeed when no special file exists
	err = fpkg.LoadPathMetadataFile(ctx)
	if err != nil {
		t.Errorf("LoadPathMetadataFile should succeed when no special file exists, got error: %v", err)
	}

	// PathMetadataEntries should remain nil or empty
	if fpkg.PathMetadataEntries == nil {
		// nil is acceptable
	} else if len(fpkg.PathMetadataEntries) != 0 {
		t.Errorf("PathMetadataEntries should be empty when no special file exists, got %d entries", len(fpkg.PathMetadataEntries))
	}
}

// TestPackage_LoadPathMetadataFile_ValidYAML tests LoadPathMetadataFile with valid YAML.
// Expected: Should parse YAML and populate PathMetadataEntries
// NOTE: This test is currently skipped because LoadPathMetadataFile internally calls
// ReadFile which requires a full package setup with file handle. A proper integration
// test would need to create an actual package file on disk.
func TestPackage_LoadPathMetadataFile_ValidYAML(t *testing.T) {
	t.Skip("Skipping - requires full package with file handle for ReadFile")

	// This test would need:
	// 1. Create actual package file on disk
	// 2. Add special file type 65001 with YAML content
	// 3. OpenPackage to get proper file handle
	// 4. Then test LoadPathMetadataFile
	//
	// For unit testing, we verify LoadPathMetadataFile is called correctly
	// in OpenPackage integration tests instead.
}

// TestPackage_LoadPathMetadataFile_MalformedYAML tests LoadPathMetadataFile with malformed YAML.
// Expected: Should return validation error
// NOTE: Skipped - requires full package setup with file handle
func TestPackage_LoadPathMetadataFile_MalformedYAML(t *testing.T) {
	t.Skip("Skipping - requires full package with file handle for ReadFile")
}

// TestPackage_LoadPathMetadataFile_EmptyPaths tests LoadPathMetadataFile with empty paths array.
// Expected: Should succeed with empty PathMetadataEntries
// NOTE: Skipped - requires full package setup with file handle
func TestPackage_LoadPathMetadataFile_EmptyPaths(t *testing.T) {
	t.Skip("Skipping - requires full package with file handle for ReadFile")
}

// TestPackage_LoadPathMetadataFile_CancelledContext tests LoadPathMetadataFile with cancelled context.
// Expected: Should return context error
func TestPackage_LoadPathMetadataFile_CancelledContext(t *testing.T) {
	cancelledCtx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true

	// LoadPathMetadataFile should fail with cancelled context
	err = fpkg.LoadPathMetadataFile(cancelledCtx)
	if err == nil {
		t.Fatal("LoadPathMetadataFile should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want ErrTypeContext", pkgErr.Type)
	}
}

// TestPackage_LoadPathMetadataFile_InvalidEntryValidation tests LoadPathMetadataFile with entry that fails validation.
// Expected: Should return validation error
// NOTE: Skipped - requires full package setup with file handle
func TestPackage_LoadPathMetadataFile_InvalidEntryValidation(t *testing.T) {
	t.Skip("Skipping - requires full package with file handle for ReadFile")
}

// =============================================================================
// TEST: UpdateSpecialMetadataFlags
// =============================================================================

// TestPackage_UpdateSpecialMetadataFlags_NoSpecialFiles tests UpdateSpecialMetadataFlags with no special files.
// Expected: Bit 6 should be cleared
func TestPackage_UpdateSpecialMetadataFlags_NoSpecialFiles(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()

	// Set flags initially to test clearing
	fpkg.header.Flags = fileformat.FlagHasSpecialMetadata | fileformat.FlagHasPerFileTags

	// UpdateSpecialMetadataFlags should clear bits 5 and 6
	err = fpkg.UpdateSpecialMetadataFlags(ctx)
	if err != nil {
		t.Errorf("UpdateSpecialMetadataFlags failed: %v", err)
	}

	// Verify bit 6 is cleared
	if fpkg.header.Flags&fileformat.FlagHasSpecialMetadata != 0 {
		t.Error("FlagHasSpecialMetadata should be cleared when no special files exist")
	}
	// Verify bit 5 is cleared
	if fpkg.header.Flags&fileformat.FlagHasPerFileTags != 0 {
		t.Error("FlagHasPerFileTags should be cleared when no path metadata tags exist")
	}
	// Verify Info flags match
	if fpkg.Info.HasMetadataFiles {
		t.Error("Info.HasMetadataFiles should be false")
	}
}

// TestPackage_UpdateSpecialMetadataFlags_WithSpecialFiles tests UpdateSpecialMetadataFlags with special files.
// Expected: Bit 6 should be set
func TestPackage_UpdateSpecialMetadataFlags_WithSpecialFiles(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()

	// Add a special file
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)
	specialFile := metadata.NewFileEntry()
	specialFile.Type = 65001
	fpkg.SpecialFiles[65001] = specialFile

	// Initialize PathMetadataEntries as empty (no tags)
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{}

	// UpdateSpecialMetadataFlags should set bit 6
	err = fpkg.UpdateSpecialMetadataFlags(ctx)
	if err != nil {
		t.Errorf("UpdateSpecialMetadataFlags failed: %v", err)
	}

	// Verify bit 6 is set
	if fpkg.header.Flags&fileformat.FlagHasSpecialMetadata == 0 {
		t.Error("FlagHasSpecialMetadata should be set when special files exist")
	}
	// Verify Info flag matches
	if !fpkg.Info.HasMetadataFiles {
		t.Error("Info.HasMetadataFiles should be true")
	}
}

// TestPackage_UpdateSpecialMetadataFlags_WithPerFileTags tests UpdateSpecialMetadataFlags with per-file tags.
// Expected: Bit 5 should be set
func TestPackage_UpdateSpecialMetadataFlags_WithPerFileTags(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.SpecialFiles = make(map[uint16]*metadata.FileEntry)

	// Add PathMetadataEntry with tags
	pme := &metadata.PathMetadataEntry{
		Path: generics.PathEntry{Path: "file.txt", PathLength: 8},
		Type: metadata.PathMetadataTypeFile,
		Properties: []*generics.Tag[any]{
			{Key: "category", Value: "test", Type: generics.TagValueTypeString},
		},
	}
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{pme}

	// UpdateSpecialMetadataFlags should set bit 5
	err = fpkg.UpdateSpecialMetadataFlags(ctx)
	if err != nil {
		t.Errorf("UpdateSpecialMetadataFlags failed: %v", err)
	}

	// Verify bit 5 is set
	if fpkg.header.Flags&fileformat.FlagHasPerFileTags == 0 {
		t.Error("FlagHasPerFileTags should be set when path metadata has tags")
	}
	// Note: HasPerFileTags is tracked via header flags but not currently exposed in PackageInfo
}

// TestPackage_UpdateSpecialMetadataFlags_CancelledContext tests UpdateSpecialMetadataFlags with cancelled context.
// Expected: Should return context error
func TestPackage_UpdateSpecialMetadataFlags_CancelledContext(t *testing.T) {
	cancelledCtx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.Info = metadata.NewPackageInfo()

	// UpdateSpecialMetadataFlags should fail with cancelled context
	err = fpkg.UpdateSpecialMetadataFlags(cancelledCtx)
	if err == nil {
		t.Fatal("UpdateSpecialMetadataFlags should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want ErrTypeContext", pkgErr.Type)
	}
}
