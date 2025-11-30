@domain:basic_ops @m2 @REQ-API_BASIC-082 @spec(api_basic_operations.md#92-error-handling)
Feature: Basic Operations: Error Handling Best Practices (Context Wrapping)

  @REQ-API_BASIC-082 @happy
  Scenario: Errors are wrapped with context information
    Given a package operation that fails
    When error is returned
    Then error is wrapped with operation context
    And error includes file paths if applicable
    And error includes operation names
    And error includes parameter values
    And error provides better debugging information

  @REQ-API_BASIC-082 @happy
  Scenario: Specific error types are handled appropriately
    Given package operations that may fail
    When different error types occur
    Then validation errors receive user-friendly messages
    And security errors are logged appropriately
    And I/O errors trigger retry logic
    And error handling strategy matches error type

  @REQ-API_BASIC-082 @happy
  Scenario: Structured error system determines handling strategy
    Given a package error
    When error type is inspected
    Then error category is determined
    And appropriate handling strategy is selected
    And error type guides response approach

  @REQ-API_BASIC-082 @error
  Scenario: Error wrapping preserves original error information
    Given a package operation that fails
    When error is wrapped with context
    Then original error is preserved
    And error chain can be unwrapped
    And full error context is available
