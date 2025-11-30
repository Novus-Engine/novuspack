@domain:signatures @m2 @REQ-SIG-045 @spec(api_signatures.md#35-verification-performance)
Feature: Verification Performance

  @REQ-SIG-045 @happy
  Scenario: NovusPack signature verification is fast
    Given a NovusPack package
    When signature verification is performed
    Then verification is fast with optimized hash calculation
    And verification performance is acceptable for package operations
    And verification enables efficient package processing

  @REQ-SIG-045 @happy
  Scenario: PGP signature verification performance
    Given a NovusPack package
    When PGP signature verification is performed
    Then verification is fast with established algorithms
    And verification performance matches industry standards
    And verification enables efficient package processing

  @REQ-SIG-045 @happy
  Scenario: X.509 signature verification performance
    Given a NovusPack package
    When X.509 signature verification is performed
    Then verification is fast with certificate chain validation
    And verification includes chain validation overhead
    And verification performance is acceptable for package operations

  @REQ-SIG-045 @happy
  Scenario: Quantum-safe signature verification performance
    Given a NovusPack package
    When quantum-safe signature verification is performed
    Then ML-DSA verification is fast
    And SLH-DSA verification is fast
    And verification performance enables practical use

  @REQ-SIG-045 @error
  Scenario: Verification performance handles errors efficiently
    Given a NovusPack package
    When signature verification fails
    Then error is returned quickly
    And error does not impact verification performance unnecessarily
    And error follows structured error format
