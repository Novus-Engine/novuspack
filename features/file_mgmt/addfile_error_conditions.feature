@domain:file_mgmt @m2 @REQ-FILEMGMT-054 @spec(api_file_mgmt_addition.md#2110-addfile-error-conditions)
Feature: AddFile error conditions cover package state and I/O errors

  @REQ-FILEMGMT-054 @happy
  Scenario: AddFile returns structured errors for package state and I/O failures
    Given a package and filesystem path for AddFile
    When AddFile is called and a failure occurs
    Then error conditions cover package state and I/O errors
    And returned errors are structured
    And the behavior matches the AddFile error conditions specification
