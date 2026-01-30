@domain:signatures @m2 @REQ-SIG-061 @spec(api_signatures.md#221-signatureinfo-struct)
Feature: Signature typed representation uses Option as source of truth

  @REQ-SIG-061 @happy
  Scenario: Signature information uses Option for typed representation
    Given a SignatureInfo structure with optional signature data
    When signature data is accessed or checked
    Then Option is used as the source of truth for presence
    And typed representation is consistent with the specification
    And the behavior matches the SignatureInfo structure specification
    And callers can distinguish present from absent signature data
