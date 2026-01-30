@domain:file_mgmt @m2 @REQ-FILEMGMT-189 @spec(api_file_mgmt_updates.md#1-updatefile-operations)
Feature: Update file operations provide file modification capabilities

  @REQ-FILEMGMT-189 @happy
  Scenario: Update file operations provide file modification capabilities
    Given a package with file entries
    When update file operations are used
    Then file modification capabilities are provided as specified
    And the behavior matches the UpdateFile operations specification
    And UpdateFileMetadata, AddFilePath, RemoveFilePath, AddFileHash are available
