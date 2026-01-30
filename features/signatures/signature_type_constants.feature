@domain:signatures @m2 @v2 @REQ-SIG-025 @REQ-SIG-044 @spec(api_signatures.md#21-signature-type-constants)
Feature: Signature Type Constants

  @REQ-SIG-025 @happy
  Scenario: Signature type constants define supported signature types
    Given a NovusPack package
    When signature type constants are examined
    Then SignatureTypeNone is 0x00
    And SignatureTypeMLDSA is 0x01 for ML-DSA (CRYSTALS-Dilithium)
    And SignatureTypeSLHDSA is 0x02 for SLH-DSA (SPHINCS+)
    And SignatureTypePGP is 0x03 for PGP (OpenPGP)
    And SignatureTypeX509 is 0x04 for X.509/PKCS#7
    And values 0x05-0xFF are reserved for future signature types

  @REQ-SIG-025 @happy
  Scenario: Signature type constants support quantum-safe algorithms
    Given a NovusPack package
    When signature type constants are examined
    Then ML-DSA signature type is defined
    And SLH-DSA signature type is defined
    And both quantum-safe signature types are available

  @REQ-SIG-025 @error
  Scenario: Invalid signature type returns error
    Given a NovusPack package
    And an invalid signature type value
    When signature operation is attempted with invalid type
    Then ErrUnsupportedSignatureType error is returned
    And error indicates unsupported signature algorithm
    And error follows structured error format

  @REQ-SIG-044 @happy
  Scenario: Signature size comparison compares signature sizes
    Given a NovusPack package
    When signature sizes are compared
    Then ML-DSA signatures are approximately 2,420-4,595 bytes
    And SLH-DSA signatures are approximately 7,856-17,088 bytes
    And PGP signatures are approximately 100-1,000 bytes
    And X.509 signatures are approximately 200-2,000 bytes

  @REQ-SIG-044 @happy
  Scenario: Quantum-safe signatures are larger than traditional signatures
    Given a NovusPack package
    When signature sizes are compared
    Then quantum-safe signatures (ML-DSA, SLH-DSA) are larger than traditional signatures
    And ML-DSA signatures are smaller than SLH-DSA signatures
    And size difference reflects quantum-safe security requirements
