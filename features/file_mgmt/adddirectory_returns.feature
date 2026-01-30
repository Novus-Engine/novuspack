@domain:file_mgmt @REQ-FILEMGMT-444 @REQ-FILEMGMT-313 @spec(api_file_mgmt_addition.md#adddirectory-returns) @spec(api_file_mgmt_addition.md#25-adddirectory-package-method)
Feature: AddDirectory Returns

  AddDirectory returns a slice of created FileEntry objects for all
  successfully added files, and an error. Partial success is possible:
  some files may have been added when an error occurs.

  @REQ-FILEMGMT-444 @happy
  Scenario: AddDirectory returns slice of created FileEntry on full success
    Given an open writable package
    And a filesystem directory with files
    When AddDirectory is called with dirPath and options
    Then slice of FileEntry objects is returned
    And each element corresponds to a successfully added file
    And error is nil
    And slice length equals number of discovered files

  @REQ-FILEMGMT-444 @happy
  Scenario: Partial success returns aggregated results and error
    Given an open writable package
    And a directory where some file additions succeed and some fail
    When AddDirectory is called and some failures occur
    Then slice contains FileEntry objects for successfully added files only
    And error is non-nil describing the failure
    And aggregated results are returned even when some files failed

  @REQ-FILEMGMT-444 @error
  Scenario: AddDirectory returns error when directory invalid or package not open
    Given a package not open or invalid dirPath
    When AddDirectory is called
    Then error is returned
    And returned slice may be nil or empty per implementation
