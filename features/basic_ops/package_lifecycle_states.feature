@domain:basic_ops @m2 @REQ-API_BASIC-198 @spec(api_basic_operations.md#3331-lifecycle-states)
Feature: Package lifecycle states

  @REQ-API_BASIC-198 @happy
  Scenario: Package lifecycle states define valid states for operations
    Given a package instance
    When lifecycle states are considered
    Then valid states are defined for the package lifecycle
    And state indicates whether the package is created opened or closed
    And operations depend on lifecycle state constraints
    And lifecycle state is tracked in memory
    And lifecycle state definitions align with state transitions and validation rules

