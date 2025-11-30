@domain:file_mgmt @m2 @REQ-FILEMGMT-011 @spec(api_file_management.md#42-remove-file-pattern)
Feature: Remove file pattern operations

  @happy
  Scenario: RemoveFilePattern removes multiple files matching pattern
    Given an open writable NovusPack package with multiple files
    When RemoveFilePattern is called with pattern
    Then all matching files are removed
    And file count decreases by number of matches
    And removed files are no longer accessible

  @happy
  Scenario: RemoveFilePattern returns results for each file
    Given RemoveFilePattern operation
    When pattern matching completes
    Then results indicate success or failure per file
    And removed file paths are included
    And errors are reported per file

  @happy
  Scenario: RemoveFilePattern supports wildcard patterns
    Given files with various names
    When RemoveFilePattern is called with wildcard pattern
    Then matching files are removed
    And non-matching files remain

  @error
  Scenario: RemoveFilePattern fails if pattern is invalid
    Given an invalid file pattern
    When RemoveFilePattern is called
    Then structured validation error is returned
    And error indicates invalid pattern

  @error
  Scenario: RemoveFilePattern respects context cancellation
    Given a large file set
    And a cancelled context
    When RemoveFilePattern is called
    Then structured context error is returned
    And operation is cancelled
