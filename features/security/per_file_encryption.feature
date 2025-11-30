@domain:security @m2 @REQ-SEC-024 @REQ-SEC-026 @REQ-SEC-037 @spec(security.md#33-per-file-encryption)
Feature: Per-File Encryption

  @REQ-SEC-024 @happy
  Scenario: Per-file encryption allows selective file encryption
    Given an open NovusPack package
    And a valid context
    And files with different encryption requirements
    When per-file encryption is used
    Then users can choose which files to encrypt
    And encryption selection is per file during package creation
    And packages can contain both encrypted and unencrypted files
    And unencrypted files accessed without decryption overhead

  @REQ-SEC-026 @happy
  Scenario: Per-file encryption operations support file-level encryption
    Given an open writable NovusPack package
    And a valid context
    And file requiring encryption
    When per-file encryption operations are used
    Then AddFileWithEncryption adds files with specific encryption types
    And GetFileEncryptionType retrieves encryption type for files
    And GetEncryptedFiles lists all encrypted files
    And file-level encryption provides granular control

  @REQ-SEC-037 @happy
  Scenario: Per-file security metadata provides file-level security information
    Given an open NovusPack package
    And a valid context
    And file with security metadata
    When per-file security metadata is examined
    Then security classification provides file security levels
    And access control provides file access restrictions
    And encryption selection is file-specific
    And compression selection is file-specific
    And security flags are file-specific

  @REQ-SEC-024 @happy
  Scenario: Per-file encryption provides performance optimization
    Given an open NovusPack package
    And a valid context
    And package with selective encryption
    When per-file encryption is used
    Then unencrypted files accessed without decryption overhead
    And performance optimization improves package operations
    And security granularity enables fine-grained control

  @REQ-SEC-024 @happy
  Scenario: Per-file encryption provides security granularity
    Given an open NovusPack package
    And a valid context
    And files with different security needs
    When per-file encryption is used
    Then fine-grained control over content protection is provided
    And sensitive files can be encrypted
    And non-sensitive files can remain unencrypted
    And security matches content sensitivity
