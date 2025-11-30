@domain:writing @m2 @REQ-WRITE-032 @REQ-WRITE-039 @spec(api_writing.md#52-write-operations-on-compressed-packages)
Feature: Write Operations on Compressed Packages

  @REQ-WRITE-032 @happy
  Scenario: Write operations on compressed packages decompress before operations
    Given an open NovusPack package
    And the package is compressed (compression type in header flags)
    When Write is called with the target path
    Then package is decompressed before write operations
    And modifications can be applied
    And decompression is handled correctly

  @REQ-WRITE-032 @happy
  Scenario: Write operations on compressed packages preserve compression state
    Given an open NovusPack package
    And the package is compressed with Zstd (compression type = 1)
    When Write is called with the target path
    Then original compression settings are preserved in header
    And compression type is maintained
    And header flags retain compression information

  @REQ-WRITE-032 @happy
  Scenario: Write operations on compressed packages recompress after modifications
    Given an open NovusPack package
    And the package is compressed
    And package modifications have been applied
    When Write is called with the target path
    Then package is recompressed with original compression type
    And written package maintains compression state
    And compression is applied correctly

  @REQ-WRITE-032 @happy
  Scenario: Write operations on compressed packages use SafeWrite
    Given an open NovusPack package
    And the package is compressed
    When Write is called with the target path
    Then SafeWrite is automatically selected
    And FastWrite is not used for compressed packages
    And compressed package handling is correct

  @REQ-WRITE-032 @happy
  Scenario: Write operations on compressed packages preserve header uncompressed
    Given an open NovusPack package
    And the package is compressed
    When Write is called with the target path
    Then header remains uncompressed for direct access
    And comment remains uncompressed for easy reading
    And signatures remain uncompressed for validation
    And only file entries, data, and index are compressed

  @REQ-WRITE-039 @happy
  Scenario: Write method compression handling accepts compressionType parameter
    Given an open NovusPack package
    When Write is called with compressionType parameter
    Then compressionType parameter is validated
    And compression type 0 indicates no compression
    And compression types 1-3 indicate specific compression (1=Zstd, 2=LZ4, 3=LZMA)

  @REQ-WRITE-039 @happy
  Scenario: Write method compression handling uses internal compression methods
    Given an open NovusPack package
    And compressionType parameter is non-zero
    When Write is called with the target path
    Then internal compression methods are used before writing
    And CompressPackage or CompressPackageFile methods are invoked
    And compression is applied correctly

  @REQ-WRITE-039 @happy
  Scenario: Write method compression handling selects SafeWrite for compressed packages
    Given an open NovusPack package
    And compressionType parameter is non-zero
    When Write is called with the target path
    Then SafeWrite is selected for compressed packages
    And FastWrite is not used when compression is specified
    And compression workflow is handled correctly

  @REQ-WRITE-039 @happy
  Scenario: Write method compression handling selects FastWrite for uncompressed packages
    Given an open NovusPack package
    And compressionType parameter is 0
    And the package is uncompressed
    When Write is called with the target path
    Then FastWrite can be selected for uncompressed packages
    And compression is not applied
    And write strategy selection works correctly

  @REQ-WRITE-039 @happy
  Scenario: Write method compression handling compresses file entries, data, and index
    Given an open NovusPack package
    And compressionType parameter is non-zero
    When Write is called with the target path
    Then file entries are compressed
    And data section is compressed
    And index is compressed
    And header, comment, and signatures remain uncompressed

  @REQ-WRITE-032 @error
  Scenario: Write operations on compressed packages return error when decompression fails
    Given an open NovusPack package
    And the package is compressed
    And compressed data is corrupted
    When Write is called with the target path
    Then DecompressionFailure error is returned
    And error indicates decompression failure
    And error follows structured error format

  @REQ-WRITE-039 @error
  Scenario: Write method compression handling returns error for unsupported compression types
    Given an open NovusPack package
    And compressionType parameter is unsupported (>3)
    When Write is called with the target path
    Then UnsupportedCompression error is returned
    And error indicates unsupported compression type
    And error follows structured error format
