@domain:file_mgmt @m2 @REQ-FILEMGMT-222 @REQ-FILEMGMT-225 @REQ-FILEMGMT-226 @spec(api_file_management.md#925-findentriesbytype)
Feature: FindEntriesByType

  @REQ-FILEMGMT-222 @happy
  Scenario: FindEntriesByType finds all file entries of specific type
    Given an open NovusPack package
    And a valid context
    And files of various types exist in the package
    When FindEntriesByType is called with fileType
    Then all file entries of matching type are returned
    And slice of FileEntry objects is returned
    And error is nil on success

  @REQ-FILEMGMT-225 @happy
  Scenario: FindEntriesByType returns slice of FileEntry objects
    Given an open NovusPack package
    And a valid context
    And files of various types exist
    When FindEntriesByType is called
    Then slice of FileEntry objects is returned
    And each FileEntry has the specified file type
    And all matching files are included in results

  @REQ-FILEMGMT-226 @happy
  Scenario: FindEntriesByType supports finding files of specific format
    Given an open NovusPack package
    And a valid context
    And files of different formats exist
    When FindEntriesByType is called with format file type
    Then all files of matching format are found
    And format-based file processing is enabled

  @REQ-FILEMGMT-226 @happy
  Scenario: FindEntriesByType supports type-based file processing
    Given an open NovusPack package
    And a valid context
    And files of various types exist
    When FindEntriesByType is used for file processing
    Then type-based file queries are enabled
    And files can be processed by type
    And type-based operations are supported

  @REQ-FILEMGMT-226 @happy
  Scenario: FindEntriesByType supports file organization by category
    Given an open NovusPack package
    And a valid context
    And files of various types exist
    When FindEntriesByType is used for organization
    Then files can be organized by type category
    And type-based organization is supported
    And file categorization is enabled

  @REQ-FILEMGMT-222 @error
  Scenario: FindEntriesByType handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a file type
    When FindEntriesByType is called
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-222 @happy
  Scenario: FindEntriesByType respects context cancellation
    Given an open NovusPack package
    And a context that can be cancelled
    And a file type
    When FindEntriesByType is called
    And context is cancelled
    Then operation respects context cancellation
    And structured context error is returned
