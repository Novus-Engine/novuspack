@domain:security @m2 @REQ-SEC-082 @spec(api_security.md#311-purpose)
Feature: EncryptionType Purpose and Usage

  @REQ-SEC-082 @happy
  Scenario: Encryption type definition purpose defines available encryption algorithms
    Given an open NovusPack package
    And a valid context
    And encryption type system
    When encryption type definition purpose is examined
    Then purpose defines available encryption algorithms for file encryption
    And system provides encryption type enumeration
    And system enables encryption type selection

  @REQ-SEC-082 @happy
  Scenario: Encryption type definition supports file encryption configuration
    Given an open NovusPack package
    And a valid context
    And file encryption requirements
    When encryption type definition is used
    Then encryption algorithms are available for file encryption
    And encryption types can be selected per file
    And encryption configuration is type-safe
