@domain:file_mgmt @REQ-FILEMGMT-455 @spec(api_file_mgmt_updates.md#177-convertpathstosymlinks-behavior) @spec(api_file_mgmt_updates.md#17-symlinkconvertoptions-struct)
Feature: ConvertPathsToSymlinks Behavior

  @REQ-FILEMGMT-455 @happy
  Scenario: ConvertPathsToSymlinks validates selects primary and creates symlinks
    Given an open writable package
    And a FileEntry with multiple paths
    When ConvertPathsToSymlinks is called
    Then pre-conversion validation is performed
    And a primary path is selected per options or default
    And non-primary paths are converted into SymlinkEntry values pointing to the primary path
    And FileEntry and path metadata are updated consistently

