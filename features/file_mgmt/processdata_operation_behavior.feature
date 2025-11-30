@domain:file_mgmt @m2 @REQ-FILEMGMT-107 @spec(api_file_management.md#1133-processdata-behavior)
Feature: ProcessData Operation Behavior

  @REQ-FILEMGMT-107 @happy
  Scenario: ProcessData applies compression to file data if requested
    Given an open NovusPack package
    And a valid context
    And a FileEntry with loaded data
    And compression is requested
    When ProcessData is called on the FileEntry
    Then compression is applied to file data
    And compressed data is prepared
    And file entry metadata is updated

  @REQ-FILEMGMT-107 @happy
  Scenario: ProcessData applies encryption to file data if requested
    Given an open NovusPack package
    And a valid context
    And a FileEntry with loaded data
    And encryption is requested
    When ProcessData is called on the FileEntry
    Then encryption is applied to file data
    And encrypted data is prepared
    And file entry metadata is updated

  @REQ-FILEMGMT-107 @happy
  Scenario: ProcessData updates file entry metadata
    Given an open NovusPack package
    And a valid context
    And a FileEntry with loaded data
    When ProcessData is called
    Then file entry metadata is updated
    And compression status is updated if compressed
    And encryption status is updated if encrypted
    And data is prepared for storage

  @REQ-FILEMGMT-107 @error
  Scenario: ProcessData handles compression failures
    Given an open NovusPack package
    And a valid context
    And a FileEntry with loaded data
    And compression is requested
    And compression fails
    When ProcessData is called and compression fails
    Then compression error is returned
    And error indicates compression failure
    And error follows structured error format

  @REQ-FILEMGMT-107 @error
  Scenario: ProcessData handles encryption failures
    Given an open NovusPack package
    And a valid context
    And a FileEntry with loaded data
    And encryption is requested
    And encryption fails
    When ProcessData is called and encryption fails
    Then encryption error is returned
    And error indicates encryption failure
    And error follows structured error format
