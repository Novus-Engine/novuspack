@domain:basic_ops @m1 @REQ-API_BASIC-015 @spec(api_basic_operations.md#8-error-handling)
Feature: Comprehensive error handling for basic operations

  @error
  Scenario: Not found error is structured
    Given a missing package path
    When OpenPackage is invoked
    Then a structured not found error is returned
    And error type is ErrTypeValidation
    And error context includes the path

  @error
  Scenario: Invalid format error is structured
    Given a non NovusPack file path
    When OpenPackage is invoked
    Then a structured invalid format error is returned
    And error type is ErrTypeValidation
    And error indicates format mismatch

  @error
  Scenario: Validation errors have rich context
    Given an invalid file path for package creation
    When Create is called
    Then a structured validation error is returned
    And error type is ErrTypeValidation
    And error context includes path and operation
    And error message is descriptive

  @error
  Scenario: I/O errors are properly wrapped
    Given a package operation that triggers I/O error
    When the operation is performed
    Then a structured I/O error is returned
    And error type is ErrTypeIO
    And underlying error is preserved
    And error context includes file path and operation

  @error
  Scenario: Security errors indicate permission issues
    Given an operation requiring insufficient permissions
    When the operation is performed
    Then a structured security error is returned
    And error type is ErrTypeSecurity
    And error context includes path and user information

  @error
  Scenario: Unsupported errors indicate version mismatch
    Given a package with unsupported version
    When the package is opened
    Then a structured unsupported error is returned
    And error type is ErrTypeUnsupported
    And error indicates version incompatibility

  @error
  Scenario: Context errors indicate cancellation or timeout
    Given a long-running operation
    And a cancelled context
    When the operation is performed
    Then a structured context error is returned
    And error type is ErrTypeContext
    And error indicates cancellation

  @error
  Scenario: Corruption errors indicate data integrity issues
    Given a corrupted package file
    When the package is opened
    Then a structured corruption error is returned
    And error type is ErrTypeCorruption
    And error indicates integrity violation

  @error
  Scenario: Error wrapping preserves error chain
    Given an operation that fails with underlying error
    When the error is wrapped
    Then structured error preserves underlying cause
    And error unwrapping works correctly
    And error chain is accessible

  @error
  Scenario: Error inspection identifies error types
    Given a structured error
    When error is inspected
    Then error type can be determined
    And error context can be accessed
    And error message is available
    And IsPackageError identifies structured errors

  @error
  Scenario: Sentinel errors can be wrapped
    Given a sentinel error
    When error is wrapped with WrapError
    Then structured error is created
    And underlying sentinel error is preserved
    And error.Is() matches sentinel error correctly

  @error
  Scenario: Error context can be added
    Given a structured error
    When WithContext is called with additional context
    Then error context is extended
    And all context keys are accessible
    And error context provides debugging information

  @REQ-API_BASIC-016 @REQ-API_BASIC-017 @error
  Scenario: Path validation errors indicate invalid input
    Given an invalid path (empty or whitespace-only)
    When package operation is called with invalid path
    Then structured validation error is returned
    And error type is ErrTypeValidation
    And error message indicates path issue
    And error context includes invalid path value

  @REQ-API_BASIC-016 @REQ-API_BASIC-019 @error
  Scenario: Context timeout errors are structured
    Given a long-running operation
    And a context with timeout
    When operation exceeds timeout
    Then structured context error is returned
    And error type is ErrTypeContext
    And error indicates timeout occurred
