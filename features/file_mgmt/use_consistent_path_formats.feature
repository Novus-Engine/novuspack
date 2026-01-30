@domain:file_mgmt @m2 @REQ-FILEMGMT-126 @spec(api_file_mgmt_best_practices.md#1311-use-consistent-path-formats) @spec(api_file_mgmt_file_entry.md#1311-use-consistent-path-formats)
Feature: Use Consistent Path Formats

  @REQ-FILEMGMT-126 @happy
  Scenario: Consistent path formats use forward slashes
    Given an open NovusPack package
    And a valid context
    When file path operations are performed
    Then forward slashes are used as path separators
    And path format is consistent
    And path format is standardized

  @REQ-FILEMGMT-126 @happy
  Scenario: Consistent path formats require leading slashes for storage
    Given an open NovusPack package
    And a valid context
    When file paths are stored internally
    Then all stored paths start with leading slash
    And leading slash indicates package root
    And storage format follows best practices
    But when paths are displayed to users
    Then leading slash is stripped from displayed paths

  @REQ-FILEMGMT-126 @error
  Scenario: Inconsistent path formats are rejected
    Given an open NovusPack package
    And a valid context
    And a path with mixed separators
    When path operations are attempted with inconsistent format
    Then path validation fails
    And appropriate error is returned
    And error indicates invalid path format
