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
	"fmt"
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
func (o *OptionalDataEntry) Validate() error {
	if len(o.Data) == 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "optional data cannot be nil or empty", nil, pkgerrors.ValidationErrorContext{
			Field:    "Data",
			Value:    nil,
			Expected: "non-empty data",
		})
	}

	if uint16(len(o.Data)) != o.DataLength {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "data length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "DataLength",
			Value:    o.DataLength,
			Expected: fmt.Sprintf("%d", len(o.Data)),
		})
	}

	return nil
}

// Size returns the total size of the OptionalDataEntry in bytes.
//
// Specification: package_file_format.md: 4.1.4.4 Optional Data
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
// Specification: package_file_format.md: 4.1.4.4 Optional Data
func (o *OptionalDataEntry) ReadFrom(r io.Reader) (int64, error) {
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
	if dataLength > 0 {
		data := make([]byte, dataLength)
		n, err := io.ReadFull(r, data)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read data", pkgerrors.ValidationErrorContext{
				Field:    "Data",
				Value:    dataLength,
				Expected: "data bytes",
			})
		}
		if uint16(n) != dataLength {
			return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete data read", nil, pkgerrors.ValidationErrorContext{
				Field:    "Data",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", dataLength),
			})
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
// Specification: package_file_format.md: 4.1.4.4 Optional Data
func (o *OptionalDataEntry) WriteTo(w io.Writer) (int64, error) {
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
	if o.DataLength > 0 {
		if uint16(len(o.Data)) != o.DataLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "data length mismatch", nil, pkgerrors.ValidationErrorContext{
				Field:    "DataLength",
				Value:    len(o.Data),
				Expected: fmt.Sprintf("%d", o.DataLength),
			})
		}
		n, err := w.Write(o.Data)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write data", pkgerrors.ValidationErrorContext{
				Field:    "Data",
				Value:    o.Data,
				Expected: "written successfully",
			})
		}
		if uint16(n) != o.DataLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete data write", nil, pkgerrors.ValidationErrorContext{
				Field:    "Data",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", o.DataLength),
			})
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}
