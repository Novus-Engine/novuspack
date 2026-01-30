@domain:core @m2 @REQ-CORE-119 @spec(api_core.md#1193-validate-returns)
Feature: Validate returns define error return for validation failures

  @REQ-CORE-119 @happy
  Scenario: Validate returns an error for validation failures
    Given an opened package
    And the package fails validation
    When Validate is called
    Then an error is returned
    And the error indicates the validation failure
