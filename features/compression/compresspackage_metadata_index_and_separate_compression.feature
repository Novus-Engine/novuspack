@domain:compression @m2 @REQ-COMPR-154 @REQ-COMPR-155 @REQ-COMPR-156 @REQ-COMPR-157 @REQ-COMPR-158 @REQ-COMPR-159 @spec(api_package_compression.md#111-compressed-content)
Feature: CompressPackage Metadata Index and Separate Compression

  @REQ-COMPR-154 @happy
  Scenario: File entry metadata compressed individually using LZ4
    Given an open NovusPack package
    And a valid context
    And file entries in the package
    When CompressPackage is called with compression type
    Then file entry metadata is compressed individually
    And each file entry uses LZ4 compression
    And metadata compression is separate from data compression

  @REQ-COMPR-155 @happy
  Scenario: File data compressed individually using specified compression type
    Given an open NovusPack package
    And a valid context
    And file data in the package
    When CompressPackage is called with CompressionZstd
    Then file data is compressed individually using Zstd
    When CompressPackage is called with CompressionLZ4
    Then file data is compressed individually using LZ4
    When CompressPackage is called with CompressionLZMA
    Then file data is compressed individually using LZMA

  @REQ-COMPR-156 @happy
  Scenario: File index compressed as single block using LZ4
    Given an open NovusPack package
    And a valid context
    And a file index in the package
    When CompressPackage is called
    Then file index is compressed as single block
    And file index uses LZ4 compression
    And file index compression is separate from metadata compression

  @REQ-COMPR-157 @happy
  Scenario: Special metadata files compressed with LZ4 for fast access
    Given an open NovusPack package
    And a valid context
    And special metadata files (types 65000-65535) in the package
    When CompressPackage is called
    Then special metadata file entry metadata uses LZ4
    And special metadata file data (YAML content) uses LZ4
    And LZ4 provides fast access to special metadata files

  @REQ-COMPR-158 @happy
  Scenario: Metadata index remains uncompressed for direct access
    Given an open NovusPack package
    And a valid context
    When CompressPackage is called
    Then metadata index is created
    And metadata index remains uncompressed
    And metadata index enables direct access to compressed blocks
    And metadata index can be read without decompression

  @REQ-COMPR-159 @happy
  Scenario: Metadata index located at fixed offset 112 bytes after header
    Given an open NovusPack package
    And a valid context
    When CompressPackage is called
    Then metadata index is created
    And metadata index is written at offset 112 bytes
    And offset 112 bytes is immediately after header
    And metadata index location is fixed when compression enabled

  @REQ-COMPR-160 @happy
  Scenario: CompressPackage compresses file entry metadata individually using LZ4
    Given an open NovusPack package
    And a valid context
    And file entries in the package
    When CompressPackage is called
    Then each file entry metadata is compressed individually
    And LZ4 compression is used for metadata
    And metadata compression occurs before data compression

  @REQ-COMPR-161 @happy
  Scenario: CompressPackage compresses file data individually using specified type
    Given an open NovusPack package
    And a valid context
    And file data in the package
    When CompressPackage is called with CompressionZstd
    Then each file data is compressed individually
    And Zstd compression is used for file data
    And file data compression is separate from metadata compression

  @REQ-COMPR-162 @happy
  Scenario: CompressPackage compresses file index with LZ4 as single block
    Given an open NovusPack package
    And a valid context
    And a file index in the package
    When CompressPackage is called
    Then file index is compressed as single block
    And LZ4 compression is used for file index
    And file index compression occurs after file data compression

  @REQ-COMPR-163 @happy
  Scenario: CompressPackage creates metadata index for fast access
    Given an open NovusPack package
    And a valid context
    When CompressPackage is called
    Then metadata index is created
    And metadata index provides fast access to compressed blocks
    And metadata index enables selective decompression
    And metadata index allows metadata access without full decompression

  @REQ-COMPR-164 @happy
  Scenario: CompressPackage writes metadata index at fixed offset 112 bytes
    Given an open NovusPack package
    And a valid context
    When CompressPackage is called
    Then metadata index is written
    And metadata index is written at offset 112 bytes
    And offset 112 bytes is immediately after header
    And metadata index location is fixed and predictable

  @REQ-COMPR-165 @happy
  Scenario: DecompressPackage removes metadata index when decompressed
    Given an open NovusPack package
    And package is compressed
    And metadata index exists
    And a valid context
    When DecompressPackage is called
    Then metadata index is removed
    And metadata index is no longer present
    And package is fully decompressed
    And all compressed content is restored

  @REQ-COMPR-163 @error
  Scenario: CompressPackage returns error if metadata index creation fails
    Given an open NovusPack package
    And a valid context
    And conditions that prevent metadata index creation
    When CompressPackage is called
    Then *PackageError is returned
    And error indicates metadata index creation failure

  @REQ-COMPR-159 @error
  Scenario: CompressPackage returns error if metadata index cannot be written at fixed offset
    Given an open NovusPack package
    And a valid context
    And conditions that prevent writing at offset 112 bytes
    When CompressPackage is called
    Then *PackageError is returned
    And error indicates write failure
