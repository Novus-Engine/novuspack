package generics

import (
	"errors"
	"testing"
)

// TestOption_String tests Option with string type
func TestOption_String(t *testing.T) {
	var opt Option[string]

	// Initially not set
	if opt.IsSet() {
		t.Error("Option should not be set initially")
	}

	val, ok := opt.Get()
	if ok {
		t.Error("Get should return false when not set")
	}
	if val != "" {
		t.Errorf("Get should return zero value, got %q", val)
	}

	// Set value
	opt.Set("hello")
	if !opt.IsSet() {
		t.Error("Option should be set after Set")
	}

	val, ok = opt.Get()
	if !ok {
		t.Error("Get should return true when set")
	}
	if val != "hello" {
		t.Errorf("Get should return set value, got %q", val)
	}

	// GetOrDefault
	if opt.GetOrDefault("default") != "hello" {
		t.Error("GetOrDefault should return set value")
	}

	// Clear
	opt.Clear()
	if opt.IsSet() {
		t.Error("Option should not be set after Clear")
	}
	if opt.GetOrDefault("default") != "default" {
		t.Error("GetOrDefault should return default after Clear")
	}
}

// TestOption_Int tests Option with int type
func TestOption_Int(t *testing.T) {
	var opt Option[int]

	// Initially not set
	if opt.IsSet() {
		t.Error("Option should not be set initially")
	}

	val, ok := opt.Get()
	if ok {
		t.Error("Get should return false when not set")
	}
	if val != 0 {
		t.Errorf("Get should return zero value, got %d", val)
	}

	// Set value
	opt.Set(42)
	if !opt.IsSet() {
		t.Error("Option should be set after Set")
	}

	val, ok = opt.Get()
	if !ok {
		t.Error("Get should return true when set")
	}
	if val != 42 {
		t.Errorf("Get should return set value, got %d", val)
	}

	// GetOrDefault
	if opt.GetOrDefault(0) != 42 {
		t.Error("GetOrDefault should return set value")
	}

	// Clear
	opt.Clear()
	if opt.IsSet() {
		t.Error("Option should not be set after Clear")
	}
	if opt.GetOrDefault(100) != 100 {
		t.Error("GetOrDefault should return default after Clear")
	}
}

// CustomType is a custom type for testing
type CustomType struct {
	ID   int
	Name string
}

// TestOption_CustomType tests Option with custom type
func TestOption_CustomType(t *testing.T) {
	var opt Option[CustomType]

	// Initially not set
	if opt.IsSet() {
		t.Error("Option should not be set initially")
	}

	val, ok := opt.Get()
	if ok {
		t.Error("Get should return false when not set")
	}
	if val.ID != 0 || val.Name != "" {
		t.Errorf("Get should return zero value, got %+v", val)
	}

	// Set value
	expected := CustomType{ID: 1, Name: "test"}
	opt.Set(expected)
	if !opt.IsSet() {
		t.Error("Option should be set after Set")
	}

	val, ok = opt.Get()
	if !ok {
		t.Error("Get should return true when set")
	}
	if val.ID != expected.ID || val.Name != expected.Name {
		t.Errorf("Get should return set value, got %+v", val)
	}

	// GetOrDefault
	defaultVal := CustomType{ID: 0, Name: "default"}
	if opt.GetOrDefault(defaultVal).ID != expected.ID {
		t.Error("GetOrDefault should return set value")
	}

	// Clear
	opt.Clear()
	if opt.IsSet() {
		t.Error("Option should not be set after Clear")
	}
	if opt.GetOrDefault(defaultVal).ID != defaultVal.ID {
		t.Error("GetOrDefault should return default after Clear")
	}
}

// TestResult_String tests Result with string type
func TestResult_String(t *testing.T) {
	// Ok result
	result := Ok("success")
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.IsErr() {
		t.Error("Result should not be Err")
	}

	val, err := result.Unwrap()
	if err != nil {
		t.Errorf("Unwrap should return nil error for Ok, got %v", err)
	}
	if val != "success" {
		t.Errorf("Unwrap should return value, got %q", val)
	}

	// Err result
	testErr := errors.New("test error")
	result = Err[string](testErr)
	if result.IsOk() {
		t.Error("Result should not be Ok")
	}
	if !result.IsErr() {
		t.Error("Result should be Err")
	}

	val, err = result.Unwrap()
	if err == nil {
		t.Error("Unwrap should return error for Err")
	}
	if err != testErr {
		t.Errorf("Unwrap should return set error, got %v", err)
	}
	if val != "" {
		t.Errorf("Unwrap should return zero value for Err, got %q", val)
	}
}

// TestResult_Int tests Result with int type
func TestResult_Int(t *testing.T) {
	// Ok result
	result := Ok(42)
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.IsErr() {
		t.Error("Result should not be Err")
	}

	val, err := result.Unwrap()
	if err != nil {
		t.Errorf("Unwrap should return nil error for Ok, got %v", err)
	}
	if val != 42 {
		t.Errorf("Unwrap should return value, got %d", val)
	}

	// Err result
	testErr := errors.New("test error")
	result = Err[int](testErr)
	if result.IsOk() {
		t.Error("Result should not be Ok")
	}
	if !result.IsErr() {
		t.Error("Result should be Err")
	}

	val, err = result.Unwrap()
	if err == nil {
		t.Error("Unwrap should return error for Err")
	}
	if err != testErr {
		t.Errorf("Unwrap should return set error, got %v", err)
	}
	if val != 0 {
		t.Errorf("Unwrap should return zero value for Err, got %d", val)
	}
}

// TestResult_CustomType tests Result with custom type
func TestResult_CustomType(t *testing.T) {
	// Ok result
	expected := CustomType{ID: 1, Name: "test"}
	result := Ok(expected)
	if !result.IsOk() {
		t.Error("Result should be Ok")
	}
	if result.IsErr() {
		t.Error("Result should not be Err")
	}

	val, err := result.Unwrap()
	if err != nil {
		t.Errorf("Unwrap should return nil error for Ok, got %v", err)
	}
	if val.ID != expected.ID || val.Name != expected.Name {
		t.Errorf("Unwrap should return value, got %+v", val)
	}

	// Err result
	testErr := errors.New("test error")
	result = Err[CustomType](testErr)
	if result.IsOk() {
		t.Error("Result should not be Ok")
	}
	if !result.IsErr() {
		t.Error("Result should be Err")
	}

	val, err = result.Unwrap()
	if err == nil {
		t.Error("Unwrap should return error for Err")
	}
	if err != testErr {
		t.Errorf("Unwrap should return set error, got %v", err)
	}
	if val.ID != 0 || val.Name != "" {
		t.Errorf("Unwrap should return zero value for Err, got %+v", val)
	}
}

// TestOption_EdgeCases tests edge cases for Option
func TestOption_EdgeCases(t *testing.T) {
	// Zero value
	var opt Option[string]
	opt.Set("")
	if !opt.IsSet() {
		t.Error("Option should be set even with zero value")
	}
	val, ok := opt.Get()
	if !ok || val != "" {
		t.Error("Option should handle zero value correctly")
	}

	// Multiple Set calls
	opt.Set("first")
	opt.Set("second")
	val, _ = opt.Get()
	if val != "second" {
		t.Error("Multiple Set calls should update value")
	}

	// Clear and Set again
	opt.Clear()
	opt.Set("third")
	if !opt.IsSet() {
		t.Error("Option should be set after Clear and Set")
	}
}

// TestResult_EdgeCases tests edge cases for Result
func TestResult_EdgeCases(t *testing.T) {
	// Nil error
	result := Ok("value")
	if !result.IsOk() {
		t.Error("Result with nil error should be Ok")
	}

	// Zero value in Ok
	resultInt := Ok(0)
	valInt, err := resultInt.Unwrap()
	if err != nil || valInt != 0 {
		t.Error("Result should handle zero value in Ok")
	}

	// Empty string in Ok
	resultStr := Ok("")
	valStr, err := resultStr.Unwrap()
	if err != nil || valStr != "" {
		t.Error("Result should handle empty string in Ok")
	}
}
