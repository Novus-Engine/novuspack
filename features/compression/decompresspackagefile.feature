@domain:compression @m2 @REQ-COMPR-131 @spec(api_package_compression.md#62-decompresspackagefile)
Feature: DecompressPackageFile

  @REQ-COMPR-131 @happy
  Scenario: DecompressPackageFile decompresses package and writes to file
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    When DecompressPackageFile is called with path
    Then package content is decompressed in memory
    And uncompressed package is written to specified path
    And file is created successfully

  @REQ-COMPR-131 @happy
  Scenario: DecompressPackageFile creates new file by default
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path that does not exist
    And overwrite is set to false
    When DecompressPackageFile is called
    Then new file is created at specified path
    And uncompressed package is written to file
    And existing file is not overwritten

  @REQ-COMPR-131 @happy
  Scenario: DecompressPackageFile overwrites existing file when requested
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path that exists
    And overwrite is set to true
    When DecompressPackageFile is called
    Then existing file is overwritten
    And uncompressed package is written to file
    And new content replaces old content

  @REQ-COMPR-131 @happy
  Scenario: DecompressPackageFile decompresses all compressed content
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path
    When DecompressPackageFile is called
    Then all compressed content is decompressed
    And decompressed content is written to file
    And package state reflects decompression

  @REQ-COMPR-131 @error
  Scenario: DecompressPackageFile returns error when package is not compressed
    Given an open NovusPack package
    And package is not compressed
    And a valid context
    And a target file path
    When DecompressPackageFile is called
    Then validation error is returned
    And error indicates package is not compressed
    And error follows structured error format

  @REQ-COMPR-131 @error
  Scenario: DecompressPackageFile returns error when file exists and overwrite is false
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path that exists
    And overwrite is set to false
    When DecompressPackageFile is called
    Then validation error is returned
    And error indicates file already exists
    And error follows structured error format

  @REQ-COMPR-131 @error
  Scenario: DecompressPackageFile returns error for I/O failures
    Given an open NovusPack package
    And package is compressed
    And a valid context
    And a target file path with I/O errors
    When DecompressPackageFile is called
    Then I/O error is returned
    And error indicates file operation failure
    And error follows structured error format
