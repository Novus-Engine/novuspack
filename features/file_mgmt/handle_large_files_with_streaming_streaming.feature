@domain:file_mgmt @m2 @REQ-FILEMGMT-133 @spec(api_file_mgmt_best_practices.md#1332-handle-large-files-with-streaming)
Feature: Handle Large Files with Streaming

  @REQ-FILEMGMT-133 @happy
  Scenario: Large files are handled efficiently with streaming
    Given an open NovusPack package
    And a valid context
    And a large file to process
    When file operations are performed with large file
    Then streaming is used for memory efficiency
    And memory usage is controlled
    And large files are handled without excessive memory consumption

  @REQ-FILEMGMT-133 @happy
  Scenario: Streaming supports efficient file processing
    Given an open NovusPack package
    And a valid context
    And a large file to process
    When streaming is used for file operations
    Then file data is processed incrementally
    And memory footprint is minimized
    And processing completes efficiently

  @REQ-FILEMGMT-133 @happy
  Scenario: Streaming works when reading from filesystem paths
    Given an open NovusPack package
    And a valid context
    And a filesystem file path
    And a large file to process
    When AddFile is called
    Then streaming is used for large files when needed
    And memory usage remains reasonable
