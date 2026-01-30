@domain:core @m2 @REQ-CORE-150 @spec(api_core.md#1255-fastwrite-detailed-behavior) @spec(api_writing.md#22-fastwrite-implementation-strategy)
Feature: FastWrite detailed behavior references in-place package updates specification

  @REQ-CORE-150 @happy
  Scenario: FastWrite detailed behavior follows in-place updates spec
    Given a package opened for writing at an existing path
    When FastWrite is called
    Then the detailed behavior follows the in-place package updates specification
    And in-place semantics are as specified
    And callers can rely on the documented behavior
