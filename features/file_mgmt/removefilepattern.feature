@domain:file_mgmt @m2 @REQ-FILEMGMT-325 @REQ-FILEMGMT-327 @spec(api_file_mgmt_removal.md#42-removefilepattern) @spec(api_file_mgmt_removal.md#partial-failure-handling)
Feature: RemoveFilePattern

  @REQ-FILEMGMT-327 @happy
  Scenario: RemoveFilePattern returns count of removed files
    Given an open writable package
    And multiple files matching pattern "*.txt" exist
    When RemoveFilePattern is called with pattern "*.txt"
    Then slice of removed FileEntry objects is returned
    And count matches number of removed files
    And removal operation completes successfully

  @REQ-FILEMGMT-328 @happy
  Scenario: RemoveFilePattern unstages files matching pattern
    Given an open writable package
    And files "doc1.txt", "doc2.txt", "data.bin" exist
    When RemoveFilePattern is called with pattern "*.txt"
    Then "doc1.txt" is removed
    And "doc2.txt" is removed
    And "data.bin" remains
    And pattern matching works correctly

  @REQ-FILEMGMT-328 @happy
  Scenario: RemoveFilePattern scans package for matching files
    Given an open writable package
    And files exist in various directories
    When RemoveFilePattern is called with pattern "documents/**/*.pdf"
    Then package is scanned for matching files
    And files matching pattern are identified
    And matching files are removed

  @REQ-FILEMGMT-328 @happy
  Scenario: RemoveFilePattern marks file data as deleted
    Given an open writable package
    And matching files with data exist
    When RemoveFilePattern removes files
    Then file data is marked as deleted
    And space is reclaimed during defragmentation
    And package metadata is updated

  @REQ-FILEMGMT-328 @happy
  Scenario: RemoveFilePattern updates package metadata and file count
    Given an open writable package
    And 10 files matching pattern exist
    When RemoveFilePattern removes all matching files
    Then package metadata is updated
    And file count is decremented by 10
    And package integrity is preserved

  @REQ-FILEMGMT-328 @happy
  Scenario: RemoveFilePattern preserves package integrity and signatures
    Given a signed writable package
    And matching files exist
    When RemoveFilePattern removes files
    Then package integrity is preserved
    And signatures remain valid
    And package structure is maintained

  @REQ-FILEMGMT-329 @error
  Scenario: RemoveFilePattern returns error for invalid pattern
    Given an open writable package
    When RemoveFilePattern is called with invalid pattern
    Then ErrInvalidPattern error is returned
    And error indicates invalid pattern
    And error follows structured error format

  @REQ-FILEMGMT-329 @error
  Scenario: RemoveFilePattern returns error when no files match
    Given an open writable package
    And no files match the pattern
    When RemoveFilePattern is called
    Then ErrNoFilesFound error is returned
    And error indicates no matches
    And error follows structured error format

  @REQ-FILEMGMT-327 @REQ-FILEMGMT-329 @error
  Scenario: RemoveFilePattern handles partial removal on error
    Given an open writable package
    And some files fail to remove
    When RemoveFilePattern encounters error
    Then some files may have been removed successfully
    And error indicates partial completion
    And removed entries are returned in result

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: RemoveFilePattern respects context cancellation
    Given an open writable package
    And a cancelled context
    When RemoveFilePattern is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
