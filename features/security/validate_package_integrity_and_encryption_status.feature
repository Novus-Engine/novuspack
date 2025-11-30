@domain:security @security @m2 @REQ-SEC-001 @spec(api_security.md#1-package-validation)
Feature: Validate package integrity and encryption status

  @security
  Scenario: Validation covers signatures, encryption, and checksums
    Given a signed and encrypted package
    When I validate the package
    Then validation should succeed and report all checks passing

  @security
  Scenario: Package validation checks all security aspects
    Given a package
    When comprehensive validation is performed
    Then signatures are validated
    And encryption status is checked
    And checksums are verified
    And integrity is confirmed

  @security
  Scenario: SecurityStatus structure contains comprehensive information
    Given a package validation result
    When SecurityStatus is examined
    Then signature validation status is included
    And encryption status is included
    And checksum validation status is included
    And overall security status is provided

  @error
  Scenario: Validation fails if security checks fail
    Given a package with security issues
    When validation is performed
    Then structured security error is returned
    And specific security issue is identified
