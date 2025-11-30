@domain:file_mgmt @m2 @REQ-FILEMGMT-126 @spec(api_file_management.md#1311-use-consistent-path-formats)
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
  Scenario: Consistent path formats avoid leading slashes
    Given an open NovusPack package
    And a valid context
    When file paths are used
    Then paths do not start with leading slash
    And relative paths are preferred
    And path format follows best practices

  @REQ-FILEMGMT-126 @error
  Scenario: Inconsistent path formats are rejected
    Given an open NovusPack package
    And a valid context
    And a path with mixed separators or leading slash
    When path operations are attempted with inconsistent format
    Then path validation fails
    And appropriate error is returned
    And error indicates invalid path format
