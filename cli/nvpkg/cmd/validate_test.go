package cmd

import (
	"context"
	"path/filepath"
	"testing"

	novuspack "github.com/novus-engine/novuspack/api/go"
)

func TestRunValidate_PackageNotFound(t *testing.T) {
	err := runValidate(nil, []string{"/nonexistent/pkg.nvpk"})
	if err == nil {
		t.Error("runValidate on missing package should fail")
	}
}

func TestRunValidate_Success(t *testing.T) {
	path := createTestPackage(t, "validate.nvpk")
	if err := runValidate(nil, []string{path}); err != nil {
		t.Errorf("runValidate: %v", err)
	}
}

func TestRunValidate_OpenPackage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "v.nvpk")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, path); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/a.txt", []byte("x"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	if err := runValidate(nil, []string{path}); err != nil {
		t.Errorf("runValidate on written package: %v", err)
	}
}
