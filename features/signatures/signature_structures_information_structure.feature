@domain:signatures @m2 @v2 @REQ-SIG-027 @REQ-SIG-028 @spec(api_signatures.md#221-signatureinfo-struct)
Feature: Signature Structures

  @REQ-SIG-027 @happy
  Scenario: SignatureInfo struct provides signature information structure
    Given a NovusPack package
    And a signature has been added to the package
    When SignatureInfo structure is retrieved
    Then structure contains Type field with signature type identifier
    And structure contains Size field with size of signature data in bytes
    And structure contains Offset field with offset to signature data
    And structure contains Flags field with signature-specific flags
    And structure contains Timestamp field with Unix timestamp
    And structure contains Data field with raw signature data
    And structure contains Algorithm field with algorithm name
    And structure contains SecurityLevel field with algorithm security level
    And structure contains Valid field indicating validation status
    And structure contains Error field with error message if validation failed

  @REQ-SIG-027 @happy
  Scenario: SignatureInfo provides complete signature metadata
    Given a NovusPack package
    And a signature has been added to the package
    When SignatureInfo structure is retrieved
    Then Type identifies the signature algorithm type
    And Size indicates the signature data size
    And Offset points to signature data location
    And Flags contain signature-specific configuration
    And Timestamp indicates when signature was created
    And Algorithm identifies the cryptographic algorithm used

  @REQ-SIG-027 @happy
  Scenario: SignatureInfo validation status reflects signature validity
    Given a NovusPack package
    And a valid signature has been added to the package
    When SignatureInfo structure is retrieved
    Then Valid field indicates signature is valid
    And Error field is empty when signature is valid
    And SecurityLevel indicates signature algorithm security level

  @REQ-SIG-027 @error
  Scenario: SignatureInfo indicates invalid signature
    Given a NovusPack package
    And an invalid signature has been added to the package
    When SignatureInfo structure is retrieved
    Then Valid field indicates signature is invalid
    And Error field contains error message describing validation failure
    And structure provides details for debugging

  @REQ-SIG-028 @happy
  Scenario: SignatureValidationResult struct provides validation result structure
    Given a NovusPack package
    And signature validation has been performed
    When SignatureValidationResult structure is retrieved
    Then structure contains Index field with signature index in package
    And structure contains Type field with signature type identifier
    And structure contains Valid field indicating validation status
    And structure contains Trusted field indicating trust status
    And structure contains Error field with error message if validation failed
    And structure contains Timestamp field indicating when signature was created
    And structure contains PublicKey field with public key used for validation

  @REQ-SIG-028 @happy
  Scenario: SignatureValidationResult provides per-signature validation status
    Given a NovusPack package
    And multiple signatures have been validated
    When SignatureValidationResult structures are retrieved
    Then each result contains Index identifying signature position
    And each result contains validation status (Valid/Invalid)
    And each result contains trust status (Trusted/Untrusted)
    And results enable per-signature validation assessment

  @REQ-SIG-028 @happy
  Scenario: SignatureValidationResult includes public key information when available
    Given a NovusPack package
    And signature validation has been performed with public key
    When SignatureValidationResult structure is retrieved
    Then PublicKey field contains public key used for validation
    And PublicKey information enables key verification
    And structure supports key-based validation workflows
