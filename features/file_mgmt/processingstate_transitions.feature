@domain:file_mgmt @REQ-PIPELINE-013 @spec(api_file_mgmt_file_entry.md#15-processingstate-type) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: ProcessingState Transitions

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingState transitions during file addition with compress and encrypt
    Given a raw file being added to package
    And AddFileOptions with compression and encryption
    Then initial ProcessingState is Raw
    When compress stage completes
    Then ProcessingState transitions to Compressed
    When encrypt stage completes
    Then ProcessingState transitions to CompressedAndEncrypted
    And file ready for package write

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingState transitions during file extraction with decrypt and decompress
    Given a file in package stored as CompressedAndEncrypted
    Then initial ProcessingState is CompressedAndEncrypted
    When decrypt stage completes
    Then ProcessingState transitions to Compressed
    When decompress stage completes
    Then ProcessingState transitions to Raw
    And file ready for user consumption

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingState reflects current data state at each stage
    Given a transformation pipeline
    When each stage completes
    Then ProcessingState is updated to reflect current transformations
    And ProcessingState accurately describes data in CurrentSource
    And subsequent operations use ProcessingState to determine needed transformations

  @REQ-PIPELINE-013 @happy
  Scenario: Extraction reverses addition transformations
    Given addition pipeline: Raw => Compressed => CompressedAndEncrypted
    And extraction pipeline: CompressedAndEncrypted => Compressed => Raw
    Then extraction reverses addition sequence
    And ProcessingState transitions are inverse
    And final state matches original (Raw)

  @REQ-PIPELINE-013 @happy
  Scenario: Single transformation updates ProcessingState appropriately
    Given a file with compression only (no encryption)
    When compress operation completes
    Then ProcessingState transitions from Raw to Compressed
    And ProcessingState skips CompressedAndEncrypted
    And reflects actual transformations applied
