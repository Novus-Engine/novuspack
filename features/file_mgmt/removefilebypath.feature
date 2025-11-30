@domain:file_mgmt @m2 @REQ-FILEMGMT-079 @spec(api_file_management.md#3122-removefilebypath)
Feature: RemoveFileByPath

  @REQ-FILEMGMT-079 @happy
  Scenario: Removefilebypath removes files by path
    Given an open writable package
    And file exists at path "documents/data.txt"
    When RemoveFileByPath is called with the path
    Then file is located by path
    And file is removed from package
    And removal operation completes successfully

  @REQ-FILEMGMT-079 @happy
  Scenario: RemoveFileByPath locates file entry by virtual path
    Given an open writable package
    And file exists at virtual path "documents/file.txt"
    When RemoveFileByPath is called with the path
    Then file entry is located by path
    And file entry is removed
    And removal succeeds

  @REQ-FILEMGMT-079 @happy
  Scenario: RemoveFileByPath removes file using path string
    Given an open writable package
    And file exists with primary path "data.bin"
    When RemoveFileByPath is called with "data.bin"
    Then file with that path is removed
    And file entry is removed from index
    And package metadata is updated

  @REQ-FILEMGMT-079 @happy
  Scenario: RemoveFileByPath handles nested directory paths
    Given an open writable package
    And file exists at path "nested/deep/structure/file.txt"
    When RemoveFileByPath is called with the nested path
    Then file is located correctly
    And file is removed successfully
    And path resolution works correctly

  @REQ-FILEMGMT-079 @error
  Scenario: RemoveFileByPath returns error when path not found
    Given an open writable package
    And file does not exist at specified path
    When RemoveFileByPath is called with non-existent path
    Then ErrFileNotFound error is returned
    And error indicates file not found
    And error follows structured error format

  @REQ-FILEMGMT-079 @error
  Scenario: RemoveFileByPath returns error for invalid path format
    Given an open writable package
    When RemoveFileByPath is called with invalid path format
    Then ErrInvalidPath error is returned
    And error indicates invalid path
    And error follows structured error format

  @REQ-FILEMGMT-079 @error
  Scenario: RemoveFileByPath respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFileByPath is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
