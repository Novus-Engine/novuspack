@domain:basic_ops @m1 @REQ-API_BASIC-012 @spec(api_basic_operations.md#73-package-information)
Feature: Package information retrieval

  @happy
  Scenario: GetInfo returns comprehensive package information
    Given an open NovusPack package
    When GetInfo is called
    Then PackageInfo structure is returned
    And file count is included
    And package sizes are included
    And VendorID and AppID are included
    And package comment is included
    And signature information is included
    And security status is included
    And compression status is included
    And timestamps are included
    And feature flags are included

  @happy
  Scenario: GetInfo returns basic package information
    Given an open NovusPack package
    When GetInfo is called
    Then file count is accurate
    And total package size is accurate
    And compressed size is accurate if applicable
    And package version information is included

  @happy
  Scenario: GetInfo returns package identity information
    Given an open NovusPack package with VendorID and AppID
    When GetInfo is called
    Then VendorID is included in PackageInfo
    And AppID is included in PackageInfo
    And package identity is complete

  @happy
  Scenario: GetInfo returns package comment
    Given an open NovusPack package with a comment
    When GetInfo is called
    Then package comment is included in PackageInfo
    And comment is accessible

  @happy
  Scenario: GetInfo returns signature information
    Given a signed open NovusPack package
    When GetInfo is called
    Then signature information is included in PackageInfo
    And signature count is accurate
    And signature details are included

  @happy
  Scenario: GetInfo returns security status
    Given an open NovusPack package
    When GetInfo is called
    Then security status is included in PackageInfo
    And encryption status is included
    And validation status is included

  @happy
  Scenario: GetInfo returns compression status
    Given an open NovusPack package
    When GetInfo is called
    Then compression status is included in PackageInfo
    And compression type is included
    And compression ratio is included if applicable

  @happy
  Scenario: GetInfo returns timestamps and feature flags
    Given an open NovusPack package
    When GetInfo is called
    Then creation timestamp is included
    And modification timestamp is included
    And feature flags are included

  @error
  Scenario: GetInfo fails if package is not open
    Given a closed NovusPack package
    When GetInfo is called
    Then a structured validation error is returned

  @error
  Scenario: GetInfo fails if package is not open
    Given a closed NovusPack package
    When GetInfo is called
    Then a structured validation error is returned
