@domain:compression @m2 @REQ-COMPR-130 @spec(api_package_compression.md#614-compresspackagefile-error-conditions)
Feature: CompressPackageFile Error Conditions

  @REQ-COMPR-130 @error
  Scenario: CompressPackageFile returns error when package is already signed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a target file path
    And a compression type
    When CompressPackageFile is called
    Then security error is returned
    And error indicates package is already signed
    And error follows structured error format

  @REQ-COMPR-130 @error
  Scenario: CompressPackageFile returns error for invalid compression type
    Given an open NovusPack package
    And an invalid compression type is provided
    And a valid context
    And a target file path
    When CompressPackageFile is called
    Then validation error is returned
    And error indicates invalid compression type
    And error follows structured error format

  @REQ-COMPR-130 @error
  Scenario: CompressPackageFile returns error when file exists and overwrite is false
    Given an open NovusPack package
    And a valid context
    And a target file path that exists
    And overwrite flag is set to false
    When CompressPackageFile is called
    Then validation error is returned
    And error indicates file already exists
    And error follows structured error format

  @REQ-COMPR-130 @error
  Scenario: CompressPackageFile returns error for I/O operation failures
    Given an open NovusPack package
    And a valid context
    And a target file path with I/O errors
    And a compression type
    When CompressPackageFile is called
    Then I/O error is returned
    And error indicates file operation failure
    And error follows structured error format

  @REQ-COMPR-130 @error
  Scenario: CompressPackageFile handles compression operation failures
    Given an open NovusPack package
    And compression operation fails
    And a valid context
    And a target file path
    When CompressPackageFile is called
    Then compression error is returned
    And error indicates compression failure
    And error follows structured error format
