@domain:file_mgmt @REQ-FILEMGMT-454 @spec(api_file_mgmt_updates.md#176-convertpathstosymlinks-returns) @spec(api_file_mgmt_updates.md#17-symlinkconvertoptions-struct)
Feature: ConvertPathsToSymlinks Returns

  @REQ-FILEMGMT-454 @happy
  Scenario: ConvertPathsToSymlinks returns updated FileEntry and created SymlinkEntry slice
    Given an open writable package
    And a FileEntry with multiple paths
    When ConvertPathsToSymlinks completes successfully
    Then an updated FileEntry with PathCount = 1 is returned
    And a slice of created SymlinkEntry values is returned
    And error is nil

