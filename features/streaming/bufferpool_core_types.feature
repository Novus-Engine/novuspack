@domain:streaming @m2 @REQ-STREAM-041 @spec(api_streaming.md#22-core-types)
Feature: BufferPool Core Types

  @REQ-STREAM-041 @happy
  Scenario: Core types define BufferPool struct
    Given an open NovusPack package
    And a valid context
    And buffer management system
    When BufferPool struct is examined
    Then struct contains buffers map (buffer IDs to buffer data)
    And struct contains lastUsed map (last access time for each buffer)
    And struct contains accessCount map (access count for each buffer)
    And struct provides buffer pool functionality

  @REQ-STREAM-041 @happy
  Scenario: Core types define BufferConfig struct
    Given an open NovusPack package
    And a valid context
    And buffer management system
    When BufferConfig struct is examined
    Then struct contains MaxTotalSize field (maximum total size of all buffers)
    And struct contains MaxBufferSize field (maximum size of a single buffer)
    And struct contains EvictionPolicy field ("lru" or "time" eviction policy)
    And struct contains EvictionTimeout field (time after which unused buffers are evicted)

  @REQ-STREAM-041 @happy
  Scenario: Core types define generic BufferPool
    Given an open NovusPack package
    And a valid context
    And buffer management system
    When generic BufferPool is examined
    Then BufferPool manages buffers of any type
    And buffers map stores buffers by ID
    And lastUsed map tracks access times
    And accessCount map tracks access frequency

  @REQ-STREAM-041 @happy
  Scenario: Core types provide type-safe buffer management
    Given an open NovusPack package
    And a valid context
    And buffer management system
    When core types are used
    Then type-safe buffer management is provided
    And generic BufferPool enables flexible buffer types
    And BufferConfig enables buffer pool configuration
