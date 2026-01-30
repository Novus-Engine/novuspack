@domain:signatures @m2 @v2 @REQ-SIG-055 @spec(api_signatures.md#53-structured-error-examples)
Feature: Signature Structured Error Examples

  @REQ-SIG-055 @happy
  Scenario: Structured error examples demonstrate error creation patterns
    Given a NovusPack package
    When signature errors are created
    Then errors use NewTypedPackageError with error type
    And errors include typed context structures
    And errors follow structured error creation pattern
    And errors provide detailed context information

  @REQ-SIG-055 @happy
  Scenario: Signature validation error includes typed context
    Given a NovusPack package
    When signature validation fails
    Then error uses ErrTypeSignature error type
    And error uses ErrSignatureValidationFailed error message
    And error includes SignatureErrorContext with signature index
    And error includes SignatureErrorContext with algorithm name
    And error includes SignatureErrorContext with operation name

  @REQ-SIG-055 @happy
  Scenario: Unsupported signature type error includes typed context
    Given a NovusPack package
    When unsupported signature type is used
    Then error uses ErrTypeUnsupported error type
    And error uses ErrUnsupportedSignatureType error message
    And error includes UnsupportedErrorContext with signature type
    And error includes UnsupportedErrorContext with supported types list
    And error includes UnsupportedErrorContext with operation name

  @REQ-SIG-055 @happy
  Scenario: Key not found error includes typed context
    Given a NovusPack package
    When signing key is not found
    Then error uses ErrTypeSecurity error type
    And error uses ErrKeyNotFound error message
    And error includes SecurityErrorContext with key ID
    And error includes SecurityErrorContext with key type
    And error includes SecurityErrorContext with operation name

  @REQ-SIG-055 @error
  Scenario: Error examples handle error creation correctly
    Given a NovusPack package
    When error is created with invalid context
    Then error still follows structured error format
    And error provides meaningful error information
