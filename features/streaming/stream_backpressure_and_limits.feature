@domain:streaming @m3 @spec(api_streaming.md#3-streaming-concurrency-patterns)
Feature: Stream Backpressure and Limits

  @REQ-STREAM-003 @performance
  Scenario: Stream honors max concurrency and limits
    Given a NovusPack package
    And a configured maximum concurrency limit
    When concurrent reads exceed the limit
    Then backpressure is applied
    And concurrency limits are respected
    And resource exhaustion is prevented
    And streaming operations remain controlled

  @REQ-STREAM-003 @performance
  Scenario: Backpressure handling prevents resource exhaustion
    Given a NovusPack package
    And streaming concurrency configuration
    When concurrent operations exceed limits
    Then backpressure mechanisms activate
    And operations are throttled appropriately
    And system resources are protected

  @REQ-STREAM-003 @performance
  Scenario: Streaming honors resource limits
    Given a NovusPack package
    And configured resource limits
    When streaming operations are performed
    Then max concurrency limits are enforced
    And memory limits are respected
    And limits prevent resource exhaustion

  @REQ-STREAM-003 @error
  Scenario: Backpressure handles limit violations gracefully
    Given a NovusPack package
    When resource limits are exceeded
    Then appropriate error or throttling occurs
    And operations are handled gracefully
    And error follows structured error format
