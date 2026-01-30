@domain:file_mgmt @m2 @REQ-FILEMGMT-255 @spec(api_file_mgmt_queries.md#3-multi-entry-queries)
Feature: Multi-entry queries return slice of FileEntry values for tag, type, or pattern-based queries

  @REQ-FILEMGMT-255 @happy
  Scenario: Multi-entry queries return FileEntry slice
    Given an open NovusPack package with multiple file entries
    When multi-entry queries are performed by tag, type, or pattern
    Then a slice of FileEntry values is returned
    And results match the query criteria
    And the behavior matches the multi-entry-queries specification
    And error is returned when package is not open
