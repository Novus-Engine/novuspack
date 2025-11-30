@domain:streaming @m2 @REQ-STREAM-052 @spec(api_streaming.md#26-default-configuration)
Feature: Streaming Default Configuration

  @REQ-STREAM-052 @happy
  Scenario: Default configuration provides default buffer settings
    Given a NovusPack package
    When default buffer configuration is used
    Then default configuration provides sensible defaults
    And MaxTotalSize defaults to 1GB
    And MaxBufferSize defaults to 1MB
    And EvictionPolicy defaults to LRU
    And EvictionTimeout defaults to 5 minutes

  @REQ-STREAM-052 @happy
  Scenario: Default configuration enables easy buffer pool setup
    Given a NovusPack package
    When DefaultBufferConfig is used
    Then buffer pool can be created without configuration
    And default settings work for most use cases
    And configuration can be customized if needed

  @REQ-STREAM-052 @happy
  Scenario: Default configuration provides balanced performance
    Given a NovusPack package
    When default configuration is used
    Then configuration balances memory usage
    And configuration balances performance
    And configuration provides reasonable defaults for streaming
