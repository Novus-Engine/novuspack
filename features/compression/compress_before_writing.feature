@domain:compression @m2 @REQ-COMPR-037 @spec(api_package_compression.md#1131-option-1-compress-before-writing)
Feature: Compress Before Writing

  @REQ-COMPR-037 @happy
  Scenario: CompressPackage compresses package content in memory first
    Given an open NovusPack package
    And a valid context
    And a compression type
    When CompressPackage is called
    Then package content is compressed in memory
    And compression occurs before writing
    And package is ready to write

  @REQ-COMPR-037 @happy
  Scenario: Write is called with CompressionNone after compression
    Given an open NovusPack package
    And package is compressed using CompressPackage
    And a valid context
    And an output file path
    When Write is called with CompressionNone
    Then already-compressed package is written to output file
    And no additional compression is applied
    And compressed package is written as-is

  @REQ-COMPR-037 @happy
  Scenario: Compress before writing workflow separates compression and writing steps
    Given an open NovusPack package
    And a valid context
    And a compression type
    And an output file path
    When compress before writing workflow is followed
    Then CompressPackage compresses package in memory first
    And Write writes compressed package to file
    And workflow separates compression and writing steps

  @REQ-COMPR-037 @happy
  Scenario: Compress before writing allows compression verification before writing
    Given an open NovusPack package
    And a valid context
    And a compression type
    When CompressPackage is called first
    Then compression can be verified before writing
    And compression state can be checked
    And Write only occurs after successful compression

  @REQ-COMPR-037 @error
  Scenario: Compress before writing handles compression errors before writing
    Given an open NovusPack package
    And a valid context
    And a compression type
    And compression operation fails
    When CompressPackage is called
    Then compression error is returned
    And Write is not called
    And package remains uncompressed
