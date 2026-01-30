@domain:file_mgmt @m2 @REQ-FILEMGMT-262 @spec(api_file_mgmt_queries.md#42-listcompressedfiles)
Feature: Package ListCompressedFiles method returns all compressed file entries

  @REQ-FILEMGMT-262 @happy
  Scenario: ListCompressedFiles returns compressed file entries
    Given an open NovusPack package with compressed files
    When ListCompressedFiles is called
    Then all compressed file entries are returned
    And the behavior matches the ListCompressedFiles specification
    And FileEntry slice and error are returned
    And empty slice is returned when no compressed files exist
