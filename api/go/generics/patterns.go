// This file implements generic function patterns for collection operations,
// factory functions, and utility patterns. It contains helper functions that
// work with generic types. This file should contain generic utility functions
// as specified in api_generics.md Section 2.
//
// Specification: api_generics.md: 1. Core Generic Types

// Package generics provides generic types and patterns for the NovusPack API.
package generics

import (
	"context"
	"fmt"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// Strategy defines a generic strategy pattern for processing different data types.
//
// Strategy[T, U] represents a pluggable behavior pattern where T is the input type
// and U is the output type. Strategies can be used for compression, encryption,
// validation, and other processing operations.
//
// Type Parameters:
//   - T: The input type
//   - U: The output type
//
// Example:
//
//	type CompressionStrategy struct{}
//	func (s *CompressionStrategy) Process(ctx context.Context, input []byte) ([]byte, error) {
//	    // compress input
//	}
//	func (s *CompressionStrategy) Name() string { return "zstd" }
//	func (s *CompressionStrategy) Type() string { return "compression" }
type Strategy[T any, U any] interface {
	// Process processes the input value and returns the output value or an error.
	Process(ctx context.Context, input T) (U, error)

	// Name returns the name of the strategy.
	Name() string

	// Type returns the type/category of the strategy.
	Type() string
}

// Validator defines a generic validation interface for type-safe validation.
//
// Validator[T] represents a validation operation that can validate values of type T.
// Validators can be composed to create complex validation rules.
//
// Type Parameters:
//   - T: The type of value to validate
//
// Example:
//
//	type StringLengthValidator struct {
//	    Min, Max int
//	}
//	func (v *StringLengthValidator) Validate(ctx context.Context, value string) error {
//	    if len(value) < v.Min || len(value) > v.Max {
//	        return errors.New("invalid length")
//	    }
//	    return nil
//	}
type Validator[T any] interface {
	// Validate validates the given value and returns an error if validation fails.
	// This is a pure in-memory operation and does not require context.
	Validate(value T) error
}

// ValidationRule represents a single validation rule with a predicate function.
//
// ValidationRule[T] provides a simple way to create validators from predicate functions.
// It implements the Validator[T] interface.
//
// Type Parameters:
//   - T: The type of value to validate
//
// Example:
//
//	rule := &ValidationRule[string]{
//	    Name:      "non-empty",
//	    Predicate: func(s string) bool { return len(s) > 0 },
//	    Message:   "string cannot be empty",
//	}
//	err := rule.Validate("")
type ValidationRule[T any] struct {
	// Name is the name of the validation rule.
	Name string

	// Predicate is the function that determines if the value is valid.
	// It should return true if the value is valid, false otherwise.
	Predicate func(T) bool

	// Message is the error message to return if validation fails.
	Message string
}

// Validate validates the given value using the predicate function.
//
// Returns a *PackageError with ValidationErrorContext if the predicate returns false.
// This is a pure in-memory operation and does not require context.
//
// Returns:
//   - error: *PackageError with ValidationErrorContext if validation fails, nil otherwise
func (r *ValidationRule[T]) Validate(value T) error {
	if r.Predicate == nil {
		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, "validation rule predicate is nil", nil, pkgerrors.ValidationErrorContext{
			Field:    "Predicate",
			Value:    nil,
			Expected: "non-nil predicate function",
		})
	}

	if !r.Predicate(value) {
		message := r.Message
		if message == "" {
			message = "validation failed"
		}

		expected := "value that satisfies predicate"
		if r.Name != "" {
			expected = fmt.Sprintf("value that satisfies %s", r.Name)
		}

		return pkgerrors.NewTypedPackageError(pkgerrors.ErrTypeValidation, message, nil, pkgerrors.ValidationErrorContext{
			Field:    r.Name,
			Value:    value,
			Expected: expected,
		})
	}

	return nil
}
