@domain:file_mgmt @m2 @REQ-FILEMGMT-256 @spec(api_file_mgmt_queries.md#331-package-findentriesbypathpatterns-method)
Feature: Package FindEntriesByPathPatterns method gets files matching patterns

  @REQ-FILEMGMT-256 @happy
  Scenario: FindEntriesByPathPatterns gets files matching patterns
    Given an open NovusPack package with file entries
    When FindEntriesByPathPatterns is called with path patterns
    Then file entries matching patterns are returned
    And the behavior matches the FindEntriesByPathPatterns method specification
    And pattern matching follows the specification
    And error is returned for invalid patterns
