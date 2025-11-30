@domain:security @m2 @REQ-SEC-091 @spec(api_security.md#45-package-file-encryption-operations)
Feature: Package File Encryption Operations

  @REQ-SEC-091 @happy
  Scenario: Package file encryption operations provide EncryptFile method
    Given an open writable NovusPack package
    And a valid context
    And file path and data
    And encryption handler
    And encryption key
    When EncryptFile is called
    Then file is encrypted using handler and key
    And encryption uses security API patterns
    And encrypted file is stored in package
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-091 @happy
  Scenario: Package file encryption operations provide DecryptFile method
    Given an open NovusPack package
    And a valid context
    And encrypted file path
    And encryption handler
    And encryption key
    When DecryptFile is called
    Then file is decrypted using handler and key
    And decryption uses security API patterns
    And decrypted data is returned
    And original data is restored

  @REQ-SEC-091 @happy
  Scenario: Package file encryption operations provide ValidateFileEncryption method
    Given an open NovusPack package
    And a valid context
    And encrypted file path
    And encryption handler
    When ValidateFileEncryption is called
    Then file encryption is validated
    And validation uses security API patterns
    And encryption integrity is verified
    And validation confirms encryption is correct

  @REQ-SEC-091 @happy
  Scenario: Package file encryption operations provide GetFileEncryptionInfo method
    Given an open NovusPack package
    And a valid context
    And file path
    When GetFileEncryptionInfo is called
    Then EncryptionConfig is returned
    And encryption information uses security API patterns
    And encryption configuration is provided
    And encryption type and settings are available

  @REQ-SEC-091 @happy
  Scenario: Package file encryption operations support type-safe encryption
    Given an open NovusPack package
    And a valid context
    And generic encryption handler
    When package file encryption operations are used
    Then operations are type-safe
    And encryption key is type-safe
    And encrypted data is type-safe
    And encryption config is type-safe

  @REQ-SEC-091 @error
  Scenario: Package file encryption operations handle encryption failures
    Given an open NovusPack package
    And a valid context
    And invalid encryption key or handler
    When package file encryption operation fails
    Then structured error is returned
    And error indicates specific encryption failure
    And error follows structured error format
