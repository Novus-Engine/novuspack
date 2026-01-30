// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains unit tests for package reader operations:
// GetInfo, GetMetadata, ReadFile, ListFiles, and Validate.
package novus_package

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/fileformat/testutil"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestPackage_GetInfo_Basic tests basic package information retrieval.
func TestPackage_GetInfo_Basic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Test: GetInfo should return package information
	info, err := pkg.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() failed: %v", err)
	}
	if info == nil {
		t.Fatal("GetInfo() returned nil")
	}
}

// TestPackage_GetInfo_WithContext tests GetInfo with context scenarios.
func TestPackage_GetInfo_WithContext(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetInfo is a pure in-memory operation and does not accept context
	info, err := pkg.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() failed: %v", err)
	}
	if info == nil {
		t.Error("GetInfo() returned nil")
	}
}

// TestPackage_GetInfo_AfterCreate tests GetInfo after Create and Open.
func TestPackage_GetInfo_AfterCreate(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Open the package
	pkg, err = OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() failed: %v", err)
	}

	if info == nil {
		t.Fatal("GetInfo() returned nil info")
	}
}

// TestPackage_GetInfo_OnNew tests GetInfo on newly created package.
func TestPackage_GetInfo_OnNew(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// GetInfo on new package should succeed (metadata is initialized)
	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() failed on newly created package: %v", err)
	}

	if info == nil {
		t.Fatal("GetInfo() returned nil")
	}

	// Verify format version is set
	if info.FormatVersion == 0 {
		t.Error("FormatVersion should be non-zero for newly created package")
	}
}

// TestPackage_GetInfo_OnClosedPackage tests GetInfo after Close.
func TestPackage_GetInfo_OnClosedPackage(t *testing.T) {
	ctx := context.Background()
	// Open a package to ensure metadata is loaded, then Close it.
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	// Close the package (in-memory metadata should remain available)
	if err := pkg.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// GetInfo on closed package should succeed if metadata is still in memory
	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() should succeed on a closed package with cached metadata, got error: %v", err)
	}
	if info == nil {
		t.Fatal("GetInfo() returned nil info")
	}
}

// TestPackage_GetInfo_WithCancelledContext tests GetInfo (no longer uses context).
func TestPackage_GetInfo_WithCancelledContext(t *testing.T) {
	ctx := context.Background()
	// Create and open a package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg2, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// GetInfo is a pure in-memory operation and does not accept context
	info, err := pkg2.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() failed: %v", err)
	}
	if info == nil {
		t.Error("GetInfo() returned nil")
	}
}

// TestPackage_GetInfo_WithNilInfo tests GetInfo when Info field is nil.
func TestPackage_GetInfo_WithNilInfo(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Manually set Info to nil to test error handling
	fpkg := pkg.(*filePackage)
	fpkg.Info = nil

	_, err = pkg.GetInfo()
	if err == nil {
		t.Error("GetInfo() should return error when Info is nil")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_GetMetadata tests the GetMetadata method.
func TestPackage_GetMetadata(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// GetMetadata should return PackageMetadata
	pm, err := pkg.GetMetadata()
	if err != nil {
		t.Fatalf("GetMetadata() failed: %v", err)
	}
	if pm == nil {
		t.Fatal("GetMetadata() returned nil")
	}
	if pm.PackageInfo == nil {
		t.Error("PackageMetadata.PackageInfo should not be nil")
	}
}

// TestPackage_ReadFile tests the ReadFile method.
func TestPackage_ReadFile(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// ReadFile on a package that hasn't been opened should return validation error
	data, err := pkg.ReadFile(ctx, "test.txt")
	if err == nil {
		t.Fatal("ReadFile() should return error when package is not open")
	}
	if data != nil {
		t.Errorf("ReadFile() should return nil data when error occurs")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_ListFiles tests the ListFiles method.
func TestPackage_ListFiles(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// ListFiles on an empty package should return empty list
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}
	if files == nil {
		t.Fatal("ListFiles() should not return nil")
	}
	if len(files) != 0 {
		t.Errorf("ListFiles() on empty package should return empty list, got %d files", len(files))
	}
}

// TestPackage_ListFiles_WithFiles tests ListFiles with files in package.
// Note: This test is simplified to test the basic functionality.
// Full file entry loading is tested in OpenPackage tests.
func TestPackage_ListFiles_WithFiles(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Create a minimal package file
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open package
	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// ListFiles should return empty list for empty package
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}
	if files == nil {
		t.Fatal("ListFiles() should not return nil")
	}
	// Empty package should return empty list
	if len(files) != 0 {
		t.Errorf("ListFiles() on empty package should return empty list, got %d files", len(files))
	}
}

// TestPackage_ListFiles_ClosedPackage tests ListFiles on closed package.
func TestPackage_ListFiles_ClosedPackage(t *testing.T) {
	ctx := context.Background()
	// Open a package to ensure metadata is loaded, then Close it.
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	// Close the package (in-memory metadata should remain available)
	if err := pkg.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// ListFiles on closed package should succeed if metadata is still in memory
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() should succeed on a closed package with cached metadata, got error: %v", err)
	}
	if files == nil {
		t.Fatal("ListFiles() should not return nil")
	}
}

// TestPackage_ListFiles_StableSorting tests that ListFiles returns stable results.
func TestPackage_ListFiles_StableSorting(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Create a minimal package file
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open package
	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Call ListFiles multiple times
	files1, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}

	files2, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}

	// Results should be identical
	if len(files1) != len(files2) {
		t.Errorf("ListFiles() returned different lengths: %d vs %d", len(files1), len(files2))
	}

	for i := range files1 {
		if files1[i].PrimaryPath != files2[i].PrimaryPath {
			t.Errorf("ListFiles() returned different order at index %d: %q vs %q", i, files1[i].PrimaryPath, files2[i].PrimaryPath)
		}
	}
}

// TestPackage_ListFiles_FileInfoFields tests that FileInfo contains all required fields.
func TestPackage_ListFiles_FileInfoFields(t *testing.T) {
	// Setup: Create a new package in memory
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true                                         // Mark package as open
	fpkg.FilePath = "/test/package.nvpk"                       // Set package path
	fpkg.Info = metadata.NewPackageInfo()                      // Initialize Info so metadata is considered loaded
	fpkg.FileEntries = []*metadata.FileEntry{}                 // Initialize FileEntries slice
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{} // Initialize PathMetadataEntries

	// Create test FileEntry directly with all properties we want to verify
	testEntry := &metadata.FileEntry{
		FileID:          12345,
		Type:            1000, // Custom file type
		OriginalSize:    4096,
		StoredSize:      1024, // Compressed size
		CompressionType: 1,    // Zstd
		EncryptionType:  0,    // Not encrypted
		RawChecksum:     0x12345678,
		StoredChecksum:  0x87654321,
		PathCount:       2, // Multi-path file
		FileVersion:     5,
		MetadataVersion: 3,
		Paths: []generics.PathEntry{
			{Path: "/textures/diffuse.dds"},
			{Path: "/assets/texture.dds"},
		},
		OptionalData: []metadata.OptionalDataEntry{
			{
				DataType: metadata.OptionalDataTagsData,
				Data:     []byte{0x01}, // Non-empty tags data
			},
		},
	}

	// Add the test entry to package
	fpkg.FileEntries = append(fpkg.FileEntries, testEntry)

	// Test: Call ListFiles
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}

	// Find our test file
	var testFileInfo *FileInfo
	for i := range files {
		if files[i].FileID == 12345 {
			testFileInfo = &files[i]
			break
		}
	}

	if testFileInfo == nil {
		t.Fatal("Test file not found in ListFiles() results")
	}

	// Verify all FileInfo fields
	t.Run("BasicIdentification", func(t *testing.T) {
		// PrimaryPath should be first lexicographically (without leading /)
		if testFileInfo.PrimaryPath != "assets/texture.dds" {
			t.Errorf("PrimaryPath = %q, want %q", testFileInfo.PrimaryPath, "assets/texture.dds")
		}

		// Paths array should contain all paths (sorted, without leading /)
		if len(testFileInfo.Paths) != 2 {
			t.Fatalf("len(Paths) = %d, want 2", len(testFileInfo.Paths))
		}
		wantPaths := []string{"assets/texture.dds", "textures/diffuse.dds"}
		for i, want := range wantPaths {
			if testFileInfo.Paths[i] != want {
				t.Errorf("Paths[%d] = %q, want %q", i, testFileInfo.Paths[i], want)
			}
		}

		// FileID should match
		if testFileInfo.FileID != 12345 {
			t.Errorf("FileID = %d, want 12345", testFileInfo.FileID)
		}

		// FileType should match
		if testFileInfo.FileType != 1000 {
			t.Errorf("FileType = %d, want 1000", testFileInfo.FileType)
		}

		// FileTypeName should be populated (currently "Unknown" until type system is implemented)
		if testFileInfo.FileTypeName == "" {
			t.Error("FileTypeName is empty, should be populated")
		}
	})

	t.Run("SizeInformation", func(t *testing.T) {
		// Size should be original size
		if testFileInfo.Size != 4096 {
			t.Errorf("Size = %d, want 4096", testFileInfo.Size)
		}

		// StoredSize should be compressed size
		if testFileInfo.StoredSize != 1024 {
			t.Errorf("StoredSize = %d, want 1024", testFileInfo.StoredSize)
		}
	})

	t.Run("ProcessingStatus", func(t *testing.T) {
		// IsCompressed should be true (CompressionType != 0)
		if !testFileInfo.IsCompressed {
			t.Error("IsCompressed = false, want true")
		}

		// IsEncrypted should be false
		if testFileInfo.IsEncrypted {
			t.Error("IsEncrypted = true, want false")
		}

		// CompressionType should be Zstd (1)
		if testFileInfo.CompressionType != 1 {
			t.Errorf("CompressionType = %d, want 1", testFileInfo.CompressionType)
		}
	})

	t.Run("ContentVerification", func(t *testing.T) {
		// RawChecksum should match
		if testFileInfo.RawChecksum != 0x12345678 {
			t.Errorf("RawChecksum = 0x%08X, want 0x12345678", testFileInfo.RawChecksum)
		}

		// StoredChecksum should match
		if testFileInfo.StoredChecksum != 0x87654321 {
			t.Errorf("StoredChecksum = 0x%08X, want 0x87654321", testFileInfo.StoredChecksum)
		}
	})

	t.Run("MultiPathSupport", func(t *testing.T) {
		// PathCount should match number of paths
		if testFileInfo.PathCount != 2 {
			t.Errorf("PathCount = %d, want 2", testFileInfo.PathCount)
		}
	})

	t.Run("VersionTracking", func(t *testing.T) {
		// FileVersion should match
		if testFileInfo.FileVersion != 5 {
			t.Errorf("FileVersion = %d, want 5", testFileInfo.FileVersion)
		}

		// MetadataVersion should match
		if testFileInfo.MetadataVersion != 3 {
			t.Errorf("MetadataVersion = %d, want 3", testFileInfo.MetadataVersion)
		}
	})

	t.Run("MetadataIndicators", func(t *testing.T) {
		// HasTags should be true (we added OptionalData with TagsData)
		if !testFileInfo.HasTags {
			t.Error("HasTags = false, want true")
		}
	})
}

// TestPackage_ListFiles_FileInfoFields_Uncompressed tests FileInfo for uncompressed files.
func TestPackage_ListFiles_FileInfoFields_Uncompressed(t *testing.T) {
	// Setup: Create a new package in memory
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true                                         // Mark package as open
	fpkg.FilePath = "/test/package_uncompressed.nvpk"          // Set package path
	fpkg.Info = metadata.NewPackageInfo()                      // Initialize Info so metadata is considered loaded
	fpkg.FileEntries = []*metadata.FileEntry{}                 // Initialize FileEntries slice
	fpkg.PathMetadataEntries = []*metadata.PathMetadataEntry{} // Initialize PathMetadataEntries

	testEntry := &metadata.FileEntry{
		FileID:          99999,
		Type:            2000,
		OriginalSize:    2048,
		StoredSize:      2048, // Same as original (not compressed)
		CompressionType: 0,    // No compression
		EncryptionType:  0,    // No encryption
		RawChecksum:     0xAABBCCDD,
		StoredChecksum:  0xAABBCCDD, // Same as raw (not processed)
		PathCount:       1,
		FileVersion:     1,
		MetadataVersion: 1,
		Paths: []generics.PathEntry{
			{Path: "/data/config.json"},
		},
		OptionalData: []metadata.OptionalDataEntry{}, // No tags
	}

	fpkg.FileEntries = append(fpkg.FileEntries, testEntry)

	// Test: Call ListFiles
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}

	// Find our test file
	var testFileInfo *FileInfo
	for i := range files {
		if files[i].FileID == 99999 {
			testFileInfo = &files[i]
			break
		}
	}

	if testFileInfo == nil {
		t.Fatal("Test file not found in ListFiles() results")
	}

	// Verify uncompressed file properties
	if testFileInfo.IsCompressed {
		t.Error("IsCompressed = true, want false for uncompressed file")
	}

	if testFileInfo.IsEncrypted {
		t.Error("IsEncrypted = true, want false for unencrypted file")
	}

	if testFileInfo.CompressionType != 0 {
		t.Errorf("CompressionType = %d, want 0 for uncompressed file", testFileInfo.CompressionType)
	}

	if testFileInfo.Size != testFileInfo.StoredSize {
		t.Errorf("Size (%d) != StoredSize (%d), should be equal for uncompressed file", testFileInfo.Size, testFileInfo.StoredSize)
	}

	if testFileInfo.RawChecksum != testFileInfo.StoredChecksum {
		t.Errorf("RawChecksum (0x%08X) != StoredChecksum (0x%08X), should be equal for unprocessed file",
			testFileInfo.RawChecksum, testFileInfo.StoredChecksum)
	}

	if testFileInfo.HasTags {
		t.Error("HasTags = true, want false for file without tags")
	}

	if testFileInfo.PathCount != 1 {
		t.Errorf("PathCount = %d, want 1 for single-path file", testFileInfo.PathCount)
	}

	if len(testFileInfo.Paths) != 1 {
		t.Errorf("len(Paths) = %d, want 1 for single-path file", len(testFileInfo.Paths))
	}
}

// TestPackage_ListFiles_EdgeCases tests edge cases in ListFiles for coverage.
func TestPackage_ListFiles_EdgeCases(t *testing.T) {
	// Setup: Create a new package in memory
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.FilePath = "/test/package_edgecases.nvpk"
	fpkg.Info = metadata.NewPackageInfo()
	fpkg.FileEntries = []*metadata.FileEntry{}

	// Test 1: Nil entry (should be skipped)
	fpkg.FileEntries = append(fpkg.FileEntries, nil)

	// Test 2: Entry with no paths (should be skipped)
	entryNoPaths := &metadata.FileEntry{
		FileID:       11111,
		Type:         100,
		OriginalSize: 1000,
		StoredSize:   1000,
		PathCount:    0,
		Paths:        []generics.PathEntry{},
	}
	fpkg.FileEntries = append(fpkg.FileEntries, entryNoPaths)

	// Test 3: Entry with invalid paths (should be skipped)
	entryInvalidPaths := &metadata.FileEntry{
		FileID:       22222,
		Type:         200,
		OriginalSize: 2000,
		StoredSize:   2000,
		PathCount:    1,
		Paths: []generics.PathEntry{
			{Path: ""}, // Empty path is invalid
		},
	}
	fpkg.FileEntries = append(fpkg.FileEntries, entryInvalidPaths)

	// Test 4: Entry with OptionalData but no tags
	entryNoTags := &metadata.FileEntry{
		FileID:       33333,
		Type:         300,
		OriginalSize: 3000,
		StoredSize:   3000,
		PathCount:    1,
		Paths: []generics.PathEntry{
			{Path: "/data/file1.txt"},
		},
		OptionalData: []metadata.OptionalDataEntry{
			{DataType: 0x01, Data: []byte{0x01, 0x02}}, // Not tags data
		},
	}
	fpkg.FileEntries = append(fpkg.FileEntries, entryNoTags)

	// Test 5: Entry with tags in middle of OptionalData
	entryTagsMiddle := &metadata.FileEntry{
		FileID:       44444,
		Type:         400,
		OriginalSize: 4000,
		StoredSize:   4000,
		PathCount:    1,
		Paths: []generics.PathEntry{
			{Path: "/data/file2.txt"},
		},
		OptionalData: []metadata.OptionalDataEntry{
			{DataType: 0x01, Data: []byte{0x01}},                          // Other data
			{DataType: metadata.OptionalDataTagsData, Data: []byte{0x01}}, // Tags data
			{DataType: 0x02, Data: []byte{0x02}},                          // More data
		},
	}
	fpkg.FileEntries = append(fpkg.FileEntries, entryTagsMiddle)

	// Test 6: Entry with empty tags data (should not set HasTags)
	entryEmptyTags := &metadata.FileEntry{
		FileID:       55555,
		Type:         500,
		OriginalSize: 5000,
		StoredSize:   5000,
		PathCount:    1,
		Paths: []generics.PathEntry{
			{Path: "/data/file3.txt"},
		},
		OptionalData: []metadata.OptionalDataEntry{
			{DataType: metadata.OptionalDataTagsData, Data: []byte{}}, // Empty tags
		},
	}
	fpkg.FileEntries = append(fpkg.FileEntries, entryEmptyTags)

	// Test: Call ListFiles
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}

	// Verify results
	// Should have 3 valid files (entryNoTags, entryTagsMiddle, entryEmptyTags)
	if len(files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(files))
	}

	// Verify file 1 (entryNoTags): FileID 33333, HasTags should be false
	var file1 *FileInfo
	for i := range files {
		if files[i].FileID == 33333 {
			file1 = &files[i]
			break
		}
	}
	if file1 == nil {
		t.Fatal("File with FileID 33333 not found")
	}
	if file1.HasTags {
		t.Error("File 33333: HasTags should be false (has non-tag OptionalData)")
	}
	if file1.PrimaryPath != "data/file1.txt" {
		t.Errorf("File 33333: PrimaryPath = %q, want %q", file1.PrimaryPath, "data/file1.txt")
	}

	// Verify file 2 (entryTagsMiddle): FileID 44444, HasTags should be true
	var file2 *FileInfo
	for i := range files {
		if files[i].FileID == 44444 {
			file2 = &files[i]
			break
		}
	}
	if file2 == nil {
		t.Fatal("File with FileID 44444 not found")
	}
	if !file2.HasTags {
		t.Error("File 44444: HasTags should be true (has tags in middle of OptionalData)")
	}

	// Verify file 3 (entryEmptyTags): FileID 55555, HasTags should be false
	var file3 *FileInfo
	for i := range files {
		if files[i].FileID == 55555 {
			file3 = &files[i]
			break
		}
	}
	if file3 == nil {
		t.Fatal("File with FileID 55555 not found")
	}
	if file3.HasTags {
		t.Error("File 55555: HasTags should be false (empty tags data)")
	}
}

// TestPackage_ListFiles_NoMetadataLoaded tests ListFiles when metadata is not loaded.
func TestPackage_ListFiles_NoMetadataLoaded(t *testing.T) {
	tests := []struct {
		name         string
		setupPackage func() *filePackage
		wantErr      bool
	}{
		{
			name: "Info is nil",
			setupPackage: func() *filePackage {
				pkg, _ := NewPackage()
				fpkg := pkg.(*filePackage)
				fpkg.Info = nil
				fpkg.FilePath = "/test/package.nvpk"
				fpkg.FileEntries = []*metadata.FileEntry{}
				return fpkg
			},
			wantErr: true,
		},
		{
			name: "FilePath is empty",
			setupPackage: func() *filePackage {
				pkg, _ := NewPackage()
				fpkg := pkg.(*filePackage)
				fpkg.Info = metadata.NewPackageInfo()
				fpkg.FilePath = ""
				fpkg.FileEntries = []*metadata.FileEntry{}
				return fpkg
			},
			wantErr: false, // Changed: FilePath can be empty for newly created packages
		},
		{
			name: "FileEntries is nil",
			setupPackage: func() *filePackage {
				pkg, _ := NewPackage()
				fpkg := pkg.(*filePackage)
				fpkg.Info = metadata.NewPackageInfo()
				fpkg.FilePath = "/test/package.nvpk"
				fpkg.FileEntries = nil
				return fpkg
			},
			wantErr: true,
		},
		{
			name: "All fields valid",
			setupPackage: func() *filePackage {
				pkg, _ := NewPackage()
				fpkg := pkg.(*filePackage)
				fpkg.Info = metadata.NewPackageInfo()
				fpkg.FilePath = "/sp/package.nvpk"
				fpkg.FileEntries = []*metadata.FileEntry{}
				return fpkg
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpkg := tt.setupPackage()
			pkg := Package(fpkg)
			defer func() { _ = pkg.Close() }()

			files, err := pkg.ListFiles()
			if tt.wantErr {
				if err == nil {
					t.Errorf("ListFiles() expected error, got nil")
				}
				if files != nil {
					t.Errorf("ListFiles() expected nil result on error, got %v", files)
				}
			} else {
				if err != nil {
					t.Errorf("ListFiles() unexpected error: %v", err)
				}
				if files == nil {
					t.Error("ListFiles() should not return nil on success")
				}
			}
		})
	}
}

// TestPackage_Validate_Basic tests basic package validation.
func TestPackage_Validate_Basic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create a valid package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open and validate
	pkg2, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// Test: Validate should succeed for valid package
	err = pkg2.Validate(ctx)
	if err != nil {
		t.Errorf("Validate() failed for valid package: %v", err)
	}
}

// TestPackage_Validate_RequiresOpen tests that Validate requires open package.
func TestPackage_Validate_RequiresOpen(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	_ = pkg.Close()

	// Test: Validate on closed package should fail
	err = pkg.Validate(ctx)
	if err == nil {
		t.Error("Validate should fail on closed package")
	}
}

// TestPackage_Validate_WithContext tests Validate with context scenarios.
func TestPackage_Validate_WithContext(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(context.Background(), pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg2, err := OpenPackage(context.Background(), pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// Test with cancelled context
	cancelledCtx := testhelpers.CancelledContext()
	err = pkg2.Validate(cancelledCtx)
	if err == nil {
		t.Error("Validate should fail with cancelled context")
	}
}

// TestPackage_Validate_WhenClosed tests Validate on closed package.
func TestPackage_Validate_WhenClosed(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Validate on closed package
	err = pkg.Validate(ctx)
	if err == nil {
		t.Error("Validate() should return error when package is closed")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("Validate() should return PackageError")
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestPackage_Validate_WithOpenPackage tests Validate on a properly opened package.
func TestPackage_Validate_WithOpenPackage(t *testing.T) {
	ctx := context.Background()
	// Create and open a package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	pkg2, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// Validate should succeed
	if err := pkg2.Validate(ctx); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}
}

// TestPackage_Validate_WithNilHeader tests Validate when header is nil.
func TestPackage_Validate_WithNilHeader(t *testing.T) {
	ctx := context.Background()
	// Create a package and manually set header to nil
	// Note: Need to use type assertion to access internal fields
	pkg, _ := NewPackage()
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = nil

	err := pkg.Validate(ctx)
	if err == nil {
		t.Error("Validate() should return error when header is nil")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("Validate() should return PackageError")
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestPackage_Validate_WithNilIndex tests Validate with nil index (should succeed).
func TestPackage_Validate_WithNilIndex(t *testing.T) {
	ctx := context.Background()
	// Create a package with nil index
	// Note: Need to use type assertion to access internal fields
	pkg, _ := NewPackage()
	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	fpkg.index = nil

	// Validate should succeed (nil index is allowed)
	if err := pkg.Validate(ctx); err != nil {
		t.Errorf("Validate() should succeed with nil index: %v", err)
	}
}

// TestPackage_Validate_WithInvalidHeader tests Validate when header.Validate() fails.
func TestPackage_Validate_WithInvalidHeader(t *testing.T) {
	ctx := context.Background()
	// Create a package with invalid header
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	fpkg.isOpen = true
	fpkg.header = fileformat.NewPackageHeader()
	// Corrupt the header to make Validate() fail
	fpkg.header.FormatVersion = 999 // Invalid version that will fail validation
	fpkg.index = nil                // Set to nil to avoid index validation

	// Validate should fail
	err = pkg.Validate(ctx)
	if err == nil {
		t.Fatal("Validate() should fail with invalid header")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_Validate_WithInvalidIndex tests Validate when index is invalid.
func TestPackage_Validate_WithInvalidIndex(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create and open a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	if err := fpkg1.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Corrupt the index by setting invalid entry count
	fpkg := pkg.(*filePackage)
	if fpkg.index != nil {
		// Set entry count to max value which will fail validation
		fpkg.index.EntryCount = 0xFFFFFFFF
	}

	// Validate should fail
	err = pkg.Validate(ctx)
	if err == nil {
		t.Fatal("Validate() should fail with invalid index")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_Validate_WithInvalidIndexV2 tests Validate when index.Validate() fails.
func TestPackage_Validate_WithInvalidIndexV2(t *testing.T) {
	ctx := context.Background()
	// Create a package and manually set an invalid index
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Set up an invalid index that will fail validation
	// Create index with mismatched entry count
	// Need type assertion to access internal fields
	fpkg := pkg.(*filePackage)
	fpkg.index = &fileformat.FileIndex{
		EntryCount:       100,                       // Claims 100 entries
		FirstEntryOffset: 0,                         // Invalid offset
		Entries:          []fileformat.IndexEntry{}, // But has no entries
	}
	fpkg.isOpen = true

	// Validate should fail
	err = fpkg.Validate(ctx)
	if err == nil {
		t.Fatal("Validate() should fail with invalid index")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_Validate_WithCancelledContext tests Validate with cancelled context.
func TestPackage_Validate_WithCancelledContext(t *testing.T) {
	ctx := context.Background()
	// Create and open a package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	pkg2, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// Test with cancelled context
	cancelledCtx := testhelpers.CancelledContext()
	err = pkg2.Validate(cancelledCtx)
	if err == nil {
		t.Error("Validate() should fail with cancelled context")
	}
}

// TestPackage_ReadFile_ContextCancelDuringIO tests ReadFile with context cancelled during I/O.
func TestPackage_ReadFile_ContextCancelDuringIO(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Package is already "open" after NewPackage, no need to call Open()
	defer func() { _ = pkg.Close() }()

	// Add file using AddFileFromMemory (which sets IsDataLoaded=true)
	ctx := context.Background()
	testContent := []byte("test content for context cancellation")
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory() failed: %v", err)
	}

	// Cancel context and try to read
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = pkg.ReadFile(cancelledCtx, "/test.txt")
	if err == nil {
		t.Error("ReadFile() should fail with cancelled context")
	}
}

// TestPackage_ReadFile_SourceFileNil tests ReadFile when SourceFile is nil.
// This tests the case where data is loaded in memory (IsDataLoaded=true).
func TestPackage_ReadFile_SourceFileNil(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	ctx := context.Background()

	// Create package to open it (required for ReadFile)
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.Create(ctx, tmpPkg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Add file using AddFileFromMemory (which sets IsDataLoaded=true)
	testContent := []byte("test content with nil SourceFile")
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory() failed: %v", err)
	}

	// Verify entry has data loaded (ReadFile should use Data field, not SourceFile)
	if !entry.IsDataLoaded {
		t.Error("AddFileFromMemory should set IsDataLoaded=true")
	}

	// ReadFile should work when IsDataLoaded=true (uses Data field, not SourceFile)
	data, err := pkg.ReadFile(ctx, "/test.txt")
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("ReadFile() content mismatch: got %q, want %q", string(data), string(testContent))
	}
}

// TestPackage_ReadFile_SourceOffsetZero tests ReadFile when SourceOffset is zero.
// For AddFile, SourceOffset is set to 0 initially, but ReadFile requires SourceOffset > 0
// for files that aren't in memory. This test verifies AddFileFromMemory works correctly.
func TestPackage_ReadFile_SourceOffsetZero(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	ctx := context.Background()

	// Create package to open it (required for ReadFile)
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.Create(ctx, tmpPkg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Add file using AddFileFromMemory (which sets IsDataLoaded=true, avoiding SourceOffset check)
	testContent := []byte("test content")
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory() failed: %v", err)
	}

	// Verify SourceOffset is zero (initial value)
	if entry.SourceOffset != 0 {
		t.Errorf("SourceOffset = %d, want 0", entry.SourceOffset)
	}

	// ReadFile should work correctly because IsDataLoaded=true (bypasses SourceOffset check)
	data, err := pkg.ReadFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("ReadFile() content mismatch: got %q, want %q", string(data), string(testContent))
	}
}
