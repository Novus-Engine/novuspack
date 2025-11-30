@domain:writing @m2 @REQ-WRITE-010 @spec(api_writing.md#53-compression-operations)
Feature: Compression Operations

  @REQ-WRITE-010 @happy
  Scenario: Compression operations validate compression type parameters
    Given an open NovusPack package
    When compression operation is performed
    Then compression type parameters are validated
    And compression types 0-3 are supported (0=none, 1=Zstd, 2=LZ4, 3=LZMA)
    And unsupported compression types are rejected

  @REQ-WRITE-010 @happy
  Scenario: In-memory compression methods compress package content in memory
    Given an open NovusPack package
    When CompressPackage is called
    Then package content is compressed in memory
    And package state is updated
    And compression is applied to file entries, data, and index

  @REQ-WRITE-010 @happy
  Scenario: In-memory decompression methods decompress package in memory
    Given an open NovusPack package
    And the package is compressed
    When DecompressPackage is called
    Then package is decompressed in memory
    And package state is updated
    And package becomes uncompressed

  @REQ-WRITE-010 @happy
  Scenario: File-based compression methods compress and write to specified path
    Given an open NovusPack package
    When CompressPackageFile is called with target path
    Then package content is compressed
    And compressed package is written to specified path
    And in-memory package state is not affected

  @REQ-WRITE-010 @happy
  Scenario: File-based decompression methods decompress and write to specified path
    Given an open NovusPack package
    And the package is compressed
    When DecompressPackageFile is called with target path
    Then package is decompressed
    And decompressed package is written to specified path
    And in-memory package state is not affected

  @REQ-WRITE-010 @happy
  Scenario: Compression operations compress file entries, data, and index
    Given an open NovusPack package
    And compressionType parameter is non-zero
    When compression operation is performed
    Then file entries are compressed
    And data section is compressed
    And index is compressed
    And header, comment, and signatures remain uncompressed

  @REQ-WRITE-010 @error
  Scenario: Compression operations return error for unsupported compression types
    Given an open NovusPack package
    And compressionType parameter is unsupported (>3)
    When compression operation is performed
    Then UnsupportedCompression error is returned
    And error indicates unsupported compression type
    And error follows structured error format

  @REQ-WRITE-010 @error
  Scenario: Compression operations return error when compression fails
    Given an open NovusPack package
    And compression operation encounters an error
    When compression operation is performed
    Then CompressionFailure error is returned
    And error indicates compression failure
    And error follows structured error format
