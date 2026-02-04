package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRunComment_GetNoComment(t *testing.T) {
	path := createTestPackage(t, "c1.nvpk")
	commentSet = ""
	commentClear = false
	defer func() { commentSet = ""; commentClear = false }()
	err := runComment(commentCmd, []string{path})
	if err != nil {
		t.Fatalf("runComment get: %v", err)
	}
}

func TestRunComment_GetWithComment(t *testing.T) {
	path := createTestPackage(t, "c2.nvpk")
	commentSet = "hello"
	commentClear = false
	defer func() { commentSet = ""; commentClear = false }()
	if err := runComment(commentCmd, []string{path}); err != nil {
		t.Fatalf("runComment set: %v", err)
	}
	commentSet = ""
	if err := runComment(commentCmd, []string{path}); err != nil {
		t.Fatalf("runComment get after set: %v", err)
	}
}

func TestRunComment_Set(t *testing.T) {
	path := createTestPackage(t, "c3.nvpk")
	commentSet = "test comment"
	commentClear = false
	defer func() { commentSet = ""; commentClear = false }()
	if err := runComment(commentCmd, []string{path}); err != nil {
		t.Fatalf("runComment set: %v", err)
	}
}

func TestRunComment_Clear(t *testing.T) {
	path := createTestPackage(t, "c4.nvpk")
	commentSet = "x"
	commentClear = false
	defer func() { commentSet = ""; commentClear = false }()
	if err := runComment(commentCmd, []string{path}); err != nil {
		t.Fatalf("runComment set: %v", err)
	}
	commentSet = ""
	commentClear = true
	defer func() { commentClear = false }()
	if err := runComment(commentCmd, []string{path}); err != nil {
		t.Fatalf("runComment clear: %v", err)
	}
}

func TestRunComment_BothSetAndClear(t *testing.T) {
	path := createTestPackage(t, "c5.nvpk")
	commentSet = "x"
	commentClear = true
	defer func() { commentSet = ""; commentClear = false }()
	err := runComment(commentCmd, []string{path})
	if err == nil {
		t.Fatal("runComment with both --set and --clear should fail")
	}
	if !strings.Contains(err.Error(), "both") {
		t.Errorf("runComment both: want error containing 'both', got %v", err)
	}
}

func TestRunComment_NoSuchPackage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.nvpk")
	commentSet = ""
	commentClear = false
	err := runComment(commentCmd, []string{path})
	if err == nil {
		t.Fatal("runComment on nonexistent package should fail")
	}
}
