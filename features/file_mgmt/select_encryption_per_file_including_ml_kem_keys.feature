@domain:file_mgmt @m2 @security @spec(api_file_mgmt_file_entry.md#9-fileentry-encryption)
Feature: Select encryption per file including ML-KEM keys

  @REQ-FILEMGMT-004 @happy
  Scenario: Set encryption settings using ML-KEM key management rules
    Given an open package
    And available ML-KEM keys
    When a file is added with encryption enabled
    Then the file metadata includes the selected encryption type
    And key references are stored in metadata
    And ML-KEM key management rules are followed

  @REQ-FILEMGMT-004 @happy
  Scenario: Select ML-KEM encryption per file
    Given an open package
    And ML-KEM key at security level 3
    When a file is added with ML-KEM encryption
    Then file is encrypted using ML-KEM
    And encryption type is set to EncryptionMLKEM
    And key reference is stored in file metadata

  @REQ-FILEMGMT-004 @happy
  Scenario: Select AES-256-GCM encryption per file
    Given an open package
    And AES encryption keys are available
    When a file is added with AES-256-GCM encryption
    Then file is encrypted using AES-256-GCM
    And encryption type is set to EncryptionAES256GCM
    And key reference is stored in file metadata

  @REQ-FILEMGMT-004 @happy
  Scenario: Select no encryption per file
    Given an open package
    When a file is added without encryption
    Then file is not encrypted
    And encryption type is set to EncryptionNone
    And no key references are stored

  @REQ-FILEMGMT-004 @happy
  Scenario: Per-file encryption selection supports mixed encryption types
    Given an open package
    When multiple files are added with different encryption types
    Then ML-KEM encrypted files can coexist with AES-256-GCM encrypted files
    And unencrypted files can coexist with encrypted files
    And each file maintains its own encryption settings

  @REQ-FILEMGMT-004 @error
  Scenario: Encryption selection requires available keys
    Given an open package
    And no encryption keys are available
    When a file is added with encryption enabled
    Then structured encryption error is returned
    And error indicates missing keys
    And error follows structured error format

  @REQ-FILEMGMT-004 @error
  Scenario: Encryption selection respects context cancellation
    Given an open package
    And a cancelled context
    When encryption selection operation is called
    Then structured context error is returned
    And error follows structured error format
