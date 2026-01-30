@domain:core @m2 @REQ-CORE-122 @spec(api_core.md#1196-validate-concurrency)
Feature: Validate concurrency defines safe concurrent access

  @REQ-CORE-122 @happy
  Scenario: Validate is safe for concurrent calls
    Given an opened package
    When Validate is called concurrently from multiple goroutines
    Then each call completes without data races
    And results are consistent for an unchanged package state
    And no concurrency errors occur during validation
