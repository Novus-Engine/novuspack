@domain:basic_ops @m2 @REQ-API_BASIC-041 @spec(api_basic_operations.md#512-openpackage-behavior)
Feature: OpenPackage Behavior

  @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage validates package file exists and is readable
    Given a Package instance
    And a path to existing package file
    When OpenPackage is called
    Then package file existence is validated
    And file readability is checked
    And file format is verified

  @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage loads package header and metadata
    Given a valid package file path
    When OpenPackage is called
    Then package header is loaded
    And package metadata is loaded
    And file entries are indexed

  @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage reads package structure into memory
    Given a valid package file
    When OpenPackage is called
    Then package structure is read into memory
    And package is ready for operations
    And package state reflects opened status

  @REQ-API_BASIC-041 @error
  Scenario: OpenPackage validates package file format
    Given a Package instance
    And a file path
    When OpenPackage is called
    Then package format is validated
    And invalid format is rejected
    And validation error is returned for invalid format
