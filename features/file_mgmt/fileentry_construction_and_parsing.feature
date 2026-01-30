@domain:file_mgmt @m2 @REQ-FILEMGMT-022 @spec(api_file_mgmt_file_entry.md#111-file-entry-properties)
Feature: FileEntry construction and parsing

  @happy
  Scenario: NewFileEntry constructs new FileEntry
    Given file entry parameters
    When NewFileEntry is called with parameters
    Then FileEntry is constructed
    And file entry fields are set correctly
    And file entry is valid

  @happy
  Scenario: ParseFileEntry parses raw data into FileEntry
    Given raw file entry data
    When ParseFileEntry is called with data
    Then FileEntry is parsed from data
    And file entry fields are populated correctly
    And file entry matches data

  @happy
  Scenario: LoadFileEntry loads FileEntry from raw data
    Given raw file entry data
    When LoadFileEntry is called with data
    Then FileEntry is loaded from data
    And file entry is complete
    And file entry can be used

  @error
  Scenario: ParseFileEntry fails with invalid data
    Given invalid raw file entry data
    When ParseFileEntry is called with data
    Then structured validation error is returned
    And error indicates parsing failure

  @error
  Scenario: LoadFileEntry fails with corrupted data
    Given corrupted file entry data
    When LoadFileEntry is called with data
    Then structured validation error is returned
    And error indicates data corruption
