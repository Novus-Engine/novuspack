@domain:file_format @m1 @spec(package_file_format.md#41-file-entry-binary-format-specification)
Feature: Validate FileEntry fixed structure fields

  # Covers 4.1.1 static section fields: sizes, ranges, and zero-reserved
  @happy
  Scenario: Fixed-size fields parse with valid ranges
    Given a NovusPack file with FileEntries present
    When a file entry is parsed
    Then FileID is a non-zero unsigned 64-bit value
    And OriginalSize and StoredSize are non-negative
    And RawChecksum and StoredChecksum are 32-bit values
    And FileVersion and MetadataVersion are >= 1
    And PathCount is >= 1
    And Type is a 16-bit value
    And CompressionType, CompressionLevel, and EncryptionType are 8-bit values
    And HashCount is an 8-bit value
    And HashDataOffset and OptionalDataOffset are 32-bit values
    And HashDataLen and OptionalDataLen are 16-bit values
    And Reserved equals 0

  # Error paths for invalid ranges
  @error
  Scenario Outline: Invalid field ranges are rejected
    Given a file entry with <Field> set to <Value>
    When the file entry is validated
    Then a structured invalid file entry error is returned

    Examples:
      | Field            | Value |
      | PathCount        | 0     |
      | ReservedNonZero  | 1     |
      | FileVersion      | 0     |
      | MetadataVersion  | 0     |
      | FileID           | 0     |

  @error
  Scenario: FileEntry with zero FileID is invalid
    Given a file entry with FileID equals 0
    When the file entry is validated
    Then a structured invalid file entry error is returned

  @error
  Scenario: FileEntry with StoredSize greater than OriginalSize when uncompressed is suspicious
    Given a file entry with CompressionType equals 0
    When StoredSize is greater than OriginalSize
    Then a structured invalid file entry error may be returned
    And validation flags a potential data corruption issue

  @happy
  Scenario: FileEntry size fields handle compression correctly
    Given a file entry with compression applied
    When OriginalSize and StoredSize are compared
    Then StoredSize is less than or equal to OriginalSize
    And both sizes are non-negative

  @error
  Scenario: FileEntry with invalid checksum values is flagged
    Given a file entry
    When RawChecksum or StoredChecksum is zero and should not be
    Then validation may flag missing checksums
    And checksums are validated if present

  @happy
  Scenario: FileEntry fixed structure is exactly 64 bytes
    Given a file entry
    When the fixed structure is serialized
    Then the fixed structure is exactly 64 bytes
    And field offsets match the specification

  @error
  Scenario: FileEntry with PathCount mismatch is invalid
    Given a file entry
    When PathCount does not match the number of path entries
    Then a structured invalid file entry error is returned

  @error
  Scenario: FileEntry with HashCount mismatch is invalid
    Given a file entry
    When HashCount does not match the number of hash entries
    Then a structured invalid file entry error is returned

  @happy
  Scenario: FileEntry validates checksum consistency
    Given a file entry with data
    When RawChecksum is calculated
    Then RawChecksum matches the stored value
    When StoredChecksum is calculated
    Then StoredChecksum matches the stored value

  @happy
  Scenario: NewFileEntry creates entry with zero values
    Given NewFileEntry is called
    Then a FileEntry is returned
    And FileEntry is in initialized state
    And file entry all fields are zero or empty

  @happy
  Scenario: WriteTo serializes file entry to binary format
    Given a FileEntry with values
    When file entry WriteTo is called with writer
    Then file entry is written to writer
    And fixed structure is written first (64 bytes)
    And variable-length data follows
    And written data matches file entry content

  @happy
  Scenario: ReadFrom deserializes file entry from binary format
    Given a reader with valid file entry data
    When file entry ReadFrom is called with reader
    Then file entry is read from reader
    And file entry fields match reader data
    And file entry is valid

  @happy
  Scenario: File entry round-trip serialization preserves all fields
    Given a FileEntry with all fields set
    When file entry WriteTo is called with writer
    And ReadFrom is called with written data
    Then all file entry fields are preserved
    And file entry is valid

  @error
  Scenario: ReadFrom fails with incomplete fixed structure
    Given a reader with less than 64 bytes of file entry data
    When file entry ReadFrom is called with reader
    Then structured IO error is returned
    And error indicates read failure
