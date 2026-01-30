@domain:compression @m2 @REQ-COMPR-167 @spec(api_file_mgmt_transform_pipelines.md#2113-processingstate-transitions)
Feature: Compression pipeline stage sequencing

  @REQ-COMPR-167 @happy
  Scenario: Compression acts as a stage in the documented transformation sequence
    Given a transformation pipeline for file addition
    When compression and encryption are both required
    Then the stage sequence is compress then encrypt
    Given a transformation pipeline for file extraction
    When decryption and decompression are both required
    Then the stage sequence is decrypt then decompress
    And stage sequencing follows documented processing state transitions

