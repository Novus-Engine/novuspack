@domain:basic_ops @m2 @REQ-API_BASIC-069 @spec(api_basic_operations.md#832-error-inspection)
Feature: Error inspection

  @REQ-API_BASIC-069 @happy
  Scenario: Error inspection checks error types
    Given a package error
    When error is inspected
    Then error type can be checked
    And error type determines handling strategy
    And error inspection enables programmatic handling

  @REQ-API_BASIC-069 @happy
  Scenario: Error inspection retrieves error messages
    Given a structured package error
    When error is inspected
    Then error message can be retrieved
    And message provides human-readable description
    And message aids in debugging

  @REQ-API_BASIC-069 @happy
  Scenario: Error inspection accesses error context
    Given a structured package error with context
    When error is inspected
    Then error context can be accessed
    And context provides additional debugging information
    And context includes operation details

  @REQ-API_BASIC-069 @happy
  Scenario: Error inspection enables appropriate logging
    Given package errors from different operations
    When errors are inspected
    Then error types guide logging strategy
    And error context is included in logs
    And logging provides comprehensive error information
