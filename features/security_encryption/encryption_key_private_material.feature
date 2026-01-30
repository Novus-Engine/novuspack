@domain:security_encryption @m2 @REQ-CRYPTO-004 @spec(api_security.md#412-secure-encryption-key-operations-with-runtimesecret)
Feature: All EncryptionKey keys must be treated as private key material

  @REQ-CRYPTO-004 @happy
  Scenario: EncryptionKey keys are treated as private key material
    Given an EncryptionKey or key access operation
    When keys are accessed or used
    Then keys are treated as private key material
    And secure handling follows the private-key material policy
    And the behavior matches the runtime secret protection specification
