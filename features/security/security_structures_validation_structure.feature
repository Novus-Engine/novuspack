@skip @domain:security @m2 @spec(api_security.md#21-securityvalidationresult-struct)
Feature: Security Structures

# This feature captures high-level security status structure expectations from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: SecurityValidationResult reports checksum validation status
    Given a package includes checksums
    When a security validation result is produced
    Then the result indicates whether checksums are present
    And the result indicates whether checksums validate successfully

  @documentation
  Scenario: SecurityValidationResult includes validation errors for diagnostics
    Given package validation finds one or more issues
    When the security validation result is produced
    Then the result includes a list of validation error strings
    And callers can surface those errors to users or logs
