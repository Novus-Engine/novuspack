@skip @domain:compression @m2 @REQ-COMPR-108 @spec(api_package_compression.md#412-compresspackage-parameters)
Feature: Compression Parameter Specification

# This feature captures parameter-level constraints for key compression operations.
# Detailed runnable scenarios live in the dedicated compression feature files.

  @REQ-COMPR-108 @validation
  Scenario: CompressPackage parameters include context and compression type
    Given a package in memory that contains file entries and file data
    When CompressPackage is invoked
    Then it accepts a context for cancellation and deadlines
    And it accepts a compression type that must be validated against supported algorithms

  @REQ-COMPR-113 @validation
  Scenario: DecompressPackage parameters include context
    Given a package in memory that is compressed
    When DecompressPackage is invoked
    Then it accepts a context for cancellation and deadlines
    And it returns a structured error if the context is cancelled during decompression
