@domain:file_mgmt @m2 @REQ-FILEMGMT-251 @REQ-FILEMGMT-210 @REQ-FILEMGMT-211 @spec(api_file_mgmt_queries.md#24-getfilebyhash)
Feature: GetFileByHash

  @REQ-FILEMGMT-251 @happy
  Scenario: GetFileByHash finds file entry by content hash
    Given an open NovusPack package
    And a valid context
    And files with hashes exist
    When GetFileByHash is called with hashType and hashData
    Then file entry with matching hash is returned
    And boolean true is returned if found
    And boolean false is returned if not found

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
