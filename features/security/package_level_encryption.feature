@domain:security @m2 @REQ-SEC-040 @spec(security.md#52-package-level-security)
Feature: Package-Level Encryption

  @REQ-SEC-040 @happy
  Scenario: Package-level security provides security flags
    Given an open NovusPack package
    And a valid context
    And package with security flags configured
    When package-level security flags are examined
    Then Bit 7 indicates multiple signatures enabled
    And Bit 6 indicates quantum-safe signatures present
    And Bit 5 indicates traditional signatures present
    And Bit 4 indicates timestamps present
    And Bit 3 indicates metadata present
    And Bit 2 indicates chain validation enabled
    And Bit 1 indicates revocation support
    And Bit 0 indicates expiration support

  @REQ-SEC-040 @happy
  Scenario: Package-level security provides vendor identification
    Given an open NovusPack package
    And a valid context
    And package with vendor information
    When package-level security is examined
    Then VendorID contains storefront/platform identifier
    And VendorID enables trusted source identification
    And VendorID supports security verification

  @REQ-SEC-040 @happy
  Scenario: Package-level security provides application identification
    Given an open NovusPack package
    And a valid context
    And package with application information
    When package-level security is examined
    Then AppID contains application/game identifier
    And AppID enables package association
    And AppID supports application-specific security policies

  @REQ-SEC-040 @happy
  Scenario: Package-level security provides locale identification
    Given an open NovusPack package
    And a valid context
    And package with locale information
    When package-level security is examined
    Then LocaleID contains locale identifier
    And LocaleID enables path encoding configuration
    And LocaleID supports locale-specific security settings

  @REQ-SEC-040 @happy
  Scenario: Package-level security provides creator identification
    Given an open NovusPack package
    And a valid context
    And package with creator information
    When package-level security is examined
    Then CreatorID contains creator identifier
    And CreatorID enables package attribution
    And CreatorID supports creator-based security policies

  @REQ-SEC-040 @happy
  Scenario: Package-level security provides comprehensive security configuration
    Given an open NovusPack package
    And a valid context
    And package with full security configuration
    When package-level security is examined
    Then security flags provide signature feature configuration
    And vendor/application identification provides package identification
    And locale/creator identification provides attribution
    And package-wide security settings are consistent
