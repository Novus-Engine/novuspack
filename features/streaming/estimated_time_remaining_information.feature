@domain:streaming @m2 @REQ-STREAM-032 @spec(api_streaming.md#1523-estimatedtimeremaining-returns)
Feature: Estimated Time Remaining Information

  @REQ-STREAM-032 @happy
  Scenario: EstimatedTimeRemaining returns time estimate
    Given a NovusPack package
    And an open FileStream
    When EstimatedTimeRemaining method is called
    Then estimated time remaining for completion is returned
    And estimate is based on current read speed
    And estimate is based on remaining bytes
    And estimate enables time-based progress display

  @REQ-STREAM-032 @happy
  Scenario: EstimatedTimeRemaining calculates estimate from progress
    Given a NovusPack package
    And an open FileStream
    And data has been read from stream
    When EstimatedTimeRemaining is called
    Then estimate uses bytesRead and totalBytes from Progress
    And estimate uses readSpeed from Progress
    And estimate calculates (totalBytes - bytesRead) / readSpeed
    And estimate provides time duration

  @REQ-STREAM-032 @happy
  Scenario: EstimatedTimeRemaining updates dynamically
    Given a NovusPack package
    And an open FileStream
    When EstimatedTimeRemaining is called multiple times
    Then estimate updates as reading progresses
    And estimate reflects current read speed
    And estimate becomes more accurate as more data is read

  @REQ-STREAM-032 @error
  Scenario: EstimatedTimeRemaining handles edge cases
    Given a NovusPack package
    And an open FileStream
    When EstimatedTimeRemaining is called
    Then estimate handles zero read speed gracefully
    And estimate handles completed stream correctly
    And estimate returns appropriate value for edge cases
