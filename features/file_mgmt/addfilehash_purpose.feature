@domain:file_mgmt @m2 @REQ-FILEMGMT-173 @spec(api_file_mgmt_updates.md#161-addfilehash-purpose)
Feature: AddFileHash purpose is to add hash to file entry

  @REQ-FILEMGMT-173 @happy
  Scenario: AddFileHash adds hash to file entry
    Given a FileEntry and a hash to add
    When AddFileHash is called
    Then the purpose is to add the hash to the file entry
    And the behavior matches the AddFileHash purpose specification
    And hash metadata is updated consistently
