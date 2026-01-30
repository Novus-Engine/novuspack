@domain:core @m2 @REQ-CORE-171 @spec(api_core.md#1033-error-context-functions)
Feature: Error context functions provide AddErrorContext for type-safe context

  @REQ-CORE-171 @happy
  Scenario: AddErrorContext adds type-safe context to errors
    Given an error and context to add
    When AddErrorContext is used
    Then type-safe context is added to the error
    And the error chain is preserved
    And the behavior matches the error context functions specification
