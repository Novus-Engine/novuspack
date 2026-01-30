@domain:file_mgmt @REQ-FILEMGMT-441 @REQ-FILEMGMT-200 @spec(api_file_mgmt_addition.md#addfilefrommemory-fileentry-field-effects) @spec(api_file_mgmt_addition.md#22-packageaddfilefrommemory-method)
Feature: AddFileFromMemory FileEntry Field Effects

  AddFileFromMemory follows the same FileEntry creation sequence as AddFile,
  with data-source differences: validate path, write raw data to temp file,
  set CurrentSource, then deduplication, encryption processing, allocation,
  runtime finalization.

  @REQ-FILEMGMT-441 @happy
  Scenario: Data source and validation use provided byte data
    Given an open writable package
    And raw byte data and package path
    When AddFileFromMemory is called with path, data, and options
    Then path is validated per package path rules
    And OriginalSize is calculated from len(data)
    And raw data is written to temporary file for pipeline
    And FileEntry.CurrentSource is set to temp-file FileSource with offset 0
    And file type is determined from path and content

  @REQ-FILEMGMT-441 @happy
  Scenario: Deduplication and encryption processing match AddFile
    Given an open writable package
    And byte data and options with or without encryption
    When AddFileFromMemory is called
    Then deduplication check runs same as AddFile
    And encryption processing applies when options require it
    And FileEntry allocation and runtime finalization follow AddFile sequence

  @REQ-FILEMGMT-441 @happy
  Scenario: Runtime field finalization sets CurrentSource and ProcessingState
    Given an open writable package
    And byte data
    When AddFileFromMemory completes successfully
    Then FileEntry.CurrentSource is set and valid
    And FileEntry.CurrentSource.IsTempFile is set appropriately
    And ProcessingState reflects data state in CurrentSource
