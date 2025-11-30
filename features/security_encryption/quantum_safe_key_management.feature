@domain:security_encryption @security @m2 @REQ-CRYPTO-003 @spec(security.md#31-quantum-safe-encryption)
Feature: Quantum-safe key management

  @security
  Scenario: Quantum-safe keys are generated correctly
    Given quantum-safe key generation
    When ML-KEM keys are generated
    Then keys are generated with appropriate security level
    And key format is correct
    And keys are ready for use

  @security
  Scenario: Multiple security levels are supported
    Given quantum-safe key generation
    When keys are generated at different security levels
    Then security levels 1-5 are supported
    And key size matches security level
    And security level is preserved
