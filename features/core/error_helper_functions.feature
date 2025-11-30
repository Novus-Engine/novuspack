@domain:core @m2 @REQ-CORE-021 @spec(api_core.md#113-error-helper-functions)
Feature: Error Helper Functions

  @REQ-CORE-021 @happy
  Scenario: NewPackageError creates a new structured error
    Given an error operation
    And an error type
    And a message
    And an optional cause error
    When NewPackageError is called
    Then new structured PackageError is created
    And error has specified type, message, and cause
    And error can have context added

  @REQ-CORE-021 @happy
  Scenario: WrapError wraps an existing error with structured information
    Given an existing error
    And an error type
    And a message
    When WrapError is called
    Then existing error is wrapped with structured information
    And wrapped error is PackageError
    And original error is preserved as cause

  @REQ-CORE-021 @happy
  Scenario: IsPackageError checks if an error is a PackageError
    Given an error
    When IsPackageError is called
    Then function checks if error is PackageError
    And returns PackageError and true if it is
    And returns nil and false if it is not

  @REQ-CORE-021 @happy
  Scenario: GetErrorType returns the error type if error is PackageError
    Given an error
    When GetErrorType is called
    Then function returns error type and true if PackageError
    And returns zero and false if not PackageError
    And error type can be used for error handling

  @REQ-CORE-021 @happy
  Scenario: Error helper functions support type-safe context operations
    Given an error operation
    When error helper functions are used
    Then WithTypedContext adds type-safe context to errors
    And GetTypedContext retrieves type-safe context from errors
    And type safety is maintained for context operations
