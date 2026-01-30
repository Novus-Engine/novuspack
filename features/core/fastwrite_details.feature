@domain:core @m2 @REQ-CORE-158 @spec(api_core.md#1275-fastwrite-details) @spec(api_writing.md#2-fastwrite---in-place-package-updates)
Feature: FastWrite details define in-place update requirements

  @REQ-CORE-158 @happy
  Scenario: FastWrite details define in-place update requirements
    Given a package opened for writing at an existing path
    When FastWrite is called
    Then the in-place update requirements are met
    And the details match the FastWrite specification
    And callers can rely on the documented in-place behavior
