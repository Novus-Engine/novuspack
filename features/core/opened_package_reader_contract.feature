@domain:core @m2 @REQ-CORE-047 @spec(api_core.md#111-opened-package-reader-contract) @spec(api_core.md#filepackage-struct) @spec(api_core.md#filepackage-field-descriptions)
Feature: Opened package reader contract defines reader scope

  @REQ-CORE-047 @happy
  Scenario: PackageReader scope is read-only for opened packages
    Given a package has been opened successfully
    When the caller uses the PackageReader interface
    Then only read-only operations are available
    And the package exposes required reader methods through the FilePackage structure
    And the opened package state defines the assumptions for reader method behavior
