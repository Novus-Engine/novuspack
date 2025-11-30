@domain:signatures @security @m2 @REQ-SIG-010 @spec(api_signatures.md#27-signature-validation)
Feature: Signature validation methods

  @happy
  Scenario: ValidateAllSignatures validates all signatures in order
    Given a package with multiple signatures
    When ValidateAllSignatures is called
    Then all signatures are validated in order
    And validation results are returned for each signature
    And validation status indicates success or failure

  @REQ-SIG-011 @happy
  Scenario: ValidateSignatureType validates signatures of specific type
    Given a package with signatures of different types
    When ValidateSignatureType is called with signature type
    Then only signatures of that type are validated
    And validation results are returned
    And other signature types are not validated

  @REQ-SIG-012 @happy
  Scenario: ValidateSignatureIndex validates signature by index
    Given a package with multiple signatures
    When ValidateSignatureIndex is called with index
    Then signature at that index is validated
    And validation result is returned
    And other signatures are not validated

  @REQ-SIG-013 @happy
  Scenario: ValidateSignatureWithKey validates signature with specific key
    Given a package with signature
    When ValidateSignatureWithKey is called with public key
    Then signature is validated with that key
    And validation result indicates if key matches
    And validation fails if key does not match

  @REQ-SIG-014 @happy
  Scenario: ValidateSignatureChain validates signature chain integrity
    Given a package with signature chain
    When ValidateSignatureChain is called
    Then signature chain is validated
    And chain integrity is verified
    And validation result indicates chain validity

  @error
  Scenario: ValidateSignatureIndex fails with invalid index
    Given a package with signatures
    When ValidateSignatureIndex is called with invalid index
    Then structured validation error is returned

  @error
  Scenario: ValidateSignatureWithKey fails with invalid key
    Given a package with signature
    When ValidateSignatureWithKey is called with invalid key
    Then structured validation error is returned

  @REQ-SIG-015 @REQ-SIG-017 @error
  Scenario: ValidateSignatureIndex validates index parameter
    Given a package with signatures
    When ValidateSignatureIndex is called with negative index
    Then structured validation error is returned
    And error indicates invalid index

  @REQ-SIG-015 @REQ-SIG-018 @error
  Scenario: ValidateSignatureWithKey validates public key parameter
    Given a package with signature
    When ValidateSignatureWithKey is called with nil key
    Then structured validation error is returned
    And error indicates invalid key

  @REQ-SIG-015 @REQ-SIG-019 @error
  Scenario: Signature validation methods respect context cancellation
    Given a package with signatures
    And a cancelled context
    When signature validation method is called
    Then structured context error is returned
    And error type is context cancellation
