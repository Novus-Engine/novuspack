@domain:file_mgmt @m2 @REQ-FILEMGMT-057 @spec(api_file_mgmt_addition.md#223-returns)
Feature: AddFilePattern Returns Slice of Created FileEntry Objects

  @REQ-FILEMGMT-057 @happy
  Scenario: AddFilePattern returns slice of created fileentry objects
    Given an open NovusPack package
    And a valid context
    And file system has files matching pattern
    When AddFilePattern is called with context, pattern, and options
    Then slice of created FileEntry objects is returned
    And each FileEntry has all metadata, compression status, encryption details, and checksums
    And slice contains entries for all matching files

  @REQ-FILEMGMT-057 @happy
  Scenario: AddFilePattern creates FileEntry for each matching file
    Given an open NovusPack package
    And a valid context
    And file system has files matching pattern "*.txt"
    When AddFilePattern is called with pattern "*.txt"
    Then FileEntry is created for each matching file
    And each FileEntry contains complete metadata
    And each FileEntry includes compression and encryption details

  @REQ-FILEMGMT-057 @happy
  Scenario: AddFilePattern handles pattern scanning and bulk file addition
    Given an open NovusPack package
    And a valid context
    And file system has multiple files matching pattern
    When AddFilePattern is called
    Then pattern scanning identifies matching files
    And bulk file addition adds all matching files
    And all files are processed efficiently

  @REQ-FILEMGMT-057 @happy
  Scenario: AddFilePattern applies pattern-specific filters
    Given an open NovusPack package
    And a valid context
    And AddFileOptions includes exclude patterns and max file size
    When AddFilePattern is called with options
    Then pattern-specific filters are applied
    And excluded files are not added
    And oversized files are not added
    And filter behavior matches options

  @REQ-FILEMGMT-057 @happy
  Scenario: AddFilePattern preserves directory structure when requested
    Given an open NovusPack package
    And a valid context
    And AddFileOptions includes directory structure preservation
    When AddFilePattern is called with options
    Then directory structure is preserved
    And file paths reflect original directory hierarchy
    And structure preservation matches options

  @REQ-FILEMGMT-057 @error
  Scenario: AddFilePattern returns error when package is not open
    Given a NovusPack package that is not open
    And a valid context
    When AddFilePattern is called
    Then package not open error is returned
    And error indicates package must be open

  @REQ-FILEMGMT-057 @error
  Scenario: AddFilePattern returns error for invalid pattern
    Given an open NovusPack package
    And a valid context
    And an invalid or malformed pattern
    When AddFilePattern is called
    Then structured invalid pattern error is returned
    And error indicates pattern validation failure

  @REQ-FILEMGMT-057 @error
  Scenario: AddFilePattern returns error when no files match pattern
    Given an open NovusPack package
    And a valid context
    And no files match the pattern
    When AddFilePattern is called
    Then no files found error may be returned
    And error indicates pattern matched no files
