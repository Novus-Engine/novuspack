// This file implements package-level comment management operations.
// It contains methods for setting, getting, clearing, and checking package comments
// as specified in api_metadata.md Section 1. This file should contain only
// Package-level comment operations (SetComment, GetComment, ClearComment, HasComment).
//
// Specification: api_metadata.md: 1. Comment Management

package novus_package

import (
	"fmt"
	"unicode/utf8"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// SetComment sets or updates the package comment.
//
// Sets the package comment, updates header flags (bit 4), CommentSize, and
// CommentStart fields. The comment is validated before being set.
// This is a pure in-memory operation and does not require context.
//
// Parameters:
//   - comment: Comment string to set
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 1. Comment Management
func (p *filePackage) SetComment(comment string) error {

	pc, err := buildPackageComment(comment)
	if err != nil {
		return err
	}

	// Update header fields
	p.header.CommentSize = pc.CommentLength
	// Note: CommentStart will be set when package is written to disk
	// For now, we set it to 0 as a placeholder
	p.header.CommentStart = 0

	// Update header flags (bit 4 = FlagHasPackageComment)
	if pc.CommentLength > 0 {
		p.header.Flags |= fileformat.FlagHasPackageComment
	} else {
		p.header.Flags &^= fileformat.FlagHasPackageComment
	}

	// Update PackageInfo
	// HasComment should be true only if there's actual comment text (not just null terminator)
	commentText := extractCommentText(pc)
	if p.Info != nil {
		p.Info.HasComment = commentText != ""
		p.Info.Comment = commentText
		// Increment MetadataVersion in PackageInfo (metadata changed)
		p.Info.MetadataVersion++
	}

	return nil
}

// GetComment retrieves the current package comment.
//
// Returns the comment text without the null terminator, or empty string if no comment is set.
// This is a pure data access method and does not require context.
//
// Returns:
//   - string: Current package comment, or empty string if no comment is set
//
// Specification: api_metadata.md: 1. Comment Management
func (p *filePackage) GetComment() string {
	if p.Info == nil {
		return ""
	}
	return p.Info.Comment
}

// ClearComment removes the package comment.
//
// Clears the package comment, updates header flags, and resets CommentSize and CommentStart.
// This is a pure in-memory operation and does not require context.
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 1.1 Package-Level Comment Methods
func (p *filePackage) ClearComment() error {

	// Clear header fields
	p.header.CommentSize = 0
	p.header.CommentStart = 0

	// Clear header flags (bit 4 = FlagHasPackageComment)
	p.header.Flags &^= fileformat.FlagHasPackageComment

	// Update PackageInfo
	if p.Info != nil {
		p.Info.HasComment = false
		p.Info.Comment = ""
		// Increment MetadataVersion in PackageInfo (metadata changed)
		p.Info.MetadataVersion++
	}

	return nil
}

// HasComment checks if the package has a comment.
//
// Returns true if the package has a non-empty comment, false otherwise.
// This is a pure data access method and does not require context.
//
// Returns:
//   - bool: True if package has a comment, false otherwise
//
// Specification: api_metadata.md: 1.1 Package-Level Comment Methods
func (p *filePackage) HasComment() bool {
	if p.Info == nil {
		return false
	}
	return p.Info.HasComment
}

func buildPackageComment(comment string) (*metadata.PackageComment, error) {
	pc := metadata.NewPackageComment()
	if comment != "" && !utf8.ValidString(comment) {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment is not valid UTF-8", nil, pkgerrors.ValidationErrorContext{
			Field:    "Comment",
			Value:    comment,
			Expected: "valid UTF-8 string",
		})
	}

	commentBytes := []byte(comment)
	if len(commentBytes) > 0 && commentBytes[len(commentBytes)-1] == 0x00 {
		commentBytes = commentBytes[:len(commentBytes)-1]
	}
	for i, b := range commentBytes {
		if b == 0x00 {
			return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("comment contains embedded null character at position %d", i), nil, pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    i,
				Expected: "no embedded null characters",
			})
		}
	}

	length := uint32(len(commentBytes) + 1)
	if length > metadata.MaxCommentLength {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length exceeds maximum", nil, pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    length,
			Expected: fmt.Sprintf("<= %d", metadata.MaxCommentLength),
		})
	}

	if len(commentBytes) > 0 {
		pc.Comment = string(commentBytes) + "\x00"
	} else {
		pc.Comment = "\x00"
	}
	pc.CommentLength = length
	pc.Reserved = [3]uint8{0, 0, 0}

	if err := pc.Validate(); err != nil {
		return nil, err
	}

	return pc, nil
}

func extractCommentText(pc *metadata.PackageComment) string {
	if pc == nil || pc.CommentLength == 0 || pc.Comment == "" {
		return ""
	}

	commentBytes := []byte(pc.Comment)
	if len(commentBytes) > 0 && commentBytes[len(commentBytes)-1] == 0x00 {
		return string(commentBytes[:len(commentBytes)-1])
	}

	return pc.Comment
}
