@domain:core @m1 @REQ-CORE-013 @spec(api_core.md#9-package-metadata-management)
Feature: Package metadata management

  @happy
  Scenario: SetMetadata sets package metadata
    Given an open writable NovusPack package
    When SetMetadata is called with metadata map
    Then metadata is stored in package
    And metadata is accessible

  @happy
  Scenario: GetMetadata retrieves package metadata
    Given an open NovusPack package with metadata
    When GetMetadata is called
    Then metadata map is returned
    And all metadata is included

  @happy
  Scenario: UpdateMetadata updates package metadata
    Given an open writable NovusPack package with metadata
    When UpdateMetadata is called with updates
    Then existing metadata is updated
    And new metadata is added
    And unchanged metadata remains

  @happy
  Scenario: ValidateMetadata validates metadata structure
    Given package metadata
    When ValidateMetadata is called
    Then metadata structure is validated
    And validation errors are reported if invalid

  @happy
  Scenario: HasMetadata checks if package has metadata
    Given an open NovusPack package
    When HasMetadata is called
    Then true is returned if metadata exists
    And false is returned if no metadata

  @happy
  Scenario: SetAppID sets application identifier
    Given an open writable NovusPack package
    When SetAppID is called with app ID
    Then AppID is set in package header
    And AppID is accessible via GetInfo

  @happy
  Scenario: GetAppID gets current application identifier
    Given an open NovusPack package with AppID
    When GetAppID is called
    Then AppID value is returned
    And value matches header

  @happy
  Scenario: ClearAppID clears application identifier
    Given an open writable NovusPack package with AppID
    When ClearAppID is called
    Then AppID is set to 0
    And AppID is cleared

  @happy
  Scenario: HasAppID checks if application identifier is set
    Given an open NovusPack package
    When HasAppID is called
    Then true is returned if AppID is non-zero
    And false is returned if AppID is zero

  @happy
  Scenario: SetVendorID sets vendor identifier
    Given an open writable NovusPack package
    When SetVendorID is called with vendor ID
    Then VendorID is set in package header
    And VendorID is accessible via GetInfo

  @happy
  Scenario: GetVendorID gets current vendor identifier
    Given an open NovusPack package with VendorID
    When GetVendorID is called
    Then VendorID value is returned
    And value matches header

  @happy
  Scenario: ClearVendorID clears vendor identifier
    Given an open writable NovusPack package with VendorID
    When ClearVendorID is called
    Then VendorID is set to 0
    And VendorID is cleared

  @happy
  Scenario: HasVendorID checks if vendor identifier is set
    Given an open NovusPack package
    When HasVendorID is called
    Then true is returned if VendorID is non-zero
    And false is returned if VendorID is zero

  @happy
  Scenario: PackageInfo structure contains comprehensive information
    Given an open NovusPack package
    When PackageInfo is retrieved
    Then file count is included
    And package sizes are included
    And identity information is included
    And metadata information is included
    And signature information is included
    And security status is included

  @error
  Scenario: Metadata operations fail if package is read-only
    Given a read-only open NovusPack package
    When SetMetadata is called
    Then a structured validation error is returned

  @error
  Scenario: Metadata operations respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When metadata operation is called
    Then a structured context error is returned
