@domain:file_format @m1 @REQ-FILEFMT-015 @spec(package_file_format.md#414-variable-length-data-follows-fixed-structure)
Feature: FileEntry variable-length data structure

  @happy
  Scenario: Variable-length data follows fixed structure
    Given a file entry with paths, hashes, and optional data
    When the file entry is serialized
    Then fixed structure comes first (64 bytes)
    And variable-length data follows immediately after
    And variable-length data ordering is: paths, hashes, optional data

  @happy
  Scenario: Path entries are at offset 0 in variable-length data
    Given a file entry with PathCount > 0
    When variable-length data is structured
    Then path entries start at offset 0
    And all PathCount paths are present

  @happy
  Scenario: Hash data is located at HashDataOffset
    Given a file entry with HashCount > 0
    When variable-length data is structured
    Then hash data starts at HashDataOffset
    And HashDataLen matches the actual hash data length
    And all HashCount hash entries are present

  @happy
  Scenario: Optional data is located at OptionalDataOffset
    Given a file entry with OptionalDataLen > 0
    When variable-length data is structured
    Then optional data starts at OptionalDataOffset
    And OptionalDataLen matches the actual optional data length

  @error
  Scenario: Variable-length data with invalid offsets causes validation failure
    Given a file entry
    When HashDataOffset or OptionalDataOffset points outside variable-length section
    Then a structured invalid file entry error is returned

  @happy
  Scenario Outline: Path entry structure matches specification
    Given a path entry
    When the path entry is serialized
    Then PathLength (2 bytes) comes first
    And Path (UTF-8, variable) follows
    And Mode (4 bytes), UserID (4 bytes), GroupID (4 bytes) follow
    And ModTime (8 bytes), CreateTime (8 bytes), AccessTime (8 bytes) follow
    And total path entry size matches PathLength + metadata size

    Examples:
      | PathLength | TotalSize |
      | 10         | 48        |
      | 100        | 138       |

  @happy
  Scenario Outline: Hash entry structure matches specification
    Given a hash entry
    When the hash entry is serialized
    Then HashType (1 byte) comes first
    And HashPurpose (1 byte) follows
    And HashLength (2 bytes) follows
    And HashData (variable) follows
    And hash data length matches HashLength

    Examples:
      | HashType | HashLength | TotalSize |
      | 0x00     | 32         | 36        |
      | 0x01     | 64         | 68        |

  @happy
  Scenario Outline: Optional data entry structure matches specification
    Given an optional data entry
    When the optional data entry is serialized
    Then DataType (1 byte) comes first
    And DataLength (2 bytes) follows
    And Data (variable) follows
    And data length matches DataLength

    Examples:
      | DataType | DataLength | TotalSize |
      | 0x00     | 100        | 103       |
      | 0x01     | 1          | 4         |
