@domain:security_encryption @m2 @REQ-CRYPTO-005 @spec(api_security.md#412-secure-encryption-key-operations-with-runtimesecret)
Feature: GetKey must execute within runtime/secret.Do for all keys

  @REQ-CRYPTO-005 @happy
  Scenario: GetKey executes within runtime/secret.Do
    Given an EncryptionKey and a secret access callback
    When GetKey is called
    Then GetKey must execute within runtime/secret.Do
    And key material is only accessible inside the Do callback
    And the behavior matches the secure key operations specification
