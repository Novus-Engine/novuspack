@domain:compression @m2 @REQ-COMPR-051 @spec(api_package_compression.md#1136-option-6-custom-memory-management)
Feature: Custom Memory Management

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management supports specific memory constraints
    Given compression operations with specific memory constraints
    When custom memory management is used
    Then specific memory constraints are supported
    And performance requirements are met
    And custom memory setup is configured

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management allows specific chunk size configuration
    Given compression operations requiring controlled chunk processing
    When ChunkSize is set to specific value (e.g. 512MB)
    Then controlled chunk processing is achieved
    And chunk size matches configuration
    And processing is predictable

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management allows custom temp directory
    Given compression operations requiring specific temp directory
    When custom TempDir path is specified
    Then temporary file storage uses custom directory
    And temp file location is controlled
    And temp file management is customized

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management enforces strict memory limits
    Given compression operations requiring strict memory control
    When MaxMemoryUsage is set to specific limit (e.g. 1GB)
    And MemoryStrategyCustom is used
    Then strict memory control is enforced
    And memory usage stays within specified limit
    And explicit MaxMemoryUsage value is utilized

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management disables adaptive chunking for predictability
    Given compression operations requiring predictable memory usage
    When AdaptiveChunking is disabled
    Then automatic chunk size adjustments are prevented
    And memory usage is predictable
    And chunk size remains constant

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management configures buffer pool size
    Given compression operations requiring predictable buffer allocation
    When BufferPoolSize is set to specific number
    Then predictable memory usage is achieved
    And buffer pool size matches configuration
    And buffer allocation is controlled

  @REQ-COMPR-051 @happy
  Scenario: Custom memory management limits temporary file sizes
    Given compression operations requiring temp file size control
    When MaxTempFileSize is configured to limit individual files
    Then individual temporary file sizes are limited
    And temp file rotation occurs when limit exceeded
    And temp file management is controlled
