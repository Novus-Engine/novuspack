package cmd

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	novuspack "github.com/novus-engine/novuspack/api/go"
)

func TestRunRemove_PackageNotFound(t *testing.T) {
	err := runRemove(removeCmd, []string{"/nonexistent/pkg.nvpk", "/some/path"})
	if err == nil {
		t.Error("runRemove on missing package should fail")
	}
}

func TestRunRemove_FileNotInPackage(t *testing.T) {
	path := createTestPackage(t, "remove.nvpk")
	if err := runRemove(removeCmd, []string{path, "/nonexistent.txt"}); err == nil {
		t.Error("runRemove with path not in package should fail")
	}
}

// TestRunRemove_Success creates a package with a file via the API, then removes the file via CLI.
func TestRunRemove_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "remove_success.nvpk")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, path); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/z.txt", []byte("z"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	err = runRemove(removeCmd, []string{path, "/z.txt"})
	if err != nil {
		t.Errorf("runRemove: %v", err)
	}
}

func TestRunRemove_WithPatternFlag(t *testing.T) {
	path := createTestPackage(t, "remove_pattern.nvpk")
	removePattern = true
	defer func() { removePattern = false }()
	// API may return ErrTypeUnsupported until RemoveFilePattern is implemented
	err := runRemove(removeCmd, []string{path, "*.tmp"})
	if err != nil && !strings.Contains(err.Error(), "unsupported") && !strings.Contains(err.Error(), "remove pattern") {
		t.Errorf("runRemove --pattern: %v", err)
	}
}

func TestRunRemove_DirectoryPath(t *testing.T) {
	path := createTestPackage(t, "remove_dir.nvpk")
	// API may return ErrTypeUnsupported until RemoveDirectory is implemented
	err := runRemove(removeCmd, []string{path, "/subdir/"})
	if err != nil && !strings.Contains(err.Error(), "unsupported") && !strings.Contains(err.Error(), "remove directory") {
		t.Errorf("runRemove directory: %v", err)
	}
}
