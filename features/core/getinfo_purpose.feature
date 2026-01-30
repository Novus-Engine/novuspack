@domain:core @m2 @REQ-CORE-100 @spec(api_core.md#1171-getinfo-purpose)
Feature: GetInfo purpose defines lightweight package information retrieval

  @REQ-CORE-100 @happy
  Scenario: GetInfo retrieves lightweight package information
    Given an opened package
    When GetInfo is called
    Then lightweight package information is retrieved
    And the purpose is to expose header-derived and computed stats
    And no full per-file metadata is loaded
