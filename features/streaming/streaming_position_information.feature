@domain:streaming @m2 @REQ-STREAM-027 @spec(api_streaming.md#1513-position-returns)
Feature: Streaming Position Information

  @REQ-STREAM-027 @happy
  Scenario: Position returns current stream position
    Given a NovusPack package
    And an open FileStream
    When Position method is called
    Then current read position in stream is returned
    And position reflects bytes read so far
    And position is zero at stream start
    And position increases as data is read

  @REQ-STREAM-027 @happy
  Scenario: Position tracks stream progress
    Given a NovusPack package
    And an open FileStream
    When data is read from stream
    Then Position reflects bytes read
    And position can be used with Size to calculate progress
    And position enables progress percentage calculation

  @REQ-STREAM-027 @happy
  Scenario: Position updates with Seek operations
    Given a NovusPack package
    And an open FileStream
    When Seek is called to change position
    Then Position reflects new position after Seek
    And position matches Seek offset
    And position updates correctly for forward and backward seeks

  @REQ-STREAM-027 @error
  Scenario: Position handles closed stream correctly
    Given a NovusPack package
    And a closed FileStream
    When Position method is called
    Then position reflects final position at closure
    And position remains available after closure
