@domain:file_mgmt @m2 @REQ-FILEMGMT-297 @spec(api_file_mgmt_file_entry.md#41-basic-data-operations)
Feature: Basic data operations define LoadData, UnloadData, GetData, and SetData methods

  @REQ-FILEMGMT-297 @happy
  Scenario: Basic data operations define LoadData UnloadData GetData SetData
    Given a FileEntry with file content
    When basic data operations are used
    Then LoadData, UnloadData, GetData, and SetData methods are defined as specified
    And the behavior matches the basic-data-operations specification
    And LoadData loads file content into memory
    And GetData returns loaded content when available
