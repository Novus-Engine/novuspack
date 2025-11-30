@domain:compression @m2 @REQ-COMPR-104 @spec(api_package_compression.md#33-compression-streaming-interface)
Feature: Compression Streaming Interface

  @REQ-COMPR-104 @happy
  Scenario: CompressPackageStream provides streaming compression for large packages
    Given a large package requiring compression
    When CompressPackageStream is called with compression type and config
    Then package is compressed using streaming interface
    And memory usage is controlled for large packages
    And streaming handles large data efficiently

  @REQ-COMPR-104 @happy
  Scenario: DecompressPackageStream provides streaming decompression
    Given a compressed package requiring decompression
    When DecompressPackageStream is called with config
    Then package is decompressed using streaming interface
    And memory usage is controlled
    And streaming handles large compressed data

  @REQ-COMPR-104 @happy
  Scenario: Streaming interface supports configurable streaming options
    Given compression or decompression operations
    When StreamConfig is provided
    Then streaming options are configured
    And chunk size, temp directory, and memory limits are applied
    And streaming behavior matches configuration

  @REQ-COMPR-104 @error
  Scenario: Streaming interface handles compression errors
    Given streaming compression operations
    When compression errors occur
    Then structured compression errors are returned
    And error context includes streaming information
    And errors provide details about streaming failure
