@domain:compression @m2 @REQ-COMPR-050 @spec(api_package_compression.md#11354-execution)
Feature: Advanced Streaming Compression Execution

  @REQ-COMPR-050 @happy
  Scenario: CompressPackageStream is called with ZSTD compression type and configured settings
    Given an open NovusPack package
    And a StreamConfig with advanced streaming settings configured
    And a valid context
    And ZSTD compression type
    When CompressPackageStream is called
    Then compression uses ZSTD algorithm
    And advanced streaming configuration is applied
    And compression proceeds with configured settings

  @REQ-COMPR-050 @happy
  Scenario: Write is called with CompressionNone after streaming compression
    Given an open NovusPack package
    And package is compressed using CompressPackageStream
    And a valid context
    And an output file path
    When Write is called with CompressionNone
    Then compressed package is written to output file
    And no additional compression is applied
    And already-compressed package is written as-is

  @REQ-COMPR-050 @happy
  Scenario: Execution workflow combines streaming compression and file writing
    Given an open NovusPack package
    And advanced streaming configuration
    And a valid context
    When execution workflow is followed
    Then CompressPackageStream compresses package using streaming
    And Write writes compressed package to file
    And workflow handles large packages efficiently

  @REQ-COMPR-050 @happy
  Scenario: Execution uses configured advanced settings during compression
    Given an open NovusPack package
    And StreamConfig with performance and memory settings
    And a valid context
    When CompressPackageStream is executed
    Then configured chunk size, memory limits, and workers are used
    And adaptive chunking and disk buffering are applied as configured
    And performance optimizations are enabled

  @REQ-COMPR-050 @error
  Scenario: Execution returns error when compression fails
    Given an open NovusPack package
    And a StreamConfig with advanced settings
    And a valid context
    And compression operation fails
    When CompressPackageStream is executed
    Then compression error is returned
    And error indicates compression failure
    And error follows structured error format

  @REQ-COMPR-050 @error
  Scenario: Execution returns error when file write fails
    Given an open NovusPack package
    And package is compressed using CompressPackageStream
    And a valid context
    And output file path causes write failure
    When Write is called
    Then I/O error is returned
    And error indicates file write failure
    And error follows structured error format
