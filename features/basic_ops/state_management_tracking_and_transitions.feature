@domain:basic_ops @m2 @REQ-API_BASIC-123 @spec(api_basic_operations.md#333-state-management)
Feature: State management

  @REQ-API_BASIC-123 @happy
  Scenario: Package state is tracked and transitions are managed consistently
    Given a package lifecycle (create, open, modify, close)
    When operations are performed that change package state
    Then package state is tracked in memory
    And state transitions are well-defined between lifecycle phases
    And state-dependent operations behave consistently with current state
    And invalid state transitions are prevented by validation
    And state management supports predictable API behavior

