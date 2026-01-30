@domain:core @m2 @REQ-CORE-166 @spec(api_core.md#10311-newpackageerror-function)
Feature: NewPackageError function creates structured error with type-safe context

  @REQ-CORE-166 @happy
  Scenario: NewPackageError creates a structured error with context
    Given an error type and optional context
    When NewPackageError is called
    Then a PackageError is created with type-safe context
    And the error is suitable for return from package operations
    And the behavior matches the NewPackageError specification
