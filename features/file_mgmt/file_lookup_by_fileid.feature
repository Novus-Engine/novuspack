@domain:file_mgmt @m2 @REQ-FILEMGMT-062 @REQ-FILEMGMT-063 @REQ-FILEMGMT-064 @REQ-FILEMGMT-065 @REQ-FILEMGMT-250 @spec(api_file_mgmt_queries.md#23-getfilebyfileid)
Feature: File lookup by FileID

  @REQ-FILEMGMT-063 @REQ-FILEMGMT-064 @happy
  Scenario: GetFileByFileID finds file by unique identifier
    Given an open package with file "data/config.json" having FileID 12345
    When GetFileByFileID is called with fileID 12345
    Then the file entry for "data/config.json" is returned
    And the found flag is true

  @REQ-FILEMGMT-063 @REQ-FILEMGMT-064 @happy
  Scenario: GetFileByFileID returns false when file not found
    Given an open package
    When GetFileByFileID is called with non-existent fileID 99999
    Then nil file entry is returned
    And the found flag is false

  @REQ-FILEMGMT-250 @happy
  Scenario: GetFileByFileID provides stable file reference across modifications
    Given an open package with file "important.txt"
    And the file has FileID 5000
    When GetFileByFileID is called with fileID 5000
    Then the file entry is returned
    When the file path is changed to "updated.txt"
    And GetFileByFileID is called again with fileID 5000
    Then the same file entry is returned with updated path

  @REQ-FILEMGMT-250 @happy
  Scenario: GetFileByFileID supports database-style lookups
    Given an open package with multiple files
    And each file has a unique FileID
    When GetFileByFileID is called for each fileID
    Then the correct file entry is returned for each ID
    And all lookups complete successfully

  @REQ-FILEMGMT-064 @error
  Scenario: GetFileByFileID respects context cancellation
    Given an open package
    And a cancelled context
    When GetFileByFileID is called
    Then a structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-065 @happy
  Scenario: GetFileByFileID returns complete FileEntry structure
    Given an open package with file containing metadata, compression, and encryption info
    When GetFileByFileID is called
    Then the returned FileEntry contains all metadata
    And compression status is included
    And encryption details are included
    And checksums are included
    And timestamps are included
