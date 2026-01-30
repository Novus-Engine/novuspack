@domain:file_mgmt @m2 @REQ-FILEMGMT-162 @spec(api_file_mgmt_updates.md#142-addfilepath-parameters)
Feature: AddFilePath parameters include context, entry, and path

  @REQ-FILEMGMT-162 @happy
  Scenario: AddFilePath accepts context, entry, and path
    Given a FileEntry and a path to add
    When AddFilePath is called
    Then parameters include context, entry, and path
    And the parameter contract matches the specification
    And the behavior matches the AddFilePath parameters specification
