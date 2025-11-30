@domain:compression @m2 @REQ-COMPR-100 @spec(api_package_compression.md#22-built-in-compression-strategies)
Feature: Built-In Compression Strategies

  @REQ-COMPR-100 @happy
  Scenario: Built-in compression strategies provide Zstandard implementation
    Given a NovusPack package
    And Zstandard compression type is selected
    When compression is performed
    Then ZstandardStrategy is used
    And Zstandard algorithm is applied
    And compression succeeds

  @REQ-COMPR-100 @happy
  Scenario: Built-in compression strategies provide LZ4 implementation
    Given a NovusPack package
    And LZ4 compression type is selected
    When compression is performed
    Then LZ4Strategy is used
    And LZ4 algorithm is applied
    And fast compression is achieved

  @REQ-COMPR-100 @happy
  Scenario: Built-in compression strategies provide LZMA implementation
    Given a NovusPack package
    And LZMA compression type is selected
    When compression is performed
    Then LZMAStrategy is used
    And LZMA algorithm is applied
    And maximum compression ratio is achieved

  @REQ-COMPR-100 @happy
  Scenario: Built-in compression strategies support generic types
    Given a compression operation with generic data type
    When built-in strategy is used
    Then strategy supports generic type parameter
    And compression works with various data types
    And type safety is maintained

  @REQ-COMPR-100 @happy
  Scenario: Built-in compression strategies support compression levels
    Given a NovusPack package
    And a specific compression level is configured
    When built-in strategy is used
    Then compression level is applied
    And strategy respects level configuration
    And compression ratio matches level
