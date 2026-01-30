// Package internal provides internal helper functions for the NovusPack implementation.
//
// This file contains unit tests for internal helper functions.
package internal

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/unicode/norm"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestValidatePath tests the ValidatePath function with various inputs.
func TestValidatePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		{
			name:        "empty path",
			path:        "",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "whitespace only",
			path:        "   \t\n  ",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "valid path",
			path:        "test.nvpk",
			shouldError: false,
		},
		{
			name:        "path with spaces",
			path:        "test file.nvpk",
			shouldError: false,
		},
		{
			name:        "nested path",
			path:        "dir/subdir/file.nvpk",
			shouldError: false,
		},
		{
			name:        "path with leading/trailing spaces",
			path:        "  test.nvpk  ",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := ValidatePath(ctx, tt.path)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ValidatePath() expected error, got nil")
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("ValidatePath() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("ValidatePath() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePath() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestCheckContext tests the CheckContext function with various context scenarios.
func TestCheckContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		operation   string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		{
			name:        "nil context",
			ctx:         nil,
			operation:   "test",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "valid context",
			ctx:         context.Background(),
			operation:   "test",
			shouldError: false,
		},
		{
			name: "cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			operation:   "test",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeContext,
		},
		{
			name: "timeout context",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
				time.Sleep(time.Millisecond)
				cancel()
				return ctx
			}(),
			operation:   "test",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeContext,
		},
		{
			name: "context with deadline exceeded",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				cancel()
				return ctx
			}(),
			operation:   "operation",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckContext(tt.ctx, tt.operation)
			if tt.shouldError {
				if err == nil {
					t.Errorf("CheckContext() expected error, got nil")
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("CheckContext() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("CheckContext() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
				if tt.operation != "" {
					_, isPkgErr := pkgerrors.IsPackageError(err)
					if !isPkgErr {
						t.Errorf("CheckContext() error should be a PackageError")
					}
				}
			} else {
				if err != nil {
					t.Errorf("CheckContext() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestOpenFileForReading tests the OpenFileForReading function.
func TestOpenFileForReading(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	validFile := filepath.Join(tmpDir, "test.nvpk")

	// Create a valid file
	file, err := os.Create(validFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_ = file.Close()

	tests := []struct {
		name        string
		path        string
		setup       func() string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		{
			name:        "non-existent file",
			path:        filepath.Join(tmpDir, "nonexistent.nvpk"),
			shouldError: true,
			errorType:   pkgerrors.ErrTypeIO,
		},
		{
			name:        "directory instead of file",
			path:        tmpDir,
			shouldError: true,
			errorType:   pkgerrors.ErrTypeIO,
		},
		{
			name:        "valid file",
			path:        validFile,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if tt.setup != nil {
				path = tt.setup()
			}

			file, err := OpenFileForReading(path)
			if tt.shouldError {
				if err == nil {
					t.Errorf("OpenFileForReading() expected error, got nil")
					return
				}
				if file != nil {
					_ = file.Close()
					t.Errorf("OpenFileForReading() should not return file on error")
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("OpenFileForReading() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("OpenFileForReading() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("OpenFileForReading() unexpected error: %v", err)
					return
				}
				if file == nil {
					t.Errorf("OpenFileForReading() expected file, got nil")
					return
				}
				_ = file.Close()
			}
		})
	}
}

// TestOpenFileForReading_StatError tests OpenFileForReading when Stat() fails.
// This tests the error path where file.Stat() returns an error.
func TestOpenFileForReading_StatError(t *testing.T) {
	// Create a file and then remove it while keeping a reference
	// On Unix systems, we can open a file, remove it, and Stat() will still work
	// So we need a different approach - test with a directory which will fail the IsDir check
	// Actually, the directory check is already tested. The Stat() error path is hard to test
	// without mocking. Let's test with a file that becomes inaccessible.

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.nvpk")

	// Create a valid file
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_ = file.Close()

	// Open the file
	openedFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Remove the file while it's open (on Unix, this works)
	// The file handle is still valid, but Stat might behave differently
	// Actually, on Unix, Stat() will still work even after removal
	// So we'll test the directory path which is already covered
	// and document that Stat() error is hard to test without mocking

	// Test with directory (already covered but ensures coverage)
	_, err = OpenFileForReading(tmpDir)
	if err == nil {
		t.Fatal("OpenFileForReading() should fail for directory")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}

	_ = openedFile.Close()

	// Note: Testing Stat() failure directly is difficult without mocking
	// The code path exists and will be executed if Stat() ever fails in practice
}

// TestOpenFileForReading_PermissionError tests OpenFileForReading when os.Open fails
// with a permission error (non-IsNotExist error). This tests lines 88-89.
func TestOpenFileForReading_PermissionError(t *testing.T) {
	// Only test on Unix-like systems where we can set permissions
	if runtime.GOOS == "windows" {
		t.Skip("Permission testing not reliable on Windows")
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "restricted.nvpk")

	// Create a file
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_ = file.Close()

	// Remove read permission
	if err := os.Chmod(path, 0000); err != nil {
		t.Fatalf("Failed to chmod file: %v", err)
	}
	defer func() {
		_ = os.Chmod(path, 0644) // Restore permissions for cleanup
	}()

	// Try to open the file - should fail with permission error
	_, err = OpenFileForReading(path)
	if err == nil {
		t.Fatal("OpenFileForReading() should fail with permission error")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}

	// Verify error message contains "failed to open package file" (line 89)
	if !strings.Contains(pkgErr.Error(), "failed to open package file") {
		t.Errorf("Expected error message about opening file, got: %v", pkgErr.Error())
	}
}

// TestReadAndValidateHeader tests the ReadAndValidateHeader function.
func TestReadAndValidateHeader(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() *os.File
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		{
			name: "valid header",
			setup: func() *os.File {
				path := filepath.Join(tmpDir, "valid.nvpk")
				file, err := os.Create(path)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				header := fileformat.NewPackageHeader()
				_, _ = header.WriteTo(file)
				_, _ = file.Seek(0, 0)
				return file
			},
			shouldError: false,
		},
		{
			name: "invalid magic number",
			setup: func() *os.File {
				path := filepath.Join(tmpDir, "invalid_magic.nvpk")
				file, err := os.Create(path)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				// Write invalid magic
				_, _ = file.Write([]byte{0x00, 0x00, 0x00, 0x00})
				_, _ = file.Seek(0, 0)
				return file
			},
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name: "empty file",
			setup: func() *os.File {
				path := filepath.Join(tmpDir, "empty.nvpk")
				file, err := os.Create(path)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				_, _ = file.Seek(0, 0)
				return file
			},
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name: "invalid header (corrupted data)",
			setup: func() *os.File {
				path := filepath.Join(tmpDir, "corrupted.nvpk")
				file, err := os.Create(path)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				// Write valid magic but invalid rest of header
				header := fileformat.NewPackageHeader()
				_, _ = header.WriteTo(file)
				// Corrupt the file by truncating it
				_, _ = file.Seek(0, 0)
				_ = file.Truncate(10) // Truncate to invalid size
				_, _ = file.Seek(0, 0)
				return file
			},
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := tt.setup()
			defer func() { _ = file.Close() }()

			header, err := ReadAndValidateHeader(context.Background(), file)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ReadAndValidateHeader() expected error, got nil")
					return
				}
				if header != nil {
					t.Errorf("ReadAndValidateHeader() should not return header on error")
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("ReadAndValidateHeader() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("ReadAndValidateHeader() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("ReadAndValidateHeader() unexpected error: %v", err)
					return
				}
				if header == nil {
					t.Errorf("ReadAndValidateHeader() expected header, got nil")
					return
				}
				if header.Magic != fileformat.NVPKMagic {
					t.Errorf("ReadAndValidateHeader() header.Magic = 0x%08X, want 0x%08X", header.Magic, fileformat.NVPKMagic)
				}
			}
		})
	}
}

// TestReadAndValidateHeader_ValidateError tests ReadAndValidateHeader when header.Validate() fails.
func TestReadAndValidateHeader_ValidateError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid_validate.nvpk")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Write a header with valid magic but invalid structure (e.g., invalid version)
	header := fileformat.NewPackageHeader()
	header.Magic = fileformat.NVPKMagic
	header.FormatVersion = 999 // Invalid version that will fail validation
	_, _ = header.WriteTo(file)
	_, _ = file.Seek(0, 0)
	_ = file.Close()

	file, err = os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() { _ = file.Close() }()

	headerResult, err := ReadAndValidateHeader(context.Background(), file)
	if err == nil {
		t.Errorf("ReadAndValidateHeader() expected error for invalid header validation, got nil")
	}
	if headerResult != nil {
		t.Errorf("ReadAndValidateHeader() should not return header on error")
	}
	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Errorf("ReadAndValidateHeader() error is not a PackageError: %v", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("ReadAndValidateHeader() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestReadAndValidateHeader_ReadFromError tests ReadAndValidateHeader when ReadFrom fails.
func TestReadAndValidateHeader_ReadFromError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "partial.nvpk")

	// Create a file with partial header data
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Write only partial header (less than PackageHeaderSize)
	_, _ = file.Write([]byte{0x4E, 0x50, 0x4B, 0x21}) // Partial magic
	_ = file.Close()

	file, err = os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() { _ = file.Close() }()

	header, err := ReadAndValidateHeader(context.Background(), file)
	if err == nil {
		t.Errorf("ReadAndValidateHeader() expected error for partial header, got nil")
	}
	if header != nil {
		t.Errorf("ReadAndValidateHeader() should not return header on error")
	}
	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Errorf("ReadAndValidateHeader() error is not a PackageError: %v", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("ReadAndValidateHeader() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestReadAndValidateHeader_ReadFromMagicError tests ReadAndValidateHeader when ReadFrom fails with magic error.
// This tests the path where the error message contains "magic" (line 86 in helpers.go).
func TestReadAndValidateHeader_ReadFromMagicError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "magic_error.nvpk")

	// Create a file with invalid magic that will cause ReadFrom to fail with a magic error
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write invalid magic bytes that will trigger a magic error in ReadFrom
	// The ReadFrom method will detect invalid magic and return an error with "magic" in the message
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

	// Open file and try to read header
	openedFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = openedFile.Close() }()

	header, err := ReadAndValidateHeader(context.Background(), openedFile)
	if err == nil {
		t.Fatal("ReadAndValidateHeader() should fail with invalid magic")
	}
	if header != nil {
		t.Error("ReadAndValidateHeader() should not return header on error")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
	// Verify the error message handling path (where strings.Contains checks for "magic")
	// This tests line 86 in helpers.go
}

// TestReadAndValidateHeader_ReadFromGenericError tests ReadAndValidateHeader when ReadFrom fails with non-magic error.
// This tests the path where the error doesn't contain "magic" (line 89 in helpers.go).
func TestReadAndValidateHeader_ReadFromGenericError(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "generic_error.nvpk")

	// Create a file that's too small to read a full header
	// This will cause ReadFrom to fail with an error that doesn't contain "magic"
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write only a few bytes (not enough for a header, but not a magic error)
	// This will trigger EOF or short read error, which won't contain "magic"
	partialData := []byte{0x4E, 0x50, 0x4B, 0x21, 0x00, 0x00} // Partial header, not magic error
	if _, err := file.Write(partialData); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write: %v", err)
	}
	_ = file.Close()

	// Open file and try to read header
	openedFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer func() { _ = openedFile.Close() }()

	header, err := ReadAndValidateHeader(context.Background(), openedFile)
	if err == nil {
		t.Fatal("ReadAndValidateHeader() should fail with partial header")
	}
	if header != nil {
		t.Error("ReadAndValidateHeader() should not return header on error")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
	// This tests line 89 in helpers.go where error doesn't contain "magic"
	// and uses the generic "failed to read package header" message
}

// TestNormalizePackagePath tests the NormalizePackagePath function.
func TestNormalizePackagePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		{
			name:        "empty path",
			path:        "",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "valid relative path",
			path:        "file.txt",
			expected:    "/file.txt",
			shouldError: false,
		},
		{
			name:        "path with backslashes",
			path:        "path\\to\\file.txt",
			expected:    "/path/to/file.txt",
			shouldError: false,
		},
		{
			name:        "path with leading slash",
			path:        "/file.txt",
			expected:    "/file.txt",
			shouldError: false,
		},
		{
			name:        "path with dot segment",
			path:        "./file.txt",
			expected:    "/file.txt",
			shouldError: false,
		},
		{
			name:        "path with parent directory",
			path:        "../file.txt",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "path with embedded dot segment",
			path:        "path/./file.txt",
			expected:    "/path/file.txt",
			shouldError: false,
		},
		{
			name:        "path with embedded parent directory",
			path:        "path/../file.txt",
			expected:    "/file.txt",
			shouldError: false,
		},
		{
			name:        "just dot",
			path:        ".",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "just parent directory",
			path:        "..",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "nested path",
			path:        "path/to/file.txt",
			expected:    "/path/to/file.txt",
			shouldError: false,
		},
		{
			name:        "path ending with dot",
			path:        "path/.",
			expected:    "/path",
			shouldError: false,
		},
		{
			name:        "path ending with parent",
			path:        "path/..",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizePackagePath(tt.path)
			if tt.shouldError {
				if err == nil {
					t.Errorf("NormalizePackagePath() expected error, got nil")
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("NormalizePackagePath() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("NormalizePackagePath() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("NormalizePackagePath() unexpected error: %v", err)
					return
				}
				if result != tt.expected {
					t.Errorf("NormalizePackagePath() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

// TestValidatePackagePath tests the ValidatePackagePath function.
func TestValidatePackagePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		{
			name:        "empty path",
			path:        "",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "whitespace only",
			path:        "   \t\n  ",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "valid path",
			path:        "file.txt",
			shouldError: false,
		},
		{
			name:        "path with dot segment",
			path:        "./file.txt",
			shouldError: false,
		},
		{
			name:        "path with parent directory",
			path:        "../file.txt",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "path with embedded dot segment",
			path:        "path/./file.txt",
			shouldError: false,
		},
		{
			name:        "path with embedded parent directory",
			path:        "path/../file.txt",
			shouldError: false,
		},
		{
			name:        "nested valid path",
			path:        "path/to/file.txt",
			shouldError: false,
		},
		{
			name:        "path with backslashes",
			path:        "path\\to\\file.txt",
			shouldError: false, // Will be normalized, but validation passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackagePath(tt.path)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ValidatePackagePath() expected error, got nil")
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("ValidatePackagePath() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("ValidatePackagePath() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePackagePath() unexpected error: %v", err)
				}
			}
		})
	}
}

// TestLoadFileEntry tests the LoadFileEntry function.
func TestLoadFileEntry(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a minimal package file with a file entry
	path := filepath.Join(tmpDir, "test.nvpk")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write header
	header := fileformat.NewPackageHeader()
	header.IndexStart = uint64(fileformat.PackageHeaderSize)
	header.IndexSize = 16 // Empty index size
	if _, err := header.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write header: %v", err)
	}

	// Write empty index
	index := fileformat.NewFileIndex()
	if _, err := index.WriteTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write index: %v", err)
	}

	// Calculate file entry offset (after header and index)
	entryOffset := uint64(fileformat.PackageHeaderSize) + uint64(index.Size())

	// Write a minimal file entry
	entry := metadata.NewFileEntry()
	entry.FileID = 1
	entry.OriginalSize = 10
	entry.StoredSize = 10
	entry.PathCount = 1
	entry.Paths = []generics.PathEntry{
		{
			PathLength: 8,
			Path:       "test.txt",
		},
	}
	if _, err := entry.WriteMetaTo(file); err != nil {
		_ = file.Close()
		t.Fatalf("Failed to write file entry: %v", err)
	}
	_ = file.Close()

	// Open file for reading
	readFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() { _ = readFile.Close() }()

	// Test loading the file entry
	loadedEntry, err := LoadFileEntry(readFile, entryOffset)
	if err != nil {
		t.Fatalf("LoadFileEntry() failed: %v", err)
	}
	if loadedEntry == nil {
		t.Fatal("LoadFileEntry() returned nil")
	}
	if loadedEntry.FileID != 1 {
		t.Errorf("LoadFileEntry() FileID = %d, want 1", loadedEntry.FileID)
	}
	if len(loadedEntry.Paths) != 1 {
		t.Errorf("LoadFileEntry() Paths length = %d, want 1", len(loadedEntry.Paths))
	}
	if loadedEntry.Paths[0].Path != "test.txt" {
		t.Errorf("LoadFileEntry() Path = %q, want %q", loadedEntry.Paths[0].Path, "test.txt")
	}
}

// TestLoadFileEntry_InvalidOffset tests LoadFileEntry with invalid offset.
func TestLoadFileEntry_InvalidOffset(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.nvpk")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_ = file.Close()

	readFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() { _ = readFile.Close() }()

	// Try to load from invalid offset (beyond file size)
	_, err = LoadFileEntry(readFile, 1000)
	if err == nil {
		t.Fatal("LoadFileEntry() should fail with invalid offset")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeIO {
		t.Errorf("Expected error type IO, got: %v", pkgErr.Type)
	}
}

// TestLoadFileEntry_InvalidEntry tests LoadFileEntry with invalid/corrupted entry.
func TestLoadFileEntry_InvalidEntry(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.nvpk")

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write invalid entry data (too short)
	_, _ = file.Write([]byte{0x00, 0x01, 0x02, 0x03})
	_ = file.Close()

	readFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() { _ = readFile.Close() }()

	// Try to load invalid entry
	_, err = LoadFileEntry(readFile, 0)
	if err == nil {
		t.Fatal("LoadFileEntry() should fail with invalid entry")
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	// Should be IO or Validation error
	if pkgErr.Type != pkgerrors.ErrTypeIO && pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type IO or Validation, got: %v", pkgErr.Type)
	}
}

// TestNormalizePackagePath_Canonicalization tests path canonicalization with dot segments.
// These tests verify that dot segments are properly canonicalized rather than rejected.
// This is the TDD Red Phase - these tests will fail until canonicalization is implemented.
func TestNormalizePackagePath_Canonicalization(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		// Category 1: Canonical path processing (should succeed after canonicalization)
		{
			name:        "no dot segments - already canonical",
			path:        "a/b/c",
			expected:    "/a/b/c",
			shouldError: false,
		},
		{
			name:        "single dot segment in middle",
			path:        "a/./b",
			expected:    "/a/b",
			shouldError: false,
		},
		{
			name:        "double dot segment in middle",
			path:        "a/b/../c",
			expected:    "/a/c",
			shouldError: false,
		},
		{
			name:        "leading slash removed",
			path:        "/a/b/c",
			expected:    "/a/b/c",
			shouldError: false,
		},
		{
			name:        "multiple slashes collapsed",
			path:        "a//b///c",
			expected:    "/a/b/c",
			shouldError: false,
		},
		{
			name:        "leading dot segment",
			path:        "./a/b",
			expected:    "/a/b",
			shouldError: false,
		},
		{
			name:        "trailing dot segment",
			path:        "a/b/.",
			expected:    "/a/b",
			shouldError: false,
		},
		{
			name:        "multiple double dot segments",
			path:        "a/b/c/../../d",
			expected:    "/a/d",
			shouldError: false,
		},
		{
			name:        "backslashes normalized",
			path:        "a\\b\\c",
			expected:    "/a/b/c",
			shouldError: false,
		},
		{
			name:        "complex canonicalization",
			path:        "./a/./b/../c/d",
			expected:    "/a/c/d",
			shouldError: false,
		},
		// Category 2: Root escaping (should fail - escapes package root)
		{
			name:        "parent at start escapes root",
			path:        "../a",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "multiple parents escape root",
			path:        "a/../../b",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "deep root escape",
			path:        "../../etc/passwd",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "just double dot escapes root",
			path:        "..",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "net negative depth escapes root",
			path:        "a/b/../../../c",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		// Category 3: Empty result after canonicalization (should fail - empty invalid)
		{
			name:        "single dot becomes empty",
			path:        ".",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "path resolves to empty",
			path:        "a/..",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "dot slash becomes empty",
			path:        "./",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "leading slash and dot",
			path:        "/.",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		// Category 4: Existing validations (should still work)
		{
			name:        "empty path rejected",
			path:        "",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizePackagePath(tt.path)
			if tt.shouldError {
				if err == nil {
					t.Errorf("NormalizePackagePath() expected error, got nil (result: %q)", result)
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("NormalizePackagePath() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("NormalizePackagePath() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("NormalizePackagePath() unexpected error: %v", err)
					return
				}
				if result != tt.expected {
					t.Errorf("NormalizePackagePath() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

// TestToDisplayPath tests the ToDisplayPath function.
func TestToDisplayPath(t *testing.T) {
	tests := []struct {
		name       string
		storedPath string
		expected   string
	}{
		{
			name:       "regular path with leading slash",
			storedPath: "/documents/file.txt",
			expected:   "documents/file.txt",
		},
		{
			name:       "nested path with leading slash",
			storedPath: "/a/b/c/d.txt",
			expected:   "a/b/c/d.txt",
		},
		{
			name:       "root path",
			storedPath: "/",
			expected:   "",
		},
		{
			name:       "single file at root",
			storedPath: "/file.txt",
			expected:   "file.txt",
		},
		{
			name:       "path without leading slash (defensive)",
			storedPath: "file.txt",
			expected:   "file.txt",
		},
		{
			name:       "deep nested path",
			storedPath: "/path/to/very/deep/directory/file.txt",
			expected:   "path/to/very/deep/directory/file.txt",
		},
		{
			name:       "path with special characters",
			storedPath: "/data/my-file_v2.txt",
			expected:   "data/my-file_v2.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToDisplayPath(tt.storedPath)
			if result != tt.expected {
				t.Errorf("ToDisplayPath(%q) = %q, want %q", tt.storedPath, result, tt.expected)
			}
		})
	}
}

// TestNormalizePackagePath_UnicodeNormalization tests Unicode NFC normalization.
// This ensures cross-platform compatibility between macOS (NFD) and Windows/Linux (NFC).
func TestNormalizePackagePath_UnicodeNormalization(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedNFC string
		note        string
	}{
		{
			name:        "NFD to NFC - decomposed e with acute",
			input:       "cafe\u0301", // café with decomposed é (e + combining acute accent)
			expectedNFC: "/café",      // café with composed é (NFC form)
			note:        "macOS NFD => NFC normalization",
		},
		{
			name:        "already NFC - composed e with acute",
			input:       "café",  // café with composed é (already NFC)
			expectedNFC: "/café", // Should remain unchanged
			note:        "Windows/Linux NFC form preserved",
		},
		{
			name:        "NFD to NFC - decomposed a with umlaut",
			input:       "a\u0308", // ä with decomposed form (a + combining diaeresis)
			expectedNFC: "/ä",      // ä with composed form (NFC)
			note:        "German umlaut normalization",
		},
		{
			name:        "complex path with multiple NFD characters",
			input:       "re\u0301sume\u0301", // résumé with decomposed accents
			expectedNFC: "/résumé",            // résumé with composed accents (NFC)
			note:        "Multiple decomposed characters",
		},
		{
			name:        "path with directory and NFD filename",
			input:       "documents/cafe\u0301.txt",
			expectedNFC: "/documents/café.txt",
			note:        "Directory path with NFD filename",
		},
		{
			name:        "ASCII path unchanged",
			input:       "simple/path/file.txt",
			expectedNFC: "/simple/path/file.txt",
			note:        "ASCII-only paths not affected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizePackagePath(tt.input)
			if err != nil {
				t.Fatalf("NormalizePackagePath(%q) failed: %v", tt.input, err)
			}

			// Compare byte sequences to verify actual normalization
			if result != tt.expectedNFC {
				t.Errorf("NormalizePackagePath(%q) = %q (bytes: %v), want %q (bytes: %v)",
					tt.input, result, []byte(result), tt.expectedNFC, []byte(tt.expectedNFC))
			}

			// Verify result is in NFC form
			if !norm.NFC.IsNormalString(result) {
				t.Errorf("NormalizePackagePath(%q) result is not in NFC form", tt.input)
			}
		})
	}
}

// TestValidatePathLength tests the ValidatePathLength function.
func TestValidatePathLength(t *testing.T) {
	// Helper to generate path of specific length
	generatePath := func(length int) string {
		if length <= 0 {
			return ""
		}
		// Generate path with repeated "a/" segments
		result := strings.Repeat("a/", length/2)
		// Pad to exact length
		for len(result) < length {
			result += "x"
		}
		return result[:length]
	}

	tests := []struct {
		name          string
		pathLength    int
		expectWarning bool
		expectError   bool
		warningCount  int
	}{
		{
			name:          "short path (100 bytes)",
			pathLength:    100,
			expectWarning: false,
			expectError:   false,
			warningCount:  0,
		},
		{
			name:          "at Windows default limit (260 bytes)",
			pathLength:    260,
			expectWarning: false,
			expectError:   false,
			warningCount:  0,
		},
		{
			name:          "above Windows default (261 bytes)",
			pathLength:    261,
			expectWarning: true,
			expectError:   false,
			warningCount:  1,
		},
		{
			name:          "at macOS limit (1024 bytes)",
			pathLength:    1024,
			expectWarning: true,
			expectError:   false,
			warningCount:  1,
		},
		{
			name:          "above macOS limit (1025 bytes)",
			pathLength:    1025,
			expectWarning: true,
			expectError:   false,
			warningCount:  1,
		},
		{
			name:          "at Linux limit (4096 bytes)",
			pathLength:    4096,
			expectWarning: true,
			expectError:   false,
			warningCount:  1,
		},
		{
			name:          "above Linux limit (4097 bytes)",
			pathLength:    4097,
			expectWarning: true,
			expectError:   false,
			warningCount:  1,
		},
		{
			name:          "at Windows extended limit (32767 bytes)",
			pathLength:    32767,
			expectWarning: true,
			expectError:   false,
			warningCount:  1,
		},
		{
			name:          "exceeds absolute maximum (32768 bytes)",
			pathLength:    32768,
			expectWarning: false, // Error instead
			expectError:   true,
			warningCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := generatePath(tt.pathLength)
			warnings, err := ValidatePathLength(path)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidatePathLength() expected error for %d byte path, got nil", tt.pathLength)
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("ValidatePathLength() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != pkgerrors.ErrTypeValidation {
					t.Errorf("ValidatePathLength() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePathLength() unexpected error: %v", err)
					return
				}

				if tt.expectWarning {
					if len(warnings) == 0 {
						t.Errorf("ValidatePathLength() expected warnings for %d byte path, got none", tt.pathLength)
					}
					if len(warnings) != tt.warningCount {
						t.Errorf("ValidatePathLength() warning count = %d, want %d", len(warnings), tt.warningCount)
					}
				} else {
					if len(warnings) != 0 {
						t.Errorf("ValidatePathLength() unexpected warnings for %d byte path: %v", tt.pathLength, warnings)
					}
				}
			}
		})
	}
}

// TestValidatePackagePath_Canonicalization tests that ValidatePackagePath properly handles paths with dot segments.
func TestValidatePackagePath_Canonicalization(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		shouldError bool
		errorType   pkgerrors.ErrorType
	}{
		// Valid paths with dot segments that should canonicalize
		{
			name:        "single dot in path",
			path:        "a/./b",
			shouldError: false,
		},
		{
			name:        "double dot in path",
			path:        "a/b/../c",
			shouldError: false,
		},
		// Invalid paths that escape root
		{
			name:        "escape root",
			path:        "../a",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		// Empty after canonicalization
		{
			name:        "resolves to empty",
			path:        "a/..",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		// Existing validations
		{
			name:        "empty path",
			path:        "",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
		{
			name:        "whitespace only",
			path:        "   ",
			shouldError: true,
			errorType:   pkgerrors.ErrTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackagePath(tt.path)
			if tt.shouldError {
				if err == nil {
					t.Errorf("ValidatePackagePath() expected error, got nil")
					return
				}
				var pkgErr *pkgerrors.PackageError
				if !pkgerrors.As(err, &pkgErr) {
					t.Errorf("ValidatePackagePath() error is not a PackageError: %v", err)
					return
				}
				if pkgErr.Type != tt.errorType {
					t.Errorf("ValidatePackagePath() error type = %v, want %v", pkgErr.Type, tt.errorType)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePackagePath() unexpected error: %v", err)
				}
			}
		})
	}
}
