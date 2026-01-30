@domain:core @m2 @REQ-CORE-185 @spec(api_core.md#111-packagereader-contract)
Feature: PackageReader contract defines read-only interface for opened packages

  @REQ-CORE-185 @happy
  Scenario: PackageReader contract defines read-only interface
    Given a package has been opened
    When the PackageReader interface is used
    Then only read-only operations are available
    And the contract defines the method set and behavior
    And the behavior matches the PackageReader contract specification
