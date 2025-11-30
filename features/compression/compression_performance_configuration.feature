@domain:compression @m2 @REQ-COMPR-048 @spec(api_package_compression.md#11352-performance-configuration)
Feature: Compression Performance Configuration

  @REQ-COMPR-048 @happy
  Scenario: Performance configuration limits temporary file sizes
    Given a compression operation with temporary files
    When MaxTempFileSize is configured
    Then individual temporary file sizes are limited
    And temporary file size constraints are applied
    And performance configuration controls resource usage

  @REQ-COMPR-048 @happy
  Scenario: Performance configuration optimizes memory usage
    Given compression operations requiring memory management
    When performance configuration is applied
    Then memory usage is optimized
    And temporary file size limits reduce memory pressure
    And resource consumption is controlled

  @REQ-COMPR-048 @happy
  Scenario: Performance configuration supports streaming operations
    Given advanced compression streaming features
    When performance configuration is used
    Then streaming operations are optimized
    And temporary file management is configured
    And performance settings enhance streaming
