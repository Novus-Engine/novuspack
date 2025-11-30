@domain:streaming @m2 @REQ-STREAM-025 @spec(api_streaming.md#1511-purpose)
Feature: Buffer Pool Information Purpose and Usage

  @REQ-STREAM-025 @happy
  Scenario: Buffer pool information purpose provides buffer information access
    Given an open NovusPack package
    And a valid context
    And buffer pool instance
    When buffer pool information purpose is examined
    Then purpose provides buffer information access
    And buffer statistics are accessible
    And buffer pool state information is available

  @REQ-STREAM-025 @happy
  Scenario: Buffer pool information supports buffer management
    Given an open NovusPack package
    And a valid context
    And buffer pool with buffers
    When buffer pool information is accessed
    Then buffer pool management is supported
    And buffer state monitoring is enabled
    And buffer pool optimization is facilitated
