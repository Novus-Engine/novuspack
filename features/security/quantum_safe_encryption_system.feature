@skip @domain:security @m2 @spec(security.md#3-encryption-system)
Feature: Quantum-Safe Encryption System

# This feature captures high-level expectations for the encryption system from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: Encryption system supports quantum-safe and traditional encryption
    Given a package contains sensitive assets
    When encryption is configured for the package
    Then ML-KEM is available for quantum-safe key exchange
    And AES-256-GCM is available for compatibility and performance

  @documentation
  Scenario: Per-file encryption selection supports mixed packages
    Given a package contains files with different sensitivity levels
    When the caller selects encryption per file
    Then encrypted and unencrypted files can coexist in the same package
    And unencrypted files are accessed without decryption overhead
