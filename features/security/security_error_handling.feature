@domain:security @m2 @REQ-SEC-099 @REQ-SEC-105 @spec(api_security.md#524-error-conditions)
Feature: Security Error Handling

  @REQ-SEC-099 @error
  Scenario: ML-KEM key generation fails with invalid security level
    Given an open NovusPack package
    And invalid security level (not 1-5)
    And a valid context
    When GenerateMLKEMKey is called with invalid level
    Then ErrInvalidSecurityLevel error is returned
    And error indicates security level must be 1-5
    And error follows structured error format

  @REQ-SEC-099 @error
  Scenario: ML-KEM key generation fails if key generation fails
    Given an open NovusPack package
    And system conditions preventing key generation
    And a valid context
    When GenerateMLKEMKey is called
    Then ErrKeyGenerationFailed error is returned
    And error indicates key generation failure
    And error follows structured error format

  @REQ-SEC-099 @REQ-SEC-011 @error
  Scenario: ML-KEM key generation fails with cancelled context
    Given an open NovusPack package
    And a cancelled context
    When GenerateMLKEMKey is called
    Then ErrContextCancelled error is returned
    And error indicates context cancellation
    And error follows structured error format

  @REQ-SEC-099 @REQ-SEC-011 @error
  Scenario: ML-KEM key generation fails with context timeout
    Given an open NovusPack package
    And a context with timeout
    And timeout expires during key generation
    When GenerateMLKEMKey is called
    Then ErrContextTimeout error is returned
    And error indicates context timeout
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption fails if encryption fails
    Given an open NovusPack package
    And ML-KEM key
    And system conditions preventing encryption
    And a valid context
    When Encrypt is called
    Then ErrEncryptionFailed error is returned
    And error indicates encryption failure
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM decryption fails if decryption fails
    Given an open NovusPack package
    And ML-KEM key
    And ciphertext
    And system conditions preventing decryption
    And a valid context
    When Decrypt is called
    Then ErrDecryptionFailed error is returned
    And error indicates decryption failure
    And error follows structured error format

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption fails with invalid key
    Given an open NovusPack package
    And invalid or corrupted ML-KEM key
    And a valid context
    When Encrypt or Decrypt is called
    Then ErrInvalidKey error is returned
    And error indicates key is invalid or corrupted
    And error follows structured error format

  @REQ-SEC-105 @REQ-SEC-011 @error
  Scenario: ML-KEM encryption operations fail with cancelled context
    Given an open NovusPack package
    And ML-KEM key
    And a cancelled context
    When Encrypt or Decrypt is called
    Then ErrContextCancelled error is returned
    And error indicates context cancellation
    And error follows structured error format

  @REQ-SEC-105 @REQ-SEC-011 @error
  Scenario: ML-KEM encryption operations fail with context timeout
    Given an open NovusPack package
    And ML-KEM key
    And a context with timeout
    And timeout expires during encryption operation
    When Encrypt or Decrypt is called
    Then ErrContextTimeout error is returned
    And error indicates context timeout
    And error follows structured error format
