@domain:security @m2 @REQ-SEC-042 @spec(security.md#522-vendor-and-application-identification)
Feature: Vendor and Application Identification

  @REQ-SEC-042 @happy
  Scenario: Vendor and application identification provides VendorID
    Given an open NovusPack package
    When vendor identification is examined
    Then VendorID field contains storefront/platform identifier
    And VendorID identifies trusted sources
    And VendorID provides package identification

  @REQ-SEC-042 @happy
  Scenario: Vendor and application identification provides AppID
    Given an open NovusPack package
    When application identification is examined
    Then AppID field contains application/game identifier
    And AppID provides package association
    And AppID enables package identification

  @REQ-SEC-042 @happy
  Scenario: Vendor and application identification provides LocaleID
    Given an open NovusPack package
    When locale identification is examined
    Then LocaleID field contains locale identifier
    And LocaleID provides path encoding information
    And LocaleID supports localization

  @REQ-SEC-042 @happy
  Scenario: Vendor and application identification provides CreatorID
    Given an open NovusPack package
    When creator identification is examined
    Then CreatorID field contains creator identifier
    And CreatorID provides package attribution
    And CreatorID enables creator identification

  @REQ-SEC-042 @happy
  Scenario: Vendor and application identification enables trusted source verification
    Given an open NovusPack package
    And package with VendorID and AppID
    When trusted source verification is performed
    Then VendorID enables storefront verification
    And AppID enables application verification
    And identification supports trust mechanisms
