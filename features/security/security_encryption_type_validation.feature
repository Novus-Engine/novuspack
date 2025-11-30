@domain:security @m2 @REQ-SEC-004 @spec(api_security.md#32-encryption-type-validation)
Feature: Encryption type validation

  @happy
  Scenario: IsValidEncryptionType validates encryption type values
    Given an encryption type value
    When IsValidEncryptionType is called with type
    Then true is returned if type is valid
    And false is returned if type is invalid
    And validation matches supported types

  @happy
  Scenario: GetEncryptionTypeName returns human-readable name
    Given a valid encryption type
    When GetEncryptionTypeName is called with type
    Then human-readable encryption type name is returned
    And name is descriptive
    And name matches encryption type

  @error
  Scenario: GetEncryptionTypeName fails with invalid type
    Given an invalid encryption type
    When GetEncryptionTypeName is called with type
    Then structured validation error is returned
    And error indicates invalid encryption type
