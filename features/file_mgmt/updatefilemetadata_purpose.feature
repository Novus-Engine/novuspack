@domain:file_mgmt @m2 @REQ-FILEMGMT-154 @spec(api_file_mgmt_updates.md#131-updatefilemetadata-purpose)
Feature: UpdateFileMetadata purpose is to update metadata without changing content

  @REQ-FILEMGMT-154 @happy
  Scenario: UpdateFileMetadata updates metadata without changing content
    Given a FileEntry and metadata to update
    When UpdateFileMetadata is called
    Then metadata is updated without changing file content
    And the purpose matches the UpdateFileMetadata specification
    And the behavior matches the UpdateFileMetadata purpose specification
