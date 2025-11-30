@domain:signatures @m2 @REQ-SIG-049 @spec(api_signatures.md#43-generic-signature-validation)
Feature: Generic Signature Validation

  @REQ-SIG-049 @happy
  Scenario: SignatureValidator provides type-safe signature validation
    Given a NovusPack package
    When SignatureValidator structure is examined
    Then structure extends Validator for type-safe validation
    And structure contains signatureRules field with validation rules
    And structure provides AddSignatureRule method for adding rules
    And structure provides ValidateSignatureData method for data validation
    And structure provides ValidateSignatureKey method for key validation
    And structure provides ValidateSignatureFormat method for format validation

  @REQ-SIG-049 @happy
  Scenario: SignatureValidationRule enables custom validation rules
    Given a NovusPack package
    When SignatureValidationRule is used
    Then rule is alias for generic ValidationRule
    And rule enables type-safe signature validation
    And rule supports custom validation logic
    And rule integrates with generic validation system

  @REQ-SIG-049 @happy
  Scenario: Generic signature validation enables comprehensive validation
    Given a NovusPack package
    When signature validation is performed
    Then ValidateSignatureData validates signature data
    And ValidateSignatureKey validates signature key
    And ValidateSignatureFormat validates signature format
    And validation rules are applied in sequence
    And validation provides type-safe error information

  @REQ-SIG-049 @error
  Scenario: Generic signature validation handles validation failures
    Given a NovusPack package
    When signature validation fails
    Then structured validation error is returned
    And error indicates which validation rule failed
    And error provides context for validation failure
    And error follows structured error format
