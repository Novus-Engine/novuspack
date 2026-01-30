@domain:file_mgmt @m2 @REQ-FILEMGMT-061 @spec(api_file_mgmt_file_entry.md#10-fileentry-data-management)
Feature: File entry data management provides data operation support

  @REQ-FILEMGMT-061 @happy
  Scenario: File entry data management provides data operation support
    Given a FileEntry with file data
    When data management operations are used
    Then data operation support is provided as specified
    And LoadData, ProcessData, and related methods follow the contract
    And the behavior matches the file entry data management specification
