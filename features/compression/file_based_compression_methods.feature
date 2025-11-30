@domain:compression @m2 @REQ-COMPR-005 @spec(api_package_compression.md#61-compresspackagefile)
Feature: File-based compression methods

  @happy
  Scenario: CompressPackageFile compresses package file
    Given a package file
    When CompressPackageFile is called with output path
    Then compressed package file is created
    And compression preserves package structure
    And output file is valid

  @happy
  Scenario: DecompressPackageFile decompresses package file
    Given a compressed package file
    When DecompressPackageFile is called with output path
    Then decompressed package file is created
    And original package structure is restored
    And output file is valid

  @error
  Scenario: File compression operations handle I/O errors
    Given a package file
    And file system I/O errors occur
    When compression operation is called
    Then structured I/O error is returned

  @REQ-COMPR-013 @REQ-COMPR-014 @error
  Scenario: Compression operations validate compression type parameter
    Given an open writable package
    When compression operation is called with invalid compression type
    Then structured validation error is returned
    And error indicates unsupported compression type

  @REQ-COMPR-013 @REQ-COMPR-016 @error
  Scenario: Compression operations respect context cancellation
    Given an open writable package
    And a cancelled context
    When compression operation is called
    Then structured context error is returned
    And error type is context cancellation
    And compression operation stops gracefully
