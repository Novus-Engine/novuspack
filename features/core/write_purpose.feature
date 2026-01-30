@domain:core @m2 @REQ-CORE-130 @spec(api_core.md#1231-write-purpose) @spec(api_writing.md#3-write-strategy-selection)
Feature: Write purpose defines automatic write strategy selection

  @REQ-CORE-130 @happy
  Scenario: Write selects write strategy automatically
    Given a package opened for writing
    When Write is called
    Then the write strategy is selected automatically
    And the purpose is to simplify caller choice of SafeWrite vs FastWrite
    And the selected strategy matches the specification
