@domain:core @m2 @REQ-CORE-079 @spec(api_core.md#1156-readfile-concurrency)
Feature: ReadFile is safe for concurrent reads

  @REQ-CORE-079 @happy
  Scenario: ReadFile supports concurrent read calls from multiple goroutines
    Given an opened package
    When ReadFile is called concurrently from multiple goroutines
    Then each call succeeds or fails independently
    And concurrent reads do not corrupt shared package state
    And no concurrency errors occur during read operations
