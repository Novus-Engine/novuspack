@domain:file_mgmt @m2 @REQ-FILEMGMT-336 @spec(api_file_mgmt_removal.md#47-removedirectory-error-conditions)
Feature: RemoveDirectory error conditions handle invalid directory paths and package state errors

  @REQ-FILEMGMT-336 @happy
  Scenario: RemoveDirectory error conditions handle invalid paths and state
    Given an open NovusPack package
    When RemoveDirectory is called with invalid path or bad state
    Then error conditions handle invalid directory paths and package state errors
    And the behavior matches the RemoveDirectory error-conditions specification
    And structured error is returned for invalid path
    And structured error is returned for package not open
