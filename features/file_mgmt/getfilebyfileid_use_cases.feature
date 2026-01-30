@domain:file_mgmt @m2 @REQ-FILEMGMT-250 @spec(api_file_mgmt_queries.md#235-getfilebyfileid-use-cases)
Feature: GetFileByFileID Use Cases

  @REQ-FILEMGMT-250 @happy
  Scenario: GetFileByFileID provides stable file references across package modifications
    Given an open NovusPack package
    And a valid context
    And files exist in the package with FileIDs
    When GetFileByFileID is called with a FileID
    Then file entry is returned with stable FileID
    And FileID persists when file path changes
    And FileID persists when file content is updated
    And stable references enable reliable file tracking

  @REQ-FILEMGMT-250 @happy
  Scenario: GetFileByFileID supports database-style lookups by primary key
    Given an open NovusPack package
    And a valid context
    And files exist in the package with unique FileIDs
    When GetFileByFileID is called with known FileIDs
    Then correct file entries are returned for each FileID
    And lookup performance is optimized
    And database-style access patterns are supported

  @REQ-FILEMGMT-250 @happy
  Scenario: GetFileByFileID supports file tracking and management systems
    Given an open NovusPack package
    And a valid context
    And a file tracking system
    When file tracking uses FileID for references
    Then file references remain valid across modifications
    And file tracking system maintains stable references
    And file management operations use FileID for consistency
