@domain:file_mgmt @m2 @REQ-FILEMGMT-055 @spec(api_file_mgmt_addition.md#221-purpose)
Feature: AddFilePattern Purpose and Usage

  @REQ-FILEMGMT-055 @happy
  Scenario: AddFilePattern purpose is to add multiple files via pattern matching
    Given an open NovusPack package
    And a valid context
    And a file pattern
    When AddFilePattern is used
    Then multiple files are added to package
    And files matching pattern are added
    And pattern-based file addition is enabled

  @REQ-FILEMGMT-055 @happy
  Scenario: AddFilePattern scans file system for matching files
    Given an open NovusPack package
    And a valid context
    And a file pattern
    When AddFilePattern is called
    Then file system is scanned for matching files
    And pattern matching is performed
    And matching files are identified

  @REQ-FILEMGMT-055 @happy
  Scenario: AddFilePattern returns created FileEntry objects
    Given an open NovusPack package
    And a valid context
    And a file pattern matching files
    When AddFilePattern is called
    Then slice of created FileEntry objects is returned
    And each FileEntry represents added file
    And FileEntry objects contain complete metadata
