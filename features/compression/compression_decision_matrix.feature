@domain:compression @m2 @REQ-COMPR-035 @spec(api_package_compression.md#112-compression-decision-matrix)
Feature: Compression Decision Matrix

  @REQ-COMPR-035 @happy
  Scenario: Compression decision matrix recommends LZ4 for real-time processing
    Given a compression use case
    And use case is real-time processing
    When compression decision matrix is consulted
    Then LZ4 is recommended
    And reason is speed priority
    And LZ4 provides fastest compression and decompression

  @REQ-COMPR-035 @happy
  Scenario: Compression decision matrix recommends Zstandard for archival storage
    Given a compression use case
    And use case is archival storage
    When compression decision matrix is consulted
    Then Zstandard is recommended
    And reason is balanced performance
    And Zstandard provides best compression ratio with moderate CPU usage

  @REQ-COMPR-035 @happy
  Scenario: Compression decision matrix recommends LZMA for maximum compression
    Given a compression use case
    And use case requires maximum compression
    When compression decision matrix is consulted
    Then LZMA is recommended
    And reason is size priority
    And LZMA provides highest compression ratio

  @REQ-COMPR-035 @happy
  Scenario: Compression decision matrix recommends Zstandard for network transfer
    Given a compression use case
    And use case is network transfer
    When compression decision matrix is consulted
    Then Zstandard is recommended
    And reason is good balance
    And Zstandard balances compression ratio and speed

  @REQ-COMPR-035 @happy
  Scenario: Compression decision matrix provides guidance for compression type selection
    Given a compression use case
    When compression decision matrix is consulted
    Then matrix provides use case to type mapping
    And matrix explains reasoning for recommendations
    And matrix helps users select appropriate compression type
