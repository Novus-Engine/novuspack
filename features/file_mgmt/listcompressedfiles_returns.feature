@domain:file_mgmt @m2 @REQ-FILEMGMT-265 @spec(api_file_mgmt_queries.md#424-listcompressedfiles-returns)
Feature: ListCompressedFiles returns define FileEntry slice and error return

  @REQ-FILEMGMT-265 @happy
  Scenario: ListCompressedFiles returns define results
    Given an open NovusPack package
    When ListCompressedFiles completes
    Then returns define FileEntry slice and error as specified
    And the behavior matches the returns specification
    And slice contains only compressed file entries
    And error indicates failure conditions
