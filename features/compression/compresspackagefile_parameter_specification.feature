@domain:compression @m2 @REQ-COMPR-128 @spec(api_package_compression.md#612-compresspackagefile-parameters)
Feature: CompressPackageFile Parameter Specification

  @REQ-COMPR-128 @happy
  Scenario: CompressPackageFile accepts context parameter
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When CompressPackageFile is called with context
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-COMPR-128 @happy
  Scenario: CompressPackageFile accepts path parameter for target file
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When CompressPackageFile is called with path
    Then path parameter specifies target file location
    And path is validated before writing
    And compressed package is written to specified path

  @REQ-COMPR-128 @happy
  Scenario: CompressPackageFile accepts compressionType parameter
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type 1-3
    And an overwrite flag
    When CompressPackageFile is called with compression type
    Then compression type parameter specifies algorithm
    And specified compression algorithm is used
    And compression type is validated

  @REQ-COMPR-128 @happy
  Scenario: CompressPackageFile accepts overwrite parameter
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When CompressPackageFile is called with overwrite flag
    Then overwrite parameter controls file handling
    And overwrite true allows file replacement
    And overwrite false prevents overwriting existing files

  @REQ-COMPR-128 @error
  Scenario: CompressPackageFile validates path parameter
    Given an open NovusPack package
    And a valid context
    And an invalid file path
    And a compression type
    When CompressPackageFile is called
    Then validation error is returned
    And error indicates invalid path
    And error follows structured error format

  @REQ-COMPR-128 @error
  Scenario: CompressPackageFile handles context cancellation
    Given an open NovusPack package
    And a cancelled context
    And a target file path
    And a compression type
    When CompressPackageFile is called
    Then context cancellation error is returned
    And error type is context cancellation
