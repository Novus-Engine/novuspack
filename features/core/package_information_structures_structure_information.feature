@domain:core @m2 @REQ-CORE-044 @spec(api_core.md#93-package-information-structures)
Feature: Package information structures

  @REQ-CORE-044 @happy
  Scenario: PackageInfo structure provides comprehensive package information
    Given a NovusPack package
    When PackageInfo is retrieved
    Then comprehensive package information is provided
    And basic package information includes file count and sizes
    And package identity includes VendorID and AppID
    And package comment information is included
    And digital signature information is included
    And security information includes security level and immutability
    And timestamps include creation and modification times
    And package features indicate metadata, encryption, and compression status

  @REQ-CORE-044 @happy
  Scenario: SignatureInfo structure provides detailed signature information
    Given a signed NovusPack package
    When SignatureInfo is retrieved
    Then detailed signature information is provided
    And signature information includes signature type and algorithm
    And signature validation status is included
    And signature metadata is available
    And signature details enable signature verification

  @REQ-CORE-044 @happy
  Scenario: SecurityStatus structure provides current security status
    Given a NovusPack package
    When SecurityStatus is retrieved
    Then current security status is provided
    And security validation results are included
    And security level information is available
    And security status enables security assessment
