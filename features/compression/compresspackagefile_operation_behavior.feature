@domain:compression @m2 @REQ-COMPR-129 @spec(api_package_compression.md#613-compresspackagefile-behavior)
Feature: CompressPackageFile Operation Behavior

  @REQ-COMPR-129 @happy
  Scenario: CompressPackageFile compresses package content in memory then writes to file
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When CompressPackageFile is called
    Then package content is compressed in memory first
    And compressed package is written to specified path
    And file operations complete successfully

  @REQ-COMPR-129 @happy
  Scenario: CompressPackageFile creates new file by default
    Given an open NovusPack package
    And a valid context
    And a target file path that does not exist
    And an overwrite flag set to false
    When CompressPackageFile is called
    Then new file is created at specified path
    And compressed package is written to new file
    And existing files are not overwritten

  @REQ-COMPR-129 @happy
  Scenario: CompressPackageFile overwrites existing file when overwrite is true
    Given an open NovusPack package
    And a valid context
    And a target file path that exists
    And an overwrite flag set to true
    When CompressPackageFile is called
    Then existing file is overwritten
    And compressed package replaces existing content
    And file update completes successfully

  @REQ-COMPR-129 @happy
  Scenario: CompressPackageFile compresses file entries, data, and index only
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    When CompressPackageFile is called
    Then file entries are compressed
    And file data is compressed
    And package index is compressed
    And header, comment, and signatures remain uncompressed

  @REQ-COMPR-129 @error
  Scenario: CompressPackageFile returns error if package is signed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a target file path
    When CompressPackageFile is called
    Then security error is returned
    And error indicates package is signed
    And error follows structured error format
