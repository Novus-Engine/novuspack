@domain:security @m2 @REQ-SEC-075 @spec(api_security.md#11-multiple-signature-validation-incremental)
Feature: Multiple Signature Validation

  @REQ-SEC-075 @happy
  Scenario: Multiple signature validation validates all signatures incrementally
    Given an open NovusPack package
    And a valid context
    And package with multiple signatures
    When ValidateAllSignatures is called
    Then all signatures are validated in order
    And each signature validates complete package state up to its position
    And SignatureValidationResult is returned for each signature
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-075 @happy
  Scenario: Multiple signature validation supports signature type filtering
    Given an open NovusPack package
    And a valid context
    And package with multiple signature types
    When ValidateSignatureType is called with specific type
    Then only signatures of specified type are validated
    And SignatureValidationResult array is returned
    And results contain validation status for each matching signature

  @REQ-SEC-075 @happy
  Scenario: Multiple signature validation supports signature index validation
    Given an open NovusPack package
    And a valid context
    And package with multiple signatures
    When ValidateSignatureIndex is called with valid index
    Then signature at specified index is validated
    And SignatureValidationResult is returned
    And result contains index, type, validity, and trust status

  @REQ-SEC-075 @happy
  Scenario: Multiple signature validation validates signature chain integrity
    Given an open NovusPack package
    And a valid context
    And package with signature chain
    When ValidateSignatureChain is called
    Then chain validation ensures no signatures were removed
    And chain validation ensures no signatures were modified
    And SignatureValidationResult array is returned for entire chain

  @REQ-SEC-075 @error
  Scenario: Multiple signature validation handles invalid signature index
    Given an open NovusPack package
    And a valid context
    And package with limited signatures
    When ValidateSignatureIndex is called with out-of-range index
    Then index out of range error is returned
    And error indicates invalid signature index
    And error follows structured error format

  @REQ-SEC-075 @happy
  Scenario: Multiple signature validation provides incremental validation process
    Given an open NovusPack package
    And a valid context
    And package with sequentially added signatures
    When incremental validation is performed
    Then first signature validates initial package state
    And second signature validates package state after first signature
    And each subsequent signature validates complete state up to its position
    And validation process ensures signature chain integrity
