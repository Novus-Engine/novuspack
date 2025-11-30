@domain:streaming @m3 @spec(api_streaming.md#2-buffer-management-system)
Feature: Reuse Buffers Across Operations

  @REQ-STREAM-002 @performance
  Scenario: Buffer pool prevents excessive allocations
    Given a NovusPack package
    And a BufferPool with enabled buffer pooling
    When multiple read operations are processed
    Then buffers are reused from pool
    And number of allocations remains within expected bounds
    And buffer pool prevents excessive memory allocations
    And memory usage is optimized through buffer reuse

  @REQ-STREAM-002 @performance
  Scenario: Buffer pool reuses buffers efficiently
    Given a NovusPack package
    And a BufferPool
    When buffers are requested and returned
    Then buffers are returned to pool after use
    And buffers are reused for subsequent operations
    And buffer reuse reduces allocation overhead

  @REQ-STREAM-002 @performance
  Scenario: Buffer pool manages buffer lifecycle
    Given a NovusPack package
    And a BufferPool
    When buffers are allocated and released
    Then Get retrieves buffer from pool
    And Put returns buffer to pool
    And buffer lifecycle is managed efficiently

  @REQ-STREAM-002 @error
  Scenario: Buffer pool handles resource constraints
    Given a NovusPack package
    And a BufferPool with memory limits
    When memory limit is exceeded
    Then eviction policies are triggered
    And buffer pool manages resources within limits
    And error handling prevents resource exhaustion
