@domain:compression @m2 @REQ-COMPR-032 @spec(api_package_compression.md#1111-zstandard-compression-1)
Feature: LZ4 compression

  @REQ-COMPR-032 @happy
  Scenario: LZ4 compression provides fastest compression and decompression
    Given compression operations requiring speed
    When LZ4 compression type (2) is used
    Then compression is fastest among compression types
    And decompression is fastest among compression types
    And speed is prioritized over compression ratio

  @REQ-COMPR-032 @happy
  Scenario: LZ4 compression has lower compression ratio but faster performance
    Given compression operations
    When LZ4 compression is used
    Then compression ratio is lower than other types
    And compression speed is fastest
    And decompression speed is fastest

  @REQ-COMPR-032 @happy
  Scenario: LZ4 compression is good for real-time applications
    Given compression operations for real-time use cases
    When LZ4 compression is selected
    Then compression is suitable for real-time applications
    And low latency is achieved
    And performance meets real-time requirements

  @REQ-COMPR-032 @happy
  Scenario: LZ4 compression provides low CPU usage
    Given compression operations
    When LZ4 compression is used
    Then CPU usage is low
    And compression is efficient
    And system resources are used minimally
