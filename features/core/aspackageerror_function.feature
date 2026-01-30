@domain:core @m2 @REQ-CORE-169 @spec(api_core.md#10321-aspackageerror-function)
Feature: AsPackageError function checks if error is PackageError

  @REQ-CORE-169 @happy
  Scenario: AsPackageError checks whether an error is PackageError
    Given an error value
    When AsPackageError is called
    Then the function returns the PackageError if the error is a PackageError
    And returns false or nil when the error is not a PackageError
    And the behavior matches the AsPackageError specification
