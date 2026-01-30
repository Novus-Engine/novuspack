@domain:file_mgmt @m2 @REQ-FILEMGMT-259 @spec(api_file_mgmt_queries.md#334-findentriesbypathpatterns-returns)
Feature: FindEntriesByPathPatterns returns define FileEntry slice and error return

  @REQ-FILEMGMT-259 @happy
  Scenario: FindEntriesByPathPatterns returns define results
    Given an open NovusPack package
    When FindEntriesByPathPatterns completes
    Then returns define FileEntry slice and error as specified
    And the behavior matches the returns specification
    And empty slice is returned when no matches
    And error indicates failure conditions
