@skip @domain:streaming @m2 @spec(api_streaming.md#1515-example-usage)
Feature: Streaming Usage Examples

# This feature captures usage-oriented scenarios derived from the streaming spec.
# Detailed runnable scenarios live in the dedicated streaming feature files.

  @REQ-STREAM-029 @documentation
  Scenario: Stream information methods report size and position
    Given an open FileStream
    When the caller queries Size and Position
    Then the caller receives the total size and current read position
    And the caller can use the information to display progress

  @REQ-STREAM-033 @documentation
  Scenario: Progress monitoring reports completion and estimated time remaining
    Given an active streaming read operation
    When the caller queries Progress and EstimatedTimeRemaining
    Then the caller receives a progress value between 0 and 1
    And the caller receives an estimated time remaining value when available
