@domain:file_mgmt @m2 @REQ-FILEMGMT-064 @spec(api_file_management.md#242-filepathsource)
Feature: FilePathSource Implementation

  @REQ-FILEMGMT-064 @happy
  Scenario: Filepathsource provides file system path-based data source
    Given a valid file system path
    When NewFilePathSource is called with the path
    Then FilePathSource is created
    And FileSource interface is implemented
    And file system path-based data source is provided

  @REQ-FILEMGMT-064 @happy
  Scenario: FilePathSource reads from filesystem files with streaming support
    Given a FilePathSource created from file path
    When file content is read
    Then content is read from filesystem file
    And streaming support is provided
    And large files are handled efficiently

  @REQ-FILEMGMT-064 @happy
  Scenario: FilePathSource automatically determines file type
    Given a FilePathSource created from file path
    When file is processed
    Then file type is automatically determined from extension
    And file type is automatically determined from content
    And appropriate processing is applied

  @REQ-FILEMGMT-064 @happy
  Scenario: FilePathSource handles large files efficiently
    Given a FilePathSource created from large file path
    When large file is processed
    Then memory management is efficient
    And streaming is used for large files
    And performance is optimized

  @REQ-FILEMGMT-064 @error
  Scenario: FilePathSource returns error for invalid path
    Given an invalid file system path
    When NewFilePathSource is called
    Then structured error is returned
    And error indicates invalid path
    And error follows structured error format

  @REQ-FILEMGMT-064 @error
  Scenario: FilePathSource respects context cancellation
    Given a valid file system path
    And a cancelled context
    When NewFilePathSource is called
    Then structured context error is returned
    And error follows structured error format
