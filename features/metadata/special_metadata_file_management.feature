@domain:metadata @m2 @REQ-META-090 @spec(api_metadata.md#83-special-metadata-file-management)
Feature: Special Metadata File Management

  @REQ-META-090 @happy
  Scenario: Special metadata file management provides special file operations
    Given a NovusPack package
    And a valid context
    When special metadata file management is used
    Then special file requirements must be met
    And special files must be saved with specific flags
    And file types must be properly set
    And context supports cancellation

  @REQ-META-090 @happy
  Scenario: Special metadata files must meet file type requirements
    Given a NovusPack package
    And a valid context
    When special metadata files are created
    Then special file types must be used
    And reserved file names must be used
    And files must be uncompressed for FastWrite compatibility
    And proper package header flags must be set

  @REQ-META-090 @happy
  Scenario: Special metadata files must meet FileEntry requirements
    Given a NovusPack package
    And a valid context
    When special metadata files are created
    Then Type field is set to appropriate special file type
    And CompressionType is set to 0 (no compression)
    And EncryptionType is set to 0x00 (no encryption)
    And Tags include file_type=special_metadata

  @REQ-META-090 @happy
  Scenario: SavePathMetadataFile creates special metadata file
    Given a NovusPack package
    And a valid context
    And path metadata entries
    When SavePathMetadataFile is called
    Then special metadata file is created
    And file meets all special file requirements
    And package header flags are updated
    And context supports cancellation

  @REQ-META-090 @error
  Scenario: Special metadata file management handles errors
    Given a NovusPack package
    When special file requirements are not met
    Then validation detects requirement violations
    And appropriate errors are returned
