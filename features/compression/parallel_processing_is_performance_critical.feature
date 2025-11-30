@domain:compression @m2 @REQ-COMPR-067 @spec(api_package_compression.md#1312-parallel-processing-performance-critical)
Feature: Parallel processing is performance critical

  @REQ-COMPR-067 @happy
  Scenario: Multi-core utilization automatically detects and uses CPU cores
    Given compression operations requiring parallel processing
    When parallel processing is enabled
    Then available CPU cores are automatically detected
    And multiple CPU cores are utilized
    And multi-core processing improves performance

  @REQ-COMPR-067 @happy
  Scenario: Worker pool management provides configurable worker count
    Given compression operations with worker pool
    When worker pool is configured
    Then worker count is configurable for optimal performance
    And worker pool manages concurrent workers
    And resource usage is optimized

  @REQ-COMPR-067 @happy
  Scenario: Load balancing distributes chunks across workers
    Given compression operations with multiple workers
    When chunks are processed
    Then chunks are distributed across workers for maximum throughput
    And load is balanced across workers
    And parallel processing efficiency is maximized

  @REQ-COMPR-067 @happy
  Scenario: Memory isolation ensures each worker operates within limits
    Given compression operations with parallel workers
    When workers process chunks
    Then each worker operates within memory limits
    And memory isolation prevents resource conflicts
    And worker memory usage is controlled
