// This file contains unit tests for package comment management operations.
// It tests SetComment, GetComment, ClearComment, and HasComment methods
// from package_comment.go.
//
// Specification: api_metadata.md: 1. Comment Management

package novus_package

import (
	"strings"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/metadata"
)

// =============================================================================
// TEST: SetComment
// =============================================================================

const testCommentStr = "Test comment"

// TestPackage_SetComment_Basic tests basic SetComment operation.
func TestPackage_SetComment_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set a valid comment
	comment := "Test package comment"
	err = fpkg.SetComment(comment)
	if err != nil {
		t.Errorf("SetComment() failed: %v", err)
	}

	// Verify comment was set
	if !fpkg.HasComment() {
		t.Error("HasComment() should return true after SetComment")
	}

	retrieved := fpkg.GetComment()
	if retrieved != comment {
		t.Errorf("GetComment() = %q, want %q", retrieved, comment)
	}
}

// TestPackage_SetComment_EmptyComment tests setting empty comment.
func TestPackage_SetComment_EmptyComment(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set empty comment
	err = fpkg.SetComment("")
	if err != nil {
		t.Errorf("SetComment() with empty string should succeed, got error: %v", err)
	}

	// Empty comment should not set HasComment flag
	if fpkg.HasComment() {
		t.Error("HasComment() should return false for empty comment")
	}
}

// TestPackage_SetComment_UpdatesHeaderFlags tests that SetComment updates header flags.
func TestPackage_SetComment_UpdatesHeaderFlags(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set comment
	comment := testCommentStr
	err = fpkg.SetComment(comment)
	if err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	// Verify header flag is set
	if (fpkg.header.Flags & fileformat.FlagHasPackageComment) == 0 {
		t.Error("Header flag FlagHasPackageComment should be set after SetComment")
	}

	// Verify CommentSize is set
	if fpkg.header.CommentSize == 0 {
		t.Error("Header CommentSize should be set after SetComment")
	}
}

// TestPackage_SetComment_InvalidUTF8 tests SetComment with invalid UTF-8.
func TestPackage_SetComment_InvalidUTF8(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Invalid UTF-8 sequence
	invalidUTF8 := string([]byte{0xff, 0xfe, 0xfd})
	err = fpkg.SetComment(invalidUTF8)
	if err == nil {
		t.Error("SetComment should fail with invalid UTF-8")
	}
}

// TestPackage_SetComment_TooLong tests SetComment with comment exceeding maximum length.
func TestPackage_SetComment_TooLong(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Create comment exceeding maximum length (MaxCommentLength is typically 1MB)
	// For testing, we'll use a smaller but still invalid length
	// Note: Actual max length check happens in PackageComment.SetComment
	longComment := strings.Repeat("a", int(metadata.MaxCommentLength)+1)
	err = fpkg.SetComment(longComment)
	if err == nil {
		t.Error("SetComment should fail with comment exceeding maximum length")
	}
}

// TestPackage_SetComment_WithContext tests SetComment (no longer applicable since SetComment doesn't take context).
// This test is kept for reference but SetComment is now a pure in-memory operation per spec.
func TestPackage_SetComment_WithContext(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// SetComment is a pure in-memory operation and doesn't take context per spec
	err = fpkg.SetComment("test")
	if err != nil {
		t.Errorf("SetComment() failed: %v", err)
	}
}

// =============================================================================
// TEST: GetComment
// =============================================================================

// TestPackage_GetComment_Basic tests basic GetComment operation.
func TestPackage_GetComment_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Initially no comment
	comment := fpkg.GetComment()
	if comment != "" {
		t.Errorf("GetComment() = %q, want empty string", comment)
	}

	// Set comment and retrieve
	testComment := testCommentStr
	err = fpkg.SetComment(testComment)
	if err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	comment = fpkg.GetComment()
	if comment != testComment {
		t.Errorf("GetComment() = %q, want %q", comment, testComment)
	}
}

// =============================================================================
// TEST: ClearComment
// =============================================================================

// TestPackage_ClearComment_Basic tests basic ClearComment operation.
func TestPackage_ClearComment_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set comment first
	err = fpkg.SetComment(testCommentStr)
	if err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	// Clear comment
	err = fpkg.ClearComment()
	if err != nil {
		t.Errorf("ClearComment() failed: %v", err)
	}

	// Verify comment is cleared
	if fpkg.HasComment() {
		t.Error("HasComment() should return false after ClearComment")
	}

	if fpkg.GetComment() != "" {
		t.Error("GetComment() should return empty string after ClearComment")
	}

	// Verify header flag is cleared
	if (fpkg.header.Flags & fileformat.FlagHasPackageComment) != 0 {
		t.Error("Header flag FlagHasPackageComment should be cleared after ClearComment")
	}

	// Verify CommentSize is cleared
	if fpkg.header.CommentSize != 0 {
		t.Errorf("Header CommentSize = %d, want 0", fpkg.header.CommentSize)
	}
}

func runClearCommentSucceeds(t *testing.T, errMsg string) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	fpkg := pkg.(*filePackage)
	err = fpkg.ClearComment()
	if err != nil {
		t.Errorf("ClearComment() %s: %v", errMsg, err)
	}
}

// TestPackage_ClearComment_NoComment tests clearing comment when no comment exists.
func TestPackage_ClearComment_NoComment(t *testing.T) {
	runClearCommentSucceeds(t, "should succeed when no comment exists")
}

// TestPackage_ClearComment_WithContext tests ClearComment (no longer applicable since ClearComment doesn't take context).
// This test is kept for reference but ClearComment is now a pure in-memory operation per spec.
func TestPackage_ClearComment_WithContext(t *testing.T) {
	runClearCommentSucceeds(t, "failed")
}

// =============================================================================
// TEST: HasComment
// =============================================================================

// TestPackage_HasComment_Basic tests basic HasComment operation.
func TestPackage_HasComment_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Initially no comment
	if fpkg.HasComment() {
		t.Error("HasComment() should return false for new package")
	}

	// Set comment
	err = fpkg.SetComment(testCommentStr)
	if err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	// Verify HasComment returns true
	if !fpkg.HasComment() {
		t.Error("HasComment() should return true after SetComment")
	}

	// Clear comment
	err = fpkg.ClearComment()
	if err != nil {
		t.Fatalf("ClearComment() failed: %v", err)
	}

	// Verify HasComment returns false
	if fpkg.HasComment() {
		t.Error("HasComment() should return false after ClearComment")
	}
}

// TestPackage_SetComment_WithNilInfo tests SetComment when Info is nil.
func TestPackage_SetComment_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Temporarily set Info to nil
	originalInfo := fpkg.Info
	fpkg.Info = nil

	// SetComment should handle nil Info gracefully
	comment := testCommentStr
	err = fpkg.SetComment(comment)
	if err != nil {
		t.Errorf("SetComment() should handle nil Info, got error: %v", err)
	}

	// Restore Info
	fpkg.Info = originalInfo
}

func runGetCommentWithNilInfo(t *testing.T) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	fpkg := pkg.(*filePackage)
	originalInfo := fpkg.Info
	fpkg.Info = nil
	defer func() { fpkg.Info = originalInfo }()
	if comment := fpkg.GetComment(); comment != "" {
		t.Errorf("GetComment() = %q, want empty string when Info is nil", comment)
	}
}

// TestPackage_GetComment_WithNilInfo tests GetComment when Info is nil.
func TestPackage_GetComment_WithNilInfo(t *testing.T) {
	runGetCommentWithNilInfo(t)
}

// TestPackage_HasComment_WithNilInfo tests HasComment when Info is nil.
func TestPackage_HasComment_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Temporarily set Info to nil
	originalInfo := fpkg.Info
	fpkg.Info = nil

	// HasComment should return false when Info is nil
	if fpkg.HasComment() {
		t.Error("HasComment() should return false when Info is nil")
	}

	// Restore Info
	fpkg.Info = originalInfo
}

// TestPackage_ClearComment_WithNilInfo tests ClearComment when Info is nil.
func TestPackage_ClearComment_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set comment first
	err = fpkg.SetComment(testCommentStr)
	if err != nil {
		t.Fatalf("SetComment() failed: %v", err)
	}

	// Temporarily set Info to nil
	originalInfo := fpkg.Info
	fpkg.Info = nil

	// ClearComment should handle nil Info gracefully
	err = fpkg.ClearComment()
	if err != nil {
		t.Errorf("ClearComment() should handle nil Info, got error: %v", err)
	}

	// Restore Info
	fpkg.Info = originalInfo
}
