// Package testutil provides testing utilities for the fileformat package.
//
// This package contains helper functions for creating test fixtures and
// mock file format structures. It should only be imported by test files.
package testutil

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
)

// CreateTestPackageFile creates a minimal valid package file for testing.
//
// This helper creates a package file with:
//   - Valid NovusPack header with magic number
//   - Empty file index
//   - Proper offsets and sizes
//
// The created package can be opened with OpenPackage() for testing.
//
// This function is exported for use in test files across the v1 API packages.
//
// Parameters:
//   - t: Testing context (will call t.Fatalf on errors)
//   - path: File path where the package should be created
//
// Example usage:
//
//	import "github.com/novus-engine/novuspack/api/go/fileformat/testutil"
//
//	tmpDir := t.TempDir()
//	pkgPath := filepath.Join(tmpDir, "test.nvpk")
//	testutil.CreateTestPackageFile(t, pkgPath)
//	pkg, err := OpenPackage(ctx, pkgPath)
func CreateTestPackageFile(t *testing.T, path string) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test package file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// Write header
	header := fileformat.NewPackageHeader()
	index := fileformat.NewFileIndex()
	index.EntryCount = 0
	index.FirstEntryOffset = uint64(fileformat.PackageHeaderSize)

	header.IndexStart = uint64(fileformat.PackageHeaderSize)
	header.IndexSize = uint64(16 + int(index.EntryCount)*fileformat.IndexEntrySize)

	if err := binary.Write(file, binary.LittleEndian, header); err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, index.EntryCount); err != nil {
		t.Fatalf("Failed to write index entry count: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, index.Reserved); err != nil {
		t.Fatalf("Failed to write index reserved: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, index.FirstEntryOffset); err != nil {
		t.Fatalf("Failed to write index first entry offset: %v", err)
	}
	for i := range index.Entries {
		if err := binary.Write(file, binary.LittleEndian, index.Entries[i]); err != nil {
			t.Fatalf("Failed to write index entry %d: %v", i, err)
		}
	}
}
