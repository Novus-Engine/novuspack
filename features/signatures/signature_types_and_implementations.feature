@domain:signatures @security @m2 @REQ-SIG-004 @spec(api_signatures.md#2-signature-types)
Feature: Signature types and implementations

  @happy
  Scenario: ML-DSA signature type is supported
    Given a package
    When ML-DSA signature is added
    Then ML-DSA signature type is used
    And signature data matches ML-DSA format
    And signature is quantum-safe

  @happy
  Scenario: SLH-DSA signature type is supported
    Given a package
    When SLH-DSA signature is added
    Then SLH-DSA signature type is used
    And signature data matches SLH-DSA format
    And signature is quantum-safe

  @happy
  Scenario: PGP signature type is supported
    Given a package
    When PGP signature is added
    Then PGP signature type is used
    And signature data matches PGP format
    And OpenPGP compatibility is maintained

  @happy
  Scenario: X.509 signature type is supported
    Given a package
    When X.509 signature is added
    Then X.509 signature type is used
    And signature data matches X.509 format
    And PKCS#7 compatibility is maintained

  @happy
  Scenario: Signature validation works for all types
    Given a package with signature of any type
    When signature validation is performed
    Then validation succeeds for valid signatures
    And validation fails for invalid signatures
    And error indicates validation failure

  @error
  Scenario: Invalid signature type is rejected
    Given a package
    When AddSignature is called with invalid signature type
    Then structured validation error is returned
    And error indicates unsupported signature type
