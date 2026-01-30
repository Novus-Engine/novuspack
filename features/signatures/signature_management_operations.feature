@domain:signatures @security @m2 @v2 @REQ-SIG-005 @spec(api_signatures.md#1-signature-management)
Feature: Signature management operations

  @happy
  Scenario: RemoveSignature removes signature by index
    Given a package with multiple signatures
    When RemoveSignature is called with index
    Then specified signature is removed
    And all later signatures are removed
    And earlier signatures remain unchanged

  @happy
  Scenario: GetSignatureCount returns total number of signatures
    Given a package with signatures
    When GetSignatureCount is called
    Then total number of signatures is returned
    And count matches actual signatures

  @happy
  Scenario: GetSignature retrieves signature by index
    Given a package with multiple signatures
    When GetSignature is called with index
    Then SignatureInfo for that index is returned
    And signature information is complete

  @happy
  Scenario: GetAllSignatures returns all signatures
    Given a package with multiple signatures
    When GetAllSignatures is called
    Then all SignatureInfo objects are returned
    And all signatures are included
    And signatures are in order

  @happy
  Scenario: ClearAllSignatures removes all signatures
    Given a package with signatures
    When ClearAllSignatures is called
    Then all signatures are removed
    And SignatureOffset is set to 0
    And flags bit 0 is cleared
    And package becomes writable again

  @error
  Scenario: RemoveSignature fails with invalid index
    Given a package with signatures
    When RemoveSignature is called with invalid index
    Then structured validation error is returned

  @error
  Scenario: GetSignature fails with invalid index
    Given a package with signatures
    When GetSignature is called with invalid index
    Then structured validation error is returned

  @error
  Scenario: Signature management operations fail if package is read-only
    Given a read-only signed package
    When signature management operation is called
    Then structured validation error is returned

  @REQ-SIG-015 @REQ-SIG-016 @error
  Scenario: Signature operations validate signature type parameter
    Given an open writable package
    When AddSignature is called with invalid signature type
    Then structured validation error is returned
    And error indicates unsupported signature type

  @REQ-SIG-015 @REQ-SIG-017 @error
  Scenario: Signature operations validate index parameter
    Given a package with signatures
    When RemoveSignature or GetSignature is called with negative index
    Then structured validation error is returned
    And error indicates invalid index

  @REQ-SIG-015 @REQ-SIG-017 @error
  Scenario: Signature operations validate index bounds
    Given a package with 3 signatures
    When GetSignature is called with index 5
    Then structured validation error is returned
    And error indicates index out of bounds

  @REQ-SIG-015 @REQ-SIG-018 @error
  Scenario: Signature operations validate public key parameter
    Given an open writable package
    When ValidateSignatureWithKey is called with nil key
    Then structured validation error is returned
    And error indicates invalid key

  @REQ-SIG-015 @REQ-SIG-019 @error
  Scenario: Signature management operations respect context cancellation
    Given an open writable package
    And a cancelled context
    When signature management operation is called
    Then structured context error is returned
    And error type is context cancellation
    And operation stops
