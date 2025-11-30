@domain:file_mgmt @m2 @REQ-FILEMGMT-021 @spec(api_file_management.md#10-file-information-and-queries)
Feature: File information and query operations

  @happy
  Scenario: GetFileByFileID retrieves file by identifier
    Given an open NovusPack package with files
    When GetFileByFileID is called with FileID
    Then FileEntry with matching FileID is returned
    And file information is complete

  @happy
  Scenario: GetFileByHash retrieves file by hash
    Given an open NovusPack package with files
    When GetFileByHash is called with hash type and data
    Then FileEntry with matching hash is returned
    And hash matching works correctly

  @happy
  Scenario: GetFileByChecksum retrieves file by checksum
    Given an open NovusPack package with files
    When GetFileByChecksum is called with checksum
    Then FileEntry with matching checksum is returned
    And checksum matching works correctly

  @happy
  Scenario: FindEntriesByTag retrieves files by tag
    Given an open NovusPack package with tagged files
    When FindEntriesByTag is called with tag key and value
    Then all FileEntries with matching tag are returned
    And tag matching is accurate

  @happy
  Scenario: FindEntriesByType retrieves files by type
    Given an open NovusPack package with files of various types
    When FindEntriesByType is called with file type
    Then all FileEntries with matching type are returned
    And type matching is accurate

  @happy
  Scenario: GetFileCount returns total file count
    Given an open NovusPack package
    When GetFileCount is called
    Then total number of files is returned
    And count is accurate

  @happy
  Scenario: File existence check works correctly
    Given an open NovusPack package
    When file existence is checked
    Then true is returned if file exists
    And false is returned if file does not exist
    And existence check is efficient

  @error
  Scenario: File queries return errors for non-existent items
    Given an open NovusPack package
    When query is performed for non-existent item
    Then appropriate error or empty result is returned
    And error indicates item not found

  @error
  Scenario: File queries respect context cancellation
    Given an open NovusPack package
    And a cancelled context
    When file query is called
    Then structured context error is returned
