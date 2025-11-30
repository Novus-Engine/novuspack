@domain:streaming @m2 @REQ-STREAM-011 @spec(api_streaming.md#2-buffer-management-system)
Feature: Buffer pool operations

  @happy
  Scenario: NewBufferPool creates buffer pool
    Given buffer pool configuration
    When NewBufferPool is called with configuration
    Then BufferPool is created
    And pool is ready for use
    And pool is not closed

  @happy
  Scenario: Get gets buffer from pool
    Given a BufferPool
    When Get is called
    Then buffer is returned from pool
    And buffer size matches configuration
    And buffer is ready for use

  @happy
  Scenario: Put returns buffer to pool
    Given a BufferPool with borrowed buffer
    When Put is called with buffer
    Then buffer is returned to pool
    And buffer is available for reuse
    And pool statistics are updated

  @happy
  Scenario: GetStats gets buffer pool statistics
    Given a BufferPool that has been used
    When GetStats is called
    Then buffer pool statistics are returned
    And active buffers count is included
    And total buffers count is included
    And pool utilization is included

  @happy
  Scenario: TotalSize returns total size of all buffers
    Given a BufferPool with buffers
    When TotalSize is called
    Then total size of all buffers is returned
    And size reflects pool state

  @happy
  Scenario: SetMaxTotalSize sets maximum total size
    Given a BufferPool
    When SetMaxTotalSize is called with max size
    Then maximum total size is set
    And pool respects maximum size limit

  @happy
  Scenario: Close closes buffer pool
    Given an open BufferPool
    When Close is called
    Then pool is closed
    And all buffers are released
    And resources are cleaned up

  @error
  Scenario: Get fails if pool is closed
    Given a closed BufferPool
    When Get is called
    Then structured validation error is returned

  @error
  Scenario: SetMaxTotalSize fails with invalid size
    Given a BufferPool
    When SetMaxTotalSize is called with invalid size
    Then structured validation error is returned

  @REQ-STREAM-013 @REQ-STREAM-015 @error
  Scenario: BufferPool validates buffer size parameter
    Given a BufferPool
    When Get is called with negative buffer size
    Then structured validation error is returned
    And error indicates invalid buffer size

  @REQ-STREAM-013 @REQ-STREAM-017 @error
  Scenario: Buffer pool operations respect context cancellation
    Given a BufferPool in use
    And a cancelled context
    When buffer pool operation is called
    Then structured context error is returned
    And error type is context cancellation
    And pool is closed
    And resources are released
