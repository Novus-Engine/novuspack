@domain:security_encryption @m2 @REQ-CRYPTO-010 @spec(api_security.md#536-secure-key-clearing-with-runtimesecret)
Feature: Key cleanup operations must execute within runtime/secret.Do

  @REQ-CRYPTO-010 @happy
  Scenario: Key cleanup uses runtime/secret.Do
    Given an EncryptionKey and a cleanup or clear operation
    When key cleanup is performed
    Then the operation must execute within runtime/secret.Do
    And key material is securely cleared
    And the behavior matches the secure key clearing specification
