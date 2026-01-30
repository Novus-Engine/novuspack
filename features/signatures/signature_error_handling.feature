@skip @domain:signatures @m2 @spec(api_signatures.md#5-error-handling)
Feature: Signature Error Handling

# This feature captures signature error handling expectations from the signatures API specification.
# Detailed runnable scenarios live in the dedicated signatures feature files.

  @REQ-SIG-050 @error
  Scenario: Signature operations use structured errors with signature-specific error types
    Given a signature operation fails
    When an error is returned
    Then the error is a structured error
    And the error type reflects signature, validation, security, unsupported, or corruption failures

  @REQ-SIG-052 @error
  Scenario: Signature errors include typed context describing the signature and operation
    Given a signature validation fails for a specific signature index
    When the structured error is created
    Then the error context includes the signature index and the operation name
    And callers can use that context to report actionable diagnostics
