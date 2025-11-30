@domain:metadata @m2 @REQ-META-087 @spec(api_metadata.md#74-package-information-methods)
Feature: Package Information Methods

  @REQ-META-087 @happy
  Scenario: Package information methods provide package information access
    Given a NovusPack package
    And a valid context
    When package information methods are used
    Then GetPackageInfo returns comprehensive package information
    And GetSecurityStatus returns current security status
    And RefreshPackageInfo refreshes package information from current state
    And context supports cancellation

  @REQ-META-087 @happy
  Scenario: GetPackageInfo returns comprehensive package information
    Given a NovusPack package
    And a valid context
    When GetPackageInfo is called
    Then PackageInfo structure is returned
    And structure contains all package information fields
    And context supports cancellation

  @REQ-META-087 @happy
  Scenario: GetSecurityStatus returns current security status
    Given a NovusPack package
    When GetSecurityStatus is called
    Then SecurityStatus structure is returned
    And structure contains signature validation results
    And structure contains checksum validation results
    And structure contains overall security level

  @REQ-META-087 @happy
  Scenario: RefreshPackageInfo refreshes package information
    Given a NovusPack package
    And a valid context
    And package state has changed
    When RefreshPackageInfo is called
    Then package information is refreshed from current state
    And information reflects latest package state
    And context supports cancellation

  @REQ-META-087 @error
  Scenario: Package information methods handle errors
    Given a NovusPack package
    When package information operations fail
    Then appropriate errors are returned
    And errors follow structured error format
