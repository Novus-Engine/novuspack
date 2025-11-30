@domain:compression @m2 @REQ-COMPR-033 @REQ-COMPR-034 @spec(api_package_compression.md#1112-lz4-compression-2,api_package_compression.md#1113-lzma-compression-3)
Feature: LZMA compression

  @REQ-COMPR-033 @REQ-COMPR-034 @happy
  Scenario: LZMA compression provides highest compression ratio
    Given compression operations requiring maximum compression
    When LZMA compression type (3) is used
    Then compression ratio is highest among compression types
    And file size reduction is maximized
    And compression ratio is optimized
    And compression efficiency is optimized

  @REQ-COMPR-033 @REQ-COMPR-034 @happy
  Scenario: LZMA compression has highest CPU usage
    Given compression operations
    When LZMA compression is used
    Then CPU usage is highest among compression types
    And compression requires intensive processing
    And CPU resources are used extensively

  @REQ-COMPR-033 @REQ-COMPR-034 @happy
  Scenario: LZMA compression is best for long-term storage
    Given compression operations for storage
    When LZMA compression is selected
    Then compression is best for long-term storage scenarios
    And maximum space savings are achieved
    And storage efficiency is optimized

  @REQ-COMPR-033 @happy
  Scenario: LZMA compression provides maximum space savings
    Given compression operations
    When LZMA compression is used
    Then maximum space savings are achieved
    And compressed file size is minimized
    And storage requirements are reduced to minimum

  @REQ-COMPR-034 @happy
  Scenario: LZMA compression prioritizes size over speed
    Given compression operations
    When LZMA compression is used
    Then size priority over speed is achieved
    And compression takes longer but produces smaller files
    And trade-off favors storage efficiency
