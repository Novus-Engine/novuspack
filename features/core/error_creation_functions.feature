@domain:core @m2 @REQ-CORE-165 @spec(api_core.md#1031-error-creation-functions)
Feature: Error creation functions provide NewPackageError and WrapErrorWithContext

  @REQ-CORE-165 @happy
  Scenario: NewPackageError and WrapErrorWithContext are available
    Given an error condition to report
    When creating or wrapping an error
    Then NewPackageError and WrapErrorWithContext are available
    And the functions create structured errors with context
    And the behavior matches the error creation specification
