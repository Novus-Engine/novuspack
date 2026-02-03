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
func (h *HashEntry) validate() error {
	return validateSliceLength(len(h.HashData), h.HashLength, "HashData", "hash data cannot be nil or empty", "non-empty hash data")
}

// Size returns the total size of the HashEntry in bytes.
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data
func (h HashEntry) size() int {
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
func (h *HashEntry) readFrom(r io.Reader) (int64, error) {
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
	hashData, n, err := readLengthPrefixedBytes(r, hashLength, "HashData", "hash data")
	if err != nil {
		return totalRead, err
	}
	totalRead += n
	h.HashData = hashData

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
func (h *HashEntry) writeTo(w io.Writer) (int64, error) {
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
	n, err := writeLengthPrefixedBytes(w, h.HashData, h.HashLength, "HashData")
	if err != nil {
		return totalWritten, err
	}
	totalWritten += n

	return totalWritten, nil
}
