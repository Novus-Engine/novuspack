@domain:compression @m2 @REQ-COMPR-046 @spec(api_package_compression.md#1135-option-5-advanced-streaming-compression)
Feature: Advanced Streaming Compression

  @REQ-COMPR-046 @happy
  Scenario: Advanced streaming compression is for extremely large packages
    Given compression operations for extremely large packages
    When advanced streaming compression is used
    Then advanced streaming compression handles large packages
    And maximum performance is achieved
    And full configuration options are available

  @REQ-COMPR-046 @happy
  Scenario: Advanced streaming compression uses intelligent defaults with auto-detection
    Given compression operations requiring advanced streaming
    When StreamConfig is created with intelligent defaults
    Then ChunkSize is set to 0 for automatic calculation
    And TempDir is set to empty string for system temp directory
    And MaxMemoryUsage is set to 0 for automatic detection
    And intelligent defaults allow auto-detection of optimal values

  @REQ-COMPR-046 @happy
  Scenario: Advanced streaming compression configures balanced memory strategy
    Given compression operations with advanced streaming
    When MemoryStrategyBalanced is selected
    Then 50% of available RAM is used for optimal performance
    And performance and system stability are balanced
    And good compression speed is achieved

  @REQ-COMPR-046 @happy
  Scenario: Advanced streaming compression enables adaptive chunking
    Given compression operations with advanced streaming
    When AdaptiveChunking is enabled
    Then system adjusts chunk size based on memory pressure
    And adaptive sizing optimizes performance
    And memory usage is dynamically managed

  @REQ-COMPR-046 @happy
  Scenario: Advanced streaming compression provides full configuration options
    Given compression operations requiring advanced features
    When advanced streaming compression is configured
    Then full configuration options are available
    And options align with modern best practices from 7zip, zstd, and tar
    And maximum performance requirements are met
