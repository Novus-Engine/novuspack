@domain:core @m2 @REQ-CORE-134 @spec(api_core.md#1235-write-detailed-behavior) @spec(api_writing.md#3-write-strategy-selection)
Feature: Write detailed behavior references write strategy selection

  @REQ-CORE-134 @happy
  Scenario: Write detailed behavior follows strategy selection rules
    Given a package opened for writing
    When Write is called
    Then the detailed behavior follows write strategy selection
    And the behavior is consistent with the specification
    And callers can rely on the documented behavior
