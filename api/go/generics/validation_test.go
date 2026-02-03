package generics

import (
	"context"
	"testing"
)

// TestValidateWith tests ValidateWith function
func TestValidateWith(t *testing.T) {
	ctx := context.Background()

	// Success case
	validator := &testValidator{shouldFail: false}
	err := ValidateWith(ctx, "test", validator)
	if err != nil {
		t.Errorf("ValidateWith should succeed, got error: %v", err)
	}

	// Failure case
	validator = &testValidator{shouldFail: true}
	err = ValidateWith(ctx, "test", validator)
	if err == nil {
		t.Error("ValidateWith should fail when validator fails")
	}

	// Nil validator
	err = ValidateWith(ctx, "test", nil)
	if err == nil {
		t.Error("ValidateWith should return error for nil validator")
	}
	expectedErrMsg := "[Validation] validator is nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("ValidateWith should return specific error for nil validator, got %q, want %q", err.Error(), expectedErrMsg)
	}
}

func TestValidateWith_DifferentTypes(t *testing.T) {
	ctx := context.Background()

	// String type
	rule := &ValidationRule[string]{
		Predicate: func(s string) bool { return s != "" },
		Message:   "string cannot be empty",
	}
	err := ValidateWith(ctx, "test", rule)
	if err != nil {
		t.Errorf("ValidateWith should succeed for valid string, got error: %v", err)
	}

	// Int type
	intRule := &ValidationRule[int]{
		Predicate: func(n int) bool { return n > 0 },
		Message:   "number must be positive",
	}
	err = ValidateWith(ctx, 42, intRule)
	if err != nil {
		t.Errorf("ValidateWith should succeed for valid int, got error: %v", err)
	}
}

// TestValidateAll tests ValidateAll function
func TestValidateAll(t *testing.T) {
	ctx := context.Background()

	// All valid
	validator := &testValidator{shouldFail: false}
	values := []string{"test1", "test2", "test3"}
	errors := ValidateAll(ctx, values, validator)
	if len(errors) != 0 {
		t.Errorf("ValidateAll should return no errors for all valid values, got %d errors", len(errors))
	}

	// Some invalid
	validator = &testValidator{shouldFail: true}
	errors = ValidateAll(ctx, values, validator)
	if len(errors) != len(values) {
		t.Errorf("ValidateAll should return error for each invalid value, got %d errors", len(errors))
	}

	// Empty slice
	errors = ValidateAll(ctx, []string{}, validator)
	if len(errors) != 0 {
		t.Errorf("ValidateAll should return no errors for empty slice, got %d errors", len(errors))
	}
}

func TestValidateAll_NilValidator(t *testing.T) {
	ctx := context.Background()
	values := []string{"test1", "test2"}

	errors := ValidateAll(ctx, values, nil)
	if len(errors) != len(values) {
		t.Errorf("ValidateAll should return error for each value when validator is nil, got %d errors", len(errors))
	}
	expectedErrMsg := "[Validation] validator is nil"
	for _, err := range errors {
		if err.Error() != expectedErrMsg {
			t.Errorf("ValidateAll should return specific error for nil validator, got %q, want %q", err.Error(), expectedErrMsg)
		}
	}
}

func TestValidateAll_MixedResults(t *testing.T) {
	ctx := context.Background()

	// Create a validator that fails for specific values
	validator := &ValidationRule[string]{
		Predicate: func(s string) bool { return s != "invalid" },
		Message:   "value is invalid",
	}

	values := []string{"valid1", "invalid", "valid2", "invalid"}
	errors := ValidateAll(ctx, values, validator)

	if len(errors) != 2 {
		t.Errorf("ValidateAll should return 2 errors, got %d", len(errors))
	}
}

// TestComposeValidators tests ComposeValidators function
func TestComposeValidators(t *testing.T) {
	// Create two validators
	validator1 := &ValidationRule[string]{
		Predicate: func(s string) bool { return s != "" },
		Message:   "string cannot be empty",
	}
	validator2 := &ValidationRule[string]{
		Predicate: func(s string) bool { return len(s) < 10 },
		Message:   "string too long",
	}

	// Compose validators
	composite := ComposeValidators(validator1, validator2)

	// Test with valid value
	err := composite.Validate("test")
	if err != nil {
		t.Errorf("ComposeValidators should succeed for valid value, got error: %v", err)
	}

	// Test with value that fails first validator
	err = composite.Validate("")
	if err == nil {
		t.Error("ComposeValidators should fail when first validator fails")
	}
	expectedErrMsg := "[Validation] string cannot be empty"
	if err.Error() != expectedErrMsg {
		t.Errorf("ComposeValidators should return first validator's error, got %q, want %q", err.Error(), expectedErrMsg)
	}

	// Test with value that fails second validator
	err = composite.Validate("this is too long")
	if err == nil {
		t.Error("ComposeValidators should fail when second validator fails")
	}
	expectedErrMsg2 := "[Validation] string too long"
	if err.Error() != expectedErrMsg2 {
		t.Errorf("ComposeValidators should return second validator's error, got %q, want %q", err.Error(), expectedErrMsg2)
	}
}

func TestComposeValidators_Empty(t *testing.T) {
	composite := ComposeValidators[string]()
	err := composite.Validate("test")
	if err != nil {
		t.Errorf("ComposeValidators with no validators should always succeed, got error: %v", err)
	}
}

func TestComposeValidators_NilValidators(t *testing.T) {
	validator := &ValidationRule[string]{
		Predicate: func(s string) bool { return s != "" },
		Message:   "string cannot be empty",
	}

	// Compose with nil validators (should be filtered out)
	composite := ComposeValidators(validator, nil, validator, nil)
	err := composite.Validate("test")
	if err != nil {
		t.Errorf("ComposeValidators should filter out nil validators, got error: %v", err)
	}

	// Should still fail for invalid value
	err = composite.Validate("")
	if err == nil {
		t.Error("ComposeValidators should fail when validator fails")
	}
}

func TestComposeValidators_MultipleTypes(t *testing.T) {
	// String validators
	strValidator1 := &ValidationRule[string]{
		Predicate: func(s string) bool { return s != "" },
		Message:   "string cannot be empty",
	}
	strValidator2 := &ValidationRule[string]{
		Predicate: func(s string) bool { return len(s) < 100 },
		Message:   "string too long",
	}
	strComposite := ComposeValidators(strValidator1, strValidator2)
	err := strComposite.Validate("test")
	if err != nil {
		t.Errorf("ComposeValidators[string] should work, got error: %v", err)
	}

	// Int validators
	intValidator1 := &ValidationRule[int]{
		Predicate: func(n int) bool { return n > 0 },
		Message:   "number must be positive",
	}
	intValidator2 := &ValidationRule[int]{
		Predicate: func(n int) bool { return n < 1000 },
		Message:   "number too large",
	}
	intComposite := ComposeValidators(intValidator1, intValidator2)
	err = intComposite.Validate(42)
	if err != nil {
		t.Errorf("ComposeValidators[int] should work, got error: %v", err)
	}
}

func TestComposeValidators_ShortCircuit(t *testing.T) {
	// Create validators that would both fail
	validator1 := &ValidationRule[string]{
		Predicate: func(s string) bool { return false },
		Message:   "first validator failed",
	}
	validator2 := &ValidationRule[string]{
		Predicate: func(s string) bool { return false },
		Message:   "second validator failed",
	}

	composite := ComposeValidators(validator1, validator2)
	err := composite.Validate("test")

	// Should return error from first validator (short-circuit)
	if err == nil {
		t.Error("ComposeValidators should fail")
	}
	expectedErrMsg := "[Validation] first validator failed"
	if err.Error() != expectedErrMsg {
		t.Errorf("ComposeValidators should short-circuit on first failure, got %q, want %q", err.Error(), expectedErrMsg)
	}
}
