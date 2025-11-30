@domain:basic_ops @m2 @REQ-API_BASIC-071 @spec(api_basic_operations.md#84-error-handling-best-practices)
Feature: Basic Operations: Error Handling Best Practices (Error Checking)

  @REQ-API_BASIC-071 @happy
  Scenario: Always check for errors after operations
    Given a package operation
    When operation completes
    Then error return value is checked
    And errors are never ignored
    And error checking prevents silent failures

  @REQ-API_BASIC-071 @happy
  Scenario: Use structured errors for better debugging
    Given package operations that may fail
    When structured errors are used
    Then errors provide rich context
    And errors aid in debugging
    And error information is comprehensive

  @REQ-API_BASIC-071 @happy
  Scenario: Use context for cancellation
    Given long-running package operations
    When context timeouts are used
    Then operations can be cancelled
    And timeouts prevent indefinite blocking
    And context cancellation is handled gracefully

  @REQ-API_BASIC-071 @happy
  Scenario: Handle different error types appropriately
    Given package operations with various error types
    When errors are handled
    Then validation errors receive user-friendly messages
    And security errors are logged
    And I/O errors trigger retry logic
    And error handling strategy matches error type

  @REQ-API_BASIC-071 @happy
  Scenario: Clean up resources properly
    Given package operations with resources
    When errors occur
    Then resources are still cleaned up
    And cleanup uses defer statements
    And resource leaks are prevented

  @REQ-API_BASIC-071 @error
  Scenario: Ignoring errors leads to silent failures
    Given a package operation that may fail
    When error return value is ignored
    Then failures go undetected
    And package state may be inconsistent
    And errors should never be ignored
