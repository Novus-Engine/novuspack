// Package testutil provides testing utilities for the fileformat package.
//
// This package contains helper functions for creating test fixtures and
// mock file format structures. It should only be imported by test files.
package testutil

import (
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
	header.IndexSize = uint64(index.Size())

	if _, err := header.WriteTo(file); err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}
	if _, err := index.WriteTo(file); err != nil {
		t.Fatalf("Failed to write index: %v", err)
	}
}
