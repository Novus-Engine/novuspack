@domain:compression @m2 @REQ-COMPR-004 @spec(api_package_compression.md#41-compresspackage)
Feature: In-memory compression methods

  @happy
  Scenario: CompressPackage compresses package in memory
    Given a package in memory
    When CompressPackage is called
    Then package content is compressed
    And compressed data is returned
    And compression preserves package structure

  @happy
  Scenario: DecompressPackage decompresses package in memory
    Given a compressed package in memory
    When DecompressPackage is called
    Then package content is decompressed
    And original package structure is restored
    And data integrity is maintained

  @error
  Scenario: Compression operations handle memory constraints
    Given a very large package
    When compression is attempted
    Then memory constraints are handled gracefully
    And error is returned if memory insufficient
