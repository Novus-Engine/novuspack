@domain:file_mgmt @m2 @REQ-FILEMGMT-335 @spec(api_file_mgmt_removal.md#46-removedirectory-behavior)
Feature: RemoveDirectory behavior includes directory validation, file discovery, removal, and metadata cleanup

  @REQ-FILEMGMT-335 @happy
  Scenario: RemoveDirectory behavior includes validation discovery removal cleanup
    Given an open NovusPack package with file entries under a directory
    When RemoveDirectory is performed
    Then behavior includes directory validation, file discovery, removal, and metadata cleanup
    And the behavior matches the RemoveDirectory behavior specification
    And directory path is validated first
    And metadata is cleaned up after removal
