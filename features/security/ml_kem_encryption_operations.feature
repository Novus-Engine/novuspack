@domain:security @m2 @REQ-SEC-101 @spec(api_security.md#53-ml-kem-encryption-operations)
Feature: ML-KEM Encryption Operations

  @REQ-SEC-101 @happy
  Scenario: MLKEMKey Encrypt performs encryption using ML-KEM key
    Given an open NovusPack package
    And ML-KEM key
    And plaintext data
    And a valid context
    When Encrypt is called with plaintext
    Then encryption is performed using ML-KEM key
    And encrypted data is returned
    And encryption uses ML-KEM algorithm

  @REQ-SEC-101 @happy
  Scenario: MLKEMKey Decrypt performs decryption using ML-KEM key
    Given an open NovusPack package
    And ML-KEM key
    And ciphertext data
    And a valid context
    When Decrypt is called with ciphertext
    Then decryption is performed using ML-KEM key
    And decrypted data is returned
    And decryption uses ML-KEM algorithm

  @REQ-SEC-101 @happy
  Scenario: ML-KEM encryption operations accept context parameter
    Given an open NovusPack package
    And ML-KEM key
    And a valid context
    When Encrypt or Decrypt is called
    Then context parameter is accepted
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-101 @happy
  Scenario: ML-KEM encryption operations accept plaintext parameter
    Given an open NovusPack package
    And ML-KEM key
    And plaintext data
    And a valid context
    When Encrypt is called
    Then plaintext parameter is accepted
    And plaintext is encrypted
    And encrypted data is returned

  @REQ-SEC-101 @happy
  Scenario: ML-KEM encryption operations accept ciphertext parameter
    Given an open NovusPack package
    And ML-KEM key
    And ciphertext data
    And a valid context
    When Decrypt is called
    Then ciphertext parameter is accepted
    And ciphertext is decrypted
    And decrypted data is returned

  @REQ-SEC-105 @REQ-SEC-011 @error
  Scenario: ML-KEM encryption fails if encryption fails
    Given an open NovusPack package
    And ML-KEM key
    And system conditions preventing encryption
    And a valid context
    When Encrypt is called
    Then ErrEncryptionFailed error is returned
    And error indicates encryption failure

  @REQ-SEC-105 @REQ-SEC-011 @error
  Scenario: ML-KEM decryption fails if decryption fails
    Given an open NovusPack package
    And ML-KEM key
    And system conditions preventing decryption
    And a valid context
    When Decrypt is called
    Then ErrDecryptionFailed error is returned
    And error indicates decryption failure

  @REQ-SEC-105 @error
  Scenario: ML-KEM encryption operations fail with invalid key
    Given an open NovusPack package
    And invalid or corrupted ML-KEM key
    And a valid context
    When Encrypt or Decrypt is called
    Then ErrInvalidKey error is returned
    And error indicates key is invalid or corrupted

  @REQ-SEC-105 @REQ-SEC-011 @error
  Scenario: ML-KEM encryption operations fail with cancelled context
    Given an open NovusPack package
    And ML-KEM key
    And a cancelled context
    When Encrypt or Decrypt is called
    Then ErrContextCancelled error is returned
    And error indicates context cancellation

  @REQ-SEC-105 @REQ-SEC-011 @error
  Scenario: ML-KEM encryption operations fail with context timeout
    Given an open NovusPack package
    And ML-KEM key
    And a context with timeout
    And timeout expires during encryption operation
    When Encrypt or Decrypt is called
    Then ErrContextTimeout error is returned
    And error indicates context timeout
