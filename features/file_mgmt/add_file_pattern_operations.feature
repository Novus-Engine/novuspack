@domain:file_mgmt @m2 @REQ-FILEMGMT-008 @spec(api_file_mgmt_addition.md#22-add-file-pattern)
Feature: Add file pattern operations

  @happy
  Scenario: AddFilePattern adds multiple files matching pattern
    Given an open writable NovusPack package
    When AddFilePattern is called with pattern and options
    Then all matching files are added
    And file count matches pattern matches
    And all files are added successfully

  @happy
  Scenario: AddFilePattern supports recursive directory traversal
    Given a directory structure with nested files
    When AddFilePattern is called with recursive option
    Then files in subdirectories are included
    And directory structure is preserved

  @happy
  Scenario: AddFilePattern respects include patterns
    Given files with various extensions
    When AddFilePattern is called with include pattern
    Then only matching files are added
    And non-matching files are excluded

  @happy
  Scenario: AddFilePattern respects exclude patterns
    Given files with various extensions
    When AddFilePattern is called with exclude pattern
    Then excluded files are not added
    And matching non-excluded files are added

  @happy
  Scenario: AddFilePattern handles symlinks according to option
    Given files including symlinks
    When AddFilePattern is called with FollowSymlinks option
    Then symlinks are followed if enabled
    And symlinks are not followed if disabled

  @happy
  Scenario: AddFilePattern returns results for each file
    Given AddFilePattern operation
    When pattern matching completes
    Then results indicate success or failure per file
    And file paths are included in results
    And errors are reported per file

  @error
  Scenario: AddFilePattern fails if pattern is invalid
    Given an invalid file pattern
    When AddFilePattern is called
    Then structured validation error is returned
    And error indicates invalid pattern

  @error
  Scenario: AddFilePattern respects context cancellation
    Given a large file set
    And a cancelled context
    When AddFilePattern is called
    Then structured context error is returned
    And operation is cancelled
