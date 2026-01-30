@domain:file_mgmt @m2 @REQ-CORE-074 @spec(api_core.md#113-readfile-method-contract) @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines) @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management) @spec(api_file_mgmt_transform_pipelines.md#16-intermediate-stage-cleanup) @spec(api_file_mgmt_transform_pipelines.md#18-error-recovery-and-resume) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions) @spec(api_file_mgmt_transform_pipelines.md#14-disk-space-management) @spec(api_file_mgmt_transform_pipelines.md#16-intermediate-stage-cleanup) @spec(api_file_mgmt_transform_pipelines.md#18-error-recovery-and-resume) @spec(api_file_mgmt_transform_pipelines.md#112-pipeline-validation) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: ReadFile extracts file content

  @REQ-CORE-074 @REQ-CORE-076 @REQ-CORE-077 @happy
  Scenario: ReadFile returns file content as byte slice
    Given an open NovusPack package
    And a file exists in the package at a specific path
    When ReadFile is called with the file path
    Then file content is returned as byte slice
    And content matches original file data
    And content is complete and unmodified

  @REQ-CORE-074 @happy
  Scenario: ReadFile handles compressed files
    Given an open NovusPack package
    And a compressed file exists in the package
    When ReadFile is called with the file path
    Then file content is automatically decompressed
    And decompressed content is returned as byte slice
    And content matches original uncompressed data

  @REQ-CORE-074 @happy
  Scenario: ReadFile handles encrypted files
    Given an open NovusPack package
    And an encrypted file exists in the package
    And encryption key is available
    When ReadFile is called with the file path
    Then file content is automatically decrypted
    And decrypted content is returned as byte slice
    And content matches original unencrypted data

  @REQ-CORE-074 @error
  Scenario: ReadFile returns error when package is not open
    Given a NovusPack package that is not open
    When ReadFile is called with a file path
    Then package not open error is returned
    And error indicates package must be open

  @REQ-CORE-074 @error
  Scenario: ReadFile returns error when file does not exist
    Given an open NovusPack package
    And a file does not exist at the specified path
    When ReadFile is called with the non-existent path
    Then file not found error is returned
    And error indicates file does not exist

  @REQ-CORE-074 @error
  Scenario: ReadFile returns error for invalid path
    Given an open NovusPack package
    When ReadFile is called with invalid or malformed path
    Then invalid path error is returned
    And error indicates path format issue

  @REQ-CORE-074 @error
  Scenario: ReadFile handles decompression failures
    Given an open NovusPack package
    And a file with corrupted compression data
    When ReadFile is called
    Then decompression failed error is returned
    And error indicates compression issue

  @REQ-CORE-074 @error
  Scenario: ReadFile handles decryption failures
    Given an open NovusPack package
    And an encrypted file with incorrect key
    When ReadFile is called
    Then decryption failed error is returned
    And error indicates decryption issue

  @REQ-CORE-074 @error
  Scenario: ReadFile respects context cancellation
    Given an open NovusPack package
    And a cancelled context
    When ReadFile is called with cancelled context
    Then context cancelled error is returned
    And error type is context cancellation

  @REQ-CORE-074 @error
  Scenario: ReadFile handles context timeout
    Given an open NovusPack package
    And a context with timeout
    And a large file that exceeds timeout
    When ReadFile is called
    Then context timeout error is returned
    And error type is context timeout

  @REQ-FILEMGMT-342 @REQ-PIPELINE-013 @happy
  Scenario: ReadFile uses multi-stage pipeline for large encrypted and compressed files
    Given an open NovusPack package
    And a large file that is both compressed and encrypted
    And the file size exceeds memory threshold
    And encryption key is available
    When ReadFile is called with the file path
    Then a multi-stage transformation pipeline is created
    And stage 1 decrypts data to temporary file
    And stage 2 decompresses data to temporary file
    And ProcessingState transitions from CompressedAndEncrypted to Compressed to Raw
    And final raw content is returned
    And all temporary files are cleaned up

  @REQ-FILEMGMT-343 @REQ-PIPELINE-015 @error
  Scenario: ReadFile checks disk space before multi-stage extraction
    Given an open NovusPack package
    And a large file requiring multi-stage pipeline
    And insufficient disk space for temporary files
    When ReadFile is called
    Then disk space error is returned
    And error indicates insufficient space
    And no temporary files are created

  @REQ-FILEMGMT-345 @REQ-PIPELINE-023 @happy
  Scenario: ReadFile cleans up temporary files on pipeline failure
    Given an open NovusPack package
    And a large file with multi-stage pipeline
    When extraction fails during stage 2
    Then stage 1 temporary file is removed
    And partial stage 2 temporary file is removed
    And no temporary files remain
    And error indicates which stage failed

  @REQ-FILEMGMT-351 @REQ-PIPELINE-032 @happy
  Scenario: ReadFile supports resume after interruption
    Given an open NovusPack package
    And a large file with multi-stage pipeline
    And extraction was interrupted after stage 1 completed
    When ResumeTransformation is called
    Then stage 1 output is reused from temporary file
    And stage 2 continues from last checkpoint
    And extraction completes successfully
    And all temporary files are cleaned up
