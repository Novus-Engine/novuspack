@domain:compression @m2 @REQ-COMPR-007 @spec(api_package_compression.md#6-file-based-compression-methods)
Feature: Decompress package file

  @happy
  Scenario: DecompressPackageFile decompresses package to file
    Given a compressed package file
    When DecompressPackageFile is called with output path
    Then package is decompressed
    And decompressed package is written to output path
    And output file is valid package

  @error
  Scenario: DecompressPackageFile fails with invalid path
    Given a compressed package
    When DecompressPackageFile is called with invalid path
    Then structured I/O error is returned

  @REQ-COMPR-013 @REQ-COMPR-015 @error
  Scenario: DecompressPackageFile validates path parameter
    Given a compressed package
    When DecompressPackageFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-COMPR-013 @REQ-COMPR-016 @error
  Scenario: DecompressPackage respects context cancellation
    Given a compressed package
    And a cancelled context
    When DecompressPackage is called
    Then structured context error is returned
    And error type is context cancellation
    And decompression operation stops
