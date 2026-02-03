// This file implements the PackageComment structure representing the optional
// package comment section. It contains the PackageComment type definition,
// NewPackageComment constructor, and methods for reading, writing, and validating
// package comments. This file should contain all code related to package comments
// as specified in api_metadata.md Section 1 and package_file_format.md Section 7.
//
// Specification: api_metadata.md: 1 Comment Management

// Package novuspack provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and api_metadata.md.
package metadata

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// NewPackageComment creates and returns a new PackageComment with zero values.
//
// The returned PackageComment has:
//   - CommentLength set to 0
//   - Comment set to empty string
//   - Reserved bytes all set to 0
//
// This is equivalent to an empty comment state.
func NewPackageComment() *PackageComment {
	return &PackageComment{
		CommentLength: 0,
		Comment:       "",
		Reserved:      [3]uint8{0, 0, 0},
	}
}

// PackageComment represents the optional package comment section.
//
// Size: Variable (4 + comment_length + 3 bytes)
//
// Specification: package_file_format.md: 7.1 Package Comment Format Specification
type PackageComment struct {
	// CommentLength is the length of comment including null terminator
	// Specification: api_metadata.md: 1.2 PackageComment Structure
	CommentLength uint32

	// Comment is the UTF-8 encoded comment string (null-terminated)
	// Specification: api_metadata.md: 1.2 PackageComment Structure
	Comment string

	// Reserved is reserved for future use (must be 0)
	// Specification: api_metadata.md: 1.2 PackageComment Structure
	Reserved [3]uint8
}

// Validate performs validation checks on the PackageComment.
//
// Validation checks:
//   - CommentLength must be > 0 if comment exists
//   - Comment must be valid UTF-8
//   - CommentLength must not exceed MaxCommentLength
//   - Comment must be null-terminated (ends with 0x00)
//   - Comment must not contain embedded null characters
//   - CommentLength must match actual comment size including null terminator
//   - Reserved bytes must all be zero
//
// Returns an error if any validation check fails.
//
// Specification: api_metadata.md: 1.3.4 PackageComment.Validate Method
func (p *PackageComment) Validate() error {
	if p.CommentLength == 0 {
		return p.validateEmptyComment()
	}
	if p.CommentLength > MaxCommentLength {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length exceeds maximum", nil, pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    p.CommentLength,
			Expected: fmt.Sprintf("<= %d", MaxCommentLength),
		})
	}

	// Comment must not be empty if CommentLength > 0
	if p.Comment == "" {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment is empty but comment length is non-zero", nil, pkgerrors.ValidationErrorContext{
			Field:    "Comment",
			Value:    "",
			Expected: "non-empty comment",
		})
	}

	// Comment must be valid UTF-8
	if !utf8.ValidString(p.Comment) {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment is not valid UTF-8", nil, pkgerrors.ValidationErrorContext{
			Field:    "Comment",
			Value:    p.Comment,
			Expected: "valid UTF-8 string",
		})
	}

	// Comment must be null-terminated
	commentBytes := []byte(p.Comment)
	if len(commentBytes) == 0 || commentBytes[len(commentBytes)-1] != 0x00 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment is not null-terminated", nil, pkgerrors.ValidationErrorContext{
			Field:    "Comment",
			Value:    p.Comment,
			Expected: "null-terminated string",
		})
	}

	// Comment must not contain embedded null characters (except at the end)
	for i := 0; i < len(commentBytes)-1; i++ {
		if commentBytes[i] == 0x00 {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("comment contains embedded null character at position %d", i), nil, pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    i,
				Expected: "no embedded null characters",
			})
		}
	}

	// CommentLength must match actual comment size including null terminator
	expectedLength := uint32(len(commentBytes))
	if p.CommentLength != expectedLength {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    p.CommentLength,
			Expected: fmt.Sprintf("%d", expectedLength),
		})
	}

	return p.validateReservedZero()
}

func (p *PackageComment) validateEmptyComment() error {
	if p.Comment != "" {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field: "CommentLength", Value: p.CommentLength, Expected: "non-zero when comment present",
		})
	}
	return p.validateReservedZero()
}

func (p *PackageComment) validateReservedZero() error {
	for i, b := range p.Reserved {
		if b != 0 {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("reserved byte %d must be zero", i), nil, pkgerrors.ValidationErrorContext{
				Field: "Reserved", Value: b, Expected: "0",
			})
		}
	}
	return nil
}

// Size returns the total size of the PackageComment in bytes.
//
// Specification: api_metadata.md: 1.3.1 PackageComment.Size Method
func (p *PackageComment) Size() int {
	return 4 + int(p.CommentLength) + 3 // Length(4) + Comment + Reserved(3)
}

// IsEmpty returns true if the comment is empty (CommentLength == 0).
//
// Specification: api_metadata.md: 1.3 PackageComment Methods
func (p *PackageComment) isEmpty() bool {
	return p.CommentLength == 0
}

// SetComment sets the comment text and updates CommentLength.
//
// The comment will be automatically null-terminated if not already.
// CommentLength will be set to the length of the comment including the null terminator.
//
// Returns an error if the comment exceeds MaxCommentLength or contains invalid UTF-8.
//
// Specification: api_metadata.md: 1.3 PackageComment Methods
func (p *PackageComment) setComment(comment string) error {
	// Validate UTF-8 before processing
	if comment != "" && !utf8.ValidString(comment) {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment is not valid UTF-8", nil, pkgerrors.ValidationErrorContext{
			Field:    "Comment",
			Value:    comment,
			Expected: "valid UTF-8 string",
		})
	}

	// Remove existing null terminator if present (before checking for embedded nulls)
	commentBytes := []byte(comment)
	if len(commentBytes) > 0 && commentBytes[len(commentBytes)-1] == 0x00 {
		commentBytes = commentBytes[:len(commentBytes)-1]
	}

	// Check for embedded null characters (after removing trailing null)
	for i, b := range commentBytes {
		if b == 0x00 {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, fmt.Sprintf("comment contains embedded null character at position %d", i), nil, pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    i,
				Expected: "no embedded null characters",
			})
		}
	}

	// Calculate length including null terminator
	length := uint32(len(commentBytes) + 1) // +1 for null terminator

	// Check maximum length
	if length > MaxCommentLength {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length exceeds maximum", nil, pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    length,
			Expected: fmt.Sprintf("<= %d", MaxCommentLength),
		})
	}

	// Set comment with null terminator
	if len(commentBytes) > 0 {
		p.Comment = string(commentBytes) + "\x00"
	} else {
		p.Comment = "\x00"
	}
	p.CommentLength = length

	// Clear reserved bytes
	p.Reserved = [3]uint8{0, 0, 0}

	return nil
}

// GetComment returns the comment text without the null terminator.
//
// Returns an empty string if the comment is empty.
//
// Specification: api_metadata.md: 1.3 PackageComment Methods
func (p *PackageComment) getComment() string {
	if p.CommentLength == 0 || p.Comment == "" {
		return ""
	}

	// Remove null terminator if present
	commentBytes := []byte(p.Comment)
	if len(commentBytes) > 0 && commentBytes[len(commentBytes)-1] == 0x00 {
		return string(commentBytes[:len(commentBytes)-1])
	}

	return p.Comment
}

// Clear removes the comment and resets all fields.
//
// Specification: api_metadata.md: 1.3 PackageComment Methods
func (p *PackageComment) clear() {
	p.CommentLength = 0
	p.Comment = ""
	p.Reserved = [3]uint8{0, 0, 0}
}

// ReadFrom reads a PackageComment from the provided io.Reader.
//
// The binary format is:
//   - CommentLength (4 bytes, little-endian uint32)
//   - Comment (CommentLength bytes, UTF-8 string with null terminator)
//   - Reserved (3 bytes, must be zero)
//
// If CommentLength is 0, the Comment field is skipped.
//
// Returns the number of bytes read and any error encountered.
//
// Specification: api_metadata.md: 1.3.3 PackageComment.ReadFrom Method
func (p *PackageComment) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read CommentLength (4 bytes)
	var commentLength uint32
	if err := binary.Read(r, binary.LittleEndian, &commentLength); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read comment length", pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    nil,
			Expected: "4 bytes",
		})
	}
	totalRead += 4
	p.CommentLength = commentLength

	// If CommentLength is 0, skip reading comment and just read reserved bytes
	if commentLength == 0 {
		p.Comment = ""
	} else {
		// Read Comment (CommentLength bytes)
		commentBytes := make([]byte, commentLength)
		n, err := io.ReadFull(r, commentBytes)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read comment", pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    commentLength,
				Expected: "comment data",
			})
		}
		if uint32(n) != commentLength {
			return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete comment read", nil, pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", commentLength),
			})
		}
		totalRead += int64(n)

		// Convert to string
		p.Comment = string(commentBytes)
	}

	// Read Reserved (3 bytes)
	reservedBytes := make([]byte, 3)
	n, err := io.ReadFull(r, reservedBytes)
	if err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read reserved bytes", pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    nil,
			Expected: "3 bytes",
		})
	}
	if n != 3 {
		return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete reserved bytes read", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    n,
			Expected: "3 bytes",
		})
	}
	totalRead += int64(n)

	// Copy reserved bytes to array
	copy(p.Reserved[:], reservedBytes)

	return totalRead, nil
}

// WriteTo writes a PackageComment to the provided io.Writer.
//
// The binary format is:
//   - CommentLength (4 bytes, little-endian uint32)
//   - Comment (CommentLength bytes, UTF-8 string with null terminator)
//   - Reserved (3 bytes, must be zero)
//
// If CommentLength is 0, the Comment field is skipped.
//
// Returns the number of bytes written and any error encountered.
//
// Specification: api_metadata.md: 1.3.2 PackageComment.WriteTo Method
func (p *PackageComment) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write CommentLength (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.CommentLength); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write comment length", pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    p.CommentLength,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	// If CommentLength > 0, write Comment
	if p.CommentLength > 0 {
		commentBytes := []byte(p.Comment)
		n, err := w.Write(commentBytes)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write comment", pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    p.Comment,
				Expected: "written successfully",
			})
		}
		if uint32(n) != p.CommentLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete comment write", nil, pkgerrors.ValidationErrorContext{
				Field:    "Comment",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", p.CommentLength),
			})
		}
		totalWritten += int64(n)
	}

	// Write Reserved (3 bytes)
	n, err := w.Write(p.Reserved[:])
	if err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write reserved bytes", pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    p.Reserved,
			Expected: "written successfully",
		})
	}
	if n != 3 {
		return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete reserved bytes write", nil, pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    n,
			Expected: "3 bytes",
		})
	}
	totalWritten += int64(n)

	return totalWritten, nil
}
