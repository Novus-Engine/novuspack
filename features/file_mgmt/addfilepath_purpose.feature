@domain:file_mgmt @m2 @REQ-FILEMGMT-161 @spec(api_file_mgmt_updates.md#141-addfilepath-purpose)
Feature: AddFilePath purpose is to add additional path to file entry

  @REQ-FILEMGMT-161 @happy
  Scenario: AddFilePath adds additional path to file entry
    Given a FileEntry and a path to add
    When AddFilePath is called
    Then the purpose is to add additional path to the file entry
    And the behavior matches the AddFilePath purpose specification
    And path metadata is updated consistently
