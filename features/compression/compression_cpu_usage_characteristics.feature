@domain:compression @m2 @REQ-COMPR-085 @spec(api_package_compression.md#135-cpu-usage)
Feature: Compression CPU Usage Characteristics

  @REQ-COMPR-085 @happy
  Scenario: LZ4 compression has lowest CPU usage
    Given a compression operation
    And LZ4 compression type is used
    When compression is performed
    Then LZ4 has lowest CPU usage among algorithms
    And CPU consumption is minimal
    And compression is fast

  @REQ-COMPR-085 @happy
  Scenario: Zstandard compression has moderate CPU usage
    Given a compression operation
    And Zstandard compression type is used
    When compression is performed
    Then Zstandard has moderate CPU usage
    And CPU consumption is balanced with compression ratio
    And compression provides good balance

  @REQ-COMPR-085 @happy
  Scenario: LZMA compression has highest CPU usage
    Given a compression operation
    And LZMA compression type is used
    When compression is performed
    Then LZMA has highest CPU usage among algorithms
    And CPU consumption is high
    And compression achieves maximum compression ratio

  @REQ-COMPR-085 @happy
  Scenario: Decompression is generally faster than compression
    Given a compression operation
    When decompression is performed
    Then decompression is generally faster than compression
    And CPU usage is lower during decompression
    And decompression consumes fewer resources

  @REQ-COMPR-085 @happy
  Scenario: LZ4 decompression is fastest
    Given a decompression operation
    And LZ4 compressed data
    When decompression is performed
    Then LZ4 decompression is fastest
    And CPU usage is minimal
    And decompression speed is optimal

  @REQ-COMPR-085 @happy
  Scenario: Zstandard decompression has moderate speed
    Given a decompression operation
    And Zstandard compressed data
    When decompression is performed
    Then Zstandard decompression has moderate speed
    And CPU usage is moderate
    And decompression provides balanced performance

  @REQ-COMPR-085 @happy
  Scenario: LZMA decompression is slowest
    Given a decompression operation
    And LZMA compressed data
    When decompression is performed
    Then LZMA decompression is slowest
    And CPU usage is higher
    And decompression trades speed for compression ratio
