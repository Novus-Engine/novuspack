@domain:file_mgmt @m2 @REQ-FILEMGMT-258 @spec(api_file_mgmt_queries.md#333-findentriesbypathpatterns-parameters)
Feature: FindEntriesByPathPatterns parameters define pattern list input

  @REQ-FILEMGMT-258 @happy
  Scenario: FindEntriesByPathPatterns parameters define pattern list
    Given an open NovusPack package
    When FindEntriesByPathPatterns is invoked with parameters
    Then parameters define pattern list input as specified
    And the behavior matches the parameters specification
    And pattern list is validated before use
    And context is accepted for cancellation
