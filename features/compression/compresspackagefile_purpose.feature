@domain:compression @m2 @REQ-COMPR-127 @spec(api_package_compression.md#611-purpose)
Feature: CompressPackageFile purpose

  @REQ-COMPR-127 @happy
  Scenario: CompressPackageFile exists to compress and write a package file
    Given an uncompressed package file on disk
    When CompressPackageFile is called
    Then the package is compressed according to configuration
    And the result is written to an output file path
    And the output is a valid compressed package file
    And the operation supports large inputs via streaming where applicable
    And the purpose aligns with documented file-based compression guidance

