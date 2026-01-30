@skip @domain:security @m2 @spec(security.md#11-security-layers)
Feature: Security Layers and Architecture

# This feature captures high-level security architecture expectations from the security specs.
# Detailed runnable scenarios live in the dedicated security feature files.

  @documentation
  Scenario: Security uses defense in depth across multiple layers
    Given a NovusPack package is distributed to consumers
    When the package is evaluated for security protections
    Then integrity checks, encryption, validation, access control, and transparency layers apply
    And no single layer is relied upon as the only protection

  @documentation
  Scenario: Security principles prioritize transparent and inspectable packages
    Given a package is subject to antivirus or security scanning
    When the package structure is inspected
    Then the package format supports transparency and inspection
    And security features do not require opaque or self-modifying behavior
