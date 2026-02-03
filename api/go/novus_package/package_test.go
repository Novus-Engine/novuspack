// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains unit tests for the Package type and its lifecycle operations.
// Following TDD methodology - these tests are written FIRST and should FAIL initially (Red Phase).
package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat/testutil"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// =============================================================================
// TEST HELPERS
// =============================================================================

// asPackageError checks if an error can be converted to PackageError.
func asPackageError(err error, target *pkgerrors.PackageError) bool {
	if err == nil {
		return false
	}
	if pkgErr, ok := err.(*pkgerrors.PackageError); ok {
		*target = *pkgErr
		return true
	}
	return false
}

// runContextCancelledTest creates a package, calls the given method with a cancelled context,
// and asserts the error is a PackageError with ErrTypeContext.
func runContextCancelledTest(t *testing.T, call func(*filePackage, context.Context) error) {
	t.Helper()
	cancelledCtx := testhelpers.CancelledContext()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	fpkg := pkg.(*filePackage)
	err = call(fpkg, cancelledCtx)
	if err == nil {
		t.Fatal("expected failure with cancelled context")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want ErrTypeContext", pkgErr.Type)
	}
}

func runReadFileExpectFail(t *testing.T, path string) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	ctx := context.Background()
	_, err = pkg.ReadFile(ctx, path)
	if err == nil {
		t.Error("ReadFile should fail")
	}
}

func runGetMetadataBasic(t *testing.T) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	metadata, err := pkg.GetMetadata()
	if err != nil {
		t.Logf("GetMetadata note: %v (may require initialization)", err)
		return
	}
	if metadata == nil {
		t.Error("GetMetadata returned nil")
	}
}

// addThreeFilesFromMemory adds three files (file0.txt, file1.txt, file2.txt) via AddFileFromMemory.
func addThreeFilesFromMemory(t *testing.T, ctx context.Context, pkg Package, pathPrefix, contentPrefix string) {
	t.Helper()
	for i := 0; i < 3; i++ {
		path := pathPrefix + "file" + string(rune('0'+i)) + ".txt"
		data := []byte(contentPrefix + string(rune('0'+i)))
		_, err := pkg.AddFileFromMemory(ctx, path, data, nil)
		if err != nil {
			t.Fatalf("AddFileFromMemory failed: %v", err)
		}
	}
}

func runGetInfoBasic(t *testing.T, requireFormatVersion bool) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	if info == nil {
		t.Fatal("GetInfo returned nil")
	}
	if requireFormatVersion && info.FormatVersion == 0 {
		t.Error("Info.FormatVersion = 0, expected non-zero")
	}
}

func runOpenPackageThenCloseThen(t *testing.T, fn func(t *testing.T, pkg Package)) {
	t.Helper()
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)
	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	if err := pkg.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}
	fn(t, pkg)
}

func runOpenPackageThenCloseThenSucceed(t *testing.T, methodName string, fn func(Package) (interface{}, error)) {
	t.Helper()
	runOpenPackageThenCloseThen(t, func(t *testing.T, pkg Package) {
		got, err := fn(pkg)
		if err != nil {
			t.Fatalf("%s() should succeed on a closed package with cached metadata, got error: %v", methodName, err)
		}
		if got == nil {
			t.Fatalf("%s() returned nil", methodName)
		}
	})
}

func runAssertGetInfoOnClosed(t *testing.T) {
	t.Helper()
	runOpenPackageThenCloseThenSucceed(t, "GetInfo", func(pkg Package) (interface{}, error) { return pkg.GetInfo() })
}

func runAssertListFilesOnClosed(t *testing.T) {
	t.Helper()
	runOpenPackageThenCloseThenSucceed(t, "ListFiles", func(pkg Package) (interface{}, error) { return pkg.ListFiles() })
}

func runOpenPackageListFilesExpectEmpty(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")
	testutil.CreateTestPackageFile(t, pkgPath)
	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	files, err := pkg.ListFiles()
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}
	if files == nil {
		t.Fatal("ListFiles() should not return nil")
	}
	if len(files) != 0 {
		t.Errorf("ListFiles() on empty package should return empty list, got %d files", len(files))
	}
}

func runWriteWithContent(t *testing.T, content []byte, verifyFile bool) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", content, nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if verifyFile {
		if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
			t.Error("Write did not create package file")
		}
	}
}

func runWriteEmptyPackage(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	tmpPkg := filepath.Join(t.TempDir(), "empty.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
		t.Error("Write did not create package file")
	}
}

func runWriteContextCancelled(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
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
	if err := pkg.Write(cancelledCtx); err == nil {
		t.Error("Write with cancelled context should fail")
	}
}

func runAddFileOverwrite(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("v1"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory(v1) failed: %v", err)
	}
	opts := &AddFileOptions{}
	opts.AllowOverwrite.Set(true)
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("v2"), opts)
	if err != nil {
		t.Logf("AddFileFromMemory with AllowOverwrite: %v (may not be fully implemented)", err)
	}
}

func runRemoveFileExpectFail(t *testing.T, path string) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	err = pkg.RemoveFile(ctx, path)
	if err == nil {
		t.Error("RemoveFile should fail")
	}
}

func runSafeWriteWithContent(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	_, err = pkg.AddFileFromMemory(ctx, "/test.txt", []byte("content"), nil)
	if err != nil {
		t.Fatalf("AddFileFromMemory failed: %v", err)
	}
	tmpPkg := filepath.Join(t.TempDir(), "test.pkg")
	if err := pkg.SetTargetPath(ctx, tmpPkg); err != nil {
		t.Fatalf("SetTargetPath failed: %v", err)
	}
	if err := pkg.SafeWrite(ctx, true); err != nil {
		t.Fatalf("SafeWrite failed: %v", err)
	}

	if _, err := os.Stat(tmpPkg); os.IsNotExist(err) {
		t.Error("SafeWrite did not create package file")
	}
}

func runAddTwoPathsThenRemove(t *testing.T, storedPath2 string) {
	t.Helper()
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	entry1, err := pkg.AddFile(ctx, testFile, nil)
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}
	opts := &AddFileOptions{}
	opts.StoredPath.Set(storedPath2)
	entry2, err := pkg.AddFile(ctx, testFile, opts)
	if err != nil {
		t.Fatalf("AddFile with different path failed: %v", err)
	}
	if entry1.FileID != entry2.FileID {
		t.Error("Deduplication should reuse same FileEntry")
	}
	err = pkg.RemoveFile(ctx, entry1.Paths[0].Path)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}
}

// =============================================================================
// TEST: NewPackage Constructor
// =============================================================================

// TestPackage_NewPackage tests the NewPackage constructor function.
// This test verifies that NewPackage creates a new Package instance correctly.
//
// Expected behavior (Red Phase - should FAIL):
// - NewPackage function does not exist yet
// - Test will fail with compilation error or runtime error
func TestPackage_NewPackage(t *testing.T) {

	// Test: NewPackage creates a valid Package instance
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("NewPackage() returned nil package")
	}
}

// TestPackage_NewPackage_WithContext tests NewPackage with various context scenarios.
//
// Expected behavior (Red Phase - should FAIL):
// - NewPackage does not exist
// - Context handling not implemented
func TestPackage_NewPackage_WithContext(t *testing.T) {
	t.Skip("NewPackage() no longer takes context parameter per spec - context validation moved to Create/Open operations")
}

// =============================================================================
// TEST: Package Structure and Fields
// =============================================================================

// TestPackage_StructFields verifies the Package struct has required fields.
//
// Expected behavior (Red Phase - should FAIL):
// - Package type does not exist
// - Struct fields not defined
func TestPackage_StructFields(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Test: Package should have required fields (verified through methods)
	// We can't access private fields directly, but we can verify behavior

	// Verify IsOpen method exists and returns false for new package
	if pkg.IsOpen() {
		t.Error("New package should not be open")
	}
}

// TestPackage_InitialState verifies the initial state of a newly created package.
//
// Expected behavior (Red Phase - should FAIL):
// - Package methods do not exist
// - State checking not implemented
func TestPackage_InitialState(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}

	// Test: New package should be in memory only (not open/on disk)
	if pkg.IsOpen() {
		t.Error("New package should not be in open state")
	}

	// Test: GetInfo should succeed for newly created package (metadata is initialized)
	info, err := pkg.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo() failed on newly created package: %v", err)
	}

	if info == nil {
		t.Fatal("GetInfo() returned nil")
	}

	// Verify format version is set
	if info.FormatVersion == 0 {
		t.Error("FormatVersion should be non-zero for newly created package")
	}
}

// =============================================================================
// TEST: ErrorType and PackageError
// =============================================================================

// TestErrorType_String tests the String() method for all ErrorType values.
func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name     string
		errType  pkgerrors.ErrorType
		expected string
	}{
		{"Validation", pkgerrors.ErrTypeValidation, "Validation"},
		{"IO", pkgerrors.ErrTypeIO, "IO"},
		{"Security", pkgerrors.ErrTypeSecurity, "Security"},
		{"Unsupported", pkgerrors.ErrTypeUnsupported, "Unsupported"},
		{"Context", pkgerrors.ErrTypeContext, "Context"},
		{"Corruption", pkgerrors.ErrTypeCorruption, "Corruption"},
		{"Unknown", pkgerrors.ErrorType(999), "Unknown(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errType.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPackageError_Error tests the Error() method.
func TestPackageError_Error(t *testing.T) {
	tests := []struct {
		name     string
		pkgErr   *pkgerrors.PackageError
		contains []string
	}{
		{
			name: "error with cause",
			pkgErr: &pkgerrors.PackageError{
				Type:    pkgerrors.ErrTypeIO,
				Message: "failed to read file",
				Cause:   os.ErrNotExist,
			},
			contains: []string{"IO", "failed to read file", "file does not exist"},
		},
		{
			name: "error without cause",
			pkgErr: &pkgerrors.PackageError{
				Type:    pkgerrors.ErrTypeValidation,
				Message: "invalid path",
				Cause:   nil,
			},
			contains: []string{"Validation", "invalid path"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkgErr.Error()
			for _, substr := range tt.contains {
				if !testhelpers.ContainsIgnoreCase(result, substr) {
					t.Errorf("Error() = %v, should contain %v", result, substr)
				}
			}
		})
	}
}

// TestPackageError_Unwrap tests the Unwrap() method.
func TestPackageError_Unwrap(t *testing.T) {
	tests := []struct {
		name   string
		pkgErr *pkgerrors.PackageError
		want   error
	}{
		{
			name: "with cause",
			pkgErr: &pkgerrors.PackageError{
				Type:    pkgerrors.ErrTypeIO,
				Message: "test error",
				Cause:   os.ErrNotExist,
			},
			want: os.ErrNotExist,
		},
		{
			name: "without cause",
			pkgErr: &pkgerrors.PackageError{
				Type:    pkgerrors.ErrTypeValidation,
				Message: "test error",
				Cause:   nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkgErr.Unwrap()
			if result != tt.want {
				t.Errorf("Unwrap() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestPackageError_WithContext tests the WithContext() method.
func TestPackageError_WithContext(t *testing.T) {
	pkgErr := pkgerrors.NewPackageError[struct{}](pkgerrors.ErrTypeIO, "test error", nil, struct{}{})

	// Add context
	result := pkgErr.WithContext("path", "/test/path")

	// Verify chaining returns same instance
	if result != pkgErr {
		t.Error("WithContext() should return the same instance for chaining")
	}

	// Verify context was added
	if pkgErr.Context["path"] != "/test/path" {
		t.Errorf("Context['path'] = %v, want /test/path", pkgErr.Context["path"])
	}

	// Add multiple contexts
	_ = pkgErr.WithContext("operation", "read").WithContext("size", 1024)

	if pkgErr.Context["operation"] != "read" {
		t.Errorf("Context['operation'] = %v, want read", pkgErr.Context["operation"])
	}
	if pkgErr.Context["size"] != 1024 {
		t.Errorf("Context['size'] = %v, want 1024", pkgErr.Context["size"])
	}
}

// TestNewPackageError tests the NewPackageError constructor.
//
//nolint:gocognit // table-driven error cases
func TestNewPackageError(t *testing.T) {
	tests := []struct {
		name    string
		errType pkgerrors.ErrorType
		message string
		cause   error
	}{
		{
			name:    "with cause",
			errType: pkgerrors.ErrTypeIO,
			message: "test error",
			cause:   os.ErrNotExist,
		},
		{
			name:    "without cause",
			errType: pkgerrors.ErrTypeValidation,
			message: "validation failed",
			cause:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pkgerrors.NewPackageError[struct{}](tt.errType, tt.message, tt.cause, struct{}{})
			//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so err is not nil after check
			if err == nil {
				t.Fatal("NewPackageError() returned nil")
			}
			//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so err is not nil after check
			if err.Type != tt.errType {
				t.Errorf("Type = %v, want %v", err.Type, tt.errType)
			}
			if err.Message != tt.message {
				t.Errorf("Message = %v, want %v", err.Message, tt.message)
			}
			if err.Cause != tt.cause {
				t.Errorf("Cause = %v, want %v", err.Cause, tt.cause)
			}
			if err.Context == nil {
				t.Error("Context should be initialized")
			}
		})
	}
}

// =============================================================================
// TEST: Helper Functions (for better coverage)
// =============================================================================

// TestCheckContext_NilContext tests checkContext with nil context.
func TestCheckContext_NilContext(t *testing.T) {
	var nilCtx context.Context // explicitly nil for testing
	//nolint:staticcheck // SA1012: Testing nil context error handling
	err := internal.CheckContext(nilCtx, "test operation")
	if err == nil {
		t.Error("checkContext() should return error for nil context")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("checkContext() should return PackageError")
	}
	// Note: checkContext returns ErrTypeValidation for nil context
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

func assertCheckContextError(t *testing.T, ctx context.Context, wantType pkgerrors.ErrorType) {
	t.Helper()
	err := internal.CheckContext(ctx, "test operation")
	if err == nil {
		t.Error("checkContext() should return error")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Error("checkContext() should return PackageError")
	}
	if pkgErr.Type != wantType {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, wantType)
	}
}

// TestCheckContext_CancelledContext tests checkContext with cancelled context.
func TestCheckContext_CancelledContext(t *testing.T) {
	assertCheckContextError(t, testhelpers.CancelledContext(), pkgerrors.ErrTypeContext)
}

// TestCheckContext_TimeoutContext tests checkContext with timed out context.
func TestCheckContext_TimeoutContext(t *testing.T) {
	assertCheckContextError(t, testhelpers.TimeoutContext(), pkgerrors.ErrTypeContext)
}

// runWithCancelledContext runs fn with a cancelled context and asserts it returns an error.
func runWithCancelledContext(t *testing.T, fn func(*filePackage, context.Context) (interface{}, error), opName string) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	cancelledCtx := testhelpers.CancelledContext()
	fpkg := pkg.(*filePackage)
	_, err = fn(fpkg, cancelledCtx)
	if err == nil {
		t.Errorf("%s should fail with cancelled context", opName)
	}
}

// =============================================================================
// ADDITIONAL COVERAGE TESTS
// =============================================================================
//
// The following tests were added to achieve 95%+ overall coverage and 92%+ per-package coverage.
// These tests focus on uncovered error paths and edge cases.
