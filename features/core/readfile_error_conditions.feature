@domain:core @m2 @REQ-CORE-078 @spec(api_core.md#1155-readfile-error-conditions)
Feature: ReadFile error conditions handle missing files and processing errors

  @REQ-CORE-078 @happy
  Scenario: ReadFile returns structured errors for missing files and processing failures
    Given an opened package
    When ReadFile is called for a missing path
    Then a structured error is returned
    And the error indicates the file was not found
    And processing errors are reported as structured errors when they occur
