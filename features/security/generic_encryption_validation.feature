@domain:security @m2 @REQ-SEC-089 @spec(api_security.md#43-generic-encryption-validation)
Feature: Generic Encryption Validation

  @REQ-SEC-089 @happy
  Scenario: EncryptionValidator provides type-safe encryption validation
    Given an open NovusPack package
    And encryption validation requirements
    When EncryptionValidator is created
    Then type-safe encryption validation is provided
    And validator extends Validator base type
    And validator is generic over data type

  @REQ-SEC-089 @happy
  Scenario: EncryptionValidator supports encryption validation rules
    Given an open NovusPack package
    And encryption validation requirements
    When EncryptionValidator is configured
    Then encryption validation rules can be added
    And EncryptionValidationRule is alias for ValidationRule
    And rules provide type-safe validation

  @REQ-SEC-089 @happy
  Scenario: EncryptionValidator validates encryption data
    Given an open NovusPack package
    And EncryptionValidator with rules
    And encryption data
    And a valid context
    When ValidateEncryptionData is called
    Then encryption data is validated
    And validation rules are applied
    And validation result is returned

  @REQ-SEC-089 @happy
  Scenario: EncryptionValidator validates decryption data
    Given an open NovusPack package
    And EncryptionValidator with rules
    And decryption data
    And a valid context
    When ValidateDecryptionData is called
    Then decryption data is validated
    And validation rules are applied
    And validation result is returned

  @REQ-SEC-089 @happy
  Scenario: EncryptionValidator validates encryption keys
    Given an open NovusPack package
    And EncryptionValidator with rules
    And encryption key
    And a valid context
    When ValidateEncryptionKey is called
    Then encryption key is validated
    And validation rules are applied
    And validation result is returned

  @REQ-SEC-089 @happy
  Scenario: EncryptionValidator supports adding encryption rules
    Given an open NovusPack package
    And EncryptionValidator
    And encryption validation rule
    When AddEncryptionRule is called
    Then encryption rule is added to validator
    And rule is included in validation process
    And type-safe rule configuration is provided

  @REQ-SEC-011 @error
  Scenario: EncryptionValidator operations respect context cancellation
    Given an open NovusPack package
    And EncryptionValidator
    And a cancelled context
    When encryption validation operation is called
    Then structured context error is returned
    And error type is context cancellation
