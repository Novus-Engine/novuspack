@domain:compression @m2 @REQ-COMPR-052 @spec(api_package_compression.md#11361-custom-configuration)
Feature: Custom Memory Management Configuration

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows explicit chunk size specification
    Given a compression operation with specific memory constraints
    When StreamConfig is created with custom configuration
    And ChunkSize is set to specific value such as 512MB
    Then specified chunk size is used
    And auto-calculation is bypassed
    And chunk processing is controlled

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows custom temporary directory path
    Given a compression operation with specific storage requirements
    When StreamConfig is created with custom TempDir path
    Then custom temporary directory is used for temporary files
    And temporary files are stored at specified location
    And temp directory path is respected

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows explicit memory usage limit
    Given a compression operation with strict memory constraints
    When StreamConfig is created with explicit MaxMemoryUsage such as 1GB
    Then strict memory control is enforced
    And MaxMemoryUsage limit is respected
    And memory usage stays within specified limit

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration uses MemoryStrategyCustom for explicit values
    Given a compression operation with specific memory requirements
    When StreamConfig is created with MemoryStrategyCustom
    Then explicit MaxMemoryUsage value is utilized
    And automatic detection is overridden
    And custom memory limit is applied

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows disabling adaptive chunking
    Given a compression operation requiring predictable behavior
    When StreamConfig is created with AdaptiveChunking disabled
    Then automatic chunk size adjustments are prevented
    And chunk size remains fixed
    And behavior is predictable

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows explicit buffer pool size
    Given a compression operation with specific memory constraints
    When StreamConfig is created with explicit BufferPoolSize
    Then buffer pool size is set to specified number of buffers
    And memory usage is predictable
    And buffer pool is controlled

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows limiting temporary file size
    Given a compression operation with storage constraints
    When StreamConfig is created with MaxTempFileSize limit
    Then individual temporary file sizes are limited
    And file rotation occurs when limit is reached
    And storage usage is controlled

  @REQ-COMPR-052 @happy
  Scenario: Custom configuration allows explicit worker count
    Given a compression operation with specific performance requirements
    When StreamConfig is created with explicit MaxWorkers number
    Then worker count is set to specified number
    And auto-detection is bypassed
    And concurrent workers are controlled
