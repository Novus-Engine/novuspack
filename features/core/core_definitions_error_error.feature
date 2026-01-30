@skip @domain:core @m2 @REQ-CORE-025 @spec(api_core.md#1043-error-propagation)
Feature: Core Definitions

# This feature captures core error propagation expectations from the core API specification.
# Detailed runnable scenarios for structured errors are covered in dedicated core error feature files.

  @REQ-CORE-025 @error
  Scenario: Operations propagate failures as structured errors with typed context
    Given an operation fails due to invalid input
    When the implementation returns the error to the caller
    Then the error is returned as a structured error
    And the structured error includes a descriptive message and typed context for the failing operation

  @REQ-CORE-025 @error
  Scenario: Wrapped errors preserve the original error and add domain context
    Given an underlying I/O error occurs during a package operation
    When the error is wrapped for propagation
    Then the returned error preserves the original error as the cause
    And the returned error adds domain-specific context such as the virtual path being processed
