@domain:file_mgmt @m2 @REQ-FILEMGMT-009 @REQ-FILEMGMT-026 @spec(api_file_mgmt_addition.md#21-addfile-package-method)
Feature: File data source derived from filesystem path

  @happy
  Scenario: AddFile derives file data source from filesystem path
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then file content is read from filesystem path

  @happy
  Scenario: AddFile reads file content using streaming when needed
    Given an open writable package
    And a filesystem file path
    When AddFile is called
    Then streaming is used for large files when needed
    And memory is managed efficiently

  @error
  Scenario: AddFile returns structured error when filesystem path is invalid
    Given an open writable package
    And an invalid filesystem file path
    When AddFile is called
    Then structured validation error is returned
    And error indicates invalid path
