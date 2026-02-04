package cmd

import (
	"path/filepath"
	"testing"
)

func runIdentityWithFlags(t *testing.T, vendorID uint32, appID uint64, path string) error {
	t.Helper()
	identityVendorID = vendorID
	identityAppID = appID
	defer func() { identityVendorID = 0; identityAppID = 0 }()
	return runIdentity(identityCmd, []string{path})
}

func TestRunIdentity(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		if err := runIdentityWithFlags(t, 0, 0, createTestPackage(t, "id1.nvpk")); err != nil {
			t.Fatalf("runIdentity get: %v", err)
		}
	})
	t.Run("set", func(t *testing.T) {
		if err := runIdentityWithFlags(t, 1, 100, createTestPackage(t, "id2.nvpk")); err != nil {
			t.Fatalf("runIdentity set: %v", err)
		}
	})
	t.Run("no_such_package", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "nonexistent.nvpk")
		err := runIdentityWithFlags(t, 0, 0, path)
		if err == nil {
			t.Fatal("runIdentity on nonexistent package should fail")
		}
	})
}
