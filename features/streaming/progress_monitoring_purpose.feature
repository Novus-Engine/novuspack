@domain:streaming @m2 @REQ-STREAM-030 @spec(api_streaming.md#1521-purpose)
Feature: Progress monitoring purpose defines progress tracking

  @REQ-STREAM-030 @happy
  Scenario: Progress monitoring purpose defines progress tracking
    Given a FileStream or streaming operation
    When progress monitoring is used
    Then the purpose defines progress tracking for the stream
    And progress information is available as specified
    And the behavior matches the progress monitoring purpose specification
    And callers can track streaming progress
