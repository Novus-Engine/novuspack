@domain:core @m2 @REQ-CORE-175 @spec(api_core.md#10413-basicerrorcontext-structure)
Feature: BasicErrorContext structure defines minimal error context

  @REQ-CORE-175 @happy
  Scenario: BasicErrorContext provides minimal error context
    Given an error that needs minimal context
    When BasicErrorContext is used
    Then a minimal context structure is available
    And the structure matches the BasicErrorContext specification
    And the context is type-safe for basic use
