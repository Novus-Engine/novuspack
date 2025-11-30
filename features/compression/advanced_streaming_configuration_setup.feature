@domain:compression @m2 @REQ-COMPR-047 @spec(api_package_compression.md#11351-configuration-setup)
Feature: Advanced Streaming Configuration Setup

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig is created with intelligent defaults for auto-detection
    Given a compression operation requiring advanced streaming
    When StreamConfig is created with intelligent defaults
    Then system auto-detects optimal values
    And ChunkSize is set to 0 for automatic calculation
    And MaxMemoryUsage is set to 0 for automatic detection

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig uses system temporary directory when TempDir is empty
    Given a compression operation requiring advanced streaming
    When StreamConfig is created with empty TempDir string
    Then system temporary directory is utilized
    And temporary files are stored in system temp location
    And temp directory is automatically managed

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig uses MemoryStrategyBalanced for optimal performance
    Given a compression operation requiring advanced streaming
    When StreamConfig is created with MemoryStrategyBalanced
    Then 50% of available RAM is used for compression
    And optimal performance is achieved
    And system responsiveness is maintained

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig enables AdaptiveChunking for dynamic chunk sizing
    Given a compression operation requiring advanced streaming
    When StreamConfig is created with AdaptiveChunking enabled
    Then system adjusts chunk size based on memory pressure
    And chunk size adapts dynamically during compression
    And memory usage stays within limits

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig performance configuration enables parallel processing
    Given a compression operation requiring advanced streaming
    When StreamConfig is configured with performance settings
    Then UseParallelProcessing can be enabled
    And MaxWorkers is set to 0 for automatic CPU core detection
    And multi-core processing is utilized

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig performance configuration enables disk buffering
    Given a compression operation requiring advanced streaming
    When StreamConfig is configured with UseDiskBuffering enabled
    Then intermediate buffering uses disk when memory limits are reached
    And memory usage is controlled
    And large packages can be processed

  @REQ-COMPR-047 @happy
  Scenario: StreamConfig performance configuration enables temp file cleanup
    Given a compression operation requiring advanced streaming
    When StreamConfig is configured with CleanupTempFiles set to true
    Then temporary files are automatically cleaned up after completion
    And disk space is freed
    And cleanup occurs automatically
