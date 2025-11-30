@domain:core @m2 @REQ-CORE-043 @spec(api_core.md#92-appidvendorid-management)
Feature: AppID/VendorID Management

  @REQ-CORE-043 @happy
  Scenario: SetAppID sets the application identifier
    Given a NovusPack package
    When SetAppID is called with an application identifier
    Then application identifier is set
    And AppID is stored in package metadata
    And AppID can be retrieved with GetAppID

  @REQ-CORE-043 @happy
  Scenario: GetAppID retrieves the current application identifier
    Given a NovusPack package with AppID set
    When GetAppID is called
    Then current application identifier is retrieved
    And AppID value is returned
    And AppID information is available

  @REQ-CORE-043 @happy
  Scenario: ClearAppID clears the application identifier
    Given a NovusPack package with AppID set
    When ClearAppID is called
    Then application identifier is cleared
    And AppID is set to 0
    And HasAppID returns false

  @REQ-CORE-043 @happy
  Scenario: HasAppID checks if application identifier is set
    Given a NovusPack package
    When HasAppID is called
    Then check returns true if AppID is set (non-zero)
    And check returns false if AppID is not set (zero)
    And AppID status is determined

  @REQ-CORE-043 @happy
  Scenario: SetVendorID sets the vendor/platform identifier
    Given a NovusPack package
    When SetVendorID is called with a vendor identifier
    Then vendor/platform identifier is set
    And VendorID is stored in package metadata
    And VendorID can be retrieved with GetVendorID

  @REQ-CORE-043 @happy
  Scenario: GetVendorID retrieves the current vendor identifier
    Given a NovusPack package with VendorID set
    When GetVendorID is called
    Then current vendor identifier is retrieved
    And VendorID value is returned
    And VendorID information is available

  @REQ-CORE-043 @happy
  Scenario: ClearVendorID clears the vendor identifier
    Given a NovusPack package with VendorID set
    When ClearVendorID is called
    Then vendor identifier is cleared
    And VendorID is set to 0
    And HasVendorID returns false

  @REQ-CORE-043 @happy
  Scenario: HasVendorID checks if vendor identifier is set
    Given a NovusPack package
    When HasVendorID is called
    Then check returns true if VendorID is set (non-zero)
    And check returns false if VendorID is not set (zero)
    And VendorID status is determined
