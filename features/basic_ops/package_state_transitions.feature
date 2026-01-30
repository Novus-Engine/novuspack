@domain:basic_ops @m2 @REQ-API_BASIC-199 @spec(api_basic_operations.md#3332-state-transitions)
Feature: Package state transitions

  @REQ-API_BASIC-199 @happy
  Scenario: Package moves between lifecycle states via defined transitions
    Given a package lifecycle
    When the package is created opened and closed
    Then transitions between states are defined
    And transitions occur only through valid operations
    And invalid transitions are prevented
    And transitions preserve internal invariants and resource lifecycle rules
    And transitions align with state validation and state-dependent operations

@domain:basic_ops @m2 @REQ-API_BASIC-199 @spec(api_basic_operations.md#3332-state-transitions)
Feature: Package state transitions

  @REQ-API_BASIC-199 @happy
  Scenario: Package moves between lifecycle states via defined transitions
    Given a package lifecycle
    When the package is created opened and closed
    Then transitions between states are defined
    And transitions occur only through valid operations
    And invalid transitions are prevented
    And transitions preserve internal invariants and resource lifecycle rules
    And transitions align with state validation and state-dependent operations

