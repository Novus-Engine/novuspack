@domain:file_mgmt @m2 @REQ-FILEMGMT-108 @REQ-FILEMGMT-432 @spec(api_file_mgmt_file_entry.md#106-fileentry-data-management-error-conditions)
Feature: FileEntry data management error conditions define data operation errors

  @REQ-FILEMGMT-108 @REQ-FILEMGMT-432 @happy
  Scenario: FileEntry data management error conditions define errors
    Given a FileEntry with data management operations
    When LoadData or ProcessData fails
    Then data management error conditions define data operation errors as specified
    And the behavior matches the data management error conditions specification
    And structured errors are returned for nil source, I/O failure, etc.
    And error type indicates failure cause
