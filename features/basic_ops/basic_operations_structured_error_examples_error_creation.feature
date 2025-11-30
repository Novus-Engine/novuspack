@domain:basic_ops @m2 @REQ-API_BASIC-067 @spec(api_basic_operations.md#83-structured-error-examples)
Feature: Basic Operations: Structured Error Examples

  @REQ-API_BASIC-067 @happy
  Scenario: Creating structured errors demonstrates error creation pattern
    Given a package operation that fails
    When structured error is created
    Then error is created with error type
    And error includes message
    And error can wrap underlying cause
    And error supports additional context

  @REQ-API_BASIC-067 @happy
  Scenario: Error inspection demonstrates error checking pattern
    Given a package error
    When error type is checked
    Then error type can be determined
    And error message can be retrieved
    And error context can be accessed
    And error handling is type-safe

  @REQ-API_BASIC-067 @happy
  Scenario: Common error scenarios demonstrate error handling patterns
    Given different package error scenarios
    When structured errors are used
    Then validation errors are handled appropriately
    And I/O errors are handled appropriately
    And security errors are handled appropriately
    And error patterns provide consistent handling

  @REQ-API_BASIC-067 @error
  Scenario: Error examples cover various error conditions
    Given package operations with potential failures
    When error examples are applied
    Then different error types are demonstrated
    And error handling strategies are shown
    And error context usage is illustrated
