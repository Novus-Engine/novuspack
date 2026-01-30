@domain:file_mgmt @m2 @REQ-FILEMGMT-065 @spec(api_file_mgmt_queries.md#234-getfilebyfileid-returns)
Feature: GetFileByFileID Return Values

  @REQ-FILEMGMT-065 @happy
  Scenario: GetFileByFileID returns FileEntry and boolean when found
    Given an open NovusPack package
    And a valid context
    And a file exists with known FileID
    When GetFileByFileID is called with existing fileID
    Then FileEntry with matching FileID is returned
    And boolean true is returned indicating found
    And FileEntry contains complete file information

  @REQ-FILEMGMT-065 @happy
  Scenario: GetFileByFileID returns nil and false when not found
    Given an open NovusPack package
    And a valid context
    And a non-existent FileID
    When GetFileByFileID is called with non-existent fileID
    Then nil FileEntry is returned
    And boolean false is returned indicating not found
    And no error is returned for not found case

  @REQ-FILEMGMT-065 @happy
  Scenario: GetFileByFileID returns stable file references
    Given an open NovusPack package
    And a valid context
    And files exist in the package
    When GetFileByFileID is called with FileID
    Then returned FileEntry has stable FileID
    And FileID persists across package modifications
    And FileID can be used for database-style lookups

  @REQ-FILEMGMT-065 @happy
  Scenario: GetFileByFileID returns complete FileEntry information
    Given an open NovusPack package
    And a valid context
    And a file exists with known FileID
    When GetFileByFileID is called
    Then returned FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums and timestamps

  @REQ-FILEMGMT-065 @error
  Scenario: GetFileByFileID handles context cancellation errors
    Given an open NovusPack package
    And a context that is cancelled
    And a FileID
    When GetFileByFileID is called
    And context is cancelled
    Then structured context error is returned
    And error follows structured error format

  @REQ-FILEMGMT-065 @error
  Scenario: GetFileByFileID handles context timeout errors
    Given an open NovusPack package
    And a context with timeout
    And a FileID
    When GetFileByFileID is called
    And timeout is exceeded
    Then structured context timeout error is returned
    And error follows structured error format
