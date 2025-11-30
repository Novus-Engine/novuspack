@domain:streaming @m2 @REQ-STREAM-049 @spec(api_streaming.md#253-setmaxtotalsize-parameters)
Feature: SetMaxTotalSize Parameter Specification

  @REQ-STREAM-049 @happy
  Scenario: SetMaxTotalSize parameters define size limit configuration
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called
    Then maxSize parameter specifies maximum total size in bytes
    And parameter defines memory limit for buffer pool
    And configuration enables dynamic memory management

  @REQ-STREAM-049 @happy
  Scenario: MaxSize parameter specifies memory limit
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with maxSize
    Then maxSize parameter defines maximum total size
    And value is specified in bytes
    And limit constrains total memory usage of all buffers

  @REQ-STREAM-049 @happy
  Scenario: MaxSize enables dynamic limit adjustment
    Given a NovusPack package
    And a BufferPool with existing limit
    When SetMaxTotalSize is called with new maxSize
    Then new limit replaces previous configuration
    And adjustment takes effect immediately
    And dynamic adjustment supports runtime configuration

  @REQ-STREAM-049 @happy
  Scenario: MaxSize parameter validation
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with valid maxSize
    Then parameter is accepted and limit is set
    And validation ensures reasonable values
    And valid parameters enable proper operation

  @REQ-STREAM-049 @error
  Scenario: SetMaxTotalSize handles invalid maxSize parameter
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with invalid maxSize
    Then structured error is returned
    And error indicates invalid parameter value
    And error follows structured error format

  @REQ-STREAM-049 @error
  Scenario: SetMaxTotalSize handles edge case parameters
    Given a NovusPack package
    And a BufferPool
    When SetMaxTotalSize is called with edge case maxSize
    Then zero value is handled appropriately
    And very large values are handled appropriately
    And edge cases do not cause panics or undefined behavior
