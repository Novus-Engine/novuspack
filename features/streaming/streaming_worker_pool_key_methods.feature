@domain:streaming @m2 @REQ-STREAM-057 @REQ-STREAM-058 @spec(api_streaming.md#33-key-methods)
Feature: Streaming Worker Pool Key Methods

  @REQ-STREAM-057 @happy
  Scenario: Streaming worker pool key methods provide Start method
    Given an open NovusPack package
    And a valid context
    And streaming worker pool instance
    When Start is called
    Then worker pool is initialized and started
    And workers are ready to process streaming jobs
    And context supports cancellation
    And context supports timeout handling

  @REQ-STREAM-057 @happy
  Scenario: Streaming worker pool key methods provide Stop method
    Given an open NovusPack package
    And a valid context
    And running streaming worker pool
    When Stop is called
    Then worker pool gracefully shuts down
    And workers finish current jobs
    And resources are properly cleaned up

  @REQ-STREAM-057 @happy
  Scenario: Streaming worker pool key methods provide SubmitStreamingJob
    Given an open NovusPack package
    And a valid context
    And streaming worker pool
    And streaming job
    When SubmitStreamingJob is called
    Then job is submitted to worker pool
    And job is queued for processing
    And job processing respects context cancellation

  @REQ-STREAM-057 @happy
  Scenario: Streaming worker pool key methods provide ProcessStreamsConcurrently
    Given an open NovusPack package
    And a valid context
    And multiple file streams
    And processor function
    And streaming concurrency config
    When ProcessStreamsConcurrently is called
    Then multiple streams are processed concurrently
    And processing respects max concurrency limits
    And errors are collected and returned

  @REQ-STREAM-057 @happy
  Scenario: Streaming worker pool key methods provide GetStreamingStats
    Given an open NovusPack package
    And a valid context
    And streaming worker pool
    When GetStreamingStats is called
    Then current streaming worker pool statistics are returned
    And statistics include worker information
    And statistics include job processing information

  @REQ-STREAM-058 @happy
  Scenario: Streaming worker pool features provide concurrent stream processing
    Given an open NovusPack package
    And streaming worker pool
    When streaming worker pool features are examined
    Then concurrent stream processing enables multiple streams processed simultaneously
    And thread-safe operations ensure thread safety
    And resource management provides intelligent allocation and cleanup
    And adaptive chunking provides dynamic chunk size based on system load
