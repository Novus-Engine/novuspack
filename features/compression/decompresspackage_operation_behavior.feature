@domain:compression @m2 @REQ-COMPR-114 @spec(api_package_compression.md#423-decompresspackage-behavior)
Feature: DecompressPackage Operation Behavior

  @REQ-COMPR-114 @happy
  Scenario: DecompressPackage decompresses all compressed content
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then all compressed content is decompressed
    And file entries are decompressed
    And file data is decompressed
    And package index is decompressed

  @REQ-COMPR-114 @happy
  Scenario: DecompressPackage updates package compression state in memory
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then package compression state is updated in memory
    And IsCompressed returns false after decompression
    And package state reflects decompressed status

  @REQ-COMPR-114 @happy
  Scenario: DecompressPackage clears package header compression flags
    Given an open NovusPack package
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then package header compression flags are cleared
    And header reflects uncompressed state
    And compression flags are reset

  @REQ-COMPR-114 @happy
  Scenario: DecompressPackage preserves all other package data
    Given an open NovusPack package
    And package is compressed
    And package has files and metadata
    And a valid context
    When DecompressPackage is called
    Then all other package data is preserved
    And files remain intact
    And metadata is preserved
    And only compression state changes

  @REQ-COMPR-114 @happy
  Scenario: DecompressPackage operates entirely in memory
    Given an open NovusPack package in memory
    And package is compressed
    And a valid context
    When DecompressPackage is called
    Then decompression occurs in memory
    And no file I/O operations are performed
    And package remains in memory after decompression
