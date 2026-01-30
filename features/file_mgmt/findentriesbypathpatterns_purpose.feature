@domain:file_mgmt @m2 @REQ-FILEMGMT-257 @spec(api_file_mgmt_queries.md#332-findentriesbypathpatterns-purpose)
Feature: FindEntriesByPathPatterns purpose defines file entry lookup by path pattern matching

  @REQ-FILEMGMT-257 @happy
  Scenario: FindEntriesByPathPatterns purpose defines pattern lookup
    Given an open NovusPack package
    When FindEntriesByPathPatterns purpose is applied
    Then file entry lookup by path pattern matching is defined
    And the behavior matches the purpose specification
    And pattern list input is accepted
    And FileEntry slice and error are returned as specified
