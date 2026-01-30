@domain:file_mgmt @m2 @REQ-FILEMGMT-333 @spec(api_file_mgmt_removal.md#44-removedirectoryoptions-struct)
Feature: RemoveDirectoryOptions structure configures directory removal behavior

  @REQ-FILEMGMT-333 @happy
  Scenario: RemoveDirectoryOptions configures removal behavior
    Given RemoveDirectoryOptions for directory removal
    When RemoveDirectory is called with options
    Then RemoveDirectoryOptions structure configures directory removal behavior
    And the behavior matches the RemoveDirectoryOptions struct specification
    And options control recursive and safety behavior
    And defaults are safe when options are nil
