@domain:compression @m2 @REQ-COMPR-148 @spec(api_package_compression.md#91-compression-specific-configuration)
Feature: Compression-Specific Configuration

  @REQ-COMPR-148 @happy
  Scenario: Compression-specific configuration provides algorithm-specific settings
    Given a compression operation
    And a specific compression algorithm is selected
    When compression-specific configuration is applied
    Then algorithm-specific settings are used
    And configuration matches algorithm requirements
    And algorithm performs optimally

  @REQ-COMPR-148 @happy
  Scenario: Compression-specific configuration supports Zstandard settings
    Given a compression operation
    And Zstandard compression type is selected
    When compression-specific configuration is applied
    Then Zstandard-specific settings are available
    And compression level 1-22 can be configured
    And Zstandard-specific optimizations are applied

  @REQ-COMPR-148 @happy
  Scenario: Compression-specific configuration supports LZ4 settings
    Given a compression operation
    And LZ4 compression type is selected
    When compression-specific configuration is applied
    Then LZ4-specific settings are available
    And compression level 1-9 can be configured
    And LZ4-specific optimizations are applied

  @REQ-COMPR-148 @happy
  Scenario: Compression-specific configuration supports LZMA settings
    Given a compression operation
    And LZMA compression type is selected
    When compression-specific configuration is applied
    Then LZMA-specific settings are available
    And compression level 1-9 can be configured
    And LZMA-specific optimizations are applied

  @REQ-COMPR-148 @happy
  Scenario: Compression-specific configuration allows auto-selection of compression level
    Given a compression operation
    And compression level is set to 0 for auto-selection
    When compression-specific configuration is applied
    Then optimal compression level is automatically selected
    And algorithm chooses best level for data
    And balance between speed and ratio is achieved
