@domain:security_encryption @m2 @REQ-CRYPTO-006 @spec(api_security.md#412-secure-encryption-key-operations-with-runtimesecret)
Feature: SetKey must execute within runtime/secret.Do for all keys

  @REQ-CRYPTO-006 @happy
  Scenario: SetKey executes within runtime/secret.Do
    Given an EncryptionKey and key material to set
    When SetKey is called
    Then SetKey must execute within runtime/secret.Do
    And key material is only provided inside the Do callback
    And the behavior matches the secure key operations specification
