// This file contains edge case tests for PackageWriter operations.
// It tests error paths, edge conditions, and uncommon scenarios.
//
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing

package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPackage_WritePackageToFile_StreamFromDisk(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file from disk (streaming)
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content from disk"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = pkg.AddFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}

func TestPackage_WritePackageToFile_ContextCancelDuringWrite(t *testing.T) {
	runWriteContextCancelled(t)
}

func TestPackage_WritePackageToFile_NilFileEntry(t *testing.T) {
	// This test would require creating a package with nil FileEntry
	// This is an internal error condition that shouldn't occur with AddFile
	t.Skip("Internal error condition - AddFile prevents nil FileEntry")
}

func TestPackage_WritePackageToFile_NoDataSource(t *testing.T) {
	// This test would require a FileEntry without SourceFile or Data
	// AddFileFromMemory always sets Data, AddFile always sets SourceFile
	t.Skip("AddFile/AddFileFromMemory always provide data source")
}

func TestPackage_SafeWrite_TempFileCloseError(t *testing.T) {
	// This test would require simulating a temp file close error
	// Difficult to test without mocking file operations
	t.Skip("Requires file operation mocking")
}

func TestPackage_SafeWrite_RenameError(t *testing.T) {
	// This test would require simulating a rename error
	// Difficult to test without mocking file operations
	t.Skip("Requires file operation mocking")
}

func TestPackage_WritePackageToFile_LargeFileStreaming(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Create large file on disk
	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "large.bin")
	largeData := make([]byte, 5*1024*1024) // 5MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	if err := os.WriteFile(largeFile, largeData, 0o644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	_, err = pkg.AddFile(ctx, largeFile, nil)
	if err != nil {
		t.Fatalf("AddFile with large file failed: %v", err)
	}

	tmpPkg := filepath.Join(t.TempDir(), "large.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}

func TestPackage_WriteFile_WithMultiplePaths(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file with first path
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	entry1, err := pkg.AddFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Add same file with different path (deduplication)
	opts := &AddFileOptions{}
	opts.StoredPath.Set("/path2.txt")
	entry2, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile with different path failed: %v", err)
	}

	// Should be same entry with multiple paths
	if entry1.FileID != entry2.FileID {
		t.Error("Deduplication should reuse same FileEntry")
	}

	if entry1.PathCount < 2 {
		t.Errorf("PathCount: got %d, want >= 2", entry1.PathCount)
	}
}

func TestPackage_ReadFile_StreamingPath(t *testing.T) {
	runReadFileFromDisk(t, []byte("streaming content"))
}

func TestPackage_ReadFile_CompressedFile(t *testing.T) {
	t.Skip("TODO(Priority 5): Compression not implemented yet")
}

func TestPackage_ReadFile_EncryptedFile(t *testing.T) {
	t.Skip("TODO(Priority 6): Encryption not implemented yet")
}

func TestPackage_Write_MultipleWriteCycles(t *testing.T) {
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

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// First write
	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("First Write failed: %v", err)
	}

	// Second write (overwrite)
	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Second Write failed: %v", err)
	}
}
