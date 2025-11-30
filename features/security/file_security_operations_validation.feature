@domain:security @m2 @REQ-SEC-033 @REQ-SEC-090 @spec(security.md#42-file-validation)
Feature: File Security Operations

  @REQ-SEC-033 @happy
  Scenario: File validation validates individual file integrity
    Given an open NovusPack package
    And file entries in package
    And a valid context
    When file validation is performed
    Then individual file integrity is validated
    And each file is validated separately
    And file validation results are returned

  @REQ-SEC-033 @happy
  Scenario: File validation validates file checksums
    Given an open NovusPack package
    And file entries with checksums
    And a valid context
    When file validation is performed
    Then file checksums are validated
    And checksum validation ensures integrity
    And invalid checksums are detected

  @REQ-SEC-033 @happy
  Scenario: File validation validates file metadata
    Given an open NovusPack package
    And file entries with metadata
    And a valid context
    When file validation is performed
    Then file metadata is validated
    And metadata validation ensures correctness
    And invalid metadata is detected

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide file-level encryption
    Given an open writable NovusPack package
    And file entry
    And encryption handler
    And encryption key
    And a valid context
    When EncryptFile is called
    Then file-level encryption is performed
    And file is encrypted using handler
    And encryption key is used

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide file-level decryption
    Given an open NovusPack package
    And encrypted file entry
    And encryption handler
    And encryption key
    And a valid context
    When DecryptFile is called
    Then file-level decryption is performed
    And file is decrypted using handler
    And encryption key is used
    And decrypted data is returned

  @REQ-SEC-090 @happy
  Scenario: File encryption operations validate file encryption
    Given an open NovusPack package
    And encrypted file entry
    And encryption handler
    And a valid context
    When ValidateFileEncryption is called
    Then file encryption is validated
    And encryption validation ensures correctness
    And validation result is returned

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide encryption information
    Given an open NovusPack package
    And encrypted file entry
    And a valid context
    When GetFileEncryptionInfo is called
    Then EncryptionConfig structure is returned
    And encryption information is accessible
    And encryption details are provided

  @REQ-SEC-007 @REQ-SEC-011 @error
  Scenario: File encryption operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When file encryption operation is called
    Then structured context error is returned
    And error type is context cancellation
