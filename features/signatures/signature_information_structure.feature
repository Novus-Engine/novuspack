@domain:signatures @m2 @v2 @REQ-SIG-026 @spec(api_signatures.md#22-signature-information-structure)
Feature: Signature Information Structure

  @REQ-SIG-026 @happy
  Scenario: Signature information structure defines signature metadata format
    Given a NovusPack package
    When signature information structure is examined
    Then structure defines signature metadata format
    And structure contains signature type identifier
    And structure contains signature data size
    And structure contains signature data offset
    And structure contains signature-specific flags
    And structure contains signature creation timestamp

  @REQ-SIG-026 @happy
  Scenario: Signature information structure includes validation status
    Given a NovusPack package
    When signature information structure is retrieved
    Then structure contains validation status (Valid field)
    And structure contains error message if validation failed
    And structure contains algorithm name
    And structure contains security level information

  @REQ-SIG-026 @happy
  Scenario: Signature information structure enables signature management
    Given a NovusPack package
    When signature information structure is used
    Then structure enables signature retrieval
    And structure enables signature validation
    And structure enables signature metadata access
    And structure supports signature management operations

  @REQ-SIG-026 @error
  Scenario: Signature information structure handles invalid signatures
    Given a NovusPack package
    And an invalid signature
    When signature information structure is retrieved
    Then Valid field indicates signature is invalid
    And Error field contains validation error message
    And structure provides details for debugging
