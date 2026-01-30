@domain:basic_ops @m2 @REQ-API_BASIC-143 @spec(api_basic_operations.md#953-package-getpath-method)
Feature: Package.GetPath method

  @REQ-API_BASIC-143 @happy
  Scenario: GetPath returns the current package file path
    Given a package opened from a file path
    When GetPath is called
    Then it returns the current package file path
    And the returned path matches the path used to open the package
    When the package is created at a target path
    Then GetPath returns that target path for subsequent operations

