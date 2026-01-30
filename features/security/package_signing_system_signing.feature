@skip @domain:security @m2 @spec(security.md#2-package-signing-system)
Feature: Package Signing System

# This feature captures high-level signing system expectations from the security specs.
# Detailed runnable scenarios live in the dedicated signatures and security feature files.

  @documentation
  Scenario: Signed packages enforce immutability after the first signature
    Given a package file contains one or more signatures
    When a caller attempts a mutating write operation
    Then the operation is rejected due to signed package immutability
    And only signature addition operations remain allowed

  @documentation
  Scenario: Signature management and validation are specified in the signatures API
    Given a caller needs to sign or validate signatures
    When the caller selects an API for those operations
    Then the caller uses the Digital Signature API
    And the Security spec provides cross-references for the signing system
