@domain:file_mgmt @m2 @REQ-FILEMGMT-014 @spec(api_file_management.md#63-update-file-metadata)
Feature: Update file metadata operations

  @happy
  Scenario: UpdateFileMetadata updates file metadata without changing data
    Given an open writable NovusPack package with existing file
    When UpdateFileMetadata is called with metadata updates
    Then file metadata is updated
    And metadata version is incremented
    And file data is unchanged
    And file data version is unchanged

  @happy
  Scenario: UpdateFileMetadata supports partial metadata updates
    Given an open writable NovusPack package with existing file
    When UpdateFileMetadata is called with partial updates
    Then specified metadata is updated
    And unspecified metadata remains unchanged
    And updates are applied correctly

  @happy
  Scenario: FileMetadataUpdate structure contains all updatable fields
    Given FileMetadataUpdate structure
    When structure is examined
    Then Paths field exists for path updates
    And Tags field exists for tag updates
    And CompressionType field exists for compression updates
    And EncryptionType field exists for encryption updates
    And other metadata fields exist

  @error
  Scenario: UpdateFileMetadata fails if file does not exist
    Given an open writable NovusPack package
    When UpdateFileMetadata is called with non-existent path
    Then structured validation error is returned

  @error
  Scenario: UpdateFileMetadata validates metadata structure
    Given invalid metadata updates
    When UpdateFileMetadata is called
    Then structured validation error is returned
    And error indicates invalid metadata
