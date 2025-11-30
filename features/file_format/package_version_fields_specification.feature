@domain:file_format @m1 @REQ-FILEFMT-004 @spec(package_file_format.md#22-package-version-fields-specification)
Feature: Package version fields specification

  @happy
  Scenario: PackageDataVersion tracks data changes correctly
    Given a new NovusPack package
    When the package is created
    Then PackageDataVersion equals 1
    And PackageDataVersion is a 32-bit unsigned integer
    When a file is added to the package
    Then PackageDataVersion increments
    When a file is removed from the package
    Then PackageDataVersion increments again
    When a file data is modified
    Then PackageDataVersion increments again

  @happy
  Scenario: MetadataVersion tracks metadata changes correctly
    Given a new NovusPack package
    When the package is created
    Then MetadataVersion equals 1
    And MetadataVersion is a 32-bit unsigned integer
    When the package comment is modified
    Then MetadataVersion increments
    When package metadata is updated
    Then MetadataVersion increments again

  @happy
  Scenario: PackageCRC can be zero when skipped
    Given a NovusPack package with PackageCRC set to 0
    When the header is validated
    Then the header is valid
    And PackageCRC indicates calculation was skipped

  @happy
  Scenario: PackageCRC validates package content integrity
    Given a NovusPack package with PackageCRC calculated
    When the header is parsed
    Then PackageCRC is a 32-bit unsigned integer
    And PackageCRC represents CRC32 of all content excluding header and signatures
    When package content is verified
    Then PackageCRC matches the calculated CRC32

  @error
  Scenario: PackageCRC mismatch indicates corruption
    Given a NovusPack package with a calculated PackageCRC
    When package content is modified
    And PackageCRC is not updated
    Then package validation detects CRC mismatch
    And a structured package corruption error is returned
