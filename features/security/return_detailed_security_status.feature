@domain:security @security @m4 @spec(api_security.md#2-security-status-structure)
Feature: Return Detailed Security Status

  @REQ-SEC-002 @happy
  Scenario: SecurityValidationResult struct contains signature count information
    Given an open NovusPack package
    And package with signatures
    When SecurityValidationResult is examined
    Then SignatureCount contains number of signatures
    And ValidSignatures contains number of valid signatures
    And TrustedSignatures contains number of trusted signatures
    And signature information is populated consistently

  @REQ-SEC-002 @happy
  Scenario: SecurityValidationResult struct contains individual signature results
    Given an open NovusPack package
    And package with multiple signatures
    When SecurityValidationResult is examined
    Then SignatureResults array contains individual results
    And each result provides validation status
    And individual results are populated consistently

  @REQ-SEC-002 @happy
  Scenario: SecurityValidationResult struct contains checksum information
    Given an open NovusPack package
    And package with checksums
    When SecurityValidationResult is examined
    Then HasChecksums indicates checksums presence
    And ChecksumsValid indicates checksums validity
    And checksum information is populated consistently

  @REQ-SEC-002 @happy
  Scenario: SecurityValidationResult struct contains security level and errors
    Given an open NovusPack package
    And package validation results
    When SecurityValidationResult is examined
    Then SecurityLevel contains overall security level
    And ValidationErrors contains list of validation errors
    And security information is populated consistently

  @REQ-SEC-002 @happy
  Scenario: SecurityStatus struct provides comprehensive security information
    Given an open NovusPack package
    And package with mixed security states
    When GetSecurityStatus is called
    Then SecurityStatus structure is returned
    And all fields are populated consistently
    And security status provides complete information

  @REQ-SEC-002 @happy
  Scenario: SignatureValidationResult struct contains signature validation details
    Given an open NovusPack package
    And package with signatures
    When SignatureValidationResult is examined
    Then Index contains signature index
    And Type contains signature type identifier
    And Valid indicates signature validity
    And Trusted indicates signature trust status
    And Error contains error message if validation failed
    And Timestamp contains signature creation time
    And PublicKey contains public key data if available

  @REQ-SEC-002 @happy
  Scenario: SecurityLevel enum provides security level classification
    Given an open NovusPack package
    And package with security validation
    When SecurityLevel is examined
    Then SecurityLevelNone indicates no security
    And SecurityLevelLow indicates low security
    And SecurityLevelMedium indicates medium security
    And SecurityLevelHigh indicates high security
    And SecurityLevelMaximum indicates maximum security

  @REQ-SEC-002 @error
  Scenario: Security status handles validation failures consistently
    Given an open NovusPack package
    And package with validation failures
    When security status is queried
    Then ValidationErrors contains error details
    And SecurityLevel reflects validation state
    And status is populated consistently even with errors
