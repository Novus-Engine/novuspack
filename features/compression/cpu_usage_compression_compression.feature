@domain:compression @m2 @REQ-COMPR-086 @spec(api_package_compression.md#1351-cpu-usage-compression)
Feature: CPU Usage Compression

  @REQ-COMPR-086 @happy
  Scenario: LZ4 compression has lowest CPU usage
    Given compression operations with LZ4 compression type
    When compression is performed
    Then CPU usage is lowest among compression types
    And compression is fastest
    And CPU resources are used efficiently

  @REQ-COMPR-086 @happy
  Scenario: Zstandard compression has moderate CPU usage
    Given compression operations with Zstandard compression type
    When compression is performed
    Then CPU usage is moderate
    And compression provides balanced performance
    And CPU usage is reasonable for compression ratio

  @REQ-COMPR-086 @happy
  Scenario: LZMA compression has highest CPU usage
    Given compression operations with LZMA compression type
    When compression is performed
    Then CPU usage is highest among compression types
    And compression provides maximum compression ratio
    And CPU resources are used intensively

  @REQ-COMPR-086 @happy
  Scenario: CPU usage scales with compression level and package size
    Given compression operations
    When compression is performed
    Then CPU usage scales with compression level
    And CPU usage scales with package size
    And higher levels require more CPU resources
