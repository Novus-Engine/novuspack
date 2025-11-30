@domain:writing @m2 @REQ-WRITE-030 @REQ-WRITE-031 @spec(api_writing.md#5-compressed-package-write-operations)
Feature: Compressed Package Write Operations

  @REQ-WRITE-030 @happy
  Scenario: Compressed package write operations support compressed packages
    Given an open NovusPack package
    And the package is compressed (compression type in header flags Bits 15-8)
    When Write is called with the target path
    Then compressed package write is supported
    And SafeWrite is used for compressed packages
    And compression state is preserved in written package
    And header flags maintain compression type

  @REQ-WRITE-030 @happy
  Scenario: Compressed package write handles decompression and recompression
    Given an open NovusPack package
    And the package is compressed with Zstd (compression type = 1)
    When Write is called with the target path
    Then package is decompressed before writing operations
    And package modifications are applied
    And package is recompressed with original compression type
    And written package maintains compression state

  @REQ-WRITE-030 @happy
  Scenario: Compressed package write preserves header and signatures uncompressed
    Given an open NovusPack package
    And the package is compressed
    When Write is called with the target path
    Then header remains uncompressed for direct access
    And comment remains uncompressed for easy reading
    And signatures remain uncompressed for validation
    And only file entries, data, and index are compressed

  @REQ-WRITE-030 @happy
  Scenario: Compressed package write uses streaming for large packages
    Given an open NovusPack package
    And the package is compressed
    And the package is large (>100MB)
    When Write is called with the target path
    Then streaming is used for large compressed packages
    And memory usage is controlled
    And decompression/recompression uses streaming

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection checks header flags Bits 15-8
    Given an open NovusPack package
    When compressed package detection is performed
    Then package compression field is checked (Bits 15-8 in header flags)
    And IsPackageCompressed flag is determined
    And compression type is identified (0=none, 1=Zstd, 2=LZ4, 3=LZMA)

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies Zstd compression
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 1 (Zstd)
    When compressed package detection is performed
    Then IsPackageCompressed returns true
    And compression type is identified as Zstd
    And compression state is correctly determined

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies LZ4 compression
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 2 (LZ4)
    When compressed package detection is performed
    Then IsPackageCompressed returns true
    And compression type is identified as LZ4
    And compression state is correctly determined

  @REQ-WRITE-031 @happy
  Scenario: Compressed package detection identifies uncompressed package
    Given an open NovusPack package
    And header flags Bits 15-8 indicate compression type 0 (none)
    When compressed package detection is performed
    Then IsPackageCompressed returns false
    And compression type is identified as none
    And uncompressed state is correctly determined

  @REQ-WRITE-030 @error
  Scenario: Compressed package write returns error when compression type is unsupported
    Given an open NovusPack package
    And header flags indicate unsupported compression type (>3)
    When Write is called with the target path
    Then UnsupportedCompression error is returned
    And error indicates unsupported compression type
    And error follows structured error format

  @REQ-WRITE-030 @error
  Scenario: Compressed package write returns error when compressed data is corrupted
    Given an open NovusPack package
    And the package is compressed
    And compressed data is corrupted
    When Write is called with the target path
    Then CorruptedCompressedData error is returned
    And error indicates data corruption
    And error follows structured error format
