// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains unit tests for read-only package operations,
// specifically testing the readOnlyPackage wrapper type.
package novus_package

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat/testutil"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestReadOnlyPackage_CreateWithOptions tests that CreateWithOptions returns a security error.
func TestReadOnlyPackage_CreateWithOptions(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open package in read-only mode
	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Attempt CreateWithOptions - should fail with security error
	newPath := filepath.Join(tmpDir, "new.nvpk")
	err = pkg.CreateWithOptions(ctx, newPath, nil)
	if err == nil {
		t.Fatal("CreateWithOptions() should fail on read-only package")
	}

	// Verify it's a PackageError with Security type
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeSecurity {
		t.Errorf("Expected error type Security, got: %v", pkgErr.Type)
	}

	// Verify error message mentions the operation
	if pkgErr.Message == "" {
		t.Error("Error message should not be empty")
	}
}

// TestReadOnlyPackage_CloseWithCleanup tests that CloseWithCleanup delegates to inner package.
func TestReadOnlyPackage_CloseWithCleanup(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open package in read-only mode
	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}

	// Verify package is open
	if !pkg.IsOpen() {
		t.Error("Package should be open")
	}

	// Call CloseWithCleanup - should succeed
	err = pkg.CloseWithCleanup(ctx)
	if err != nil {
		t.Errorf("CloseWithCleanup() failed: %v", err)
	}

	// Verify package is now closed
	if pkg.IsOpen() {
		t.Error("Package should be closed after CloseWithCleanup()")
	}
}

// TestReadOnlyPackage_GetPath tests that GetPath returns the correct path.
func TestReadOnlyPackage_GetPath(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open package in read-only mode
	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Get path - should return the package path
	returnedPath := pkg.GetPath()
	if returnedPath != pkgPath {
		t.Errorf("GetPath() = %q, want %q", returnedPath, pkgPath)
	}
}

// TestReadOnlyPackage_IsReadOnly tests that IsReadOnly returns true.
func TestReadOnlyPackage_IsReadOnly(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	// Open package in read-only mode
	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Verify IsReadOnly returns true
	if !pkg.IsReadOnly() {
		t.Error("IsReadOnly() should return true for read-only package")
	}
}

// TestReadOnlyPackage_WriteOperationsAllFail tests that all write operations fail with security errors.
func TestReadOnlyPackage_WriteOperationsAllFail(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	tests := []struct {
		name string
		op   func() error
	}{
		{
			name: "Create",
			op: func() error {
				return pkg.Create(ctx, filepath.Join(tmpDir, "new.nvpk"))
			},
		},
		// TODO(Priority 2): Re-enable these tests after AddFile/RemoveFile are fully implemented
		// {
		// 	name: "StageFile",
		// 	op: func() error {
		// 		return pkg.StageFile(ctx, "test.txt", []byte("data"), nil)
		// 	},
		// },
		// {
		// 	name: "UnstageFile",
		// 	op: func() error {
		// 		return pkg.UnstageFile(ctx, "test.txt")
		// 	},
		// },
		{
			name: "Write",
			op: func() error {
				return pkg.Write(ctx)
			},
		},
		{
			name: "SafeWrite",
			op: func() error {
				return pkg.SafeWrite(ctx, false)
			},
		},
		{
			name: "FastWrite",
			op: func() error {
				return pkg.FastWrite(ctx)
			},
		},
		{
			name: "Defragment",
			op: func() error {
				return pkg.Defragment(ctx)
			},
		},
		{
			name: "SetComment",
			op: func() error {
				return pkg.SetComment("test comment")
			},
		},
		{
			name: "ClearComment",
			op: func() error {
				return pkg.ClearComment()
			},
		},
		{
			name: "SetAppID",
			op: func() error {
				return pkg.SetAppID(0x12345678)
			},
		},
		{
			name: "ClearAppID",
			op: func() error {
				return pkg.ClearAppID()
			},
		},
		{
			name: "SetVendorID",
			op: func() error {
				return pkg.SetVendorID(0x87654321)
			},
		},
		{
			name: "ClearVendorID",
			op: func() error {
				return pkg.ClearVendorID()
			},
		},
		{
			name: "SetPackageIdentity",
			op: func() error {
				return pkg.SetPackageIdentity(0x87654321, 0x12345678)
			},
		},
		{
			name: "ClearPackageIdentity",
			op: func() error {
				return pkg.ClearPackageIdentity()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.op()
			if err == nil {
				t.Fatalf("%s should fail on read-only package", tt.name)
			}

			// Verify it's a security error
			pkgErr := &pkgerrors.PackageError{}
			if !asPackageError(err, pkgErr) {
				t.Fatalf("Expected PackageError, got: %T", err)
			}
			if pkgErr.Type != pkgerrors.ErrTypeSecurity {
				t.Errorf("Expected error type Security, got: %v", pkgErr.Type)
			}
		})
	}
}

// TestReadOnlyPackage_ReadOperationsWork tests that read operations still work on read-only packages.
func TestReadOnlyPackage_ReadOperationsWork(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Test GetInfo
	info, err := pkg.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() failed: %v", err)
	}
	if info == nil {
		t.Error("GetInfo() returned nil")
	}

	// Test GetMetadata
	metadata, err := pkg.GetMetadata()
	if err != nil {
		t.Errorf("GetMetadata() failed: %v", err)
	}
	if metadata == nil {
		t.Error("GetMetadata() returned nil")
	}

	// Test ListFiles
	files, err := pkg.ListFiles()
	if err != nil {
		t.Errorf("ListFiles() failed: %v", err)
	}
	if files == nil {
		t.Error("ListFiles() returned nil")
	}

	// Test Validate
	err = pkg.Validate(ctx)
	if err != nil {
		t.Errorf("Validate() failed: %v", err)
	}

	// Test IsOpen
	if !pkg.IsOpen() {
		t.Error("IsOpen() should return true")
	}

	// Test GetComment (read operation)
	_ = pkg.GetComment()

	// Test HasComment (read operation)
	_ = pkg.HasComment()

	// Test GetAppID (read operation)
	_ = pkg.GetAppID()

	// Test HasAppID (read operation)
	_ = pkg.HasAppID()

	// Test GetVendorID (read operation)
	_ = pkg.GetVendorID()

	// Test HasVendorID (read operation)
	_ = pkg.HasVendorID()

	// Test GetPackageIdentity (read operation)
	_, _ = pkg.GetPackageIdentity()
}
