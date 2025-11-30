@domain:core @m2 @REQ-CORE-027 @spec(api_core.md#1145-error-logging-and-debugging)
Feature: Error Logging and Debugging

  @REQ-CORE-027 @happy
  Scenario: Structured errors provide comprehensive logging information
    Given a PackageError
    And an operation name
    When error logging is performed
    Then error type is logged
    And error message is logged
    And error context is logged
    And cause error is logged if present

  @REQ-CORE-027 @happy
  Scenario: Error logging includes full context for debugging
    Given a PackageError with context
    And an operation name
    When error logging is performed
    Then all context fields are included in log
    And context provides debugging details
    And logging supports troubleshooting

  @REQ-CORE-027 @happy
  Scenario: Error logging supports debugging workflows
    Given a PackageError
    When error logging is performed
    Then logged information enables debugging
    And error details help identify problem source
    And logging supports error diagnosis

  @REQ-CORE-027 @happy
  Scenario: Structured errors provide more information for logs than standard errors
    Given a PackageError
    And a standard error
    When error logging is performed
    Then structured errors provide more information
    And structured errors include type, message, context, and cause
    And logging is more comprehensive
