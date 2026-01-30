@domain:file_mgmt @m2 @REQ-FILEMGMT-200 @spec(api_file_mgmt_addition.md#22-addfilefrommemory) @spec(api_file_mgmt_transform_pipelines.md#1-multi-stage-transformation-pipelines) @spec(api_file_mgmt_transform_pipelines.md#112-pipeline-validation) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions)
Feature: AddFileFromMemory Basic

  @REQ-FILEMGMT-200 @REQ-FILEMGMT-201 @happy
  Scenario: Add file from byte data
    Given in-memory byte data for a file
    When AddFileFromMemory is called with path and data
    Then file is added to package
    And FileEntry is created with provided data
    And file content comes from memory not filesystem

  @REQ-FILEMGMT-202 @happy
  Scenario: Package-relative path handling
    Given path "/config/settings.json"
    When AddFileFromMemory is called
    Then path is treated as package-relative
    And path is normalized to NFC
    And path follows package path rules
    And leading slash is required for storage

  @REQ-FILEMGMT-202 @happy
  Scenario: Path normalization follows package semantics
    Given path without leading slash "data/file.txt"
    When AddFileFromMemory is called
    Then path is normalized with leading slash
    And stored path is "/data/file.txt"
    And package path semantics are followed

  @REQ-FILEMGMT-203 @happy
  Scenario: Returns FileEntry with metadata and checksums
    Given byte data for file
    When AddFileFromMemory successfully adds file
    Then FileEntry is returned
    And FileEntry contains complete metadata
    And FileEntry has file checksums computed
    And FileEntry has compression status
    And FileEntry has encryption details if encrypted

  @REQ-FILEMGMT-204 @happy
  Scenario: Validates path and uses provided byte data
    Given package-relative path and byte data
    When AddFileFromMemory is called
    Then path is validated for correctness
    And path format is checked
    And provided byte data is used as raw content
    And data is assumed uncompressed and unencrypted

  @REQ-FILEMGMT-204 @happy
  Scenario: Data is raw file content
    Given uncompressed unencrypted byte data
    When AddFileFromMemory is called
    Then data is treated as raw file content
    And data is not interpreted as compressed
    And data is not interpreted as encrypted
    And data is used directly for file entry

  @REQ-FILEMGMT-201 @happy
  Scenario: AddFileFromMemory parameters
    Given context for cancellation
    And package-relative path
    And pointer to byte slice data
    And optional AddFileOptions
    When AddFileFromMemory is called
    Then all parameters are processed correctly
    And context is respected for cancellation
    And options control processing behavior

  @REQ-FILEMGMT-200 @happy
  Scenario: Use case - adding generated content
    Given content generated in memory
    When content needs to be added without filesystem
    Then AddFileFromMemory provides direct API
    And filesystem is not involved
    And content added efficiently from memory

  @REQ-FILEMGMT-342 @happy
  Scenario: AddFileFromMemory uses pipeline for large data with transformations
    Given an open writable package
    And large byte array exceeding memory threshold
    And AddFileOptions with compression and encryption
    When AddFileFromMemory is called
    Then multi-stage pipeline processes data
    And temporary files used for intermediate stages
    And ProcessingState transitions through pipeline stages
    And all temporary files cleaned up after write
