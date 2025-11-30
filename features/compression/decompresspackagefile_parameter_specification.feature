@domain:compression @m2 @REQ-COMPR-133 @spec(api_package_compression.md#622-decompresspackagefile-parameters)
Feature: DecompressPackageFile Parameter Specification

  @REQ-COMPR-133 @happy
  Scenario: DecompressPackageFile accepts context parameter
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    And an overwrite flag
    When DecompressPackageFile is called with context
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-COMPR-133 @happy
  Scenario: DecompressPackageFile accepts path parameter for target file
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    And an overwrite flag
    When DecompressPackageFile is called with path
    Then path parameter specifies target file location
    And path is validated before writing
    And uncompressed package is written to specified path

  @REQ-COMPR-133 @happy
  Scenario: DecompressPackageFile accepts overwrite parameter
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    And an overwrite flag
    When DecompressPackageFile is called with overwrite flag
    Then overwrite parameter controls file handling
    And overwrite true allows file replacement
    And overwrite false prevents overwriting existing files

  @REQ-COMPR-133 @error
  Scenario: DecompressPackageFile validates path parameter
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And an invalid file path
    When DecompressPackageFile is called
    Then validation error is returned
    And error indicates invalid path
    And error follows structured error format

  @REQ-COMPR-133 @error
  Scenario: DecompressPackageFile handles context cancellation
    Given an open NovusPack package
    And package is compressed
    And a cancelled context
    And a target file path
    When DecompressPackageFile is called
    Then context cancellation error is returned
    And error type is context cancellation
