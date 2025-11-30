@domain:compression @m2 @REQ-COMPR-107 @spec(api_package_compression.md#411-compresspackage-purpose)
Feature: CompressPackage Purpose and Usage

  @REQ-COMPR-107 @happy
  Scenario: CompressPackage purpose is to compress package in memory
    Given an open NovusPack package
    And a valid context
    And a compression type
    When CompressPackage is called
    Then package content is compressed in memory
    And compression occurs without writing to disk
    And package remains in memory after compression

  @REQ-COMPR-107 @happy
  Scenario: CompressPackage handles compression and decompression of in-memory packages
    Given an open NovusPack package
    And a valid context
    When CompressPackage is used
    Then method handles compression of in-memory packages
    And method supports decompression of in-memory packages
    And operations occur entirely in memory

  @REQ-COMPR-107 @happy
  Scenario: CompressPackage operates on packages without file I/O
    Given an open NovusPack package in memory
    And a valid context
    And a compression type
    When CompressPackage is called
    Then no file I/O operations are performed
    And compression occurs in memory only
    And package state is updated in memory
