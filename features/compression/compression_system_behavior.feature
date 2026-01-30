@skip @domain:compression @m2 @REQ-COMPR-109 @spec(api_package_compression.md#413-compresspackage-behavior)
Feature: Compression System Behavior

# This feature captures core package compression behaviors from the compression API specification.
# More detailed runnable scenarios for compression operations live in dedicated compression feature files.

  @REQ-COMPR-109 @behavior
  Scenario: CompressPackage compresses the package index as a single block
    Given an uncompressed package in memory
    When CompressPackage is called
    Then file entry metadata remains individually addressable
    And the file index is stored as a single compressed block

  @REQ-COMPR-109 @constraint
  Scenario: Compression writes a metadata index at a fixed offset when package compression is enabled
    Given an uncompressed package in memory
    When CompressPackage is called
    Then a metadata index exists for fast access to compressed blocks
    And the metadata index is written at the fixed offset defined by the compression constraints
