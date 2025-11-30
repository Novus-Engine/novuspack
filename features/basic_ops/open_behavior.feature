@domain:basic_ops @m2 @REQ-API_BASIC-041 @spec(api_basic_operations.md#512-open-behavior)
Feature: Open Behavior

  @REQ-API_BASIC-041 @happy
  Scenario: Open validates package file exists and is readable
    Given a Package instance
    And a path to existing package file
    When Open is called
    Then package file existence is validated
    And file readability is checked
    And file format is verified

  @REQ-API_BASIC-041 @happy
  Scenario: Open loads package header and metadata
    Given a valid package file path
    When Open is called
    Then package header is loaded
    And package metadata is loaded
    And file entries are indexed

  @REQ-API_BASIC-041 @happy
  Scenario: Open reads package structure into memory
    Given a valid package file
    When Open is called
    Then package structure is read into memory
    And package is ready for operations
    And package state reflects opened status

  @REQ-API_BASIC-041 @happy
  Scenario: Open supports read-only mode
    Given a package file to open
    When Open is called
    Then package can be opened in read-only mode
    And read-only flag is set appropriately
    And write operations are prevented in read-only mode

  @REQ-API_BASIC-041 @error
  Scenario: Open validates package file format
    Given a Package instance
    And a file path
    When Open is called
    Then package format is validated
    And invalid format is rejected
    And validation error is returned for invalid format
