@domain:security @m2 @REQ-SEC-080 @spec(api_security.md#24-securitylevel-enum)
Feature: SecurityLevel Enum

  @REQ-SEC-080 @happy
  Scenario: SecurityLevel enum provides SecurityLevelNone value
    Given an open NovusPack package
    And a valid context
    And package with no security features
    When SecurityLevel is examined
    Then SecurityLevelNone indicates no security
    And SecurityLevelNone is used for unsecured packages
    And SecurityLevelNone provides baseline classification

  @REQ-SEC-080 @happy
  Scenario: SecurityLevel enum provides SecurityLevelLow value
    Given an open NovusPack package
    And a valid context
    And package with basic security features
    When SecurityLevel is examined
    Then SecurityLevelLow indicates low security
    And SecurityLevelLow is used for basic security packages
    And SecurityLevelLow provides low-level classification

  @REQ-SEC-080 @happy
  Scenario: SecurityLevel enum provides SecurityLevelMedium value
    Given an open NovusPack package
    And a valid context
    And package with moderate security features
    When SecurityLevel is examined
    Then SecurityLevelMedium indicates medium security
    And SecurityLevelMedium is used for moderate security packages
    And SecurityLevelMedium provides medium-level classification

  @REQ-SEC-080 @happy
  Scenario: SecurityLevel enum provides SecurityLevelHigh value
    Given an open NovusPack package
    And a valid context
    And package with high security features
    When SecurityLevel is examined
    Then SecurityLevelHigh indicates high security
    And SecurityLevelHigh is used for high security packages
    And SecurityLevelHigh provides high-level classification

  @REQ-SEC-080 @happy
  Scenario: SecurityLevel enum provides SecurityLevelMaximum value
    Given an open NovusPack package
    And a valid context
    And package with maximum security features
    When SecurityLevel is examined
    Then SecurityLevelMaximum indicates maximum security
    And SecurityLevelMaximum is used for maximum security packages
    And SecurityLevelMaximum provides maximum-level classification

  @REQ-SEC-080 @happy
  Scenario: SecurityLevel enum is used in SecurityStatus structure
    Given an open NovusPack package
    And a valid context
    And package with security validation
    When GetSecurityStatus is called
    Then SecurityStatus.SecurityLevel contains SecurityLevel value
    And SecurityLevel reflects overall package security
    And SecurityLevel classification is consistent
