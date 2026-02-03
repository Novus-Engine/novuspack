package novus_package

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestSetTargetPath_Success tests successful target path setting.
func TestSetTargetPath_Success(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Test setting target path to valid location
	targetPath := filepath.Join(tmpDir, "output.nvpk")
	err = pkg.SetTargetPath(ctx, targetPath)
	if err != nil {
		t.Errorf("SetTargetPath(%q) failed: %v", targetPath, err)
	}

	// Verify the path was set
	if pkg.GetPath() != targetPath {
		t.Errorf("GetPath() = %q, want %q", pkg.GetPath(), targetPath)
	}
}

// TestSetTargetPath_ErrorCases tests error conditions for SetTargetPath.
//
//nolint:gocognit // table-driven error cases
func TestSetTargetPath_ErrorCases(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupPath   func(t *testing.T) string
		wantErrType pkgerrors.ErrorType
		wantErrMsg  string
	}{
		{
			name: "Empty path",
			setupPath: func(t *testing.T) string {
				return ""
			},
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "target path cannot be empty",
		},
		{
			name: "Non-existent directory",
			setupPath: func(t *testing.T) string {
				return "/nonexistent/directory/file.nvpk"
			},
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "target directory does not exist",
		},
		{
			name: "Non-writable directory",
			setupPath: func(t *testing.T) string {
				// Create a read-only directory
				tmpDir := t.TempDir()
				readOnlyDir := filepath.Join(tmpDir, "readonly")
				if err := os.Mkdir(readOnlyDir, 0o444); err != nil {
					t.Fatalf("Failed to create read-only directory: %v", err)
				}
				return filepath.Join(readOnlyDir, "file.nvpk")
			},
			wantErrType: pkgerrors.ErrTypeSecurity,
			wantErrMsg:  "target directory is not writable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := NewPackage()
			if err != nil {
				t.Fatalf("NewPackage failed: %v", err)
			}

			targetPath := tt.setupPath(t)
			err = pkg.SetTargetPath(ctx, targetPath)
			if err == nil {
				t.Errorf("SetTargetPath(%q) succeeded, want error", targetPath)
				return
			}

			// Check error type
			pkgErr, ok := err.(*pkgerrors.PackageError)
			if !ok {
				t.Errorf("SetTargetPath(%q) error type = %T, want *pkgerrors.PackageError", targetPath, err)
				return
			}

			if pkgErr.Type != tt.wantErrType {
				t.Errorf("SetTargetPath(%q) error type = %v, want %v", targetPath, pkgErr.Type, tt.wantErrType)
			}

			if pkgErr.Message != tt.wantErrMsg {
				t.Errorf("SetTargetPath(%q) error message = %q, want %q", targetPath, pkgErr.Message, tt.wantErrMsg)
			}
		})
	}
}

// TestSetTargetPath_ContextCancellation tests context cancellation handling.
func TestSetTargetPath_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "output.nvpk")

	err = pkg.SetTargetPath(ctx, targetPath)
	if err == nil {
		t.Error("SetTargetPath with cancelled context should fail")
		return
	}

	pkgErr, ok := err.(*pkgerrors.PackageError)
	if !ok {
		t.Errorf("Error type = %T, want *pkgerrors.PackageError", err)
		return
	}

	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeContext)
	}
}

// TestSetTargetPath_UpdatesPath tests that SetTargetPath correctly updates the internal path.
func TestSetTargetPath_UpdatesPath(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tmpDir := t.TempDir()

	// Set first path
	path1 := filepath.Join(tmpDir, "output1.nvpk")
	if err := pkg.SetTargetPath(ctx, path1); err != nil {
		t.Fatalf("SetTargetPath(%q) failed: %v", path1, err)
	}

	if pkg.GetPath() != path1 {
		t.Errorf("GetPath() = %q, want %q", pkg.GetPath(), path1)
	}

	// Set second path
	path2 := filepath.Join(tmpDir, "output2.nvpk")
	if err := pkg.SetTargetPath(ctx, path2); err != nil {
		t.Fatalf("SetTargetPath(%q) failed: %v", path2, err)
	}

	if pkg.GetPath() != path2 {
		t.Errorf("GetPath() = %q, want %q after second SetTargetPath", pkg.GetPath(), path2)
	}
}

// TestSetTargetPath_PathCleaning tests that paths are cleaned/normalized.
func TestSetTargetPath_PathCleaning(t *testing.T) {
	ctx := context.Background()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tmpDir := t.TempDir()

	// Test path with extra slashes and dots
	dirtyPath := filepath.Join(tmpDir, "subdir", "..", "output.nvpk")
	expectedClean := filepath.Join(tmpDir, "output.nvpk")

	// Create subdir so the parent validation works
	if err := os.Mkdir(filepath.Join(tmpDir, "subdir"), 0o755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	if err := pkg.SetTargetPath(ctx, dirtyPath); err != nil {
		t.Fatalf("SetTargetPath(%q) failed: %v", dirtyPath, err)
	}

	got := pkg.GetPath()
	if got != expectedClean {
		t.Errorf("GetPath() = %q, want %q (cleaned path)", got, expectedClean)
	}
}
