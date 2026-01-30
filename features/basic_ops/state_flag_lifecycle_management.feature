@domain:basic_ops @m2 @REQ-API_BASIC-216 @spec(api_basic_operations.md#3373-state-flag-lifecycle)
Feature: State flag lifecycle

  @REQ-API_BASIC-216 @happy
  Scenario: State flag lifecycle defines how state flags are managed
    Given internal state flags used by the package
    When operations change the package state
    Then state flags are set and cleared consistently
    And state flags reflect the current lifecycle state
    And state flags support validation of state-dependent operations
    And state flags remain consistent under error conditions
    And state flag lifecycle aligns with documented state management

