@domain:compression @m2 @REQ-COMPR-134 @spec(api_package_compression.md#623-decompresspackagefile-behavior)
Feature: DecompressPackageFile Operation Behavior

  @REQ-COMPR-134 @happy
  Scenario: DecompressPackageFile decompresses package content in memory then writes to file
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    And an overwrite flag
    When DecompressPackageFile is called
    Then package content is decompressed in memory first
    And uncompressed package is written to specified path
    And file operations complete successfully

  @REQ-COMPR-134 @happy
  Scenario: DecompressPackageFile creates new file by default
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path that does not exist
    And overwrite flag is set to false
    When DecompressPackageFile is called
    Then new file is created at specified path
    And uncompressed package is written to new file
    And existing files are not overwritten

  @REQ-COMPR-134 @happy
  Scenario: DecompressPackageFile overwrites existing file when overwrite is true
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path that exists
    And overwrite flag is set to true
    When DecompressPackageFile is called
    Then existing file is overwritten
    And uncompressed package replaces existing content
    And file update completes successfully

  @REQ-COMPR-134 @happy
  Scenario: DecompressPackageFile decompresses all compressed content
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    When DecompressPackageFile is called
    Then all compressed content is decompressed
    And decompressed content is written to file
    And package state reflects decompression
