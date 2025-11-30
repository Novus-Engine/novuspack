@domain:writing @m2 @REQ-WRITE-050 @spec(api_writing.md#571-compression-errors)
Feature: Writing Compression Errors

  @REQ-WRITE-050 @error
  Scenario: Compression errors define compression-specific errors
    Given a NovusPack package
    When compression errors are examined
    Then CompressionFailure is returned when compression operation fails
    And DecompressionFailure is returned when decompression operation fails
    And UnsupportedCompression is returned for unsupported compression types
    And CorruptedCompressedData is returned when compressed data is corrupted
    And CompressSignedPackageError is returned when attempting to compress signed package

  @REQ-WRITE-050 @error
  Scenario: CompressionFailure indicates compression operation failure
    Given a NovusPack package
    And compression operation failure during write
    When compression is attempted
    Then CompressionFailure error is returned
    And error indicates compression operation failed
    And error follows structured error format

  @REQ-WRITE-050 @error
  Scenario: DecompressionFailure indicates decompression operation failure
    Given a NovusPack package
    And decompression operation failure during write
    When decompression is attempted
    Then DecompressionFailure error is returned
    And error indicates decompression operation failed
    And error follows structured error format

  @REQ-WRITE-050 @error
  Scenario: UnsupportedCompression indicates unsupported compression type
    Given a NovusPack package
    And unsupported compression type
    When compression with unsupported type is attempted
    Then UnsupportedCompression error is returned
    And error indicates unsupported compression type
    And error follows structured error format

  @REQ-WRITE-050 @error
  Scenario: CorruptedCompressedData indicates corrupted compressed data
    Given a NovusPack package
    And corrupted compressed data
    When compressed data is processed
    Then CorruptedCompressedData error is returned
    And error indicates compressed data is corrupted
    And error follows structured error format

  @REQ-WRITE-050 @error
  Scenario: CompressSignedPackageError indicates signed package compression attempt
    Given a NovusPack package
    And a signed package
    When compression is attempted on signed package
    Then CompressSignedPackageError is returned
    And error indicates signed package cannot be compressed
    And error follows structured error format
