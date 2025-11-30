@domain:file_mgmt @m2 @REQ-FILEMGMT-169 @REQ-FILEMGMT-016 @spec(api_file_management.md#553-returns)
Feature: RemoveFilePath

  @REQ-FILEMGMT-169 @REQ-FILEMGMT-016 @happy
  Scenario: Removefilepath returns updated fileentry
    Given an open writable package
    And file entry with multiple paths exists
    When RemoveFilePath is called with entry and path string
    Then updated FileEntry is returned
    And specified path is removed from entry
    And PathCount is decremented
    And other paths remain unchanged

  @REQ-FILEMGMT-169 @REQ-FILEMGMT-016 @happy
  Scenario: RemoveFilePath removes path from file entry while preserving content
    Given an open writable package
    And file entry has paths "original.txt" and "copy.txt"
    When RemoveFilePath is called with "copy.txt"
    Then "copy.txt" path is removed
    And "original.txt" path remains
    And file content is preserved
    And file entry remains valid

  @REQ-FILEMGMT-169 @REQ-FILEMGMT-016 @happy
  Scenario: RemoveFilePath updates PathCount after removal
    Given an open writable package
    And file entry has PathCount of 3
    When RemoveFilePath removes one path
    Then PathCount is decremented to 2
    And updated FileEntry reflects new count
    And path array is updated

  @REQ-FILEMGMT-169 @REQ-FILEMGMT-016 @happy
  Scenario: RemoveFilePath preserves remaining paths
    Given an open writable package
    And file entry has paths "path1.txt", "path2.txt", "path3.txt"
    When RemoveFilePath is called with "path2.txt"
    Then "path1.txt" remains
    And "path3.txt" remains
    And removed path is no longer present
    And file entry is updated correctly

  @REQ-FILEMGMT-169 @REQ-FILEMGMT-016 @happy
  Scenario: RemoveFilePath increments MetadataVersion
    Given an open writable package
    And file entry with MetadataVersion 5
    When RemoveFilePath is called
    Then MetadataVersion is incremented to 6
    And metadata change is tracked
    And version reflects path removal

  @REQ-FILEMGMT-169 @REQ-FILEMGMT-171 @error
  Scenario: RemoveFilePath returns error when path not found
    Given an open writable package
    And file entry exists
    And specified path does not exist in entry
    When RemoveFilePath is called with non-existent path
    Then structured error is returned
    And error indicates path not found
    And error follows structured error format

  @REQ-FILEMGMT-169 @error
  Scenario: RemoveFilePath respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFilePath is called
    Then structured context error is returned
    And error follows structured error format
