package generics

import (
	"context"
	"errors"
	"testing"
)

// TestStrategy tests the Strategy interface
type testStrategy struct {
	name         string
	strategyType string
}

func (s *testStrategy) Process(ctx context.Context, input string) (string, error) {
	if input == "error" {
		return "", errors.New("processing error")
	}
	return "processed: " + input, nil
}

func (s *testStrategy) Name() string {
	return s.name
}

func (s *testStrategy) Type() string {
	return s.strategyType
}

func TestStrategy(t *testing.T) {
	strategy := &testStrategy{
		name:         "test-strategy",
		strategyType: "test",
	}

	// Test Name
	if strategy.Name() != "test-strategy" {
		t.Errorf("Name should return 'test-strategy', got %q", strategy.Name())
	}

	// Test Type
	if strategy.Type() != "test" {
		t.Errorf("Type should return 'test', got %q", strategy.Type())
	}

	// Test Process - success
	ctx := context.Background()
	result, err := strategy.Process(ctx, "input")
	if err != nil {
		t.Errorf("Process should succeed, got error: %v", err)
	}
	if result != "processed: input" {
		t.Errorf("Process should return processed value, got %q", result)
	}

	// Test Process - error
	_, err = strategy.Process(ctx, "error")
	if err == nil {
		t.Error("Process should return error for 'error' input")
	}
}

// TestValidator tests the Validator interface
type testValidator struct {
	shouldFail bool
}

func (v *testValidator) Validate(value string) error {
	if v.shouldFail {
		return errors.New("validation failed")
	}
	return nil
}

func TestValidator(t *testing.T) {
	// Success case
	validator := &testValidator{shouldFail: false}
	err := validator.Validate("test")
	if err != nil {
		t.Errorf("Validate should succeed, got error: %v", err)
	}

	// Failure case
	validator = &testValidator{shouldFail: true}
	err = validator.Validate("test")
	if err == nil {
		t.Error("Validate should fail when shouldFail is true")
	}
}

// TestValidationRule tests ValidationRule
func TestValidationRule_String(t *testing.T) {
	// Valid case
	rule := &ValidationRule[string]{
		Name:      "non-empty",
		Predicate: func(s string) bool { return len(s) > 0 },
		Message:   "string cannot be empty",
	}

	err := rule.Validate("test")
	if err != nil {
		t.Errorf("ValidationRule should pass for non-empty string, got error: %v", err)
	}

	// Invalid case
	err = rule.Validate("")
	if err == nil {
		t.Error("ValidationRule should fail for empty string")
	}
	expectedErrMsg := "[Validation] string cannot be empty"
	if err.Error() != expectedErrMsg {
		t.Errorf("ValidationRule should return configured message, got %q, want %q", err.Error(), expectedErrMsg)
	}
}

func TestValidationRule_Int(t *testing.T) {
	// Valid case - positive number
	rule := &ValidationRule[int]{
		Name:      "positive",
		Predicate: func(n int) bool { return n > 0 },
		Message:   "number must be positive",
	}

	err := rule.Validate(42)
	if err != nil {
		t.Errorf("ValidationRule should pass for positive number, got error: %v", err)
	}

	// Invalid case - zero
	err = rule.Validate(0)
	if err == nil {
		t.Error("ValidationRule should fail for zero")
	}

	// Invalid case - negative
	err = rule.Validate(-1)
	if err == nil {
		t.Error("ValidationRule should fail for negative number")
	}
}

func TestValidationRule_CustomType(t *testing.T) {
	// Valid case
	rule := &ValidationRule[CustomType]{
		Name: "has-id",
		Predicate: func(ct CustomType) bool {
			return ct.ID > 0
		},
		Message: "ID must be greater than 0",
	}

	err := rule.Validate(CustomType{ID: 1, Name: "test"})
	if err != nil {
		t.Errorf("ValidationRule should pass for valid CustomType, got error: %v", err)
	}

	// Invalid case
	err = rule.Validate(CustomType{ID: 0, Name: "test"})
	if err == nil {
		t.Error("ValidationRule should fail for invalid CustomType")
	}
}

func TestValidationRule_NilPredicate(t *testing.T) {
	rule := &ValidationRule[string]{
		Name:      "nil-predicate",
		Predicate: nil,
		Message:   "test",
	}

	err := rule.Validate("test")
	if err == nil {
		t.Error("ValidationRule should return error when predicate is nil")
	}
	expectedErrMsg := "[Validation] validation rule predicate is nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("ValidationRule should return specific error for nil predicate, got %q, want %q", err.Error(), expectedErrMsg)
	}
}

func TestValidationRule_EmptyMessage(t *testing.T) {
	rule := &ValidationRule[string]{
		Name:      "empty-message",
		Predicate: func(s string) bool { return false },
		Message:   "",
	}

	err := rule.Validate("test")
	if err == nil {
		t.Error("ValidationRule should fail when predicate returns false")
	}
	expectedErrMsg := "[Validation] validation failed"
	if err.Error() != expectedErrMsg {
		t.Errorf("ValidationRule should return default message when Message is empty, got %q, want %q", err.Error(), expectedErrMsg)
	}
}

func TestValidationRule_InterfaceCompatibility(t *testing.T) {
	// Ensure ValidationRule implements Validator interface
	var _ Validator[string] = &ValidationRule[string]{}
	var _ Validator[int] = &ValidationRule[int]{}
	var _ Validator[CustomType] = &ValidationRule[CustomType]{}
}
