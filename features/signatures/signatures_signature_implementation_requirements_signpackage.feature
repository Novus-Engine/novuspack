@domain:signatures @m2 @v2 @REQ-SIG-034 @spec(api_signatures.md#281-implementation-requirements)
Feature: Signatures: Signature Implementation Requirements (SignPackage)

  @REQ-SIG-034 @happy
  Scenario: SignPackage functions follow AddSignature implementation requirements
    Given a NovusPack package
    And a private key
    When SignPackage function is called
    Then signature bit is set if this is first signature
    And SignatureOffset is set if this is first signature
    And header state is validated before signing
    And signature is generated using private key
    And signature data is appended to package
    And immutability is maintained after signing

  @REQ-SIG-034 @happy
  Scenario: SignPackage functions validate header state
    Given a NovusPack package
    And a private key
    When SignPackage function is called
    Then package header is validated for valid signing state
    And content modifications are checked
    And package signing eligibility is verified
    And header validation occurs before signature generation

  @REQ-SIG-034 @happy
  Scenario: SignPackage functions append signature data
    Given a NovusPack package
    And a private key
    And signature data has been generated
    When SignPackage function completes
    Then signature metadata header is added (18 bytes)
    And signature comment is added if provided
    And signature data is appended
    And file size and offsets are updated

  @REQ-SIG-034 @error
  Scenario: SignPackage functions handle invalid header state
    Given a NovusPack package
    And an invalid header state
    When SignPackage function is called
    Then structured error is returned
    And error indicates invalid header state
    And error follows structured error format
