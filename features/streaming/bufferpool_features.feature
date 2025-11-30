@domain:streaming @m2 @REQ-STREAM-045 @spec(api_streaming.md#24-features)
Feature: BufferPool Features

  @REQ-STREAM-045 @happy
  Scenario: BufferPool features define buffer management capabilities
    Given a NovusPack package
    And a BufferPool with configuration
    When buffer pool features are examined
    Then size-based pools separate pools for different buffer sizes
    And LRU eviction uses least recently used eviction policy
    And time-based eviction automatically cleans up unused buffers
    And memory limits provide configurable total memory usage limits
    And access tracking provides statistics on buffer usage patterns
    And thread safety enables concurrent access with proper synchronization

  @REQ-STREAM-045 @happy
  Scenario: Size-based pools organize buffers by size
    Given a NovusPack package
    And a BufferPool with configuration
    When buffers of different sizes are requested
    Then buffers are organized into separate pools by size
    And buffer retrieval is optimized for size matching
    And size-based organization improves allocation efficiency

  @REQ-STREAM-045 @happy
  Scenario: LRU eviction policy manages buffer lifecycle
    Given a NovusPack package
    And a BufferPool with LRU eviction policy
    When memory limit is reached
    Then least recently used buffers are evicted first
    And eviction frees memory for new buffer allocations
    And LRU policy maintains frequently used buffers in memory

  @REQ-STREAM-045 @happy
  Scenario: Time-based eviction automatically cleans unused buffers
    Given a NovusPack package
    And a BufferPool with time-based eviction configured
    When buffers remain unused beyond eviction timeout
    Then unused buffers are automatically evicted
    And eviction timeout prevents memory leaks
    And time-based cleanup maintains memory efficiency

  @REQ-STREAM-045 @happy
  Scenario: Memory limits prevent excessive memory usage
    Given a NovusPack package
    And a BufferPool with memory limit configured
    When buffer allocations approach memory limit
    Then eviction policies are triggered automatically
    And memory usage stays within configured limit
    And memory management prevents resource exhaustion

  @REQ-STREAM-045 @happy
  Scenario: Access tracking provides buffer usage statistics
    Given a NovusPack package
    And a BufferPool with access tracking enabled
    When buffers are accessed multiple times
    Then access count statistics are maintained per buffer
    And access patterns enable optimization decisions
    And statistics provide insights into buffer usage

  @REQ-STREAM-045 @happy
  Scenario: Thread safety enables concurrent buffer operations
    Given a NovusPack package
    And a BufferPool with concurrent access
    When multiple goroutines access buffer pool simultaneously
    Then all operations are thread-safe
    And concurrent access is properly synchronized
    And thread safety prevents data races
