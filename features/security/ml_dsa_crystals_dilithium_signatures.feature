@domain:security @m2 @v2 @REQ-SEC-017 @spec(security.md#231-ml-dsa-crystals-dilithium)
Feature: ML-DSA Crystals Dilithium Signatures

  @REQ-SEC-017 @happy
  Scenario: ML-DSA provides quantum-safe signature algorithm
    Given an open NovusPack package
    And a valid context
    And ML-DSA signature implementation
    When ML-DSA signature is examined
    Then algorithm uses NIST PQC Standard ML-DSA
    And algorithm provides quantum-safe signatures
    And algorithm supports all three security levels

  @REQ-SEC-017 @happy
  Scenario: ML-DSA supports security level 2 (128-bit security)
    Given an open NovusPack package
    And a valid context
    And ML-DSA with security level 2
    When ML-DSA signature is created
    Then signature size is approximately 2420 bytes
    And signature provides 128-bit security
    And signature follows ML-DSA specifications

  @REQ-SEC-017 @happy
  Scenario: ML-DSA supports security level 3 (192-bit security)
    Given an open NovusPack package
    And a valid context
    And ML-DSA with security level 3
    When ML-DSA signature is created
    Then signature size is approximately 3293 bytes
    And signature provides 192-bit security
    And signature follows ML-DSA specifications

  @REQ-SEC-017 @happy
  Scenario: ML-DSA supports security level 5 (256-bit security)
    Given an open NovusPack package
    And a valid context
    And ML-DSA with security level 5
    When ML-DSA signature is created
    Then signature size is approximately 4595 bytes
    And signature provides 256-bit security
    And signature follows ML-DSA specifications

  @REQ-SEC-017 @happy
  Scenario: ML-DSA provides optimized performance for package signing
    Given an open NovusPack package
    And a valid context
    And ML-DSA signature implementation
    When ML-DSA signing is performed
    Then signing performance is optimized for archive signing
    And signature key management is secure
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-017 @happy
  Scenario: ML-DSA provides secure key generation and storage
    Given an open NovusPack package
    And a valid context
    And ML-DSA key generation requirements
    When ML-KEM key is generated for ML-DSA
    Then key generation uses cryptographically secure random number generators
    And key storage implements secure storage mechanisms
    And key management follows security best practices
