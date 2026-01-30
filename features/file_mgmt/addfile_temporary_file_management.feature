@domain:file_mgmt @REQ-FILEMGMT-437 @spec(api_file_mgmt_addition.md#temporary-file-management) @spec(api_file_mgmt_addition.md#21-packageaddfile-method)
Feature: AddFile Temporary File Management

  AddFile does not delete temporary files before returning; they remain
  open until Write completes. AddFile cleans up temp files on failure.
  Package close cleans up all tracked temp files.

  @REQ-FILEMGMT-437 @happy
  Scenario: AddFile does not delete temp files before returning
    Given an open writable package
    And AddFileOptions that require a temp file (e.g. encryption)
    And a filesystem file
    When AddFile completes successfully
    Then temporary files are not deleted before return
    And temp files remain open and accessible
    And Write operations can read from temp files

  @REQ-FILEMGMT-437 @happy
  Scenario: Temp files are available until Write completes
    Given an open writable package
    And a file added with encryption (using temp file)
    When Write has not yet been called
    Then FileEntry.CurrentSource references temp file
    And temp file is open and readable
    When Write completes
    Then temp file can be closed and removed per spec

  @REQ-FILEMGMT-437 @error
  Scenario: AddFile cleans up temp file on failure
    Given an open writable package
    And AddFileOptions that create a temp file
    And a scenario that causes AddFile to fail after temp file creation
    When AddFile returns an error
    Then temporary file created during AddFile is cleaned up
    And no temp file handle is leaked
    And error is returned to caller

  @REQ-FILEMGMT-437 @happy
  Scenario: Package close cleans up tracked temp files
    Given an open writable package
    And one or more FileEntry objects with temp files (CurrentSource.IsTempFile true)
    When Package is closed before Write completes
    Then Package close operation cleans up all tracked temporary files
    And no temp files are left on disk from FileEntry staging
