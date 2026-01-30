@domain:file_mgmt @m2 @REQ-FILEMGMT-170 @spec(api_file_mgmt_updates.md#15-packageremovefilepath-method)
Feature: RemoveFilePath behavior removes path from file entry

  @REQ-FILEMGMT-170 @happy
  Scenario: RemoveFilePath removes path from file entry
    Given a FileEntry with multiple paths
    When RemoveFilePath is called for a path
    Then the path is removed from the file entry
    And the behavior matches the RemoveFilePath specification
    And path metadata is updated as specified
