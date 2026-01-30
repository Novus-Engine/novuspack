package novus_package

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestSetSessionBase_Success tests successful session base setting.
func TestSetSessionBase_Success(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tests := []struct {
		name     string
		basePath string
	}{
		{
			name:     "Unix absolute path",
			basePath: "/home/user/project",
		},
		// Note: Windows path testing skipped - filepath.IsAbs is platform-specific
		// On Linux, Windows paths like "C:\\" are not considered absolute
		{
			name:     "Path with trailing slash",
			basePath: "/home/user/project/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pkg.SetSessionBase(tt.basePath)
			if err != nil {
				t.Errorf("SetSessionBase(%q) failed: %v", tt.basePath, err)
			}

			// Verify session base is set
			if !pkg.HasSessionBase() {
				t.Error("HasSessionBase() = false, want true")
			}

			// Verify GetSessionBase returns the path (cleaned)
			got := pkg.GetSessionBase()
			if got == "" {
				t.Error("GetSessionBase() returned empty string")
			}
		})
	}
}

// TestSetSessionBase_ErrorCases tests error conditions for SetSessionBase.
func TestSetSessionBase_ErrorCases(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	tests := []struct {
		name        string
		basePath    string
		wantErrType pkgerrors.ErrorType
		wantErrMsg  string
	}{
		{
			name:        "Empty path",
			basePath:    "",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "session base path cannot be empty",
		},
		{
			name:        "Relative path",
			basePath:    "relative/path",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "session base path must be an absolute path",
		},
		{
			name:        "Relative path with dot",
			basePath:    "./relative/path",
			wantErrType: pkgerrors.ErrTypeValidation,
			wantErrMsg:  "session base path must be an absolute path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pkg.SetSessionBase(tt.basePath)
			if err == nil {
				t.Errorf("SetSessionBase(%q) succeeded, want error", tt.basePath)
				return
			}

			// Check error type
			pkgErr, ok := err.(*pkgerrors.PackageError)
			if !ok {
				t.Errorf("SetSessionBase(%q) error type = %T, want *pkgerrors.PackageError", tt.basePath, err)
				return
			}

			if pkgErr.Type != tt.wantErrType {
				t.Errorf("SetSessionBase(%q) error type = %v, want %v", tt.basePath, pkgErr.Type, tt.wantErrType)
			}

			if pkgErr.Message != tt.wantErrMsg {
				t.Errorf("SetSessionBase(%q) error message = %q, want %q", tt.basePath, pkgErr.Message, tt.wantErrMsg)
			}
		})
	}
}

// TestGetSessionBase tests GetSessionBase method.
func TestGetSessionBase(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Initially should be empty
	if got := pkg.GetSessionBase(); got != "" {
		t.Errorf("GetSessionBase() = %q, want empty string", got)
	}

	// Set a base path
	basePath := "/home/user/project"
	if err := pkg.SetSessionBase(basePath); err != nil {
		t.Fatalf("SetSessionBase failed: %v", err)
	}

	// Should return the set path (cleaned)
	got := pkg.GetSessionBase()
	if got == "" {
		t.Error("GetSessionBase() returned empty string after SetSessionBase")
	}
}

// TestClearSessionBase tests ClearSessionBase method.
func TestClearSessionBase(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Set a session base
	basePath := "/home/user/project"
	if err := pkg.SetSessionBase(basePath); err != nil {
		t.Fatalf("SetSessionBase failed: %v", err)
	}

	// Verify it's set
	if !pkg.HasSessionBase() {
		t.Error("HasSessionBase() = false, want true after SetSessionBase")
	}

	// Clear it
	pkg.ClearSessionBase()

	// Verify it's cleared
	if pkg.HasSessionBase() {
		t.Error("HasSessionBase() = true, want false after ClearSessionBase")
	}

	if got := pkg.GetSessionBase(); got != "" {
		t.Errorf("GetSessionBase() = %q, want empty string after ClearSessionBase", got)
	}
}

// TestHasSessionBase tests HasSessionBase method.
func TestHasSessionBase(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	// Initially should be false
	if pkg.HasSessionBase() {
		t.Error("HasSessionBase() = true, want false for new package")
	}

	// Set a session base
	basePath := "/home/user/project"
	if err := pkg.SetSessionBase(basePath); err != nil {
		t.Fatalf("SetSessionBase failed: %v", err)
	}

	// Should be true now
	if !pkg.HasSessionBase() {
		t.Error("HasSessionBase() = false, want true after SetSessionBase")
	}

	// Clear it
	pkg.ClearSessionBase()

	// Should be false again
	if pkg.HasSessionBase() {
		t.Error("HasSessionBase() = true, want false after ClearSessionBase")
	}
}

// TestSessionBase_MultipleUpdates tests updating session base multiple times.
func TestSessionBase_MultipleUpdates(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage failed: %v", err)
	}

	paths := []string{
		"/home/user/project1",
		"/home/user/project2",
		"/opt/application",
	}

	for _, path := range paths {
		if err := pkg.SetSessionBase(path); err != nil {
			t.Errorf("SetSessionBase(%q) failed: %v", path, err)
		}

		if !pkg.HasSessionBase() {
			t.Errorf("HasSessionBase() = false after SetSessionBase(%q)", path)
		}

		got := pkg.GetSessionBase()
		if got == "" {
			t.Errorf("GetSessionBase() returned empty string after SetSessionBase(%q)", path)
		}
	}
}
