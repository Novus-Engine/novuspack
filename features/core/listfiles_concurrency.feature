@domain:core @m2 @REQ-CORE-090 @spec(api_core.md#1166-listfiles-concurrency)
Feature: ListFiles concurrency defines safe concurrent access

  @REQ-CORE-090 @happy
  Scenario: ListFiles is safe for concurrent calls
    Given an opened package
    When ListFiles is called concurrently from multiple goroutines
    Then each call completes without data races
    And results are consistent for an unchanged package state
    And no concurrency errors occur during listing
