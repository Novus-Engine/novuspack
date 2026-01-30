@skip @domain:core @m2 @spec(api_core.md#102-packageerror-structure)
Feature: Core Error Handling

# This feature captures the core PackageError structure and helper expectations.
# Detailed runnable scenarios for structured errors are covered in dedicated core error feature files.

  @REQ-CORE-020 @error
  Scenario: PackageError provides a typed category and a human-readable message
    Given an operation fails with a structured error
    When the returned error is inspected as a PackageError
    Then the error has an ErrorType category
    And the error has a human-readable Message describing the failure

  @REQ-CORE-020 @REQ-CORE-021 @error
  Scenario: PackageError supports unwrapping and error matching via the cause chain
    Given a PackageError with an underlying cause error
    When a caller unwraps the error
    Then the underlying cause error is available for inspection
    And callers can match errors using standard error matching semantics
