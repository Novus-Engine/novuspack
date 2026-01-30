@domain:core @m2 @REQ-CORE-023 @spec(api_core.md#1041-creating-structured-errors)
Feature: Core: Creating Structured Errors

  @REQ-CORE-023 @happy
  Scenario: NewPackageError creates new validation errors with context
    Given package operations that may fail with validation errors
    When NewPackageError is called
    Then new validation error is created with ErrTypeValidation type
    And error message is provided
    And error context can be added with WithContext
    And structured error is ready for use

  @REQ-CORE-023 @happy
  Scenario: NewPackageError creates errors with different error types
    Given package operations that may fail with various errors
    When NewPackageError is called with different error types
    Then errors are created with appropriate error types (Validation, IO, Security, etc.)
    And error types categorize errors appropriately
    And error handling can respond based on type

  @REQ-CORE-023 @happy
  Scenario: WrapError wraps existing errors with structured information
    Given existing errors from package operations
    When WrapError is called
    Then existing errors are wrapped with structured information
    And error type is assigned to wrapped error
    And error message is added to wrapped error
    And error chain is maintained through Cause field

  @REQ-CORE-023 @happy
  Scenario: WithContext adds additional context to errors
    Given structured errors
    When WithContext is called with key-value pairs
    Then additional context is added to errors
    And context includes relevant operation details
    And context aids in debugging
    And context information is preserved

  @REQ-CORE-023 @error
  Scenario: Error creation provides structured information for error handling
    Given package operations that fail
    When structured errors are created
    Then structured information is provided for error handling
    And error details include type, message, and context
    And error handling can respond appropriately
