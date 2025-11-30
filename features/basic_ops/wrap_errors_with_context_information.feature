@domain:basic_ops @m2 @REQ-API_BASIC-083 @spec(api_basic_operations.md#921-wrap-errors-with-context)
Feature: Wrap Errors with Context

  @REQ-API_BASIC-083 @happy
  Scenario: Errors are wrapped with additional context information
    Given a package operation that returns an error
    When error is wrapped with context
    Then error includes additional context information
    And context provides better debugging information
    And error message is enhanced with details

  @REQ-API_BASIC-083 @happy
  Scenario: Error context includes relevant details like file paths
    Given a package operation involving file paths
    And operation returns an error
    When error is wrapped with context
    Then file paths are included in error context
    And context helps identify source of problem
    And error message provides better information

  @REQ-API_BASIC-083 @happy
  Scenario: Error context includes operation names
    Given a package operation with a specific name
    And operation returns an error
    When error is wrapped with context
    Then operation name is included in error context
    And context helps identify which operation failed
    And error message indicates operation context

  @REQ-API_BASIC-083 @happy
  Scenario: Error context includes parameter values
    Given a package operation with parameters
    And operation returns an error
    When error is wrapped with context
    Then parameter values are included in error context
    And context helps identify problematic parameters
    And error message provides parameter details

  @REQ-API_BASIC-083 @happy
  Scenario: Error wrapping uses structured error system
    Given a package operation that returns an error
    When error is wrapped with context
    Then structured error system is used
    And error follows PackageError structure
    And WithContext method is used to add context

  @REQ-API_BASIC-083 @happy
  Scenario: Wrapped errors provide better error messages to users
    Given a package operation that returns an error
    When error is wrapped with context
    Then error message is user-friendly
    And error message includes helpful context
    And error message is more informative than unwrapped error

  @REQ-API_BASIC-083 @error
  Scenario: Error wrapping preserves original error information
    Given a package operation that returns an error
    When error is wrapped with context
    Then original error is preserved as cause
    And error unwrapping still works
    And error inspection can access original error
