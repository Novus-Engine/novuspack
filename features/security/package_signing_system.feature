@domain:security @m2 @REQ-SEC-014 @spec(security.md#2-package-signing-system)
Feature: Package signing system provides digital signature support

  @REQ-SEC-014 @happy
  Scenario: Package signing system provides digital signature support
    Given a package that may be signed
    When package signing system is used
    Then digital signature support is provided as specified (v2)
    And the behavior matches the package signing system specification
    And signature operations are available when implemented
    And signing and verification workflows are supported
