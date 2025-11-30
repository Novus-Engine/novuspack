@domain:basic_ops @m2 @REQ-API_BASIC-058 @spec(api_basic_operations.md#731-getinfo-error-conditions)
Feature: GetInfo error conditions

  @REQ-API_BASIC-058 @error
  Scenario: GetInfo returns validation error when package is not open
    Given a NovusPack package that is not open
    When GetInfo is called
    Then validation error is returned
    And error indicates package is not currently open
    And package information cannot be retrieved

  @REQ-API_BASIC-058 @error
  Scenario: GetInfo returns context error on cancellation
    Given an open NovusPack package
    And a cancelled context
    When GetInfo is called with cancelled context
    Then context error is returned
    And error type is context cancellation

  @REQ-API_BASIC-058 @error
  Scenario: GetInfo returns context error on timeout
    Given an open NovusPack package
    And a context with expired timeout
    When GetInfo is called
    Then context error is returned
    And error type is context timeout
