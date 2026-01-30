@domain:file_mgmt @REQ-PIPELINE-009 @REQ-PIPELINE-010 @spec(api_file_mgmt_transform_pipelines.md#11-pipeline-structure) @spec(api_file_mgmt_extraction.md#5142-multi-stage-pipeline-extraction-large-files) @spec(api_file_mgmt_file_entry.md#1-fileentry-structure) @spec(api_file_mgmt_file_entry.md#45-multi-stage-transformation-pipeline)
Feature: Transformation Pipeline Stages

  @REQ-PIPELINE-009 @happy
  Scenario: TransformType identifies transformation types
    Given transformation type enumeration
    Then TransformTypeNone is 0x00 (no transformation)
    And TransformTypeCompress is 0x01 (compression)
    And TransformTypeDecompress is 0x02 (decompression)
    And TransformTypeEncrypt is 0x03 (encryption)
    And TransformTypeDecrypt is 0x04 (decryption)
    And TransformTypeVerify is 0x05 (checksum verification)
    And TransformTypeCustom is 0xFF (custom transformation)

  @REQ-PIPELINE-010 @happy
  Scenario: ExecuteTransformStage executes specific stage by index
    Given a FileEntry with TransformPipeline
    And pipeline has 3 stages
    When ExecuteTransformStage is called with stage index 1
    Then stage 1 is executed
    And stage reads from InputSource
    And stage writes to OutputSource (new temporary file)
    And stage completion status is updated
    And CurrentSource is updated to stage output

  @REQ-PIPELINE-011 @happy
  Scenario: Stage reads from InputSource and writes to OutputSource
    Given a transformation stage
    When stage executes
    Then data is read from InputSource location
    And data is transformed according to StageType
    And transformed data is written to OutputSource
    And OutputSource is a temporary file for intermediate stages

  @REQ-PIPELINE-014 @happy
  Scenario: System determines transformation sequence automatically
    Given a file with CompressionType Zstd and EncryptionType AES256
    When adding file to package
    Then system creates pipeline with compress stage first
    And system creates pipeline with encrypt stage second
    When extracting file from package
    Then system creates pipeline with decrypt stage first
    And system creates pipeline with decompress stage second
