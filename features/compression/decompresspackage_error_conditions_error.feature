@domain:compression @m2 @REQ-COMPR-115 @spec(api_package_compression.md#424-decompresspackage-error-conditions)
Feature: DecompressPackage Error Conditions

  @REQ-COMPR-115 @error
  Scenario: DecompressPackage returns error when package is not compressed
    Given an open NovusPack package
    And package is not compressed
    And a valid context
    When DecompressPackage is called
    Then validation error is returned
    And error indicates package is not compressed
    And error follows structured error format

  @REQ-COMPR-115 @error
  Scenario: DecompressPackage returns error on decompression failure
    Given an open NovusPack package
    And package has corrupted compressed data
    And a valid context
    When DecompressPackage is called
    Then compression error is returned
    And error indicates decompression operation failed
    And error follows structured error format

  @REQ-COMPR-115 @error
  Scenario: DecompressPackage handles context cancellation
    Given an open NovusPack package
    And package is compressed
    And a cancelled context
    When DecompressPackage is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-COMPR-115 @error
  Scenario: DecompressPackage handles corrupted compressed data
    Given an open NovusPack package
    And compressed data is corrupted
    And a valid context
    When DecompressPackage is called
    Then corruption error is returned
    And error indicates compressed data is corrupted
    And error follows structured error format

  @REQ-COMPR-115 @error
  Scenario: DecompressPackage error conditions use structured error system
    Given an open NovusPack package
    And an error condition occurs
    When DecompressPackage is called
    Then structured error is returned
    And error provides categorization
    And error provides context information
    And error provides debugging capabilities
