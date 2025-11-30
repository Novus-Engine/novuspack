@domain:streaming @m2 @REQ-STREAM-050 @spec(api_streaming.md#254-behavior)
Feature: SetMaxTotalSize Operation Behavior

  @REQ-STREAM-050 @happy
  Scenario: SetMaxTotalSize behavior defines size limit enforcement
    Given a NovusPack package
    And a BufferPool with buffers
    When SetMaxTotalSize behavior is examined
    Then TotalSize returns current memory usage of pool
    And SetMaxTotalSize dynamically adjusts memory limit
    And when limit is exceeded eviction policies are triggered
    And behavior maintains memory within configured limits

  @REQ-STREAM-050 @happy
  Scenario: TotalSize returns current memory usage
    Given a NovusPack package
    And a BufferPool with allocated buffers
    When TotalSize is called
    Then current memory usage of pool is returned
    And returned value reflects actual buffer memory usage
    And size can be used for memory monitoring

  @REQ-STREAM-050 @happy
  Scenario: SetMaxTotalSize dynamically adjusts memory limit
    Given a NovusPack package
    And a BufferPool with initial configuration
    When SetMaxTotalSize is called with new limit
    Then memory limit is dynamically adjusted
    And adjustment takes effect immediately
    And new limit is enforced on subsequent operations

  @REQ-STREAM-050 @happy
  Scenario: Eviction policies trigger when limit exceeded
    Given a NovusPack package
    And a BufferPool with memory limit
    And buffers exceeding current limit
    When SetMaxTotalSize sets lower limit
    Then eviction policies are triggered automatically
    And buffers are evicted until limit is satisfied
    And eviction respects configured eviction policy

  @REQ-STREAM-050 @happy
  Scenario: Behavior maintains memory within limits
    Given a NovusPack package
    And a BufferPool with memory limit configured
    When buffer operations are performed
    Then memory usage stays within configured limit
    And eviction prevents memory from exceeding limit
    And behavior ensures resource constraints are respected

  @REQ-STREAM-050 @error
  Scenario: Behavior handles edge cases correctly
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with edge case values
    Then behavior handles zero limit appropriately
    And behavior handles very large limits appropriately
    And error handling prevents invalid configurations
