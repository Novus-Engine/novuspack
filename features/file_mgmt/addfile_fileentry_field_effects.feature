@domain:file_mgmt @REQ-FILEMGMT-434 @spec(api_file_mgmt_addition.md#fileentry-field-effects) @spec(api_file_mgmt_addition.md#21-packageaddfile-method)
Feature: AddFile FileEntry Field Effects

  AddFile creates or updates a FileEntry in the in-memory package state
  following the specified sequence: filesystem validation, deduplication
  check, conditional encryption processing, FileEntry allocation (unique),
  and runtime field finalization.

  @REQ-FILEMGMT-434 @happy
  Scenario: AddFile executes filesystem validation and metadata read first
    Given an open writable package
    And a valid filesystem file path
    When AddFile is called with the path
    Then filesystem path is validated
    And source file is opened and staged as FileEntry.CurrentSource
    And file metadata is read
    And file type is determined
    And OriginalSize is calculated
    And stored package path is derived

  @REQ-FILEMGMT-434 @happy
  Scenario: AddFile performs deduplication check before processing
    Given an open writable package
    And AddFileOptions.AllowDuplicate is false
    And a file whose raw content may match an existing entry
    When AddFile is called
    Then deduplication check runs before encryption or compression
    And if duplicate found, path is added to existing FileEntry.Paths
    And PathCount and MetadataVersion are updated
    And no new FileEntry is allocated for duplicate

  @REQ-FILEMGMT-434 @happy
  Scenario: AddFile allocates new FileEntry for unique files
    Given an open writable package
    And a unique file (no matching raw content)
    When AddFile is called
    Then new FileEntry is allocated
    And FileID is assigned
    And Paths contains derived stored path
    And PathCount is 1
    And OriginalSize, RawChecksum, CompressionType, EncryptionType are set
    And runtime CurrentSource and ProcessingState are finalized

  @REQ-FILEMGMT-434 @happy
  Scenario: AddFile does not populate Data or IsDataLoaded
    Given an open writable package
    And a filesystem file
    When AddFile completes successfully
    Then FileEntry.Data is not populated
    And FileEntry.IsDataLoaded is false
    And CurrentSource and offset/size are used for Write operations
