@domain:compression @m2 @REQ-COMPR-087 @spec(api_package_compression.md#1352-cpu-usage-decompression)
Feature: CPU Usage Decompression

  @REQ-COMPR-087 @happy
  Scenario: Decompression is generally faster than compression
    Given compression operations
    When decompression is performed
    Then decompression is generally faster than compression
    And decompression requires less CPU time
    And decompression is more efficient

  @REQ-COMPR-087 @happy
  Scenario: LZ4 decompression is fastest
    Given decompression operations with LZ4 compressed data
    When decompression is performed
    Then LZ4 decompression is fastest
    And CPU usage is lowest for decompression
    And decompression speed is optimal

  @REQ-COMPR-087 @happy
  Scenario: Zstandard decompression has moderate speed
    Given decompression operations with Zstandard compressed data
    When decompression is performed
    Then Zstandard decompression has moderate speed
    And CPU usage is reasonable
    And decompression provides balanced performance

  @REQ-COMPR-087 @happy
  Scenario: LZMA decompression is slowest
    Given decompression operations with LZMA compressed data
    When decompression is performed
    Then LZMA decompression is slowest
    And CPU usage is higher for decompression
    And decompression takes more time but provides maximum compression
