@domain:basic_ops @m2 @REQ-API_BASIC-040 @spec(api_basic_operations.md#511-open-parameters)
Feature: Open Parameter Specification

  @REQ-API_BASIC-040 @happy
  Scenario: Open accepts context parameter for cancellation and timeout
    Given a NovusPack package instance
    And a valid context
    When Open is called with context parameter
    Then context is accepted as first parameter
    And context supports cancellation handling
    And context supports timeout handling

  @REQ-API_BASIC-040 @happy
  Scenario: Open accepts file path parameter
    Given a NovusPack package instance
    And a valid context
    And a valid package file path
    When Open is called with path parameter
    Then path parameter is accepted as second parameter
    And path specifies file system location of package file

  @REQ-API_BASIC-040 @happy
  Scenario: Open method signature matches specification
    Given the Open method definition
    When method signature is examined
    Then method accepts ctx context.Context as first parameter
    And method accepts path string as second parameter
    And method returns error

  @REQ-API_BASIC-040 @error
  Scenario: Open validates path parameter
    Given a NovusPack package instance
    And a valid context
    And an invalid or empty path
    When Open is called with invalid path
    Then validation error is returned
    And error indicates path format issue

  @REQ-API_BASIC-040 @error
  Scenario: Open handles context cancellation
    Given a NovusPack package instance
    And a cancelled context
    When Open is called with cancelled context
    Then context cancellation error is returned
    And error type is context cancellation

  @REQ-API_BASIC-040 @error
  Scenario: Open handles context timeout
    Given a NovusPack package instance
    And a context with timeout
    And operation exceeds timeout duration
    When Open is called
    Then context timeout error is returned
    And error type is context timeout
