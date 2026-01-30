// This file implements core generic types: Option[T], Result[T], Strategy[T, U],
// and Validator[T]. It contains the type definitions and methods for these
// fundamental generic abstractions. This file should contain the core generic
// type implementations as specified in api_generics.md Section 1.
//
// Specification: api_generics.md: 1. Core Generic Types

// Package generics provides generic types and patterns for the NovusPack API.
//
// This package implements type-safe generic abstractions including Option,
// Result, Strategy, Validator, and concurrency patterns as specified in
// the NovusPack technical specifications.
package generics

// Option provides type-safe optional configuration values.
//
// Option[T] represents an optional value that may or may not be set.
// It provides methods to set, get, and check the presence of a value.
//
// Type Parameters:
//   - T: The type of the optional value
//
// Example:
//
//	var opt Option[string]
//	opt.Set("hello")
//	if opt.IsSet() {
//	    val, _ := opt.Get()
//	    fmt.Println(val)
//	}
//	value := opt.GetOrDefault("default")
type Option[T any] struct {
	value T
	set   bool
}

// Set sets the option value and marks it as set.
func (o *Option[T]) Set(value T) {
	o.value = value
	o.set = true
}

// Get returns the option value and a boolean indicating if it was set.
// If the option is not set, returns the zero value of T and false.
func (o *Option[T]) Get() (T, bool) {
	return o.value, o.set
}

// GetOrDefault returns the option value if set, otherwise returns the default value.
func (o *Option[T]) GetOrDefault(defaultValue T) T {
	if o.set {
		return o.value
	}
	return defaultValue
}

// IsSet returns true if the option value has been set.
func (o *Option[T]) IsSet() bool {
	return o.set
}

// Clear clears the option value and marks it as unset.
func (o *Option[T]) Clear() {
	var zero T
	o.value = zero
	o.set = false
}

// Result provides type-safe error handling for operations that may fail.
//
// Result[T] encapsulates either a successful value or an error.
// It provides methods to check success/failure and unwrap the value or error.
//
// Type Parameters:
//   - T: The type of the successful value
//
// Example:
//
//	result := Ok("success")
//	if result.IsOk() {
//	    value, _ := result.Unwrap()
//	    fmt.Println(value)
//	}
//
//	result := Err[string](errors.New("failed"))
//	if result.IsErr() {
//	    _, err := result.Unwrap()
//	    fmt.Println(err)
//	}
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result with the given value.
func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		err:   nil,
	}
}

// Err creates a failed Result with the given error.
func Err[T any](err error) Result[T] {
	var zero T
	return Result[T]{
		value: zero,
		err:   err,
	}
}

// Unwrap returns the value and error from the Result.
// If the Result is Ok, returns the value and nil error.
// If the Result is Err, returns the zero value of T and the error.
func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// IsOk returns true if the Result represents a successful operation.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result represents a failed operation.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}
