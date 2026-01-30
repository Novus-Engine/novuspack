// This file implements package-level comment management operations.
// It contains methods for setting, getting, clearing, and checking package comments
// as specified in api_metadata.md Section 1. This file should contain only
// Package-level comment operations (SetComment, GetComment, ClearComment, HasComment).
//
// Specification: api_metadata.md: 1. Comment Management

package novus_package

import (
	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/metadata"
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

	// Create PackageComment instance
	pc := metadata.NewPackageComment()

	// Set comment using PackageComment.SetComment which handles validation
	if err := pc.SetComment(comment); err != nil {
		return err
	}

	// Validate the comment
	if err := pc.Validate(); err != nil {
		return err
	}

	// Update header fields
	p.header.CommentSize = pc.CommentLength
	// Note: CommentStart will be set when package is written to disk
	// For now, we set it to 0 as a placeholder
	p.header.CommentStart = 0

	// Update header flags (bit 4 = FlagHasPackageComment)
	if pc.CommentLength > 0 {
		p.header.SetFeature(fileformat.FlagHasPackageComment)
	} else {
		p.header.ClearFeature(fileformat.FlagHasPackageComment)
	}

	// Update PackageInfo
	// HasComment should be true only if there's actual comment text (not just null terminator)
	commentText := pc.GetComment()
	if p.Info != nil {
		p.Info.HasComment = len(commentText) > 0
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
	p.header.ClearFeature(fileformat.FlagHasPackageComment)

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
