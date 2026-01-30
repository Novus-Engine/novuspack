// This file implements the structured error system: PackageError type, ErrorType
// constants, and error helper functions. It contains NewPackageError, WrapError,
// AsPackageError, and generic error context functions. This file should contain
// all code related to structured error handling as specified in api_core.md
// Section 10.
//
// Specification: api_core.md: 1. Core Interfaces

// Package pkgerrors provides error handling domain structures for the NovusPack implementation.
//
// This package contains error types and structures for structured error handling
// as specified in api_core.md.
package pkgerrors

import (
	"errors"
	"fmt"
)

// ErrorType represents the type of error that occurred.
//
// ErrorType is used to categorize errors into different classes, allowing
// callers to handle errors differently based on their nature. For example,
// validation errors might be retried with corrected input, while corruption
// errors might require data recovery procedures.
//
// Specification: api_core.md: 10 Structured Error System
type ErrorType int

const (
	// ErrTypeValidation indicates a validation error (invalid input, bad format, etc.)
	// Examples: empty path, invalid magic number, malformed header
	ErrTypeValidation ErrorType = iota

	// ErrTypeIO indicates an I/O error (file not found, permission denied, etc.)
	// Examples: failed to open file, failed to read/write, disk full
	ErrTypeIO

	// ErrTypeSecurity indicates a security error (unauthorized access, signature failure, etc.)
	// Examples: signature verification failed, hash mismatch, permission denied
	ErrTypeSecurity

	// ErrTypeUnsupported indicates an unsupported operation error
	// Examples: unsupported format version, unsupported compression algorithm
	ErrTypeUnsupported

	// ErrTypeContext indicates a context error (cancellation, timeout, deadline exceeded)
	// Examples: context cancelled, context deadline exceeded
	ErrTypeContext

	// ErrTypeCorruption indicates data corruption error (checksum failure, invalid structure)
	// Examples: invalid checksum, corrupted file entry, malformed index
	ErrTypeCorruption

	// ErrTypeEncryption indicates an encryption/decryption error (key management, encryption failures)
	// Examples: encryption failed, decryption failed, key management error
	ErrTypeEncryption

	// ErrTypeCompression indicates a compression/decompression error (algorithm errors, compression failures)
	// Examples: compression failed, decompression failed, algorithm error
	ErrTypeCompression

	// ErrTypeSignature indicates a digital signature error (validation failures, signing errors)
	// Examples: signature verification failed, signing failed, signature validation error
	ErrTypeSignature
)

// String returns the string representation of the ErrorType.
func (e ErrorType) String() string {
	switch e {
	case ErrTypeValidation:
		return "Validation"
	case ErrTypeIO:
		return "IO"
	case ErrTypeSecurity:
		return "Security"
	case ErrTypeUnsupported:
		return "Unsupported"
	case ErrTypeContext:
		return "Context"
	case ErrTypeCorruption:
		return "Corruption"
	case ErrTypeEncryption:
		return "Encryption"
	case ErrTypeCompression:
		return "Compression"
	case ErrTypeSignature:
		return "Signature"
	default:
		return fmt.Sprintf("Unknown(%d)", e)
	}
}

// PackageError represents a structured error with type and context information.
//
// PackageError wraps errors with additional metadata to provide better error
// handling and diagnostics. It supports error wrapping through the Unwrap method,
// allowing use with errors.Is() and errors.As().
//
// Example usage:
//
//	pkg, err := novuspack.OpenPackage(ctx, "file.nvpk")
//	if err != nil {
//	    var pkgErr *errors.PackageError
//	    if errors.As(err, &pkgErr) {
//	        switch pkgErr.Type {
//	        case errors.ErrTypeValidation:
//	            // Handle validation error
//	        case errors.ErrTypeIO:
//	            // Handle I/O error
//	        }
//	    }
//	}
//
// Specification: api_core.md: 10 Structured Error System
type PackageError struct {
	// Type categorizes the error for programmatic handling
	Type ErrorType

	// Message provides a human-readable description
	Message string

	// Cause is the underlying error that was wrapped
	Cause error

	// Context provides additional diagnostic information
	Context map[string]any
}

// Error implements the error interface.
// Returns a formatted string combining the message and underlying cause.
func (e *PackageError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap returns the underlying error for error wrapping support.
// This allows PackageError to work with errors.Is() and errors.As().
func (e *PackageError) Unwrap() error {
	return e.Cause
}

// Is implements error matching for error comparison.
// This allows PackageError to participate in Go's standard error matching patterns.
// If a cause error exists, delegates to errors.Is(e.Cause, target) to check if the cause matches the target error.
// If no cause exists, returns false.
func (e *PackageError) Is(target error) bool {
	if e.Cause != nil {
		return errors.Is(e.Cause, target)
	}
	return false
}

// NewPackageError creates a structured error with type-safe context.
//
// This is the primary function for creating PackageError instances.
// All errors must include typed context for compile-time type safety.
//
// Type Parameters:
//   - T: The type of the context value
//
// Parameters:
//   - errType: The type of error that occurred
//   - message: A human-readable error message
//   - cause: The underlying error (can be nil)
//   - context: The typed context value
//
// Returns:
//   - *PackageError: A new PackageError instance with typed context
//
// Specification: api_core.md: 10 Structured Error System
func NewPackageError[T any](errType ErrorType, message string, cause error, context T) *PackageError {
	return NewTypedPackageError(errType, message, cause, context)
}

// WithContext adds context information to the error.
// Returns the error for method chaining.
func (e *PackageError) WithContext(key string, value any) *PackageError {
	e.Context[key] = value
	return e
}

// WrapError wraps an existing error with structured information.
//
// This provides a convenient way to convert standard errors to PackageError.
// It creates a PackageError with the provided error type and message, using
// the existing error as the cause parameter. For typed context, use WrapErrorWithContext.
//
// Parameters:
//   - err: The error to wrap
//   - errType: The type of error that occurred
//   - message: A human-readable error message
//
// Returns:
//   - *PackageError: A new PackageError that wraps the original error
//
// Specification: api_core.md: 10 Structured Error System
func WrapError(err error, errType ErrorType, message string) *PackageError {
	// Use empty struct as minimal context for backward compatibility
	return NewPackageError(errType, message, err, struct{}{})
}

// IsPackageError checks if an error is a PackageError.
//
// Uses Go's errors.As function to attempt type assertion. If the error is a
// PackageError (or wraps one), returns the PackageError pointer and true.
// If the error is not a PackageError, returns nil and false.
//
// This enables safe error type checking:
//
//	if pkgErr, ok := IsPackageError(err); ok {
//	    switch pkgErr.Type {
//	    case ErrTypeValidation:
//	        // Handle validation error
//	    }
//	}
//
// Parameters:
//   - err: The error to check
//
// Returns:
//   - *PackageError: The PackageError if found, nil otherwise
//   - bool: true if the error is a PackageError, false otherwise
//
// Specification: api_core.md: 10 Structured Error System
func IsPackageError(err error) (*PackageError, bool) {
	if err == nil {
		return nil, false
	}
	var pkgErr *PackageError
	if As(err, &pkgErr) {
		return pkgErr, true
	}
	return nil, false
}

// As is a wrapper around the standard errors.As for PackageError.
// This allows the errors package to be self-contained.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// GetErrorType returns the error type if the error is a PackageError.
//
// Calls IsPackageError to check if the error is a PackageError. If successful,
// returns the Type field from the PackageError and true. If the error is not
// a PackageError, returns 0 (zero value for ErrorType) and false.
//
// This enables error type checking:
//
//	if errType, ok := GetErrorType(err); ok && errType == ErrTypeValidation {
//	    // Handle validation error
//	}
//
// Parameters:
//   - err: The error to inspect
//
// Returns:
//   - ErrorType: The error type if found, 0 otherwise
//   - bool: true if the error is a PackageError, false otherwise
//
// Specification: api_core.md: 10 Structured Error System
func GetErrorType(err error) (ErrorType, bool) {
	pkgErr, ok := IsPackageError(err)
	if !ok {
		return 0, false
	}
	return pkgErr.Type, true
}

// NewTypedPackageError creates a structured error with type-safe context.
//
// Combines error creation with type-safe context in a single operation. The
// context parameter is stored in the error's Context map with a special key
// that preserves type information.
//
// Parameters:
//   - errType: The type of error that occurred
//   - message: A human-readable error message
//   - cause: The underlying error (can be nil)
//   - context: The typed context value
//
// Returns:
//   - *PackageError: A new PackageError instance with typed context
//
// Specification: api_core.md: 10 Structured Error System
func NewTypedPackageError[T any](errType ErrorType, message string, cause error, context T) *PackageError {
	pkgErr := &PackageError{
		Type:    errType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]any),
	}
	// Store typed context in a way that preserves type information
	// Using a special key prefix for typed contexts
	pkgErr.Context["_typed_context"] = context
	return pkgErr
}

// MapError transforms an error with a generic mapper function.
//
// Enables error transformation patterns with type safety. If the error is a
// PackageError, applies the mapper function to transform it. Otherwise, returns
// the original error unchanged.
//
// Parameters:
//   - err: The error to transform
//   - mapper: The transformation function
//
// Returns:
//   - error: The transformed error, or the original error if transformation is not applicable
//
// Specification: api_core.md: 10 Structured Error System
func MapError[T any, U any](err error, mapper func(T) U) error {
	if err == nil {
		return nil
	}
	// Try to extract typed context and apply mapper
	typedCtx, ok := GetErrorContext[T](err, "_typed_context")
	if !ok {
		// If no typed context, return original error
		return err
	}
	// Apply mapper to get new typed context
	newCtx := mapper(typedCtx)
	// Create new error with transformed context
	pkgErr, isPkgErr := IsPackageError(err)
	if !isPkgErr {
		return NewTypedPackageError(ErrTypeValidation, err.Error(), err, newCtx)
	}
	// Copy error and update context
	newErr := *pkgErr
	newErr.Context = make(map[string]any)
	for k, v := range pkgErr.Context {
		newErr.Context[k] = v
	}
	newErr.Context["_typed_context"] = newCtx
	return &newErr
}

// WrapErrorWithContext wraps an error with type-safe context.
//
// Provides a convenient way to wrap errors with typed contextual information.
// If the error is already a PackageError, adds the context to it. Otherwise,
// creates a new PackageError wrapping the error.
//
// Parameters:
//   - err: The error to wrap
//   - errType: The type of error that occurred
//   - message: A human-readable error message
//   - context: The typed context value
//
// Returns:
//   - *PackageError: A PackageError with typed context
//
// Specification: api_core.md: 10 Structured Error System
func WrapErrorWithContext[T any](err error, errType ErrorType, message string, context T) *PackageError {
	if err == nil {
		return NewTypedPackageError(errType, message, nil, context)
	}
	pkgErr, ok := IsPackageError(err)
	if !ok {
		// Create new PackageError with minimal context
		pkgErr = &PackageError{
			Type:    errType,
			Message: message,
			Cause:   err,
			Context: make(map[string]any),
		}
	} else {
		// Update message if provided
		if message != "" {
			pkgErr.Message = message
		}
		// Update type if different
		if pkgErr.Type != errType {
			pkgErr.Type = errType
		}
	}
	pkgErr.Context["_typed_context"] = context
	return pkgErr
}

// AddErrorContext adds type-safe context to errors.
//
// Provides compile-time type safety for error context values. If the error is
// a PackageError, adds the typed context value. If not, wraps the error as
// a PackageError with ErrTypeValidation.
//
// Parameters:
//   - err: The error to add context to
//   - key: The context key
//   - value: The typed context value
//
// Returns:
//   - error: The error with added context (may be a new PackageError if err was not one)
//
// Specification: api_core.md: 10 Structured Error System
func AddErrorContext[T any](err error, key string, value T) error {
	if err == nil {
		return nil
	}
	pkgErr, ok := IsPackageError(err)
	if !ok {
		// Wrap non-PackageError as PackageError with minimal context
		pkgErr = NewPackageError(ErrTypeValidation, err.Error(), err, struct{}{})
	}
	pkgErr.Context[key] = value
	return pkgErr
}

// GetErrorContext retrieves type-safe context from errors.
//
// Enables type-safe access to error context values. If the error is a PackageError
// and the key exists, returns the typed value and true. Otherwise returns the
// zero value and false.
//
// Parameters:
//   - err: The error to retrieve context from
//   - key: The context key
//
// Returns:
//   - T: The typed context value if found, zero value otherwise
//   - bool: true if the context value was found and is of type T, false otherwise
//
// Specification: api_core.md: 10 Structured Error System
func GetErrorContext[T any](err error, key string) (T, bool) {
	var zero T
	if err == nil {
		return zero, false
	}
	pkgErr, ok := IsPackageError(err)
	if !ok {
		return zero, false
	}
	value, exists := pkgErr.Context[key]
	if !exists {
		return zero, false
	}
	typedValue, ok := value.(T)
	if !ok {
		return zero, false
	}
	return typedValue, true
}
