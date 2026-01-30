@domain:file_mgmt @REQ-FILEMGMT-442 @REQ-FILEMGMT-200 @spec(api_file_mgmt_addition.md#addfilefrommemory-data-management) @spec(api_file_mgmt_addition.md#22-packageaddfilefrommemory-method)
Feature: AddFileFromMemory Data Management

  AddFileFromMemory treats the provided data as raw content: calculates
  RawChecksum and OriginalSize, writes to temp file for pipeline, uses temp
  file as CurrentSource for Write. Caller retains data ownership until
  Write completes or data is processed.

  @REQ-FILEMGMT-442 @happy
  Scenario: Data is treated as raw content
    Given an open writable package
    And raw byte data (uncompressed, unencrypted)
    When AddFileFromMemory is called with path, data, and options
    Then RawChecksum is calculated from raw data for deduplication
    And OriginalSize is set from len(data)
    And raw data is written to temporary file for consistent pipeline
    And temp file serves as FileEntry.CurrentSource for Write operations

  @REQ-FILEMGMT-442 @happy
  Scenario: Encryption and compression applied per options during processing
    Given an open writable package
    And options with encryption or compression
    When AddFileFromMemory is called
    Then encryption replaces temp file with encrypted output when used
    And compression and encryption follow same rules as AddFile
    And temporary files are managed by package lifecycle and cleaned on Close or after Write

  @REQ-FILEMGMT-442 @happy
  Scenario: Caller retains data ownership until processing complete
    Given an open writable package
    And a byte slice from caller
    When AddFileFromMemory is called with that slice
    Then only slice header is copied (reference type)
    And caller must keep data valid until Write completes or data is processed for encryption
