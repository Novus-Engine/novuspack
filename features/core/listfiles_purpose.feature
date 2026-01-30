@domain:core @m2 @REQ-CORE-085 @spec(api_core.md#1161-listfiles-purpose)
Feature: ListFiles purpose defines file information retrieval

  @REQ-CORE-085 @happy
  Scenario: ListFiles retrieves file information for package contents
    Given an opened package with files
    When ListFiles is called
    Then file information is retrieved for package contents
    And the purpose is to expose lightweight file metadata
    And the behavior matches the ListFiles purpose specification
