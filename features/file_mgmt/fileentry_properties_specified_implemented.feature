@domain:file_mgmt @m2 @REQ-FILEMGMT-232 @spec(api_file_mgmt_file_entry.md#7-fileentry-properties)
Feature: 11.1 File Entry Properties is specified and implemented

  @REQ-FILEMGMT-232 @happy
  Scenario: File Entry Properties are specified and implemented
    Given an open NovusPack package with file entries
    When FileEntry properties are accessed
    Then FileEntry properties specification is implemented
    And property access matches the specification
    And the behavior matches the fileentry-properties specification
    And all required properties are available
