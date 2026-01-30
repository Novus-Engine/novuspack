@domain:file_mgmt @m2 @REQ-FILEMGMT-331 @spec(api_file_mgmt_removal.md#41-removedirectory-purpose)
Feature: RemoveDirectory purpose is to remove all files within a directory path from package

  @REQ-FILEMGMT-331 @happy
  Scenario: RemoveDirectory purpose removes directory contents
    Given an open NovusPack package with file entries under a directory path
    When RemoveDirectory purpose is applied
    Then purpose is to remove all files within the directory path from package
    And the behavior matches the RemoveDirectory purpose specification
    And directory path is validated before removal
    And all files under the path are removed
