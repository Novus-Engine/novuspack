@domain:core @m2 @REQ-CORE-162 @spec(api_core.md#1022-packageerror-error-method) @spec(api_core.md#packageerrorerror-method)
Feature: PackageError Error method implements error interface

  @REQ-CORE-162 @happy
  Scenario: PackageError implements the error interface
    Given a PackageError value
    When Error is called
    Then the method returns a string suitable for the error interface
    And the return value is suitable for logging and display
    And the behavior matches the PackageError Error method specification
