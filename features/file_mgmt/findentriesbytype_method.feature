@domain:file_mgmt @m2 @REQ-FILEMGMT-290 @spec(api_file_mgmt_queries.md#321-package-findentriesbytype-method)
Feature: Package FindEntriesByType method finds all file entries with a specific file type

  @REQ-FILEMGMT-290 @happy
  Scenario: FindEntriesByType finds file entries by type
    Given an open NovusPack package with file entries of various types
    When FindEntriesByType is called with a file type
    Then all file entries with the specific file type are returned
    And the behavior matches the FindEntriesByType method specification
    And empty slice is returned when no matches
    And error is returned on failure
