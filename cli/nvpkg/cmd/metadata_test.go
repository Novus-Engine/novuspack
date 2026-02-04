package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func runMetadataWithFlags(t *testing.T, useJSON, readOnly bool, path string) error {
	t.Helper()
	metadataJSON = useJSON
	metadataReadOnly = readOnly
	defer func() { metadataJSON = false; metadataReadOnly = false }()
	return runMetadata(metadataCmd, []string{path})
}

func TestRunMetadata(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if err := runMetadataWithFlags(t, false, false, createTestPackage(t, "meta1.nvpk")); err != nil {
			t.Fatalf("runMetadata: %v", err)
		}
	})
	t.Run("json", func(t *testing.T) {
		if err := runMetadataWithFlags(t, true, false, createTestPackage(t, "meta2.nvpk")); err != nil {
			t.Fatalf("runMetadata --json: %v", err)
		}
	})
	t.Run("read_only", func(t *testing.T) {
		if err := runMetadataWithFlags(t, false, true, createTestPackage(t, "meta3.nvpk")); err != nil {
			t.Fatalf("runMetadata --read-only: %v", err)
		}
	})
	t.Run("no_such_package", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "nonexistent.nvpk")
		err := runMetadataWithFlags(t, false, false, path)
		if err == nil {
			t.Fatal("runMetadata on nonexistent package should fail")
		}
		if !strings.Contains(err.Error(), "open") {
			t.Errorf("runMetadata no package: want error containing 'open', got %v", err)
		}
	})
}
