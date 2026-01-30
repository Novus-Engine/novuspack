@domain:file_mgmt @m2 @REQ-FILEMGMT-263 @spec(api_file_mgmt_queries.md#422-listcompressedfiles-purpose)
Feature: ListCompressedFiles purpose defines compressed file entry retrieval

  @REQ-FILEMGMT-263 @happy
  Scenario: ListCompressedFiles purpose defines retrieval
    Given an open NovusPack package
    When ListCompressedFiles purpose is applied
    Then compressed file entry retrieval is defined as specified
    And the behavior matches the purpose specification
    And no input parameters are required
    And FileEntry slice and error are returned
