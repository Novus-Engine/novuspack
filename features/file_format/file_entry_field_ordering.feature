@domain:file_format @m2 @REQ-FILEFMT-054 @spec(package_file_format.md#4131-field-ordering)
Feature: File Entry Field Ordering

  @REQ-FILEFMT-054 @happy
  Scenario: Field ordering defines field arrangement
    Given a file entry
    When field ordering is examined
    Then fields are ordered by size (largest to smallest)
    And field arrangement minimizes padding
    And field ordering optimizes structure layout

  @REQ-FILEFMT-054 @happy
  Scenario: 8-byte fields are ordered first
    Given a file entry
    When field ordering is examined
    Then 8-byte fields come first
    And FileID is 8 bytes and ordered first
    And OriginalSize is 8 bytes and ordered next
    And StoredSize is 8 bytes and ordered after OriginalSize

  @REQ-FILEFMT-054 @happy
  Scenario: 4-byte fields follow 8-byte fields
    Given a file entry
    When field ordering is examined
    Then 4-byte fields follow 8-byte fields
    And RawChecksum, StoredChecksum, FileVersion, MetadataVersion are 4 bytes
    And HashDataOffset, OptionalDataOffset, Reserved are 4 bytes
    And 4-byte fields are grouped together

  @REQ-FILEFMT-054 @happy
  Scenario: 2-byte fields follow 4-byte fields
    Given a file entry
    When field ordering is examined
    Then 2-byte fields follow 4-byte fields
    And PathCount, Type are 2 bytes
    And HashDataLen, OptionalDataLen are 2 bytes
    And 2-byte fields are grouped together

  @REQ-FILEFMT-054 @happy
  Scenario: 1-byte fields are ordered last
    Given a file entry
    When field ordering is examined
    Then 1-byte fields are ordered last
    And CompressionType, CompressionLevel are 1 byte
    And EncryptionType, HashCount are 1 byte
    And 1-byte fields are grouped together

  @REQ-FILEFMT-054 @happy
  Scenario: Field ordering minimizes padding
    Given a file entry
    When fixed structure is serialized
    Then field ordering minimizes padding
    And structure is optimized for 8-byte alignment
    And performance is improved on modern systems
