@domain:core @m2 @REQ-CORE-173 @spec(api_core.md#10411-error-context-type-definitions)
Feature: Error context type definitions define typed error context structures

  @REQ-CORE-173 @happy
  Scenario: Error context types define typed structures for context
    Given a need to attach context to errors
    When error context types are used
    Then typed structures are available for context
    And the type definitions match the specification
    And type safety is preserved for context access
