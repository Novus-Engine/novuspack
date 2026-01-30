@domain:security @m2 @v2 @REQ-SEC-018 @spec(security.md#232-slh-dsa-sphincs)
Feature: SLH-DSA SPHINCS Signatures

  @REQ-SEC-018 @happy
  Scenario: SLH-DSA provides quantum-safe hash-based signatures
    Given an open NovusPack package
    And a valid context
    And SLH-DSA signature implementation
    When SLH-DSA signature is examined
    Then algorithm uses NIST PQC Standard SLH-DSA
    And algorithm provides quantum-safe hash-based signatures
    And algorithm supports all three security levels

  @REQ-SEC-018 @happy
  Scenario: SLH-DSA supports security level 1 (128-bit security)
    Given an open NovusPack package
    And a valid context
    And SLH-DSA with security level 1
    When SLH-DSA signature is created
    Then signature size is approximately 7856 bytes
    And signature provides 128-bit security
    And signature follows SLH-DSA specifications

  @REQ-SEC-018 @happy
  Scenario: SLH-DSA supports security level 3 (192-bit security)
    Given an open NovusPack package
    And a valid context
    And SLH-DSA with security level 3
    When SLH-DSA signature is created
    Then signature size is approximately 12272 bytes
    And signature provides 192-bit security
    And signature follows SLH-DSA specifications

  @REQ-SEC-018 @happy
  Scenario: SLH-DSA supports security level 5 (256-bit security)
    Given an open NovusPack package
    And a valid context
    And SLH-DSA with security level 5
    When SLH-DSA signature is created
    Then signature size is approximately 17088 bytes
    And signature provides 256-bit security
    And signature follows SLH-DSA specifications

  @REQ-SEC-018 @happy
  Scenario: SLH-DSA provides stateless hash-based signatures
    Given an open NovusPack package
    And a valid context
    And SLH-DSA signature implementation
    When SLH-DSA signing is performed
    Then signatures are stateless
    And signature verification is stateless
    And stateless design simplifies key management

  @REQ-SEC-018 @happy
  Scenario: SLH-DSA provides single-use key generation
    Given an open NovusPack package
    And a valid context
    And SLH-DSA key generation requirements
    When SLH-DSA key is generated
    Then key generation uses single-use approach
    And single-use keys provide additional security
    And key management follows SLH-DSA requirements
