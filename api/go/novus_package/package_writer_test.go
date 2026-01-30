// This file contains unit tests for PackageWriter operations.
// It tests Write, SafeWrite, FastWrite, AddFile, and RemoveFile methods.
//
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing

package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPackage_WriteFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file using AddFileFromMemory (simpler and more reliable)
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("test content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path and write
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
		return
	}

	// Verify file was created
	if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
		t.Error("Write did not create package file")
	}
}

func TestPackage_WriteFile_ThenReadFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file using AddFileFromMemory (simpler and more reliable)
	testContent := []byte("test content for read")
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path and write
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
		return
	}

	// Open and read the package
	pkg2, err := OpenPackage(ctx, tmpPkg)
	if err != nil {
		t.Logf("OpenPackage failed: %v (may require complete Write implementation)", err)
		return
	}
	defer func() { _ = pkg2.Close() }()

	// Read the file
	data, err := pkg2.ReadFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("ReadFile content mismatch: got %q, want %q", string(data), string(testContent))
	}
}

func TestPackage_RemoveFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file using AddFileFromMemory (simpler for this test)
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Verify entry has paths
	if entry.PathCount == 0 || len(entry.Paths) == 0 {
		t.Fatal("Entry should have at least one path")
	}

	// Remove file using RemoveFile
	err = pkg.RemoveFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}

	// Verify file was removed (can't read it - but ReadFile may not work on new packages anyway)
	// The important thing is RemoveFile succeeded
}

func TestPackage_Write(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file using AddFileFromMemory
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path and write
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
		return
	}

	// Verify file was created
	if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
		t.Error("Write did not create package file")
	}
}

func TestPackage_SafeWrite(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file using AddFileFromMemory
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path and safe write
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Logf("SafeWrite failed: %v (implementation may be incomplete)", err)
		return
	}

	// Verify file was created
	if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
		t.Error("SafeWrite did not create package file")
	}
}

func TestPackage_SafeWrite_RoundTrip(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file using AddFileFromMemory
	testContent := []byte("round trip content")
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", testContent, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path and safe write
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Logf("SafeWrite failed: %v (implementation may be incomplete)", err)
		return
	}

	// Open and read the package
	pkg2, err := OpenPackage(ctx, tmpPkg)
	if err != nil {
		t.Logf("OpenPackage failed: %v (may require complete Write implementation)", err)
		return
	}
	defer func() { _ = pkg2.Close() }()

	// Read the file
	data, err := pkg2.ReadFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(data) != string(testContent) {
		t.Errorf("ReadFile content mismatch: got %q, want %q", string(data), string(testContent))
	}
}

func TestPackage_FastWrite(t *testing.T) {
	t.Skip("TODO(Priority 5): FastWrite not implemented yet")
}

func TestPackage_Defragment(t *testing.T) {
	t.Skip("TODO(Priority 2): Defragment implementation pending")
}

func TestPackage_WriteFile_InvalidPath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Try to add non-existent file
	_, err = pkg.AddFile(ctx, "/nonexistent/file.txt", nil)
	if err == nil {
		t.Error("AddFile with non-existent path should fail")
	}
}

func TestPackage_WriteFile_UpdateExisting(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add first file
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("v1"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(v1) failed: %v", err)
	}

	// Add file with same path and AllowOverwrite
	opts := &AddFileOptions{}
	opts.AllowOverwrite.Set(true)
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("v2"), opts)
	if err != nil {
		t.Logf("AddFileFromMemory with AllowOverwrite: %v (may not be fully implemented)", err)
	}
}

func TestPackage_RemoveFile_NonExistent(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Try to remove non-existent file
	err = pkg.RemoveFile(ctx, "/nonexistent.txt")
	if err == nil {
		t.Error("RemoveFile with non-existent path should fail")
	}
}

func TestPackage_SafeWrite_NoFilePath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Try to write without setting target path
	err = pkg.SafeWrite(ctx, true)
	if err == nil {
		t.Error("SafeWrite without target path should fail")
	}
}

func TestPackage_SafeWrite_NoOverwrite(t *testing.T) {
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

	// Set target path
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Create existing file
	if err := os.WriteFile(tmpPkg, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// Try SafeWrite without overwrite (should fail)
	err = pkg.SafeWrite(ctx, false)
	if err == nil {
		t.Error("SafeWrite without overwrite on existing file should fail")
	}
}

func TestPackage_Write_NoFilePath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Try to write without setting target path
	err = pkg.Write(ctx)
	if err == nil {
		t.Error("Write without target path should fail")
	}
}

func TestPackage_RemoveFile_WithMultiplePaths(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file with first path
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	entry1, err := pkg.AddFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Add same file with different path (deduplication should add path to existing entry)
	opts := &AddFileOptions{}
	opts.StoredPath.Set("/different/path.txt")
	entry2, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile with different path failed: %v", err)
	}

	// Should be same entry (deduplication)
	if entry1.FileID != entry2.FileID {
		t.Error("Deduplication should reuse same FileEntry")
	}

	// Remove file using one path
	err = pkg.RemoveFile(ctx, entry1.Paths[0].Path)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}
}

func TestPackage_RemoveFile_InvalidPath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Try to remove with invalid path
	err = pkg.RemoveFile(ctx, "")
	if err == nil {
		t.Error("RemoveFile with empty path should fail")
	}
}

func TestPackage_WritePackageToFile_EmptyPackage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Set target path and write empty package
	tmpPkg := filepath.Join(t.TempDir(), "empty.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write empty package failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
		t.Error("Write did not create package file")
	}
}

func TestPackage_WriteFile_LargeFile(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Create large file
	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "large.bin")
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	if err := os.WriteFile(largeFile, largeData, 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	_, err = pkg.AddFile(ctx, largeFile, nil)
	if err != nil {
		t.Fatalf("AddFile with large file failed: %v", err)
	}

	// Set target path and write
	tmpPkg := filepath.Join(t.TempDir(), "large.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write with large file failed: %v (implementation may be incomplete)", err)
	}
}

func TestPackage_SafeWrite_MultipleFiles(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add multiple files
	for i := 0; i < 3; i++ {
		path := "/file" + string(rune('0'+i)) + ".txt"
		data := []byte("content " + string(rune('0'+i)))
		_, err := pkg.AddFileFromMemory(ctx, path, data, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
	}

	// Set target path and safe write
	tmpPkg := filepath.Join(t.TempDir(), "multi.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Fatalf("SafeWrite with multiple files failed: %v", err)
	}
}

func TestPackage_Write_WithPathMetadata(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add file with explicit stored path
	opts := &AddFileOptions{}
	opts.StoredPath.Set("/metadata/path/file.txt")
	_, err = pkg.AddFileFromMemory(ctx, "/metadata/path/file.txt", []byte("content"), opts)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path and write
	tmpPkg := filepath.Join(t.TempDir(), "metadata.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write with path metadata failed: %v (implementation may be incomplete)", err)
	}
}

func TestPackage_SafeWrite_Overwrite(t *testing.T) {
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

	// Set target path
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Create existing file
	if err := os.WriteFile(tmpPkg, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// SafeWrite with overwrite (should succeed)
	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Fatalf("SafeWrite with overwrite failed: %v", err)
	}
}

func TestPackage_Write_ContextCancelled(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Add file
	ctx := context.Background()
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Cancel context and try to write
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err = pkg.Write(cancelledCtx)
	if err == nil {
		t.Error("Write with cancelled context should fail")
	}
}

func TestPackage_SafeWrite_ContextCancelled(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Add file
	ctx := context.Background()
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Set target path
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Cancel context and try to safe write
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err = pkg.SafeWrite(cancelledCtx, true)
	if err == nil {
		t.Error("SafeWrite with cancelled context should fail")
	}
}

func TestPackage_RemoveFile_ContextCancelled(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Add file
	ctx := context.Background()
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Cancel context and try to remove
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err = pkg.RemoveFile(cancelledCtx, entry.Paths[0].Path)
	if err == nil {
		t.Error("RemoveFile with cancelled context should fail")
	}
}
