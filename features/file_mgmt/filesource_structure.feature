@domain:file_mgmt @m2 @REQ-FILEMGMT-338 @spec(api_file_mgmt_file_entry.md#16-filesource-structure)
Feature: FileSource structure defines unified source location tracking for file data

  @REQ-FILEMGMT-338 @happy
  Scenario: FileSource structure defines unified source tracking
    Given a FileEntry with file data source
    When FileSource is used for source location tracking
    Then FileSource structure defines unified source location for file data (original, intermediate, or final stage)
    And the behavior matches the FileSource structure specification
    And source type and location are tracked
    And CurrentSource and OriginalSource are available
