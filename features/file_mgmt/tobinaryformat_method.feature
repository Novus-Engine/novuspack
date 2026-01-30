@domain:file_mgmt @m2 @REQ-FILEMGMT-098 @spec(api_file_mgmt_file_entry.md#6-marshaling)
Feature: ToBinaryFormat Method

  @REQ-FILEMGMT-098 @happy
  Scenario: Tobinaryformat returns binary representation
    Given a file entry
    When ToBinaryFormat is called
    Then FileEntryBinaryFormat structure is returned
    And binary representation is created
    And conversion operation completes successfully

  @REQ-FILEMGMT-098 @happy
  Scenario: ToBinaryFormat converts FileEntry to binary format
    Given a FileEntry with all fields populated
    When ToBinaryFormat is called
    Then FileEntryBinaryFormat is created
    And all fields are converted to binary format
    And structure matches binary format specification

  @REQ-FILEMGMT-098 @happy
  Scenario: ToBinaryFormat includes variable-length data sections
    Given a FileEntry with paths, hashes, and optional data
    When ToBinaryFormat is called
    Then path entries are included in binary format
    And hash data is included in binary format
    And optional data is included in binary format
    And variable-length sections are correctly serialized

  @REQ-FILEMGMT-098 @happy
  Scenario: ToBinaryFormat enables storage and transmission
    Given a FileEntry
    When ToBinaryFormat is called
    Then binary format is suitable for storage
    And binary format is suitable for transmission
    And format is efficient and compact

  @REQ-FILEMGMT-098 @happy
  Scenario: ToBinaryFormat preserves all file entry information
    Given a FileEntry with complete metadata
    When ToBinaryFormat is called
    Then all static fields are preserved
    And all variable-length data is preserved
    And all metadata is preserved
    And information integrity is maintained

  @REQ-FILEMGMT-098 @error
  Scenario: ToBinaryFormat returns error for invalid FileEntry
    Given an invalid or incomplete FileEntry
    When ToBinaryFormat is called
    Then structured conversion error is returned
    And error indicates conversion failure
    And error follows structured error format

  @REQ-FILEMGMT-098 @error
  Scenario: ToBinaryFormat handles serialization errors
    Given a FileEntry that fails during serialization
    When ToBinaryFormat is called
    Then structured serialization error is returned
    And error indicates serialization issue
    And error follows structured error format
