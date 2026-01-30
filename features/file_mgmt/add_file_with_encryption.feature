@domain:file_mgmt @m2 @REQ-FILEMGMT-023 @spec(api_file_mgmt_addition.md#216-convenience-wrapper-addfilewithencryption) @spec(api_file_mgmt_transform_pipelines.md#12-processingstate-transitions) @spec(api_file_mgmt_transform_pipelines.md#17-temporary-file-security) @spec(file_validation.md#11-file-name-validation) @spec(api_file_mgmt_file_entry.md#9-fileentry-encryption)
Feature: Add file with encryption

  @happy
  Scenario: AddFileWithEncryption adds file with encryption key
    Given an open writable package
    And a valid encryption key
    And a filesystem file path
    When AddFileWithEncryption is called with path, encryption key, and nil options
    Then file is added with encryption enabled
    And encryption key is associated with file entry
    And file can be decrypted with the same key
    And AddFileWithEncryption is a convenience wrapper around AddFile
    And returned FileEntry matches the added file

  @happy
  Scenario: GetFileEncryptionType returns encryption type
    Given an open package with encrypted file
    When GetFileEncryptionType is called with file path
    Then encryption type is returned
    And type matches file entry encryption type

  @happy
  Scenario: GetEncryptedFiles returns list of encrypted files
    Given an open package with encrypted and unencrypted files
    When GetEncryptedFiles is called
    Then list of encrypted file paths is returned
    And all encrypted files are included
    And unencrypted files are excluded

  @error
  Scenario: AddFileWithEncryption fails with invalid encryption key
    Given an open writable package
    When AddFileWithEncryption is called with invalid encryption key
    Then structured validation error is returned
    And error indicates invalid or expired encryption key

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: Encryption-aware operations validate path parameter
    Given an open writable package
    When AddFileWithEncryption or GetFileEncryptionType is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-040 @error
  Scenario: Encryption-aware operations validate encryption type parameter
    Given an open writable package
    When AddFileWithEncryption is called with invalid encryption type
    Then structured validation error is returned
    And error indicates unsupported encryption type

  @happy
  Scenario: AddFileWithEncryption accepts additional options
    Given an open writable package
    And a valid encryption key
    And a filesystem file path
    And AddFileOptions with compression enabled
    When AddFileWithEncryption is called with path, encryption key, and options
    Then file is added with encryption and compression enabled
    And encryption key is associated with file entry
    And file is compressed and encrypted

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: Encryption-aware operations respect context cancellation
    Given an open writable package
    And a cancelled context
    When encryption-aware operation is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-CRYPTO-012 @REQ-PIPELINE-013 @happy
  Scenario: Encryption works as pipeline stage with compression
    Given an open writable package
    And a large file
    And AddFileOptions with compression and encryption enabled
    When AddFile is called
    Then compression stage executes first (Raw to Compressed)
    And encryption stage executes second (Compressed to CompressedAndEncrypted)
    And pipeline stages execute sequentially
    And final encrypted file written to package

  @REQ-CRYPTO-013 @REQ-PIPELINE-025 @happy
  Scenario: Encrypted temp files use context-aware security
    Given an open writable package
    And AddFileOptions with encryption enabled
    When large file is added using pipeline
    Then intermediate temp files are secured appropriately
    And temp file security follows context-aware rules
    And sensitive data is protected
