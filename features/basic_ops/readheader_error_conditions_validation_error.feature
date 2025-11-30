@domain:basic_ops @m2 @REQ-API_BASIC-062 @spec(api_basic_operations.md#743-readheader-error-conditions)
Feature: ReadHeader Error Conditions

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns validation error for invalid header format
    Given a valid context
    And a reader with invalid package header format
    When ReadHeader is called
    Then validation error is returned
    And error indicates invalid package header format
    And error follows structured error format

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns unsupported error for unsupported version
    Given a valid context
    And a reader with unsupported package version
    When ReadHeader is called
    Then unsupported error is returned
    And error indicates package version not supported
    And error follows structured error format

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns context cancellation error
    Given a cancelled context
    And a reader for package header
    When ReadHeader is called
    Then context cancellation error is returned
    And error type is context cancellation
    And error follows structured error format

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns context timeout error
    Given a context with timeout
    And a reader that exceeds timeout duration
    When ReadHeader is called
    Then context timeout error is returned
    And error type is context timeout
    And error follows structured error format

  @REQ-API_BASIC-062 @error
  Scenario: ReadHeader returns error for truncated header data
    Given a valid context
    And a reader with insufficient header data
    When ReadHeader is called
    Then validation error is returned
    And error indicates header data is incomplete
    And error follows structured error format
