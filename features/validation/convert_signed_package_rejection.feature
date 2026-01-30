@domain:validation @REQ-FILEMGMT-352 @spec(api_file_mgmt_updates.md#17-convertpathstosymlinks-package-method)
Feature: Reject conversion on signed packages

  As a package system
  I want to prevent modification of signed packages
  So that package integrity and signatures remain valid

  @REQ-FILEMGMT-352 @error
  Scenario: Reject conversion on signed package
    Given a signed package
    And the package has a file entry with multiple paths
    When I attempt to convert the duplicate paths to symlinks
    Then the conversion should fail with ErrTypePackageState
    And the error message should indicate the package is signed
    And the package structure should remain unchanged

  @REQ-FILEMGMT-352 @error
  Scenario: Reject batch conversion on signed package
    Given a signed package
    And the package has multiple file entries with multiple paths
    When I attempt to call ConvertAllPathsToSymlinks
    Then the conversion should fail with ErrTypePackageState
    And the error message should indicate the package is signed
    And no file entries should be modified
