@domain:metadata @m2 @REQ-META-086 @spec(api_metadata.md#73-securitystatus-structure)
Feature: Metadata: SecurityStatus Structure

  @REQ-META-086 @happy
  Scenario: SecurityStatus structure provides security status information
    Given a NovusPack package
    When SecurityStatus structure is examined
    Then structure contains signature count information
    And structure contains signature validation results
    And structure contains checksum information
    And structure contains validation errors

  @REQ-META-086 @happy
  Scenario: SecurityStatus contains signature information
    Given a NovusPack package
    And SecurityStatus structure
    When signature information is examined
    Then SignatureCount contains number of signatures
    And ValidSignatures contains number of valid signatures
    And TrustedSignatures contains number of trusted signatures
    And SignatureResults contains individual validation results

  @REQ-META-086 @happy
  Scenario: SecurityStatus contains checksum information
    Given a NovusPack package
    And SecurityStatus structure
    When checksum information is examined
    Then HasChecksums indicates if checksums are present
    And ChecksumsValid indicates if checksums are valid

  @REQ-META-086 @error
  Scenario: SecurityStatus handles security violations
    Given a NovusPack package
    When security violations are detected
    Then ValidationErrors contains violation details
    And appropriate security status is reported
