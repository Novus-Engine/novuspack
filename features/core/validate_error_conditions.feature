@domain:core @m2 @REQ-CORE-121 @spec(api_core.md#1195-validate-error-conditions)
Feature: Validate error conditions reference common error mapping table

  @REQ-CORE-121 @happy
  Scenario: Validate errors use the common error mapping table
    Given an opened package that may have validation failures
    When Validate is called and encounters an error
    Then the error is mapped using the common error mapping table
    And the returned error is structured
    And error types follow the core error mapping rules
