@domain:file_mgmt @m2 @REQ-FILEMGMT-150 @REQ-FILEMGMT-013 @spec(api_file_mgmt_updates.md#123-updatefilepattern-returns) @spec(api_file_mgmt_file_entry.md#613-updatefile-returns) @spec(api_file_mgmt_file_entry.md#614-updatefile-behavior) @spec(api_file_mgmt_file_entry.md#615-updatefile-error-conditions)
Feature: UpdateFilePattern

  @REQ-FILEMGMT-150 @REQ-FILEMGMT-013 @happy
  Scenario: Updatefilepattern returns slice of updated fileentry objects
    Given an open writable package
    And multiple files matching pattern "*.txt" exist
    When UpdateFilePattern is called with pattern, sourceDir, and options
    Then slice of updated FileEntry objects is returned
    And each entry contains all metadata, compression status, encryption details, and checksums
    And update operation completes successfully

  @REQ-FILEMGMT-150 @REQ-FILEMGMT-013 @happy
  Scenario: UpdateFilePattern updates files matching pattern
    Given an open writable package
    And files "doc1.txt", "doc2.txt", "data.bin" exist
    And corresponding files exist in sourceDir
    When UpdateFilePattern is called with pattern "*.txt"
    Then "doc1.txt" is updated
    And "doc2.txt" is updated
    And "data.bin" remains unchanged
    And pattern matching works correctly

  @REQ-FILEMGMT-150 @REQ-FILEMGMT-151 @happy
  Scenario: UpdateFilePattern performs pattern matching and bulk updates
    Given an open writable package
    And files matching pattern exist
    And sourceDir contains updated files
    When UpdateFilePattern is called
    Then pattern matching identifies files
    And each matching file is updated with new content from filesystem
    And bulk update operation completes

  @REQ-FILEMGMT-150 @REQ-FILEMGMT-013 @happy
  Scenario: UpdateFilePattern uses sourceDir for file location
    Given an open writable package
    And files matching pattern exist
    When UpdateFilePattern is called with sourceDir "updates/"
    Then files are read from sourceDir
    And updates are applied from sourceDir files
    And sourceDir parameter is used correctly

  @REQ-FILEMGMT-150 @REQ-FILEMGMT-013 @happy
  Scenario: UpdateFilePattern applies AddFileOptions to updates
    Given an open writable package
    And files matching pattern exist
    And AddFileOptions with compression and encryption settings
    When UpdateFilePattern is called with options
    Then options are applied to each update
    And compression settings are applied
    And encryption settings are applied

  @REQ-FILEMGMT-150 @REQ-FILEMGMT-152 @error
  Scenario: UpdateFilePattern returns error for invalid pattern
    Given an open writable package
    When UpdateFilePattern is called with invalid pattern
    Then ErrInvalidPattern error is returned
    And error indicates invalid pattern
    And error follows structured error format

  @REQ-FILEMGMT-150 @error
  Scenario: UpdateFilePattern handles partial updates on error
    Given an open writable package
    And some files fail to update
    When UpdateFilePattern encounters error
    Then some files may have been updated successfully
    And error indicates partial completion
    And updated entries are returned in result

  @REQ-FILEMGMT-150 @error
  Scenario: UpdateFilePattern respects context cancellation
    Given an open writable package
    And a cancelled context
    When UpdateFilePattern is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
