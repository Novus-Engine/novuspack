@domain:compression @m2 @REQ-COMPR-058 @spec(api_package_compression.md#1212-common-error-conditions-decompression-errors)
Feature: Common Error Conditions Decompression Errors

  @REQ-COMPR-058 @error
  Scenario: Validation error when package is not compressed
    Given a NovusPack package that is not compressed
    When decompression operation is attempted
    Then validation error is returned
    And error indicates package is not compressed
    And error prevents invalid decompression

  @REQ-COMPR-058 @error
  Scenario: Validation error for invalid compressed data format
    Given a NovusPack package
    And package has invalid compressed data format
    When decompression operation is attempted
    Then validation error is returned
    And error indicates invalid compressed data format
    And error provides format details

  @REQ-COMPR-058 @error
  Scenario: Compression error when decompression operation fails
    Given a compressed NovusPack package
    When decompression operation fails
    Then compression error is returned
    And error indicates decompression operation failure
    And error provides details about failure

  @REQ-COMPR-058 @error
  Scenario: Compression error for algorithm-specific decompression failures
    Given a compressed NovusPack package
    And decompression algorithm encounters specific failure
    When decompression operation is attempted
    Then compression error is returned
    And error indicates algorithm-specific failure
    And error provides algorithm context

  @REQ-COMPR-058 @error
  Scenario: Corruption error when compressed data is corrupted
    Given a NovusPack package with corrupted compressed data
    When decompression operation is attempted
    Then corruption error is returned
    And error indicates compressed data is corrupted
    And error provides corruption details

  @REQ-COMPR-058 @error
  Scenario: Corruption error when checksum validation fails
    Given a compressed NovusPack package
    And checksum validation fails during decompression
    When decompression operation is attempted
    Then corruption error is returned
    And error indicates checksum validation failure
    And error indicates data integrity issue
