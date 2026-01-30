@domain:file_mgmt @m2 @REQ-FILEMGMT-155 @spec(api_file_mgmt_updates.md#132-updatefilemetadata-parameters)
Feature: UpdateFileMetadata parameters include context, entry, and metadata

  @REQ-FILEMGMT-155 @happy
  Scenario: UpdateFileMetadata accepts context, entry, and metadata
    Given a FileEntry and metadata to update
    When UpdateFileMetadata is called
    Then parameters include context, entry, and metadata
    And the parameter contract matches the specification
    And the behavior matches the UpdateFileMetadata parameters specification
