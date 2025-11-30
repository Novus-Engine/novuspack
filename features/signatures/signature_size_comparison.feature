@domain:signatures @m2 @REQ-SIG-044 @spec(api_signatures.md#34-signature-size-comparison)
Feature: Signature Size Comparison

  @REQ-SIG-044 @happy
  Scenario: ML-DSA signature sizes are documented
    Given a NovusPack package
    When ML-DSA signature sizes are examined
    Then ML-DSA signatures are approximately 2,420-4,595 bytes
    And size varies by security level
    And size is appropriate for quantum-safe signatures

  @REQ-SIG-044 @happy
  Scenario: SLH-DSA signature sizes are documented
    Given a NovusPack package
    When SLH-DSA signature sizes are examined
    Then SLH-DSA signatures are approximately 7,856-17,088 bytes
    And size varies by security level
    And size is larger than ML-DSA signatures
    And size is appropriate for quantum-safe hash-based signatures

  @REQ-SIG-044 @happy
  Scenario: Traditional signature sizes are documented
    Given a NovusPack package
    When traditional signature sizes are examined
    Then PGP signatures are approximately 100-1,000 bytes
    And X.509 signatures are approximately 200-2,000 bytes
    And traditional signatures are smaller than quantum-safe signatures

  @REQ-SIG-044 @happy
  Scenario: Signature size comparison enables informed selection
    Given a NovusPack package
    When signature type is selected
    Then signature size information enables selection
    And size comparison supports trade-off analysis
    And size information helps with storage planning
