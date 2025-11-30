@domain:signatures @m2 @REQ-SIG-028 @spec(api_signatures.md#222-signaturevalidationresult-struct)
Feature: Signatures: SignatureValidationResult Structure

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
  Scenario: SignatureValidationResult provides per-signature validation details
    Given a NovusPack package
    And multiple signatures have been validated
    When SignatureValidationResult structures are retrieved
    Then each result contains Index identifying signature position
    And each result contains Type identifying signature algorithm
    And each result contains Valid indicating signature validity
    And each result contains Trusted indicating signature trust status
    And results enable comprehensive signature validation assessment

  @REQ-SIG-028 @happy
  Scenario: SignatureValidationResult includes error information when validation fails
    Given a NovusPack package
    And a signature that fails validation
    When SignatureValidationResult structure is retrieved
    Then Valid field indicates signature is invalid
    And Error field contains detailed error message
    And structure provides context for validation failure

  @REQ-SIG-028 @happy
  Scenario: SignatureValidationResult includes public key information when available
    Given a NovusPack package
    And signature validation performed with public key
    When SignatureValidationResult structure is retrieved
    Then PublicKey field contains public key used for validation
    And PublicKey information enables key verification
    And structure supports key-based validation workflows

  @REQ-SIG-028 @error
  Scenario: SignatureValidationResult handles missing signature gracefully
    Given a NovusPack package
    And a signature index that does not exist
    When SignatureValidationResult structure is retrieved
    Then appropriate error is returned
    And error indicates signature not found
    And error follows structured error format
