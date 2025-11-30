@domain:signatures @m2 @REQ-SIG-040 @spec(api_signatures.md#3-comparison-with-other-signed-file-implementations)
Feature: Signature Implementation Comparison

  @REQ-SIG-040 @happy
  Scenario: Comparison with other implementations compares signature systems
    Given a NovusPack package
    When comparison with other implementations is examined
    Then comparison includes PGP Files
    And comparison includes X.509/PKCS#7
    And comparison includes Windows Authenticode
    And comparison includes macOS Code Signing
    And feature comparison is provided

  @REQ-SIG-040 @happy
  Scenario: Comparison includes signature location and metadata features
    Given a NovusPack package
    When signature location comparison is examined
    Then NovusPack places signatures at end of file
    And comparison shows signature placement differences
    And header metadata comparison shows NovusPack extended header advantage

  @REQ-SIG-040 @happy
  Scenario: Comparison includes signature types and quantum-safe features
    Given a NovusPack package
    When signature types comparison is examined
    Then NovusPack supports 4 types (ML-DSA, SLH-DSA, PGP, X.509)
    And comparison shows quantum-safe advantage (ML-DSA/SLH-DSA)
    And comparison shows cross-platform advantage

  @REQ-SIG-040 @happy
  Scenario: Comparison includes signature size and performance
    Given a NovusPack package
    When signature size and performance comparison is examined
    Then signature sizes are compared across implementations
    And verification speed is compared
    And performance characteristics are documented

  @REQ-SIG-040 @error
  Scenario: Comparison validation ensures accuracy
    Given a NovusPack package
    When comparison data is validated
    Then comparison data is verified for accuracy
    And appropriate validation is performed
