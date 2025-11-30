@domain:file_mgmt @m2 @REQ-FILEMGMT-009 @REQ-FILEMGMT-026 @spec(api_file_management.md#23-filesource-interface)
Feature: FileSource interface and implementations

  @happy
  Scenario: FileSource interface defines required methods
    Given FileSource interface
    When interface is examined
    Then Read method is defined
    And Size method is defined
    And Close method is defined
    And interface enables flexible file sources

  @happy
  Scenario: FilePathSource reads from file system
    Given a FilePathSource with file path
    When file source is used
    Then file is read from file system
    And file content is accessible
    And file handle is managed correctly

  @happy
  Scenario: MemorySource reads from byte slice
    Given a MemorySource with byte data
    When file source is used
    Then data is read from memory
    And byte slice content is accessible
    And no file system access is required

  @happy
  Scenario: FileSource implementations are interchangeable
    Given FilePathSource and MemorySource
    When same operations are performed
    Then both sources provide identical data
    And interface abstraction works correctly

  @happy
  Scenario: FileSource Read method provides data
    Given a FileSource instance
    When Read is called
    Then file data is returned
    And data matches source content
    And read operation is efficient

  @happy
  Scenario: FileSource Size method returns data size
    Given a FileSource instance
    When Size is called
    Then data size is returned
    And size matches actual data length
    And size is accurate

  @happy
  Scenario: FileSource Close method releases resources
    Given a FileSource instance
    When Close is called
    Then resources are released
    And file handles are closed
    And memory is cleaned up

  @error
  Scenario: FileSource operations handle I/O errors
    Given a FileSource with I/O error
    When Read or other operation is called
    Then structured I/O error is returned
    And error indicates source issue
