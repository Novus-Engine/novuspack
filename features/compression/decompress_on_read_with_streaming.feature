@domain:compression @m2 @REQ-COMPR-002 @spec(api_package_compression.md#52-decompresspackagestream)
Feature: Decompress on read with streaming

  @happy
  Scenario: Transparent decompression for consumers
    Given a package with compressed content
    When I read the file contents via stream
    Then I should receive decompressed bytes matching the original

  @happy
  Scenario: Streaming decompression handles large packages
    Given a large compressed package
    When DecompressPackageStream is called
    Then decompression is performed via stream
    And memory usage is controlled
    And decompressed data is streamed

  @happy
  Scenario: Streaming decompression preserves package structure
    Given a compressed package
    When streaming decompression is performed
    Then package structure is preserved
    And file entries are decompressed correctly
    And file data is decompressed correctly
    And file index is decompressed correctly

  @error
  Scenario: Streaming decompression handles errors gracefully
    Given a corrupted compressed package
    When streaming decompression is attempted
    Then structured corruption error is returned
    And stream is closed properly
