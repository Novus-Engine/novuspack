package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreate_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pkg.nvpk")
	err := runCreate(createCmd, []string{path})
	if err != nil {
		t.Fatalf("runCreate: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("package file was not created")
	}
}

func TestRunCreate_WithOptions(t *testing.T) {
	createComment = "test comment"
	createVendorID = 1
	createAppID = 100
	defer func() {
		createComment = ""
		createVendorID = 0
		createAppID = 0
	}()
	dir := t.TempDir()
	path := filepath.Join(dir, "opts.nvpk")
	err := runCreate(createCmd, []string{path})
	if err != nil {
		t.Fatalf("runCreate with options: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("package file was not created")
	}
}

func TestRunCreate_InvalidPath(t *testing.T) {
	err := runCreate(createCmd, []string{""})
	if err == nil {
		t.Error("runCreate with empty path should fail")
	}
}

func TestRunCreate_WithOptionsInvalidPath(t *testing.T) {
	createComment = "x"
	defer func() { createComment = "" }()
	err := runCreate(createCmd, []string{""})
	if err == nil {
		t.Error("runCreate with options and empty path should fail")
	}
	if err != nil && !strings.Contains(err.Error(), "create with options") && !strings.Contains(err.Error(), "create") {
		t.Errorf("runCreate options invalid path: want create error, got %v", err)
	}
}

// TestRunCreate_ThenValidate creates a package then validates it via CLI (full write round-trip).
func TestRunCreate_ThenValidate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pkg.nvpk")
	if err := runCreate(createCmd, []string{path}); err != nil {
		t.Fatalf("runCreate: %v", err)
	}
	if err := runValidate(validateCmd, []string{path}); err != nil {
		t.Fatalf("runValidate after create: %v", err)
	}
}

func TestRunCreate_WithMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mode.nvpk")
	createComment = ""
	createVendorID = 0
	createAppID = 0
	createModeStr = "0644"
	defer func() { createModeStr = "" }()
	if err := runCreate(createCmd, []string{path}); err != nil {
		t.Fatalf("runCreate with mode: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("package file was not created")
	}
}

func TestParseFileMode(t *testing.T) {
	// 0644 octal == 420 decimal; "420" is valid octal (0o420 == 272) so code parses as octal first
	mode0644 := os.FileMode(0o644)
	tests := []struct {
		s    string
		want os.FileMode
		ok   bool
	}{
		{"0644", mode0644, true},
		{"420", os.FileMode(0o420), true}, // "420" parsed as octal = 272
		{"499", os.FileMode(499), true},   // "499" invalid octal, parsed as decimal
		{"", 0, false},
		{"invalid", 0, false},
	}
	for _, tt := range tests {
		got, err := parseFileMode(tt.s)
		if tt.ok && err != nil {
			t.Errorf("parseFileMode(%q): %v", tt.s, err)
			continue
		}
		if !tt.ok && err == nil {
			t.Errorf("parseFileMode(%q): expected error", tt.s)
			continue
		}
		if tt.ok && got != tt.want {
			t.Errorf("parseFileMode(%q) = %#o, want %#o", tt.s, got, tt.want)
		}
	}
}
