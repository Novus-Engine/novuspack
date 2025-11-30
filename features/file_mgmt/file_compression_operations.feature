@domain:file_mgmt @m2 @REQ-FILEMGMT-018 @spec(api_file_management.md#7-file-compression-operations)
Feature: File compression operations

  @REQ-FILEMGMT-035 @happy
  Scenario: CompressFile compresses file content
    Given an open writable NovusPack package with uncompressed file
    When CompressFile is called with compression type and level
    Then file content is compressed
    And StoredSize reflects compressed size
    And CompressionType is set
    And CompressionLevel is set
    And file remains readable after decompression

  @REQ-FILEMGMT-036 @happy
  Scenario: DecompressFile decompresses file content
    Given an open NovusPack package with compressed file
    When DecompressFile is called
    Then file content is decompressed
    And original size is restored
    And file content matches original

  @happy
  Scenario: GetFileCompressionInfo returns compression information
    Given an open NovusPack package with compressed file
    When GetFileCompressionInfo is called
    Then compression type is returned
    And compression level is returned
    And compression ratio is returned
    And compression status is returned

  @error
  Scenario: CompressFile fails with invalid compression type
    Given an open writable NovusPack package
    When CompressFile is called with invalid type
    Then structured validation error is returned

  @error
  Scenario: Compression operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When compression operation is called
    Then structured context error is returned
