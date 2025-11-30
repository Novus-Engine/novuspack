@domain:compression @m2 @REQ-COMPR-057 @spec(api_package_compression.md#1211-common-error-conditions-compression-errors)
Feature: Common Error Conditions Compression Errors

  @REQ-COMPR-057 @error
  Scenario: Security error when package is already signed
    Given a NovusPack package that is already signed
    When compression operation is attempted
    Then security error is returned
    And error indicates package is already signed
    And error indicates signed packages cannot be compressed

  @REQ-COMPR-057 @error
  Scenario: Validation error for invalid compression algorithm
    Given a NovusPack package
    And an invalid compression algorithm is specified
    When compression operation is attempted
    Then validation error is returned
    And error indicates invalid compression algorithm
    And error specifies valid compression types

  @REQ-COMPR-057 @error
  Scenario: Validation error when package already compressed with different type
    Given a NovusPack package already compressed with one type
    When compression with different type is attempted
    Then validation error is returned
    And error indicates package is already compressed
    And error indicates different compression type conflict

  @REQ-COMPR-057 @error
  Scenario: Compression error when compression operation fails
    Given a NovusPack package
    When compression operation fails
    Then compression error is returned
    And error indicates compression operation failure
    And error provides details about failure

  @REQ-COMPR-057 @error
  Scenario: Compression error for algorithm-specific failures
    Given a NovusPack package
    And compression algorithm encounters specific failure
    When compression operation is attempted
    Then compression error is returned
    And error indicates algorithm-specific failure
    And error provides algorithm context
