@domain:security @m2 @REQ-SEC-078 @spec(api_security.md#22-securitystatus-struct)
Feature: Security: SecurityStatus Structure

  @REQ-SEC-078 @happy
  Scenario: SecurityStatus struct contains signature count information
    Given an open NovusPack package
    And package with signatures
    When SecurityStatus struct is examined
    Then SignatureCount field contains number of signatures
    And ValidSignatures field contains number of valid signatures
    And TrustedSignatures field contains number of trusted signatures

  @REQ-SEC-078 @happy
  Scenario: SecurityStatus struct contains individual signature results
    Given an open NovusPack package
    And package with multiple signatures
    When SecurityStatus struct is examined
    Then SignatureResults field contains array of SignatureValidationResult
    And individual signature validation results are accessible
    And results provide detailed validation information

  @REQ-SEC-078 @happy
  Scenario: SecurityStatus struct contains checksum information
    Given an open NovusPack package
    And package with checksums
    When SecurityStatus struct is examined
    Then HasChecksums field indicates if checksums are present
    And ChecksumsValid field indicates if checksums are valid
    And checksum information is accessible

  @REQ-SEC-078 @happy
  Scenario: SecurityStatus struct contains validation errors
    Given an open NovusPack package
    And package validation results
    When SecurityStatus struct is examined
    And ValidationErrors field contains list of validation errors
    And security assessment information is accessible

  @REQ-SEC-078 @happy
  Scenario: SecurityStatus struct is returned by GetSecurityStatus
    Given an open NovusPack package
    And package with security state
    When GetSecurityStatus is called
    Then SecurityStatus struct is returned
    And all fields are populated
    And comprehensive security status is provided

  @REQ-SEC-078 @happy
  Scenario: SecurityStatus struct provides complete security information
    Given an open NovusPack package
    And package with mixed security states
    When SecurityStatus struct is examined
    Then structure provides signature information
    And structure provides checksum information
    And structure provides validation errors
    And complete security information is available
