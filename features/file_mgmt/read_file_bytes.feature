@domain:file_mgmt @m2 @REQ-CORE-074 @spec(api_core.md#113-readfile-method-contract) @spec(api_file_mgmt_extraction.md#5142-multi-stage-pipeline-extraction-large-files) @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines) @spec(api_file_mgmt_transform_pipelines.md#11-pipeline-structure) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions) @spec(api_file_mgmt_transform_pipelines.md#112-pipeline-validation) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: Read file bytes

  @REQ-CORE-077 @happy
  Scenario: ReadFile obeys encryption and validation
    Given a package with an encrypted file "secret.txt"
    When I read the file with correct keys
    Then the read bytes should match the original content

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: ReadFile validates path parameter
    Given an open package
    When ReadFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: ReadFile respects context cancellation
    Given an open package with files
    And a cancelled context
    When ReadFile is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-342 @REQ-PIPELINE-014 @happy
  Scenario: Large file reading with automatic multi-stage pipeline
    Given an open NovusPack package
    And a large compressed and encrypted file
    When ReadFile is called
    Then system creates transformation pipeline automatically
    And pipeline includes decrypt and decompress stages
    And ProcessingState transitions through pipeline stages
    And reading completes without loading entire file in memory

  @REQ-FILEMGMT-342 @REQ-PIPELINE-011 @happy
  Scenario: Read large file uses temporary files based on threshold
    Given an open NovusPack package
    And a large file exceeding memory threshold
    When ReadFile is called
    Then file is read using transformation pipeline
    And temporary files are created for intermediate stages
    And CurrentSource tracks progression through pipeline
    And all temporary files are cleaned up after reading

  @REQ-PIPELINE-012 @REQ-PIPELINE-013 @happy
  Scenario: ProcessingState transitions correctly during multi-stage reading
    Given an open NovusPack package
    And a file stored as CompressedAndEncrypted
    When multi-stage reading begins
    Then initial ProcessingState is CompressedAndEncrypted
    When decrypt stage completes
    Then ProcessingState transitions to Compressed
    When decompress stage completes
    Then ProcessingState transitions to Raw
    And final content is available
