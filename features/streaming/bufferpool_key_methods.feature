@domain:streaming @m2 @REQ-STREAM-044 @REQ-STREAM-045 @REQ-STREAM-046 @spec(api_streaming.md#23-key-methods)
Feature: BufferPool Key Methods

  @REQ-STREAM-044 @happy
  Scenario: BufferPool key methods provide buffer operations
    Given a NovusPack package
    And a BufferPool with configuration
    When key buffer pool methods are used
    Then NewBufferPool creates buffer pool
    And Get retrieves buffer from pool
    And Put returns buffer to pool
    And GetStats provides buffer pool statistics
    And key methods enable core buffer management

  @REQ-STREAM-044 @happy
  Scenario: NewBufferPool creates buffer pool with configuration
    Given a NovusPack package
    And a BufferConfig
    When NewBufferPool is called
    Then BufferPool is created with configuration
    And pool is ready for buffer operations
    And configuration determines pool behavior

  @REQ-STREAM-044 @happy
  Scenario: Get retrieves buffer from pool
    Given a NovusPack package
    And a BufferPool with buffers
    When Get is called with size
    Then buffer of requested size is returned
    And buffer may be from pool or newly allocated
    And retrieval follows pool eviction policy

  @REQ-STREAM-044 @happy
  Scenario: Put returns buffer to pool
    Given a NovusPack package
    And a BufferPool
    And a buffer to return
    When Put is called with buffer
    Then buffer is returned to pool
    And buffer is available for reuse
    And pool manages buffer lifecycle

  @REQ-STREAM-044 @happy
  Scenario: GetStats provides buffer pool statistics
    Given a NovusPack package
    And a BufferPool with activity
    When GetStats is called
    Then buffer pool statistics are returned
    And statistics include access patterns
    And statistics enable pool optimization

  @REQ-STREAM-045 @happy
  Scenario: BufferPool features define buffer management capabilities
    Given a NovusPack package
    And a BufferPool
    When buffer pool features are examined
    Then size-based pools separate pools for different buffer sizes
    And LRU eviction uses least recently used policy
    And time-based eviction automatically cleans unused buffers
    And memory limits provide configurable total memory usage
    And access tracking provides statistics on usage patterns
    And thread safety enables concurrent access

  @REQ-STREAM-046 @happy
  Scenario: BufferPool additional methods provide extended operations
    Given a NovusPack package
    And a BufferPool
    When additional buffer pool methods are used
    Then TotalSize returns total size of all buffers
    And SetMaxTotalSize sets maximum total size
    And additional methods enable monitoring and configuration
    And extended operations complement core operations
