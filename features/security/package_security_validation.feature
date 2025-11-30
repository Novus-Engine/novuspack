@domain:security @m2 @REQ-SEC-030 @spec(security.md#41-package-validation)
Feature: Package Security Validation

  @REQ-SEC-030 @happy
  Scenario: Package validation validates package integrity
    Given an open NovusPack package
    And a valid context
    And package requiring validation
    When Validate is called
    Then package format is validated
    And package structure is validated
    And package integrity is verified
    And context supports cancellation
    And context supports timeout handling

  @REQ-SEC-030 @happy
  Scenario: Package validation validates signatures
    Given an open NovusPack package
    And a valid context
    And package with signatures
    When Validate is called
    Then signatures are validated
    And signature integrity is verified
    And signature chain is validated

  @REQ-SEC-030 @happy
  Scenario: Package validation validates encryption
    Given an open NovusPack package
    And a valid context
    And package with encrypted files
    When Validate is called
    Then encryption integrity is validated
    And encryption keys are validated
    And encrypted file integrity is verified

  @REQ-SEC-030 @happy
  Scenario: Package validation validates checksums
    Given an open NovusPack package
    And a valid context
    And package with checksums
    When ValidateIntegrity is called
    Then package integrity is validated using checksums
    And file checksums are verified
    And checksum integrity is maintained

  @REQ-SEC-030 @happy
  Scenario: Package validation provides comprehensive validation
    Given an open NovusPack package
    And a valid context
    And package requiring comprehensive validation
    When Validate is called
    Then validation covers signatures, encryption, and checksums
    And all security aspects are validated
    And validation provides complete security assessment

  @REQ-SEC-030 @error
  Scenario: Package validation handles validation failures
    Given an open NovusPack package
    And a valid context
    And package with validation issues
    When Validate is called
    Then validation errors are returned
    And error indicates specific validation failure
    And error follows structured error format
