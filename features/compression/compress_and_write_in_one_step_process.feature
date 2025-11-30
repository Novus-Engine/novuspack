@domain:compression @m2 @REQ-COMPR-038 @spec(api_package_compression.md#11311-process)
Feature: Compress and Write in One Step Process

  @REQ-COMPR-038 @happy
  Scenario: CompressPackageFile compresses and writes package in one step
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type
    And an overwrite flag
    When CompressPackageFile is called
    Then package content is compressed in memory
    And compressed package is written to specified path
    And operation completes in single step

  @REQ-COMPR-038 @happy
  Scenario: CompressPackageFile process uses target file path parameter
    Given an open NovusPack package
    And a valid context
    And a target file path
    When CompressPackageFile is called with path
    Then path specifies where compressed package is written
    And path is validated before writing
    And file is created or overwritten based on flag

  @REQ-COMPR-038 @happy
  Scenario: CompressPackageFile process uses compression type parameter
    Given an open NovusPack package
    And a valid context
    And a target file path
    And a compression type 1-3
    When CompressPackageFile is called with compression type
    Then compression type is applied during compression
    And compressed package uses specified algorithm
    And compression is performed before writing

  @REQ-COMPR-038 @happy
  Scenario: CompressPackageFile process uses overwrite flag parameter
    Given an open NovusPack package
    And a valid context
    And a target file path that exists
    And overwrite flag is set
    When CompressPackageFile is called
    Then overwrite flag determines file handling
    And file is overwritten if flag is true
    And error is returned if flag is false and file exists

  @REQ-COMPR-038 @happy
  Scenario: CompressPackageFile process combines compression and writing
    Given an open NovusPack package
    And a valid context
    And compression and writing are needed
    When CompressPackageFile is called
    Then compression and writing occur in single operation
    And no intermediate steps are required
    And process is simplified compared to separate operations

  @REQ-COMPR-038 @error
  Scenario: CompressPackageFile process returns error when package is signed
    Given an open NovusPack package
    And package has signatures
    And a valid context
    And a target file path
    When CompressPackageFile is called
    Then security error is returned
    And error indicates package cannot be compressed when signed
    And error follows structured error format

  @REQ-COMPR-038 @error
  Scenario: CompressPackageFile process returns error for invalid compression type
    Given an open NovusPack package
    And an invalid compression type
    And a valid context
    And a target file path
    When CompressPackageFile is called
    Then validation error is returned
    And error indicates invalid compression type
    And error follows structured error format
