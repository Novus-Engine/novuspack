@domain:security @m2 @REQ-SEC-044 @spec(security.md#61-comparison-with-industry-standards)
Feature: Security Comparison with Industry Standards

  @REQ-SEC-044 @happy
  Scenario: NovusPack supports multiple signatures like industry standards
    Given an open NovusPack package
    And comparison with industry standards
    When security features are compared
    Then multiple signatures are supported
    And feature matches PGP, X.509, Authenticode, and Code Signing
    And comparison demonstrates parity

  @REQ-SEC-044 @happy
  Scenario: NovusPack supports more signature types than industry standards
    Given an open NovusPack package
    And comparison with industry standards
    When signature types are compared
    Then NovusPack supports 4 signature types
    And other formats support 1 type each
    And NovusPack provides more flexibility

  @REQ-SEC-044 @happy
  Scenario: NovusPack provides quantum-safe signatures unlike most standards
    Given an open NovusPack package
    And comparison with industry standards
    When quantum-safe features are compared
    Then NovusPack supports quantum-safe signatures
    And PGP, X.509, Authenticode, and Code Signing do not
    And NovusPack provides future-proof security

  @REQ-SEC-044 @happy
  Scenario: NovusPack is cross-platform like PGP and X.509
    Given an open NovusPack package
    And comparison with industry standards
    When cross-platform support is compared
    Then NovusPack works on all major operating systems
    And feature matches PGP and X.509
    And Windows Authenticode and macOS Code Signing are platform-specific

  @REQ-SEC-044 @happy
  Scenario: NovusPack supports multiple key management approaches
    Given an open NovusPack package
    And comparison with industry standards
    When key management is compared
    Then NovusPack supports multiple key management approaches
    And flexibility exceeds platform-specific key stores
    And comparison demonstrates advantages
