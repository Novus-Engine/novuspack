@domain:file_mgmt @m2 @REQ-FILEMGMT-059 @spec(api_file_management.md#225-error-conditions)
Feature: AddFilePattern Error Conditions

  @REQ-FILEMGMT-059 @error
  Scenario: AddFilePattern handles package not open errors
    Given a closed NovusPack package
    And a valid context
    And a file pattern
    When AddFilePattern is called
    Then ErrPackageNotOpen is returned
    And error indicates package is not open
    And error follows structured error format

  @REQ-FILEMGMT-059 @error
  Scenario: AddFilePattern handles invalid pattern errors
    Given an open NovusPack package
    And a valid context
    And an invalid file pattern
    When AddFilePattern is called with invalid pattern
    Then ErrInvalidPattern is returned
    And error indicates invalid pattern
    And error follows structured error format

  @REQ-FILEMGMT-059 @error
  Scenario: AddFilePattern handles no files found errors
    Given an open NovusPack package
    And a valid context
    And a pattern matching no files
    When AddFilePattern is called
    Then ErrNoFilesFound is returned
    And error indicates no files matched pattern
    And error follows structured error format

  @REQ-FILEMGMT-059 @error
  Scenario: AddFilePattern handles I/O errors
    Given an open NovusPack package
    And a valid context
    And a file pattern
    And I/O operation failure occurs
    When AddFilePattern encounters I/O error
    Then ErrIOError is returned
    And error indicates I/O failure
    And error follows structured error format
