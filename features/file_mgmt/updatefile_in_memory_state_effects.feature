@domain:file_mgmt @REQ-FILEMGMT-450 @spec(api_file_mgmt_updates.md#117-in-memory-package-state-effects-updatefile) @spec(api_file_mgmt_updates.md#11-updatefile-package-method)
Feature: UpdateFile In-Memory Package State Effects

  @REQ-FILEMGMT-450 @happy
  Scenario: Updated FileEntry is visible to subsequent operations without disk write
    Given an open writable package
    And an existing file entry in the package
    When UpdateFile completes successfully
    Then the existing FileEntry is updated in the in-memory index
    And the updated FileEntry is visible to subsequent in-process operations
    And PackageInfo is updated to reflect the new state

