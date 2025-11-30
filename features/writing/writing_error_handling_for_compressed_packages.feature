@domain:writing @m2 @REQ-WRITE-049 @spec(api_writing.md#57-error-handling)
Feature: Writing Error Handling for Compressed Packages

  @REQ-WRITE-049 @error
  Scenario: Write error handling: CompressionFailure returned when compression fails
    Given an open NovusPack package
    And compression operation fails during write
    When Write is called with compressionType parameter
    Then CompressionFailure error is returned
    And error indicates compression failure
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: DecompressionFailure returned when decompression fails
    Given an open NovusPack package
    And the package is compressed
    And decompression operation fails during write
    When Write is called with the target path
    Then DecompressionFailure error is returned
    And error indicates decompression failure
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: UnsupportedCompression returned for unsupported types
    Given an open NovusPack package
    And compressionType parameter is unsupported (>3)
    When Write is called with unsupported compressionType
    Then UnsupportedCompression error is returned
    And error indicates unsupported compression type
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: CorruptedCompressedData returned when data is corrupted
    Given an open NovusPack package
    And the package is compressed
    And compressed data is corrupted
    When Write is called with the target path
    Then CorruptedCompressedData error is returned
    And error indicates data corruption
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: CompressSignedPackageError returned for signed packages
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    And compressionType parameter is non-zero
    When Write is called with the target path
    Then CompressSignedPackageError is returned
    And error indicates signed packages cannot be compressed
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: FastWriteOnCompressed returned for compressed packages
    Given an open NovusPack package
    And the package is compressed
    When FastWrite is called directly with the target path
    Then FastWriteOnCompressed error is returned
    And error indicates FastWrite is not supported for compressed packages
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: CompressionMismatch returned when types don't match
    Given an open NovusPack package
    And compression type doesn't match expectations
    When Write is called with the target path
    Then CompressionMismatch error is returned
    And error indicates compression type mismatch
    And error follows structured error format

  @REQ-WRITE-049 @error
  Scenario: Write error handling: MemoryInsufficient returned when memory is insufficient
    Given an open NovusPack package
    And insufficient memory is available for compression operations
    When Write is called with compressionType parameter
    Then MemoryInsufficient error is returned
    And error indicates insufficient memory
    And error follows structured error format
