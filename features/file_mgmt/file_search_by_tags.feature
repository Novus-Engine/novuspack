@domain:file_mgmt @m2 @REQ-FILEMGMT-217 @REQ-FILEMGMT-218 @REQ-FILEMGMT-219 @REQ-FILEMGMT-220 @REQ-FILEMGMT-221 @spec(api_file_mgmt_queries.md#31-findentriesbytag)
Feature: File Management: File Search by Tags

  @REQ-FILEMGMT-218 @REQ-FILEMGMT-219 @happy
  Scenario: FindEntriesByTag finds files with matching tag key
    Given an open package
    And file "image1.png" has tag "category" = "texture"
    And file "image2.jpg" has tag "category" = "texture"
    And file "data.json" has tag "category" = "config"
    When FindEntriesByTag is called with tag "category"
    Then all files with "category" tag are returned
    And the result includes "image1.png" and "image2.jpg" and "data.json"

  @REQ-FILEMGMT-218 @REQ-FILEMGMT-219 @happy
  Scenario: FindEntriesByTag returns empty slice when no matches
    Given an open package without files having tag "priority"
    When FindEntriesByTag is called with tag "priority"
    Then an empty slice is returned
    And no error occurs

  @REQ-FILEMGMT-221 @happy
  Scenario: FindEntriesByTag supports tag-based file organization
    Given an open package with files tagged by type
    And files have tags "type:texture", "type:audio", "type:config"
    When FindEntriesByTag is called with tag "type"
    Then all files with type tag are returned
    And files can be organized by tag value

  @REQ-FILEMGMT-221 @happy
  Scenario: FindEntriesByTag finds files with multiple tags
    Given an open package
    And file "asset.png" has tags "category:texture" and "priority:high"
    When FindEntriesByTag is called with tag "category"
    Then "asset.png" is included in results
    When FindEntriesByTag is called with tag "priority"
    Then "asset.png" is included in results

  @REQ-FILEMGMT-219 @error
  Scenario: FindEntriesByTag validates tag parameter
    Given an open package
    When FindEntriesByTag is called with empty tag
    Then structured validation error is returned
    And error indicates invalid tag

  @REQ-FILEMGMT-219 @error
  Scenario: FindEntriesByTag respects context cancellation
    Given an open package
    And a cancelled context
    When FindEntriesByTag is called
    Then a structured context error is returned
    And error type is context cancellation

  @REQ-FILEMGMT-220 @happy
  Scenario: FindEntriesByTag returns complete FileEntry objects
    Given an open package with tagged files
    When FindEntriesByTag is called
    Then returned FileEntry objects contain tag information
    And all file metadata is accessible
    And tags can be retrieved from FileEntry
