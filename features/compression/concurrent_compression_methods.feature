@domain:compression @m2 @REQ-COMPR-145 @spec(api_package_compression.md#83-concurrent-compression-methods)
Feature: Concurrent compression methods

  @REQ-COMPR-145 @happy
  Scenario: CompressPackageConcurrent provides parallel compression using worker pool
    Given a NovusPack package requiring compression
    When CompressPackageConcurrent is called with compression type and config
    Then compression uses worker pool for parallel processing
    And package content is compressed concurrently
    And parallel compression improves performance

  @REQ-COMPR-145 @happy
  Scenario: DecompressPackageConcurrent provides parallel decompression
    Given a compressed NovusPack package
    When DecompressPackageConcurrent is called with config
    Then decompression uses worker pool for parallel processing
    And package content is decompressed concurrently
    And parallel decompression improves performance

  @REQ-COMPR-145 @happy
  Scenario: Concurrent compression methods use worker pool for resource management
    Given concurrent compression operations
    When compression methods are called
    Then worker pool manages concurrent workers
    And resource usage is optimized
    And worker pool handles resource allocation
