@domain:file_mgmt @m2 @REQ-FILEMGMT-157 @REQ-FILEMGMT-014 @spec(api_file_mgmt_updates.md#134-updatefilemetadata-returns) @spec(api_file_mgmt_file_entry.md#633-filemetadataupdate-structure) @spec(api_file_mgmt_file_entry.md#634-updatefilemetadata-returns) @spec(api_file_mgmt_file_entry.md#635-updatefilemetadata-behavior) @spec(api_file_mgmt_file_entry.md#636-updatefilemetadata-error-conditions)
Feature: UpdateFileMetadata

  @REQ-FILEMGMT-157 @REQ-FILEMGMT-014 @happy
  Scenario: Updatefilemetadata returns updated fileentry
    Given an open writable package
    And existing file entry
    And FileMetadataUpdate structure with metadata changes
    When UpdateFileMetadata is called
    Then updated FileEntry is returned
    And metadata is updated without changing content
    And content remains unchanged

  @REQ-FILEMGMT-157 @REQ-FILEMGMT-014 @happy
  Scenario: UpdateFileMetadata updates metadata without changing content
    Given an open writable package
    And existing file entry with content
    And FileMetadataUpdate with tag updates
    When UpdateFileMetadata is called
    Then file metadata is updated
    And file content is unchanged
    And content integrity is preserved

  @REQ-FILEMGMT-157 @REQ-FILEMGMT-014 @REQ-FILEMGMT-158 @happy
  Scenario: UpdateFileMetadata increments MetadataVersion
    Given an open writable package
    And existing file entry with MetadataVersion 5
    When UpdateFileMetadata updates metadata
    Then MetadataVersion is incremented
    And metadata change is tracked
    And FileVersion remains unchanged

  @REQ-FILEMGMT-157 @REQ-FILEMGMT-014 @happy
  Scenario: UpdateFileMetadata supports partial metadata updates
    Given an open writable package
    And existing file entry
    And FileMetadataUpdate with only tags specified
    When UpdateFileMetadata is called
    Then tags are updated
    And unspecified metadata remains unchanged
    And partial updates work correctly

  @REQ-FILEMGMT-157 @REQ-FILEMGMT-156 @happy
  Scenario: FileMetadataUpdate structure contains all updatable fields
    Given FileMetadataUpdate structure
    When structure is examined
    Then Tags field exists for tag updates
    And CompressionType and CompressionLevel fields exist
    And EncryptionType field exists
    And AddPaths, RemovePaths, UpdatePaths fields exist
    And AddHashes, RemoveHashes, UpdateHashes fields exist
    And OptionalData field exists

  @REQ-FILEMGMT-157 @REQ-FILEMGMT-159 @error
  Scenario: UpdateFileMetadata handles invalid metadata errors
    Given an open writable package
    And existing file entry
    And invalid FileMetadataUpdate structure
    When UpdateFileMetadata is called
    Then structured validation error is returned
    And error indicates invalid metadata
    And error follows structured error format

  @REQ-FILEMGMT-157 @error
  Scenario: UpdateFileMetadata respects context cancellation
    Given an open writable package
    And a cancelled context
    When UpdateFileMetadata is called
    Then structured context error is returned
    And error follows structured error format
