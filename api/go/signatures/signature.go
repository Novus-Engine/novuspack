// This file implements the Signature structure and signature management operations.
// It contains the Signature type definition, signature creation and validation
// methods, and signature management functions. This file should contain all code
// related to signature operations as specified in api_signatures.md Section 1
// and Section 2.
//
// Specification: api_signatures.md: 1. Signature Management

// Package novuspack provides signatures domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to digital signatures
// as specified in docs/tech_specs/package_file_format.md and api_signatures.md.
package signatures

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// Signature represents a digital signature for package integrity verification.
//
// Size: Variable (18 bytes + comment_length + signature_size)
//
// Specification: package_file_format.md: 8.1 Signature Structure
type Signature struct {
	// SignatureType is the signature algorithm identifier
	// Specification: package_file_format.md: 1. `.nvpk` File Format Overview
	SignatureType uint32

	// SignatureSize is the size of signature data in bytes
	// Specification: package_file_format.md: 7.1 Package Comment Format Specification
	SignatureSize uint32

	// SignatureFlags contains signature-specific metadata
	// Specification: package_file_format.md: 1. `.nvpk` File Format Overview
	SignatureFlags uint32

	// SignatureTimestamp is the signature creation time (Unix nanoseconds)
	// Specification: package_file_format.md: 1. `.nvpk` File Format Overview
	SignatureTimestamp uint32

	// CommentLength is the length of signature comment
	// Specification: package_file_format.md: 1. `.nvpk` File Format Overview
	CommentLength uint16

	// SignatureComment is a human-readable comment about the signature
	// Specification: package_file_format.md: 7.1 Package Comment Format Specification
	SignatureComment string

	// SignatureData contains the raw signature bytes
	// Specification: package_file_format.md: 7.1 Package Comment Format Specification
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
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "signature type cannot be zero", nil, pkgerrors.ValidationErrorContext{
			Field:    "SignatureType",
			Value:    s.SignatureType,
			Expected: "non-zero value",
		})
	}

	if len(s.SignatureData) == 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "signature data cannot be nil or empty", nil, pkgerrors.ValidationErrorContext{
			Field:    "SignatureData",
			Value:    nil,
			Expected: "non-empty signature data",
		})
	}

	if uint32(len(s.SignatureData)) != s.SignatureSize {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "signature size mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "SignatureSize",
			Value:    s.SignatureSize,
			Expected: fmt.Sprintf("%d", len(s.SignatureData)),
		})
	}

	if s.SignatureComment != "" {
		if uint16(len(s.SignatureComment)) != s.CommentLength {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length mismatch", nil, pkgerrors.ValidationErrorContext{
				Field:    "CommentLength",
				Value:    s.CommentLength,
				Expected: fmt.Sprintf("%d", len(s.SignatureComment)),
			})
		}
	} else if s.CommentLength != 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length is non-zero but comment is empty", nil, pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    s.CommentLength,
			Expected: "0 or non-empty comment",
		})
	}

	return nil
}

// Size returns the total size of the Signature in bytes.
//
// Specification: package_file_format.md: 7.1 Package Comment Format Specification
func (s *Signature) Size() int {
	// Type(4) + Size(4) + Flags(4) + Timestamp(4) + CommentLength(2) + Comment + Data
	return 18 + int(s.CommentLength) + int(s.SignatureSize)
}

// NewSignature creates and returns a new Signature with zero values.
//
// The returned Signature is initialized with all fields set to zero or empty.
//
// Specification: package_file_format.md: 8.1 Signature Structure
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
// Specification: package_file_format.md: 8.1 Signature Structure
func (s *Signature) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read SignatureType (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureType); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, "failed to read signature type: incomplete data", pkgerrors.ValidationErrorContext{
				Field:    "SignatureType",
				Value:    totalRead,
				Expected: "4 bytes",
			})
		}
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature type", pkgerrors.ValidationErrorContext{
			Field:    "SignatureType",
			Value:    nil,
			Expected: "4 bytes",
		})
	}
	totalRead += 4

	// Read SignatureSize (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureSize); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature size", pkgerrors.ValidationErrorContext{
			Field:    "SignatureSize",
			Value:    nil,
			Expected: "4 bytes",
		})
	}
	totalRead += 4

	// Read SignatureFlags (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureFlags); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature flags", pkgerrors.ValidationErrorContext{
			Field:    "SignatureFlags",
			Value:    nil,
			Expected: "4 bytes",
		})
	}
	totalRead += 4

	// Read SignatureTimestamp (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureTimestamp); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature timestamp", pkgerrors.ValidationErrorContext{
			Field:    "SignatureTimestamp",
			Value:    nil,
			Expected: "4 bytes",
		})
	}
	totalRead += 4

	// Read CommentLength (2 bytes)
	if err := binary.Read(r, binary.LittleEndian, &s.CommentLength); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read comment length", pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    nil,
			Expected: "2 bytes",
		})
	}
	totalRead += 2

	// Read SignatureComment (CommentLength bytes)
	if s.CommentLength > 0 {
		commentBytes := make([]byte, s.CommentLength)
		n, err := io.ReadFull(r, commentBytes)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature comment", pkgerrors.ValidationErrorContext{
				Field:    "SignatureComment",
				Value:    s.CommentLength,
				Expected: "comment data",
			})
		}
		if uint16(n) != s.CommentLength {
			return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete comment read", nil, pkgerrors.ValidationErrorContext{
				Field:    "SignatureComment",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", s.CommentLength),
			})
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
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature data", pkgerrors.ValidationErrorContext{
				Field:    "SignatureData",
				Value:    s.SignatureSize,
				Expected: "signature data",
			})
		}
		if uint32(n) != s.SignatureSize {
			return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete signature data read", nil, pkgerrors.ValidationErrorContext{
				Field:    "SignatureData",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", s.SignatureSize),
			})
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
// Specification: package_file_format.md: 8.1 Signature Structure
func (s *Signature) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Update CommentLength and SignatureSize to match actual data
	s.CommentLength = uint16(len(s.SignatureComment))
	s.SignatureSize = uint32(len(s.SignatureData))

	// Write SignatureType (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureType); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature type", pkgerrors.ValidationErrorContext{
			Field:    "SignatureType",
			Value:    s.SignatureType,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	// Write SignatureSize (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureSize); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature size", pkgerrors.ValidationErrorContext{
			Field:    "SignatureSize",
			Value:    s.SignatureSize,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	// Write SignatureFlags (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureFlags); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature flags", pkgerrors.ValidationErrorContext{
			Field:    "SignatureFlags",
			Value:    s.SignatureFlags,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	// Write SignatureTimestamp (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.SignatureTimestamp); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature timestamp", pkgerrors.ValidationErrorContext{
			Field:    "SignatureTimestamp",
			Value:    s.SignatureTimestamp,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	// Write CommentLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, s.CommentLength); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write comment length", pkgerrors.ValidationErrorContext{
			Field:    "CommentLength",
			Value:    s.CommentLength,
			Expected: "written successfully",
		})
	}
	totalWritten += 2

	// Write SignatureComment (CommentLength bytes)
	if s.CommentLength > 0 {
		commentBytes := []byte(s.SignatureComment)
		if uint16(len(commentBytes)) != s.CommentLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length mismatch", nil, pkgerrors.ValidationErrorContext{
				Field:    "CommentLength",
				Value:    len(commentBytes),
				Expected: fmt.Sprintf("%d", s.CommentLength),
			})
		}
		n, err := w.Write(commentBytes)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature comment", pkgerrors.ValidationErrorContext{
				Field:    "SignatureComment",
				Value:    s.SignatureComment,
				Expected: "written successfully",
			})
		}
		if uint16(n) != s.CommentLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete comment write", nil, pkgerrors.ValidationErrorContext{
				Field:    "SignatureComment",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", s.CommentLength),
			})
		}
		totalWritten += int64(n)
	}

	// Write SignatureData (SignatureSize bytes)
	if s.SignatureSize > 0 {
		if uint32(len(s.SignatureData)) != s.SignatureSize {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "signature size mismatch", nil, pkgerrors.ValidationErrorContext{
				Field:    "SignatureSize",
				Value:    len(s.SignatureData),
				Expected: fmt.Sprintf("%d", s.SignatureSize),
			})
		}
		n, err := w.Write(s.SignatureData)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature data", pkgerrors.ValidationErrorContext{
				Field:    "SignatureData",
				Value:    s.SignatureData,
				Expected: "written successfully",
			})
		}
		if uint32(n) != s.SignatureSize {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete signature data write", nil, pkgerrors.ValidationErrorContext{
				Field:    "SignatureData",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", s.SignatureSize),
			})
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}

// HasFlag checks if a specific signature flag is set.
//
// Specification: package_file_format.md: 8.2.2 SignatureFlags Field
func (s *Signature) HasFlag(flag uint32) bool {
	return (s.SignatureFlags & flag) != 0
}

// SetFlag sets a specific signature flag.
//
// Specification: package_file_format.md: 8.2.2 SignatureFlags Field
func (s *Signature) SetFlag(flag uint32) {
	s.SignatureFlags |= flag
}

// ClearFlag clears a specific signature flag.
//
// Specification: package_file_format.md: 8.2.2 SignatureFlags Field
func (s *Signature) ClearFlag(flag uint32) {
	s.SignatureFlags &= ^flag
}
