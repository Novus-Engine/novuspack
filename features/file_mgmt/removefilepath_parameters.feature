@domain:file_mgmt @m2 @REQ-FILEMGMT-168 @spec(api_file_mgmt_updates.md#152-removefilepath-parameters)
Feature: RemoveFilePath parameters include context, entry, and path

  @REQ-FILEMGMT-168 @happy
  Scenario: RemoveFilePath accepts context, entry, and path
    Given a FileEntry and a path to remove
    When RemoveFilePath is called
    Then parameters include context, entry, and path
    And the parameter contract matches the specification
    And the behavior matches the RemoveFilePath parameters specification
