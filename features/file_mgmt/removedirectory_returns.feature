@domain:file_mgmt @m2 @REQ-FILEMGMT-334 @spec(api_file_mgmt_removal.md#45-removedirectory-returns)
Feature: RemoveDirectory returns slice of removed file paths and error

  @REQ-FILEMGMT-334 @happy
  Scenario: RemoveDirectory returns removed paths and error
    Given an open NovusPack package with file entries under a directory
    When RemoveDirectory completes
    Then returns slice of removed file paths and error as specified
    And the behavior matches the RemoveDirectory returns specification
    And slice contains paths of removed files
    And error indicates failure conditions
