@domain:core @m1 @REQ-CORE-008 @spec(api_core.md#61-package-compression-types)
Feature: Package compression type identification

  @happy
  Scenario: GetPackageCompressionType returns compression type
    Given a NovusPack package
    When GetPackageCompressionType is called
    Then compression type is returned
    And type matches header flags bits 15-8

  @happy
  Scenario: IsPackageCompressed checks compression status
    Given a NovusPack package
    When IsPackageCompressed is called
    Then true is returned if package is compressed
    And false is returned if package is not compressed
    And status matches header flags

  @happy
  Scenario: CanCompressPackage checks if compression is allowed
    Given a NovusPack package
    When CanCompressPackage is called
    Then true is returned if package is not signed
    And false is returned if package is signed
    And status matches SignatureOffset check
