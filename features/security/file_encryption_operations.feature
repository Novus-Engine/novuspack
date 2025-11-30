@domain:security @m2 @REQ-SEC-090 @spec(api_security.md#44-file-encryption-operations)
Feature: File Encryption Operations

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide FileEncryptionHandler interface
    Given an open NovusPack package
    And a valid context
    And FileEncryptionHandler implementation
    When file encryption operations are examined
    Then FileEncryptionHandler provides file-specific encryption operations
    And EncryptFile method encrypts file using handler and key
    And DecryptFile method decrypts file using handler and key
    And ValidateFileEncryption method validates file encryption

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide AES256GCMFileHandler
    Given an open NovusPack package
    And a valid context
    And AES256GCMFileHandler implementation
    When file encryption operations are used
    Then AES256GCMFileHandler provides AES-256-GCM encryption
    And handler encrypts files using AES-256-GCM algorithm
    And handler decrypts files using AES-256-GCM algorithm
    And handler validates AES-256-GCM encryption

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide MLKEMFileHandler
    Given an open NovusPack package
    And a valid context
    And MLKEMFileHandler implementation
    When file encryption operations are used
    Then MLKEMFileHandler provides ML-KEM encryption
    And handler encrypts files using ML-KEM algorithm
    And handler decrypts files using ML-KEM algorithm
    And handler validates ML-KEM encryption

  @REQ-SEC-090 @happy
  Scenario: File encryption operations support context integration
    Given an open NovusPack package
    And a valid context
    And file encryption handler
    When file encryption operations are performed
    Then all methods accept context.Context
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-090 @happy
  Scenario: File encryption operations provide type-safe encryption
    Given an open NovusPack package
    And a valid context
    And generic file encryption handler
    When file encryption operations are used
    Then operations are type-safe
    And encryption key is type-safe
    And encrypted data is type-safe

  @REQ-SEC-090 @error
  Scenario: File encryption operations handle encryption failures
    Given an open NovusPack package
    And a valid context
    And invalid encryption key or data
    When file encryption operation fails
    Then structured error is returned
    And error indicates specific encryption failure
    And error follows structured error format
