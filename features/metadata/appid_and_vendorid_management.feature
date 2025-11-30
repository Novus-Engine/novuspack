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
  Scenario: AppID and VendorID combinations work correctly
    Given an open writable package
    When AppID and VendorID are set together
    Then both identifiers are stored
    And combination identifies platform application
