@domain:compression @m2 @REQ-COMPR-112 @spec(api_package_compression.md#421-decompresspackage-purpose)
Feature: DecompressPackage Purpose and Usage

  @REQ-COMPR-112 @happy
  Scenario: DecompressPackage purpose is to decompress package content
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then package content is decompressed
    And compressed file entries are restored
    And compressed file data is restored
    And compressed index is restored

  @REQ-COMPR-112 @happy
  Scenario: DecompressPackage decompresses package content in memory
    Given an open NovusPack package in memory
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then decompression occurs in memory
    And no file I/O operations are performed
    And package remains in memory after decompression

  @REQ-COMPR-112 @happy
  Scenario: DecompressPackage restores package to uncompressed state
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then package is restored to uncompressed state
    And compression flags are cleared
    And package can be accessed without decompression overhead

  @REQ-COMPR-112 @happy
  Scenario: DecompressPackage enables access to uncompressed package content
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then uncompressed package content becomes accessible
    And files can be read directly
    And package operations proceed without compression overhead
