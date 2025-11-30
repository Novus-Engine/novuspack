@domain:file_mgmt @m2 @REQ-FILEMGMT-091 @spec(api_file_management.md#1012-fileentry-access)
Feature: FileEntry Access

  @REQ-FILEMGMT-091 @happy
  Scenario: FileEntry access returns comprehensive FileEntry objects
    Given an open NovusPack package
    And a valid context
    When file query functions are called
    Then FileEntry objects are returned
    And FileEntry objects contain all metadata
    And FileEntry objects contain compression status
    And FileEntry objects contain encryption details
    And FileEntry objects contain checksums and timestamps

  @REQ-FILEMGMT-091 @happy
  Scenario: FileEntry access replaces GetFileInfo function
    Given an open NovusPack package
    And a valid context
    When file information is needed
    Then GetFileByPath returns complete FileEntry
    And separate GetFileInfo function is not needed
    And FileEntry provides comprehensive information

  @REQ-FILEMGMT-091 @happy
  Scenario: FileEntry access provides FileEntry arrays for bulk queries
    Given an open NovusPack package
    And a valid context
    When bulk file queries are performed
    Then FileEntry arrays are returned
    And each FileEntry contains complete information
    And FileEntry access is consistent across query types
