@domain:file_format @m2 @REQ-FILEFMT-046 @spec(package_file_format.md#411-file-entry-static-section-field-encoding)
Feature: File Entry Static Section Field Encoding

  @REQ-FILEFMT-046 @happy
  Scenario: File entry static section field encoding defines field representation
    Given a file entry
    When static section field encoding is examined
    Then field representation is defined
    And each field has specified size and format
    And field encoding matches specification

  @REQ-FILEFMT-046 @happy
  Scenario: FileID field is 8 bytes (64-bit unsigned integer)
    Given a file entry
    When static section fields are examined
    Then FileID is 8 bytes
    And FileID is 64-bit unsigned integer
    And FileID provides unique file identifier

  @REQ-FILEFMT-046 @happy
  Scenario: Size fields are 8 bytes (64-bit unsigned integers)
    Given a file entry
    When static section fields are examined
    Then OriginalSize is 8 bytes (64-bit unsigned integer)
    And StoredSize is 8 bytes (64-bit unsigned integer)
    And size fields represent file sizes

  @REQ-FILEFMT-046 @happy
  Scenario: Checksum fields are 4 bytes (32-bit unsigned integers)
    Given a file entry
    When static section fields are examined
    Then RawChecksum is 4 bytes (32-bit unsigned integer)
    And StoredChecksum is 4 bytes (32-bit unsigned integer)
    And checksum fields store CRC32 values

  @REQ-FILEFMT-046 @happy
  Scenario: Version fields are 4 bytes (32-bit unsigned integers)
    Given a file entry
    When static section fields are examined
    Then FileVersion is 4 bytes (32-bit unsigned integer)
    And MetadataVersion is 4 bytes (32-bit unsigned integer)
    And version fields track changes

  @REQ-FILEFMT-046 @happy
  Scenario: PathCount and Type fields are 2 bytes (16-bit unsigned integers)
    Given a file entry
    When static section fields are examined
    Then PathCount is 2 bytes (16-bit unsigned integer)
    And Type is 2 bytes (16-bit unsigned integer)
    And 2-byte fields represent counts and identifiers

  @REQ-FILEFMT-046 @happy
  Scenario: Compression and encryption fields are 1 byte (8-bit unsigned integers)
    Given a file entry
    When static section fields are examined
    Then CompressionType is 1 byte (8-bit unsigned integer)
    And CompressionLevel is 1 byte (8-bit unsigned integer)
    And EncryptionType is 1 byte (8-bit unsigned integer)
    And HashCount is 1 byte (8-bit unsigned integer)

  @REQ-FILEFMT-046 @happy
  Scenario: Offset and length fields use appropriate sizes
    Given a file entry
    When static section fields are examined
    Then HashDataOffset and OptionalDataOffset are 4 bytes (32-bit unsigned integers)
    And HashDataLen and OptionalDataLen are 2 bytes (16-bit unsigned integers)
    And offset and length fields reference variable-length data

  @REQ-FILEFMT-046 @happy
  Scenario: Reserved field is 4 bytes and must be zero
    Given a file entry
    When static section fields are examined
    Then Reserved is 4 bytes (32-bit unsigned integer)
    And Reserved must be 0
    And Reserved field is reserved for future use
