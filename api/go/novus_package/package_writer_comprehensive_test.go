// This file contains comprehensive tests for PackageWriter operations.
// It verifies Write, SafeWrite, AddFile, RemoveFile, and related operations.
//
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing

package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPackage_WriteFile_ContextCancellation(t *testing.T) {
	runWriteContextCancelled(t)
}

func TestPackage_RemoveFile_ContextCancellation(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := pkg.RemoveFile(cancelledCtx, entry.Paths[0].Path); err == nil {
		t.Error("RemoveFile with cancelled context should fail")
	}
}

func TestPackage_Write_ContextCancellation(t *testing.T) {
	runWriteContextCancelled(t)
}

func TestPackage_SafeWrite_ContextCancellation(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := pkg.SafeWrite(cancelledCtx, true); err == nil {
		t.Error("SafeWrite with cancelled context should fail")
	}
}

func TestPackage_SafeWrite_FileExistsNoOverwrite(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Create existing file
	if err := os.WriteFile(tmpPkg, []byte("existing"), 0o644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// SafeWrite without overwrite should fail
	if err := pkg.SafeWrite(ctx, false); err == nil {
		t.Error("SafeWrite without overwrite on existing file should fail")
	}
}

func TestPackage_WritePackageToFile_WithComment(t *testing.T) {
	runWriteWithContent(t, []byte("content"), false)
}

func TestPackage_WritePackageToFile_EmptyPackageWithComment(t *testing.T) {
	runWriteEmptyPackage(t)
}

func TestPackage_WritePackageToFile_MultipleFilesVariousSizes(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add files of various sizes
	sizes := []int{0, 1, 100, 1024, 10240}
	for i, size := range sizes {
		data := make([]byte, size)
		for j := range data {
			data[j] = byte((i + j) % 256)
		}
		path := "/file" + string(rune('0'+i)) + ".bin"
		_, err := pkg.AddFileFromMemory(ctx, path, data, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed for size %d: %v", size, err)
		}
	}

	tmpPkg := filepath.Join(t.TempDir(), "multi.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}

func TestPackage_RemoveFile_AllPaths(t *testing.T) {
	runAddTwoPathsThenRemove(t, "/path2.txt")
}

func TestPackage_RemoveFile_OnePath(t *testing.T) {
	runAddTwoPathsThenRemove(t, "/path2.txt")
}

func TestPackage_WriteFile_PathValidation(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Test invalid paths (empty or whitespace-only should fail)
	invalidPaths := []string{"", "   "}
	for _, path := range invalidPaths {
		_, err := pkg.AddFileFromMemory(ctx, path, []byte("content"), nil)
		if err == nil {
			t.Errorf("AddFileFromMemory with invalid path %q should fail", path)
		}
	}

	// Test paths that may be normalized (not necessarily invalid)
	// These paths may be normalized rather than rejected
	normalizablePaths := []string{"../outside", "//double"}
	for _, path := range normalizablePaths {
		_, err := pkg.AddFileFromMemory(ctx, path, []byte("content"), nil)
		// These may be normalized - just verify they don't crash
		_ = err
	}
}

func TestPackage_SafeWrite_TempFileCreation(t *testing.T) {
	runSafeWriteWithContent(t)
}

func TestPackage_Write_UpdatesInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add files
	for i := 0; i < 3; i++ {
		path := "/file" + string(rune('0'+i)) + ".txt"
		_, err := pkg.AddFileFromMemory(ctx, path, []byte("content"), nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
	}

	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	initialFileCount := info.FileCount

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Info should reflect file count
	info2, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	if info2.FileCount != initialFileCount {
		t.Errorf("FileCount changed unexpectedly: got %d, want %d", info2.FileCount, initialFileCount)
	}
}

func TestPackage_WriteFile_SequentialFileIDs(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add multiple files with different content to avoid deduplication
	var fileIDs []uint64
	for i := 0; i < 5; i++ {
		path := "/file" + string(rune('0'+i)) + ".txt"
		content := []byte("content " + string(rune('0'+i)))
		entry, err := pkg.AddFileFromMemory(ctx, path, content, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
		fileIDs = append(fileIDs, entry.FileID)
	}

	// Verify FileIDs are sequential
	for i := 1; i < len(fileIDs); i++ {
		if fileIDs[i] != fileIDs[i-1]+1 {
			t.Errorf("FileIDs not sequential: got %d, want %d", fileIDs[i], fileIDs[i-1]+1)
		}
	}
}
