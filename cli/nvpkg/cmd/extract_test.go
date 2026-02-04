package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	novuspack "github.com/novus-engine/novuspack/api/go"
)

func TestMatchPathPrefix(t *testing.T) {
	tests := []struct {
		displayPath string
		pathPrefix  string
		want        bool
	}{
		{"a", "", true},
		{"a/b", "", true},
		{"a", "a", true},
		{"a/b", "a", true},
		{"a/b/c", "a/b", true},
		{"a", "b", false},
		{"a/b", "b", false},
		{"ab", "a", false},
	}
	for _, tt := range tests {
		got := matchPathPrefix(tt.displayPath, tt.pathPrefix)
		if got != tt.want {
			t.Errorf("matchPathPrefix(%q, %q) => %v, want %v", tt.displayPath, tt.pathPrefix, got, tt.want)
		}
	}
}

func TestResolveExtractDest(t *testing.T) {
	t.Run("empty_returns_error", func(t *testing.T) {
		_, err := resolveExtractDest("")
		if err == nil {
			t.Error("resolveExtractDest(\"\") => nil error")
		}
	})
	t.Run("valid_returns_abs", func(t *testing.T) {
		dir := t.TempDir()
		got, err := resolveExtractDest(dir)
		if err != nil {
			t.Fatalf("resolveExtractDest: %v", err)
		}
		if !filepath.IsAbs(got) {
			t.Errorf("resolveExtractDest => %q not absolute", got)
		}
	})
	t.Run("path_under_file_fails", func(t *testing.T) {
		dir := t.TempDir()
		blocker := filepath.Join(dir, "blocker")
		if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := resolveExtractDest(filepath.Join(blocker, "sub"))
		if err == nil {
			t.Error("resolveExtractDest under existing file should fail")
		}
	})
}

func TestRunExtract_NoOutputFlag(t *testing.T) {
	path := createTestPackage(t, "ext_noout.nvpk")
	extractOutput = ""
	defer func() { extractOutput = "" }()
	err := runExtract(nil, []string{path})
	if err == nil {
		t.Error("runExtract with empty output dir should fail")
	}
	if err != nil && !strings.Contains(err.Error(), "output") && !strings.Contains(err.Error(), "required") {
		t.Errorf("runExtract no output: want error about output/required, got %v", err)
	}
}

func TestRunExtract_BadPackagePath(t *testing.T) {
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	extractOutput = outDir
	defer func() { extractOutput = "" }()
	err := runExtract(nil, []string{filepath.Join(dir, "nonexistent.nvpk")})
	if err == nil {
		t.Error("runExtract with nonexistent package should fail")
	}
}

func TestRunExtract_WithPathPrefix(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "p.nvpk")
	outDir := filepath.Join(dir, "out")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	extractOutput = outDir
	defer func() { extractOutput = "" }()
	// Path prefix that matches nothing in empty package; should complete without error
	if err := runExtract(nil, []string{pkgPath, "/nonexistent"}); err != nil {
		t.Fatalf("runExtract: %v", err)
	}
}

func TestRunExtract_EmptyPackage(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "p.nvpk")
	outDir := filepath.Join(dir, "out")
	if err := runCreate(nil, []string{pkgPath}); err != nil {
		t.Fatalf("create: %v", err)
	}
	extractOutput = outDir
	defer func() { extractOutput = "" }()
	if err := runExtract(nil, []string{pkgPath}); err != nil {
		t.Fatalf("runExtract: %v", err)
	}
	// Empty package: no files extracted; out dir may be created empty
	if _, err := os.Stat(outDir); err != nil {
		t.Errorf("output dir not created: %v", err)
	}
}

// TestExtractOneFile covers extractOneFile using an in-memory package (no Write).
func TestExtractOneFile(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "mem.nvpk")
	outDir := filepath.Join(dir, "out")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/sub/b.txt", []byte("in-memory"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := extractOneFile(ctx, pkg, outDir, "sub/b.txt"); err != nil {
		t.Fatalf("extractOneFile: %v", err)
	}
	gotPath := filepath.Join(outDir, "sub", "b.txt")
	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(got) != "in-memory" {
		t.Errorf("extracted content: got %q, want %q", string(got), "in-memory")
	}
}

// TestRunExtract_OneFile covers runExtract when the package has one file on disk.
func TestRunExtract_OneFile(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "one.nvpk")
	outDir := filepath.Join(dir, "out")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/a.txt", []byte("extracted"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	if err := pkg.Write(ctx); err != nil {
		t.Skipf("Write failed (api path metadata may be incomplete): %v", err)
	}
	extractOutput = outDir
	defer func() { extractOutput = "" }()
	if err := runExtract(nil, []string{pkgPath}); err != nil {
		t.Fatalf("runExtract: %v", err)
	}
	gotPath := filepath.Join(outDir, "a.txt")
	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(got) != "extracted" {
		t.Errorf("extracted content: got %q, want %q", string(got), "extracted")
	}
}

func TestExtractOneFile_ReadError(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "empty.nvpk")
	outDir := filepath.Join(dir, "out")
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	// Package has no file at /missing.txt
	err = extractOneFile(ctx, pkg, outDir, "missing.txt")
	if err == nil {
		t.Error("extractOneFile with missing path should fail")
	}
}

func TestExtractOneFile_DestUnderFile(t *testing.T) {
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "p.nvpk")
	outDir := filepath.Join(dir, "blocker")
	if err := os.WriteFile(outDir, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	pkg, err := novuspack.NewPackage()
	if err != nil {
		t.Fatalf("NewPackage: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	if err := pkg.Create(ctx, pkgPath); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := pkg.AddFileFromMemory(ctx, "/a.txt", []byte("x"), nil); err != nil {
		t.Fatalf("AddFileFromMemory: %v", err)
	}
	// destDir is a file, so MkdirAll(filepath.Dir(destPath)) or WriteFile will fail
	err = extractOneFile(ctx, pkg, outDir, "a.txt")
	if err == nil {
		t.Error("extractOneFile with dest under file should fail")
	}
}
