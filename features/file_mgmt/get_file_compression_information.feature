@domain:file_mgmt @m2 @REQ-FILEMGMT-033 @spec(api_file_management.md#7-file-compression-operations)
Feature: Get file compression information

  @happy
  Scenario: GetFileCompressionInfo returns compression information
    Given an open package with compressed file
    When GetFileCompressionInfo is called with file path
    Then FileCompressionInfo is returned
    And compression type is included
    And compressed size is included
    And original size is included
    And compression ratio is included

  @happy
  Scenario: GetFileCompressionInfo returns info for uncompressed file
    Given an open package with uncompressed file
    When GetFileCompressionInfo is called with file path
    Then FileCompressionInfo is returned
    And compression type indicates no compression
    And compressed size equals original size

  @happy
  Scenario: CompressFile compresses existing file
    Given an open writable package with uncompressed file
    When CompressFile is called with file path and compression type
    Then file is compressed
    And compression type is set in file entry
    And file can be decompressed

  @happy
  Scenario: DecompressFile decompresses existing file
    Given an open writable package with compressed file
    When DecompressFile is called with file path
    Then file is decompressed
    And compression type is cleared in file entry
    And original content is restored

  @error
  Scenario: CompressFile fails for read-only package
    Given a read-only package with file
    When CompressFile is called
    Then structured validation error is returned

  @error
  Scenario: CompressFile fails with unsupported compression type
    Given an open writable package with file
    When CompressFile is called with unsupported compression type
    Then structured validation error is returned

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: File compression methods validate path parameter
    Given an open writable package
    When CompressFile or DecompressFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: File compression methods respect context cancellation
    Given an open writable package with file
    And a cancelled context
    When CompressFile or DecompressFile is called
    Then structured context error is returned
    And error type is context cancellation
