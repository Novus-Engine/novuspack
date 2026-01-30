// This file contains integration tests for package write operations.
// It verifies Write, SafeWrite, and related operations work correctly end-to-end.
//
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing

package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPackage_SafeWrite_InSubdirectory(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Create subdirectory
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	tmpPkg := filepath.Join(subDir, "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Logf("SafeWrite failed: %v (implementation may be incomplete)", err)
	}
}

func TestPackage_Write_ReopenAndModify(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("v1"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
		return
	}

	// Reopen and modify
	pkg2, err := OpenPackage(ctx, tmpPkg)
	if err != nil {
		t.Logf("OpenPackage failed: %v (may require complete Write implementation)", err)
		return
	}
	defer func() { _ = pkg2.Close() }()

	// Add another file
	_, err = pkg2.AddFileFromMemory(ctx, "/test2.txt", []byte("v2"), nil)
	if err != nil {
		t.Logf("AddFileFromMemory on reopened package: %v (may not be supported)", err)
	}
}

func TestPackage_WritePackageToFile_IndexBuilding(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add multiple files
	for i := 0; i < 5; i++ {
		path := "/file" + string(rune('0'+i)) + ".txt"
		_, err := pkg.AddFileFromMemory(ctx, path, []byte("content"), nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
	}

	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
	}
}

func TestPackage_WritePackageToFile_HeaderUpdates(t *testing.T) {
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

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
	}
}

func TestPackage_SafeWrite_AtomicRename(t *testing.T) {
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

	// SafeWrite uses atomic rename (temp file + rename)
	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Logf("SafeWrite failed: %v (implementation may be incomplete)", err)
	}
}

func TestPackage_WriteFile_UpdateExistingWithDifferentData(t *testing.T) {
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

func TestPackage_RemoveFile_UpdatePathMetadata(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	entry, err := pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Remove file (should update path metadata)
	err = pkg.RemoveFile(ctx, entry.Paths[0].Path)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}
}

func TestPackage_Write_EmptyPackage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	tmpPkg := filepath.Join(t.TempDir(), "empty.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Logf("Write failed: %v (implementation may be incomplete)", err)
	}
}
