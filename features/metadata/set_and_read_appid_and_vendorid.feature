@domain:metadata @m5 @spec(api_metadata.md#2-appid-management)
Feature: Set and read AppID and VendorID

  @REQ-META-002 @happy
  Scenario: AppID and VendorID constraints enforced
    Given an open package
    When I set AppID and VendorID to valid values
    Then reading them should return those values and constraints are enforced

  @REQ-META-011 @REQ-META-013 @error
  Scenario: SetAppID validates AppID parameter
    Given an open writable package
    When SetAppID is called with invalid format
    Then structured validation error is returned
    And error indicates invalid AppID format

  @REQ-META-011 @REQ-META-013 @error
  Scenario: SetVendorID validates VendorID parameter
    Given an open writable package
    When SetVendorID is called with invalid format
    Then structured validation error is returned
    And error indicates invalid VendorID format

  @REQ-META-011 @REQ-META-013 @error
  Scenario: AppID/VendorID operations validate length
    Given an open writable package
    When SetAppID or SetVendorID is called with excessive length
    Then structured validation error is returned
    And error indicates length limit exceeded

  @REQ-META-011 @REQ-META-014 @error
  Scenario: AppID/VendorID operations respect context cancellation
    Given an open writable package
    And a cancelled context
    When AppID or VendorID operation is called
    Then structured context error is returned
    And error type is context cancellation
