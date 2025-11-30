@domain:basic_ops @m2 @REQ-API_BASIC-063 @spec(api_basic_operations.md#744-readheader-example-usage)
Feature: ReadHeader Example Usage

  @REQ-API_BASIC-063 @happy
  Scenario: ReadHeader example demonstrates basic header reading
    Given a valid context
    And a file or reader with package header
    When ReadHeader is called with context and reader
    Then header is read successfully
    And Header structure is returned
    And no error is returned

  @REQ-API_BASIC-063 @happy
  Scenario: ReadHeader example includes error handling pattern
    Given a valid context
    And a reader for package header
    When ReadHeader is called with error checking
    Then error result is checked
    And error handling follows standard Go pattern
    And function returns early on error

  @REQ-API_BASIC-063 @happy
  Scenario: ReadHeader example follows standard Go error handling
    Given a code example demonstrating ReadHeader usage
    When example code is examined
    Then ReadHeader call includes error check
    And error is handled with if err != nil pattern
    And function returns error on failure

  @REQ-API_BASIC-063 @happy
  Scenario: ReadHeader example uses context parameter correctly
    Given a code example demonstrating ReadHeader usage
    When example code is examined
    Then context is passed as first parameter
    And context is used from calling function
    And context supports standard Go patterns

  @REQ-API_BASIC-063 @happy
  Scenario: ReadHeader example supports header-only inspection
    Given a valid context
    And a package file
    When ReadHeader is used for header-only inspection
    Then package header is read without opening full package
    And header metadata is accessible
    And package data is not loaded

  @REQ-API_BASIC-063 @error
  Scenario: ReadHeader example handles invalid header format errors
    Given a valid context
    And a reader with invalid header format
    When ReadHeader is called
    Then validation error is returned
    And error indicates invalid header format
    And error follows structured error format

  @REQ-API_BASIC-063 @error
  Scenario: ReadHeader example handles unsupported version errors
    Given a valid context
    And a reader with unsupported package version
    When ReadHeader is called
    Then unsupported error is returned
    And error indicates version not supported
    And error follows structured error format
