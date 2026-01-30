@domain:security @m2 @REQ-SEC-032 @spec(security.md#412-security-status-information)
Feature: Security Status Information

  @REQ-SEC-032 @happy
  Scenario: Security status information provides signature information
    Given an open NovusPack package
    And a valid context
    And package with signatures
    When GetSecurityStatus is called
    Then signature count is provided
    And valid signatures count is provided
    And trusted signatures count is provided
    And individual signature validation results are provided

  @REQ-SEC-032 @happy
  Scenario: Security status information provides checksum information
    Given an open NovusPack package
    And a valid context
    And package with checksums
    When GetSecurityStatus is called
    Then HasChecksums indicates checksums are present
    And ChecksumsValid indicates checksums are valid
    And checksum validation status is provided

  @REQ-SEC-032 @happy
  Scenario: Security status information provides validation errors
    Given an open NovusPack package
    And a valid context
    And package with validation issues
    When GetSecurityStatus is called
    Then ValidationErrors contains error details
    And errors indicate specific validation failures
    And error information supports troubleshooting

  @REQ-SEC-032 @happy
  Scenario: Security status information provides comprehensive validation results
    Given an open NovusPack package
    And a valid context
    And package requiring validation
    When GetSecurityStatus is called
    Then comprehensive security status is provided
    And status includes all security aspects
    And status provides complete security assessment
