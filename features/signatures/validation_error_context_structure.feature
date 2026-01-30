@domain:signatures @m2 @v2 @REQ-SIG-077 @spec(api_signatures.md#534-validationerrorcontext-structure)
Feature: Signature validation error context structure

  @REQ-SIG-077 @happy
  Scenario: ValidationErrorContext provides structured fields for signature validation failures
    Given a signature validation failure
    When returning a structured error context
    Then ValidationErrorContext contains fields describing the validation failure
    And the context supports diagnosis of signature validation errors
    And the context includes enough information to identify which signature failed validation

  @REQ-SIG-077 @error
  Scenario: Validation errors include ValidationErrorContext when verification fails
    Given a signature verification operation fails
    When a structured error is returned
    Then the error includes ValidationErrorContext
    And the error type indicates signature validation failure

