@domain:file_mgmt @REQ-PIPELINE-013 @spec(api_file_mgmt_file_entry.md#15-processingstate-type) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: ProcessingState Data Model

  Background:
    Given ProcessingState uses data-state model not workflow model

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingStateRaw indicates unprocessed data
    Given file data is raw (unprocessed)
    When ProcessingState is checked
    Then ProcessingState is Raw
    And no compression or encryption applied
    And data is in original form

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingStateCompressed indicates compressed data
    Given file data is compressed but not encrypted
    When ProcessingState is checked
    Then ProcessingState is Compressed
    And compression has been applied
    And encryption has not been applied
    And additional encryption may be needed

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingStateEncrypted indicates encrypted data
    Given file data is encrypted but not compressed
    When ProcessingState is checked
    Then ProcessingState is Encrypted
    And encryption has been applied
    And compression has not been applied
    And this state is rarely used

  @REQ-PIPELINE-013 @happy
  Scenario: ProcessingStateCompressedAndEncrypted indicates both transformations
    Given file data is both compressed and encrypted
    When ProcessingState is checked
    Then ProcessingState is CompressedAndEncrypted
    And both transformations have been applied
    And data is in final processed form
    And ready for package storage or no further processing needed

  @REQ-PIPELINE-013 @happy
  Scenario: Data-state model describes actual data state
    Given ProcessingState value
    Then ProcessingState indicates what transformations applied
    And ProcessingState does not indicate workflow phase
    And ProcessingState informs what processing is still needed
    And different from workflow states (Idle, Loading, Processing, etc.)
