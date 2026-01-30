@skip @domain:basic_ops @m2 @spec(api_basic_operations.md#81-structured-error-system)
Feature: Basic Operations Error Handling

# This feature captures basic operations error handling expectations and best-practice guidance.
# Detailed runnable scenarios live in the dedicated basic_ops and core error handling feature files.

  @REQ-API_BASIC-072 @error
  Scenario: Callers must always check for returned errors
    Given a package operation that returns an error
    When the caller receives a non-nil error
    Then the caller must treat that error as a failure
    And the caller must not continue as if the operation succeeded

  @REQ-API_BASIC-073 @error
  Scenario: Structured errors include actionable context for debugging
    Given a package operation fails due to invalid input
    When the operation returns an error
    Then the error is a structured error
    And the error includes enough context to identify the failing parameter or operation
