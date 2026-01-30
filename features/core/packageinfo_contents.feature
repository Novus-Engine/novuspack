@domain:core @m2 @REQ-CORE-104 @spec(api_core.md#1175-packageinfo-contents)
Feature: PackageInfo contents define header-derived and computed package-level statistics

  @REQ-CORE-104 @happy
  Scenario: PackageInfo contains header-derived and computed statistics
    Given an opened package
    When GetInfo is called
    Then PackageInfo contains header-derived fields
    And PackageInfo contains computed package-level statistics
    And the contents match the PackageInfo structure specification
