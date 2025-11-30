@domain:compression @m2 @REQ-COMPR-119 @spec(api_package_compression.md#512-compresspackagestream-parameters)
Feature: CompressPackageStream Parameter Specification

  @REQ-COMPR-119 @happy
  Scenario: CompressPackageStream accepts context parameter
    Given a CompressPackageStream operation
    When CompressPackageStream is called
    Then context parameter is accepted
    And context supports cancellation and timeout handling
    And context enables operation control

  @REQ-COMPR-119 @happy
  Scenario: CompressPackageStream accepts compression type parameter
    Given a CompressPackageStream operation
    When CompressPackageStream is called
    Then compression type parameter is accepted
    And compression type specifies algorithm (1-3)
    And compression type determines compression method

  @REQ-COMPR-119 @happy
  Scenario: CompressPackageStream accepts StreamConfig parameter
    Given a CompressPackageStream operation
    When CompressPackageStream is called
    Then StreamConfig parameter is accepted
    And config provides unified streaming configuration
    And config manages memory and optimization settings
    And config enables fine-tuned control

  @REQ-COMPR-119 @happy
  Scenario: StreamConfig provides memory management configuration
    Given CompressPackageStream with StreamConfig
    When config is provided
    Then memory management is configured
    And MaxMemoryUsage controls memory limits
    And ChunkSize controls processing chunks
    And TempDir specifies temporary file location
