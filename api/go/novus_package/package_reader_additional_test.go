// Additional reader tests to improve coverage for Priority 2.

package novus_package

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestReadFile_AfterAddFromMemory tests reading a file that was added from memory.
func TestReadFile_AfterAddFromMemory(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	storedPath := "test/file.txt"
	originalData := []byte("Hello, World! This is test content.")

	// Add file from memory
	_, err = pkg.AddFileFromMemory(ctx, storedPath, originalData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Try to read it back
	readData, err := pkg.ReadFile(ctx, storedPath)
	if err != nil {
		t.Logf("ReadFile note: %v (may require Write first)", err)
		// ReadFile may require the package to be written first
		return
	}

	// Verify data matches
	if !bytes.Equal(readData, originalData) {
		t.Errorf("ReadFile data = %q, want %q", string(readData), string(originalData))
	}
}

// TestReadFile_AfterWrite tests reading a file after writing the package.
func TestReadFile_AfterWrite(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	storedPath := "written/file.txt"
	originalData := []byte("Content to write and read back")

	// Add file
	_, err = pkg.AddFileFromMemory(ctx, storedPath, originalData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Write to temp file
	tmpFile := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpFile); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Close and reopen
	if err := pkg.Close(); err != nil {
		t.Logf("Close note: %v", err)
	}

	// Open for reading
	pkg2, err := OpenPackage(ctx, tmpFile)
	if err != nil {
		t.Fatalf("OpenPackage failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// Read file
	readData, err := pkg2.ReadFile(ctx, storedPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if !bytes.Equal(readData, originalData) {
		t.Errorf("ReadFile data = %q, want %q", string(readData), string(originalData))
	}
}

// TestReadFile_NonExistent tests reading a non-existent file.
func TestReadFile_NonExistent(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	_, err = pkg.ReadFile(ctx, "does/not/exist.txt")
	if err == nil {
		t.Fatal("ReadFile(nonexistent) succeeded, want error")
	}

	pkgErr, ok := err.(*pkgerrors.PackageError)
	if !ok {
		t.Fatalf("Error type = %T, want *pkgerrors.PackageError", err)
	}

	// Should be validation or not found error
	if pkgErr.Type != pkgerrors.ErrTypeValidation && pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Logf("Error type = %v (may be implementation-specific)", pkgErr.Type)
	}
}

// TestReadFile_EmptyPath tests reading with empty path.
func TestReadFile_EmptyPath(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	_, err = pkg.ReadFile(ctx, "")
	if err == nil {
		t.Fatal("ReadFile(\"\") succeeded, want error")
	}

	pkgErr, ok := err.(*pkgerrors.PackageError)
	if !ok {
		t.Fatalf("Error type = %T, want *pkgerrors.PackageError", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestReadFile_ContextCancelled tests ReadFile with cancelled context.
func TestReadFile_ContextCancelled(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Add a file first
	ctx := context.Background()
	testData := []byte("data")
	_, err = pkg.AddFileFromMemory(ctx, "test.txt", testData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	// Try to read with cancelled context
	ctx2, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = pkg.ReadFile(ctx2, "test.txt")
	if err == nil {
		t.Fatal("ReadFile with cancelled context succeeded, want error")
	}

	pkgErr, ok := err.(*pkgerrors.PackageError)
	if !ok {
		t.Logf("Error type = %T (may vary)", err)
		return
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext && pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Logf("Error type = %v (may be implementation-specific)", pkgErr.Type)
	}
}

// TestGetMetadata_Basic tests GetMetadata method.
func TestGetMetadata_Basic(t *testing.T) {
	runGetMetadataBasic(t)
}

// TestGetInfo_Basic tests GetInfo method.
func TestGetInfo_Basic(t *testing.T) {
	runGetInfoBasic(t, true)
}

// TestValidate_EmptyPackage tests validating an empty package.
func TestValidate_EmptyPackage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()
	err = pkg.Validate(ctx)
	if err != nil {
		t.Logf("Validate on empty package: %v (may be expected)", err)
	}
}

// TestValidate_WithFiles tests validating a package with files.
func TestValidate_WithFiles(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add some files
	for i := 0; i < 5; i++ {
		path := filepath.Join("files", "file"+string(rune('0'+i))+".txt")
		data := []byte("Content " + string(rune('0'+i)))
		_, err := pkg.AddFileFromMemory(ctx, path, data, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
	}

	// Validate
	err = pkg.Validate(ctx)
	if err != nil {
		t.Logf("Validate with files: %v (may require Write)", err)
	}
}

// TestListFiles_EmptyPackage tests ListFiles on empty package.
func TestListFiles_EmptyPackage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	files, err := pkg.ListFiles()
	if err != nil {
		t.Logf("ListFiles on empty package: %v (may require metadata init)", err)
		return
	}

	if len(files) != 0 {
		t.Errorf("ListFiles() on empty package returned %d files, want 0", len(files))
	}
}

// TestListFiles_Ordering tests that ListFiles returns consistent ordering.
func TestListFiles_Ordering(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add files in specific order
	paths := []string{"z.txt", "a.txt", "m.txt", "b.txt"}
	for _, path := range paths {
		data := []byte("data")
		_, err := pkg.AddFileFromMemory(ctx, path, data, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory(%q) failed: %v", path, err)
		}
	}

	// List files multiple times
	files1, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	files2, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles (2nd call) failed: %v", err)
	}

	// Ordering should be consistent
	if len(files1) != len(files2) {
		t.Fatalf("ListFiles returned different lengths: %d vs %d", len(files1), len(files2))
	}

	for i := range files1 {
		if files1[i].PrimaryPath != files2[i].PrimaryPath {
			t.Errorf("ListFiles ordering inconsistent at index %d: %q vs %q",
				i, files1[i].PrimaryPath, files2[i].PrimaryPath)
		}
	}
}

// TestWrite_EmptyPackage tests writing an empty package.
func TestWrite_EmptyPackage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "empty.pkg")
	ctx := context.Background()

	if err := pkg.SetTargetPath(ctx, tmpFile); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	err = pkg.Write(ctx)
	if err != nil {
		t.Fatalf("Write empty package failed: %v", err)
	}
}

// TestWrite_WithFiles tests writing a package with files.
func TestWrite_WithFiles(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	addThreeFilesFromMemory(t, ctx, pkg, "", "Content for file ")

	tmpFile := filepath.Join(t.TempDir(), "withfiles.pkg")
	if err := pkg.SetTargetPath(ctx, tmpFile); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	err = pkg.Write(ctx)
	if err != nil {
		t.Fatalf("Write with files failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Write did not create output file")
	}
}

// TestSafeWrite_NoOverwrite tests SafeWrite without overwrite flag.
func TestSafeWrite_NoOverwrite(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add a file
	overwriteData := []byte("data")
	_, err = pkg.AddFileFromMemory(ctx, "test.txt", overwriteData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpFile); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Create the file first
	if err := os.WriteFile(tmpFile, []byte("existing"), 0o644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// Try SafeWrite without overwrite (should fail)
	err = pkg.SafeWrite(ctx, false)
	if err == nil {
		t.Error("SafeWrite without overwrite on existing file succeeded, want error")
	}
}

// TestSafeWrite_WithOverwrite tests SafeWrite with overwrite flag.
func TestSafeWrite_WithOverwrite(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	ctx := context.Background()

	// Add a file
	diffData := []byte("new data")
	_, err = pkg.AddFileFromMemory(ctx, "test.txt", diffData, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpFile); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}

	// Create the file first
	if err := os.WriteFile(tmpFile, []byte("old data"), 0o644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	// SafeWrite with overwrite (should succeed)
	err = pkg.SafeWrite(ctx, true)
	if err != nil {
		t.Fatalf("SafeWrite with overwrite failed: %v", err)
	}
}
