@domain:streaming @m2 @REQ-STREAM-051 @spec(api_streaming.md#255-example-usage)
Feature: SetMaxTotalSize Example Usage

  @REQ-STREAM-051 @happy
  Scenario: BufferPool example usage demonstrates buffer management
    Given a NovusPack package
    When buffer pool is used as in example
    Then NewBufferPool creates buffer pool with configuration
    And TotalSize returns current pool size for monitoring
    And SetMaxTotalSize adjusts memory limit dynamically
    And example demonstrates proper buffer pool usage patterns

  @REQ-STREAM-051 @happy
  Scenario: BufferPool creation and configuration
    Given a NovusPack package
    And a BufferConfig
    When NewBufferPool is called with config
    Then BufferPool is created
    And pool is ready for buffer operations
    And configuration determines pool behavior

  @REQ-STREAM-051 @happy
  Scenario: TotalSize enables memory monitoring
    Given a NovusPack package
    And a BufferPool with allocated buffers
    When TotalSize is called to check memory usage
    Then current pool size is returned in bytes
    And size can be printed for monitoring
    And monitoring enables proactive memory management

  @REQ-STREAM-051 @happy
  Scenario: SetMaxTotalSize enables dynamic limit adjustment
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with new limit
    Then memory limit is adjusted to specified value
    And limit adjustment takes effect immediately
    And example demonstrates dynamic configuration

  @REQ-STREAM-051 @happy
  Scenario: Memory usage monitoring loop
    Given a NovusPack package
    And a BufferPool with configured limit
    When memory usage is monitored in loop
    Then TotalSize is checked periodically
    And high memory usage can be detected
    And monitoring enables alerting on memory pressure

  @REQ-STREAM-051 @error
  Scenario: Example usage handles error conditions
    Given a NovusPack package
    And a BufferPool with error condition
    When buffer pool operations encounter errors
    Then errors are handled gracefully
    And error handling follows example patterns
    And structured errors provide diagnostic information
