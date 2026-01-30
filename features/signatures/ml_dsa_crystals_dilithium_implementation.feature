@domain:signatures @m2 @v2 @REQ-SIG-029 @spec(api_signatures.md#23-ml-dsa-crystals-dilithium-implementation)
Feature: ML-DSA Crystals Dilithium Implementation

  @REQ-SIG-029 @happy
  Scenario: ML-DSA implementation provides quantum-safe signature support
    Given a NovusPack package
    And a valid context
    When ML-DSA implementation is used
    Then algorithm uses NIST PQC Standard ML-DSA
    And all three security levels are supported
    And performance is optimized for package signing
    And secure key generation and storage is provided
    And context supports cancellation

  @REQ-SIG-029 @happy
  Scenario: ML-DSA supports all three security levels
    Given a NovusPack package
    And a valid context
    When ML-DSA security levels are examined
    Then Level 2 provides ~2,420-byte signatures with 128-bit security
    And Level 3 provides ~3,293-byte signatures with 192-bit security
    And Level 5 provides ~4,595-byte signatures with 256-bit security

  @REQ-SIG-029 @happy
  Scenario: ML-DSA provides optimized performance
    Given a NovusPack package
    And a valid context
    When ML-DSA signing is performed
    Then signing performance is optimized for archive signing
    And verification performance is fast
    And algorithm is efficient for package format

  @REQ-SIG-029 @error
  Scenario: ML-DSA implementation handles errors
    Given a NovusPack package
    When ML-DSA operations fail
    Then appropriate errors are returned
    And errors follow structured error format
