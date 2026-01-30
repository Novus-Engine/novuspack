@domain:core @m2 @REQ-CORE-120 @spec(api_core.md#1194-validate-behavior) @spec(api_core.md#packagereadervalidate-behavior)
Feature: Validate behavior defines validation process

  @REQ-CORE-120 @happy
  Scenario: Validate performs the defined validation process
    Given an opened package
    When Validate is called
    Then the validation process runs as specified
    And format, structure, and integrity are checked
    And the behavior matches the PackageReader Validate specification
