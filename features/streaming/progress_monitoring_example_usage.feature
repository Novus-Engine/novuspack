@domain:streaming @m2 @REQ-STREAM-033 @spec(api_streaming.md#1524-example-usage)
Feature: Progress Monitoring Example Usage

  @REQ-STREAM-033 @happy
  Scenario: Progress monitoring example demonstrates progress tracking
    Given an open NovusPack package
    And a valid context
    And file stream with data
    When progress monitoring example is followed
    Then ReadChunk is called in loop
    And Progress returns bytesRead, totalBytes, speed, and elapsed
    And EstimatedTimeRemaining returns remaining time
    And progress information is formatted and displayed

  @REQ-STREAM-033 @happy
  Scenario: Progress monitoring example demonstrates error handling
    Given an open NovusPack package
    And a valid context
    And file stream with potential errors
    When progress monitoring example is followed
    Then error handling checks for read errors
    And loop breaks on error
    And error information is properly handled

  @REQ-STREAM-033 @happy
  Scenario: Progress monitoring example demonstrates context usage
    Given an open NovusPack package
    And a valid context
    And progress monitoring requirements
    When progress monitoring example is followed
    Then context.Context is used correctly
    And context cancellation is handled
    And context timeout is handled
