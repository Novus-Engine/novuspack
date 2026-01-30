@domain:file_mgmt @m2 @REQ-FILEMGMT-167 @spec(api_file_mgmt_updates.md#151-removefilepath-purpose)
Feature: RemoveFilePath purpose is to remove path from file entry

  @REQ-FILEMGMT-167 @happy
  Scenario: RemoveFilePath removes path from file entry
    Given a FileEntry with multiple paths
    When RemoveFilePath is called for a path
    Then the purpose is to remove the path from the file entry
    And the behavior matches the RemoveFilePath purpose specification
    And path metadata is updated consistently
