@domain:core @m2 @REQ-CORE-117 @spec(api_core.md#1191-validate-purpose)
Feature: Validate purpose defines package format, structure, and integrity validation

  @REQ-CORE-117 @happy
  Scenario: Validate verifies format, structure, and integrity
    Given an opened package
    When Validate is called
    Then package format is validated
    And package structure is validated
    And package integrity is validated
