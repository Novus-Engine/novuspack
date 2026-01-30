@domain:file_mgmt @m2 @REQ-FILEMGMT-284 @spec(api_file_mgmt_queries.md#224-getfilebyoffset-returns)
Feature: GetFileByOffset returns define FileEntry and error return

  @REQ-FILEMGMT-284 @happy
  Scenario: GetFileByOffset returns define results
    Given an open NovusPack package
    When GetFileByOffset completes
    Then returns define FileEntry and error as specified
    And the behavior matches the returns specification
    And nil FileEntry and error are returned when not found
    And error indicates failure conditions
