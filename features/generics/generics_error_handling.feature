@domain:generics @m2 @REQ-GEN-010 @REQ-GEN-011 @spec(api_generics.md#12-result-type)
Feature: Generics Error Handling

  @REQ-GEN-010 @happy
  Scenario: Result type provides type-safe error handling
    Given a generic Result type
    When Result is used for operations that may fail
    Then type-safe error handling is provided
    And Result wraps both value and error
    And Result indicates success or failure

  @REQ-GEN-010 @happy
  Scenario: Result Ok creates successful result
    Given a value
    When Ok is called with value
    Then Result indicates success
    And Result contains the value
    And Result error is nil

  @REQ-GEN-010 @error
  Scenario: Result Err creates error result
    Given an error
    When Err is called with error
    Then Result indicates failure
    And Result contains the error
    And Result value is zero value

  @REQ-GEN-010 @happy
  Scenario: Result Unwrap retrieves value and error
    Given a Result type
    When Unwrap is called
    Then value is returned if Result is Ok
    And error is returned if Result is Err
    And type-safe unwrapping is provided

  @REQ-GEN-010 @happy
  Scenario: Result IsOk checks for success
    Given a Result type
    When IsOk is called
    Then boolean indicates if Result is Ok
    And success state can be checked
    And type-safe state checking is provided

  @REQ-GEN-010 @error
  Scenario: Result IsErr checks for failure
    Given a Result type
    When IsErr is called
    Then boolean indicates if Result is Err
    And error state can be checked
    And type-safe error checking is provided

  @REQ-GEN-011 @happy
  Scenario: Result type usage demonstrates successful operation handling
    Given an operation that returns Result
    When operation succeeds
    Then Result IsOk returns true
    And Result Unwrap returns value and nil error
    And success path is handled type-safely

  @REQ-GEN-011 @error
  Scenario: Result type usage demonstrates error operation handling
    Given an operation that returns Result
    When operation fails
    Then Result IsErr returns true
    And Result Unwrap returns zero value and error
    And error path is handled type-safely
