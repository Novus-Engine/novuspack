@domain:basic_ops @m2 @REQ-API_BASIC-057 @spec(api_basic_operations.md#723-defragment-example-usage)
Feature: Defragment Example Usage

  @REQ-API_BASIC-057 @happy
  Scenario: Defragment example demonstrates defragmentation usage
    Given an open NovusPack package
    And a valid context
    When Defragment is called with context
    And error return value is checked
    Then defragmentation pattern is demonstrated
    And error handling is shown
    And usage example is complete

  @REQ-API_BASIC-057 @happy
  Scenario: Defragment example shows error checking pattern
    Given an open NovusPack package
    And a valid context
    When Defragment is called with context
    And error is checked with if err != nil
    Then error checking pattern is demonstrated
    And proper error handling is shown
    And example follows best practices

  @REQ-API_BASIC-057 @happy
  Scenario: Defragment example follows standard Go error handling
    Given a code example using Defragment
    When example code is examined
    Then Defragment call includes error check
    And error is handled with if err != nil pattern
    And function returns error on failure

  @REQ-API_BASIC-057 @happy
  Scenario: Defragment example uses context parameter correctly
    Given a code example using Defragment
    When example code is examined
    Then context is passed as parameter
    And context supports cancellation and timeout
    And context follows standard Go patterns

  @REQ-API_BASIC-057 @error
  Scenario: Defragment example handles validation errors
    Given an open NovusPack package
    And a valid context
    And package is in read-only mode
    When Defragment is called
    Then validation error is returned
    And error indicates package validation issue
    And error follows structured error format
