@domain:basic_ops @m2 @REQ-API_BASIC-064 @spec(api_basic_operations.md#81-structured-error-system)
Feature: Basic Operations: Structured Error System

  @REQ-API_BASIC-064 @happy
  Scenario: Structured error system provides error categorization
    Given a NovusPack package operation
    When an error occurs during operation
    Then structured error is returned with error type
    And error is categorized (validation, I/O, security, etc.)
    And error category enables appropriate handling

  @REQ-API_BASIC-064 @happy
  Scenario: Structured errors include context for debugging
    Given a package operation that fails
    When structured error is returned
    Then error contains operation context
    And error includes relevant parameters
    And error provides debugging information

  @REQ-API_BASIC-064 @happy
  Scenario: Structured errors support error inspection
    Given a structured error from package operation
    When error is inspected
    Then error type can be checked
    And error message can be retrieved
    And error context can be accessed
    And error enables programmatic handling

  @REQ-API_BASIC-064 @happy
  Scenario: Structured errors provide better error messages
    Given package operations that may fail
    When structured errors are returned
    Then error messages are descriptive
    And error messages include relevant details
    And error messages aid in debugging
