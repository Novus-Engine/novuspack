@domain:compression @m2 @REQ-COMPR-006 @spec(api_package_compression.md#4-in-memory-compression-methods)
Feature: Decompress package in memory

  @happy
  Scenario: DecompressPackage decompresses package in memory
    Given a compressed package
    When DecompressPackage is called
    Then package is decompressed in memory
    And package content is accessible
    And compression type is cleared

  @error
  Scenario: DecompressPackage fails for uncompressed package
    Given an uncompressed package
    When DecompressPackage is called
    Then structured validation error is returned
