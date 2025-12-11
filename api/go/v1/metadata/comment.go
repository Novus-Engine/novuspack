// Package novuspack provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and docs/tech_specs/api_metadata.md.
package metadata

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1 - Package Comment Format Specification
type PackageComment struct {
	// CommentLength is the length of comment including null terminator
	// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1
	CommentLength uint32

	// Comment is the UTF-8 encoded comment string (null-terminated)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1
	Comment string

	// Reserved is reserved for future use (must be 0)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1.1
func (p *PackageComment) Validate() error {
	// Check if comment exists
	if p.CommentLength == 0 {
		if p.Comment != "" {
			return fmt.Errorf("comment length is zero but comment is not empty")
		}
		// Empty comment is valid
		// Verify reserved bytes are zero
		for i, b := range p.Reserved {
			if b != 0 {
				return fmt.Errorf("reserved byte %d must be zero, got %d", i, b)
			}
		}
		return nil
	}

	// CommentLength must not exceed maximum
	if p.CommentLength > MaxCommentLength {
		return fmt.Errorf("comment length %d exceeds maximum %d", p.CommentLength, MaxCommentLength)
	}

	// Comment must not be empty if CommentLength > 0
	if p.Comment == "" {
		return fmt.Errorf("comment is empty but comment length is %d", p.CommentLength)
	}

	// Comment must be valid UTF-8
	if !utf8.ValidString(p.Comment) {
		return fmt.Errorf("comment is not valid UTF-8")
	}

	// Comment must be null-terminated
	commentBytes := []byte(p.Comment)
	if len(commentBytes) == 0 || commentBytes[len(commentBytes)-1] != 0x00 {
		return fmt.Errorf("comment is not null-terminated")
	}

	// Comment must not contain embedded null characters (except at the end)
	for i := 0; i < len(commentBytes)-1; i++ {
		if commentBytes[i] == 0x00 {
			return fmt.Errorf("comment contains embedded null character at position %d", i)
		}
	}

	// CommentLength must match actual comment size including null terminator
	expectedLength := uint32(len(commentBytes))
	if p.CommentLength != expectedLength {
		return fmt.Errorf("comment length mismatch: specified %d, actual %d", p.CommentLength, expectedLength)
	}

	// Verify reserved bytes are zero
	for i, b := range p.Reserved {
		if b != 0 {
			return fmt.Errorf("reserved byte %d must be zero, got %d", i, b)
		}
	}

	return nil
}

// Size returns the total size of the PackageComment in bytes.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1
func (p *PackageComment) Size() int {
	return 4 + int(p.CommentLength) + 3 // Length(4) + Comment + Reserved(3)
}

// IsEmpty returns true if the comment is empty (CommentLength == 0).
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1.1
func (p *PackageComment) IsEmpty() bool {
	return p.CommentLength == 0
}

// SetComment sets the comment text and updates CommentLength.
//
// The comment will be automatically null-terminated if not already.
// CommentLength will be set to the length of the comment including the null terminator.
//
// Returns an error if the comment exceeds MaxCommentLength or contains invalid UTF-8.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1.1
func (p *PackageComment) SetComment(comment string) error {
	// Validate UTF-8 before processing
	if comment != "" && !utf8.ValidString(comment) {
		return fmt.Errorf("comment is not valid UTF-8")
	}

	// Remove existing null terminator if present (before checking for embedded nulls)
	commentBytes := []byte(comment)
	if len(commentBytes) > 0 && commentBytes[len(commentBytes)-1] == 0x00 {
		commentBytes = commentBytes[:len(commentBytes)-1]
	}

	// Check for embedded null characters (after removing trailing null)
	for i, b := range commentBytes {
		if b == 0x00 {
			return fmt.Errorf("comment contains embedded null character at position %d", i)
		}
	}

	// Calculate length including null terminator
	length := uint32(len(commentBytes) + 1) // +1 for null terminator

	// Check maximum length
	if length > MaxCommentLength {
		return fmt.Errorf("comment length %d exceeds maximum %d", length, MaxCommentLength)
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1.1
func (p *PackageComment) GetComment() string {
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1.1
func (p *PackageComment) Clear() {
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1
func (p *PackageComment) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read CommentLength (4 bytes)
	var commentLength uint32
	if err := binary.Read(r, binary.LittleEndian, &commentLength); err != nil {
		return totalRead, fmt.Errorf("failed to read comment length: %w", err)
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
			return totalRead, fmt.Errorf("failed to read comment: %w", err)
		}
		if uint32(n) != commentLength {
			return totalRead, fmt.Errorf("incomplete comment read: got %d bytes, expected %d", n, commentLength)
		}
		totalRead += int64(n)

		// Convert to string
		p.Comment = string(commentBytes)
	}

	// Read Reserved (3 bytes)
	reservedBytes := make([]byte, 3)
	n, err := io.ReadFull(r, reservedBytes)
	if err != nil {
		return totalRead, fmt.Errorf("failed to read reserved bytes: %w", err)
	}
	if n != 3 {
		return totalRead, fmt.Errorf("incomplete reserved bytes read: got %d bytes, expected 3", n)
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1
func (p *PackageComment) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write CommentLength (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.CommentLength); err != nil {
		return totalWritten, fmt.Errorf("failed to write comment length: %w", err)
	}
	totalWritten += 4

	// If CommentLength > 0, write Comment
	if p.CommentLength > 0 {
		commentBytes := []byte(p.Comment)
		n, err := w.Write(commentBytes)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write comment: %w", err)
		}
		if uint32(n) != p.CommentLength {
			return totalWritten, fmt.Errorf("incomplete comment write: wrote %d bytes, expected %d", n, p.CommentLength)
		}
		totalWritten += int64(n)
	}

	// Write Reserved (3 bytes)
	n, err := w.Write(p.Reserved[:])
	if err != nil {
		return totalWritten, fmt.Errorf("failed to write reserved bytes: %w", err)
	}
	if n != 3 {
		return totalWritten, fmt.Errorf("incomplete reserved bytes write: wrote %d bytes, expected 3", n)
	}
	totalWritten += int64(n)

	return totalWritten, nil
}
