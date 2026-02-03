// This file implements the OptionalData structure representing optional data
// for file entries. It contains the OptionalData type definition, optional
// data type constants, and methods for working with optional data entries.
// This file should contain all code related to optional data as specified in
// api_file_mgmt_file_entry.md Section 1.1 and package_file_format.md Section 4.1.4.4.
//
// Specification: api_file_mgmt_file_entry.md: 1. FileEntry Structure

package metadata

import (
	"encoding/binary"
	"io"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// Optional data type constants
// Specification: package_file_format.md: 4.1.4.4 Optional Data
const (
	OptionalDataTagsData = 0x00 // Per-file tags data
)

// OptionalDataEntry represents rarely-used file attributes.
//
// Size: Variable (3 bytes + data length)
//
// Specification: package_file_format.md: 4.1.4.4 Optional Data
type OptionalDataEntry struct {
	// DataType is the optional data type identifier
	// Specification: package_file_format.md: 4.1.4.4 Optional Data
	DataType uint8

	// DataLength is the length of optional data in bytes
	// Specification: package_file_format.md: 4.1.4.4 Optional Data
	DataLength uint16

	// Data contains the actual optional data
	// Specification: package_file_format.md: 4.1.4.4 Optional Data
	Data []byte
}

// Validate performs validation checks on the OptionalDataEntry.
//
// Validation checks:
//   - DataLength must match actual Data length
//   - Data must not be nil or empty
//
// Returns an error if any validation check fails.
func (o *OptionalDataEntry) validate() error {
	return validateSliceLength(len(o.Data), o.DataLength, "Data", "optional data cannot be nil or empty", "non-empty data")
}

// Size returns the total size of the OptionalDataEntry in bytes.
//
// Specification: package_file_format.md: 4.1.4.4 Optional Data
func (o OptionalDataEntry) size() int {
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
// Specification: package_file_format.md: 4.1.4.4 Optional Data
func (o *OptionalDataEntry) readFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read DataType (1 byte)
	var dataType uint8
	if err := binary.Read(r, binary.LittleEndian, &dataType); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read data type", pkgerrors.ValidationErrorContext{
			Field:    "DataType",
			Value:    nil,
			Expected: "1 byte",
		})
	}
	totalRead += 1
	o.DataType = dataType

	// Read DataLength (2 bytes)
	var dataLength uint16
	if err := binary.Read(r, binary.LittleEndian, &dataLength); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read data length", pkgerrors.ValidationErrorContext{
			Field:    "DataLength",
			Value:    nil,
			Expected: "2 bytes",
		})
	}
	totalRead += 2
	o.DataLength = dataLength

	// Read Data (DataLength bytes)
	data, n, err := readLengthPrefixedBytes(r, dataLength, "Data", "data bytes")
	if err != nil {
		return totalRead, err
	}
	totalRead += n
	o.Data = data

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
// Specification: package_file_format.md: 4.1.4.4 Optional Data
func (o *OptionalDataEntry) writeTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write DataType (1 byte)
	if err := binary.Write(w, binary.LittleEndian, o.DataType); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write data type", pkgerrors.ValidationErrorContext{
			Field:    "DataType",
			Value:    o.DataType,
			Expected: "written successfully",
		})
	}
	totalWritten += 1

	// Write DataLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, o.DataLength); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write data length", pkgerrors.ValidationErrorContext{
			Field:    "DataLength",
			Value:    o.DataLength,
			Expected: "written successfully",
		})
	}
	totalWritten += 2

	// Write Data (DataLength bytes)
	n, err := writeLengthPrefixedBytes(w, o.Data, o.DataLength, "Data")
	if err != nil {
		return totalWritten, err
	}
	totalWritten += n

	return totalWritten, nil
}
