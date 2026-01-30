@domain:security_encryption @m2 @REQ-CRYPTO-007 @spec(api_security.md#412-secure-encryption-key-operations-with-runtimesecret)
Feature: Encrypt and Decrypt methods must execute within runtime/secret.Do when accessing keys

  @REQ-CRYPTO-007 @happy
  Scenario: Encrypt and Decrypt access keys within runtime/secret.Do
    Given an EncryptionKey and Encrypt or Decrypt operation
    When Encrypt or Decrypt is called
    Then key access must occur within runtime/secret.Do
    And key material is not exposed outside the Do callback
    And the behavior matches the secure key operations specification
