@domain:signatures @m2 @REQ-SIG-062 @spec(api_signatures.md#221-signatureinfo-struct)
Feature: Signature GetData returns (T, error) not (T, bool)

  @REQ-SIG-062 @happy
  Scenario: GetData returns value and error for signature data
    Given a SignatureInfo with optional signature data
    When GetData is invoked to retrieve signature data
    Then the method returns (T, error) not (T, bool)
    And errors indicate absence or failure clearly
    And the behavior matches the SignatureInfo structure specification
    And callers use error checking instead of boolean flags
