@domain:file_mgmt @m2 @REQ-CORE-074 @spec(api_core.md#113-readfile-method-contract) @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines) @spec(api_file_mgmt_transform_pipelines.md#11-pipeline-structure) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions) @spec(api_file_mgmt_transform_pipelines.md#112-pipeline-validation) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: Read File Operations

  @REQ-CORE-074 @REQ-CORE-076 @happy
  Scenario: ReadFile operations support reading file content from packages
    Given an open NovusPack package
    And file "documents/data.txt" exists in the package
    When ReadFile is called with the file path
    Then file content is read and returned
    And content is returned as byte slice
    And reading operation completes successfully

  @REQ-CORE-074 @happy
  Scenario: ReadFile automatically handles decompression
    Given an open package
    And compressed file "archive.zip" exists
    When ReadFile is called
    Then file content is decompressed automatically
    And decompressed content is returned
    And decompression is transparent to caller

  @REQ-CORE-074 @happy
  Scenario: ReadFile automatically handles decryption
    Given an open package
    And encrypted file "secret.txt" exists
    And decryption keys are available
    When ReadFile is called
    Then file content is decrypted automatically
    And decrypted content is returned
    And decryption is transparent to caller

  @REQ-CORE-074 @happy
  Scenario: ReadFile handles combined compression and encryption
    Given an open package
    And file exists with both compression and encryption
    When ReadFile is called
    Then file is decrypted first
    Then file is decompressed
    And final content is returned
    And operations occur in correct order

  @REQ-CORE-074 @happy
  Scenario: ReadFile locates file entry in package index
    Given an open package
    And file entry exists in package index
    When ReadFile is called with valid path
    Then file entry is located in package index
    And file content is read from data section
    And reading completes successfully

  @REQ-CORE-074 @error
  Scenario: ReadFile returns error when package not open
    Given a package that is not open
    When ReadFile is called
    Then ErrPackageNotOpen error is returned
    And error follows structured error format

  @REQ-CORE-074 @error
  Scenario: ReadFile returns error when file not found
    Given an open package
    And file does not exist at specified path
    When ReadFile is called with non-existent path
    Then ErrFileNotFound error is returned
    And error indicates file not found
    And error follows structured error format

  @REQ-CORE-074 @error
  Scenario: ReadFile returns error when decryption fails
    Given an open package
    And encrypted file exists
    And decryption keys are invalid or missing
    When ReadFile is called
    Then ErrDecryptionFailed error is returned
    And error indicates decryption failure
    And error follows structured error format

  @REQ-CORE-074 @error
  Scenario: ReadFile respects context cancellation
    Given an open package
    And a cancelled context
    When ReadFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format

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
