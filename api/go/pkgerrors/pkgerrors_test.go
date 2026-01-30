// Package pkgerrors provides error handling domain structures for the NovusPack implementation.
//
// This file contains unit tests for error handling functions and types.
package pkgerrors

import (
	"errors"
	"testing"

	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestErrorType_String tests the ErrorType String method.
func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected string
	}{
		{
			name:     "Validation",
			errType:  ErrTypeValidation,
			expected: "Validation",
		},
		{
			name:     "IO",
			errType:  ErrTypeIO,
			expected: "IO",
		},
		{
			name:     "Security",
			errType:  ErrTypeSecurity,
			expected: "Security",
		},
		{
			name:     "Unsupported",
			errType:  ErrTypeUnsupported,
			expected: "Unsupported",
		},
		{
			name:     "Context",
			errType:  ErrTypeContext,
			expected: "Context",
		},
		{
			name:     "Corruption",
			errType:  ErrTypeCorruption,
			expected: "Corruption",
		},
		{
			name:     "Encryption",
			errType:  ErrTypeEncryption,
			expected: "Encryption",
		},
		{
			name:     "Compression",
			errType:  ErrTypeCompression,
			expected: "Compression",
		},
		{
			name:     "Signature",
			errType:  ErrTypeSignature,
			expected: "Signature",
		},
		{
			name:     "Unknown",
			errType:  ErrorType(999),
			expected: "Unknown(999)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errType.String()
			if result != tt.expected {
				t.Errorf("ErrorType.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPackageError_Error tests the PackageError Error method.
func TestPackageError_Error(t *testing.T) {
	tests := []struct {
		name     string
		pkgErr   *PackageError
		contains []string
	}{
		{
			name: "error with cause",
			pkgErr: &PackageError{
				Type:    ErrTypeValidation,
				Message: "test error",
				Cause:   errors.New("underlying error"),
			},
			contains: []string{"[Validation]", "test error", "underlying error"},
		},
		{
			name: "error without cause",
			pkgErr: &PackageError{
				Type:    ErrTypeIO,
				Message: "test error",
				Cause:   nil,
			},
			contains: []string{"[IO]", "test error"},
		},
		{
			name: "error with context",
			pkgErr: &PackageError{
				Type:    ErrTypeSecurity,
				Message: "security error",
				Cause:   nil,
				Context: map[string]interface{}{
					"key": "value",
				},
			},
			contains: []string{"[Security]", "security error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkgErr.Error()
			for _, substr := range tt.contains {
				if !testhelpers.Contains(result, substr) {
					t.Errorf("PackageError.Error() = %v, should contain %v", result, substr)
				}
			}
		})
	}
}

// TestPackageError_Unwrap tests the PackageError Unwrap method.
func TestPackageError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")

	tests := []struct {
		name     string
		pkgErr   *PackageError
		expected error
	}{
		{
			name: "error with cause",
			pkgErr: &PackageError{
				Type:    ErrTypeValidation,
				Message: "test error",
				Cause:   underlyingErr,
			},
			expected: underlyingErr,
		},
		{
			name: "error without cause",
			pkgErr: &PackageError{
				Type:    ErrTypeIO,
				Message: "test error",
				Cause:   nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkgErr.Unwrap()
			if result != tt.expected {
				t.Errorf("PackageError.Unwrap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPackageError_Is tests the PackageError Is method.
func TestPackageError_Is(t *testing.T) {
	targetErr := errors.New("target error")
	otherErr := errors.New("other error")

	tests := []struct {
		name     string
		pkgErr   *PackageError
		target   error
		expected bool
	}{
		{
			name: "error with cause that matches target",
			pkgErr: &PackageError{
				Type:    ErrTypeValidation,
				Message: "test error",
				Cause:   targetErr,
			},
			target:   targetErr,
			expected: true,
		},
		{
			name: "error with cause that doesn't match target",
			pkgErr: &PackageError{
				Type:    ErrTypeIO,
				Message: "test error",
				Cause:   otherErr,
			},
			target:   targetErr,
			expected: false,
		},
		{
			name: "error without cause",
			pkgErr: &PackageError{
				Type:    ErrTypeSecurity,
				Message: "test error",
				Cause:   nil,
			},
			target:   targetErr,
			expected: false,
		},
		{
			name: "error with nil target",
			pkgErr: &PackageError{
				Type:    ErrTypeValidation,
				Message: "test error",
				Cause:   targetErr,
			},
			target:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkgErr.Is(tt.target)
			if result != tt.expected {
				t.Errorf("PackageError.Is() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNewPackageError tests the NewPackageError function.
func TestNewPackageError(t *testing.T) {
	underlyingErr := errors.New("underlying error")

	tests := []struct {
		name     string
		errType  ErrorType
		message  string
		cause    error
		validate func(*PackageError) bool
	}{
		{
			name:    "with cause",
			errType: ErrTypeValidation,
			message: "test error",
			cause:   underlyingErr,
			validate: func(e *PackageError) bool {
				return e.Type == ErrTypeValidation &&
					e.Message == "test error" &&
					e.Cause == underlyingErr &&
					e.Context != nil
			},
		},
		{
			name:    "without cause",
			errType: ErrTypeIO,
			message: "test error",
			cause:   nil,
			validate: func(e *PackageError) bool {
				return e.Type == ErrTypeIO &&
					e.Message == "test error" &&
					e.Cause == nil &&
					e.Context != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPackageError[struct{}](tt.errType, tt.message, tt.cause, struct{}{})
			if result == nil {
				t.Fatal("NewPackageError() returned nil")
			}
			if !tt.validate(result) {
				t.Errorf("NewPackageError() validation failed: %+v", result)
			}
		})
	}
}

// TestPackageError_WithContext tests the WithContext method.
func TestPackageError_WithContext(t *testing.T) {
	pkgErr := NewPackageError[struct{}](ErrTypeValidation, "test error", nil, struct{}{})

	result := pkgErr.WithContext("key1", "value1")
	if result != pkgErr {
		t.Errorf("WithContext() should return the same error instance")
	}
	if pkgErr.Context["key1"] != "value1" {
		t.Errorf("WithContext() failed to add context")
	}

	// Test chaining
	_ = pkgErr.WithContext("key2", 42).WithContext("key3", true)
	if pkgErr.Context["key2"] != 42 {
		t.Errorf("WithContext() failed to add second context")
	}
	if pkgErr.Context["key3"] != true {
		t.Errorf("WithContext() failed to add third context")
	}
}

// TestWrapError tests the WrapError function.
func TestWrapError(t *testing.T) {
	underlyingErr := errors.New("underlying error")

	tests := []struct {
		name     string
		err      error
		errType  ErrorType
		message  string
		validate func(*PackageError) bool
	}{
		{
			name:    "wrap standard error",
			err:     underlyingErr,
			errType: ErrTypeIO,
			message: "wrapped error",
			validate: func(e *PackageError) bool {
				return e.Type == ErrTypeIO &&
					e.Message == "wrapped error" &&
					e.Cause == underlyingErr
			},
		},
		{
			name:    "wrap nil error",
			err:     nil,
			errType: ErrTypeValidation,
			message: "wrapped error",
			validate: func(e *PackageError) bool {
				return e.Type == ErrTypeValidation &&
					e.Message == "wrapped error" &&
					e.Cause == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.errType, tt.message)
			if result == nil {
				t.Fatal("WrapError() returned nil")
			}
			if !tt.validate(result) {
				t.Errorf("WrapError() validation failed: %+v", result)
			}
		})
	}
}

// TestIsPackageError tests the IsPackageError function.
func TestIsPackageError(t *testing.T) {
	standardErr := errors.New("standard error")
	pkgErr := NewPackageError[struct{}](ErrTypeValidation, "package error", nil, struct{}{})
	wrappedErr := WrapError(standardErr, ErrTypeIO, "wrapped")

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "PackageError",
			err:      pkgErr,
			expected: true,
		},
		{
			name:     "wrapped PackageError",
			err:      wrappedErr,
			expected: true,
		},
		{
			name:     "standard error",
			err:      standardErr,
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := IsPackageError(tt.err)
			if ok != tt.expected {
				t.Errorf("IsPackageError() ok = %v, want %v", ok, tt.expected)
			}
			if tt.expected && result == nil {
				t.Errorf("IsPackageError() should return non-nil PackageError when ok is true")
			}
			if !tt.expected && result != nil {
				t.Errorf("IsPackageError() should return nil PackageError when ok is false")
			}
		})
	}
}

// TestAs tests the As function.
func TestAs(t *testing.T) {
	pkgErr := NewPackageError[struct{}](ErrTypeValidation, "test", nil, struct{}{})
	standardErr := errors.New("standard error")

	tests := []struct {
		name     string
		err      error
		setup    func() interface{}
		expected bool
	}{
		{
			name: "PackageError to PackageError",
			err:  pkgErr,
			setup: func() interface{} {
				var target *PackageError
				return &target
			},
			expected: true,
		},
		{
			name: "standard error to PackageError",
			err:  standardErr,
			setup: func() interface{} {
				var target *PackageError
				return &target
			},
			expected: false,
		},
		{
			name: "nil error",
			err:  nil,
			setup: func() interface{} {
				var target *PackageError
				return &target
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.setup()
			result := As(tt.err, target)
			if result != tt.expected {
				t.Errorf("As() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetErrorType tests the GetErrorType function.
func TestGetErrorType(t *testing.T) {
	pkgErr := NewPackageError[struct{}](ErrTypeSecurity, "test", nil, struct{}{})
	standardErr := errors.New("standard error")

	tests := []struct {
		name         string
		err          error
		expectedType ErrorType
		expectedOk   bool
	}{
		{
			name:         "PackageError",
			err:          pkgErr,
			expectedType: ErrTypeSecurity,
			expectedOk:   true,
		},
		{
			name:         "standard error",
			err:          standardErr,
			expectedType: 0,
			expectedOk:   false,
		},
		{
			name:         "nil error",
			err:          nil,
			expectedType: 0,
			expectedOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultType, ok := GetErrorType(tt.err)
			if ok != tt.expectedOk {
				t.Errorf("GetErrorType() ok = %v, want %v", ok, tt.expectedOk)
			}
			if resultType != tt.expectedType {
				t.Errorf("GetErrorType() type = %v, want %v", resultType, tt.expectedType)
			}
		})
	}
}

// TestAddErrorContext tests the AddErrorContext generic function.
func TestAddErrorContext(t *testing.T) {
	pkgErr := NewPackageError[struct{}](ErrTypeValidation, "test", nil, struct{}{})
	standardErr := errors.New("standard error")

	tests := []struct {
		name     string
		err      error
		key      string
		value    interface{}
		validate func(error) bool
	}{
		{
			name:  "PackageError with string context",
			err:   pkgErr,
			key:   "string_key",
			value: "string_value",
			validate: func(e error) bool {
				val, ok := GetErrorContext[string](e, "string_key")
				return ok && val == "string_value"
			},
		},
		{
			name:  "PackageError with int context",
			err:   pkgErr,
			key:   "int_key",
			value: 42,
			validate: func(e error) bool {
				val, ok := GetErrorContext[int](e, "int_key")
				return ok && val == 42
			},
		},
		{
			name:  "standard error wrapped",
			err:   standardErr,
			key:   "key",
			value: "value",
			validate: func(e error) bool {
				// Should be wrapped as PackageError
				_, ok := IsPackageError(e)
				return ok
			},
		},
		{
			name:  "nil error",
			err:   nil,
			key:   "key",
			value: "value",
			validate: func(e error) bool {
				return e == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result error
			switch v := tt.value.(type) {
			case string:
				result = AddErrorContext(tt.err, tt.key, v)
			case int:
				result = AddErrorContext(tt.err, tt.key, v)
			default:
				t.Fatalf("Unsupported value type: %T", v)
			}
			if !tt.validate(result) {
				t.Errorf("AddErrorContext() validation failed")
			}
		})
	}
}

// TestGetErrorContext tests the GetErrorContext generic function.
func TestGetErrorContext(t *testing.T) {
	pkgErr := NewPackageError[struct{}](ErrTypeValidation, "test", nil, struct{}{})
	pkgErr.Context["string_key"] = "string_value"
	pkgErr.Context["int_key"] = 42
	pkgErr.Context["wrong_type"] = "not_an_int"

	tests := []struct {
		name        string
		err         error
		key         string
		expectedVal interface{}
		expectedOk  bool
		valueType   string
	}{
		{
			name:        "string context exists",
			err:         pkgErr,
			key:         "string_key",
			expectedVal: "string_value",
			expectedOk:  true,
			valueType:   "string",
		},
		{
			name:        "int context exists",
			err:         pkgErr,
			key:         "int_key",
			expectedVal: 42,
			expectedOk:  true,
			valueType:   "int",
		},
		{
			name:        "key does not exist",
			err:         pkgErr,
			key:         "nonexistent",
			expectedVal: "",
			expectedOk:  false,
			valueType:   "string",
		},
		{
			name:        "wrong type",
			err:         pkgErr,
			key:         "wrong_type",
			expectedVal: 0,
			expectedOk:  false,
			valueType:   "int",
		},
		{
			name:        "nil error",
			err:         nil,
			key:         "key",
			expectedVal: "",
			expectedOk:  false,
			valueType:   "string",
		},
		{
			name:        "standard error",
			err:         errors.New("standard"),
			key:         "key",
			expectedVal: "",
			expectedOk:  false,
			valueType:   "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val interface{}
			var ok bool
			switch tt.valueType {
			case "string":
				val, ok = GetErrorContext[string](tt.err, tt.key)
			case "int":
				val, ok = GetErrorContext[int](tt.err, tt.key)
			}
			if ok != tt.expectedOk {
				t.Errorf("GetErrorContext() ok = %v, want %v", ok, tt.expectedOk)
			}
			if val != tt.expectedVal {
				t.Errorf("GetErrorContext() val = %v, want %v", val, tt.expectedVal)
			}
		})
	}
}

// TestNewTypedPackageError tests the NewTypedPackageError function.
func TestNewTypedPackageError(t *testing.T) {
	type TestContext struct {
		Field1 string
		Field2 int
	}

	ctx := TestContext{Field1: "test", Field2: 42}
	underlyingErr := errors.New("underlying")

	tests := []struct {
		name     string
		errType  ErrorType
		message  string
		cause    error
		context  TestContext
		validate func(*PackageError) bool
	}{
		{
			name:    "with cause",
			errType: ErrTypeValidation,
			message: "test error",
			cause:   underlyingErr,
			context: ctx,
			validate: func(e *PackageError) bool {
				typedCtx, ok := GetErrorContext[TestContext](e, "_typed_context")
				return e.Type == ErrTypeValidation &&
					e.Message == "test error" &&
					e.Cause == underlyingErr &&
					ok &&
					typedCtx.Field1 == "test" &&
					typedCtx.Field2 == 42
			},
		},
		{
			name:    "without cause",
			errType: ErrTypeIO,
			message: "test error",
			cause:   nil,
			context: ctx,
			validate: func(e *PackageError) bool {
				typedCtx, ok := GetErrorContext[TestContext](e, "_typed_context")
				return e.Type == ErrTypeIO &&
					e.Message == "test error" &&
					e.Cause == nil &&
					ok &&
					typedCtx.Field1 == "test" &&
					typedCtx.Field2 == 42
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewTypedPackageError(tt.errType, tt.message, tt.cause, tt.context)
			if result == nil {
				t.Fatal("NewTypedPackageError() returned nil")
			}
			if !tt.validate(result) {
				t.Errorf("NewTypedPackageError() validation failed: %+v", result)
			}
		})
	}
}

// TestWrapErrorWithContext tests the WrapErrorWithContext function.
func TestWrapErrorWithContext(t *testing.T) {
	type TestContext struct {
		Value string
	}

	standardErr := errors.New("standard error")
	pkgErr := NewPackageError[struct{}](ErrTypeValidation, "original", nil, struct{}{})
	ctx := TestContext{Value: "test"}

	tests := []struct {
		name     string
		err      error
		errType  ErrorType
		message  string
		context  TestContext
		validate func(*PackageError) bool
	}{
		{
			name:    "wrap standard error",
			err:     standardErr,
			errType: ErrTypeIO,
			message: "wrapped",
			context: ctx,
			validate: func(e *PackageError) bool {
				typedCtx, ok := GetErrorContext[TestContext](e, "_typed_context")
				return e.Type == ErrTypeIO &&
					e.Message == "wrapped" &&
					e.Cause == standardErr &&
					ok &&
					typedCtx.Value == "test"
			},
		},
		{
			name:    "wrap PackageError",
			err:     pkgErr,
			errType: ErrTypeSecurity,
			message: "updated",
			context: ctx,
			validate: func(e *PackageError) bool {
				typedCtx, ok := GetErrorContext[TestContext](e, "_typed_context")
				return e.Type == ErrTypeSecurity &&
					e.Message == "updated" &&
					ok &&
					typedCtx.Value == "test"
			},
		},
		{
			name:    "wrap nil error",
			err:     nil,
			errType: ErrTypeValidation,
			message: "new error",
			context: ctx,
			validate: func(e *PackageError) bool {
				typedCtx, ok := GetErrorContext[TestContext](e, "_typed_context")
				return e.Type == ErrTypeValidation &&
					e.Message == "new error" &&
					e.Cause == nil &&
					ok &&
					typedCtx.Value == "test"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapErrorWithContext(tt.err, tt.errType, tt.message, tt.context)
			if result == nil {
				t.Fatal("WrapErrorWithContext() returned nil")
			}
			if !tt.validate(result) {
				t.Errorf("WrapErrorWithContext() validation failed: %+v", result)
			}
		})
	}
}

// TestMapError tests the MapError function.
func TestMapError(t *testing.T) {
	type SourceContext struct {
		Value int
	}
	type TargetContext struct {
		Value string
	}

	sourceCtx := SourceContext{Value: 42}
	pkgErr := NewTypedPackageError(ErrTypeValidation, "test", nil, sourceCtx)
	standardErr := errors.New("standard error")

	tests := []struct {
		name     string
		err      error
		mapper   func(SourceContext) TargetContext
		validate func(error) bool
	}{
		{
			name: "map PackageError with typed context",
			err:  pkgErr,
			mapper: func(src SourceContext) TargetContext {
				return TargetContext{Value: "mapped"}
			},
			validate: func(e error) bool {
				typedCtx, ok := GetErrorContext[TargetContext](e, "_typed_context")
				return ok && typedCtx.Value == "mapped"
			},
		},
		{
			name: "map error without typed context",
			err:  standardErr,
			mapper: func(src SourceContext) TargetContext {
				return TargetContext{Value: "mapped"}
			},
			validate: func(e error) bool {
				// Should return original error
				return e == standardErr
			},
		},
		{
			name: "map nil error",
			err:  nil,
			mapper: func(src SourceContext) TargetContext {
				return TargetContext{Value: "mapped"}
			},
			validate: func(e error) bool {
				return e == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapError(tt.err, tt.mapper)
			if !tt.validate(result) {
				t.Errorf("MapError() validation failed")
			}
		})
	}
}
