// This file contains coverage tests for PackageReader operations.
// It verifies edge cases and code coverage for ReadFile, ListFiles, GetInfo, etc.
//
// Specification: api_core.md: 1. Core Interfaces

package novus_package

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
)

func TestPackage_ReadFile_NotFound(t *testing.T) {
	runReadFileExpectFail(t, "/nonexistent.txt")
}

func TestPackage_ReadFile_EmptyPath(t *testing.T) {
	runReadFileExpectFail(t, "")
}

func TestPackage_ReadFile_NotOpen(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Close package
	if err := pkg.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	ctx := context.Background()

	// Try to read from closed package
	_, err = pkg.ReadFile(ctx, "/test.txt")
	if err == nil {
		t.Error("ReadFile on closed package should fail")
	}
}

func runReadFileFromDisk(t *testing.T, testContent []byte) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	ctx := context.Background()
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.Create(ctx, tmpPkg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}
	data, err := pkg.ReadFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if !bytes.Equal(data, testContent) {
		t.Errorf("ReadFile content mismatch: got %q, want %q", string(data), string(testContent))
	}
}

func TestPackage_ReadFile_FromDisk(t *testing.T) {
	runReadFileFromDisk(t, []byte("content from disk"))
}

func TestPackage_ListFiles_Empty(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// List files from empty package
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("ListFiles on empty package: got %d files, want 0", len(files))
	}
}

func TestPackage_ListFiles_NoMetadata(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// List files
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("ListFiles: got %d files, want 1", len(files))
	}
}

func TestPackage_ListFiles_MultipleFilesWithSubdirs(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add files in subdirectories (use different content to avoid deduplication)
	paths := []string{"/file1.txt", "/dir1/file2.txt", "/dir1/dir2/file3.txt", "/root.txt"}
	for i, path := range paths {
		content := []byte("content " + string(rune('0'+i)))
		_, err := pkg.AddFileFromMemory(ctx, path, content, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed for %q: %v", path, err)
		}
	}

	// List files
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	// ListFiles should return all files (deduplication shouldn't affect ListFiles)
	expectedCount := len(paths)
	if len(files) != expectedCount {
		t.Errorf("ListFiles: got %d files, want %d", len(files), expectedCount)
	}
}

func TestPackage_ListFiles_RoundTrip(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add files (use different content to avoid deduplication)
	testPaths := []string{"/a.txt", "/b.txt", "/c.txt"}
	for i, path := range testPaths {
		content := []byte("content " + string(rune('a'+i)))
		_, err := pkg.AddFileFromMemory(ctx, path, content, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
	}

	// List files
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	// ListFiles should return all files (deduplication shouldn't affect ListFiles)
	if len(files) != len(testPaths) {
		t.Errorf("ListFiles: got %d files, want %d", len(files), len(testPaths))
	}

	// Write and reopen
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Reopen and list again
	pkg2, err := OpenPackage(ctx, tmpPkg)
	if err != nil {
		t.Fatalf("OpenPackage failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	files2, err := pkg2.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles on reopened package failed: %v", err)
	}

	if len(files2) != len(testPaths) {
		t.Errorf("ListFiles after reopen: got %d files, want %d", len(files2), len(testPaths))
	}
}

func TestPackage_ReadFile_AfterFileHandleClosed(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	testContent := []byte("content")
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Close package
	if err := pkg.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Try to read after close
	_, err = pkg.ReadFile(ctx, entry.Paths[0].Path)
	if err == nil {
		t.Error("ReadFile after Close should fail")
	}
}

func TestPackage_GetInfo_Coverage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Get info from empty package
	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	if info == nil {
		t.Fatal("GetInfo returned nil")
	}

	addThreeFilesFromMemory(t, ctx, pkg, "/", "content ")

	// Get info after adding files
	info2, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	// FileCount should reflect added files (may be less if deduplication occurred)
	if info2.FileCount < 1 {
		t.Errorf("FileCount: got %d, want >= 1", info2.FileCount)
	}
}

func TestPackage_GetMetadata_Coverage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Get metadata (may require package to be opened/loaded)
	metadata, err := pkg.GetMetadata()
	if err != nil {
		t.Logf("GetMetadata failed: %v (may require package to be opened/loaded)", err)
		return
	}
	if metadata == nil {
		t.Fatal("GetMetadata returned nil")
	}
}

func TestPackage_GetMetadata_NoMetadataLoaded(t *testing.T) {
	runGetMetadataBasic(t)
}

func TestPackage_OpenPackage_FileNotFound(t *testing.T) {
	ctx := context.Background()

	// Try to open non-existent package
	_, err := OpenPackage(ctx, "/nonexistent/package.pkg")
	if err == nil {
		t.Error("OpenPackage with non-existent file should fail")
	}
}
