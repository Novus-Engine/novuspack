@domain:file_mgmt @REQ-FILEMGMT-457 @spec(api_file_mgmt_updates.md#1712-convertsymlinkstohardlinks-behavior) @spec(api_file_mgmt_updates.md#17-symlinkconvertoptions-struct)
Feature: ConvertSymlinksToHardLinks Behavior

  @REQ-FILEMGMT-457 @happy
  Scenario: ConvertSymlinksToHardLinks removes symlink and restores hard link path
    Given an open writable package
    And one or more SymlinkEntry values in the package
    When ConvertSymlinksToHardLinks is called
    Then SymlinkEntry is removed from the package
    And the symlink source path is added to the target FileEntry.Paths
    And symlink metadata is preserved
    And FileEntry.PathCount and MetadataVersion are updated

