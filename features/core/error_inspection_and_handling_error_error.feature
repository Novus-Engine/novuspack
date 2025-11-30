@domain:core @m2 @REQ-CORE-024 @spec(api_core.md#1142-error-inspection-and-handling)
Feature: Error Inspection and Handling

  @REQ-CORE-024 @happy
  Scenario: Error inspection checks if error is a PackageError
    Given an error
    When error inspection is performed
    Then IsPackageError function checks if error is PackageError
    And PackageError can be extracted for inspection
    And error type and context can be examined

  @REQ-CORE-024 @happy
  Scenario: Error inspection enables type-based error handling
    Given an error
    When error inspection is performed
    Then error type can be checked using GetErrorType
    And different error types can be handled appropriately
    And switch statements enable type-based handling

  @REQ-CORE-024 @happy
  Scenario: Error inspection allows context extraction
    Given a PackageError with context
    When error inspection is performed
    Then context information can be extracted
    And context fields can be accessed
    And context provides debugging information

  @REQ-CORE-024 @happy
  Scenario: Error handling uses inspected error information
    Given an error
    And error inspection is performed
    When error handling is performed
    Then error type determines handling strategy
    And context information guides recovery
    And appropriate error response is selected
