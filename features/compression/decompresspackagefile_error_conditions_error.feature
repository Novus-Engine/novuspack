@domain:compression @m2 @REQ-COMPR-135 @spec(api_package_compression.md#624-decompresspackagefile-error-conditions)
Feature: DecompressPackageFile Error Conditions

  @REQ-COMPR-135 @error
  Scenario: DecompressPackageFile returns error when package is not compressed
    Given an open NovusPack package
    And package is not compressed
    And a valid context
    And a target file path
    And an overwrite flag
    When DecompressPackageFile is called
    Then validation error is returned
    And error indicates package is not compressed
    And error follows structured error format

  @REQ-COMPR-135 @error
  Scenario: DecompressPackageFile returns error when file exists and overwrite is false
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path that exists
    And overwrite flag is set to false
    When DecompressPackageFile is called
    Then validation error is returned
    And error indicates file already exists
    And error follows structured error format

  @REQ-COMPR-135 @error
  Scenario: DecompressPackageFile returns error for I/O operation failures
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path with I/O errors
    When DecompressPackageFile is called
    Then I/O error is returned
    And error indicates I/O operation failure
    And error follows structured error format

  @REQ-COMPR-135 @error
  Scenario: DecompressPackageFile handles context cancellation
    Given an open NovusPack package
    And package is compressed
    And a cancelled context
    And a target file path
    When DecompressPackageFile is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-COMPR-135 @error
  Scenario: DecompressPackageFile handles decompression operation failures
    Given an open NovusPack package
    And package has corrupted compressed data
    And a valid context
    And a target file path
    When DecompressPackageFile is called
    Then compression error is returned
    And error indicates decompression failure
    And error follows structured error format
