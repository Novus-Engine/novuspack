@domain:streaming @m2 @REQ-STREAM-031 @spec(api_streaming.md#1522-progress-returns)
Feature: Streaming Progress Information

  @REQ-STREAM-031 @happy
  Scenario: Progress returns streaming progress information
    Given a NovusPack package
    And an open FileStream
    When Progress method is called
    Then bytesRead is returned indicating bytes read so far
    And totalBytes is returned indicating total bytes to read
    And readSpeed is returned indicating current read speed in bytes per second
    And elapsed is returned indicating time elapsed since stream creation

  @REQ-STREAM-031 @happy
  Scenario: Progress enables progress percentage calculation
    Given a NovusPack package
    And an open FileStream
    When Progress method is called
    Then bytesRead and totalBytes enable percentage calculation
    And progress percentage is calculated as (bytesRead / totalBytes) * 100
    And progress information enables progress bars

  @REQ-STREAM-031 @happy
  Scenario: Progress tracks read speed
    Given a NovusPack package
    And an open FileStream
    When Progress method is called during reading
    Then readSpeed reflects current bytes per second
    And readSpeed enables performance monitoring
    And readSpeed updates dynamically during reading

  @REQ-STREAM-031 @happy
  Scenario: Progress tracks elapsed time
    Given a NovusPack package
    And an open FileStream
    When Progress method is called
    Then elapsed reflects time since stream creation
    And elapsed enables time tracking
    And elapsed can be used with readSpeed for analysis

  @REQ-STREAM-031 @error
  Scenario: Progress handles closed stream correctly
    Given a NovusPack package
    And a closed FileStream
    When Progress method is called
    Then progress information reflects final state
    And bytesRead equals totalBytes if fully read
    And elapsed reflects total time to completion
