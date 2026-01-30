@domain:basic_ops @m2 @REQ-API_BASIC-153 @spec(api_basic_operations.md#10311-error-context-type-definitions)
Feature: Error context type definitions

  @REQ-API_BASIC-153 @happy
  Scenario: Typed error context structures are defined for structured errors
    Given structured error handling for package operations
    When error contexts are included in errors
    Then typed error context structures are defined for common operation categories
    And error context types capture operation-specific details
    And error context types are usable for logging and diagnostics
    And context types are stable across API versions
    And error context type definitions are part of the public error model

