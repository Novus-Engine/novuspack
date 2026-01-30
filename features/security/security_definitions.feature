@skip @domain:security @m2 @spec(security.md#511-security-classification)
Feature: Security Definitions

# This feature captures high-level security metadata expectations from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: Files may be assigned a security classification
    Given a file is included in a package
    When security metadata is applied to that file
    Then the file may be marked for script validation or resource limits
    And the file may be marked with a security level classification

  @documentation
  Scenario: Security metadata supports access control decisions
    Given a package contains both encrypted and unencrypted files
    When the package is inspected for security metadata
    Then encryption selection and security flags can be evaluated per file
    And callers can decide which files require additional validation before use
