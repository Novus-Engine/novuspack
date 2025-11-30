@domain:file_mgmt @m2 @REQ-FILEMGMT-013 @spec(api_file_management.md#62-update-file-pattern)
Feature: Update file pattern operations

  @happy
  Scenario: UpdateFilePattern updates multiple files matching pattern
    Given an open writable NovusPack package with multiple files
    When UpdateFilePattern is called with pattern and options
    Then all matching files are updated
    And file versions are incremented
    And updates are applied correctly

  @happy
  Scenario: UpdateFilePattern returns results for each file
    Given UpdateFilePattern operation
    When pattern matching completes
    Then results indicate success or failure per file
    And updated file paths are included
    And errors are reported per file

  @error
  Scenario: UpdateFilePattern fails if pattern is invalid
    Given an invalid file pattern
    When UpdateFilePattern is called
    Then structured validation error is returned

  @error
  Scenario: UpdateFilePattern respects context cancellation
    Given a large file set
    And a cancelled context
    When UpdateFilePattern is called
    Then structured context error is returned
