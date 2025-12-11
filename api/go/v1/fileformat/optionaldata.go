package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"
)

// OptionalDataEntry represents rarely-used file attributes.
//
// Size: Variable (3 bytes + data length)
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4 - Optional Data
type OptionalDataEntry struct {
	// DataType is the optional data type identifier
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4
	DataType uint8

	// DataLength is the length of optional data in bytes
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4
	DataLength uint16

	// Data contains the actual optional data
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4
	Data []byte
}

// Validate performs validation checks on the OptionalDataEntry.
//
// Validation checks:
//   - DataLength must match actual Data length
//   - Data must not be nil or empty
//
// Returns an error if any validation check fails.
func (o *OptionalDataEntry) Validate() error {
	if len(o.Data) == 0 {
		return fmt.Errorf("optional data cannot be nil or empty")
	}

	if uint16(len(o.Data)) != o.DataLength {
		return fmt.Errorf("data length mismatch: specified %d, actual %d", o.DataLength, len(o.Data))
	}

	return nil
}

// Size returns the total size of the OptionalDataEntry in bytes.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4
func (o *OptionalDataEntry) Size() int {
	return 3 + int(o.DataLength) // Type(1) + Length(2) + Data
}

// ReadFrom reads an OptionalDataEntry from the provided io.Reader.
//
// The binary format is:
//   - DataType (1 byte)
//   - DataLength (2 bytes, little-endian uint16)
//   - Data (DataLength bytes)
//
// Returns the number of bytes read and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4 - Optional Data
func (o *OptionalDataEntry) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read DataType (1 byte)
	var dataType uint8
	if err := binary.Read(r, binary.LittleEndian, &dataType); err != nil {
		return totalRead, fmt.Errorf("failed to read data type: %w", err)
	}
	totalRead += 1
	o.DataType = dataType

	// Read DataLength (2 bytes)
	var dataLength uint16
	if err := binary.Read(r, binary.LittleEndian, &dataLength); err != nil {
		return totalRead, fmt.Errorf("failed to read data length: %w", err)
	}
	totalRead += 2
	o.DataLength = dataLength

	// Read Data (DataLength bytes)
	if dataLength > 0 {
		data := make([]byte, dataLength)
		n, err := io.ReadFull(r, data)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read data: %w", err)
		}
		if uint16(n) != dataLength {
			return totalRead, fmt.Errorf("incomplete data read: got %d bytes, expected %d", n, dataLength)
		}
		totalRead += int64(n)
		o.Data = data
	} else {
		o.Data = nil
	}

	return totalRead, nil
}

// WriteTo writes an OptionalDataEntry to the provided io.Writer.
//
// The binary format is:
//   - DataType (1 byte)
//   - DataLength (2 bytes, little-endian uint16)
//   - Data (DataLength bytes)
//
// Returns the number of bytes written and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4 - Optional Data
func (o *OptionalDataEntry) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write DataType (1 byte)
	if err := binary.Write(w, binary.LittleEndian, o.DataType); err != nil {
		return totalWritten, fmt.Errorf("failed to write data type: %w", err)
	}
	totalWritten += 1

	// Write DataLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, o.DataLength); err != nil {
		return totalWritten, fmt.Errorf("failed to write data length: %w", err)
	}
	totalWritten += 2

	// Write Data (DataLength bytes)
	if o.DataLength > 0 {
		if uint16(len(o.Data)) != o.DataLength {
			return totalWritten, fmt.Errorf("data length mismatch: specified %d, actual %d", o.DataLength, len(o.Data))
		}
		n, err := w.Write(o.Data)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write data: %w", err)
		}
		if uint16(n) != o.DataLength {
			return totalWritten, fmt.Errorf("incomplete data write: wrote %d bytes, expected %d", n, o.DataLength)
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}
