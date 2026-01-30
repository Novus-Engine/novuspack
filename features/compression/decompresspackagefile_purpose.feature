@domain:compression @m2 @REQ-COMPR-132 @spec(api_package_compression.md#621-purpose)
Feature: DecompressPackageFile purpose

  @REQ-COMPR-132 @happy
  Scenario: DecompressPackageFile exists to decompress from a package file
    Given a compressed package file on disk
    When DecompressPackageFile is called
    Then the package is decompressed according to configuration
    And the result is written to an output file path
    And the output is a valid uncompressed package file
    And the operation supports large inputs via streaming where applicable
    And the purpose aligns with documented file-based decompression guidance

