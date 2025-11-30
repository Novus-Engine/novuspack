@domain:file_mgmt @m2 @REQ-FILEMGMT-217 @REQ-FILEMGMT-220 @REQ-FILEMGMT-221 @spec(api_file_management.md#924-findentriesbytag)
Feature: FindEntriesByTag

  @REQ-FILEMGMT-217 @happy
  Scenario: FindEntriesByTag finds all file entries with specific tag
    Given an open NovusPack package
    And a valid context
    And files with tags exist in the package
    When FindEntriesByTag is called with tag string
    Then all file entries with matching tag are returned
    And slice of FileEntry objects is returned
    And error is nil on success

  @REQ-FILEMGMT-220 @happy
  Scenario: FindEntriesByTag returns slice of FileEntry objects
    Given an open NovusPack package
    And a valid context
    And files with tags exist in the package
    When FindEntriesByTag is called
    Then slice of FileEntry objects is returned
    And each FileEntry has the specified tag
    And all matching files are included in results

  @REQ-FILEMGMT-221 @happy
  Scenario: FindEntriesByTag supports finding files with specific label
    Given an open NovusPack package
    And a valid context
    And files with label tags exist
    When FindEntriesByTag is called with label tag
    Then all files with matching label are found
    And tag-based file organization is enabled

  @REQ-FILEMGMT-221 @happy
  Scenario: FindEntriesByTag supports organizing files by category
    Given an open NovusPack package
    And a valid context
    And files with category tags exist
    When FindEntriesByTag is called with category tag
    Then all files in matching category are found
    And category-based file organization is supported

  @REQ-FILEMGMT-221 @happy
  Scenario: FindEntriesByTag supports tag-based file management
    Given an open NovusPack package
    And a valid context
    And files with various tags exist
    When FindEntriesByTag is used for file management
    Then tag-based file queries are enabled
    And file management operations use tags
    And tag-based organization is supported

  @REQ-FILEMGMT-217 @error
  Scenario: FindEntriesByTag handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a tag string
    When FindEntriesByTag is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-217 @happy
  Scenario: FindEntriesByTag respects context cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    And a tag string
    When FindEntriesByTag is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned
