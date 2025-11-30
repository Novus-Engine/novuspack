@domain:file_format @m2 @REQ-FILEFMT-026 @spec(package_file_format.md#221-packagedataversion-field)
Feature: PackageDataVersion Field

  @REQ-FILEFMT-026 @happy
  Scenario: PackageDataVersion field tracks package data version
    Given a NovusPack package
    And PackageDataVersion field is present
    When PackageDataVersion field is examined
    Then PackageDataVersion is a 32-bit unsigned integer
    And PackageDataVersion tracks changes to package data content
    And PackageDataVersion enables package data change detection

  @REQ-FILEFMT-026 @happy
  Scenario: PackageDataVersion increments when files are added
    Given a NovusPack package
    And PackageDataVersion has initial value
    When a file is added to the package
    Then PackageDataVersion is incremented
    And version change indicates package data modification
    And version tracks file addition operations

  @REQ-FILEFMT-026 @happy
  Scenario: PackageDataVersion increments when files are removed
    Given a NovusPack package
    And package contains files
    And PackageDataVersion has current value
    When a file is removed from the package
    Then PackageDataVersion is incremented
    And version change indicates package data modification
    And version tracks file removal operations

  @REQ-FILEFMT-026 @happy
  Scenario: PackageDataVersion increments when file data is modified
    Given a NovusPack package
    And package contains files
    And PackageDataVersion has current value
    When file data is modified
    Then PackageDataVersion is incremented
    And version change indicates package data modification
    And version tracks file data changes

  @REQ-FILEFMT-026 @happy
  Scenario: PackageDataVersion initial value is 1 for new packages
    Given a new NovusPack package
    When PackageDataVersion is examined
    Then PackageDataVersion is set to 1
    And initial version indicates new package
    And version starts at 1 for new packages

  @REQ-FILEFMT-026 @happy
  Scenario: PackageDataVersion supports version range from 1 to 4 billion
    Given a NovusPack package
    When PackageDataVersion is examined
    Then PackageDataVersion range is 1 to 4294967295
    And version supports 4 billion versions
    And version enables long-term change tracking
