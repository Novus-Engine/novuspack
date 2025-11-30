@domain:signatures @security @m2 @REQ-SIG-001 @spec(api_signatures.md#11-multiple-signature-management-incremental-signing)
Feature: Add digital signatures incrementally

  @happy
  Scenario: Multiple signature types supported including quantum-safe
    Given a package with content to sign
    When I add a digital signature of type "ed25519"
    Then the package should include a signature entry of type "ed25519"

  @happy
  Scenario: AddSignature appends signature incrementally
    Given an unsigned package
    When AddSignature is called
    Then signature is appended to end of file
    And SignatureOffset is set if first signature
    And flags bit 0 is set if first signature

  @happy
  Scenario: First signature sets flags bit and SignatureOffset
    Given an unsigned package with SignatureOffset = 0
    When first signature is added
    Then flags bit 0 is set to 1
    And SignatureOffset is set to signature location
    And package becomes immutable

  @happy
  Scenario: Subsequent signatures are appended
    Given a signed package with existing signature
    When additional signature is added
    Then new signature is appended after existing signature
    And existing signature remains unchanged
    And all signatures validate correctly

  @happy
  Scenario: AddSignature validates header state
    Given a package
    When AddSignature is called
    Then header state is validated
    And package integrity is verified
    And validation occurs before signing

  @happy
  Scenario: RemoveSignature removes signature and later signatures
    Given a package with multiple signatures
    When RemoveSignature is called with index
    Then specified signature is removed
    And all later signatures are removed
    And earlier signatures remain

  @REQ-SIG-006 @happy
  Scenario: GetSignatureCount returns total number of signatures
    Given a package with signatures
    When GetSignatureCount is called
    Then total number of signatures is returned
    And count matches actual signatures

  @REQ-SIG-007 @happy
  Scenario: GetSignature retrieves signature by index
    Given a package with multiple signatures
    When GetSignature is called with index
    Then SignatureInfo for that index is returned
    And signature information is complete

  @REQ-SIG-008 @happy
  Scenario: GetAllSignatures retrieves all signatures
    Given a package with multiple signatures
    When GetAllSignatures is called
    Then all SignatureInfo objects are returned
    And all signatures are included

  @REQ-SIG-009 @happy
  Scenario: ClearAllSignatures removes all signatures
    Given a package with signatures
    When ClearAllSignatures is called
    Then all signatures are removed
    And SignatureOffset is set to 0
    And flags bit 0 is cleared

  @error
  Scenario: AddSignature fails if header state is invalid
    Given a package with invalid header state
    When AddSignature is called
    Then structured validation error is returned
    And signature is not added

  @error
  Scenario: Signature operations respect context cancellation
    Given a package
    And a cancelled context
    When signature operation is called
    Then structured context error is returned

  @REQ-SIG-015 @REQ-SIG-016 @error
  Scenario: AddSignature validates signature type parameter
    Given an open writable package
    When AddSignature is called with invalid signature type
    Then structured validation error is returned
    And error indicates unsupported signature type

  @REQ-SIG-015 @REQ-SIG-019 @error
  Scenario: AddSignature respects context cancellation
    Given an open writable package
    And a cancelled context
    When AddSignature is called
    Then structured context error is returned
    And error type is context cancellation
    And signature operation stops
