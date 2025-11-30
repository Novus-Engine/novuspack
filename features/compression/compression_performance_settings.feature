@domain:compression @m2 @REQ-COMPR-053 @spec(api_package_compression.md#11362-performance-settings)
Feature: Compression Performance Settings

  @REQ-COMPR-053 @happy
  Scenario: UseParallelProcessing enables multi-core utilization
    Given compression operations
    When UseParallelProcessing is enabled
    Then multi-core CPU utilization is enabled
    And parallel processing improves performance
    And compression operations use multiple cores

  @REQ-COMPR-053 @happy
  Scenario: MaxWorkers limits concurrent workers
    Given compression operations requiring worker management
    When MaxWorkers is set to specific number
    Then concurrent workers are limited to that number
    And worker pool size is controlled
    And resource usage is managed

  @REQ-COMPR-053 @happy
  Scenario: CompressionLevel provides consistent compression behavior
    Given compression operations
    When specific CompressionLevel is specified
    Then compression behavior is consistent
    And compression ratio matches level
    And performance characteristics are predictable

  @REQ-COMPR-053 @happy
  Scenario: Performance settings optimize compression workflow
    Given custom memory execution requirements
    When performance settings are configured
    Then compression workflow is optimized
    And resource usage is controlled
    And performance matches requirements
