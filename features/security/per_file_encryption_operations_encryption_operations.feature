@domain:security @m2 @REQ-SEC-026 @spec(security.md#342-per-file-encryption-operations)
Feature: Per-File Encryption Operations

  @REQ-SEC-026 @happy
  Scenario: Per-file encryption operations provide AddFileWithEncryption
    Given an open writable NovusPack package
    And a valid context
    And file data and encryption type
    When AddFileWithEncryption is called
    Then files are added with specific encryption types
    And encryption type can be None, AES-256-GCM, or ML-KEM
    And file encryption is configured per file
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-026 @happy
  Scenario: Per-file encryption operations provide GetFileEncryptionType
    Given an open NovusPack package
    And a valid context
    And file in package
    When GetFileEncryptionType is called
    Then encryption type for specific file is retrieved
    And returned type indicates None, AES-256-GCM, or ML-KEM
    And encryption type information is accurate

  @REQ-SEC-026 @happy
  Scenario: Per-file encryption operations provide GetEncryptedFiles
    Given an open NovusPack package
    And a valid context
    And package with mixed encrypted/unencrypted files
    When GetEncryptedFiles is called
    Then all encrypted files in package are listed
    And list contains only encrypted files
    And unencrypted files are excluded from list

  @REQ-SEC-026 @happy
  Scenario: Per-file encryption operations provide granular control
    Given an open writable NovusPack package
    And a valid context
    And multiple files requiring encryption
    When per-file encryption operations are used
    Then each file can have different encryption type
    And encryption selection is file-specific
    And granular control enables selective encryption

  @REQ-SEC-026 @happy
  Scenario: Per-file encryption operations support mixed packages
    Given an open writable NovusPack package
    And a valid context
    And files with different encryption requirements
    When per-file encryption operations are applied
    Then packages can contain both encrypted and unencrypted files
    And mixed encryption is supported
    And each file maintains its encryption configuration

  @REQ-SEC-026 @error
  Scenario: Per-file encryption operations handle invalid encryption types
    Given an open writable NovusPack package
    And a valid context
    And invalid encryption type
    When AddFileWithEncryption is called with invalid type
    Then validation error is returned
    And error indicates invalid encryption type
    And error follows structured error format
