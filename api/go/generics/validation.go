// This file implements generic validation patterns: Validator[T] interface
// and ValidationRule[T]. It contains type-safe validation functions including
// ValidateWith, ValidateAll, and ComposeValidators. This file should contain
// all code related to generic validation as specified in api_generics.md
// Section 1.7 and Section 2.2.
//
// Specification: api_generics.md: 1. Core Generic Types

// Package generics provides generic types and patterns for the NovusPack API.
package generics

import (
	"context"
	"fmt"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// ValidateWith validates a single value using a validator.
//
// ValidateWith[T] applies the given validator to the value and returns
// any validation error that occurs. Errors are returned as *PackageError
// with ValidationErrorContext for structured error handling.
//
// Type Parameters:
//   - T: The type of value to validate
//
// Example:
//
//	validator := &ValidationRule[string]{
//	    Predicate: func(s string) bool { return len(s) > 0 },
//	    Message:   "string cannot be empty",
//	}
//	err := ValidateWith(ctx, "", validator)
func ValidateWith[T any](ctx context.Context, value T, validator Validator[T]) error {
	// Check context cancellation before validation
	if ctx != nil {
		select {
		case <-ctx.Done():
			return pkgerrors.WrapErrorWithContext(ctx.Err(), pkgerrors.ErrTypeContext, "validation cancelled", pkgerrors.ValidationErrorContext{
				Field:    "value",
				Value:    value,
				Expected: "validation completed",
			})
		default:
		}
	}

	if validator == nil {
		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "validator is nil", nil, pkgerrors.ValidationErrorContext{
			Field:    "validator",
			Value:    nil,
			Expected: "non-nil validator",
		})
	}

	err := validator.Validate(value)
	if err == nil {
		return nil
	}

	// If error is already a PackageError, return it
	if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
		return pkgErr
	}

	// Wrap error with ValidationErrorContext
	return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, "validation failed", pkgerrors.ValidationErrorContext{
		Field:    "value",
		Value:    value,
		Expected: "valid value",
	})
}

// ValidateAll validates multiple values using a validator.
//
// ValidateAll[T] applies the given validator to each value and returns
// a slice of errors, one for each validation failure. If all validations
// pass, returns an empty slice. Errors are returned as *PackageError
// with ValidationErrorContext for structured error handling.
//
// Type Parameters:
//   - T: The type of values to validate
//
// Example:
//
//	validator := &ValidationRule[int]{
//	    Predicate: func(n int) bool { return n > 0 },
//	    Message:   "value must be positive",
//	}
//	errors := ValidateAll(ctx, []int{1, -1, 2, -2}, validator)
func ValidateAll[T any](ctx context.Context, values []T, validator Validator[T]) []error {
	if validator == nil {
		return validateAllNilValidator(values)
	}
	var validationErrors []error
	for i, value := range values {
		if ctxErr := validateAllCheckContext(ctx, i, value); ctxErr != nil {
			for j := i; j < len(values); j++ {
				validationErrors = append(validationErrors, ctxErr)
			}
			return validationErrors
		}
		if err := validator.Validate(value); err != nil {
			validationErrors = append(validationErrors, validateAllWrapError(err, i, value))
		}
	}
	return validationErrors
}

func validateAllNilValidator[T any](values []T) []error {
	err := pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "validator is nil", nil, pkgerrors.ValidationErrorContext{
		Field: "validator", Value: nil, Expected: "non-nil validator",
	})
	out := make([]error, len(values))
	for i := range out {
		out[i] = err
	}
	return out
}

func validateAllCheckContext[T any](ctx context.Context, i int, value T) error {
	if ctx == nil {
		return nil
	}
	select {
	case <-ctx.Done():
		return pkgerrors.WrapErrorWithContext(ctx.Err(), pkgerrors.ErrTypeContext, "validation cancelled", pkgerrors.ValidationErrorContext{
			Field: fmt.Sprintf("values[%d]", i), Value: value, Expected: "validation completed",
		})
	default:
		return nil
	}
}

func validateAllWrapError[T any](err error, i int, value T) error {
	if pkgErr, ok := pkgerrors.IsPackageError(err); ok {
		return pkgErr
	}
	return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation, fmt.Sprintf("validation failed at index %d", i), pkgerrors.ValidationErrorContext{
		Field: fmt.Sprintf("values[%d]", i), Value: value, Expected: "valid value",
	})
}

// ComposeValidators creates a validator that runs multiple validators.
//
// ComposeValidators[T] creates a composite validator that runs all provided
// validators in sequence. If any validator fails, the composite validator
// returns that error immediately (short-circuit evaluation).
//
// Type Parameters:
//   - T: The type of value to validate
//
// Example:
//
//	validator1 := &ValidationRule[string]{
//	    Predicate: func(s string) bool { return len(s) > 0 },
//	    Message:   "string cannot be empty",
//	}
//	validator2 := &ValidationRule[string]{
//	    Predicate: func(s string) bool { return len(s) < 100 },
//	    Message:   "string too long",
//	}
//	composite := ComposeValidators(validator1, validator2)
//	err := composite.Validate("test")
type compositeValidator[T any] struct {
	validators []Validator[T]
}

func (c *compositeValidator[T]) Validate(value T) error {
	for _, validator := range c.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

// ComposeValidators creates a validator that runs multiple validators in sequence.
func ComposeValidators[T any](validators ...Validator[T]) Validator[T] {
	// Filter out nil validators
	nonNilValidators := make([]Validator[T], 0, len(validators))
	for _, v := range validators {
		if v != nil {
			nonNilValidators = append(nonNilValidators, v)
		}
	}

	return &compositeValidator[T]{
		validators: nonNilValidators,
	}
}
