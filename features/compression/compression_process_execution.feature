@domain:compression @m2 @REQ-COMPR-045 @spec(api_package_compression.md#11342-process)
Feature: Compression Process Execution

  @REQ-COMPR-045 @happy
  Scenario: Process executes advanced streaming compression
    Given advanced streaming compression configuration
    When CompressPackageStream is called with ZSTD compression
    Then streaming compression process executes
    And package is compressed using streaming interface
    And advanced streaming features are utilized

  @REQ-COMPR-045 @happy
  Scenario: Process supports streaming for large packages
    Given a large package requiring streaming compression
    When streaming compression process is executed
    Then package is processed in chunks
    And memory usage is controlled
    And streaming interface handles large data

  @REQ-COMPR-045 @happy
  Scenario: Process applies advanced streaming configuration
    Given streaming compression configuration with advanced features
    When process executes
    Then advanced streaming features are applied
    And performance settings are used
    And process completes successfully
