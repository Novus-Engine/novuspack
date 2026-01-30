@domain:file_mgmt @m2 @REQ-FILEMGMT-428 @spec(api_file_mgmt_file_entry.md#102-fileentry-processdata-method)
Feature: FileEntry ProcessData method processes the file data (compression, encryption, etc.)

  @REQ-FILEMGMT-428 @happy
  Scenario: ProcessData processes file data
    Given a FileEntry with loaded file content
    When ProcessData is called on the FileEntry
    Then file data is processed (compression, encryption, etc.) as specified
    And the behavior matches the ProcessData method specification
    And processing applies compression or encryption per configuration
    And error is returned on processing failure
