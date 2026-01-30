@domain:file_mgmt @m2 @REQ-FILEMGMT-279 @spec(api_file_mgmt_queries.md#214-getfilebypath-returns)
Feature: GetFileByPath returns define FileEntry and error return

  @REQ-FILEMGMT-279 @happy
  Scenario: GetFileByPath returns define results
    Given an open NovusPack package
    When GetFileByPath completes
    Then returns define FileEntry and error as specified
    And the behavior matches the returns specification
    And nil FileEntry and error are returned when not found
    And error indicates failure conditions
