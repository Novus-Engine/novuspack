@domain:file_mgmt @m2 @REQ-FILEMGMT-429 @spec(api_file_mgmt_file_entry.md#103-fileentry-data-management-purpose)
Feature: FileEntry data management purpose defines data loading and processing

  @REQ-FILEMGMT-429 @happy
  Scenario: FileEntry data management purpose defines loading and processing
    Given a FileEntry with data management operations
    When data management purpose is applied
    Then data management purpose defines data loading and processing as specified
    And the behavior matches the data management purpose specification
    And LoadData, UnloadData, ProcessData are available
    And purpose is consistent with specification
