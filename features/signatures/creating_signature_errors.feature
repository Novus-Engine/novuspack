@domain:signatures @m2 @v2 @REQ-SIG-056 @spec(api_signatures.md#531-creating-signature-errors)
Feature: Creating Signature Errors

  @REQ-SIG-056 @happy
  Scenario: Signature errors are created using NewTypedPackageError
    Given a NovusPack package
    When signature error is created
    Then NewTypedPackageError is used with error type
    And NewTypedPackageError is used with error message
    And NewTypedPackageError is used with typed context
    And error follows structured error creation pattern

  @REQ-SIG-056 @happy
  Scenario: Signature validation error creation includes context
    Given a NovusPack package
    When signature validation error is created
    Then error uses ErrTypeSignature error type
    And error uses ErrSignatureValidationFailed error message
    And error includes SignatureErrorContext with signature index
    And error includes SignatureErrorContext with algorithm name
    And error includes SignatureErrorContext with operation name

  @REQ-SIG-056 @happy
  Scenario: Unsupported signature type error creation includes context
    Given a NovusPack package
    When unsupported signature type error is created
    Then error uses ErrTypeUnsupported error type
    And error uses ErrUnsupportedSignatureType error message
    And error includes UnsupportedErrorContext with signature type
    And error includes UnsupportedErrorContext with supported types
    And error includes UnsupportedErrorContext with operation name

  @REQ-SIG-056 @happy
  Scenario: Key not found error creation includes context
    Given a NovusPack package
    When key not found error is created
    Then error uses ErrTypeSecurity error type
    And error uses ErrKeyNotFound error message
    And error includes SecurityErrorContext with key ID
    And error includes SecurityErrorContext with key type
    And error includes SecurityErrorContext with operation name

  @REQ-SIG-056 @error
  Scenario: Error creation handles invalid context gracefully
    Given a NovusPack package
    When error is created with invalid context
    Then error still follows structured error format
    And error provides meaningful error information
