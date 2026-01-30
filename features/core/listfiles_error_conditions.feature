@domain:core @m2 @REQ-CORE-089 @spec(api_core.md#1165-listfiles-error-conditions)
Feature: ListFiles error conditions reference common error mapping table

  @REQ-CORE-089 @happy
  Scenario: ListFiles errors use the common error mapping table
    Given an opened package
    When ListFiles encounters an error condition
    Then the error is mapped using the common error mapping table
    And the returned error is structured
    And error types follow the core error mapping rules
