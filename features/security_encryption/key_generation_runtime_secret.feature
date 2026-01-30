@domain:security_encryption @m2 @REQ-CRYPTO-009 @spec(api_security.md#526-secure-encryption-operations-with-runtimesecret)
Feature: Key generation operations must execute within runtime/secret.Do

  @REQ-CRYPTO-009 @happy
  Scenario: Key generation uses runtime/secret.Do
    Given a key generation operation
    When keys are generated
    Then the operation must execute within runtime/secret.Do
    And generated key material is protected
    And the behavior matches the secure encryption operations specification
