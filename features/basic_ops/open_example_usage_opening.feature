@domain:basic_ops @m2 @REQ-API_BASIC-040 @REQ-API_BASIC-041 @spec(api_basic_operations.md#514-openpackage-example-usage)
Feature: OpenPackage Example Usage

  @REQ-API_BASIC-040 @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage example demonstrates basic package opening
    Given a valid context
    And an existing package file at "/path/to/existing-package.nvpk"
    When OpenPackage is called
    Then package is opened successfully
    And package is ready for operations
    And no error is returned

  @REQ-API_BASIC-040 @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage example includes error handling pattern
    Given a valid context
    And a package file path
    When OpenPackage is called with error checking
    Then error result is checked
    And error handling follows standard Go pattern
    And function returns early on error

  @REQ-API_BASIC-040 @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage example follows standard Go error handling
    Given a code example demonstrating OpenPackage usage
    When example code is examined
    Then OpenPackage call includes error check
    And error is handled with if err != nil pattern
    And function returns error on failure

  @REQ-API_BASIC-040 @REQ-API_BASIC-041 @happy
  Scenario: OpenPackage example uses context parameter correctly
    Given a code example demonstrating OpenPackage usage
    When example code is examined
    Then context is passed as first parameter
    And context is used from calling function
    And context supports standard Go patterns

  @REQ-API_BASIC-040 @REQ-API_BASIC-041 @error
  Scenario: OpenPackage example handles file not found errors
    Given a valid context
    And a non-existent package file path
    When OpenPackage is called
    Then file not found error is returned
    And error indicates package file does not exist
    And error follows structured error format

  @REQ-API_BASIC-040 @REQ-API_BASIC-041 @error
  Scenario: OpenPackage example handles invalid package format errors
    Given a valid context
    And a file with invalid package format
    When OpenPackage is called
    Then invalid format error is returned
    And error indicates format validation failure
    And error follows structured error format
