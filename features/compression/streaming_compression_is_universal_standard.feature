@domain:compression @m2 @REQ-COMPR-066 @REQ-COMPR-116 @spec(api_package_compression.md#1311-streaming-compression-universal-standard)
Feature: Streaming compression is universal standard

  @REQ-COMPR-066 @happy
  Scenario: ZSTD streaming uses industry-standard streaming functions
    Given compression operations using ZSTD streaming
    When streaming compression is performed
    Then ZSTD_compressStream2 and ZSTD_decompressStream2 are used for large files
    And streaming follows universal standard
    And industry-standard practices are followed

  @REQ-COMPR-066 @happy
  Scenario: Streaming compression provides memory efficiency
    Given compression operations with streaming
    When streaming compression is performed
    Then memory usage is constant regardless of file size
    And memory efficiency is maintained
    And large files do not require excessive memory

  @REQ-COMPR-066 @happy
  Scenario: Streaming compression enables real-time processing
    Given compression operations for large files
    When streaming compression is used
    Then real-time processing is enabled
    And files larger than available RAM can be compressed
    And processing can start before entire file is loaded

  @REQ-COMPR-066 @happy
  Scenario: Streaming compression provides progress reporting
    Given compression operations with streaming
    When streaming compression is performed
    Then industry-standard progress callbacks are provided
    And real-time progress updates are available
    And user feedback is enhanced

  @REQ-COMPR-116 @happy
  Scenario: Streaming compression methods handle large packages
    Given compression operations for large packages
    When streaming compression methods are used
    Then large packages are handled efficiently
    And memory limitations are avoided
    And packages of any size are supported

  @REQ-COMPR-116 @happy
  Scenario: CompressPackageStream and DecompressPackageStream use streaming
    Given compression operations for large packages
    When CompressPackageStream or DecompressPackageStream is called
    Then streaming is used for large package content
    And temporary files and chunked processing are employed
    And files exceeding available RAM are handled
