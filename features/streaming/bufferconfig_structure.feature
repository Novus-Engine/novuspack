@domain:streaming @m2 @REQ-STREAM-043 @spec(api_streaming.md#222-bufferconfig-struct)
Feature: BufferConfig Structure

  @REQ-STREAM-043 @happy
  Scenario: BufferConfig struct provides buffer configuration structure
    Given a NovusPack package
    When BufferConfig struct is used
    Then struct contains MaxTotalSize field for maximum total buffer size
    And struct contains MaxBufferSize field for maximum single buffer size
    And struct contains EvictionPolicy field for eviction policy
    And struct contains EvictionTimeout field for eviction timeout
    And struct provides complete buffer configuration

  @REQ-STREAM-043 @happy
  Scenario: MaxTotalSize field configures maximum total buffer size
    Given a NovusPack package
    And a BufferConfig
    When MaxTotalSize field is set
    Then maximum total size of all buffers is configured
    And value is specified in bytes
    And value limits total memory usage of buffer pool

  @REQ-STREAM-043 @happy
  Scenario: MaxBufferSize field configures maximum single buffer size
    Given a NovusPack package
    And a BufferConfig
    When MaxBufferSize field is set
    Then maximum size of single buffer is configured
    And value is specified in bytes
    And value prevents allocation of excessively large buffers

  @REQ-STREAM-043 @happy
  Scenario: EvictionPolicy field configures eviction policy
    Given a NovusPack package
    And a BufferConfig
    When EvictionPolicy field is set
    Then eviction policy is configured as "lru" or "time"
    And policy determines how buffers are evicted
    And policy affects buffer lifecycle management

  @REQ-STREAM-043 @happy
  Scenario: EvictionTimeout field configures eviction timeout
    Given a NovusPack package
    And a BufferConfig
    When EvictionTimeout field is set
    Then time after which unused buffers are evicted is configured
    And value is specified as time.Duration
    And timeout enables automatic cleanup of unused buffers

  @REQ-STREAM-043 @error
  Scenario: BufferConfig struct validates configuration values
    Given a NovusPack package
    And a BufferConfig with invalid values
    When BufferConfig is used to create BufferPool
    Then validation error is returned
    And error indicates invalid configuration field
    And error follows structured error format
