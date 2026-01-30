@domain:file_mgmt @REQ-FILEMGMT-435 @spec(api_file_mgmt_addition.md#file-data-processing-model) @spec(api_file_mgmt_addition.md#21-packageaddfile-method)
Feature: AddFile File Data Processing Model

  AddFile applies different processing based on compression and encryption
  settings: encryption-only, compression-only, both, or neither. Encryption
  is applied early; compression is deferred to Write when no encryption.

  @REQ-FILEMGMT-435 @happy
  Scenario: Encryption-only processing encrypts and stages to temp file
    Given an open writable package
    And AddFileOptions with encryption and no compression
    And a filesystem file
    When AddFile is called
    Then file data is encrypted
    And encrypted data is written to temporary file
    And FileEntry.CurrentSource points to temp file
    And FileEntry.CurrentSource.IsTempFile is true
    And Write operations read encrypted data from temp file

  @REQ-FILEMGMT-435 @happy
  Scenario: Compression-only processing defers to Write
    Given an open writable package
    And AddFileOptions with compression and no encryption
    And a filesystem file
    When AddFile is called
    Then file data is not compressed by AddFile
    And FileEntry.CurrentSource points to original source file
    And FileEntry.CurrentSource.IsTempFile is false
    And compression is applied during Write operations

  @REQ-FILEMGMT-435 @happy
  Scenario: Compression and encryption applies compress then encrypt
    Given an open writable package
    And AddFileOptions with both compression and encryption
    And a filesystem file
    When AddFile is called
    Then compression is applied first
    And encryption is applied to compressed output
    And processed data is staged for Write
    And FileEntry.CurrentSource reflects staged output

  @REQ-FILEMGMT-435 @happy
  Scenario: No compression or encryption leaves source staged
    Given an open writable package
    And AddFileOptions with no compression and no encryption
    And a filesystem file
    When AddFile is called
    Then file data is not read or processed by AddFile
    And FileEntry.CurrentSource points to original file
    And processing is deferred to Write operations
