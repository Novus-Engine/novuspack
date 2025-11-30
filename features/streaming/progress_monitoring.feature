@domain:streaming @m2 @REQ-STREAM-009 @spec(api_streaming.md#152-progress-monitoring)
Feature: Progress Monitoring

  @REQ-STREAM-009 @happy
  Scenario: Progress monitoring methods track stream progress
    Given a NovusPack package
    And an open FileStream
    When progress monitoring methods are called
    Then Progress returns detailed progress information
    And EstimatedTimeRemaining returns time estimate
    And progress monitoring enables progress tracking

  @REQ-STREAM-009 @happy
  Scenario: Progress returns comprehensive progress information
    Given a NovusPack package
    And an open FileStream
    When Progress method is called
    Then bytesRead is returned indicating bytes read so far
    And totalBytes is returned indicating total bytes to read
    And readSpeed is returned indicating current read speed in bytes per second
    And elapsed is returned indicating time elapsed since stream creation

  @REQ-STREAM-009 @happy
  Scenario: Progress enables progress percentage calculation
    Given a NovusPack package
    And an open FileStream
    When Progress is used to calculate percentage
    Then bytesRead and totalBytes enable percentage calculation
    And progress percentage is calculated as (bytesRead / totalBytes) * 100
    And progress information enables progress bars and UI updates

  @REQ-STREAM-009 @happy
  Scenario: ReadSpeed tracks current reading performance
    Given a NovusPack package
    And an open FileStream
    When Progress is called during reading
    Then readSpeed reflects current bytes per second
    And readSpeed enables performance monitoring
    And readSpeed updates dynamically during reading

  @REQ-STREAM-009 @happy
  Scenario: Elapsed time tracks stream operation duration
    Given a NovusPack package
    And an open FileStream
    When Progress is called
    Then elapsed reflects time since stream creation
    And elapsed enables time tracking
    And elapsed can be used with readSpeed for performance analysis

  @REQ-STREAM-009 @happy
  Scenario: EstimatedTimeRemaining estimates completion time
    Given a NovusPack package
    And an open FileStream
    When EstimatedTimeRemaining is called
    Then estimated time remaining for completion is returned
    And estimate is calculated based on current read speed
    And time estimate enables ETA display for users

  @REQ-STREAM-009 @error
  Scenario: Progress monitoring handles closed stream correctly
    Given a NovusPack package
    And a closed FileStream
    When Progress or EstimatedTimeRemaining is called
    Then progress information reflects final state
    And bytesRead equals totalBytes if fully read
    And elapsed reflects total time to completion
