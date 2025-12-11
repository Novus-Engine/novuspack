// Package novuspack provides signatures domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to digital signatures
// as specified in docs/tech_specs/package_file_format.md and docs/tech_specs/api_signatures.md.
package signatures

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Signature represents a digital signature for package integrity verification.
//
// Size: Variable (18 bytes + comment_length + signature_size)
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1 - Signature Structure
type Signature struct {
	// SignatureType is the signature algorithm identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.1
	SignatureType uint32

	// SignatureSize is the size of signature data in bytes
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1
	SignatureSize uint32

	// SignatureFlags contains signature-specific metadata
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.2
	SignatureFlags uint32

	// SignatureTimestamp is the signature creation time (Unix nanoseconds)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.3
	SignatureTimestamp uint32

	// CommentLength is the length of signature comment
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.4
	CommentLength uint16

	// SignatureComment is a human-readable comment about the signature
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1
	SignatureComment string

	// SignatureData contains the raw signature bytes
	// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1
	SignatureData []byte
}

// Validate performs validation checks on the Signature.
//
// Validation checks:
//   - SignatureType must be valid
//   - SignatureSize must match actual SignatureData length
//   - CommentLength must match actual SignatureComment length (if present)
//   - SignatureData must not be nil or empty
//
// Returns an error if any validation check fails.
func (s *Signature) Validate() error {
	if s.SignatureType == 0 {
		return fmt.Errorf("signature type cannot be zero")
	}

	if len(s.SignatureData) == 0 {
		return fmt.Errorf("signature data cannot be nil or empty")
	}

	if uint32(len(s.SignatureData)) != s.SignatureSize {
		return fmt.Errorf("signature size mismatch: specified %d, actual %d", s.SignatureSize, len(s.SignatureData))
	}

	if s.SignatureComment != "" {
		if uint16(len(s.SignatureComment)) != s.CommentLength {
			return fmt.Errorf("comment length mismatch: specified %d, actual %d", s.CommentLength, len(s.SignatureComment))
		}
	} else if s.CommentLength != 0 {
		return fmt.Errorf("comment length is %d but comment is empty", s.CommentLength)
	}

	return nil
}

// Size returns the total size of the Signature in bytes.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1
func (s *Signature) Size() int {
	// Type(4) + Size(4) + Flags(4) + Timestamp(4) + CommentLength(2) + Comment + Data
	return 18 + int(s.CommentLength) + int(s.SignatureSize)
}

// NewSignature creates and returns a new Signature with zero values.
//
// The returned Signature is initialized with all fields set to zero or empty.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1 - Signature Structure
func NewSignature() *Signature {
	return &Signature{}
}

// ReadFrom reads a Signature from the provided io.Reader.
//
// The binary format is:
//   - SignatureType (4 bytes, little-endian uint32)
//   - SignatureSize (4 bytes, little-endian uint32)
//   - SignatureFlags (4 bytes, little-endian uint32)
//   - SignatureTimestamp (4 bytes, little-endian uint32)
//   - CommentLength (2 bytes, little-endian uint16)
//   - SignatureComment (CommentLength bytes, UTF-8 string with null terminator)
//   - SignatureData (SignatureSize bytes)
//
// Returns the number of bytes read and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1 - Signature Structure
func (s *Signature) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read SignatureType (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureType); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return totalRead, fmt.Errorf("failed to read signature type: incomplete data: %w", err)
		}
		return totalRead, fmt.Errorf("failed to read signature type: %w", err)
	}
	totalRead += 4

	// Read SignatureSize (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureSize); err != nil {
		return totalRead, fmt.Errorf("failed to read signature size: %w", err)
	}
	totalRead += 4

	// Read SignatureFlags (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureFlags); err != nil {
		return totalRead, fmt.Errorf("failed to read signature flags: %w", err)
	}
	totalRead += 4

	// Read SignatureTimestamp (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureTimestamp); err != nil {
		return totalRead, fmt.Errorf("failed to read signature timestamp: %w", err)
	}
	totalRead += 4

	// Read CommentLength (2 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.CommentLength); err != nil {
		return totalRead, fmt.Errorf("failed to read comment length: %w", err)
	}
	totalRead += 2

	// Read SignatureComment (CommentLength bytes)
	if s.CommentLength > 0 {
		commentBytes := make([]byte, s.CommentLength)
		n, err := io.ReadFull(r, commentBytes)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read signature comment: %w", err)
		}
		if uint16(n) != s.CommentLength {
			return totalRead, fmt.Errorf("incomplete comment read: got %d bytes, expected %d", n, s.CommentLength)
		}
		totalRead += int64(n)
		s.SignatureComment = string(commentBytes)
	} else {
		s.SignatureComment = ""
	}

	// Read SignatureData (SignatureSize bytes)
	if s.SignatureSize > 0 {
		signatureData := make([]byte, s.SignatureSize)
		n, err := io.ReadFull(r, signatureData)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read signature data: %w", err)
		}
		if uint32(n) != s.SignatureSize {
			return totalRead, fmt.Errorf("incomplete signature data read: got %d bytes, expected %d", n, s.SignatureSize)
		}
		totalRead += int64(n)
		s.SignatureData = signatureData
	} else {
		s.SignatureData = nil
	}

	return totalRead, nil
}

// WriteTo writes a Signature to the provided io.Writer.
//
// The binary format is:
//   - SignatureType (4 bytes, little-endian uint32)
//   - SignatureSize (4 bytes, little-endian uint32)
//   - SignatureFlags (4 bytes, little-endian uint32)
//   - SignatureTimestamp (4 bytes, little-endian uint32)
//   - CommentLength (2 bytes, little-endian uint16)
//   - SignatureComment (CommentLength bytes, UTF-8 string with null terminator)
//   - SignatureData (SignatureSize bytes)
//
// Before writing, the method updates CommentLength and SignatureSize to match
// the actual comment and data lengths.
//
// Returns the number of bytes written and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.1 - Signature Structure
func (s *Signature) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Update CommentLength and SignatureSize to match actual data
	s.CommentLength = uint16(len(s.SignatureComment))
	s.SignatureSize = uint32(len(s.SignatureData))

	// Write SignatureType (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureType); err != nil {
		return totalWritten, fmt.Errorf("failed to write signature type: %w", err)
	}
	totalWritten += 4

	// Write SignatureSize (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureSize); err != nil {
		return totalWritten, fmt.Errorf("failed to write signature size: %w", err)
	}
	totalWritten += 4

	// Write SignatureFlags (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureFlags); err != nil {
		return totalWritten, fmt.Errorf("failed to write signature flags: %w", err)
	}
	totalWritten += 4

	// Write SignatureTimestamp (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureTimestamp); err != nil {
		return totalWritten, fmt.Errorf("failed to write signature timestamp: %w", err)
	}
	totalWritten += 4

	// Write CommentLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.CommentLength); err != nil {
		return totalWritten, fmt.Errorf("failed to write comment length: %w", err)
	}
	totalWritten += 2

	// Write SignatureComment (CommentLength bytes)
	if s.CommentLength > 0 {
		commentBytes := []byte(s.SignatureComment)
		if uint16(len(commentBytes)) != s.CommentLength {
			return totalWritten, fmt.Errorf("comment length mismatch: specified %d, actual %d", s.CommentLength, len(commentBytes))
		}
		n, err := w.Write(commentBytes)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write signature comment: %w", err)
		}
		if uint16(n) != s.CommentLength {
			return totalWritten, fmt.Errorf("incomplete comment write: wrote %d bytes, expected %d", n, s.CommentLength)
		}
		totalWritten += int64(n)
	}

	// Write SignatureData (SignatureSize bytes)
	if s.SignatureSize > 0 {
		if uint32(len(s.SignatureData)) != s.SignatureSize {
			return totalWritten, fmt.Errorf("signature size mismatch: specified %d, actual %d", s.SignatureSize, len(s.SignatureData))
		}
		n, err := w.Write(s.SignatureData)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write signature data: %w", err)
		}
		if uint32(n) != s.SignatureSize {
			return totalWritten, fmt.Errorf("incomplete signature data write: wrote %d bytes, expected %d", n, s.SignatureSize)
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}

// HasFlag checks if a specific signature flag is set.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.2 - SignatureFlags Field
func (s *Signature) HasFlag(flag uint32) bool {
	return (s.SignatureFlags & flag) != 0
}

// SetFlag sets a specific signature flag.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.2 - SignatureFlags Field
func (s *Signature) SetFlag(flag uint32) {
	s.SignatureFlags |= flag
}

// ClearFlag clears a specific signature flag.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.2 - SignatureFlags Field
func (s *Signature) ClearFlag(flag uint32) {
	s.SignatureFlags &= ^flag
}
