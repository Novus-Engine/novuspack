@domain:compression @m2 @REQ-COMPR-031 @spec(api_package_compression.md#111-compression-type-selection)
Feature: Zstandard Compression

  @REQ-COMPR-031 @happy
  Scenario: Zstandard compression provides best compression ratio among standard options
    Given compression operations requiring good compression
    When Zstandard compression type (1) is used
    Then compression ratio is best among standard options
    And compression efficiency is high
    And compression quality is optimized

  @REQ-COMPR-031 @happy
  Scenario: Zstandard compression has moderate CPU usage
    Given compression operations
    When Zstandard compression is used
    Then CPU usage is moderate
    And compression provides balanced performance
    And CPU usage is reasonable for compression ratio

  @REQ-COMPR-031 @happy
  Scenario: Zstandard compression is good for archival storage
    Given compression operations for storage
    When Zstandard compression is selected
    Then compression is good for archival storage scenarios
    And balanced performance meets storage requirements
    And compression ratio optimizes storage space

  @REQ-COMPR-031 @happy
  Scenario: Zstandard compression provides good balance for network transfer
    Given compression operations for network transfer
    When Zstandard compression is used
    Then compression provides good balance
    And transfer efficiency is improved
    And network transfer benefits from compression
