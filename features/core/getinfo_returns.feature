@domain:core @m2 @REQ-CORE-102 @spec(api_core.md#1173-getinfo-returns) @spec(api_core.md#packagereadergetinfo-returns)
Feature: GetInfo returns define lightweight package information structure

  @REQ-CORE-102 @happy
  Scenario: GetInfo returns a lightweight PackageInfo structure
    Given an opened package
    When GetInfo is called
    Then a PackageInfo structure is returned
    And the structure contains header-derived and computed statistics
    And the return type matches the PackageReader GetInfo contract
