@domain:compression @m2 @REQ-COMPR-110 @spec(api_package_compression.md#414-compresspackage-error-conditions)
Feature: CompressPackage Error Conditions

  @REQ-COMPR-110 @error
  Scenario: CompressPackage returns error when package is already signed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a compression type
    When CompressPackage is called
    Then security error is returned
    And error indicates package is already signed
    And error follows structured error format

  @REQ-COMPR-110 @error
  Scenario: CompressPackage returns error for invalid compression type
    Given an open NovusPack package
    And an invalid compression type is provided
    And a valid context
    When CompressPackage is called
    Then validation error is returned
    And error indicates invalid compression type
    And error follows structured error format

  @REQ-COMPR-110 @error
  Scenario: CompressPackage returns error when package is already compressed with different type
    Given an open NovusPack package
    And package is already compressed with Zstandard
    And different compression type is requested
    And a valid context
    When CompressPackage is called
    Then validation error is returned
    And error indicates package already compressed with different type
    And error follows structured error format

  @REQ-COMPR-110 @error
  Scenario: CompressPackage returns error on context cancellation
    Given an open NovusPack package
    And a cancelled context
    And a compression type
    When CompressPackage is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-COMPR-110 @error
  Scenario: CompressPackage returns error on compression operation failure
    Given an open NovusPack package
    And compression operation fails
    And a valid context
    And a compression type
    When CompressPackage is called
    Then compression error is returned
    And error indicates compression operation failure
    And error follows structured error format

  @REQ-COMPR-110 @error
  Scenario: CompressPackage error conditions use structured error system
    Given an open NovusPack package
    And an error condition occurs
    When CompressPackage is called
    Then structured error is returned
    And error provides categorization
    And error provides context information
    And error provides debugging capabilities
