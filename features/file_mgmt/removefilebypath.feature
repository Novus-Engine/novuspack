@domain:file_mgmt @m2 @REQ-FILEMGMT-137 @spec(api_file_mgmt_removal.md#23-removefile-parameters)
Feature: RemoveFile by path

  @REQ-FILEMGMT-137 @happy
  Scenario: RemoveFile removes files by path
    Given an open writable package
    And file exists at path "documents/data.txt"
    When RemoveFile is called with path
    Then file is located by path
    And file is removed from package
    And removal operation completes successfully

  @REQ-FILEMGMT-137 @happy
  Scenario: RemoveFile locates file entry by virtual path
    Given an open writable package
    And file exists at virtual path "documents/file.txt"
    When RemoveFile is called with path
    Then file entry is located by path
    And file entry is removed
    And removal succeeds

  @REQ-FILEMGMT-137 @happy
  Scenario: RemoveFile removes file using path string
    Given an open writable package
    And file exists with primary path "data.bin"
    When RemoveFile is called with path
    Then file with that path is removed
    And file entry is removed from index
    And package metadata is updated

  @REQ-FILEMGMT-137 @happy
  Scenario: RemoveFile handles nested directory paths
    Given an open writable package
    And file exists at path "nested/deep/structure/file.txt"
    When RemoveFile is called with path
    Then file is located correctly
    And file is removed successfully
    And path resolution works correctly

  @REQ-FILEMGMT-137 @error
  Scenario: RemoveFile returns error when path not found
    Given an open writable package
    And file does not exist at specified path
    When RemoveFile is called with non-existent path
    Then ErrFileNotFound error is returned
    And error indicates file not found
    And error follows structured error format

  @REQ-FILEMGMT-137 @error
  Scenario: RemoveFile returns error for invalid path format
    Given an open writable package
    When RemoveFile is called with invalid path format
    Then ErrInvalidPath error is returned
    And error indicates invalid path
    And error follows structured error format

  @REQ-FILEMGMT-137 @error
  Scenario: RemoveFile respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
