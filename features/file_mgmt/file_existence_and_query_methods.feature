@domain:file_mgmt @m2 @REQ-FILEMGMT-028 @spec(api_file_mgmt_queries.md#1-file-existence-and-listing)
Feature: File existence and query methods

  @happy
  Scenario: FileExists checks if file exists
    Given an open package with files
    When FileExists is called with existing file path
    Then true is returned
    And file entry information is available

  @happy
  Scenario: FileExists returns false for non-existent file
    Given an open package
    When FileExists is called with non-existent path
    Then false is returned

  @REQ-FILEMGMT-029 @happy
  Scenario: ListFiles returns all files in package
    Given an open package with multiple files
    When ListFiles is called
    Then list of all file entries is returned
    And all files are included
    And file information is complete

  @REQ-FILEMGMT-030 @happy
  Scenario: FindEntriesByPathPatterns matches path patterns
    Given an open package with files matching patterns
    When FindEntriesByPathPatterns is called with patterns
    Then file entries matching patterns are returned
    And pattern matching works correctly

  @REQ-FILEMGMT-031 @happy
  Scenario: GetFileByPath retrieves file by path
    Given an open package with file
    When GetFileByPath is called with file path
    Then FileEntry with matching path is returned
    And file information is complete

  @REQ-FILEMGMT-032 @happy
  Scenario: GetFileByOffset retrieves file by offset
    Given an open package with files
    When GetFileByOffset is called with offset
    Then FileEntry at that offset is returned
    And file information is complete

  @error
  Scenario: GetFileByPath fails for non-existent file
    Given an open package
    When GetFileByPath is called with non-existent path
    Then structured validation error is returned

  @error
  Scenario: GetFileByOffset fails with invalid offset
    Given an open package
    When GetFileByOffset is called with invalid offset
    Then structured validation error is returned

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-038 @error
  Scenario: File query methods validate path parameters
    Given an open package
    When GetFileByPath is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-FILEMGMT-037 @REQ-FILEMGMT-041 @error
  Scenario: File query methods respect context cancellation
    Given an open package with files
    And a cancelled context
    When file query method is called
    Then structured context error is returned
    And error type is context cancellation
