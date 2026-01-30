// This file contains integration tests for path canonicalization with
// ReadFile, AddFile, and RemoveFile operations.
// It verifies that the path canonicalization logic correctly handles dot segments
// in real package operations.
//
// Specification: api_core.md: 1.1.2 Package Path Semantics

package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestPathCanonicalization_AddFile_Integration tests that AddFile correctly
// canonicalizes paths with dot segments.
func TestPathCanonicalization_AddFile_Integration(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx := context.Background()

	// Test path with dot segments
	opts := &AddFileOptions{}
	opts.StoredPath.Set("./dot/../segments/./file.txt")

	entry, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Verify path was canonicalized (should remove dot segments)
	if entry.PathCount == 0 || len(entry.Paths) == 0 {
		t.Fatal("FileEntry has no paths")
	}

	// Path should be canonicalized (no dot segments)
	canonicalPath := entry.Paths[0].Path
	if canonicalPath != "/segments/file.txt" {
		t.Errorf("Path not canonicalized correctly: got %q, want %q", canonicalPath, "/segments/file.txt")
	}
}

func TestPathCanonicalization_ReadFile_Integration(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Create package to open it (required for ReadFile)
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.Create(ctx, tmpPkg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	testContent := []byte("test content")

	// Add file with path containing dot segments using AddFileFromMemory
	opts := &AddFileOptions{}
	opts.StoredPath.Set("./a/../b/./file.txt")
	entry, err := pkg.AddFileFromMemory(ctx, "./a/../b/./file.txt", testContent, opts)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Read file using canonicalized path
	data, err := pkg.ReadFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("ReadFile content mismatch: got %q, want %q", string(data), string(testContent))
	}
}

func TestPathCanonicalization_RemoveFile_Integration(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Create package to open it (required for ReadFile verification)
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.Create(ctx, tmpPkg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Add file using AddFileFromMemory (path will be normalized to /path/to/file.txt)
	entry, err := pkg.AddFileFromMemory(ctx, "/path/to/file.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Verify the stored path
	if entry.PathCount == 0 || len(entry.Paths) == 0 {
		t.Fatal("Entry should have at least one path")
	}
	storedPath := entry.Paths[0].Path
	if storedPath != "/path/to/file.txt" {
		t.Fatalf("Unexpected stored path: got %q, want %q", storedPath, "/path/to/file.txt")
	}

	// Remove file using the canonical path
	err = pkg.RemoveFile(ctx, storedPath)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}

	// Verify file was removed (can't read it)
	_, err = pkg.ReadFile(ctx, storedPath)
	if err == nil {
		t.Error("ReadFile should fail after RemoveFile")
	}
}

func TestPathCanonicalization_RoundTrip(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Create package to open it (required for ReadFile)
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.Create(ctx, tmpPkg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	testContent := []byte("round trip content")

	// Add file with dot segments using AddFileFromMemory
	opts := &AddFileOptions{}
	opts.StoredPath.Set("./a/../b/./c/file.txt")
	entry, err := pkg.AddFileFromMemory(ctx, "./a/../b/./c/file.txt", testContent, opts)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Verify canonicalized path
	canonicalPath := entry.Paths[0].Path
	if canonicalPath != "/b/c/file.txt" {
		t.Errorf("Path not canonicalized: got %q, want %q", canonicalPath, "/b/c/file.txt")
	}

	// Read file using canonical path
	data, err := pkg.ReadFile(ctx, canonicalPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("Content mismatch: got %q, want %q", string(data), string(testContent))
	}

	// Remove file using dot segments (should canonicalize to same path)
	err = pkg.RemoveFile(ctx, "./b/../b/./c/file.txt")
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}
}
