@domain:file_mgmt @m2 @REQ-FILEMGMT-176 @spec(api_file_mgmt_updates.md#16-packageaddfilehash-method)
Feature: AddFileHash behavior adds hash to file entry metadata

  @REQ-FILEMGMT-176 @happy
  Scenario: AddFileHash adds hash to file entry metadata
    Given a FileEntry and a hash to add
    When AddFileHash is called
    Then the hash is added to file entry metadata
    And the behavior matches the AddFileHash specification
    And hash metadata is updated as specified
