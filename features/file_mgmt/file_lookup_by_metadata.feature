@domain:file_mgmt @m2 @REQ-FILEMGMT-093 @spec(api_file_management.md#102-file-lookup-by-metadata)
Feature: File Lookup by Metadata

  @REQ-FILEMGMT-093 @happy
  Scenario: File lookup by metadata supports fileID lookup
    Given an open NovusPack package
    And a valid context
    And files with FileIDs exist in the package
    When GetFileByFileID is called with a fileID
    Then file entry with matching FileID is returned
    And boolean true is returned if found
    And boolean false is returned if not found

  @REQ-FILEMGMT-093 @happy
  Scenario: File lookup by metadata supports hash lookup
    Given an open NovusPack package
    And a valid context
    And files with hashes exist in the package
    When GetFileByHash is called with hashType and hashData
    Then file entry with matching hash is returned
    And boolean true is returned if found
    And boolean false is returned if not found

  @REQ-FILEMGMT-093 @happy
  Scenario: File lookup by metadata supports checksum lookup
    Given an open NovusPack package
    And a valid context
    And files with CRC32 checksums exist in the package
    When GetFileByChecksum is called with checksum
    Then file entry with matching checksum is returned
    And boolean true is returned if found
    And boolean false is returned if not found

  @REQ-FILEMGMT-093 @happy
  Scenario: File lookup by metadata supports tag-based search
    Given an open NovusPack package
    And a valid context
    And files with tags exist in the package
    When FindEntriesByTag is called with tag string
    Then all file entries with matching tag are returned
    And slice of FileEntry objects is returned
    And error is nil on success

  @REQ-FILEMGMT-093 @happy
  Scenario: File lookup by metadata supports type-based search
    Given an open NovusPack package
    And a valid context
    And files with types exist in the package
    When FindEntriesByType is called with fileType
    Then all file entries of matching type are returned
    And slice of FileEntry objects is returned
    And error is nil on success

  @REQ-FILEMGMT-093 @error
  Scenario: File lookup by metadata handles package not open errors
    Given a closed NovusPack package
    And a valid context
    When metadata lookup operations are attempted
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-093 @happy
  Scenario: File lookup by metadata respects context cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    When metadata lookup operations are performed
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned
