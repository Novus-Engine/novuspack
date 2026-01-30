@domain:basic_ops @m2 @REQ-API_BASIC-200 @spec(api_basic_operations.md#3333-state-dependent-operations)
Feature: State-dependent operations

  @REQ-API_BASIC-200 @happy
  Scenario: Operations require specific lifecycle states
    Given a package lifecycle with distinct states
    When an operation is invoked
    Then the operation requires a specific package state to proceed
    And operations validate the package state before executing
    And operations fail predictably when state requirements are not met
    And state requirements prevent invalid lifecycle usage
    And state-dependent rules align with documented lifecycle states and transitions

