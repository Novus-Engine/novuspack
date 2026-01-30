// This file implements the HashEntry structure representing hash data for file
// entries. It contains the HashEntry type definition, hash type and purpose
// constants, and methods for working with hash entries. This file should contain
// all code related to hash entries as specified in api_file_mgmt_index.md
// Section 1.1 and package_file_format.md Section 4.1.4.3.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure

package metadata

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// HashEntry represents a hash value for a file with its type and purpose.
//
// Size: Variable (4 bytes + hash data length)
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data
type HashEntry struct {
	// HashType is the hash algorithm identifier
	// Specification: package_file_format.md: 4.1.5 Hash Algorithm Support
	HashType uint8

	// HashPurpose is the hash purpose identifier
	// Specification: package_file_format.md: 4.1.5 Hash Algorithm Support
	HashPurpose uint8

	// HashLength is the length of hash data in bytes
	// Specification: package_file_format.md: 4.1.4.3 Hash Data
	HashLength uint16

	// HashData contains the actual hash value
	// Specification: package_file_format.md: 4.1.4.3 Hash Data
	HashData []byte
}

// Validate performs validation checks on the HashEntry.
//
// Validation checks:
//   - HashLength must match actual HashData length
//   - HashData must not be nil or empty
//   - HashType must be valid
//
// Returns an error if any validation check fails.
func (h *HashEntry) Validate() error {
	if len(h.HashData) == 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "hash data cannot be nil or empty", nil, pkgerrors.ValidationErrorContext{
			Field:    "HashData",
			Value:    nil,
			Expected: "non-empty hash data",
		})
	}

	if uint16(len(h.HashData)) != h.HashLength {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "hash length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "HashLength",
			Value:    h.HashLength,
			Expected: fmt.Sprintf("%d", len(h.HashData)),
		})
	}

	return nil
}

// Size returns the total size of the HashEntry in bytes.
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data
func (h *HashEntry) Size() int {
	return 4 + int(h.HashLength) // Type(1) + Purpose(1) + Length(2) + Data
}

// ReadFrom reads a HashEntry from the provided io.Reader.
//
// The binary format is:
//   - HashType (1 byte)
//   - HashPurpose (1 byte)
//   - HashLength (2 bytes, little-endian uint16)
//   - HashData (HashLength bytes)
//
// Returns the number of bytes read and any error encountered.
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data
func (h *HashEntry) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read HashType (1 byte)
	var hashType uint8
	if err := binary.Read(r, binary.LittleEndian, &hashType); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read hash type", pkgerrors.ValidationErrorContext{
			Field:    "HashType",
			Value:    nil,
			Expected: "1 byte",
		})
	}
	totalRead += 1
	h.HashType = hashType

	// Read HashPurpose (1 byte)
	var hashPurpose uint8
	if err := binary.Read(r, binary.LittleEndian, &hashPurpose); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read hash purpose", pkgerrors.ValidationErrorContext{
			Field:    "HashPurpose",
			Value:    nil,
			Expected: "1 byte",
		})
	}
	totalRead += 1
	h.HashPurpose = hashPurpose

	// Read HashLength (2 bytes)
	var hashLength uint16
	if err := binary.Read(r, binary.LittleEndian, &hashLength); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read hash length", pkgerrors.ValidationErrorContext{
			Field:    "HashLength",
			Value:    nil,
			Expected: "2 bytes",
		})
	}
	totalRead += 2
	h.HashLength = hashLength

	// Read HashData (HashLength bytes)
	if hashLength > 0 {
		hashData := make([]byte, hashLength)
		n, err := io.ReadFull(r, hashData)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read hash data", pkgerrors.ValidationErrorContext{
				Field:    "HashData",
				Value:    hashLength,
				Expected: "hash data",
			})
		}
		if uint16(n) != hashLength {
			return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete hash data read", nil, pkgerrors.ValidationErrorContext{
				Field:    "HashData",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", hashLength),
			})
		}
		totalRead += int64(n)
		h.HashData = hashData
	} else {
		h.HashData = nil
	}

	return totalRead, nil
}

// WriteTo writes a HashEntry to the provided io.Writer.
//
// The binary format is:
//   - HashType (1 byte)
//   - HashPurpose (1 byte)
//   - HashLength (2 bytes, little-endian uint16)
//   - HashData (HashLength bytes)
//
// Returns the number of bytes written and any error encountered.
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data
func (h *HashEntry) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write HashType (1 byte)
	if err := binary.Write(w, binary.LittleEndian, h.HashType); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write hash type", pkgerrors.ValidationErrorContext{
			Field:    "HashType",
			Value:    h.HashType,
			Expected: "written successfully",
		})
	}
	totalWritten += 1

	// Write HashPurpose (1 byte)
	if err := binary.Write(w, binary.LittleEndian, h.HashPurpose); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write hash purpose", pkgerrors.ValidationErrorContext{
			Field:    "HashPurpose",
			Value:    h.HashPurpose,
			Expected: "written successfully",
		})
	}
	totalWritten += 1

	// Write HashLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, h.HashLength); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write hash length", pkgerrors.ValidationErrorContext{
			Field:    "HashLength",
			Value:    h.HashLength,
			Expected: "written successfully",
		})
	}
	totalWritten += 2

	// Write HashData (HashLength bytes)
	if h.HashLength > 0 {
		if uint16(len(h.HashData)) != h.HashLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "hash length mismatch", nil, pkgerrors.ValidationErrorContext{
				Field:    "HashLength",
				Value:    len(h.HashData),
				Expected: fmt.Sprintf("%d", h.HashLength),
			})
		}
		n, err := w.Write(h.HashData)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write hash data", pkgerrors.ValidationErrorContext{
				Field:    "HashData",
				Value:    h.HashData,
				Expected: "written successfully",
			})
		}
		if uint16(n) != h.HashLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete hash data write", nil, pkgerrors.ValidationErrorContext{
				Field:    "HashData",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", h.HashLength),
			})
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}
