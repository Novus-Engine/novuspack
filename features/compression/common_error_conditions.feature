@domain:compression @m2 @REQ-COMPR-056 @spec(api_package_compression.md#121-common-error-conditions)
Feature: Common Error Conditions

  @REQ-COMPR-056 @error
  Scenario: Compression errors include security error for signed packages
    Given an open NovusPack package
    And package is already signed
    And a valid context
    When compression operation is attempted
    Then security error is returned
    And error indicates package cannot be compressed when signed
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: Compression errors include validation error for invalid algorithm
    Given an open NovusPack package
    And an invalid compression algorithm is specified
    And a valid context
    When compression operation is attempted
    Then validation error is returned
    And error indicates invalid compression algorithm
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: Compression errors include validation error when package already compressed with different type
    Given an open NovusPack package
    And package is already compressed with Zstandard
    And different compression type is requested
    And a valid context
    When compression operation is attempted
    Then validation error is returned
    And error indicates package already compressed with different type
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: Compression errors include compression error for operation failures
    Given an open NovusPack package
    And compression operation fails
    And a valid context
    When compression operation is attempted
    Then compression error is returned
    And error indicates compression operation failed
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: Decompression errors include validation error when package is not compressed
    Given an open NovusPack package
    And package is not compressed
    And a valid context
    When decompression operation is attempted
    Then validation error is returned
    And error indicates package is not compressed
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: Decompression errors include compression error for decompression failures
    Given an open NovusPack package
    And package has corrupted compressed data
    And a valid context
    When decompression operation is attempted
    Then compression error is returned
    And error indicates decompression operation failed
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: Decompression errors include corruption error for corrupted compressed data
    Given an open NovusPack package
    And compressed data is corrupted
    And a valid context
    When decompression operation is attempted
    Then corruption error is returned
    And error indicates compressed data is corrupted
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: File operation errors include validation error when file exists and overwrite is false
    Given an open NovusPack package
    And a target file path that exists
    And overwrite is set to false
    And a valid context
    When file-based compression operation is attempted
    Then validation error is returned
    And error indicates file already exists
    And error follows structured error format

  @REQ-COMPR-056 @error
  Scenario: File operation errors include I/O error for operation failures
    Given an open NovusPack package
    And file operation fails
    And a valid context
    When file-based compression operation is attempted
    Then I/O error is returned
    And error indicates I/O operation failed
    And error follows structured error format
