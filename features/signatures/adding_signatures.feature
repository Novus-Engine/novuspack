@domain:signatures @m2 @v2 @REQ-SIG-022 @spec(api_signatures.md#121-adding-subsequent-signatures)
Feature: Adding Signatures

  @REQ-SIG-022 @happy
  Scenario: Adding subsequent signatures supports incremental signing
    Given a NovusPack package
    And a valid context
    And an existing signed package
    When AddSignature is called with new signature
    Then new signature is appended incrementally
    And all previous signatures remain valid
    And signature signs all content up to that point
    And context supports cancellation

  @REQ-SIG-022 @happy
  Scenario: Adding subsequent signature appends without invalidating previous signatures
    Given a NovusPack package
    And a valid context
    And a package with existing signatures
    When subsequent signature is added
    Then new signature metadata header is appended
    And new signature comment is appended
    Then new signature data is appended
    And all previous signatures remain unchanged

  @REQ-SIG-022 @happy
  Scenario: Adding subsequent signature validates content up to that point
    Given a NovusPack package
    And a valid context
    And a package with existing signatures
    When subsequent signature is added
    Then signature validates all content up to that point
    And signature includes previous signatures in validation
    And signature validates its own metadata and comment

  @REQ-SIG-022 @error
  Scenario: Adding subsequent signatures handles errors
    Given a NovusPack package
    When invalid signature data is provided
    Then signature validation detects invalid data
    And appropriate errors are returned
    And errors follow structured error format
