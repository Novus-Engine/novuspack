@domain:core @m2 @REQ-CORE-040 @spec(api_core.md#7-digital-signatures-and-security)
Feature: Digital signatures and security

  @REQ-CORE-040 @happy
  Scenario: Digital signatures provide package integrity verification
    Given a NovusPack package
    When digital signatures are used
    Then package integrity verification is provided
    And signatures verify package authenticity
    And signatures ensure package content integrity

  @REQ-CORE-040 @happy
  Scenario: Multiple signature types are supported including quantum-safe
    Given signature operations
    When signatures are added
    Then multiple signature types are supported
    And ML-DSA (quantum-safe) signatures are supported
    And SLH-DSA (quantum-safe) signatures are supported
    And PGP signatures are supported
    And X.509/PKCS#7 signatures are supported

  @REQ-CORE-040 @happy
  Scenario: Multiple signatures can be added to a single package
    Given a NovusPack package
    When multiple signatures are added
    Then multiple signatures can be added to single package
    And incremental signing supports sequential signatures
    And each signature validates all previous content

  @REQ-CORE-040 @happy
  Scenario: Signature validation provides detailed status per signature
    Given a signed NovusPack package
    When signatures are validated
    Then validation returns detailed status per signature
    And each signature validation result is available
    And validation status indicates success or failure
    And validation details include signature information

  @REQ-CORE-040 @happy
  Scenario: Signed packages are protected from modification
    Given a signed NovusPack package
    When write operations are attempted
    Then signed packages are protected from modification
    And write operations are blocked after signing
    And package immutability is enforced
