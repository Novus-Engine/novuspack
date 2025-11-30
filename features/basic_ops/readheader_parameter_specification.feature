@domain:basic_ops @m2 @REQ-API_BASIC-061 @spec(api_basic_operations.md#742-readheader-parameters)
Feature: ReadHeader Parameter Specification

  @REQ-API_BASIC-061 @happy
  Scenario: ReadHeader accepts context parameter
    Given a valid context
    And a reader for package header
    When ReadHeader is called
    Then context is accepted as first parameter
    And context supports cancellation handling
    And context supports timeout handling

  @REQ-API_BASIC-061 @happy
  Scenario: ReadHeader accepts reader parameter
    Given a valid context
    And an io.Reader for input stream
    When ReadHeader is called with reader
    Then reader is accepted as second parameter
    And header is read from input stream
    And reader provides header data

  @REQ-API_BASIC-061 @happy
  Scenario: ReadHeader returns Header structure
    Given a valid context
    And a reader with valid package header
    When ReadHeader is called
    Then Header structure is returned
    And Header contains package metadata
    And Header follows NovusPack header format

  @REQ-API_BASIC-061 @happy
  Scenario: ReadHeader returns error on failure
    Given a valid context
    And a reader that cannot provide header
    When ReadHeader is called
    Then error is returned
    And error follows structured error format
    And error provides details about failure

  @REQ-API_BASIC-061 @error
  Scenario: ReadHeader validates context is not cancelled
    Given a cancelled context
    And a reader for package header
    When ReadHeader is called
    Then context cancellation error is returned
    And error type is context cancellation

  @REQ-API_BASIC-061 @error
  Scenario: ReadHeader validates context timeout
    Given a context with timeout
    And a reader that exceeds timeout
    When ReadHeader is called
    Then context timeout error is returned
    And error type is context timeout
