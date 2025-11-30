@domain:security @m2 @REQ-SEC-041 @spec(security.md#521-security-flags)
Feature: Security Flags

  @REQ-SEC-041 @happy
  Scenario: Security flags define signature feature bits
    Given an open NovusPack package
    And a valid context
    And package with security flags configured
    When security flags are examined
    Then Bit 7 indicates multiple signatures enabled
    And Bit 6 indicates quantum-safe signatures present
    And Bit 5 indicates traditional signatures present
    And Bit 4 indicates timestamps present
    And Bit 3 indicates metadata present
    And Bit 2 indicates chain validation enabled
    And Bit 1 indicates revocation support
    And Bit 0 indicates expiration support

  @REQ-SEC-041 @happy
  Scenario: Security flags provide comprehensive security configuration
    Given an open NovusPack package
    And a valid context
    And package with full security configuration
    When security flags are examined
    Then flags configure signature features (bits 15-8)
    And flags configure package security options
    And flags provide package-wide security settings
    And flags enable security feature detection

  @REQ-SEC-041 @happy
  Scenario: Security flags enable multiple signatures
    Given an open NovusPack package
    And a valid context
    And package with multiple signatures
    When security flags are examined
    Then Bit 7 is set when multiple signatures are enabled
    And flag indicates incremental signing support
    And flag enables multiple signature validation

  @REQ-SEC-041 @happy
  Scenario: Security flags indicate signature types
    Given an open NovusPack package
    And a valid context
    And package with quantum-safe signatures
    When security flags are examined
    Then Bit 6 is set when quantum-safe signatures are present
    When package has traditional signatures
    Then Bit 5 is set when traditional signatures are present
    And flags enable signature type detection

  @REQ-SEC-041 @happy
  Scenario: Security flags indicate security features
    Given an open NovusPack package
    And a valid context
    And package with security features
    When security flags are examined
    Then Bit 4 indicates timestamps are present
    And Bit 3 indicates metadata is present
    And Bit 2 indicates chain validation is enabled
    And Bit 1 indicates revocation support
    And Bit 0 indicates expiration support
