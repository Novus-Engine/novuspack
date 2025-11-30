@domain:streaming @m2 @REQ-STREAM-026 @spec(api_streaming.md#1512-size-returns)
Feature: Streaming Size Information

  @REQ-STREAM-026 @happy
  Scenario: Size returns file stream size information
    Given a NovusPack package
    And an open FileStream
    When Size method is called
    Then total size of stream in bytes is returned
    And size reflects total file size
    And size is available at any time
    And size remains constant during stream lifetime

  @REQ-STREAM-026 @happy
  Scenario: Size provides stream capacity information
    Given a NovusPack package
    And an open FileStream
    When Size method is called
    Then size enables progress calculation
    And size enables percentage calculation
    And size enables remaining bytes calculation

  @REQ-STREAM-026 @error
  Scenario: Size handles closed stream correctly
    Given a NovusPack package
    And a closed FileStream
    When Size method is called
    Then size information is still available
    And size reflects stream size at closure
