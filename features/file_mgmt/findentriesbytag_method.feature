@domain:file_mgmt @m2 @REQ-FILEMGMT-289 @spec(api_file_mgmt_queries.md#311-package-findentriesbytag-method)
Feature: Package FindEntriesByTag method finds all file entries with a specific tag

  @REQ-FILEMGMT-289 @happy
  Scenario: FindEntriesByTag finds file entries with tag
    Given an open NovusPack package with tagged file entries
    When FindEntriesByTag is called with a tag key and value
    Then all file entries with the specific tag are returned
    And the behavior matches the FindEntriesByTag method specification
    And empty slice is returned when no matches
    And error is returned on failure
