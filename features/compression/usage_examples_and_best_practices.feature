@domain:compression @m2 @REQ-COMPR-064 @REQ-COMPR-122 @spec(api_package_compression.md#13-modern-best-practices)
Feature: Usage Examples and Best Practices

  @REQ-COMPR-064 @happy
  Scenario: Modern best practices define recommended usage patterns
    Given a compression operation
    And modern best practices are followed
    When compression is performed
    Then recommended usage patterns are demonstrated
    And patterns align with industry standards
    And patterns optimize performance

  @REQ-COMPR-064 @happy
  Scenario: Modern best practices align with 7zip, zstd, and tar patterns
    Given a compression operation
    When modern best practices are examined
    Then practices align with 7zip patterns
    And practices align with zstd patterns
    And practices align with tar patterns
    And practices follow industry standards

  @REQ-COMPR-064 @happy
  Scenario: Modern best practices use streaming for large files
    Given a large package file
    When modern best practices are followed
    Then streaming compression is used
    And memory usage is controlled
    And temporary files are used when needed

  @REQ-COMPR-064 @happy
  Scenario: Modern best practices use parallel processing for performance
    Given a compression operation
    When modern best practices are followed
    Then parallel processing is enabled
    And multiple CPU cores are utilized
    And performance is optimized

  @REQ-COMPR-122 @happy
  Scenario: Configuration usage patterns document simple streaming configuration
    Given a compression operation
    When simple streaming configuration pattern is used
    Then basic settings are configured
    And ChunkSize is set to 0 for auto-calculation
    And MaxMemoryUsage is set to 0 for auto-detection
    And TempDir uses system temporary directory

  @REQ-COMPR-122 @happy
  Scenario: Configuration usage patterns document advanced streaming configuration
    Given a compression operation
    When advanced streaming configuration pattern is used
    Then full configuration options are specified
    And chunk size, memory usage, and workers are configured
    And parallel processing and compression level are set
    And memory strategy and adaptive chunking are enabled

  @REQ-COMPR-122 @happy
  Scenario: Configuration usage patterns demonstrate different configuration levels
    Given a compression operation
    When configuration usage patterns are examined
    Then simple usage pattern is documented
    And advanced usage pattern is documented
    And custom configuration pattern is documented
    And patterns serve different use cases
