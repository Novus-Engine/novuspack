@skip @domain:file_mgmt @m2 @spec(api_file_mgmt_errors.md#11-error-type-categories)
Feature: File Management Error Handling

# This feature captures file management error handling expectations from the file management error spec.
# Detailed runnable scenarios live in the dedicated file_mgmt error handling feature files.

  @REQ-FILEMGMT-110 @error
  Scenario: File management operations return structured errors
    Given a file management operation fails
    When the operation returns an error
    Then the error is a structured error
    And the error category reflects the failure type (validation, IO, encryption, compression, security, context)

  @REQ-FILEMGMT-111 @error
  Scenario: Structured errors include typed context for the failing file operation
    Given an operation fails while processing a specific path
    When the structured error is created
    Then the error includes typed context that includes the path and operation name
    And callers can inspect that context to diagnose the failure
