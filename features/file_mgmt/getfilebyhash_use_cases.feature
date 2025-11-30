@domain:file_mgmt @m2 @REQ-FILEMGMT-211 @spec(api_file_management.md#9224-use-cases)
Feature: GetFileByHash Use Cases

  @REQ-FILEMGMT-211 @happy
  Scenario: GetFileByHash supports content deduplication
    Given an open NovusPack package
    And a valid context
    And files with matching content hashes exist
    When GetFileByHash is called with content hash
    Then duplicate content can be identified
    And deduplication can be performed
    And storage space can be saved

  @REQ-FILEMGMT-211 @happy
  Scenario: GetFileByHash supports integrity verification
    Given an open NovusPack package
    And a valid context
    And a file with known hash
    When GetFileByHash is called with expected hash
    Then file integrity can be verified
    And content corruption can be detected
    And hash verification ensures data integrity

  @REQ-FILEMGMT-211 @happy
  Scenario: GetFileByHash supports finding files by content rather than path
    Given an open NovusPack package
    And a valid context
    And files with known content hashes exist
    When GetFileByHash is called with content hash
    Then file is found by content
    And path is not required for lookup
    And content-based file identification is enabled
