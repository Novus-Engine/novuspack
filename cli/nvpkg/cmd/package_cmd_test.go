package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPackageCommands_SuccessAndNotFound(t *testing.T) {
	tests := []struct {
		name     string
		pkgName  string
		runWith  func(path string) error
		runEmpty func() error
	}{
		{
			name:     "header",
			pkgName:  "header.nvpk",
			runWith:  func(p string) error { return runHeader(headerCmd, []string{p}) },
			runEmpty: func() error { return runHeader(headerCmd, []string{"/nonexistent/pkg.nvpk"}) },
		},
		{
			name:     "info",
			pkgName:  "info.nvpk",
			runWith:  func(p string) error { return runInfo(infoCmd, []string{p}) },
			runEmpty: func() error { return runInfo(infoCmd, []string{"/nonexistent/pkg.nvpk"}) },
		},
		{
			name:     "list",
			pkgName:  "list.nvpk",
			runWith:  func(p string) error { return runList(listCmd, []string{p}) },
			runEmpty: func() error { return runList(listCmd, []string{"/nonexistent/pkg.nvpk"}) },
		},
		{
			name:     "validate",
			pkgName:  "validate.nvpk",
			runWith:  func(p string) error { return runValidate(nil, []string{p}) },
			runEmpty: func() error { return runValidate(nil, []string{"/nonexistent/pkg.nvpk"}) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/success", func(t *testing.T) {
			path := createTestPackage(t, tt.pkgName)
			if err := tt.runWith(path); err != nil {
				t.Fatalf("%s: %v", tt.name, err)
			}
		})
		t.Run(tt.name+"/notfound", func(t *testing.T) {
			if err := tt.runEmpty(); err == nil {
				t.Errorf("%s: expected failure on missing package", tt.name)
			}
		})
	}
}

func TestRunInfo_WithComment(t *testing.T) {
	createComment = "my comment"
	defer func() { createComment = "" }()
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "comment.nvpk")
	if err := runCreate(createCmd, []string{pkgPath}); err != nil {
		t.Fatalf("runCreate: %v", err)
	}
	if err := runInfo(infoCmd, []string{pkgPath}); err != nil {
		t.Fatalf("runInfo: %v", err)
	}
}

func TestRunInfo_WithVendorAndAppId(t *testing.T) {
	createVendorID = 1
	createAppID = 100
	defer func() { createVendorID = 0; createAppID = 0 }()
	dir := t.TempDir()
	pkgPath := filepath.Join(dir, "identity.nvpk")
	if err := runCreate(createCmd, []string{pkgPath}); err != nil {
		t.Fatalf("runCreate: %v", err)
	}
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	err = runInfo(infoCmd, []string{pkgPath})
	_ = w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("runInfo: %v", err)
	}
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "Vendor ID") || !strings.Contains(string(out), "App ID") {
		t.Errorf("info output should contain Vendor ID and App ID: %s", out)
	}
}
