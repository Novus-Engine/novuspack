// This file contains tests for package comment round-trip (write and read back).
// It verifies that comments are properly serialized to disk and deserialized on read.

package novus_package

import (
	"context"
	"path/filepath"
	"testing"
)

// TestPackage_CommentRoundTrip_Basic tests writing and reading back a comment.
func TestPackage_CommentRoundTrip_Basic(t *testing.T) {
	ctx := context.Background()
	tempPath := filepath.Join(t.TempDir(), "test.nvpk")

	// Create package and set comment
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if err := pkg1.Create(ctx, tempPath); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	fpkg1 := pkg1.(*filePackage)
	expectedComment := "Test comment for round-trip"
	if err := fpkg1.SetComment(expectedComment); err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	// Write package (0 files is valid - just metadata)
	// Write and close
	if err := pkg1.Write(ctx); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}
	if err := pkg1.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Reopen and verify comment
	pkg2, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	actualComment := pkg2.GetComment()
	if actualComment != expectedComment {
		t.Errorf("GetComment() = %q, want %q", actualComment, expectedComment)
	}

	fpkg2 := pkg2.(*filePackage)
	if !fpkg2.HasComment() {
		t.Error("HasComment() = false, want true")
	}
	if !fpkg2.Info.HasComment {
		t.Error("Info.HasComment = false, want true")
	}
}

// TestPackage_CommentRoundTrip_LongComment tests round-trip with a longer comment.
func TestPackage_CommentRoundTrip_LongComment(t *testing.T) {
	ctx := context.Background()
	tempPath := filepath.Join(t.TempDir(), "test.nvpk")

	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if err := pkg1.Create(ctx, tempPath); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	fpkg1 := pkg1.(*filePackage)
	expectedComment := "This is a longer comment with multiple lines.\nLine 2 with special chars: !@#$%^&*()\nLine 3 with unicode: \u00e9\u00e8\u00ea"
	if err := fpkg1.SetComment(expectedComment); err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	// Write package (0 files is valid)
	if err := pkg1.Write(ctx); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}
	if err := pkg1.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	pkg2, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	actualComment := pkg2.GetComment()
	if actualComment != expectedComment {
		t.Errorf("GetComment() = %q, want %q", actualComment, expectedComment)
	}
}

// TestPackage_CommentRoundTrip_NoComment tests package without comment.
func TestPackage_CommentRoundTrip_NoComment(t *testing.T) {
	ctx := context.Background()
	tempPath := filepath.Join(t.TempDir(), "test.nvpk")

	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if err := pkg1.Create(ctx, tempPath); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Don't set comment, just write (0 files is valid)
	if err := pkg1.Write(ctx); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}
	if err := pkg1.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	pkg2, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg2.Close() }()

	if pkg2.HasComment() {
		t.Error("HasComment() = true, want false")
	}
	if pkg2.GetComment() != "" {
		t.Errorf("GetComment() = %q, want empty string", pkg2.GetComment())
	}
}

// TestPackage_CommentRoundTrip_ModifyComment tests modifying comment and writing again.
func TestPackage_CommentRoundTrip_ModifyComment(t *testing.T) {
	ctx := context.Background()
	tempPath := filepath.Join(t.TempDir(), "test.nvpk")

	// Create with initial comment
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if err := pkg1.Create(ctx, tempPath); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	fpkg1 := pkg1.(*filePackage)
	if err := fpkg1.SetComment("Initial comment"); err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}
	// Write package (0 files is valid)
	if err := pkg1.Write(ctx); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}
	if err := pkg1.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Reopen and modify comment
	pkg2, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	fpkg2 := pkg2.(*filePackage)
	expectedComment := "Modified comment"
	if err := fpkg2.SetComment(expectedComment); err != nil {
		t.Fatalf("SetComment() on reopened package failed: %v", err)
	}
	if err := pkg2.Write(ctx); err != nil {
		t.Fatalf("Write() after modification failed: %v", err)
	}
	if err := pkg2.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Reopen again and verify modified comment
	pkg3, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() second time failed: %v", err)
	}
	defer func() { _ = pkg3.Close() }()

	actualComment := pkg3.GetComment()
	if actualComment != expectedComment {
		t.Errorf("GetComment() after modification = %q, want %q", actualComment, expectedComment)
	}
}

// TestPackage_CommentRoundTrip_ClearComment tests clearing a comment.
func TestPackage_CommentRoundTrip_ClearComment(t *testing.T) {
	ctx := context.Background()
	tempPath := filepath.Join(t.TempDir(), "test.nvpk")

	// Create with comment
	pkg1, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if err := pkg1.Create(ctx, tempPath); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	fpkg1 := pkg1.(*filePackage)
	if err := fpkg1.SetComment("Comment to be cleared"); err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}
	// Write package (0 files is valid)
	if err := pkg1.Write(ctx); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}
	if err := pkg1.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Reopen and clear comment
	pkg2, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}

	fpkg2 := pkg2.(*filePackage)
	if err := fpkg2.ClearComment(); err != nil {
		t.Fatalf("ClearComment() failed: %v", err)
	}
	if err := pkg2.Write(ctx); err != nil {
		t.Fatalf("Write() after clear failed: %v", err)
	}
	if err := pkg2.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// Reopen and verify no comment
	pkg3, err := OpenPackage(ctx, tempPath)
	if err != nil {
		t.Fatalf("OpenPackage() second time failed: %v", err)
	}
	defer func() { _ = pkg3.Close() }()

	if pkg3.HasComment() {
		t.Error("HasComment() after clear = true, want false")
	}
	if pkg3.GetComment() != "" {
		t.Errorf("GetComment() after clear = %q, want empty string", pkg3.GetComment())
	}
}
