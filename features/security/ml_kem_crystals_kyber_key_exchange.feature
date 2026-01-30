@skip @domain:security @m2 @spec(security.md#321-ml-kem-crystals-kyber)
Feature: ML-KEM Crystals Kyber Key Exchange

# This feature captures high-level expectations for ML-KEM from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: ML-KEM provides quantum-safe key exchange with multiple security levels
    Given ML-KEM encryption is enabled
    When a key pair is generated
    Then the selected security level determines key sizes and security strength
    And keys are suitable for quantum-safe key exchange in the encryption system

  @documentation
  Scenario: ML-KEM is used for key exchange in a hybrid encryption approach
    Given a package uses hybrid encryption
    When a file is encrypted
    Then ML-KEM is used for key exchange
    And AES-256-GCM may be used for bulk data encryption
