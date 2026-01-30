@domain:core @m2 @REQ-CORE-116 @spec(api_core.md#119-validate-method-contract)
Feature: Validate method contract defines package validation interface

  @REQ-CORE-116 @happy
  Scenario: Validate provides a package validation method contract
    Given an opened package
    When Validate is called
    Then package validation is performed through a defined method contract
    And the method is part of the PackageReader interface
    And validation returns structured errors on failures
