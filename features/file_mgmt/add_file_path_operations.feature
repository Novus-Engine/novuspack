@domain:file_mgmt @m2 @REQ-FILEMGMT-015 @spec(api_file_management.md#64-add-file-path)
Feature: Add file path operations

  @happy
  Scenario: AddFilePath adds additional path to existing file
    Given an open writable NovusPack package with existing file
    When AddFilePath is called with new path
    Then additional path is added to file entry
    And PathCount is incremented
    And file content is shared between paths
    And both paths point to same content

  @happy
  Scenario: AddFilePath supports path-specific metadata
    Given an open writable NovusPack package with existing file
    When AddFilePath is called with path and metadata
    Then path is added with specified metadata
    And path permissions are set correctly
    And path timestamps are set correctly

  @error
  Scenario: AddFilePath fails if file does not exist
    Given an open writable NovusPack package
    When AddFilePath is called with non-existent file
    Then structured validation error is returned
