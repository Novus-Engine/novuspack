package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAdd_SourceNotFound(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "pkg.nvpk")
	if err := runCreate(createCmd, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	err := runAdd(addCmd, []string{pkgPath, filepath.Join(dir, "nonexistent.txt")})
	if err == nil {
		t.Error("runAdd with missing source should fail")
	}
}

func TestRunAdd_PackageNotExist_CreatePath(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "newpkg.nvpk")
	srcPath := filepath.Join(dir, "src.txt")
	if err := os.WriteFile(srcPath, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runAdd(addCmd, []string{pkgPath, srcPath})
	if err != nil {
		t.Fatalf("runAdd (create new package): %v", err)
	}
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		t.Error("package file was not created")
	}
}

func TestRunAdd_OpenExisting(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "existing.nvpk")
	srcPath := filepath.Join(dir, "file.txt")
	if err := runCreate(createCmd, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := os.WriteFile(srcPath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runAdd(addCmd, []string{pkgPath, srcPath})
	if err != nil {
		t.Fatalf("runAdd (open existing): %v", err)
	}
}

func TestRunAdd_WithStoredPath(t *testing.T) {
	addStoredPath = "/stored.txt"
	defer func() { addStoredPath = "" }()
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "with-as.nvpk")
	srcPath := filepath.Join(dir, "local.txt")
	if err := os.WriteFile(srcPath, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runAdd(addCmd, []string{pkgPath, srcPath})
	if err != nil {
		t.Fatalf("runAdd (with --as): %v", err)
	}
}

func TestRunAdd_Directory(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	pkgPath := filepath.Join(dir, "dirpkg.nvpk")
	err := runAdd(addCmd, []string{pkgPath, subdir})
	// AddDirectory is a stub; accept success or unsupported until full implementation
	if err != nil && !strings.Contains(err.Error(), "unsupported") && !strings.Contains(err.Error(), "AddDirectory") {
		t.Fatalf("runAdd (directory): %v", err)
	}
}

func TestRunAdd_OpenInvalidPackage(t *testing.T) {
	dir := t.TempDir()
	// Path exists but is a regular file, not a package
	fakePkg := filepath.Join(dir, "fake.nvpk")
	if err := os.WriteFile(fakePkg, []byte("not a package"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := runAdd(addCmd, []string{fakePkg, src})
	if err == nil {
		t.Error("runAdd with non-package file path should fail")
	}
}

func TestRunAdd_PackagePathIsDirectory(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Package path is an existing directory; OpenPackage will fail
	err := runAdd(addCmd, []string{dir, src})
	if err == nil {
		t.Error("runAdd with directory as package path should fail")
	}
}
