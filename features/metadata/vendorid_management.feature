@domain:metadata @m2 @REQ-META-060 @spec(api_metadata.md#3-vendorid-management)
Feature: VendorID Management

  @REQ-META-060 @happy
  Scenario: SetVendorID sets vendor identifier
    Given an open writable NovusPack package
    And a valid context
    When SetVendorID is called with vendor ID
    Then VendorID is set in package header
    And VendorID is accessible via GetVendorID
    And VendorID is accessible via GetInfo

  @REQ-META-060 @happy
  Scenario: GetVendorID retrieves current vendor identifier
    Given an open NovusPack package with VendorID set
    When GetVendorID is called
    Then VendorID value is returned
    And value matches header
    And value is a 32-bit unsigned integer

  @REQ-META-060 @happy
  Scenario: ClearVendorID clears vendor identifier
    Given an open writable NovusPack package with VendorID set
    And a valid context
    When ClearVendorID is called
    Then VendorID is set to 0
    And VendorID is cleared
    And HasVendorID returns false

  @REQ-META-060 @happy
  Scenario: HasVendorID checks if vendor identifier is set
    Given an open NovusPack package
    When HasVendorID is called
    Then true is returned if VendorID is non-zero
    And false is returned if VendorID is zero
    And VendorID status is determined

  @REQ-META-060 @happy
  Scenario: GetVendorIDInfo gets detailed vendor identifier information
    Given an open NovusPack package with VendorID set
    When GetVendorIDInfo is called
    Then VendorIDInfo structure is returned
    And detailed vendor information is available

  @REQ-META-013 @REQ-META-011 @error
  Scenario: SetVendorID validates vendor ID parameter
    Given an open writable NovusPack package
    And a valid context
    When SetVendorID is called with invalid vendor ID
    Then structured validation error is returned
    And error indicates validation failure

  @REQ-META-011 @REQ-META-014 @error
  Scenario: VendorID operations respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When SetVendorID or ClearVendorID is called
    Then structured context error is returned
    And error type is context cancellation

  @REQ-META-011 @error
  Scenario: VendorID operations fail if package is read-only
    Given a read-only open NovusPack package
    When SetVendorID is called
    Then structured validation error is returned
    And error indicates read-only package
