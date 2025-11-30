@domain:signatures @m2 @REQ-SIG-041 @REQ-SIG-043 @spec(api_signatures.md#31-industry-standard-comparison)
Feature: Industry Standard Comparison

  @REQ-SIG-041 @happy
  Scenario: Industry standard comparison compares with industry standards
    Given a NovusPack package
    When industry standard comparison is examined
    Then comparison includes feature comparison table
    And comparison shows NovusPack vs PGP Files
    And comparison shows NovusPack vs X.509/PKCS#7
    And comparison shows NovusPack vs Windows Authenticode
    And comparison shows NovusPack vs macOS Code Signing

  @REQ-SIG-041 @happy
  Scenario: Industry standard comparison shows feature differences
    Given a NovusPack package
    When feature comparison is examined
    Then signature location, header metadata, multiple signatures are compared
    And signature types, quantum-safe, cross-platform features are compared
    And key management, signature size, verification speed are compared
    And industry adoption status is compared

  @REQ-SIG-043 @happy
  Scenario: Industry standard compliance ensures standards alignment
    Given a NovusPack package
    When industry standard compliance is examined
    Then PGP Compatibility follows OpenPGP standard (RFC 4880)
    And X.509 Compliance follows PKCS#7 standard (RFC 2315)
    And Signature Placement follows industry standard (end of file)
    And Hash Algorithm uses SHA-256 (industry standard)
    And Key Management supports standard key formats

  @REQ-SIG-043 @happy
  Scenario: Industry standard compliance validates implementation
    Given a NovusPack package
    When compliance validation is performed
    Then OpenPGP standard compliance is verified
    And PKCS#7 standard compliance is verified
    And industry standard practices are followed
    And compliance ensures interoperability

  @REQ-SIG-041 @REQ-SIG-043 @error
  Scenario: Industry standard comparison and compliance validation
    Given a NovusPack package
    When comparison or compliance validation fails
    Then appropriate errors are returned
    And errors follow structured error format
