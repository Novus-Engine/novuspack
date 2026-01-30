@domain:basic_ops @m2 @REQ-API_BASIC-141 @spec(api_basic_operations.md#951-package-isopen-method)
Feature: Package.IsOpen method

  @REQ-API_BASIC-141 @happy
  Scenario: IsOpen reports whether the package is currently open
    Given a package instance
    When the package has not been opened
    Then IsOpen reports false
    When the package is opened successfully
    Then IsOpen reports true
    When the package is closed
    Then IsOpen reports false
    And IsOpen reflects the current lifecycle state

