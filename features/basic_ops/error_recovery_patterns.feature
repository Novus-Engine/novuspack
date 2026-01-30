@domain:basic_ops @m2 @REQ-API_BASIC-217 @spec(api_basic_operations.md#3374-error-recovery)
Feature: Error recovery patterns

  @REQ-API_BASIC-217 @happy
  Scenario: Error recovery defines handling and recovery patterns
    Given an operation that can fail partially
    When an error occurs
    Then error recovery patterns define how the system responds
    And recovery maintains internal consistency of state
    And recovery ensures resources are released appropriately
    And recovery produces structured errors with context
    And recovery behavior is predictable across operations

