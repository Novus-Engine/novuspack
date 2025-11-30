@domain:basic_ops @m2 @REQ-API_BASIC-068 @spec(api_basic_operations.md#831-creating-structured-errors)
Feature: Basic Operations: Creating Structured Errors

  @REQ-API_BASIC-068 @happy
  Scenario: Structured errors are created with error type and message
    Given a package operation that fails
    When structured error is created
    Then error is created with error type
    And error includes message
    And error can wrap underlying cause
    And error structure supports categorization

  @REQ-API_BASIC-068 @happy
  Scenario: Structured errors support additional context
    Given a package operation error
    When structured error is created
    Then error can include path context
    And error can include operation context
    And error can include parameter context
    And context aids in debugging

  @REQ-API_BASIC-068 @happy
  Scenario: Structured errors wrap existing errors
    Given an existing error from operation
    When structured error wraps the error
    Then original error is preserved
    And error chain is maintained
    And error unwrapping is supported

  @REQ-API_BASIC-068 @happy
  Scenario: Structured errors provide type-safe error creation
    Given error creation needs
    When NewPackageError is used
    Then error type is specified
    And error message is provided
    And error follows structured error pattern
