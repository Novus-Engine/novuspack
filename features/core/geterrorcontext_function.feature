@domain:core @m2 @REQ-CORE-170 @spec(api_core.md#10322-geterrorcontext-function)
Feature: GetErrorContext function retrieves type-safe context from errors

  @REQ-CORE-170 @happy
  Scenario: GetErrorContext retrieves type-safe context from an error
    Given an error that may contain context
    When GetErrorContext is called
    Then type-safe context is retrieved when present
    And the return value allows typed access to context fields
    And the behavior matches the GetErrorContext specification
