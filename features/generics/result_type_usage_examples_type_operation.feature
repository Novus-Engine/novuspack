@domain:generics @m2 @REQ-GEN-011 @spec(api_generics.md#121-result-type-usage-examples)
Feature: Result Type Usage Examples

  @REQ-GEN-011 @happy
  Scenario: Result type usage demonstrates successful operation
    Given an operation that returns Result
    When operation succeeds and returns value
    Then Ok is called with value
    And Result IsOk returns true
    And Result Unwrap returns value and nil error
    And success path is handled correctly

  @REQ-GEN-011 @error
  Scenario: Result type usage demonstrates error operation
    Given an operation that returns Result
    When operation fails and returns error
    Then Err is called with error
    And Result IsErr returns true
    And Result Unwrap returns zero value and error
    And error path is handled correctly

  @REQ-GEN-011 @happy
  Scenario: Result type usage demonstrates checking success state
    Given a Result type
    When IsOk is checked
    Then boolean indicates if Result is Ok
    And conditional logic can be applied based on state
    And type-safe state checking is provided

  @REQ-GEN-011 @error
  Scenario: Result type usage demonstrates checking error state
    Given a Result type
    When IsErr is checked
    Then boolean indicates if Result is Err
    And error handling logic can be applied
    And type-safe error checking is provided

  @REQ-GEN-011 @happy
  Scenario: Result type usage demonstrates error handling patterns
    Given generic Result type
    When Result is used in different contexts
    Then error handling patterns are demonstrated
    And type-safe error handling is shown
    And best practices are illustrated
