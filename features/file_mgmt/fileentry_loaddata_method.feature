@domain:file_mgmt @m2 @REQ-FILEMGMT-298 @spec(api_file_mgmt_file_entry.md#411-fileentry-loaddata-method)
Feature: FileEntry LoadData method loads file content into memory

  @REQ-FILEMGMT-298 @happy
  Scenario: LoadData loads file content into memory
    Given a FileEntry with a valid data source
    When LoadData is called on the FileEntry
    Then file content is loaded into memory as specified
    And the behavior matches the LoadData method specification
    And content is available via GetData after load
    And error is returned on load failure
