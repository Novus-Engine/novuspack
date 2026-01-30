@domain:core @m2 @REQ-CORE-106 @spec(api_core.md#1177-getinfo-error-conditions)
Feature: GetInfo error conditions reference common error mapping table

  @REQ-CORE-106 @happy
  Scenario: GetInfo maps errors using the common error mapping table
    Given an opened package
    When GetInfo is called
    Then errors are mapped using the common error mapping table
    And returned errors are structured
    And error types follow the core error mapping rules
