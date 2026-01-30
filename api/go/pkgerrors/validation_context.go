// This file defines typed error context structures for validation operations.
// It contains ValidationErrorContext and other typed context structures used
// with the structured error system. This file should contain only type
// definitions for error context structures.
//
// Specification: api_core.md: 10 Structured Error System

// Package pkgerrors provides error handling domain structures for the NovusPack implementation.
package pkgerrors

// ValidationErrorContext provides type-safe error context for validation operations.
//
// ValidationErrorContext is used with NewPackageError and WrapErrorWithContext
// to provide structured context information when validation fails.
//
// Example usage:
//
//	err := NewPackageError(ErrTypeValidation, "validation failed", nil, ValidationErrorContext{
//	    Field:    "path",
//	    Value:    path,
//	    Expected: "non-empty string",
//	})
//
// Specification: api_generics.md: 3.3.2 Error Context in Validation Functions
type ValidationErrorContext struct {
	// Field is the name of the field that failed validation
	Field string

	// Value is the actual value that failed validation
	Value interface{}

	// Expected describes what was expected (e.g., "non-empty string", "positive integer")
	Expected string
}
