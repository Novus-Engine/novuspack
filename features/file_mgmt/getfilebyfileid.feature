@domain:file_mgmt @m2 @REQ-FILEMGMT-202 @REQ-FILEMGMT-205 @REQ-FILEMGMT-206 @spec(api_file_management.md#921-getfilebyfileid)
Feature: GetFileByFileID

  @REQ-FILEMGMT-202 @happy
  Scenario: GetFileByFileID finds file entry by unique 64-bit identifier
    Given an open NovusPack package
    And a valid context
    And files with FileIDs exist
    When GetFileByFileID is called with fileID
    Then file entry with matching FileID is returned
    And boolean true is returned if found
    And boolean false is returned if not found

  @REQ-FILEMGMT-205 @happy
  Scenario: GetFileByFileID returns FileEntry and boolean when found
    Given an open NovusPack package
    And a valid context
    And a file exists with known FileID
    When GetFileByFileID is called with existing fileID
    Then FileEntry with matching FileID is returned
    And boolean true is returned indicating found
    And FileEntry contains complete file information

  @REQ-FILEMGMT-205 @happy
  Scenario: GetFileByFileID returns nil and false when not found
    Given an open NovusPack package
    And a valid context
    And a non-existent FileID
    When GetFileByFileID is called with non-existent fileID
    Then nil FileEntry is returned
    And boolean false is returned indicating not found
    And no error is returned for not found case

  @REQ-FILEMGMT-206 @happy
  Scenario: GetFileByFileID provides stable file references across package modifications
    Given an open NovusPack package
    And a valid context
    And files exist in the package with FileIDs
    When GetFileByFileID is called with a FileID
    Then file entry is returned with stable FileID
    And FileID persists when file path changes
    And FileID persists when file content is updated
    And stable references enable reliable file tracking

  @REQ-FILEMGMT-206 @happy
  Scenario: GetFileByFileID supports database-style lookups by primary key
    Given an open NovusPack package
    And a valid context
    And files exist in the package with unique FileIDs
    When GetFileByFileID is called with known FileIDs
    Then correct file entries are returned for each FileID
    And lookup performance is optimized
    And database-style access patterns are supported

  @REQ-FILEMGMT-206 @happy
  Scenario: GetFileByFileID supports file tracking and management systems
    Given an open NovusPack package
    And a valid context
    And a file tracking system
    When file tracking uses FileID for references
    Then file references remain valid across modifications
    And file tracking system maintains stable references
    And file management operations use FileID for consistency
