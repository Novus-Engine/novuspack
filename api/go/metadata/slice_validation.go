// This file provides shared validation for slice length used by HashEntry and
// OptionalDataEntry to avoid duplicate validation logic.
//
// Specification: package_file_format.md: 4.1.4.3 Hash Data

package metadata

import (
	"fmt"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// validateSliceLength ensures the slice is non-empty and lengthField matches slice length.
func validateSliceLength(sliceLen int, lengthField uint16, fieldName, emptyErrMsg, emptyExpected string) error {
	if sliceLen == 0 {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, emptyErrMsg, nil, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    nil,
			Expected: emptyExpected,
		})
	}
	if uint16(sliceLen) != lengthField {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    fieldName,
			Value:    lengthField,
			Expected: fmt.Sprintf("%d", sliceLen),
		})
	}
	return nil
}
