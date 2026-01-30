@domain:compression @m2 @REQ-COMPR-123 @spec(api_package_compression.md#521-purpose)
Feature: DecompressPackageStream purpose

  @REQ-COMPR-123 @happy
  Scenario: DecompressPackageStream exists to support streaming decompression
    Given a compressed package data stream
    When decompression is performed using DecompressPackageStream
    Then decompression is performed with streaming I/O
    And streaming decompression supports large packages without full buffering
    And output is a valid uncompressed package stream or file content as specified
    And decompression can be combined with other streaming stages
    And the purpose aligns with documented streaming decompression guidance

