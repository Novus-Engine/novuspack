@domain:core @m2 @REQ-CORE-167 @spec(api_core.md#10312-wraperrorwithcontext-function)
Feature: WrapErrorWithContext function wraps error with type-safe context

  @REQ-CORE-167 @happy
  Scenario: WrapErrorWithContext wraps an error with context
    Given an existing error and context
    When WrapErrorWithContext is called
    Then the error is wrapped with type-safe context
    And the underlying error is preserved for unwrapping
    And the behavior matches the WrapErrorWithContext specification
