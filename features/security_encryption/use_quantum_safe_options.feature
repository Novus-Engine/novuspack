@domain:security_encryption @security @m4 @REQ-CRYPTO-002 @spec(security.md#31-quantum-safe-encryption)
Feature: Use Quantum-Safe Options

  @REQ-CRYPTO-002 @security @happy
  Scenario: Quantum-safe options provide quantum-safe signature and key management
    Given a NovusPack package
    When quantum-safe options are used
    Then quantum-safe signatures are available (ML-DSA, SLH-DSA)
    And quantum-safe encryption is available (ML-KEM)
    And quantum-safe key management is available

  @REQ-CRYPTO-002 @security @happy
  Scenario: Quantum-safe signatures use NIST PQC standards
    Given a NovusPack package
    When quantum-safe signatures are configured
    Then ML-DSA (CRYSTALS-Dilithium) signature support is available
    And SLH-DSA (SPHINCS+) signature support is available
    And signatures use NIST PQC standard algorithms

  @REQ-CRYPTO-002 @security @happy
  Scenario: Quantum-safe encryption uses ML-KEM
    Given a NovusPack package
    When quantum-safe encryption is configured
    Then ML-KEM (CRYSTALS-Kyber) encryption support is available
    And encryption uses NIST PQC standard algorithm
    And quantum-safe key exchange is provided

  @REQ-CRYPTO-002 @security @happy
  Scenario: Quantum-safe key management provides secure key generation
    Given a NovusPack package
    When quantum-safe key management is used
    Then ML-KEM keys can be generated at specified security levels
    And keys can be encrypted and decrypted using ML-KEM
    And keys support secure storage

  @REQ-CRYPTO-002 @security @error
  Scenario: Quantum-safe options validate security levels
    Given a NovusPack package
    When invalid security levels are provided
    Then security level validation detects invalid levels
    And appropriate errors are returned
