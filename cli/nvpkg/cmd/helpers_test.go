package cmd

import (
	"path/filepath"
	"testing"
)

// createTestPackage creates an empty package at a temp path and returns the path.
func createTestPackage(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := runCreate(createCmd, []string{path}); err != nil {
		t.Fatalf("create: %v", err)
	}
	return path
}
