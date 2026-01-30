@domain:file_mgmt @m2 @REQ-FILEMGMT-190 @spec(api_file_mgmt_compression.md#612-decompressfile-parameters) @spec(api_file_mgmt_file_entry.md#713-decompressfile-parameters)
Feature: DecompressFile Parameter Specification

  @REQ-FILEMGMT-190 @happy
  Scenario: DecompressFile parameters include context and path
    Given an open NovusPack package
    And a valid context
    And a compressed file path in the package
    When DecompressFile is called
    Then context parameter supports cancellation and timeout handling
    And path parameter specifies virtual file path to decompress
    And parameters are validated

  @REQ-FILEMGMT-190 @happy
  Scenario: DecompressFile context supports cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    And a compressed file path
    When DecompressFile is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned

  @REQ-FILEMGMT-190 @happy
  Scenario: DecompressFile context supports timeout handling
    Given an open NovusPack package
    And a context with timeout
    And a compressed file path
    When DecompressFile is called
    And timeout is exceeded
    Then operation respects context timeout
    And structured context timeout error is returned

  @REQ-FILEMGMT-190 @error
  Scenario: DecompressFile handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a compressed file path
    When DecompressFile is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-190 @error
  Scenario: DecompressFile handles file not compressed errors
    Given an open NovusPack package
    And a valid context
    And an uncompressed file path
    When DecompressFile is called with uncompressed file
    Then a structured error is returned
    And error indicates file is not compressed
    And error follows structured error format
