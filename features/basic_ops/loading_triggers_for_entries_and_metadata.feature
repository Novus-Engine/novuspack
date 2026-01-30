@domain:basic_ops @m2 @REQ-API_BASIC-196 @spec(api_basic_operations.md#3323-loading-triggers)
Feature: Loading triggers

  @REQ-API_BASIC-196 @happy
  Scenario: Loading triggers define when file entries and path metadata are loaded
    Given a package opened from disk
    When operations require file entries or path metadata
    Then loading triggers determine when entries are loaded
    And loading triggers determine when path metadata is loaded
    And triggers avoid loading data that is not required by the current workflow
    And triggers preserve correctness and consistency of in-memory state
    And triggers support both eager and on-demand loading strategies

