@domain:signatures @m2 @REQ-SIG-035 @spec(api_signatures.md#282-function-usage-guide)
Feature: Signature Usage Examples

  @REQ-SIG-035 @happy
  Scenario: Function usage guide provides signature function guidance
    Given a NovusPack package
    And a valid context
    When function usage guide is examined
    Then guide explains when to use AddSignature (Low-Level)
    And guide explains when to use SignPackage* functions (High-Level)
    And guide explains implementation pattern
    And context supports cancellation

  @REQ-SIG-035 @happy
  Scenario: Usage guide explains when to use AddSignature
    Given a NovusPack package
    And a valid context
    When AddSignature usage is examined
    Then use AddSignature when you have pre-computed signature data
    And use AddSignature when you want direct control over signature addition
    And use AddSignature when implementing custom signature generation logic
    And use AddSignature when adding signatures from external signature services

  @REQ-SIG-035 @happy
  Scenario: Usage guide explains when to use SignPackage* functions
    Given a NovusPack package
    And a valid context
    When SignPackage* usage is examined
    Then use SignPackage* when you have private key and want to generate + add signature
    And use SignPackage* when you want convenience of automatic signature generation
    And use SignPackage* when using standard signature types (ML-DSA, SLH-DSA, PGP, X.509)
    And use SignPackage* when you want automatic key management

  @REQ-SIG-035 @happy
  Scenario: Usage guide explains implementation pattern
    Given a NovusPack package
    And a valid context
    When implementation pattern is examined
    Then high-level functions generate signature data using private key
    And high-level functions call AddSignature with generated signature data
    And high-level functions handle errors from signature generation or addition

  @REQ-SIG-035 @error
  Scenario: Function usage guide handles errors
    Given a NovusPack package
    When function usage errors occur
    Then appropriate errors are returned
    And errors follow structured error format
