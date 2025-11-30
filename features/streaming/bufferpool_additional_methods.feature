@domain:streaming @m2 @REQ-STREAM-046 @spec(api_streaming.md#25-additional-methods)
Feature: BufferPool Additional Methods

  @REQ-STREAM-046 @happy
  Scenario: BufferPool additional methods provide extended operations
    Given a NovusPack package
    And a BufferPool with configuration
    When additional buffer pool methods are used
    Then TotalSize returns total size of all buffers in pool
    And SetMaxTotalSize sets maximum total size for buffer pool
    And additional methods provide buffer pool management capabilities
    And additional methods enable buffer pool monitoring

  @REQ-STREAM-046 @happy
  Scenario: TotalSize returns current memory usage
    Given a NovusPack package
    And a BufferPool with buffers allocated
    When TotalSize method is called
    Then total size of all buffers currently in pool is returned
    And returned size reflects current memory usage
    And size can be used for memory monitoring

  @REQ-STREAM-046 @happy
  Scenario: SetMaxTotalSize dynamically adjusts memory limit
    Given a NovusPack package
    And a BufferPool with initial memory limit
    When SetMaxTotalSize is called with new limit
    Then memory limit is dynamically adjusted
    And new limit is enforced immediately
    And eviction policies respect updated limit

  @REQ-STREAM-046 @happy
  Scenario: TotalSize enables memory usage monitoring
    Given a NovusPack package
    And a BufferPool with configured limits
    When memory usage is monitored using TotalSize
    Then current memory usage can be tracked
    And memory usage can be compared to limits
    And monitoring enables proactive memory management

  @REQ-STREAM-046 @happy
  Scenario: SetMaxTotalSize triggers eviction when limit exceeded
    Given a NovusPack package
    And a BufferPool with buffers exceeding new limit
    When SetMaxTotalSize is called with lower limit
    Then eviction policies are triggered automatically
    And buffers are evicted until limit is satisfied
    And pool memory usage respects new limit

  @REQ-STREAM-046 @error
  Scenario: SetMaxTotalSize handles invalid parameters
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with invalid maxSize
    Then structured error is returned
    And error indicates invalid parameter
    And error follows structured error format
