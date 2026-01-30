@domain:metadata @m2 @REQ-META-007 @spec(api_metadata.md#7-package-information-structures)
Feature: Package Information Structures and Methods

  @happy
  Scenario: GetPackageInfo returns comprehensive package information
    Given an open package
    When GetPackageInfo is called
    Then PackageInfo structure is returned
    And package information is complete
    And all metadata is included

  @REQ-META-008 @happy
  Scenario: GetSecurityStatus returns current security status
    Given an open package
    When GetSecurityStatus is called
    Then SecurityStatus is returned
    And encryption status is included
    And signature status is included
    And validation status is included

  @REQ-META-009 @happy
  Scenario: RefreshPackageInfo refreshes package information cache
    Given an open package with cached information
    When RefreshPackageInfo is called
    Then package information cache is refreshed
    And current package state is reflected
    And information is up to date

  @error
  Scenario: Package information methods fail if package is not open
    Given a closed package
    When GetPackageInfo is called
    Then structured validation error is returned

  @REQ-META-011 @REQ-META-014 @error
  Scenario: GetPackageInfo fails if package is not open
    Given a closed package
    When GetPackageInfo is called
    Then structured validation error is returned
