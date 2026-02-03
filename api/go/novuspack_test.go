// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains unit tests for the wrapper functions in novuspack.go
// that re-export functionality from the generics package.
package novuspack

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/fileformat/testutil"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// =============================================================================
// TEST HELPERS
// =============================================================================

// testValidator is a test implementation of Validator for testing wrapper functions.
type testValidator struct {
	shouldFail bool
}

func (v *testValidator) Validate(value string) error {
	if v.shouldFail {
		return errors.New("validation failed")
	}
	return nil
}

// =============================================================================
// VALIDATEWITH TESTS
// =============================================================================

// TestValidateWith tests the ValidateWith wrapper function.
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

// TestValidateWith_DifferentTypes tests ValidateWith with different types.
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

	// Test failure case
	err = ValidateWith(ctx, "", rule)
	if err == nil {
		t.Error("ValidateWith should fail for invalid string")
	}

	err = ValidateWith(ctx, -1, intRule)
	if err == nil {
		t.Error("ValidateWith should fail for invalid int")
	}
}

// =============================================================================
// VALIDATEALL TESTS
// =============================================================================

// TestValidateAll tests the ValidateAll wrapper function.
func TestValidateAll(t *testing.T) {
	ctx := context.Background()

	// All valid
	validator := &testValidator{shouldFail: false}
	values := []string{"test1", "test2", "test3"}
	errs := ValidateAll(ctx, values, validator)
	if len(errs) != 0 {
		t.Errorf("ValidateAll should return no errors for all valid values, got %d errors", len(errs))
	}

	// Some invalid
	validator = &testValidator{shouldFail: true}
	errs = ValidateAll(ctx, values, validator)
	if len(errs) != len(values) {
		t.Errorf("ValidateAll should return error for each invalid value, got %d errors", len(errs))
	}

	// Empty slice
	errs = ValidateAll(ctx, []string{}, validator)
	if len(errs) != 0 {
		t.Errorf("ValidateAll should return no errors for empty slice, got %d errors", len(errs))
	}
}

// TestValidateAll_NilValidator tests ValidateAll with nil validator.
func TestValidateAll_NilValidator(t *testing.T) {
	ctx := context.Background()
	values := []string{"test1", "test2"}

	errs := ValidateAll(ctx, values, nil)
	if len(errs) != len(values) {
		t.Errorf("ValidateAll should return error for each value when validator is nil, got %d errors", len(errs))
	}
	expectedErrMsg := "[Validation] validator is nil"
	for _, err := range errs {
		if err.Error() != expectedErrMsg {
			t.Errorf("ValidateAll should return specific error for nil validator, got %q, want %q", err.Error(), expectedErrMsg)
		}
	}
}

// TestValidateAll_MixedResults tests ValidateAll with mixed valid/invalid values.
func TestValidateAll_MixedResults(t *testing.T) {
	ctx := context.Background()

	// Create a validator that fails for specific values
	validator := &ValidationRule[string]{
		Predicate: func(s string) bool { return s != "invalid" },
		Message:   "value is invalid",
	}

	values := []string{"valid1", "invalid", "valid2", "invalid"}
	errs := ValidateAll(ctx, values, validator)

	if len(errs) != 2 {
		t.Errorf("ValidateAll should return 2 errors, got %d", len(errs))
	}
}

// TestValidateAll_DifferentTypes tests ValidateAll with different types.
func TestValidateAll_DifferentTypes(t *testing.T) {
	ctx := context.Background()

	// Int type
	intRule := &ValidationRule[int]{
		Predicate: func(n int) bool { return n > 0 },
		Message:   "number must be positive",
	}

	values := []int{1, 2, 3, -1, -2}
	errs := ValidateAll(ctx, values, intRule)

	if len(errs) != 2 {
		t.Errorf("ValidateAll should return 2 errors for invalid ints, got %d", len(errs))
	}
}

// =============================================================================
// COMPOSEVALIDATORS TESTS
// =============================================================================

// TestComposeValidators tests the ComposeValidators wrapper function.
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

// TestComposeValidators_Empty tests ComposeValidators with no validators.
func TestComposeValidators_Empty(t *testing.T) {
	composite := ComposeValidators[string]()
	err := composite.Validate("test")
	if err != nil {
		t.Errorf("ComposeValidators with no validators should always succeed, got error: %v", err)
	}
}

// TestComposeValidators_MultipleTypes tests ComposeValidators with different types.
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

// TestComposeValidators_ShortCircuit tests that ComposeValidators short-circuits.
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

// =============================================================================
// OK TESTS
// =============================================================================

// TestOk tests the Ok wrapper function.
func TestOk(t *testing.T) {
	// String type
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

	// Int type
	resultInt := Ok(42)
	if !resultInt.IsOk() {
		t.Error("Result[int] should be Ok")
	}

	valInt, err := resultInt.Unwrap()
	if err != nil {
		t.Errorf("Unwrap should return nil error for Ok[int], got %v", err)
	}
	if valInt != 42 {
		t.Errorf("Unwrap should return value, got %d", valInt)
	}

	// Zero value
	resultZero := Ok(0)
	if !resultZero.IsOk() {
		t.Error("Result with zero value should be Ok")
	}

	valZero, err := resultZero.Unwrap()
	if err != nil || valZero != 0 {
		t.Error("Result should handle zero value in Ok")
	}
}

// TestOk_DifferentTypes tests Ok with various types.
func TestOk_DifferentTypes(t *testing.T) {
	// Bool type
	resultBool := Ok(true)
	if !resultBool.IsOk() {
		t.Error("Result[bool] should be Ok")
	}
	valBool, _ := resultBool.Unwrap()
	if valBool != true {
		t.Error("Result[bool] should return correct value")
	}

	// Slice type
	resultSlice := Ok([]int{1, 2, 3})
	if !resultSlice.IsOk() {
		t.Error("Result[[]int] should be Ok")
	}
	valSlice, _ := resultSlice.Unwrap()
	if len(valSlice) != 3 {
		t.Error("Result[[]int] should return correct value")
	}
}

// =============================================================================
// ERR TESTS
// =============================================================================

// TestErr tests the Err wrapper function.
func TestErr(t *testing.T) {
	testErr := errors.New("test error")

	// String type
	result := Err[string](testErr)
	if result.IsOk() {
		t.Error("Result should not be Ok")
	}
	if !result.IsErr() {
		t.Error("Result should be Err")
	}

	val, err := result.Unwrap()
	if err == nil {
		t.Error("Unwrap should return error for Err")
	}
	if err != testErr {
		t.Errorf("Unwrap should return set error, got %v", err)
	}
	if val != "" {
		t.Errorf("Unwrap should return zero value for Err, got %q", val)
	}

	// Int type
	resultInt := Err[int](testErr)
	if resultInt.IsOk() {
		t.Error("Result[int] should not be Ok")
	}
	if !resultInt.IsErr() {
		t.Error("Result[int] should be Err")
	}

	valInt, errInt := resultInt.Unwrap()
	if errInt == nil {
		t.Error("Unwrap should return error for Err[int]")
	}
	if valInt != 0 {
		t.Errorf("Unwrap should return zero value for Err[int], got %d", valInt)
	}
}

// TestErr_DifferentTypes tests Err with various types.
func TestErr_DifferentTypes(t *testing.T) {
	testErr := errors.New("test error")

	// Bool type
	resultBool := Err[bool](testErr)
	if !resultBool.IsErr() {
		t.Error("Result[bool] should be Err")
	}
	valBool, errBool := resultBool.Unwrap()
	if errBool == nil {
		t.Error("Unwrap should return error for Err[bool]")
	}
	if valBool != false {
		t.Error("Unwrap should return zero value for Err[bool]")
	}

	// Slice type
	resultSlice := Err[[]int](testErr)
	if !resultSlice.IsErr() {
		t.Error("Result[[]int] should be Err")
	}
	valSlice, errSlice := resultSlice.Unwrap()
	if errSlice == nil {
		t.Error("Unwrap should return error for Err[[]int]")
	}
	if valSlice != nil {
		t.Error("Unwrap should return nil for Err[[]int]")
	}
}

// TestErr_NilError tests Err with nil error.
func TestErr_NilError(t *testing.T) {
	// Err with nil error creates an Ok result (since err == nil means Ok)
	result := Err[string](nil)
	if !result.IsOk() {
		t.Error("Result should be Ok when Err is called with nil error")
	}
	if result.IsErr() {
		t.Error("Result should not be Err when Err is called with nil error")
	}

	val, err := result.Unwrap()
	if err != nil {
		t.Errorf("Unwrap should return nil error for Ok result, got %v", err)
	}
	if val != "" {
		t.Errorf("Unwrap should return zero value, got %q", val)
	}
}

// =============================================================================
// PACKAGE LIFECYCLE WRAPPER TESTS
// =============================================================================

// TestNewPackage tests the NewPackage wrapper function.
func TestNewPackage(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("NewPackage() returned nil package")
	}

	// Verify package is in initial state (not open)
	if pkg.IsOpen() {
		t.Error("New package should not be open")
	}
}

// TestNewBuilder tests the NewBuilder wrapper function.
func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}

	// Verify builder can be used to build a package
	ctx := context.Background()
	pkg, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("builder.Build() failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("builder.Build() returned nil package")
	}
}

// TestNewBuilder_WithOptions tests the NewBuilder with configuration options.
func TestNewBuilder_WithOptions(t *testing.T) {
	ctx := context.Background()

	builder := NewBuilder().
		WithCompression(CompressionZstd).
		WithVendorID(0x12345678).
		WithAppID(0x87654321).
		WithComment("Test package")

	pkg, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("builder.Build() with options failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("builder.Build() returned nil package")
	}

	// Package should still not be open after building
	if pkg.IsOpen() {
		t.Error("Built package should not be open until Create/Open is called")
	}
}

// TestOpenPackage tests the OpenPackage wrapper function.
func TestOpenPackage(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/test.nvpk"

	// Create a test package file
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test opening the package
	pkg, err := OpenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	if pkg == nil {
		t.Fatal("OpenPackage() returned nil package")
	}

	if !pkg.IsOpen() {
		t.Error("Opened package should be in open state")
	}
}

// TestOpenPackage_InvalidPath tests OpenPackage with invalid path.
func TestOpenPackage_InvalidPath(t *testing.T) {
	ctx := context.Background()

	// Test with non-existent file
	_, err := OpenPackage(ctx, "/nonexistent/file.nvpk")
	if err == nil {
		t.Error("OpenPackage() should fail for non-existent file")
	}
}

// TestOpenPackage_CancelledContext tests OpenPackage with cancelled context.
func TestOpenPackage_CancelledContext(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/test.nvpk"
	testutil.CreateTestPackageFile(t, pkgPath)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := OpenPackage(ctx, pkgPath)
	if err == nil {
		t.Error("OpenPackage() should fail with cancelled context")
	}
}

// TestOpenPackageReadOnly tests the OpenPackageReadOnly wrapper function.
func TestOpenPackageReadOnly(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/test.nvpk"

	// Create a test package file
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test opening the package in read-only mode
	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	if pkg == nil {
		t.Fatal("OpenPackageReadOnly() returned nil package")
	}

	if !pkg.IsOpen() {
		t.Error("Opened package should be in open state")
	}

	if !pkg.IsReadOnly() {
		t.Error("Package should be in read-only mode")
	}
}

// TestOpenPackageReadOnly_WriteOperationsFail tests that write operations fail on read-only packages.
func TestOpenPackageReadOnly_WriteOperationsFail(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/test.nvpk"

	testutil.CreateTestPackageFile(t, pkgPath)

	pkg, err := OpenPackageReadOnly(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenPackageReadOnly() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	// Attempt a write operation - should fail
	_, err = pkg.AddFile(ctx, "/tmp/test.txt", nil)
	if err == nil {
		t.Error("AddFile() should fail on read-only package")
	}
}

// TestOpenBrokenPackage tests the OpenBrokenPackage wrapper function.
func TestOpenBrokenPackage(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/broken.nvpk"

	// Create a valid package first
	testutil.CreateTestPackageFile(t, pkgPath)

	// Test opening as broken package (should still work for valid package)
	pkg, err := OpenBrokenPackage(ctx, pkgPath)
	if err != nil {
		t.Fatalf("OpenBrokenPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	if pkg == nil {
		t.Fatal("OpenBrokenPackage() returned nil package")
	}
}

// TestOpenBrokenPackage_InvalidFile tests OpenBrokenPackage with invalid file.
func TestOpenBrokenPackage_InvalidFile(t *testing.T) {
	ctx := context.Background()

	// Test with non-existent file
	_, err := OpenBrokenPackage(ctx, "/nonexistent/file.nvpk")
	if err == nil {
		t.Error("OpenBrokenPackage() should fail for non-existent file")
	}
}

// TestReadHeader tests the ReadHeader wrapper function.
func TestReadHeader(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/test.nvpk"

	testutil.CreateTestPackageFile(t, pkgPath)

	// Open file and read header
	file, err := os.Open(pkgPath)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() { _ = file.Close() }()

	header, err := ReadHeader(ctx, file)
	if err != nil {
		t.Fatalf("ReadHeader() failed: %v", err)
	}

	if header == nil {
		t.Fatal("ReadHeader() returned nil header")
	}

	// Verify header has expected magic number
	if header.Magic != fileformat.NVPKMagic {
		t.Errorf("Header magic = 0x%08X, want 0x%08X", header.Magic, fileformat.NVPKMagic)
	}
}

// TestReadHeader_InvalidReader tests ReadHeader with invalid reader.
func TestReadHeader_InvalidReader(t *testing.T) {
	ctx := context.Background()

	// Create a reader that returns error
	errReader := &testhelpers.ErrorReader{}

	_, err := ReadHeader(ctx, errReader)
	if err == nil {
		t.Error("ReadHeader() should fail with error reader")
	}
}

// TestReadHeaderFromPath tests the ReadHeaderFromPath wrapper function.
func TestReadHeaderFromPath(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	pkgPath := tmpDir + "/test.nvpk"

	testutil.CreateTestPackageFile(t, pkgPath)

	header, err := ReadHeaderFromPath(ctx, pkgPath)
	if err != nil {
		t.Fatalf("ReadHeaderFromPath() failed: %v", err)
	}

	if header == nil {
		t.Fatal("ReadHeaderFromPath() returned nil header")
	}

	// Verify header has expected magic number
	if header.Magic != fileformat.NVPKMagic {
		t.Errorf("Header magic = 0x%08X, want 0x%08X", header.Magic, fileformat.NVPKMagic)
	}
}

// TestReadHeaderFromPath_InvalidPath tests ReadHeaderFromPath with invalid path.
func TestReadHeaderFromPath_InvalidPath(t *testing.T) {
	ctx := context.Background()

	_, err := ReadHeaderFromPath(ctx, "/nonexistent/file.nvpk")
	if err == nil {
		t.Error("ReadHeaderFromPath() should fail for non-existent file")
	}
}

// TestReadHeaderFromPath_EmptyPath tests ReadHeaderFromPath with empty path.
func TestReadHeaderFromPath_EmptyPath(t *testing.T) {
	ctx := context.Background()

	_, err := ReadHeaderFromPath(ctx, "")
	if err == nil {
		t.Error("ReadHeaderFromPath() should fail for empty path")
	}
}
