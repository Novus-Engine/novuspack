@domain:file_mgmt @m2 @REQ-FILEMGMT-009 @REQ-FILEMGMT-026 @spec(api_file_mgmt_addition.md#21-addfile-package-method)
Feature: AddFile reads data from filesystem paths

  @happy
  Scenario: AddFile uses filesystem path input
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then file content is read from filesystem path

  @happy
  Scenario: AddFile supports cancellation via context
    Given an open writable package
    And a filesystem file path
    And a cancelled context
    When AddFile is called
    Then ErrContextCancelled error is returned
    And error follows structured error format
