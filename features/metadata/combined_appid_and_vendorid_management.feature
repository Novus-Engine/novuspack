@domain:metadata @m2 @REQ-META-061 @spec(api_metadata.md#4-combined-management)
Feature: Combined AppID and VendorID Management

  @REQ-META-061 @happy
  Scenario: Combined management provides combined AppID/VendorID operations
    Given a NovusPack package
    When combined management operations are used
    Then SetPackageIdentity sets both VendorID and AppID
    And GetPackageIdentity gets both VendorID and AppID
    And ClearPackageIdentity clears both VendorID and AppID
    And GetPackageInfo gets comprehensive package information

  @REQ-META-061 @happy
  Scenario: SetPackageIdentity sets both identifiers atomically
    Given a NovusPack package
    And a valid context
    And a VendorID value
    And an AppID value
    When SetPackageIdentity is called
    Then both VendorID and AppID are set together
    And operation completes successfully
    And context supports cancellation

  @REQ-META-061 @happy
  Scenario: GetPackageIdentity retrieves both identifiers
    Given a NovusPack package
    And a package with VendorID and AppID set
    When GetPackageIdentity is called
    Then VendorID value is returned
    And AppID value is returned
    And both values are returned together

  @REQ-META-061 @happy
  Scenario: ClearPackageIdentity clears both identifiers
    Given a NovusPack package
    And a valid context
    And a package with VendorID and AppID set
    When ClearPackageIdentity is called
    Then both VendorID and AppID are cleared
    And VendorID is set to 0
    And AppID is set to 0
    And context supports cancellation

  @REQ-META-061 @error
  Scenario: Combined management validates parameters
    Given a NovusPack package
    When invalid parameters are provided
    Then parameter validation detects invalid values
    And appropriate error is returned
