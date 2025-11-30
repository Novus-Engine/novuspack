@domain:core @m1 @REQ-CORE-014 @spec(api_core.md#10-file-validation-requirements)
Feature: File validation requirements

  @happy
  Scenario: File path validation passes for valid paths
    Given a file path
    When path validation is performed
    Then path is not empty
    And path does not contain only whitespace
    And path is normalized correctly

  @error
  Scenario: File path validation rejects empty paths
    Given an empty file path
    When path validation is performed
    Then validation fails
    And structured validation error is returned

  @error
  Scenario: File path validation rejects whitespace-only paths
    Given a file path containing only whitespace
    When path validation is performed
    Then validation fails
    And structured validation error is returned

  @happy
  Scenario: File path normalization removes redundant separators
    Given a file path with redundant separators
    When path is normalized
    Then redundant separators are removed
    And path is consistent

  @happy
  Scenario: File path normalization resolves relative references
    Given a file path with relative references
    When path is normalized
    Then relative references are resolved
    And path is absolute

  @happy
  Scenario: File data validation allows empty files
    Given a file with empty data (len = 0)
    When data validation is performed
    Then validation passes
    And empty file is accepted

  @error
  Scenario: File data validation rejects nil data
    Given nil file data
    When data validation is performed
    Then validation fails
    And structured validation error is returned

  @error
  Scenario: Invalid files are rejected with appropriate error messages
    Given an invalid file
    When file validation is performed
    Then validation fails
    And structured validation error is returned
    And error message describes the issue
