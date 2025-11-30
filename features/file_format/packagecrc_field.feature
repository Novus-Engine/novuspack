@domain:file_format @m2 @REQ-FILEFMT-028 @spec(package_file_format.md#223-packagecrc-field)
Feature: PackageCRC Field

  @REQ-FILEFMT-028 @happy
  Scenario: PackageCRC field stores package checksum
    Given a NovusPack package
    And package has CRC calculated
    When PackageCRC field is examined
    Then PackageCRC field stores package checksum
    And checksum is CRC32 value
    And checksum enables integrity verification

  @REQ-FILEFMT-028 @happy
  Scenario: PackageCRC of 0 indicates checksum is skipped
    Given a NovusPack package
    And package CRC calculation is skipped
    When PackageCRC is examined
    Then PackageCRC is 0
    And checksum calculation was skipped
    And checksum is not available

  @REQ-FILEFMT-028 @happy
  Scenario: PackageCRC enables package integrity verification
    Given a NovusPack package
    And package has PackageCRC calculated
    When PackageCRC is verified
    Then package integrity can be verified
    And CRC mismatch indicates corruption
    And checksum validates package content
