@domain:file_mgmt @m2 @REQ-FILEMGMT-019 @spec(api_file_management.md#8-file-encryption-operations)
Feature: File Management: File Encryption Operations

  @happy
  Scenario: EncryptFile encrypts file content
    Given an open writable NovusPack package with unencrypted file
    When EncryptFile is called with encryption type and key
    Then file content is encrypted
    And EncryptionType is set
    And file is protected
    And file can be decrypted with correct key

  @happy
  Scenario: DecryptFile decrypts file content
    Given an open NovusPack package with encrypted file
    When DecryptFile is called with correct key
    Then file content is decrypted
    And original content is restored
    And file content matches original

  @happy
  Scenario: File encryption supports multiple encryption types
    Given file encryption functionality
    When encryption types are examined
    Then AES-256-GCM encryption is supported
    And quantum-safe encryption (ML-KEM + ML-DSA) is supported
    And encryption types are configurable

  @error
  Scenario: DecryptFile fails with incorrect key
    Given an open NovusPack package with encrypted file
    When DecryptFile is called with incorrect key
    Then structured encryption error is returned
    And error indicates decryption failure

  @error
  Scenario: Encryption operations respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When encryption operation is called
    Then structured context error is returned
