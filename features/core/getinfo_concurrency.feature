@domain:core @m2 @REQ-CORE-107 @spec(api_core.md#1178-getinfo-concurrency)
Feature: GetInfo concurrency is safe for concurrent access

  @REQ-CORE-107 @happy
  Scenario: GetInfo can be called concurrently without inconsistent results
    Given an opened package
    When GetInfo is called concurrently from multiple goroutines
    Then each call returns consistent PackageInfo results
    And no concurrency errors occur
    And results reflect a stable view of package state
