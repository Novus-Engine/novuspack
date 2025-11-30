@domain:file_mgmt @m2 @REQ-FILEMGMT-125 @spec(api_file_management.md#131-file-path-management)
Feature: File Path Management

  @REQ-FILEMGMT-125 @happy
  Scenario: File path management uses consistent path formats
    Given an open NovusPack package
    And a valid context
    When file path operations are performed
    Then consistent path formats are used
    And path format is standardized
    And path normalization is applied

  @REQ-FILEMGMT-125 @happy
  Scenario: File path management validates paths before use
    Given an open NovusPack package
    And a valid context
    When file path operations are performed
    Then paths are validated before use
    And path format is checked
    And path validity is verified
    And invalid paths are rejected

  @REQ-FILEMGMT-125 @happy
  Scenario: File path management supports path operations
    Given an open NovusPack package
    And a valid context
    And files with paths exist in the package
    When path operations are performed
    Then paths can be queried
    And paths can be resolved
    And path relationships can be determined

  @REQ-FILEMGMT-125 @error
  Scenario: File path management handles invalid paths
    Given an open NovusPack package
    And a valid context
    And an invalid path
    When path operations are attempted with invalid path
    Then a structured error is returned
    And error indicates invalid path
    And error follows structured error format

  @REQ-FILEMGMT-125 @error
  Scenario: File path management handles package not open errors
    Given a closed NovusPack package
    And a valid context
    When path operations are attempted
    Then a structured error is returned
    And error indicates package is not open
    And error follows structured error format
