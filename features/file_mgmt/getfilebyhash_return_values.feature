@domain:file_mgmt @m2 @REQ-FILEMGMT-210 @spec(api_file_management.md#9223-returns)
Feature: GetFileByHash Return Values

  @REQ-FILEMGMT-210 @happy
  Scenario: GetFileByHash returns FileEntry and boolean when found
    Given an open NovusPack package
    And a valid context
    And a file exists with known hash
    When GetFileByHash is called with matching hashType and hashData
    Then FileEntry with matching hash is returned
    And boolean true is returned indicating found
    And FileEntry contains complete file information

  @REQ-FILEMGMT-210 @happy
  Scenario: GetFileByHash returns nil and false when not found
    Given an open NovusPack package
    And a valid context
    And a non-existent hash
    When GetFileByHash is called with non-existent hash
    Then nil FileEntry is returned
    And boolean false is returned indicating not found
    And no error is returned for not found case

  @REQ-FILEMGMT-210 @happy
  Scenario: GetFileByHash supports content-based lookup
    Given an open NovusPack package
    And a valid context
    And files exist in the package
    When GetFileByHash is called with content hash
    Then file is found by content rather than path
    And hash-based deduplication is supported
    And integrity verification is enabled

  @REQ-FILEMGMT-210 @happy
  Scenario: GetFileByHash returns complete FileEntry information
    Given an open NovusPack package
    And a valid context
    And a file exists with known hash
    When GetFileByHash is called
    Then returned FileEntry contains all metadata
    And FileEntry contains compression status
    And FileEntry contains encryption details
    And FileEntry contains checksums and timestamps

  @REQ-FILEMGMT-210 @happy
  Scenario: GetFileByHash supports multiple hash algorithms
    Given an open NovusPack package
    And a valid context
    And files with different hash types exist
    When GetFileByHash is called with different hash types
    Then SHA-256 hash lookup succeeds
    And SHA-512 hash lookup succeeds
    And BLAKE3 hash lookup succeeds
    And XXH3 hash lookup succeeds

  @REQ-FILEMGMT-210 @error
  Scenario: GetFileByHash handles context cancellation errors
    Given an open NovusPack package
    And a context that is cancelled
    And a hash type and hash data
    When GetFileByHash is called
    And context is cancelled
    Then structured context error is returned
    And error follows structured error format

  @REQ-FILEMGMT-210 @error
  Scenario: GetFileByHash handles context timeout errors
    Given an open NovusPack package
    And a context with timeout
    And a hash type and hash data
    When GetFileByHash is called
    And timeout is exceeded
    Then structured context timeout error is returned
    And error follows structured error format
