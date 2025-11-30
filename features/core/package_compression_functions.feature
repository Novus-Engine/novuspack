@domain:core @m2 @REQ-CORE-038 @spec(api_core.md#62-package-compression-functions)
Feature: Package Compression Functions

  @REQ-CORE-038 @happy
  Scenario: Package compression functions provide CompressPackage and DecompressPackage
    Given an open NovusPack package
    And a valid context
    When package compression functions are used
    Then CompressPackage function is available
    And DecompressPackage function is available
    And functions compress or decompress package in memory

  @REQ-CORE-038 @happy
  Scenario: Package compression functions provide CompressPackageFile and DecompressPackageFile
    Given an open NovusPack package
    And a valid context
    And a file path
    When package compression functions are used
    Then CompressPackageFile function is available
    And DecompressPackageFile function is available
    And functions handle compression and file I/O together

  @REQ-CORE-038 @happy
  Scenario: Package compression functions provide compression information methods
    Given an open NovusPack package
    And a valid context
    When package compression functions are used
    Then GetPackageCompressionInfo function is available
    And IsPackageCompressed function is available
    And GetPackageCompressionType function is available

  @REQ-CORE-038 @happy
  Scenario: Package compression functions link to Package Compression API
    Given an open NovusPack package
    And compression operations are needed
    When package compression functions are examined
    Then functions reference Package Compression API documentation
    And detailed method signatures are provided
    And implementation details are documented
