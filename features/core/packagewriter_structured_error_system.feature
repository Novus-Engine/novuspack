@domain:core @m2 @REQ-CORE-184 @spec(api_core.md#105-packagewriter-structured-error-system)
Feature: PackageWriter structured error system uses structured errors exclusively

  @REQ-CORE-184 @happy
  Scenario: PackageWriter methods return only structured errors
    Given a package opened for writing
    When a PackageWriter method returns an error
    Then the error is a structured error
    And no sentinel or untyped errors are returned
    And the behavior matches the PackageWriter structured error specification
