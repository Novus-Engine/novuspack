@domain:compression @m2 @REQ-COMPR-030 @spec(api_package_compression.md#11-compression-strategy-selection)
Feature: Compression Strategy Selection

  @REQ-COMPR-030 @happy
  Scenario: Zstandard compression provides best compression ratio with moderate CPU usage
    Given a compression operation
    And Zstandard compression type is selected
    When compression is performed
    Then best compression ratio is achieved
    And CPU usage is moderate
    And compression is good for archival storage

  @REQ-COMPR-030 @happy
  Scenario: LZ4 compression provides fastest compression and decompression
    Given a compression operation
    And LZ4 compression type is selected
    When compression is performed
    Then fastest compression and decompression is achieved
    And compression ratio is lower
    And compression is good for real-time applications

  @REQ-COMPR-030 @happy
  Scenario: LZMA compression provides highest compression ratio with highest CPU usage
    Given a compression operation
    And LZMA compression type is selected
    When compression is performed
    Then highest compression ratio is achieved
    And CPU usage is highest
    And compression is best for long-term storage

  @REQ-COMPR-030 @happy
  Scenario: Compression type selection guides algorithm choice based on use case
    Given a compression use case
    When compression type selection guidance is consulted
    Then appropriate algorithm is recommended
    And recommendation considers compression ratio needs
    And recommendation considers CPU usage constraints
    And recommendation considers speed requirements

  @REQ-COMPR-030 @happy
  Scenario: Compression strategy selection considers performance trade-offs
    Given a compression operation
    When compression strategy is selected
    Then trade-offs between compression ratio and speed are considered
    And trade-offs between compression ratio and CPU usage are considered
    And appropriate strategy is chosen for use case
