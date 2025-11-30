@domain:streaming @m2 @REQ-STREAM-005 @spec(api_streaming.md#3-streaming-concurrency-patterns)
Feature: Backpressure handling

  @happy
  Scenario: Backpressure prevents memory overflow
    Given a streaming operation with slow consumer
    When backpressure is applied
    Then memory usage is controlled
    And producer slows down
    And system remains stable

  @happy
  Scenario: Backpressure mechanism is automatic
    Given streaming operations
    When memory pressure increases
    Then backpressure is applied automatically
    And operations are throttled appropriately
