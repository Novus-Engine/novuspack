@domain:file_mgmt @m2 @REQ-FILEMGMT-300 @spec(api_file_mgmt_file_entry.md#21-newfileentry-function)
Feature: NewFileEntry creates a new FileEntry with proper tag synchronization

  @REQ-FILEMGMT-300 @happy
  Scenario: NewFileEntry creates FileEntry with tag sync
    Given a package context and file entry parameters
    When NewFileEntry is called
    Then a new FileEntry is created with proper tag synchronization
    And the behavior matches the NewFileEntry function specification
    And tags are synchronized with underlying storage
    And FileEntry is ready for use
