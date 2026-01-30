@domain:core @m2 @REQ-CORE-133 @spec(api_core.md#1234-write-behavior) @spec(api_writing.md#3-write-strategy-selection) @spec(api_writing.md#533-packagewrite-method)
Feature: Write behavior defines automatic strategy selection and durability

  @REQ-CORE-133 @happy
  Scenario: Write behavior follows automatic strategy and durability rules
    Given a package opened for writing
    When Write is called
    Then the write strategy is selected automatically
    And durability is applied as specified
    And the behavior matches the Write behavior specification
