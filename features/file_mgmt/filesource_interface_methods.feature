@domain:file_mgmt @m2 @REQ-FILEMGMT-062 @REQ-FILEMGMT-063 @REQ-FILEMGMT-009 @spec(api_file_management.md#232-methods)
Feature: FileSource Interface Methods

  @REQ-FILEMGMT-062 @REQ-FILEMGMT-063 @REQ-FILEMGMT-009 @happy
  Scenario: Filesource interface methods provide data access operations
    Given FileSource interface
    When interface methods are examined
    Then Read method reads file data into buffer
    And Size method returns total size of file data
    And Close method closes source and releases resources
    And data access operations are provided

  @REQ-FILEMGMT-062 @REQ-FILEMGMT-063 @REQ-FILEMGMT-009 @happy
  Scenario: Filesource implementations provide concrete data source types
    Given FileSource interface
    When implementations are examined
    Then FilePathSource implements FileSource for filesystem files
    And MemorySource implements FileSource for in-memory data
    And custom implementations can be created
    And concrete data source types are provided

  @REQ-FILEMGMT-009 @happy
  Scenario: FileSource interface provides unified file data sources
    Given FileSource interface
    When file operations use FileSource
    Then unified interface is used
    And various data sources are supported
    And source abstraction enables flexibility

  @REQ-FILEMGMT-009 @happy
  Scenario: FileSource Read method supports streaming
    Given FileSource implementation
    And context for cancellation
    When Read is called with buffer
    Then data is read into buffer
    And streaming is supported
    And context cancellation is respected

  @REQ-FILEMGMT-009 @happy
  Scenario: FileSource Size method provides file size information
    Given FileSource implementation
    And context for cancellation
    When Size is called
    Then total size of file data is returned
    And size information is accurate
    And context cancellation is respected

  @REQ-FILEMGMT-009 @happy
  Scenario: FileSource Close method releases resources
    Given FileSource implementation
    And context for cancellation
    When Close is called
    Then source is closed
    And resources are released
    And cleanup is performed

  @REQ-FILEMGMT-009 @error
  Scenario: FileSource methods handle context cancellation
    Given FileSource implementation
    And a cancelled context
    When Read, Size, or Close is called
    Then context cancellation error is returned
    And error follows structured error format
