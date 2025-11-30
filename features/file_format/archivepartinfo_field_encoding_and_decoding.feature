@domain:file_format @m1 @REQ-FILEFMT-009 @spec(package_file_format.md#26-archivepartinfo-field-specification)
Feature: ArchivePartInfo field encoding and decoding

  @happy
  Scenario: ArchivePartInfo defaults to single archive (part 1 of 1)
    Given a new NovusPack package
    When the package is created
    Then ArchivePartInfo equals 0x00010000
    And part number equals 1
    And total parts equals 1

  @happy
  Scenario: ArchivePartInfo encodes part number in bits 31-16
    Given a NovusPack package
    When ArchivePartInfo is set with part number 2 and total parts 3
    Then bits 31-16 equal 0x0002
    And part number equals 2

  @happy
  Scenario: ArchivePartInfo encodes total parts in bits 15-0
    Given a NovusPack package
    When ArchivePartInfo is set with part number 2 and total parts 3
    Then bits 15-0 equal 0x0003
    And total parts equals 3

  @happy
  Scenario Outline: ArchivePartInfo packs and unpacks correctly
    Given a NovusPack package
    When ArchivePartInfo is set to <ArchivePartInfo>
    Then part number equals <Part>
    And total parts equals <Total>
    And ArchivePartInfo can be decoded correctly

    Examples:
      | ArchivePartInfo | Part | Total |
      | 0x00010000      | 1    | 1     |
      | 0x00020003      | 2    | 3     |
      | 0x0005000A      | 5    | 10    |
      | 0x00100064      | 16   | 100   |

  @error
  Scenario: ArchivePartInfo with zero part number is invalid for multi-part archives
    Given a NovusPack package
    When ArchivePartInfo is set with part number 0 and total parts greater than 1
    Then a structured invalid archive part info error is returned

  @error
  Scenario: ArchivePartInfo with part number greater than total parts is invalid
    Given a NovusPack package
    When ArchivePartInfo is set with part number greater than total parts
    Then a structured invalid archive part info error is returned

  @happy
  Scenario: ArchivePartInfo supports split archive scenarios
    Given a large NovusPack archive split across multiple files
    When ArchivePartInfo is set for each part
    Then each part correctly identifies its position and total parts
    And ArchiveChainID links related archive parts
