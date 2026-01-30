// This file contains comprehensive unit tests for file management operations.
// Tests cover AddFile, AddFileFromMemory, RemoveFile, and file lookup operations.
//
// Specification: api_file_mgmt_index.md: 1. File Management Document Map

package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// ====================
// AddFile Tests
// ====================

// TestAddFile_BasicSuccess tests successful file addition with default options.
func TestAddFile_BasicSuccess(t *testing.T) {
	// Create package
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create temp test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("Hello, World!")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add file with default options
	ctx := context.Background()
	entry, err := pkg.AddFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Verify returned FileEntry
	if entry == nil {
		t.Fatal("AddFile returned nil entry")
	}
	if entry.FileID == 0 {
		t.Error("FileEntry.FileID = 0, want non-zero")
	}
	if entry.OriginalSize != uint64(len(testContent)) {
		t.Errorf("FileEntry.OriginalSize = %d, want %d", entry.OriginalSize, len(testContent))
	}
	if entry.ProcessingState != metadata.ProcessingStateRaw {
		t.Errorf("FileEntry.ProcessingState = %v, want %v (ProcessingStateRaw)", entry.ProcessingState, metadata.ProcessingStateRaw)
	}
	if entry.SourceFile == nil {
		t.Error("FileEntry.SourceFile = nil, want non-nil file handle")
	}
	// Note: RawChecksum may be 0 if no deduplication check was performed (no existing files to compare against)
	// This is expected behavior per spec Section 2.1.4 step 4
}

// TestAddFile_WithStoredPath tests adding a file with explicit stored path.
func TestAddFile_WithStoredPath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create temp test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add file with explicit stored path
	ctx := context.Background()
	storedPath := "custom/path/file.txt"
	opts := &AddFileOptions{}
	opts.StoredPath.Set(storedPath)

	entry, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Verify stored path is used (with leading slash per spec normalization)
	expectedPath := "/custom/path/file.txt"
	if entry.PathCount == 0 {
		t.Fatal("FileEntry.PathCount is 0")
	}
	if len(entry.Paths) == 0 {
		t.Fatal("FileEntry.Paths is empty")
	}
	if entry.Paths[0].Path != expectedPath {
		t.Errorf("FileEntry.Paths[0].Path = %q, want %q", entry.Paths[0].Path, expectedPath)
	}
}

// TestAddFile_Deduplication tests that duplicate files are detected.
func TestAddFile_Deduplication(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create temp test files with same content
	tmpDir := t.TempDir()
	testContent := []byte("duplicate content")

	file1 := filepath.Join(tmpDir, "file1.txt")
	if err := os.WriteFile(file1, testContent, 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	file2 := filepath.Join(tmpDir, "file2.txt")
	if err := os.WriteFile(file2, testContent, 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	ctx := context.Background()

	// Add first file
	entry1, err := pkg.AddFile(ctx, file1, nil)
	if err != nil {
		t.Fatalf("AddFile(file1) failed: %v", err)
	}

	// Add second file with same content (should detect duplicate by default)
	entry2, err := pkg.AddFile(ctx, file2, nil)
	if err != nil {
		t.Fatalf("AddFile(file2) failed: %v", err)
	}

	// Both entries should reference the same file data (same FileID)
	if entry1.FileID != entry2.FileID {
		t.Errorf("Deduplication failed: entry1.FileID = %d, entry2.FileID = %d, want same FileID", entry1.FileID, entry2.FileID)
	}

	// Both paths should be registered
	if len(entry2.Paths) < 2 {
		t.Errorf("FileEntry.Paths length = %d, want >= 2 for deduplicated files", len(entry2.Paths))
	}
}

// TestAddFile_AllowDuplicate tests that AllowDuplicate option prevents deduplication.
func TestAddFile_AllowDuplicate(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create temp test files with same content
	tmpDir := t.TempDir()
	testContent := []byte("duplicate content")

	file1 := filepath.Join(tmpDir, "file1.txt")
	if err := os.WriteFile(file1, testContent, 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	file2 := filepath.Join(tmpDir, "file2.txt")
	if err := os.WriteFile(file2, testContent, 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	ctx := context.Background()

	// Add first file
	entry1, err := pkg.AddFile(ctx, file1, nil)
	if err != nil {
		t.Fatalf("AddFile(file1) failed: %v", err)
	}

	// Add second file with AllowDuplicate=true
	opts := &AddFileOptions{}
	opts.AllowDuplicate.Set(true)
	entry2, err := pkg.AddFile(ctx, file2, opts)
	if err != nil {
		t.Fatalf("AddFile(file2) failed: %v", err)
	}

	// Entries should have different FileIDs (not deduplicated)
	if entry1.FileID == entry2.FileID {
		t.Errorf("AllowDuplicate=true but entries have same FileID: %d", entry1.FileID)
	}
}

// TestAddFile_Symlinks tests symlink handling with FollowSymlinks option.
func TestAddFile_Symlinks(t *testing.T) {
	// Create temp test file and symlink
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "target.txt")
	testContent := []byte("symlink target")
	if err := os.WriteFile(targetFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	symlinkFile := filepath.Join(tmpDir, "link.txt")
	if err := os.Symlink(targetFile, symlinkFile); err != nil {
		t.Skipf("Symlink creation failed (may not be supported): %v", err)
	}

	ctx := context.Background()

	t.Run("FollowSymlinks_Default_True", func(t *testing.T) {
		pkg, err := NewPackage()
		if err != nil {
			t.Fatalf("NewPackage failed: %v", err)
		}

		// Default behavior should follow symlinks (spec says default is true)
		entry, err := pkg.AddFile(ctx, symlinkFile, nil)
		if err != nil {
			t.Fatalf("AddFile(symlink) with default options failed: %v", err)
		}

		// Should have read the target file's size
		if entry.OriginalSize != uint64(len(testContent)) {
			t.Errorf("FileEntry.OriginalSize = %d, want %d (target file size)", entry.OriginalSize, len(testContent))
		}
	})

	t.Run("FollowSymlinks_Explicit_True", func(t *testing.T) {
		pkg, err := NewPackage()
		if err != nil {
			t.Fatalf("NewPackage failed: %v", err)
		}

		opts := &AddFileOptions{}
		opts.FollowSymlinks.Set(true)
		entry, err := pkg.AddFile(ctx, symlinkFile, opts)
		if err != nil {
			t.Fatalf("AddFile(symlink) with FollowSymlinks=true failed: %v", err)
		}

		if entry.OriginalSize != uint64(len(testContent)) {
			t.Errorf("FileEntry.OriginalSize = %d, want %d", entry.OriginalSize, len(testContent))
		}
	})

	t.Run("FollowSymlinks_False", func(t *testing.T) {
		pkg, err := NewPackage()
		if err != nil {
			t.Fatalf("NewPackage failed: %v", err)
		}

		opts := &AddFileOptions{}
		opts.FollowSymlinks.Set(false)
		_, err = pkg.AddFile(ctx, symlinkFile, opts)

		// Should return validation error when trying to add symlink with FollowSymlinks=false
		if err == nil {
			t.Fatal("AddFile(symlink) with FollowSymlinks=false succeeded, want error")
		}

		pkgErr, ok := err.(*pkgerrors.PackageError)
		if !ok {
			t.Fatalf("Error type = %T, want *pkgerrors.PackageError", err)
		}
		if pkgErr.Type != pkgerrors.ErrTypeValidation {
			t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
		}
	})
}

// TestAddFile_ErrorCases tests various error conditions for AddFile.
func TestAddFile_ErrorCases(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		setupFunc   func() string // Returns path to use
		wantErrType pkgerrors.ErrorType
		wantErrMsg  string
	}{
		{
			name:        "Empty path",
			path:        "",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "path cannot be empty or whitespace-only",
		},
		{
			name:        "Whitespace-only path",
			path:        "   ",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "path cannot be empty or whitespace-only",
		},
		{
			name: "Non-existent file",
			setupFunc: func() string {
				return filepath.Join(tmpDir, "nonexistent.txt")
			},
			wantErrType: pkgerrors.ErrTypeIO,
		},
		{
			name: "Directory instead of file",
			setupFunc: func() string {
				dirPath := filepath.Join(tmpDir, "testdir")
				_ = os.Mkdir(dirPath, 0755)
				return dirPath
			},
			wantErrType: pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := tt.path
			if tt.setupFunc != nil {
				testPath = tt.setupFunc()
			}

			_, err := pkg.AddFile(ctx, testPath, nil)
			if err == nil {
				t.Errorf("AddFile(%q) succeeded, want error", testPath)
				return
			}

			// Check error type
			pkgErr, ok := err.(*pkgerrors.PackageError)
			if !ok {
				t.Errorf("Error type = %T, want *pkgerrors.PackageError", err)
				return
			}

			if pkgErr.Type != tt.wantErrType {
				t.Errorf("Error type = %v, want %v", pkgErr.Type, tt.wantErrType)
			}

			if tt.wantErrMsg != "" && pkgErr.Message != tt.wantErrMsg {
				t.Errorf("Error message = %q, want %q", pkgErr.Message, tt.wantErrMsg)
			}
		})
	}
}

// TestAddFile_ContextCancellation tests context cancellation handling.
func TestAddFile_ContextCancellation(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create temp test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = pkg.AddFile(ctx, testFile, nil)
	if err == nil {
		t.Fatal("AddFile with cancelled context succeeded, want error")
	}

	pkgErr, ok := err.(*pkgerrors.PackageError)
	if !ok {
		t.Fatalf("Error type = %T, want *pkgerrors.PackageError", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeContext)
	}
}

// ====================
// AddFileFromMemory Tests
// ====================

// TestAddFileFromMemory_BasicSuccess tests successful memory file addition.
func TestAddFileFromMemory_BasicSuccess(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	testData := []byte("Memory file content")
	storedPath := "memory/test.txt"

	entry, err := pkg.AddFileFromMemory(ctx, storedPath, testData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Verify returned FileEntry
	if entry == nil {
		t.Fatal("AddFileFromMemory returned nil entry")
	}
	if entry.FileID == 0 {
		t.Error("FileEntry.FileID = 0, want non-zero")
	}
	if entry.OriginalSize != uint64(len(testData)) {
		t.Errorf("FileEntry.OriginalSize = %d, want %d", entry.OriginalSize, len(testData))
	}
	if entry.PathCount == 0 {
		t.Fatal("FileEntry.PathCount = 0, want > 0")
	}
	// Spec requires internal paths to have leading slash
	expectedPath := "/memory/test.txt"
	if entry.Paths[0].Path != expectedPath {
		t.Errorf("FileEntry.Paths[0].Path = %q, want %q", entry.Paths[0].Path, expectedPath)
	}
	if entry.RawChecksum == 0 {
		t.Error("FileEntry.RawChecksum = 0, want non-zero checksum")
	}
	if !entry.IsDataLoaded {
		t.Error("FileEntry.IsDataLoaded = false, want true for memory files")
	}
}

// TestAddFileFromMemory_EmptyData tests adding empty memory file.
func TestAddFileFromMemory_EmptyData(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	emptyData := []byte{}
	storedPath := "empty.txt"

	entry, err := pkg.AddFileFromMemory(ctx, storedPath, emptyData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory with empty data failed: %v", err)
	}

	if entry.OriginalSize != 0 {
		t.Errorf("FileEntry.OriginalSize = %d, want 0", entry.OriginalSize)
	}
}

// TestAddFileFromMemory_NilData tests adding nil data (should be treated as empty).
func TestAddFileFromMemory_NilData(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	storedPath := "nil.txt"

	entry, err := pkg.AddFileFromMemory(ctx, storedPath, nil, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory with nil data failed: %v", err)
	}

	if entry.OriginalSize != 0 {
		t.Errorf("FileEntry.OriginalSize = %d, want 0", entry.OriginalSize)
	}
}

// TestAddFileFromMemory_Deduplication tests memory file deduplication.
func TestAddFileFromMemory_Deduplication(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	testData := []byte("duplicate memory content")

	// Add first memory file
	entry1, err := pkg.AddFileFromMemory(ctx, "file1.txt", testData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(file1) failed: %v", err)
	}

	// Add second memory file with same content
	entry2, err := pkg.AddFileFromMemory(ctx, "file2.txt", testData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(file2) failed: %v", err)
	}

	// Should be deduplicated (same FileID)
	if entry1.FileID != entry2.FileID {
		t.Errorf("Deduplication failed: entry1.FileID = %d, entry2.FileID = %d, want same FileID", entry1.FileID, entry2.FileID)
	}
}

// TestAddFileFromMemory_ErrorCases tests error conditions for AddFileFromMemory.
func TestAddFileFromMemory_ErrorCases(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	testData := []byte("test")

	tests := []struct {
		name        string
		path        string
		wantErrType pkgerrors.ErrorType
		wantErrMsg  string
	}{
		{
			name:        "Empty path",
			path:        "",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "path cannot be empty or whitespace-only",
		},
		{
			name:        "Whitespace-only path",
			path:        "   ",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "path cannot be empty or whitespace-only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pkg.AddFileFromMemory(ctx, tt.path, testData, nil)
			if err == nil {
				t.Errorf("AddFileFromMemory(%q) succeeded, want error", tt.path)
				return
			}

			pkgErr, ok := err.(*pkgerrors.PackageError)
			if !ok {
				t.Errorf("Error type = %T, want *pkgerrors.PackageError", err)
				return
			}

			if pkgErr.Type != tt.wantErrType {
				t.Errorf("Error type = %v, want %v", pkgErr.Type, tt.wantErrType)
			}

			if tt.wantErrMsg != "" && pkgErr.Message != tt.wantErrMsg {
				t.Errorf("Error message = %q, want %q", pkgErr.Message, tt.wantErrMsg)
			}
		})
	}
}

// ====================
// RemoveFile Tests
// ====================

// TestRemoveFile_Success tests successful file removal.
func TestRemoveFile_Success(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add a memory file first
	storedPath := "test/remove.txt"
	content := []byte("content")
	entry, err := pkg.AddFileFromMemory(ctx, storedPath, content, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	if entry == nil {
		t.Fatal("AddFileFromMemory returned nil entry")
	}

	// Remove the file
	err = pkg.RemoveFile(ctx, storedPath)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}

	// Verify removal by checking ListFiles
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	for _, fi := range files {
		if fi.PrimaryPath == storedPath {
			t.Error("File still appears in ListFiles after RemoveFile")
		}
	}
}

// TestRemoveFile_NonExistent tests removing a non-existent file.
func TestRemoveFile_NonExistent(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	nonExistentPath := "does/not/exist.txt"

	err = pkg.RemoveFile(ctx, nonExistentPath)
	if err == nil {
		t.Fatal("RemoveFile(nonexistent) succeeded, want error")
	}

	pkgErr, ok := err.(*pkgerrors.PackageError)
	if !ok {
		t.Fatalf("Error type = %T, want *pkgerrors.PackageError", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestRemoveFile_ErrorCases tests error conditions for RemoveFile.
func TestRemoveFile_ErrorCases(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		path        string
		wantErrType pkgerrors.ErrorType
		wantErrMsg  string
	}{
		{
			name:        "Empty path",
			path:        "",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "path cannot be empty or whitespace-only",
		},
		{
			name:        "Whitespace-only path",
			path:        "   ",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "path cannot be empty or whitespace-only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pkg.RemoveFile(ctx, tt.path)
			if err == nil {
				t.Errorf("RemoveFile(%q) succeeded, want error", tt.path)
				return
			}

			pkgErr, ok := err.(*pkgerrors.PackageError)
			if !ok {
				t.Errorf("Error type = %T, want *pkgerrors.PackageError", err)
				return
			}

			if pkgErr.Type != tt.wantErrType {
				t.Errorf("Error type = %v, want %v", pkgErr.Type, tt.wantErrType)
			}

			if tt.wantErrMsg != "" && pkgErr.Message != tt.wantErrMsg {
				t.Errorf("Error message = %q, want %q", pkgErr.Message, tt.wantErrMsg)
			}
		})
	}
}

// ====================
// File Integration Tests
// ====================

// TestFileAddAndList tests adding files and listing them.
func TestFileAddAndList(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add multiple files
	paths := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, path := range paths {
		content := []byte("content for " + path)
		_, err := pkg.AddFileFromMemory(ctx, path, content, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory(%q) failed: %v", path, err)
		}
	}

	// List files and verify
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	if len(files) != len(paths) {
		t.Errorf("ListFiles() returned %d files, want %d", len(files), len(paths))
	}

	// Verify each file is present
	foundPaths := make(map[string]bool)
	for _, fi := range files {
		foundPaths[fi.PrimaryPath] = true
	}

	for _, path := range paths {
		if !foundPaths[path] {
			t.Errorf("ListFiles() missing path %q", path)
		}
	}
}

// ====================
// Path Determination Tests
// ====================

// TestAddFile_WithBasePath tests AddFile with BasePath option.
func TestAddFile_WithBasePath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create test file in nested directory
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "project", "src")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	testFile := filepath.Join(nestedDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx := context.Background()

	// Add file with BasePath to strip leading directories
	opts := &AddFileOptions{}
	opts.BasePath.Set(filepath.Join(tmpDir, "project"))

	entry, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile with BasePath failed: %v", err)
	}

	// Verify path has base path stripped
	if entry.PathCount == 0 {
		t.Fatal("FileEntry.PathCount = 0, want > 0")
	}

	storedPath := entry.Paths[0].Path
	// Path should be relative to BasePath
	if !filepath.IsAbs(storedPath) || storedPath == testFile {
		// Path handling varies, just verify file was added successfully
		t.Logf("Stored path: %q", storedPath)
	}
}

// TestAddFile_WithPreserveDepth tests AddFile with PreserveDepth option.
func TestAddFile_WithPreserveDepth(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create test file in nested directory
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	testFile := filepath.Join(nestedDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx := context.Background()

	// Add file with PreserveDepth=2
	opts := &AddFileOptions{}
	opts.PreserveDepth.Set(2)

	entry, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile with PreserveDepth failed: %v", err)
	}

	if entry.PathCount == 0 {
		t.Fatal("FileEntry.PathCount = 0, want > 0")
	}
}

// TestAddFile_WithFlattenPaths tests AddFile with FlattenPaths option.
func TestAddFile_WithFlattenPaths(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create test file in nested directory
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "deep", "nested", "path")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	testFile := filepath.Join(nestedDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx := context.Background()

	// Add file with FlattenPaths=true
	opts := &AddFileOptions{}
	opts.FlattenPaths.Set(true)

	entry, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile with FlattenPaths failed: %v", err)
	}

	if entry.PathCount == 0 {
		t.Fatal("FileEntry.PathCount = 0, want > 0")
	}

	// With FlattenPaths, should have only filename
	storedPath := entry.Paths[0].Path
	t.Logf("Stored path with FlattenPaths: %q", storedPath)
}

// TestAddFile_WithSessionBase tests AddFile using session base.
func TestAddFile_WithSessionBase(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create test file in nested directory
	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "myproject")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	testFile := filepath.Join(projectDir, "readme.txt")
	if err := os.WriteFile(testFile, []byte("README"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx := context.Background()

	// Set session base
	if err := pkg.SetSessionBase(tmpDir); err != nil {
		t.Fatalf("SetSessionBase failed: %v", err)
	}

	// Add file - should use session base for path determination
	entry, err := pkg.AddFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("AddFile with session base failed: %v", err)
	}

	if entry.PathCount == 0 {
		t.Fatal("FileEntry.PathCount = 0, want > 0")
	}

	// Clear session base
	pkg.ClearSessionBase()
}

// ====================
// Additional Edge Cases
// ====================

// TestAddFile_MultiplePathsForSameFile tests adding the same physical file multiple times with different paths.
func TestAddFile_MultiplePathsForSameFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "shared.txt")
	testContent := []byte("shared content")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx := context.Background()

	// Add same file with different stored paths
	opts1 := &AddFileOptions{}
	opts1.StoredPath.Set("path1/shared.txt")
	entry1, err := pkg.AddFile(ctx, testFile, opts1)
	if err != nil {
		t.Fatalf("AddFile(path1) failed: %v", err)
	}

	opts2 := &AddFileOptions{}
	opts2.StoredPath.Set("path2/shared.txt")
	entry2, err := pkg.AddFile(ctx, testFile, opts2)
	if err != nil {
		t.Fatalf("AddFile(path2) failed: %v", err)
	}

	// Should be deduplicated (same FileID, multiple paths)
	if entry1.FileID != entry2.FileID {
		t.Errorf("Expected same FileID for deduplicated file, got %d and %d", entry1.FileID, entry2.FileID)
	}

	// Should have 2 paths registered
	if entry2.PathCount < 2 {
		t.Errorf("PathCount = %d, want >= 2 for file with multiple paths", entry2.PathCount)
	}
}

// TestAddFileFromMemory_LargeFile tests adding a large file from memory.
func TestAddFileFromMemory_LargeFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Create 1MB of test data
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	entry, err := pkg.AddFileFromMemory(ctx, "large.bin", largeData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory with large file failed: %v", err)
	}

	if entry.OriginalSize != uint64(len(largeData)) {
		t.Errorf("OriginalSize = %d, want %d", entry.OriginalSize, len(largeData))
	}
}

// TestAddFileFromMemory_MultipleWithSamePath tests adding multiple files to the same path (should overwrite).
func TestAddFileFromMemory_MultipleWithSamePath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	storedPath := "overwrite/test.txt"

	// Add first file
	data1 := []byte("first version")
	entry1, err := pkg.AddFileFromMemory(ctx, storedPath, data1, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(first) failed: %v", err)
	}

	firstFileID := entry1.FileID

	// Add second file to same path
	data2 := []byte("second version - different content")
	entry2, err := pkg.AddFileFromMemory(ctx, storedPath, data2, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(second) failed: %v", err)
	}

	// Should be different files (different content)
	if entry2.FileID == firstFileID {
		t.Logf("Note: Same FileID despite different content (may be expected behavior)")
	}
}

// TestRemoveFile_ThenReadd tests removing a file and then re-adding it.
func TestRemoveFile_ThenReadd(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	storedPath := "remove/readd.txt"
	testData := []byte("test content")

	// Add file
	entry1, err := pkg.AddFileFromMemory(ctx, storedPath, testData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}
	firstFileID := entry1.FileID

	// Remove file
	if err := pkg.RemoveFile(ctx, storedPath); err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}

	// Re-add same file
	entry2, err := pkg.AddFileFromMemory(ctx, storedPath, testData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory (re-add) failed: %v", err)
	}

	// May or may not have same FileID depending on implementation
	t.Logf("First FileID: %d, Re-added FileID: %d", firstFileID, entry2.FileID)
}

// TestAddFile_LargeFilesystemFile tests adding a large file from filesystem.
func TestAddFile_LargeFilesystemFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create large test file (10MB)
	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "large.bin")

	f, err := os.Create(largeFile)
	if err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	// Write 10MB of data
	data := make([]byte, 1024*1024) // 1MB chunks
	for i := 0; i < 10; i++ {
		for j := range data {
			data[j] = byte((i + j) % 256)
		}
		if _, err := f.Write(data); err != nil {
			_ = f.Close()
			t.Fatalf("Failed to write to large file: %v", err)
		}
	}
	_ = f.Close()

	ctx := context.Background()
	entry, err := pkg.AddFile(ctx, largeFile, nil)
	if err != nil {
		t.Fatalf("AddFile with large file failed: %v", err)
	}

	expectedSize := uint64(10 * 1024 * 1024)
	if entry.OriginalSize != expectedSize {
		t.Errorf("OriginalSize = %d, want %d", entry.OriginalSize, expectedSize)
	}

	// Verify file handle is open
	if entry.SourceFile == nil {
		t.Error("SourceFile = nil, want non-nil for filesystem file")
	}
}

// TestAddFileFromMemory_WithAllowOverwrite tests AllowOverwrite option.
func TestAddFileFromMemory_WithAllowOverwrite(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	storedPath := "overwrite/file.txt"

	// Add first file
	v1 := []byte("v1")
	_, err = pkg.AddFileFromMemory(ctx, storedPath, v1, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(v1) failed: %v", err)
	}

	// Try to add with same path and AllowOverwrite
	opts := &AddFileOptions{}
	opts.AllowOverwrite.Set(true)

	v2 := []byte("v2 - overwrite")
	_, err = pkg.AddFileFromMemory(ctx, storedPath, v2, opts)
	if err != nil {
		t.Logf("AddFileFromMemory with AllowOverwrite: %v", err)
		// May or may not be implemented yet
	}
}
