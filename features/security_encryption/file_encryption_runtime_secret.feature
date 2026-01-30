@domain:security_encryption @m2 @REQ-CRYPTO-008 @spec(api_security.md#461-secure-file-encryption-operations-with-runtimesecret)
Feature: File encryption operations must execute within runtime/secret.Do

  @REQ-CRYPTO-008 @happy
  Scenario: File encryption operations use runtime/secret.Do
    Given a file encryption or decryption operation
    When the operation accesses encryption keys
    Then the operation must execute within runtime/secret.Do
    And key material is protected during the operation
    And the behavior matches the secure file encryption specification
