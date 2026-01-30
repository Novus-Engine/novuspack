@domain:file_mgmt @m2 @REQ-FILEMGMT-165 @spec(api_file_mgmt_updates.md#14-packageaddfilepath-method)
Feature: AddFilePath error conditions handle invalid paths

  @REQ-FILEMGMT-165 @happy
  Scenario: AddFilePath returns error for invalid paths
    Given a FileEntry and an invalid path
    When AddFilePath is called
    Then error conditions handle invalid paths as specified
    And returned errors are structured
    And the behavior matches the AddFilePath error conditions specification
