@domain:metadata @m2 @REQ-META-004 @spec(api_metadata.md#2-appid-management)
Feature: AppID and VendorID management

  @happy
  Scenario: SetAppID sets application identifier
    Given an open writable package
    When SetAppID is called with app ID
    Then AppID is set in package header
    And AppID is accessible via GetInfo

  @happy
  Scenario: SetVendorID sets vendor identifier
    Given an open writable package
    When SetVendorID is called with vendor ID
    Then VendorID is set in package header
    And VendorID is accessible via GetInfo

  @happy
  Scenario: GetAppID retrieves application identifier
    Given an open package with AppID set
    When GetAppID is called
    Then AppID value is returned
    And value matches the set AppID

  @happy
  Scenario: GetVendorID retrieves vendor identifier
    Given an open package with VendorID set
    When GetVendorID is called
    Then VendorID value is returned
    And value matches the set VendorID

  @happy
  Scenario: ClearAppID removes application identifier
    Given an open writable package with AppID set
    When ClearAppID is called
    Then AppID is removed from package
    And GetAppID returns zero value
    And HasAppID returns false

  @happy
  Scenario: ClearVendorID removes vendor identifier
    Given an open writable package with VendorID set
    When ClearVendorID is called
    Then VendorID is removed from package
    And GetVendorID returns zero value
    And HasVendorID returns false

  @happy
  Scenario: HasAppID checks application identifier existence
    Given an open package
    When HasAppID is called
    Then true is returned if AppID is set
    And false is returned if AppID is not set

  @happy
  Scenario: HasVendorID checks vendor identifier existence
    Given an open package
    When HasVendorID is called
    Then true is returned if VendorID is set
    And false is returned if VendorID is not set

  @happy
  Scenario: AppID and VendorID combinations work correctly
    Given an open writable package
    When AppID and VendorID are set together
    Then both identifiers are stored
    And combination identifies platform application
    And GetAppID and GetVendorID return correct values

  @error
  Scenario: SetAppID fails on read-only package
    Given a read-only open package
    When SetAppID is called
    Then structured security error is returned
    And error indicates package is read-only

  @error
  Scenario: SetVendorID fails on read-only package
    Given a read-only open package
    When SetVendorID is called
    Then structured security error is returned
    And error indicates package is read-only
