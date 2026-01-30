@domain:file_mgmt @REQ-FILEMGMT-448 @spec(api_file_mgmt_updates.md#113-updatefile-returns) @spec(api_file_mgmt_updates.md#11-updatefile-package-method)
Feature: UpdateFile Returns

  @REQ-FILEMGMT-448 @happy
  Scenario: UpdateFile returns updated FileEntry and nil error on success
    Given an open writable package
    And an existing file in the package
    And a filesystem source file path
    When UpdateFile completes successfully
    Then an updated FileEntry is returned
    And error is nil

  @REQ-FILEMGMT-448 @error
  Scenario: UpdateFile returns non-nil error on failure
    Given an open writable package
    And an invalid storedPath or invalid sourceFilePath
    When UpdateFile is called
    Then error is returned

