// This file implements the Signature structure and signature management operations.
// It contains the Signature type definition, signature creation and validation
// methods, and signature management functions. This file should contain all code
// related to signature operations as specified in api_signatures.md Section 1
// and Section 2.
//
// Specification: api_signatures.md: 4.1.4.1 Signature Struct

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
	// Specification: package_file_format.md: 8.1 Signature Structure
	SignatureType uint32

	// SignatureSize is the size of signature data in bytes
	// Specification: package_file_format.md: 8.1 Signature Structure
	SignatureSize uint32

	// SignatureFlags contains signature-specific metadata
	// Specification: package_file_format.md: 8.1 Signature Structure
	SignatureFlags uint32

	// SignatureTimestamp is the signature creation time (Unix nanoseconds)
	// Specification: package_file_format.md: 8.1 Signature Structure
	SignatureTimestamp uint32

	// CommentLength is the length of signature comment
	// Specification: package_file_format.md: 8.1 Signature Structure
	CommentLength uint16

	// SignatureComment is a human-readable comment about the signature
	// Specification: package_file_format.md: 8.1 Signature Structure
	SignatureComment string

	// SignatureData contains the raw signature bytes
	// Specification: package_file_format.md: 8.1 Signature Structure
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
func (s *Signature) validate() error {
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
func (s *Signature) size() int {
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
// readSignatureHeader reads the 18-byte fixed header into s; returns bytes read and error.
func readSignatureHeader(r io.Reader, s *Signature) (int64, error) {
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureType); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeCorruption, "failed to read signature type: incomplete data", pkgerrors.ValidationErrorContext{
				Field: "SignatureType", Value: int64(0), Expected: "4 bytes",
			})
		}
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature type", pkgerrors.ValidationErrorContext{
			Field: "SignatureType", Value: nil, Expected: "4 bytes",
		})
	}
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureSize); err != nil {
		return 4, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature size", pkgerrors.ValidationErrorContext{
			Field: "SignatureSize", Value: nil, Expected: "4 bytes",
		})
	}
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureFlags); err != nil {
		return 8, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature flags", pkgerrors.ValidationErrorContext{
			Field: "SignatureFlags", Value: nil, Expected: "4 bytes",
		})
	}
	if err := binary.Read(r, binary.LittleEndian, &s.SignatureTimestamp); err != nil {
		return 12, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature timestamp", pkgerrors.ValidationErrorContext{
			Field: "SignatureTimestamp", Value: nil, Expected: "4 bytes",
		})
	}
	if err := binary.Read(r, binary.LittleEndian, &s.CommentLength); err != nil {
		return 16, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read comment length", pkgerrors.ValidationErrorContext{
			Field: "CommentLength", Value: nil, Expected: "2 bytes",
		})
	}
	return 18, nil
}

// readSignatureComment reads CommentLength bytes into s.SignatureComment; returns bytes read and error.
func readSignatureComment(r io.Reader, s *Signature) (int64, error) {
	if s.CommentLength == 0 {
		s.SignatureComment = ""
		return 0, nil
	}
	commentBytes := make([]byte, s.CommentLength)
	n, err := io.ReadFull(r, commentBytes)
	if err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature comment", pkgerrors.ValidationErrorContext{
			Field: "SignatureComment", Value: s.CommentLength, Expected: "comment data",
		})
	}
	if uint16(n) != s.CommentLength {
		return int64(n), pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete comment read", nil, pkgerrors.ValidationErrorContext{
			Field: "SignatureComment", Value: n, Expected: fmt.Sprintf("%d bytes", s.CommentLength),
		})
	}
	s.SignatureComment = string(commentBytes)
	return int64(n), nil
}

// readSignatureData reads SignatureSize bytes into s.SignatureData; returns bytes read and error.
func readSignatureData(r io.Reader, s *Signature) (int64, error) {
	if s.SignatureSize == 0 {
		s.SignatureData = nil
		return 0, nil
	}
	signatureData := make([]byte, s.SignatureSize)
	n, err := io.ReadFull(r, signatureData)
	if err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read signature data", pkgerrors.ValidationErrorContext{
			Field: "SignatureData", Value: s.SignatureSize, Expected: "signature data",
		})
	}
	if uint32(n) != s.SignatureSize {
		return int64(n), pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete signature data read", nil, pkgerrors.ValidationErrorContext{
			Field: "SignatureData", Value: n, Expected: fmt.Sprintf("%d bytes", s.SignatureSize),
		})
	}
	s.SignatureData = signatureData
	return int64(n), nil
}

// Returns the number of bytes read and any error encountered.
//
// Specification: package_file_format.md: 8.1 Signature Structure
func (s *Signature) readFrom(r io.Reader) (int64, error) {
	n, err := readSignatureHeader(r, s)
	if err != nil {
		return n, err
	}
	totalRead := n
	n, err = readSignatureComment(r, s)
	if err != nil {
		return totalRead, err
	}
	totalRead += n
	n, err = readSignatureData(r, s)
	if err != nil {
		return totalRead, err
	}
	return totalRead + n, nil
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
// writeSignatureHeader writes the 18-byte fixed header; s.CommentLength and s.SignatureSize must be set.
func writeSignatureHeader(w io.Writer, s *Signature) (int64, error) {
	if err := binary.Write(w, binary.LittleEndian, s.SignatureType); err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature type", pkgerrors.ValidationErrorContext{
			Field: "SignatureType", Value: s.SignatureType, Expected: "written successfully",
		})
	}
	if err := binary.Write(w, binary.LittleEndian, s.SignatureSize); err != nil {
		return 4, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature size", pkgerrors.ValidationErrorContext{
			Field: "SignatureSize", Value: s.SignatureSize, Expected: "written successfully",
		})
	}
	if err := binary.Write(w, binary.LittleEndian, s.SignatureFlags); err != nil {
		return 8, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature flags", pkgerrors.ValidationErrorContext{
			Field: "SignatureFlags", Value: s.SignatureFlags, Expected: "written successfully",
		})
	}
	if err := binary.Write(w, binary.LittleEndian, s.SignatureTimestamp); err != nil {
		return 12, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature timestamp", pkgerrors.ValidationErrorContext{
			Field: "SignatureTimestamp", Value: s.SignatureTimestamp, Expected: "written successfully",
		})
	}
	if err := binary.Write(w, binary.LittleEndian, s.CommentLength); err != nil {
		return 16, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write comment length", pkgerrors.ValidationErrorContext{
			Field: "CommentLength", Value: s.CommentLength, Expected: "written successfully",
		})
	}
	return 18, nil
}

// writeSignatureComment writes s.SignatureComment (s.CommentLength bytes); returns bytes written and error.
func writeSignatureComment(w io.Writer, s *Signature) (int64, error) {
	if s.CommentLength == 0 {
		return 0, nil
	}
	commentBytes := []byte(s.SignatureComment)
	if uint16(len(commentBytes)) != s.CommentLength {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "comment length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field: "CommentLength", Value: len(commentBytes), Expected: fmt.Sprintf("%d", s.CommentLength),
		})
	}
	n, err := w.Write(commentBytes)
	if err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature comment", pkgerrors.ValidationErrorContext{
			Field: "SignatureComment", Value: s.SignatureComment, Expected: "written successfully",
		})
	}
	if uint16(n) != s.CommentLength {
		return int64(n), pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete comment write", nil, pkgerrors.ValidationErrorContext{
			Field: "SignatureComment", Value: n, Expected: fmt.Sprintf("%d bytes", s.CommentLength),
		})
	}
	return int64(n), nil
}

// writeSignatureData writes s.SignatureData (s.SignatureSize bytes); returns bytes written and error.
func writeSignatureData(w io.Writer, s *Signature) (int64, error) {
	if s.SignatureSize == 0 {
		return 0, nil
	}
	if uint32(len(s.SignatureData)) != s.SignatureSize {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "signature size mismatch", nil, pkgerrors.ValidationErrorContext{
			Field: "SignatureSize", Value: len(s.SignatureData), Expected: fmt.Sprintf("%d", s.SignatureSize),
		})
	}
	n, err := w.Write(s.SignatureData)
	if err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write signature data", pkgerrors.ValidationErrorContext{
			Field: "SignatureData", Value: s.SignatureData, Expected: "written successfully",
		})
	}
	if uint32(n) != s.SignatureSize {
		return int64(n), pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete signature data write", nil, pkgerrors.ValidationErrorContext{
			Field: "SignatureData", Value: n, Expected: fmt.Sprintf("%d bytes", s.SignatureSize),
		})
	}
	return int64(n), nil
}

// Specification: package_file_format.md: 8.1 Signature Structure
func (s *Signature) writeTo(w io.Writer) (int64, error) {
	s.CommentLength = uint16(len(s.SignatureComment))
	s.SignatureSize = uint32(len(s.SignatureData))
	n, err := writeSignatureHeader(w, s)
	if err != nil {
		return n, err
	}
	totalWritten := n
	n, err = writeSignatureComment(w, s)
	if err != nil {
		return totalWritten, err
	}
	totalWritten += n
	n, err = writeSignatureData(w, s)
	if err != nil {
		return totalWritten, err
	}
	return totalWritten + n, nil
}

// HasFlag checks if a specific signature flag is set.
//
// Specification: package_file_format.md: 8.2.2 SignatureFlags Field
func (s *Signature) hasFlag(flag uint32) bool {
	return (s.SignatureFlags & flag) != 0
}

// SetFlag sets a specific signature flag.
//
// Specification: package_file_format.md: 8.2.2 SignatureFlags Field
func (s *Signature) setFlag(flag uint32) {
	s.SignatureFlags |= flag
}

// ClearFlag clears a specific signature flag.
//
// Specification: package_file_format.md: 8.2.2 SignatureFlags Field
func (s *Signature) clearFlag(flag uint32) {
	s.SignatureFlags &= ^flag
}
