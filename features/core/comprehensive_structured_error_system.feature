@domain:core @m1 @REQ-CORE-001 @REQ-CORE-010 @spec(api_core.md#11-structured-error-system)
Feature: Comprehensive structured error system

  @happy
  Scenario: Structured error on invalid operation
    Given an operation precondition is not met
    When the operation is invoked
    Then a structured error with a code is returned
    And error type is categorized correctly

  @happy
  Scenario: Error types categorize errors correctly
    Given the structured error system
    When error types are examined
    Then ErrTypeValidation exists for validation errors
    And ErrTypeIO exists for I/O errors
    And ErrTypeSecurity exists for security errors
    And ErrTypeCorruption exists for corruption errors
    And ErrTypeUnsupported exists for unsupported features
    And ErrTypeContext exists for context errors
    And ErrTypeEncryption exists for encryption errors
    And ErrTypeCompression exists for compression errors
    And ErrTypeSignature exists for signature errors

  @happy
  Scenario: PackageError structure contains all required fields
    Given a PackageError instance
    When error structure is examined
    Then Type field exists for error category
    And Message field exists for human-readable message
    And Cause field exists for underlying error
    And Context field exists for additional context

  @happy
  Scenario: NewPackageError creates structured error
    Given error creation functionality
    When NewPackageError is called with type, message, and cause
    Then a PackageError is created
    And error type is set correctly
    And error message is set correctly
    And underlying cause is preserved
    And context map is initialized

  @happy
  Scenario: WrapError wraps existing error
    Given an existing error
    When WrapError is called
    Then a PackageError is created
    And underlying error is preserved as Cause
    And error type is specified
    And error message is provided

  @happy
  Scenario: IsPackageError identifies structured errors
    Given a PackageError instance
    When IsPackageError is called
    Then true is returned
    And error instance is returned
    And error can be inspected

  @happy
  Scenario: GetErrorType extracts error type
    Given a PackageError instance
    When GetErrorType is called
    Then error type is returned
    And true indicates PackageError
    And false indicates non-PackageError

  @happy
  Scenario: WithContext adds error context
    Given a PackageError instance
    When WithContext is called with key and value
    Then context is added to error
    And context key is accessible
    And context value is accessible

  @happy
  Scenario: Error unwrapping preserves error chain
    Given a wrapped PackageError
    When Unwrap is called
    Then underlying error is returned
    And error chain is preserved
    And errors.Is() works correctly

  @happy
  Scenario: Error implements error interface
    Given a PackageError instance
    When Error() method is called
    Then string representation is returned
    And message includes cause if present

  @happy
  Scenario: Error matching works with Is method
    Given a PackageError wrapping a sentinel error
    When errors.Is is called with target
    Then error matching works correctly
    And underlying sentinel errors are matched
