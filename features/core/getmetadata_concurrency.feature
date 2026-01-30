@domain:core @m2 @REQ-CORE-115 @spec(api_core.md#1188-getmetadata-concurrency)
Feature: GetMetadata concurrency is safe for concurrent access

  @REQ-CORE-115 @happy
  Scenario: GetMetadata can be called concurrently without inconsistent results
    Given an opened package
    When GetMetadata is called concurrently from multiple goroutines
    Then each call returns consistent PackageMetadata results
    And no concurrency errors occur
    And returned metadata reflects a stable view of package state
