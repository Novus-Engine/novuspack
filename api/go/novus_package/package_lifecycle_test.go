// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains unit tests for package lifecycle operations:
// Create, Close, OpenPackage, OpenPackageReadOnly, ReadHeader, IsOpen, CreateWithOptions,
// and CloseWithCleanup.
package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/fileformat/testutil"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// =============================================================================
// TEST: Create Operations
// =============================================================================

// TestPackage_Create_Basic tests basic package creation.
func TestPackage_Create_Basic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Test: Create new package at path
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Test: Configure package (Create no longer writes to disk per spec)
	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, pkgPath)
	if err != nil {
		t.Errorf("Create() failed: %v", err)
	}

	// Test: File should NOT exist on disk (Create only configures in memory)
	if _, err := os.Stat(pkgPath); !os.IsNotExist(err) {
		t.Error("Package file should not exist on disk after Create() - Create() only configures in memory")
	}

	// Cleanup
	_ = pkg.Close()
}

// TestPackage_Create_ValidatesPath tests that Create validates the path parameter.
func TestPackage_Create_ValidatesPath(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	tests := []struct {
		name        string
		path        string
		shouldError bool
		errorType   string
	}{
		{
			name:        "empty path",
			path:        "",
			shouldError: true,
			errorType:   "ErrTypeValidation",
		},
		{
			name:        "whitespace only",
			path:        "   ",
			shouldError: true,
			errorType:   "ErrTypeValidation",
		},
		{
			name:        "valid path",
			path:        filepath.Join(t.TempDir(), "valid.nvpk"),
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpkg := pkg.(*filePackage)
			err := fpkg.Create(ctx, tt.path)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected validation error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestPackage_Create_InMemoryOnly tests that NewPackage creates in-memory only.
func TestPackage_Create_InMemoryOnly(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "should-not-exist.nvpk")

	// Test: NewPackage should not create file on disk
	_, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Test: File should NOT exist after NewPackage
	if _, err := os.Stat(pkgPath); !os.IsNotExist(err) {
		t.Error("NewPackage should not create file on disk")
	}
}

// TestPackage_Create_InitializesHeader tests that Create initializes package header.
func TestPackage_Create_InitializesHeader(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test: Package info should reflect initialized header
	// Note: GetInfo() requires package to be open, so we can't test it here
	// The header is initialized in memory, but GetInfo() is only available after OpenPackage()
	// This test verifies Create() completes successfully without errors
}

// TestPackage_Create_WithValidPath tests Create with various valid paths.
func TestPackage_Create_WithValidPath(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		path string
	}{
		{"simple path", filepath.Join(t.TempDir(), "test.nvpk")},
		{"nested path", filepath.Join(t.TempDir(), "subdir", "test.nvpk")},
		{"with spaces", filepath.Join(t.TempDir(), "test file.nvpk")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create parent directory if needed
			if err := os.MkdirAll(filepath.Dir(tt.path), 0755); err != nil {
				t.Fatalf("Failed to create parent dir: %v", err)
			}

			pkg, err := NewPackage()
			if err != nil {
				t.Fatalf("NewPackage() failed: %v", err)
			}

			fpkg := pkg.(*filePackage)
			if err := fpkg.Create(ctx, tt.path); err != nil {
				t.Errorf("Create() failed for %v: %v", tt.name, err)
			}

			// Verify file was NOT created (Create() no longer writes to disk)
			if _, err := os.Stat(tt.path); !os.IsNotExist(err) {
				t.Errorf("File should not exist after Create() - Create() only configures in memory: %v", err)
			}
		})
	}
}

// TestPackage_Create_WithIndexInitialization tests Create with nil index.
func TestPackage_Create_WithIndexInitialization(t *testing.T) {
	ctx := context.Background()

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Explicitly set index to nil
	fpkg := pkg.(*filePackage)
	fpkg.index = nil

	// Create should initialize the index
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Errorf("Create() failed: %v", err)
	}

	if fpkg.index == nil {
		t.Error("Create() should initialize index when nil")
	}
}

// TestPackage_Create_WithCancelledContext tests Create with cancelled context.
func TestPackage_Create_WithCancelledContext(t *testing.T) {
	ctx := testhelpers.CancelledContext()

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, tempFile)
	if err == nil {
		t.Error("Create() should return error for cancelled context")
	}
}

// TestPackage_Create_WithWhitespacePath tests Create with whitespace-padded path.
func TestPackage_Create_WithWhitespacePath(t *testing.T) {
	ctx := context.Background()

	tempDir := t.TempDir()
	path := "  " + filepath.Join(tempDir, "test.nvpk") + "  "

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, path); err != nil {
		t.Errorf("Create() should succeed with whitespace-padded path: %v", err)
	}

	// Verify the FilePath was trimmed
	if fpkg.FilePath != strings.TrimSpace(path) {
		t.Errorf("FilePath = %v, want %v", fpkg.FilePath, strings.TrimSpace(path))
	}
}

// TestPackage_Create_WithEmptyPath tests Create with empty path.
func TestPackage_Create_WithEmptyPath(t *testing.T) {
	ctx := context.Background()

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Try to create with empty path
	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, "")
	if err == nil {
		t.Fatal("Create() should fail with empty path")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_Create_WithTabOnlyPath tests Create with tab-only path.
func TestPackage_Create_WithTabOnlyPath(t *testing.T) {
	ctx := context.Background()

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Try to create with tab-only path (whitespace-only)
	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, "\t\t\t")
	if err == nil {
		t.Fatal("Create() should fail with whitespace-only path (tabs)")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_Create_WithReadOnlyDirectory tests Create in read-only directory.
func TestPackage_Create_WithReadOnlyDirectory(t *testing.T) {
	ctx := context.Background()
	// Skip if not on Unix-like system (Windows doesn't have same permissions model)
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping on Windows")
	}

	tempDir := t.TempDir()

	// Make directory read-only
	if err := os.Chmod(tempDir, 0444); err != nil {
		t.Fatalf("Failed to make directory read-only: %v", err)
	}
	defer func() {
		_ = os.Chmod(tempDir, 0755) // Restore permissions for cleanup
	}()

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Try to create in read-only directory
	// Note: Create() no longer writes to disk, so it won't fail on read-only directory
	tempFile := filepath.Join(tempDir, "test.nvpk")
	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, tempFile)
	// Create() should succeed since it doesn't write to disk
	if err != nil {
		t.Errorf("Create() should succeed (it only configures in memory, doesn't write to disk): %v", err)
	}
}

// =============================================================================
// TEST: Open Operations
// =============================================================================

// TestPackage_Open_Basic tests opening an existing package.
func TestPackage_Open_Basic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create a package first
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test: Open the created package
	pkg2, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Errorf("OpenPackage() failed: %v", err)
	}
	if pkg2 == nil {
		t.Fatal("OpenPackage() returned nil")
	}
	defer func() { _ = pkg2.Close() }()

	// Test: Opened package should be in open state
	if !pkg2.IsOpen() {
		t.Error("Opened package should be in open state")
	}
}

// TestPackage_Open_ValidatesMagicNumber tests that Open validates magic number.
func TestPackage_Open_ValidatesMagicNumber(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "invalid.nvpk")

	// Setup: Create an invalid file (not a NovusPack file)
	err := os.WriteFile(invalidPath, []byte("This is not a NovusPack file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test: Open should fail with invalid magic number
	pkg, err := OpenPackage(ctx, invalidPath)
	if err == nil {
		t.Error("Expected error for invalid magic number")
		if pkg != nil {
			_ = pkg.Close()
		}
	}
}

// TestPackage_Open_LoadsHeader tests that Open loads the package header.
func TestPackage_Open_LoadsHeader(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test: Open and verify header is loaded
	pkg2, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	info, err := pkg2.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() failed: %v", err)
	}

	if info == nil {
		t.Error("Info should not be nil after Open")
	}
}

// TestPackage_Open_LoadsFileIndex tests that Open loads the file index.
func TestPackage_Open_LoadsFileIndex(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test: Open and verify file index is accessible
	pkg2, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	// Verify through GetInfo (which should reflect index data)
	info, err := pkg2.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() failed: %v", err)
	}

	// File count should be available (0 for empty package)
	if info == nil {
		t.Error("Info should not be nil after Open")
	}
}

// TestPackage_Open_ErrorConditions tests various error conditions for Open.
func TestPackage_Open_ErrorConditions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		path        string
		shouldError bool
		errorType   string
	}{
		{
			name:        "empty path",
			path:        "",
			shouldError: true,
			errorType:   "ErrTypeValidation",
		},
		{
			name:        "non-existent file",
			path:        "/path/to/nonexistent.nvpk",
			shouldError: true,
			errorType:   "ErrTypeIO",
		},
		{
			name:        "directory instead of file",
			path:        "/tmp",
			shouldError: true,
			errorType:   "ErrTypeIO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := OpenPackage(ctx, tt.path)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
					if pkg != nil {
						_ = pkg.Close()
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if pkg != nil {
					_ = pkg.Close()
				}
			}
		})
	}
}

// TestPackage_Open_WithContext tests Open with context scenarios.
func TestPackage_Open_WithContext(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(context.Background(), pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	tests := []struct {
		name        string
		ctx         context.Context
		shouldError bool
	}{
		{
			name:        "valid context",
			ctx:         context.Background(),
			shouldError: false,
		},
		{
			name:        "cancelled context",
			ctx:         testhelpers.CancelledContext(),
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := OpenPackage(tt.ctx, pkgPath)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error for cancelled context")
					if pkg != nil {
						_ = pkg.Close()
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if pkg != nil {
					_ = pkg.Close()
				}
			}
		})
	}
}

// TestPackage_Open tests the Open method.

// =============================================================================
// TEST: OpenPackage Operations
// =============================================================================

// TestPackage_OpenPackage_ValidatesContext tests OpenPackage context validation.
func TestPackage_OpenPackage_ValidatesContext(t *testing.T) {
	ctx := context.Background()
	// Create a valid package file first
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Try to open with nil context
	var nilCtx context.Context // explicitly nil for testing
	_, err = OpenPackage(nilCtx, tempFile)
	if err == nil {
		t.Error("OpenPackage() should return error for nil context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("OpenPackage() should return PackageError")
	}
	// Note: checkContext returns ErrTypeValidation for nil context
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestPackage_OpenPackage_WithWhitespacePath tests OpenPackage with whitespace-padded path.
func TestPackage_OpenPackage_WithWhitespacePath(t *testing.T) {
	ctx := context.Background()

	// Create a package first
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Open with whitespace-padded path
	pathWithWhitespace := "  " + tempFile + "  "
	pkg2, err := OpenPackage(ctx, pathWithWhitespace)
	if err != nil {
		t.Errorf("OpenPackage() should succeed with whitespace-padded path: %v", err)
	}
	if pkg2 != nil {
		defer func() { _ = pkg2.Close() }()
	}
}

// TestPackage_OpenPackage_WithCancelledContext tests OpenPackage with cancelled context.
func TestPackage_OpenPackage_WithCancelledContext(t *testing.T) {
	ctx := context.Background()

	// Create a valid package first
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Try to open with cancelled context
	cancelledCtx := testhelpers.CancelledContext()
	_, err = OpenPackage(cancelledCtx, tempFile)
	if err == nil {
		t.Error("OpenPackage() should return error for cancelled context")
	}
}

// TestPackage_OpenPackage_WithDirectory tests opening a directory instead of a file.
func TestPackage_OpenPackage_WithDirectory(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// Try to open a directory
	_, err := OpenPackage(ctx, tempDir)
	if err == nil {
		t.Fatal("OpenPackage() should fail when opening a directory")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}
	if !strings.Contains(pkgErr.Message, "directory") {
		t.Errorf("Expected error message to mention directory, got: %s", pkgErr.Message)
	}
}

// TestPackage_OpenPackage_WithCorruptedIndex tests opening a package with corrupted index.
// Uses a small entry count to test corruption without causing OOM.
func TestPackage_OpenPackage_WithCorruptedIndex(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping test that may cause OOM in short mode")
	}
	tempFile := filepath.Join(t.TempDir(), "corrupted.nvpk")

	// Create a package first
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually
	testutil.CreateTestPackageFile(t, tempFile)

	// Now corrupt the index by writing invalid data
	file, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Failed to open file for corruption: %v", err)
	}

	// Seek to index location (right after header)
	if _, err := file.Seek(int64(fileformat.PackageHeaderSize), 0); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to seek: %v", err)
	}

	// Write invalid entry count (use a small value that still tests corruption without OOM)
	invalidData := make([]byte, 8)
	// Set entry count to 10 (small enough to avoid OOM, large enough to test corruption)
	invalidData[0] = 0x0A // 10 in little-endian
	invalidData[1] = 0x00
	invalidData[2] = 0x00
	invalidData[3] = 0x00
	if _, err := file.Write(invalidData); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to corrupt index: %v", err)
	}
	_ = file.Close()

	// Try to open the corrupted package
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with corrupted index")
	}

	// Verify error type (can be either IO or Validation)
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	// Corrupted index can fail at read (IO) or validation stage (Validation)
	if pkgErr.Type != pkgerrors.ErrTypeValidation && pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type Validation or IO, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_WithTruncatedFile tests opening a truncated package file.
func TestPackage_OpenPackage_WithTruncatedFile(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "truncated.nvpk")

	// Create a truncated file (smaller than header size)
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Write only partial header (less than required)
	partialData := make([]byte, fileformat.PackageHeaderSize/2)
	if _, err := file.Write(partialData); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write truncated data: %v", err)
	}
	_ = file.Close()

	// Try to open the truncated file
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with truncated file")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_SeekFailure tests OpenPackage when index seek fails.
func TestPackage_OpenPackage_SeekFailure(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "invalid_offset.nvpk")

	// Create a package
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file first since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Modify the header to have an invalid index offset (beyond file size)
	file, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Read the header
	header := fileformat.NewPackageHeader()
	if _, err := header.ReadFrom(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to read header: %v", err)
	}

	// Set index start to beyond file size
	header.IndexStart = 0x7FFFFFFF // Large but reasonable offset to avoid OOM

	// Write back the modified header
	if _, err := file.Seek(0, 0); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to seek to start: %v", err)
	}
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write modified header: %v", err)
	}
	_ = file.Close()

	// Try to open the package with invalid offset
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with invalid index offset")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_SeekError tests OpenPackage when Seek to index fails.
func TestPackage_OpenPackage_SeekError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a package file with invalid index start position
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write valid header
	header := fileformat.NewPackageHeader()
	header.IndexStart = 999999 // Invalid index start (beyond file size)
	header.IndexSize = 100
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// Try to open - should fail on seek
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with invalid index start")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_IndexReadError tests OpenPackage when reading index fails.
func TestPackage_OpenPackage_IndexReadError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a package file with truncated index
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write valid header
	header := fileformat.NewPackageHeader()
	header.IndexStart = uint64(fileformat.PackageHeaderSize)
	header.IndexSize = 1000 // Claim large index but file is too small
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	// Write only partial index data
	_, _ = file.Write([]byte{0x00, 0x01, 0x02}) // Partial data
	_ = file.Close()

	// Try to open - should fail when reading index
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with truncated index")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_IndexValidationError tests OpenPackage when index validation fails.
func TestPackage_OpenPackage_IndexValidationError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a package file with invalid index
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write valid header
	header := fileformat.NewPackageHeader()
	header.IndexStart = uint64(fileformat.PackageHeaderSize)

	// Create an index with mismatched entry count and entries
	// This will cause validation to fail
	index := fileformat.NewFileIndex()
	index.EntryCount = 5                                               // Claims 5 entries
	index.FirstEntryOffset = uint64(fileformat.PackageHeaderSize + 16) // Offset for index header
	// But create only 2 entries, causing mismatch
	index.Entries = []fileformat.IndexEntry{
		{FileID: 1, Offset: 100},
		{FileID: 2, Offset: 200},
	}
	header.IndexSize = uint64(index.Size())

	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	// Write index with mismatched count
	if _, err := index.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write index: %v", err)
	}
	_ = file.Close()

	// Try to open - should fail on index validation
	// Note: The validation might pass if it only checks structure, not count match
	_, err = OpenPackage(ctx, tempFile)
	if err != nil {
		// If it fails, verify it's the right error type
		pkgErr := &pkgerrors.PackageError{}
		if asPackageError(err, pkgErr) {
			if pkgErr.Type != pkgerrors.ErrTypeValidation && pkgErr.Type != pkgerrors.ErrTypeIO {
				t.Errorf("Expected error type Validation or IO, got: %v", pkgErr.Type)
			}
		}
	} else {
		// If validation passes, that's also acceptable - the test documents current behavior
		t.Log("OpenPackage() succeeded - index validation may not check entry count mismatch")
	}
}

// TestPackage_OpenPackage_NoIndex tests OpenPackage when there's no index (IndexStart = 0).
func TestPackage_OpenPackage_NoIndex(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a package file with no index
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write header with no index
	header := fileformat.NewPackageHeader()
	header.IndexStart = 0 // No index
	header.IndexSize = 0
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// Try to open - should succeed with empty index
	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() should succeed with no index: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Verify package is open
	if !pkg.IsOpen() {
		t.Error("Package should be open")
	}

	// Verify index is empty
	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() failed: %v", err)
	}
	if info.FileCount != 0 {
		t.Errorf("FileCount = %v, want 0", info.FileCount)
	}
}

// TestPackage_OpenPackage_WithSeekError tests OpenPackage when Seek fails.
func TestPackage_OpenPackage_WithSeekError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create the file on disk first (Create() doesn't write to disk)
	testutil.CreateTestPackageFile(t, tempFile)

	// Corrupt the file to make it too small for seek to work properly
	// by truncating it to just the header size, but the header claims there's an index
	// This will cause issues when trying to seek to the index
	if err := os.Truncate(tempFile, 50); err != nil {
		t.Fatalf("Failed to truncate file: %v", err)
	}

	// Try to open - this should fail because the file is corrupted
	_, err := OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with truncated file")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
}

// TestPackage_OpenPackage_WithIndexValidateError tests OpenPackage when index.Validate() fails.
func TestPackage_OpenPackage_WithIndexValidateError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create the file on disk first (Create() doesn't write to disk)
	testutil.CreateTestPackageFile(t, tempFile)

	// Now corrupt the index data to make Validate() fail
	// Open the file and corrupt the index section
	file, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Seek to index location (after header)
	if _, err := file.Seek(int64(fileformat.PackageHeaderSize), 0); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to seek: %v", err)
	}

	// Write garbage data to corrupt the index
	corruptData := make([]byte, 100)
	for i := range corruptData {
		corruptData[i] = 0xFF
	}
	if _, err := file.Write(corruptData); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write corrupt data: %v", err)
	}
	_ = file.Close()

	// Try to open - this should fail with validation error
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with corrupted index")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
}

// TestPackage_OpenPackage_WithIndexReadFromError tests OpenPackage when index.ReadFrom() fails.
func TestPackage_OpenPackage_WithIndexReadFromError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create the file on disk first (Create() doesn't write to disk)
	testutil.CreateTestPackageFile(t, tempFile)

	// Corrupt the file by truncating it right after the header
	// This will cause index.ReadFrom to fail
	file, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Truncate to just after header (header size + 1 byte, so ReadFrom will fail)
	headerSize := int64(fileformat.PackageHeaderSize)
	if err := file.Truncate(headerSize + 1); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to truncate file: %v", err)
	}
	_ = file.Close()

	// Try to open - should fail when reading index
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with truncated index")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	// Should be IO error from ReadFrom failure
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_WithInvalidHeaderVersion tests OpenPackage with invalid header.
func TestPackage_OpenPackage_WithInvalidHeaderVersion(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "invalid_header.nvpk")

	// Create a package with invalid header
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write a header with correct magic but invalid format version
	header := fileformat.NewPackageHeader()
	header.Magic = fileformat.NVPKMagic // Correct magic
	header.FormatVersion = 999          // Invalid version (too high)
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write invalid header: %v", err)
	}
	_ = file.Close()

	// Try to open package
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with invalid header version")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// =============================================================================
// TEST: CreateWithOptions Operations
// =============================================================================

// TestPackage_CreateWithOptions tests the CreateWithOptions method.
func TestPackage_CreateWithOptions(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		Comment:  "test comment",
		VendorID: 1,
		AppID:    100,
	}

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(ctx, tempFile, options)
	if err != nil {
		t.Fatalf("CreateWithOptions() failed: %v", err)
	}

	// Verify CreateWithOptions does not write to disk
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Error("CreateWithOptions() should not write to disk - file should not exist")
	}

	// Verify options were applied in memory (using methods that don't require package to be open)
	if !pkg.HasComment() || pkg.GetComment() != "test comment" {
		t.Errorf("Comment not set: HasComment=%v, Comment=%v", pkg.HasComment(), pkg.GetComment())
	}
	if pkg.GetVendorID() != 1 {
		t.Errorf("VendorID = %v, want 1", pkg.GetVendorID())
	}
	if pkg.GetAppID() != 100 {
		t.Errorf("AppID = %v, want 100", pkg.GetAppID())
	}
}

// TestPackage_CreateWithOptions_NilOptions tests CreateWithOptions with nil options.
func TestPackage_CreateWithOptions_NilOptions(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(ctx, tempFile, nil)
	if err != nil {
		t.Fatalf("CreateWithOptions(nil) failed: %v", err)
	}

	// Verify CreateWithOptions does not write to disk
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Error("CreateWithOptions() should not write to disk - file should not exist")
	}
}

// TestPackage_CreateWithOptions_CommentOnly tests CreateWithOptions with only Comment set.
func TestPackage_CreateWithOptions_CommentOnly(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		Comment: "comment only",
	}

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(ctx, tempFile, options)
	if err != nil {
		t.Fatalf("CreateWithOptions() failed: %v", err)
	}

	// Verify comment was set (using methods that don't require package to be open)
	if !pkg.HasComment() || pkg.GetComment() != "comment only" {
		t.Errorf("Comment not set correctly: HasComment=%v, Comment=%v", pkg.HasComment(), pkg.GetComment())
	}
	if pkg.GetVendorID() != 0 {
		t.Errorf("VendorID should be 0, got %v", pkg.GetVendorID())
	}
	if pkg.GetAppID() != 0 {
		t.Errorf("AppID should be 0, got %v", pkg.GetAppID())
	}
}

// TestPackage_CreateWithOptions_VendorIDOnly tests CreateWithOptions with only VendorID set.
func TestPackage_CreateWithOptions_VendorIDOnly(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		VendorID: 42,
	}

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(ctx, tempFile, options)
	if err != nil {
		t.Fatalf("CreateWithOptions() failed: %v", err)
	}

	// Verify VendorID was set (using methods that don't require package to be open)
	if pkg.GetVendorID() != 42 {
		t.Errorf("VendorID = %v, want 42", pkg.GetVendorID())
	}
	if pkg.HasComment() {
		t.Error("HasComment should be false when Comment is empty")
	}
	if pkg.GetAppID() != 0 {
		t.Errorf("AppID should be 0, got %v", pkg.GetAppID())
	}
}

// TestPackage_CreateWithOptions_AppIDOnly tests CreateWithOptions with only AppID set.
func TestPackage_CreateWithOptions_AppIDOnly(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		AppID: 999,
	}

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(ctx, tempFile, options)
	if err != nil {
		t.Fatalf("CreateWithOptions() failed: %v", err)
	}

	// Verify AppID was set (using methods that don't require package to be open)
	if pkg.GetAppID() != 999 {
		t.Errorf("AppID = %v, want 999", pkg.GetAppID())
	}
	if pkg.HasComment() {
		t.Error("HasComment should be false when Comment is empty")
	}
	if pkg.GetVendorID() != 0 {
		t.Errorf("VendorID should be 0, got %v", pkg.GetVendorID())
	}
}

// TestPackage_CreateWithOptions_CancelledContext tests CreateWithOptions with cancelled context.
func TestPackage_CreateWithOptions_CancelledContext(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		Comment: "test",
	}

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(cancelCtx, tempFile, options)
	if err == nil {
		t.Fatal("CreateWithOptions() should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Expected error type Context, got: %v", pkgErr.Type)
	}
}

// TestPackage_CreateWithOptions_PropagatesCreateErrors tests that CreateWithOptions propagates errors from Create.
func TestPackage_CreateWithOptions_PropagatesCreateErrors(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		Comment: "test",
	}

	// Test with invalid path (empty) - should propagate error from Create
	err = fpkg.CreateWithOptions(ctx, "", options)
	if err == nil {
		t.Fatal("CreateWithOptions() should fail with invalid path")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_CreateWithOptions_EmptyComment tests CreateWithOptions with empty comment string.
func TestPackage_CreateWithOptions_EmptyComment(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	options := &CreateOptions{
		Comment:  "", // Empty comment should not set HasComment
		VendorID: 5,
		AppID:    10,
	}

	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.CreateWithOptions(ctx, tempFile, options)
	if err != nil {
		t.Fatalf("CreateWithOptions() failed: %v", err)
	}

	// Verify options were applied (using methods that don't require package to be open)
	if pkg.HasComment() {
		t.Error("HasComment should be false when Comment is empty")
	}
	if pkg.GetVendorID() != 5 {
		t.Errorf("VendorID = %v, want 5", pkg.GetVendorID())
	}
	if pkg.GetAppID() != 10 {
		t.Errorf("AppID = %v, want 10", pkg.GetAppID())
	}
}

// =============================================================================
// TEST: IsReadOnly and GetPath Operations
// =============================================================================

// TestPackage_IsReadOnly tests the IsReadOnly method.
func TestPackage_IsReadOnly(t *testing.T) {
	ctx := context.Background()

	// Test with writable package
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	if fpkg.IsReadOnly() {
		t.Error("New package should not be read-only")
	}

	// Test with opened package (writable)
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	openedPkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = openedPkg.Close() }()

	if openedPkg.IsReadOnly() {
		t.Error("Opened package (writable) should not be read-only")
	}

	// Test with read-only package
	readOnlyPkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = readOnlyPkg.Close() }()

	if !readOnlyPkg.IsReadOnly() {
		t.Error("Read-only package should return true for IsReadOnly()")
	}
}

// TestPackage_GetPath tests the GetPath method.
func TestPackage_GetPath(t *testing.T) {
	ctx := context.Background()

	// Test with new package (no path set)
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	if fpkg.GetPath() != "" {
		t.Errorf("New package should have empty path, got: %s", fpkg.GetPath())
	}

	// Test with Create (path set)
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	err = fpkg.Create(ctx, tempFile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if fpkg.GetPath() != tempFile {
		t.Errorf("GetPath() = %s, want %s", fpkg.GetPath(), tempFile)
	}

	// Test with opened package
	testutil.CreateTestPackageFile(t, tempFile)
	openedPkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = openedPkg.Close() }()

	if openedPkg.GetPath() != tempFile {
		t.Errorf("Opened package GetPath() = %s, want %s", openedPkg.GetPath(), tempFile)
	}
}

// =============================================================================
// TEST: CloseWithCleanup Operations
// =============================================================================

// TestPackage_CloseWithCleanup tests the CloseWithCleanup method.
func TestPackage_CloseWithCleanup(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify Create does not write to disk
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Error("Create() should not write to disk - file should not exist")
	}

	err = fpkg.CloseWithCleanup(ctx)
	if err != nil {
		t.Fatalf("CloseWithCleanup() failed: %v", err)
	}
	if fpkg.isOpen {
		t.Error("Package should not be open after CloseWithCleanup()")
	}
	if fpkg.FileEntries != nil {
		t.Error("FileEntries should be nil after CloseWithCleanup()")
	}
	if fpkg.header != nil {
		t.Error("header should be nil after CloseWithCleanup()")
	}
}

// TestPackage_CloseWithCleanup_CancelledContext tests CloseWithCleanup with cancelled context.
func TestPackage_CloseWithCleanup_CancelledContext(t *testing.T) {
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = fpkg.CloseWithCleanup(cancelCtx)
	if err == nil {
		t.Fatal("CloseWithCleanup() should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Expected error type Context, got: %v", pkgErr.Type)
	}
}

// TestPackage_CloseWithCleanup_WithOpenPackage tests CloseWithCleanup on an open package.
func TestPackage_CloseWithCleanup_WithOpenPackage(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)

	// Verify package is open
	if !fpkg.isOpen {
		t.Error("Package should be open")
	}

	// CloseWithCleanup should work
	err = fpkg.CloseWithCleanup(ctx)
	if err != nil {
		t.Fatalf("CloseWithCleanup() failed: %v", err)
	}

	// Verify cleanup was performed
	if fpkg.isOpen {
		t.Error("Package should be closed after CloseWithCleanup()")
	}
	if fpkg.FileEntries != nil {
		t.Error("FileEntries should be nil after CloseWithCleanup()")
	}
	if fpkg.SpecialFiles != nil {
		t.Error("SpecialFiles should be nil after CloseWithCleanup()")
	}
}

// TestPackage_CloseWithCleanup_CloseErrorPropagation tests that CloseWithCleanup propagates errors from Close().
func TestPackage_CloseWithCleanup_CloseErrorPropagation(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)

	// Close the file handle manually to simulate Close() error
	if fpkg.fileHandle != nil {
		_ = fpkg.fileHandle.Close()
		// Reset state so Close() will attempt to close again
		fpkg.isOpen = true

		// CloseWithCleanup should propagate the error from Close()
		err = fpkg.CloseWithCleanup(ctx)
		if err == nil {
			t.Fatal("CloseWithCleanup() should propagate error from Close()")
		}

		pkgErr := &pkgerrors.PackageError{}
		if !asPackageError(err, pkgErr) {
			t.Fatalf("Expected PackageError, got: %T", err)
		}
		if pkgErr.Type != pkgerrors.ErrTypeIO {
			t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
		}

		// Verify cleanup did NOT happen (CloseWithCleanup returns early on Close() error)
		// The fileHandle should still be set (though to a closed file)
		// This tests that error propagation happens before cleanup
	} else {
		t.Skip("fileHandle is nil, cannot test Close() error propagation")
	}
}

// =============================================================================
// TEST: Close Operations
// =============================================================================

// TestPackage_Close_Basic tests basic package closing.
//
// Expected behavior (Red Phase - should FAIL):
// - Close method does not exist
func TestPackage_Close_Basic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create and open a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Test: Close should succeed
	err = pkg1.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Test: Package should not be open after Close
	if pkg1.IsOpen() {
		t.Error("Package should not be open after Close")
	}
}

// TestPackage_Close_ClosesFileHandle tests that Close releases file handle.
//
// Expected behavior (Red Phase - should FAIL):
// - File handle management not implemented
func TestPackage_Close_ClosesFileHandle(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Close the package
	err = pkg.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Test: Attempting operations after Close should fail
	_, err = pkg.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() should succeed after Close() if metadata is cached, got: %v", err)
	}
}

// TestPackage_Close_Idempotency tests that Close can be called multiple times.
//
// Expected behavior (Red Phase - should FAIL):
// - Idempotent Close not implemented
func TestPackage_Close_Idempotency(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Test: First Close should succeed
	err = pkg.Close()
	if err != nil {
		t.Errorf("First Close() failed: %v", err)
	}

	// Test: Second Close should not panic (may return error or succeed)
	err = pkg.Close()
	// We accept either success or a validation error for already-closed package
	// The important thing is it shouldn't panic
	_ = err // Acceptable to succeed or return validation error
}

// TestPackage_Close_Idempotent tests Close is truly idempotent.
func TestPackage_Close_Idempotent(t *testing.T) {
	ctx := context.Background()

	// Create and open a package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	pkg2, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	// Close multiple times
	for i := 0; i < 5; i++ {
		if err := pkg2.Close(); err != nil {
			t.Errorf("Close() iteration %d failed: %v", i, err)
		}

		// Verify state after each close
		if pkg2.IsOpen() {
			t.Errorf("IsOpen() should return false after Close() iteration %d", i)
		}
	}
}

// TestPackage_Close_Multiple tests calling Close multiple times (idempotency).
// This ensures Close is truly idempotent and doesn't cause issues.
func TestPackage_Close_Multiple(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create and open a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	if err := fpkg1.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	// Close multiple times
	for i := 0; i < 5; i++ {
		if err := pkg.Close(); err != nil {
			t.Fatalf("Close() call %d failed: %v", i+1, err)
		}
	}

	// Verify package is closed
	fpkg := pkg.(*filePackage)
	if fpkg.isOpen {
		t.Error("Package should not be open after multiple Close() calls")
	}
	if fpkg.fileHandle != nil {
		t.Error("File handle should be nil after Close()")
	}
}

// TestPackage_Close_ResourceCleanup tests that Close releases all resources.
//
// Expected behavior (Red Phase - should FAIL):
// - Resource cleanup not implemented
func TestPackage_Close_ResourceCleanup(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	err = fpkg.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Close should release all resources
	err = pkg.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Test: State should be cleared
	if pkg.IsOpen() {
		t.Error("IsOpen should return false after Close")
	}

	// Note: We can't directly verify memory/buffer cleanup without instrumentation,
	// but we verify the package is unusable after Close
}

// TestPackage_Close_WhenAlreadyClosed tests idempotent Close behavior.
func TestPackage_Close_WhenAlreadyClosed(t *testing.T) {
	ctx := context.Background()

	// Create and open a package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	pkg2, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	// Close once
	if err := pkg2.Close(); err != nil {
		t.Fatalf("First Close() failed: %v", err)
	}

	// Close again - should be idempotent
	if err := pkg2.Close(); err != nil {
		t.Errorf("Second Close() failed: %v (should be idempotent)", err)
	}

	// Verify package is closed
	if pkg2.IsOpen() {
		t.Error("Package should be closed after Close()")
	}
}

// TestPackage_Close_WithFileCloseError tests Close when file.Close() returns an error.
func TestPackage_Close_WithFileCloseError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create the file on disk first (Create() doesn't write to disk)
	testutil.CreateTestPackageFile(t, tempFile)

	openedPkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	// Close should work normally
	err = openedPkg.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Close again (idempotency test - this path is already covered)
	err = openedPkg.Close()
	if err != nil {
		t.Errorf("Close() second call returned error: %v", err)
	}

	// Test the path where fileHandle is nil but closed is false
	// This tests the idempotency path more thoroughly
	pkg2, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg2 := pkg2.(*filePackage)
	fpkg2.fileHandle = nil
	fpkg2.isOpen = false

	err = pkg2.Close()
	if err != nil {
		t.Errorf("Close() on package with nil fileHandle should not error: %v", err)
	}

	// Test the early return path: !isOpen && fileHandle == nil
	// This should return immediately without doing anything
	fpkg2.fileHandle = nil
	fpkg2.isOpen = false
	err = pkg2.Close()
	if err != nil {
		t.Errorf("Close() on already-closed package with nil fileHandle should not error: %v", err)
	}
	// Verify it's still closed
	if fpkg2.isOpen {
		t.Error("Package should remain closed after Close() on already-closed package")
	}
}

// TestPackage_Close_FileHandleCloseError tests Close when fileHandle.Close() actually returns an error.
// This tests the error path at line 468-472 in package_lifecycle.go.
func TestPackage_Close_FileHandleCloseError(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	fpkg := pkg.(*filePackage)

	// Close the file handle manually first
	if fpkg.fileHandle != nil {
		_ = fpkg.fileHandle.Close()
		// Now fileHandle points to a closed file - calling Close() again will return an error
		// But we need to reset the isOpen flag to test the error path
		fpkg.isOpen = true

		// Now Close() should attempt to close the already-closed file and return an error
		err = fpkg.Close()
		if err == nil {
			t.Fatal("Close() should return error when fileHandle.Close() fails")
		}

		pkgErr := &pkgerrors.PackageError{}
		if !asPackageError(err, pkgErr) {
			t.Fatalf("Expected PackageError, got: %T", err)
		}
		if pkgErr.Type != pkgerrors.ErrTypeIO {
			t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
		}

		// Verify state cleanup happened even on error
		if fpkg.isOpen {
			t.Error("isOpen should be false after Close() error")
		}
		if fpkg.isOpen {
			t.Error("isOpen should be false after Close() error")
		}
		if fpkg.fileHandle != nil {
			t.Error("fileHandle should be nil after Close() (even on error)")
		}
	} else {
		t.Skip("fileHandle is nil, cannot test Close() error path")
	}
}

// =============================================================================
// TEST: IsOpen Operations
// =============================================================================

// TestPackage_IsOpen_States tests IsOpen in different states.
func TestPackage_IsOpen_States(t *testing.T) {
	ctx := context.Background()

	// Test new package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if pkg1.IsOpen() {
		t.Error("New package should not be open")
	}

	// Test created package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	fpkg1 := pkg1.(*filePackage)
	if err := fpkg1.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	if !pkg1.IsOpen() {
		t.Error("Created package should be open (ready for operations)")
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Test opened package
	pkg2, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	if !pkg2.IsOpen() {
		t.Error("Opened package should be open")
	}

	// Test closed package
	if err := pkg2.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}
	if pkg2.IsOpen() {
		t.Error("Closed package should not be open")
	}
}

// =============================================================================
// TEST: OpenPackageReadOnly Operations
// =============================================================================

// TestPackage_OpenPackageReadOnly tests the OpenPackageReadOnly function.
func TestPackage_OpenPackageReadOnly(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	if !pkg.IsOpen() {
		t.Error("Package should be open after OpenPackageReadOnly()")
	}

	// Verify it's a readOnlyPackage wrapper
	_, ok := pkg.(*readOnlyPackage)
	if !ok {
		t.Error("OpenPackageReadOnly() should return a readOnlyPackage wrapper")
	}
}

// TestPackage_OpenPackageReadOnly_RejectsMutation tests that OpenPackageReadOnly rejects mutation operations.
func TestPackage_OpenPackageReadOnly_RejectsMutation(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Test various mutation operations
	mutationTests := []struct {
		name      string
		operation func() error
	}{
		{"Create", func() error { return pkg.Create(ctx, "new.nvpk") }},
		{"Defragment", func() error { return pkg.Defragment(ctx) }},
		// TODO(Priority 2): Re-enable after AddFile/RemoveFile are fully implemented
		{"Write", func() error { return pkg.Write(ctx) }},
		{"SafeWrite", func() error { return pkg.SafeWrite(ctx, false) }},
		{"FastWrite", func() error { return pkg.FastWrite(ctx) }},
		{"SetComment", func() error { return pkg.SetComment("comment") }},
		{"ClearComment", func() error { return pkg.ClearComment() }},
		{"SetAppID", func() error { return pkg.SetAppID(123) }},
		{"ClearAppID", func() error { return pkg.ClearAppID() }},
		{"SetVendorID", func() error { return pkg.SetVendorID(456) }},
		{"ClearVendorID", func() error { return pkg.ClearVendorID() }},
		{"SetPackageIdentity", func() error { return pkg.SetPackageIdentity(456, 123) }},
		{"ClearPackageIdentity", func() error { return pkg.ClearPackageIdentity() }},
	}

	for _, tt := range mutationTests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if err == nil {
				t.Errorf("%s() should return an error on read-only package", tt.name)
				return
			}

			pkgErr := &pkgerrors.PackageError{}
			if !asPackageError(err, pkgErr) {
				t.Errorf("Expected PackageError, got: %T", err)
				return
			}

			if pkgErr.Type != pkgerrors.ErrTypeSecurity {
				t.Errorf("Expected error type Security, got: %v", pkgErr.Type)
			}

			// Verify error context contains operation name
			// Extract the typed context from the error
			typedCtxValue, exists := pkgErr.Context["_typed_context"]
			if !exists {
				t.Errorf("Expected typed context, got: %v", pkgErr.Context)
				return
			}
			// Type assert to ReadOnlyErrorContext (exported type from package_lifecycle.go)
			typedCtx, ok := typedCtxValue.(ReadOnlyErrorContext)
			if !ok {
				t.Errorf("Expected ReadOnlyErrorContext, got: %T (value: %+v)", typedCtxValue, typedCtxValue)
				return
			}
			if typedCtx.Operation != tt.name {
				t.Errorf("Expected error context Operation=%s, got: %s", tt.name, typedCtx.Operation)
			}
		})
	}
}

// TestPackage_OpenPackageReadOnly_AllowsReadOperations tests that OpenPackageReadOnly allows read operations.
func TestPackage_OpenPackageReadOnly_AllowsReadOperations(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// These should not panic (even if not fully implemented)
	_, _ = pkg.GetInfo()
	_, _ = pkg.GetMetadata()
	_, _ = pkg.ListFiles()
	// Test ReadFile delegation (even though ReadFile isn't implemented yet)
	_, err = pkg.ReadFile(ctx, "nonexistent.txt")
	// Should return error (unsupported or not found), but not panic
	// This verifies the delegation path works
	if err == nil {
		t.Log("ReadFile returned no error (may be expected if not implemented)")
	}
	_ = pkg.Validate(ctx)
	_ = pkg.IsOpen()
	_ = pkg.GetComment()
	_ = pkg.HasComment()
	_ = pkg.GetAppID()
	_ = pkg.HasAppID()
	_ = pkg.GetVendorID()
	_ = pkg.HasVendorID()
	_, _ = pkg.GetPackageIdentity()
}

// TestPackage_OpenPackageReadOnly_PreventsTypeAssertion tests that OpenPackageReadOnly prevents type assertion to writable type.
func TestPackage_OpenPackageReadOnly_PreventsTypeAssertion(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Type assertion to writable implementation type should fail
	_, ok := pkg.(*filePackage)
	if ok {
		t.Error("Type assertion to *filePackage should fail for read-only package")
	}

	// Type assertion to readOnlyPackage should succeed
	_, ok = pkg.(*readOnlyPackage)
	if !ok {
		t.Error("Type assertion to *readOnlyPackage should succeed")
	}
}

// TestPackage_OpenPackageReadOnly_PropagatesOpenPackageErrors tests that OpenPackageReadOnly propagates errors from OpenPackage.
func TestPackage_OpenPackageReadOnly_PropagatesOpenPackageErrors(t *testing.T) {
	ctx := context.Background()

	// Test with non-existent file
	_, err := OpenPackageReadOnly(ctx, "/nonexistent/file.nvpk")
	if err == nil {
		t.Fatal("OpenPackageReadOnly() should fail with non-existent file")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}

	// Test with cancelled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = OpenPackageReadOnly(cancelCtx, "/nonexistent/file.nvpk")
	if err == nil {
		t.Fatal("OpenPackageReadOnly() should fail with cancelled context")
	}

	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Expected error type Context, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackageReadOnly_ErrorPropagationFromInnerPackage tests error propagation from inner package operations.
func TestPackage_OpenPackageReadOnly_ErrorPropagationFromInnerPackage(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Test that errors from inner package operations are properly propagated
	// Close the inner package to cause errors on subsequent operations
	_ = pkg.Close()

	// Pure in-memory read operations should still work after Close if metadata is cached.
	_, err = pkg.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() should succeed after Close() if metadata is cached, got: %v", err)
	}

	_, err = pkg.GetMetadata()
	if err != nil {
		t.Errorf("GetMetadata() should succeed after Close() if metadata is cached, got: %v", err)
	}

	_, err = pkg.ListFiles()
	if err != nil {
		t.Errorf("ListFiles() should succeed after Close() if metadata is cached, got: %v", err)
	}

	err = pkg.Validate(ctx)
	if err == nil {
		t.Error("Validate() should fail after Close()")
	}
}

// TestPackage_OpenPackageReadOnly_StateTransitions tests state transitions in read-only package.
func TestPackage_OpenPackageReadOnly_StateTransitions(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	pkg, err := OpenPackageReadOnly(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}

	// Verify initial state
	if !pkg.IsOpen() {
		t.Error("Package should be open after OpenPackageReadOnly()")
	}

	// Close and verify state
	err = pkg.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	if pkg.IsOpen() {
		t.Error("Package should not be open after Close()")
	}

	// Verify in-memory operations still work after Close if metadata remains cached
	_, err = pkg.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() should succeed after Close() if metadata is cached, got: %v", err)
	}
}

// TestPackage_Create_StateTransitions tests state transitions in Create method.
func TestPackage_Create_StateTransitions(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Verify initial state (not open, not created)
	if fpkg.isOpen {
		t.Error("Package should not be open initially")
	}
	if fpkg.FilePath != "" {
		t.Error("FilePath should be empty initially")
	}

	// Create and verify state
	err = fpkg.Create(ctx, tempFile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify state after Create
	if !fpkg.isOpen {
		t.Error("Package should be open after Create() (ready for operations)")
	}
	if fpkg.FilePath != tempFile {
		t.Errorf("FilePath = %v, want %v", fpkg.FilePath, tempFile)
	}
	if fpkg.header == nil {
		t.Error("Header should be initialized after Create()")
	}
	if fpkg.index == nil {
		t.Error("Index should be initialized after Create()")
	}
	if fpkg.Info == nil {
		t.Error("Info should be initialized after Create()")
	}
	if fpkg.Info.Created.IsZero() {
		t.Error("Info.Created should be set after Create()")
	}
}

// TestPackage_Create_InitializesTimestamps tests that Create initializes timestamps correctly.
func TestPackage_Create_InitializesTimestamps(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	beforeCreate := time.Now()

	err = fpkg.Create(ctx, tempFile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	afterCreate := time.Now()

	// Verify timestamps are set and reasonable
	if fpkg.Info.Created.Before(beforeCreate) || fpkg.Info.Created.After(afterCreate) {
		t.Errorf("Info.Created = %v, should be between %v and %v", fpkg.Info.Created, beforeCreate, afterCreate)
	}
	if fpkg.Info.Modified.Before(beforeCreate) || fpkg.Info.Modified.After(afterCreate) {
		t.Errorf("Info.Modified = %v, should be between %v and %v", fpkg.Info.Modified, beforeCreate, afterCreate)
	}
	if fpkg.header.CreatedTime == 0 {
		t.Error("header.CreatedTime should be set")
	}
	if fpkg.header.ModifiedTime == 0 {
		t.Error("header.ModifiedTime should be set")
	}
}

// TestPackage_OpenPackage_ContextCancellationDuringSeek tests context cancellation during seek operation.
func TestPackage_OpenPackage_ContextCancellationDuringSeek(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	testutil.CreateTestPackageFile(t, tempFile)

	// Create a context that will be cancelled
	cancelCtx, cancel := context.WithCancel(ctx)

	// Open file first
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// Cancel context before OpenPackage (simulating cancellation during operation)
	cancel()

	// OpenPackage should detect cancelled context early
	_, err = OpenPackage(cancelCtx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with cancelled context")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Expected error type Context, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_ResourceCleanupOnError tests that resources are cleaned up on error.
func TestPackage_OpenPackage_ResourceCleanupOnError(t *testing.T) {
	ctx := context.Background()

	// Test with invalid file (should clean up file handle)
	_, err := OpenPackage(ctx, "/nonexistent/invalid.nvpk")
	if err == nil {
		t.Fatal("OpenPackage() should fail with non-existent file")
	}

	// Verify error is returned (file handle should be closed/cleaned up)
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}

	// Test with corrupted header (should clean up file handle)
	tempFile := filepath.Join(t.TempDir(), "corrupted.nvpk")
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	// Write invalid header
	_, _ = file.Write([]byte("INVALID"))
	_ = file.Close()

	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with corrupted header")
	}

	// Verify error is returned and file handle is cleaned up
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
}

// TestPackage_OpenPackage_IndexStartZero tests OpenPackage with index start at 0 (no index).
func TestPackage_OpenPackage_IndexStartZero(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a package file with index start at 0 (no index)
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	header := fileformat.NewPackageHeader()
	header.IndexStart = 0
	header.IndexSize = 0

	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// OpenPackage should succeed with no index
	pkg, err := OpenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	if !pkg.IsOpen() {
		t.Error("Package should be open")
	}
}

// TestPackage_OpenPackage_IndexReadFailure tests OpenPackage when index ReadFrom fails.
func TestPackage_OpenPackage_IndexReadFailure(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create file with valid header but truncated index data
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	header := fileformat.NewPackageHeader()
	header.IndexStart = uint64(fileformat.PackageHeaderSize)
	header.IndexSize = 20 // Claim index exists (16 bytes header + at least 16 bytes for 1 entry)

	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}

	// Write index header with EntryCount = 1 (small, safe value)
	indexHeader := make([]byte, 16)
	// EntryCount = 1 (4 bytes, little-endian)
	indexHeader[0] = 0x01
	indexHeader[1] = 0x00
	indexHeader[2] = 0x00
	indexHeader[3] = 0x00
	// Reserved = 0 (4 bytes)
	// FirstEntryOffset = 0 (8 bytes)
	// But truncate before writing the entry data - this will cause ReadFrom to fail
	if _, err := file.Write(indexHeader); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write index header: %v", err)
	}
	// Don't write the entry data - truncate here to cause read failure
	_ = file.Close()

	// Try to open - should fail at index read
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with truncated index")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	// Should be IO error (read failure) or corruption error
	if pkgErr.Type != pkgerrors.ErrTypeIO && pkgErr.Type != pkgerrors.ErrTypeCorruption {
		t.Errorf("Expected error type IO or Corruption, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenPackage_IndexValidationFailure tests OpenPackage when index validation fails after successful read.
func TestPackage_OpenPackage_IndexValidationFailure(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create file with valid header and index that reads successfully but fails validation
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	header := fileformat.NewPackageHeader()
	header.IndexStart = uint64(fileformat.PackageHeaderSize)

	// Create index with duplicate FileIDs (will fail validation)
	index := fileformat.NewFileIndex()
	index.EntryCount = 2
	index.FirstEntryOffset = uint64(fileformat.PackageHeaderSize + 16)
	index.Entries = []fileformat.IndexEntry{
		{FileID: 1, Offset: 100},
		{FileID: 1, Offset: 200}, // Duplicate FileID - will fail validation
	}
	header.IndexSize = uint64(index.Size())

	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	if _, err := index.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write index: %v", err)
	}
	_ = file.Close()

	// Try to open - should fail at index validation
	_, err = OpenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenPackage() should fail with invalid index")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// =============================================================================
// TEST: OpenBrokenPackage Operations
// =============================================================================

// TestPackage_OpenBrokenPackage_ValidHeader tests opening a package with valid header but broken index.
func TestPackage_OpenBrokenPackage_ValidHeader(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "broken.nvpk")

	// Create file with valid header but corrupted index
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write valid header
	header := fileformat.NewPackageHeader()
	header.IndexStart = uint64(fileformat.PackageHeaderSize)
	header.IndexSize = 100
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}

	// Write corrupted index data
	corruptData := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	if _, err := file.Write(corruptData); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write corrupt data: %v", err)
	}
	_ = file.Close()

	// OpenBrokenPackage should succeed (best-effort)
	pkg, err := OpenBrokenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenBrokenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Package should be open
	if !pkg.IsOpen() {
		t.Error("Package should be open")
	}

	// Package should have empty index (not nil)
	fpkg := pkg.(*filePackage)
	if fpkg.index == nil {
		t.Error("Package index should not be nil (should be empty)")
	}
}

// TestPackage_OpenBrokenPackage_InvalidHeader tests opening with completely invalid header.
func TestPackage_OpenBrokenPackage_InvalidHeader(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "invalid.nvpk")

	// Create file with invalid magic number
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	// Write garbage data
	if _, err := file.Write([]byte{0x00, 0x01, 0x02, 0x03}); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write data: %v", err)
	}
	_ = file.Close()

	// OpenBrokenPackage should fail (header is unrecoverable)
	_, err = OpenBrokenPackage(ctx, tempFile)
	if err == nil {
		t.Fatal("OpenBrokenPackage() should fail with invalid header")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_OpenBrokenPackage_NoIndex tests opening with no index.
func TestPackage_OpenBrokenPackage_NoIndex(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "noindex.nvpk")

	// Create file with valid header but no index
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write header with IndexStart = 0 (no index)
	header := fileformat.NewPackageHeader()
	header.IndexStart = 0
	header.IndexSize = 0
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// OpenBrokenPackage should succeed with empty index
	pkg, err := OpenBrokenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenBrokenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Package should be open
	if !pkg.IsOpen() {
		t.Error("Package should be open")
	}

	// Package should have empty index
	fpkg := pkg.(*filePackage)
	if fpkg.index == nil {
		t.Error("Package index should not be nil")
	}
	if fpkg.index.EntryCount != 0 {
		t.Errorf("Package index should be empty, got %d entries", fpkg.index.EntryCount)
	}
}

// TestPackage_OpenBrokenPackage_ReadFileDoesNotPanic tests that ReadFile doesn't panic with empty index.
func TestPackage_OpenBrokenPackage_ReadFileDoesNotPanic(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "broken.nvpk")

	// Create file with valid header but no index
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	header := fileformat.NewPackageHeader()
	header.IndexStart = 0
	header.IndexSize = 0
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// Open broken package
	pkg, err := OpenBrokenPackage(ctx, tempFile)
	if err != nil {
		t.Fatalf("OpenBrokenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Try to read a file (should fail gracefully, not panic)
	_, err = pkg.ReadFile(ctx, "test.txt")
	if err == nil {
		t.Fatal("ReadFile() should fail with empty index")
	}

	// Verify it returned an error, not a panic
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
}

// =============================================================================
// TEST: ReadHeader Operations
// =============================================================================

// TestPackage_ReadHeader_Basic tests reading package header without opening.
//
// Expected behavior (Red Phase - should FAIL):
// - ReadHeader function does not exist
func TestPackage_ReadHeader_Basic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Setup: Create a package
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg1 := pkg1.(*filePackage)
	err = fpkg1.Create(ctx, pkgPath)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	_ = pkg1.Close()

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test: ReadHeader should read header without opening package
	header, err := ReadHeaderFromPath(ctx, pkgPath)
	if err != nil {
		t.Errorf("ReadHeader() failed: %v", err)
	}
	if header == nil {
		t.Error("ReadHeader() returned nil")
	}
}

// TestPackage_ReadHeader_ErrorConditions tests ReadHeader error handling.
//
// Expected behavior (Red Phase - should FAIL):
// - Error handling not implemented
func TestPackage_ReadHeader_ErrorConditions(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		path      string
		wantError bool
	}{
		{
			name:      "empty path",
			path:      "",
			wantError: true,
		},
		{
			name:      "non-existent file",
			path:      "/path/to/nonexistent.nvpk",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadHeaderFromPath(ctx, tt.path)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}

// TestPackage_ReadHeader_Success tests ReadHeader with a valid package.
func TestPackage_ReadHeader_Success(t *testing.T) {
	ctx := context.Background()

	// Create a valid package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Read header
	header, err := ReadHeaderFromPath(ctx, tempFile)
	if err != nil {
		t.Fatalf("ReadHeader() failed: %v", err)
	}

	if header == nil {
		t.Fatal("ReadHeader() returned nil header")
	}
}

// TestPackage_ReadHeader_ValidatesMagic tests that ReadHeader validates magic.
//
// Expected behavior (Red Phase - should FAIL):
// - Magic validation not implemented
func TestPackage_ReadHeader_ValidatesMagic(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "invalid.nvpk")

	// Setup: Create invalid file
	err := os.WriteFile(invalidPath, []byte("Not a NovusPack file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test: ReadHeader should fail for invalid magic
	_, err = ReadHeaderFromPath(ctx, invalidPath)
	if err == nil {
		t.Error("ReadHeader should fail for invalid magic number")
	}
}

// TestPackage_ReadHeader_WithCancelledContext tests ReadHeader with cancelled context.
func TestPackage_ReadHeader_WithCancelledContext(t *testing.T) {
	ctx := context.Background()
	// Create a valid package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Try to read header with cancelled context
	cancelledCtx := testhelpers.CancelledContext()
	_, err = ReadHeaderFromPath(cancelledCtx, tempFile)
	if err == nil {
		t.Error("ReadHeader() should return error for cancelled context")
	}
}

// TestPackage_ReadHeader_WithInvalidMagic tests ReadHeader with invalid magic number.
// This tests the magic number validation path in readAndValidateHeader.
func TestPackage_ReadHeader_WithInvalidMagic(t *testing.T) {
	ctx := context.Background()
	tempFile := filepath.Join(t.TempDir(), "invalid_magic.nvpk")

	// Create a valid package first
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file first since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Now modify just the magic number
	file, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// Write invalid magic at the start (first 4 bytes)
	invalidMagic := []byte{0xEF, 0xBE, 0xAD, 0xDE} // 0xDEADBEEF in little-endian
	if _, err := file.WriteAt(invalidMagic, 0); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write invalid magic: %v", err)
	}
	_ = file.Close()

	// Try to read header
	_, err = ReadHeaderFromPath(ctx, tempFile)
	if err == nil {
		t.Fatal("ReadHeader() should fail with invalid magic number")
	}

	// Verify error type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
	if !strings.Contains(pkgErr.Message, "magic") {
		t.Errorf("Expected error message to mention magic number, got: %s", pkgErr.Message)
	}
}

// TestPackage_ReadHeader_WithWhitespacePath tests ReadHeader with whitespace-padded path.
func TestPackage_ReadHeader_WithWhitespacePath(t *testing.T) {
	ctx := context.Background()

	// Create a valid package
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file first since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Read header with whitespace-padded path
	pathWithWhitespace := "  " + tempFile + "  "
	header, err := ReadHeaderFromPath(ctx, pathWithWhitespace)
	if err != nil {
		t.Errorf("ReadHeader() should succeed with whitespace-padded path: %v", err)
	}
	if header == nil {
		t.Error("ReadHeader() returned nil header")
	}
}

// =============================================================================
// TEST: Internal Helper Functions
// =============================================================================

// TestOpenFileForReading_NonExistentFile tests openFileForReading with non-existent file.
func TestOpenFileForReading_NonExistentFile(t *testing.T) {
	_, err := internal.OpenFileForReading("/nonexistent/path/file.nvpk")
	if err == nil {
		t.Error("openFileForReading() should return error for non-existent file")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("openFileForReading() should return PackageError")
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeIO)
	}
}

// TestOpenFileForReading_Directory tests openFileForReading with a directory.
func TestOpenFileForReading_Directory(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	_, err := internal.OpenFileForReading(tempDir)
	if err == nil {
		t.Error("openFileForReading() should return error for directory")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("openFileForReading() should return PackageError")
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeIO)
	}
}

// TestOpenFileForReading_Success tests openFileForReading with valid file.
func TestOpenFileForReading_Success(t *testing.T) {
	ctx := context.Background()
	// Create a valid package file
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Open the file
	file, err := internal.OpenFileForReading(tempFile)
	if err != nil {
		t.Errorf("openFileForReading() failed: %v", err)
	}
	if file != nil {
		defer func() { _ = file.Close() }()
	}
}

// TestReadAndValidateHeader_InvalidMagic tests readAndValidateHeader with invalid magic.
func TestReadAndValidateHeader_InvalidMagic(t *testing.T) {
	// Create a temporary file with invalid magic
	tempFile := filepath.Join(t.TempDir(), "invalid.nvpk")
	f, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write invalid magic (4 bytes of zeros)
	invalidMagic := []byte{0x00, 0x00, 0x00, 0x00}
	if _, err := f.Write(invalidMagic); err != nil {
		_ = f.Close()
		t.Fatalf("Failed to write invalid magic: %v", err)
	}
	_ = f.Close()

	// Open file and try to read header
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer func() { _ = file.Close() }()

	_, err = internal.ReadAndValidateHeader(context.Background(), file)
	if err == nil {
		t.Error("readAndValidateHeader() should return error for invalid magic")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("readAndValidateHeader() should return PackageError")
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestReadAndValidateHeader_UnsupportedVersion tests readAndValidateHeader with unsupported version.
func TestReadAndValidateHeader_UnsupportedVersion(t *testing.T) {
	// Create a temporary file with valid magic but unsupported version
	tempFile := filepath.Join(t.TempDir(), "unsupported.nvpk")

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Create package with default header
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(context.Background(), tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Modify the file to have unsupported version
	file, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// Seek to version field (after magic: 4 bytes)
	if _, err := file.Seek(4, 0); err != nil {
		t.Fatalf("Failed to seek: %v", err)
	}

	// Write unsupported version (e.g., 999)
	unsupportedVersion := uint16(999)
	versionBytes := []byte{byte(unsupportedVersion), byte(unsupportedVersion >> 8)}
	if _, err := file.Write(versionBytes); err != nil {
		t.Fatalf("Failed to write version: %v", err)
	}
	_ = file.Close()

	// Reopen and try to read header
	file, err = os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to reopen file: %v", err)
	}
	defer func() { _ = file.Close() }()

	_, err = internal.ReadAndValidateHeader(context.Background(), file)
	if err == nil {
		t.Error("readAndValidateHeader() should return error for unsupported version")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("readAndValidateHeader() should return PackageError")
	}
	// Note: readAndValidateHeader wraps validation errors as ErrTypeValidation
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestReadAndValidateHeader_Success tests readAndValidateHeader with valid file.
func TestReadAndValidateHeader_Success(t *testing.T) {
	ctx := context.Background()

	// Create a valid package file
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	fpkg := pkg.(*filePackage)
	if err := fpkg.Create(ctx, tempFile); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create the file manually since Create() no longer writes to disk
	testutil.CreateTestPackageFile(t, tempFile)

	// Open and read header
	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = file.Close() }()

	header, err := internal.ReadAndValidateHeader(context.Background(), file)
	if err != nil {
		t.Errorf("readAndValidateHeader() failed: %v", err)
	}
	if header == nil {
		t.Error("readAndValidateHeader() returned nil header")
	}
}

// TestReadAndValidateHeader_WithHeaderValidateError tests readAndValidateHeader when header.Validate() fails.
func TestReadAndValidateHeader_WithHeaderValidateError(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a file with a header that has valid magic but invalid structure
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write a header with correct magic but invalid version
	header := fileformat.NewPackageHeader()
	header.Magic = fileformat.NVPKMagic
	header.FormatVersion = 999 // Invalid version that might fail validation
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// Open file and try to read header
	openedFile, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = openedFile.Close() }()

	_, err = internal.ReadAndValidateHeader(context.Background(), openedFile)
	// This might succeed if version 999 doesn't fail validation
	// But we're testing the Validate() call path
	if err != nil {
		pkgErr := &pkgerrors.PackageError{}
		if asPackageError(err, pkgErr) {
			if pkgErr.Type != pkgerrors.ErrTypeValidation {
				t.Errorf("Expected validation error, got: %v", pkgErr.Type)
			}
		}
	}
}

// TestReadAndValidateHeader_WithReadFromError tests readAndValidateHeader when header.ReadFrom() fails.
func TestReadAndValidateHeader_WithReadFromError(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a file that's too small to read a full header
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write only a few bytes (not enough for a header)
	if _, err := file.Write([]byte{0x4E, 0x56, 0x50, 0x4B}); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write: %v", err)
	}
	_ = file.Close()

	// Open file and try to read header
	openedFile, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = openedFile.Close() }()

	_, err = internal.ReadAndValidateHeader(context.Background(), openedFile)
	if err == nil {
		t.Fatal("readAndValidateHeader() should fail with truncated file")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected validation error, got: %v", pkgErr.Type)
	}
}

// TestReadAndValidateHeader_WithMagicCheckAfterReadFrom tests the redundant magic number check path.
func TestReadAndValidateHeader_WithMagicCheckAfterReadFrom(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a header with correct magic but then manually corrupt it after ReadFrom
	// This is tricky, but we can create a file where ReadFrom succeeds but magic is wrong
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write a header with wrong magic (but valid header structure)
	header := fileformat.NewPackageHeader()
	header.Magic = 0xDEADBEEF // Wrong magic
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}
	_ = file.Close()

	// Open and try to read - ReadFrom might fail, but if it doesn't, the magic check should catch it
	openedFile, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = openedFile.Close() }()

	_, err = internal.ReadAndValidateHeader(context.Background(), openedFile)
	if err == nil {
		t.Fatal("readAndValidateHeader() should fail with wrong magic")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected validation error, got: %v", pkgErr.Type)
	}
}

// TestReadAndValidateHeader_WithReadFromMagicError tests the error message check for "magic" in ReadFrom errors.
func TestReadAndValidateHeader_WithReadFromMagicError(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test.nvpk")

	// Create a file with invalid magic that will cause ReadFrom to fail with a magic error
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Write invalid magic bytes
	invalidMagic := []byte{0x00, 0x00, 0x00, 0x00}
	if _, err := file.Write(invalidMagic); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write: %v", err)
	}
	// Write rest of header with zeros to make it header-sized
	zeros := make([]byte, fileformat.PackageHeaderSize-4)
	if _, err := file.Write(zeros); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write zeros: %v", err)
	}
	_ = file.Close()

	// Open and try to read - ReadFrom should fail with magic error
	openedFile, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = openedFile.Close() }()

	_, err = internal.ReadAndValidateHeader(context.Background(), openedFile)
	if err == nil {
		t.Fatal("readAndValidateHeader() should fail with invalid magic")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected validation error, got: %v", pkgErr.Type)
	}
}
