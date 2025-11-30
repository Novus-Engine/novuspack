@domain:streaming @m2 @REQ-STREAM-048 @spec(api_streaming.md#252-totalsize-returns)
Feature: Buffer Pool Total Size

  @REQ-STREAM-048 @happy
  Scenario: TotalSize returns total buffer pool size
    Given a NovusPack package
    And a BufferPool with buffers
    When TotalSize method is called
    Then total size of all buffers in pool is returned
    And size reflects current memory usage
    And size includes all buffers currently in pool

  @REQ-STREAM-048 @happy
  Scenario: TotalSize enables memory monitoring
    Given a NovusPack package
    And a BufferPool
    When TotalSize is monitored
    Then current memory usage can be tracked
    And memory usage can be compared to limits
    And TotalSize enables memory management

  @REQ-STREAM-048 @happy
  Scenario: TotalSize updates as buffers are managed
    Given a NovusPack package
    And a BufferPool
    When buffers are added or removed
    Then TotalSize reflects current state
    And TotalSize increases when buffers are added
    And TotalSize decreases when buffers are evicted

  @REQ-STREAM-048 @error
  Scenario: TotalSize handles empty pool correctly
    Given a NovusPack package
    And an empty BufferPool
    When TotalSize is called
    Then zero is returned
    And TotalSize reflects empty pool state
