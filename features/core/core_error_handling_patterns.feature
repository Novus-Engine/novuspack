@domain:core @m2 @REQ-CORE-022 @spec(api_core.md#114-error-handling-patterns)
Feature: Core: Error Handling Patterns

  @REQ-CORE-022 @happy
  Scenario: Creating structured errors with context
    Given package operations that may fail
    When structured errors are created
    Then NewPackageError creates new validation errors with context
    And WrapError wraps existing errors with structured information
    And WithContext adds additional context to errors
    And error creation pattern is followed

  @REQ-CORE-022 @happy
  Scenario: Error inspection and handling checks error types
    Given package operations returning errors
    When errors are inspected
    Then IsPackageError checks if error is a PackageError
    And GetErrorType returns error type if error is PackageError
    And error types are checked for appropriate handling
    And switch statements handle different error types

  @REQ-CORE-022 @happy
  Scenario: Error propagation wraps errors with context
    Given package operations that propagate errors
    When errors are propagated
    Then errors are wrapped with additional context
    And error context includes path, operation, and relevant details
    And WrapError wraps errors with structured information
    And error chain is maintained

  @REQ-CORE-022 @happy
  Scenario: Sentinel error compatibility maintains backward compatibility
    Given code using sentinel errors
    When structured errors are used
    Then sentinel errors are still supported and can be wrapped
    And sentinel errors can be converted to structured errors
    And backward compatibility is maintained

  @REQ-CORE-022 @happy
  Scenario: Error logging provides full context for debugging
    Given structured errors with context
    When errors are logged
    Then error logging includes full context information
    And error type, message, and context are logged
    And cause information is included if available
    And logging enables better debugging
