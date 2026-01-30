@domain:file_mgmt @m2 @REQ-FILEMGMT-051 @REQ-FILEMGMT-053 @spec(api_file_mgmt_addition.md#21-addfile-package-method)
Feature: Filesystem path validation for AddFile

  @happy
  Scenario: AddFile reads from a valid filesystem path
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then file content is read from filesystem path
    And file is added to package

  @error
  Scenario: AddFile returns structured validation error for invalid filesystem path
    Given an open writable package
    And an invalid filesystem file path
    When AddFile is called
    Then structured validation error is returned
    And error follows structured error format
