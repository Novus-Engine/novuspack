// This file provides shared read/write of length-prefixed byte slices used by
// HashEntry and OptionalDataEntry to avoid duplicate I/O logic.
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data

package metadata

import (
	"fmt"
	"io"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// readLengthPrefixedBytes reads length bytes from r into a new slice.
func readLengthPrefixedBytes(r io.Reader, length uint16, fieldName, expectedDesc string) (data []byte, n int64, err error) {
	if length == 0 {
		return nil, 0, nil
	}
	data = make([]byte, length)
	var readN int
	readN, err = io.ReadFull(r, data)
	if err != nil {
		return nil, 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read "+fieldName, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    length,
			Expected: expectedDesc,
		})
	}
	if uint16(readN) != length {
		return nil, 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete "+fieldName+" read", nil, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    readN,
			Expected: fmt.Sprintf("%d bytes", length),
		})
	}
	return data, int64(readN), nil
}

// writeLengthPrefixedBytes validates length matches data then writes data to w.
func writeLengthPrefixedBytes(w io.Writer, data []byte, length uint16, fieldName string) (int64, error) {
	if length == 0 {
		return 0, nil
	}
	if uint16(len(data)) != length {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    len(data),
			Expected: fmt.Sprintf("%d", length),
		})
	}
	n, err := w.Write(data)
	if err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write "+fieldName, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    data,
			Expected: "written successfully",
		})
	}
	if uint16(n) != length {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete "+fieldName+" write", nil, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    n,
			Expected: fmt.Sprintf("%d bytes", length),
		})
	}
	return int64(n), nil
}
