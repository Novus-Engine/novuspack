@domain:core @m2 @REQ-CORE-168 @spec(api_core.md#1032-error-query-functions)
Feature: Error query functions provide AsPackageError and GetErrorContext

  @REQ-CORE-168 @happy
  Scenario: AsPackageError and GetErrorContext are available for error inspection
    Given an error that may be a PackageError
    When AsPackageError or GetErrorContext is called
    Then the error can be inspected for PackageError type
    And type-safe context can be retrieved when present
    And the behavior matches the error query specification
