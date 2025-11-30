@domain:streaming @m2 @REQ-STREAM-053 @spec(api_streaming.md#261-defaultbufferconfig-bufferconfig)
Feature: DefaultBufferConfig Buffer Configuration

  @REQ-STREAM-053 @happy
  Scenario: DefaultBufferConfig provides default buffer configuration
    Given a NovusPack package
    When DefaultBufferConfig function is called
    Then BufferConfig structure is returned
    And MaxTotalSize is set to 1GB (1 << 30)
    And MaxBufferSize is set to 1MB (1 << 20)
    And EvictionPolicy is set to "lru"
    And EvictionTimeout is set to 5 minutes

  @REQ-STREAM-053 @happy
  Scenario: DefaultBufferConfig provides ready-to-use configuration
    Given a NovusPack package
    When DefaultBufferConfig is used to create buffer pool
    Then buffer pool uses default settings
    And default settings are appropriate for most scenarios
    And configuration enables immediate buffer pool usage

  @REQ-STREAM-053 @happy
  Scenario: DefaultBufferConfig values are documented and consistent
    Given a NovusPack package
    When DefaultBufferConfig is examined
    Then MaxTotalSize of 1GB provides reasonable memory limit
    And MaxBufferSize of 1MB provides reasonable buffer size
    And LRU eviction policy provides efficient buffer reuse
    And 5 minute timeout provides reasonable cleanup interval
