@domain:metadata @m2 @v2 @REQ-META-085 @spec(api_metadata.md#72-signatureinfo-structure)
Feature: SignatureInfo Structure

  @REQ-META-085 @happy
  Scenario: SignatureInfo structure provides signature information
    Given a NovusPack package
    When SignatureInfo structure is examined
    Then structure contains signature index
    And structure contains signature type
    And structure contains signature size and offset
    And structure contains signature flags and timestamp
    And structure contains signature comment
    And structure contains algorithm and security level
    And structure contains validation status

  @REQ-META-085 @happy
  Scenario: SignatureInfo contains signature metadata
    Given a NovusPack package
    And SignatureInfo structure
    When signature metadata is examined
    Then Index contains signature index in package
    And Type contains signature type (ML-DSA, SLH-DSA, PGP, X.509)
    And Size contains size of signature data in bytes
    And Offset contains offset to signature data

  @REQ-META-085 @happy
  Scenario: SignatureInfo contains signature flags and timestamp
    Given a NovusPack package
    And SignatureInfo structure
    When signature flags and timestamp are examined
    Then Flags contains signature-specific flags
    And Timestamp contains Unix timestamp when signature was created
    And Comment contains signature comment if any

  @REQ-META-085 @happy
  Scenario: SignatureInfo contains validation information
    Given a NovusPack package
    And SignatureInfo structure
    When validation information is examined
    Then Algorithm contains algorithm name/description
    And SecurityLevel contains algorithm security level
    And Valid indicates whether signature is valid
    And Trusted indicates whether signature is trusted
    And Error contains error message if validation failed

  @REQ-META-085 @error
  Scenario: SignatureInfo handles invalid signatures
    Given a NovusPack package
    When invalid signature information is provided
    Then Valid is set to false
    And Error contains validation error message
    And Trusted is set to false
