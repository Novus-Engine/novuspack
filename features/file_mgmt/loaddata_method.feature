@domain:file_mgmt @m2 @REQ-FILEMGMT-427 @spec(api_file_mgmt_file_entry.md#101-fileentry-loaddata-method)
Feature: LoadData method loads file content into memory

  @REQ-FILEMGMT-427 @happy
  Scenario: LoadData method loads file content
    Given a FileEntry with a valid data source
    When LoadData method is called on the FileEntry
    Then file content is loaded into memory as specified
    And the behavior matches the LoadData method specification
    And content is available via GetData after load
    And error is returned on load failure
