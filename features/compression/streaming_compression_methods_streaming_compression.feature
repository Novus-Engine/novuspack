@domain:compression @m2 @REQ-COMPR-116 @spec(api_package_compression.md#5-streaming-compression-methods)
Feature: Streaming Compression Methods

  @REQ-COMPR-116 @happy
  Scenario: Streaming compression methods handle large packages using streaming
    Given a compression operation
    And package is large
    When streaming compression methods are used
    Then streaming is used to avoid memory limitations
    And large packages are handled efficiently
    And temporary files are used when needed

  @REQ-COMPR-116 @happy
  Scenario: CompressPackageStream compresses large packages using streaming
    Given an open NovusPack package
    And package is large
    And a valid context
    And a compression type
    And a StreamConfig
    When CompressPackageStream is called
    Then streaming compression is used
    And chunked processing handles large files
    And memory usage is controlled

  @REQ-COMPR-116 @happy
  Scenario: DecompressPackageStream decompresses large packages using streaming
    Given an open NovusPack package
    And package is compressed and large
    And a valid context
    And a StreamConfig
    When DecompressPackageStream is called
    Then streaming decompression is used
    And chunked processing handles large files
    And memory usage is controlled

  @REQ-COMPR-116 @happy
  Scenario: Streaming compression methods use temporary files for memory management
    Given a compression operation
    And package exceeds available memory
    When streaming compression methods are used
    Then temporary files are created when needed
    And memory management uses disk storage
    And operations continue despite memory constraints
