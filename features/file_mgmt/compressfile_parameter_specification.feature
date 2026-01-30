@domain:file_mgmt @m2 @REQ-FILEMGMT-180 @spec(api_file_mgmt_compression.md#711-compressfile-parameters) @spec(api_file_mgmt_file_entry.md#712-compressfile-parameters)
Feature: CompressFile Parameter Specification

  @REQ-FILEMGMT-180 @happy
  Scenario: CompressFile parameters include context, path, and compressionType
    Given an open NovusPack package
    And a valid context
    And a file path in the package
    And a compression type
    When CompressFile is called
    Then context parameter supports cancellation and timeout handling
    And path parameter specifies virtual file path
    And compressionType parameter specifies compression algorithm
    And parameters are validated

  @REQ-FILEMGMT-180 @happy
  Scenario: CompressFile context supports cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    And a file path
    And a compression type
    When CompressFile is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned

  @REQ-FILEMGMT-180 @happy
  Scenario: CompressFile context supports timeout handling
    Given an open NovusPack package
    And a context with timeout
    And a file path
    And a compression type
    When CompressFile is called
    And timeout is exceeded
    Then operation respects context timeout
    And structured context timeout error is returned

  @REQ-FILEMGMT-180 @error
  Scenario: CompressFile handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a file path
    And a compression type
    When CompressFile is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-180 @error
  Scenario: CompressFile handles invalid compression type
    Given an open NovusPack package
    And a valid context
    And a file path
    And an invalid compression type
    When CompressFile is called with invalid compression type
    Then a structured error is returned
    And error indicates unsupported compression type
    And error follows structured error format
