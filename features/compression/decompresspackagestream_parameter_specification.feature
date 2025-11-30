@domain:compression @m2 @REQ-COMPR-124 @spec(api_package_compression.md#522-decompresspackagestream-parameters)
Feature: DecompressPackageStream Parameter Specification

  @REQ-COMPR-124 @happy
  Scenario: DecompressPackageStream accepts context parameter
    Given a DecompressPackageStream operation
    When DecompressPackageStream is called
    Then context parameter is accepted
    And context supports cancellation and timeout handling
    And context enables operation control

  @REQ-COMPR-124 @happy
  Scenario: DecompressPackageStream accepts StreamConfig parameter
    Given a DecompressPackageStream operation
    When DecompressPackageStream is called
    Then StreamConfig parameter is accepted
    And config provides streaming configuration for memory management
    And config enables fine-tuned control over decompression

  @REQ-COMPR-124 @happy
  Scenario: StreamConfig manages memory during decompression
    Given DecompressPackageStream with StreamConfig
    When config is provided
    Then memory management is configured for decompression
    And MaxMemoryUsage controls memory limits
    And ChunkSize controls processing chunks
    And TempDir specifies temporary file location for streaming
