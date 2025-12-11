package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"
)

// HashEntry represents a hash value for a file with its type and purpose.
//
// Size: Variable (4 bytes + hash data length)
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3 - Hash Data
type HashEntry struct {
	// HashType is the hash algorithm identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.5
	HashType uint8

	// HashPurpose is the hash purpose identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.5
	HashPurpose uint8

	// HashLength is the length of hash data in bytes
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3
	HashLength uint16

	// HashData contains the actual hash value
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3
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
		return fmt.Errorf("hash data cannot be nil or empty")
	}

	if uint16(len(h.HashData)) != h.HashLength {
		return fmt.Errorf("hash length mismatch: specified %d, actual %d", h.HashLength, len(h.HashData))
	}

	return nil
}

// Size returns the total size of the HashEntry in bytes.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3 - Hash Data
func (h *HashEntry) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read HashType (1 byte)
	var hashType uint8
	if err := binary.Read(r, binary.LittleEndian, &hashType); err != nil {
		return totalRead, fmt.Errorf("failed to read hash type: %w", err)
	}
	totalRead += 1
	h.HashType = hashType

	// Read HashPurpose (1 byte)
	var hashPurpose uint8
	if err := binary.Read(r, binary.LittleEndian, &hashPurpose); err != nil {
		return totalRead, fmt.Errorf("failed to read hash purpose: %w", err)
	}
	totalRead += 1
	h.HashPurpose = hashPurpose

	// Read HashLength (2 bytes)
	var hashLength uint16
	if err := binary.Read(r, binary.LittleEndian, &hashLength); err != nil {
		return totalRead, fmt.Errorf("failed to read hash length: %w", err)
	}
	totalRead += 2
	h.HashLength = hashLength

	// Read HashData (HashLength bytes)
	if hashLength > 0 {
		hashData := make([]byte, hashLength)
		n, err := io.ReadFull(r, hashData)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read hash data: %w", err)
		}
		if uint16(n) != hashLength {
			return totalRead, fmt.Errorf("incomplete hash data read: got %d bytes, expected %d", n, hashLength)
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3 - Hash Data
func (h *HashEntry) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write HashType (1 byte)
	if err := binary.Write(w, binary.LittleEndian, h.HashType); err != nil {
		return totalWritten, fmt.Errorf("failed to write hash type: %w", err)
	}
	totalWritten += 1

	// Write HashPurpose (1 byte)
	if err := binary.Write(w, binary.LittleEndian, h.HashPurpose); err != nil {
		return totalWritten, fmt.Errorf("failed to write hash purpose: %w", err)
	}
	totalWritten += 1

	// Write HashLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, h.HashLength); err != nil {
		return totalWritten, fmt.Errorf("failed to write hash length: %w", err)
	}
	totalWritten += 2

	// Write HashData (HashLength bytes)
	if h.HashLength > 0 {
		if uint16(len(h.HashData)) != h.HashLength {
			return totalWritten, fmt.Errorf("hash length mismatch: specified %d, actual %d", h.HashLength, len(h.HashData))
		}
		n, err := w.Write(h.HashData)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write hash data: %w", err)
		}
		if uint16(n) != h.HashLength {
			return totalWritten, fmt.Errorf("incomplete hash data write: wrote %d bytes, expected %d", n, h.HashLength)
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}
