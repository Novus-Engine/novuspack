@domain:file_mgmt @m2 @spec(api_file_management.md#41-remove-file)
Feature: Remove file entry

  @REQ-FILEMGMT-002 @happy
  Scenario: Removing a file updates directory state
    Given a package containing a file "old.bin"
    When I remove the file "old.bin"
    Then the directory state should mark "old.bin" as removed

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: RemoveFile validates path parameter
    Given an open writable package
    When RemoveFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: RemoveFile respects context cancellation
    Given an open writable package with files
    And a cancelled context
    When RemoveFile is called
    Then structured context error is returned
    And error type is context cancellation
