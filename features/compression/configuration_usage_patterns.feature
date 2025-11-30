@domain:compression @m2 @REQ-COMPR-122 @spec(api_package_compression.md#515-configuration-usage-patterns)
Feature: Configuration Usage Patterns

  @REQ-COMPR-122 @happy
  Scenario: Simple usage pattern uses basic settings only
    Given compression operations requiring basic configuration
    When simple usage pattern is applied
    Then ChunkSize is set to 0 for auto-calculate
    And MaxMemoryUsage is set to 0 for auto-detect
    And TempDir is set to empty string for system temp
    And configuration uses intelligent defaults

  @REQ-COMPR-122 @happy
  Scenario: Simple usage pattern relies on automatic detection
    Given StreamConfig with simple usage pattern
    When compression operation is performed
    Then system auto-calculates chunk size
    And system auto-detects memory usage
    And system uses default temporary directory
    And intelligent defaults optimize configuration

  @REQ-COMPR-122 @happy
  Scenario: Advanced usage pattern provides full configuration
    Given compression operations requiring fine-tuned control
    When advanced usage pattern is applied
    Then ChunkSize can be set to specific value (e.g. 1GB)
    And MaxMemoryUsage can be set to specific limit (e.g. 8GB)
    And UseParallelProcessing can be enabled
    And MaxWorkers can be configured
    And CompressionLevel can be specified
    And UseSolidCompression can be enabled
    And MemoryStrategy can be set
    And AdaptiveChunking can be enabled

  @REQ-COMPR-122 @happy
  Scenario: Advanced usage pattern enables maximum optimization
    Given StreamConfig with advanced usage pattern
    When compression operation is performed
    Then full configuration is applied
    And performance is optimized
    And resource usage is controlled precisely
    And advanced features enhance operation
