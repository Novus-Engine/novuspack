@domain:file_mgmt @m2 @REQ-FILEMGMT-023 @spec(api_file_management.md#21-add-file)
Feature: Add file with encryption

  @happy
  Scenario: AddFileWithEncryption adds file with specific encryption type
    Given an open writable package
    When AddFileWithEncryption is called with file source and encryption type
    Then file is added with specified encryption type
    And encryption type is stored in file entry
    And file can be decrypted with appropriate key

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
  Scenario: AddFileWithEncryption fails with invalid encryption type
    Given an open writable package
    When AddFileWithEncryption is called with invalid encryption type
    Then structured validation error is returned
    And error indicates unsupported encryption type

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

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: Encryption-aware operations respect context cancellation
    Given an open writable package
    And a cancelled context
    When encryption-aware operation is called
    Then structured context error is returned
    And error type is context cancellation
