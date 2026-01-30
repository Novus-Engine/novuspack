@domain:file_mgmt @m2 @REQ-FILEMGMT-156 @spec(api_file_mgmt_updates.md#133-filemetadataupdate-structure) @spec(api_file_mgmt_file_entry.md#633-filemetadataupdate-structure)
Feature: FileMetadataUpdate Structure

  @REQ-FILEMGMT-156 @happy
  Scenario: FileMetadataUpdate structure defines basic metadata fields
    Given an open NovusPack package
    And a valid context
    And a FileEntry to update
    When FileMetadataUpdate structure is used
    Then Tags field contains file tags
    And CompressionType field specifies compression type
    And CompressionLevel field specifies compression level
    And EncryptionType field specifies encryption type

  @REQ-FILEMGMT-156 @happy
  Scenario: FileMetadataUpdate structure defines path management fields
    Given an open NovusPack package
    And a valid context
    And a FileEntry to update
    When FileMetadataUpdate structure is used
    Then AddPaths field adds additional paths
    And RemovePaths field removes paths by path string
    And UpdatePaths field updates existing path metadata

  @REQ-FILEMGMT-156 @happy
  Scenario: FileMetadataUpdate structure defines hash management fields
    Given an open NovusPack package
    And a valid context
    And a FileEntry to update
    When FileMetadataUpdate structure is used
    Then AddHashes field adds additional hash entries
    And RemoveHashes field removes hash entries by type
    And UpdateHashes field updates existing hash entries

  @REQ-FILEMGMT-156 @happy
  Scenario: FileMetadataUpdate structure defines optional data field
    Given an open NovusPack package
    And a valid context
    And a FileEntry to update
    When FileMetadataUpdate structure is used
    Then OptionalData field contains structured optional data
    And optional data supports multiple data types
    And optional data can update tags, compression, and extended attributes
