@domain:file_mgmt @m2 @REQ-FILEMGMT-164 @spec(api_file_mgmt_updates.md#14-packageaddfilepath-method)
Feature: AddFilePath behavior adds path to existing file entry

  @REQ-FILEMGMT-164 @happy
  Scenario: AddFilePath adds path to existing file entry
    Given a FileEntry and a path to add
    When AddFilePath is called
    Then the path is added to the existing file entry
    And the behavior matches the AddFilePath specification
    And path metadata is updated as specified
