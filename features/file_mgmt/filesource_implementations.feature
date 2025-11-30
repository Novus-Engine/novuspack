@domain:file_mgmt @m2 @REQ-FILEMGMT-026 @spec(api_file_management.md#24-filesource-implementations)
Feature: FileSource Implementations

  @REQ-FILEMGMT-026 @happy
  Scenario: FilePathSource implements FileSource for filesystem files
    Given an open NovusPack package
    And a valid context
    And a filesystem file path
    When NewFilePathSource is called with path
    Then FilePathSource is created
    And FilePathSource reads from filesystem files
    And FilePathSource supports streaming for large files
    And file type is determined from extension and content

  @REQ-FILEMGMT-026 @happy
  Scenario: MemorySource implements FileSource for in-memory data
    Given an open NovusPack package
    And a valid context
    And in-memory byte data
    When NewMemorySource is called with data
    Then MemorySource is created
    And MemorySource reads from in-memory data
    And MemorySource is suitable for small files
    And MemorySource has no filesystem overhead

  @REQ-FILEMGMT-026 @happy
  Scenario: FileSource implementations provide unified interface
    Given an open NovusPack package
    And a valid context
    When FileSource implementations are used with AddFile
    Then FilePathSource works with AddFile
    And MemorySource works with AddFile
    And custom FileSource implementations work with AddFile
    And unified interface enables flexibility

  @REQ-FILEMGMT-026 @happy
  Scenario: FilePathSource handles large files efficiently
    Given an open NovusPack package
    And a valid context
    And a large filesystem file
    When FilePathSource is used with large file
    Then streaming is used for memory efficiency
    And memory management is efficient
    And large files are handled without excessive memory usage

  @REQ-FILEMGMT-026 @error
  Scenario: FilePathSource handles file access errors
    Given an open NovusPack package
    And a valid context
    And an inaccessible file path
    When NewFilePathSource is called with inaccessible path
    Then appropriate error is returned
    And error indicates file access issue
    And error follows structured error format
