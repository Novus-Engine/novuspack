package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	novuspack "github.com/novus-engine/novuspack/api/go"
)

func TestRunRead_PackageNotFound(t *testing.T) {
	err := runRead(readCmd, []string{"/nonexistent/pkg.nvpk", "/path"})
	if err == nil {
		t.Error("runRead on missing package should fail")
	}
}

func TestRunRead_FileNotInPackage(t *testing.T) {
	path := createTestPackage(t, "read.nvpk")
	if err := runRead(readCmd, []string{path, "/nonexistent.txt"}); err == nil {
		t.Error("runRead with path not in package should fail")
	}
}

func TestRunRead_OutputToDirFails(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "p.nvpk")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/x", []byte("x"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	_ = pkg.Close()
	readOutput = dir // output is a directory; WriteFile will fail
	defer func() { readOutput = "" }()
	err = runRead(readCmd, []string{pkgPath, "/x"})
	if err == nil {
		t.Error("runRead with output=dir should fail")
	}
	if err != nil && !strings.Contains(err.Error(), "write output") {
		t.Errorf("runRead output=dir: want 'write output' error, got %v", err)
	}
}

func TestRunRead_OutputToFile(t *testing.T) {
	readOutput = "/tmp/nvpkg-read-test-out"
	defer func() { readOutput = "" }()
	dir := t.TempDir()
	path := filepath.Join(dir, "read.nvpk")
	if err := runCreate(createCmd, []string{path}); err != nil {
		t.Fatalf("create: %v", err)
	}
	err := runRead(readCmd, []string{path, "/nonexistent.txt"})
	if err == nil {
		t.Error("runRead with missing file should fail")
	}
	_ = os.Remove(readOutput)
}

// TestRunRead_Success uses the api to create a package with a file; when api path
// metadata write is complete this test will cover runRead success (stdout and -o).
func TestRunRead_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "withfile.nvpk")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, path); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/x.txt", []byte("content"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	readOutput = ""
	err = runRead(readCmd, []string{path, "/x.txt"})
	if err != nil {
		t.Errorf("runRead: %v", err)
	}
}

func TestRunRead_Success_OutputToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "withfile2.nvpk")
	outPath := filepath.Join(dir, "out.txt")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, path); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/y.txt", []byte("data"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	readOutput = outPath
	defer func() { readOutput = "" }()
	err = runRead(readCmd, []string{path, "/y.txt"})
	if err != nil {
		t.Errorf("runRead: %v", err)
	}
	got, _ := os.ReadFile(outPath)
	if string(got) != "data" {
		t.Errorf("output file: got %q, want %q", string(got), "data")
	}
}
