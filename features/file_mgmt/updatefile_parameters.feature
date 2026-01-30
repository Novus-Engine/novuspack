@domain:file_mgmt @REQ-FILEMGMT-447 @spec(api_file_mgmt_updates.md#112-updatefile-parameters) @spec(api_file_mgmt_updates.md#11-updatefile-package-method)
Feature: UpdateFile Parameters

  @REQ-FILEMGMT-447 @happy
  Scenario: UpdateFile accepts storedPath sourceFilePath and options
    Given an open writable package
    And a stored package path identifying an existing file
    And a filesystem source file path with new content
    When UpdateFile is called with ctx, storedPath, sourceFilePath, and options
    Then ctx is used for cancellation and timeouts
    And storedPath identifies the package file to update
    And sourceFilePath identifies the filesystem source content
    And options may be nil for defaults

