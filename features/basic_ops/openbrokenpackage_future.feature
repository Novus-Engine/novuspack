@skip @domain:basic_ops @m2 @REQ-API_BASIC-091 @spec(api_basic_operations.md#53-openbrokenpackage)
Feature: OpenBrokenPackage future behavior

  @REQ-API_BASIC-091 @happy
  Scenario: OpenBrokenPackage exists to support repair workflows
    Given a valid context
    And a broken package file
    When OpenBrokenPackage is called
    Then a package instance is returned or a structured error is returned
    And the result can be used for future repair workflows
