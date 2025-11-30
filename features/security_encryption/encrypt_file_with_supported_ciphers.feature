@domain:security_encryption @security @m2 @REQ-CRYPTO-001 @spec(security.md#32-encryption-algorithms)
Feature: Encrypt file with supported ciphers

  @security
  Scenario: Supported ciphers and key sizes enforced
    Given a file to encrypt
    When I encrypt using a supported cipher and key size
    Then encryption should succeed and metadata should reflect the cipher and key size

  @security
  Scenario: AES-256-GCM encryption is supported
    Given a file to encrypt
    When AES-256-GCM encryption is applied
    Then encryption succeeds
    And EncryptionType is set correctly
    And file is protected

  @security
  Scenario: Quantum-safe encryption (ML-KEM + ML-DSA) is supported
    Given a file to encrypt
    When quantum-safe encryption is applied
    Then encryption succeeds
    And quantum-safe encryption type is set
    And file is protected against quantum attacks

  @security
  Scenario: Encryption key management follows security rules
    Given file encryption functionality
    When encryption keys are managed
    Then ML-KEM key management rules are followed
    And keys are generated securely
    And keys are stored securely

  @error
  Scenario: Invalid encryption parameters are rejected
    Given file encryption functionality
    When encryption is attempted with invalid parameters
    Then structured validation error is returned
    And error indicates invalid encryption configuration
