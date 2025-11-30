@domain:signatures @m2 @REQ-SIG-021 @spec(api_signatures.md#12-incremental-signing-implementation)
Feature: Incremental Signature Operations

  @REQ-SIG-021 @happy
  Scenario: Incremental signing implementation provides sequential signature support
    Given a NovusPack package
    And a valid context
    When incremental signing is used
    Then first signature signs all content up to its metadata and comment
    And first signature is appended at end of file
    And subsequent signatures sign all content up to that point including previous signatures
    And subsequent signatures are appended
    And each signature validates all content up to that point
    And context supports cancellation

  @REQ-SIG-021 @happy
  Scenario: First signature in incremental signing
    Given a NovusPack package
    And a valid context
    When first signature is added
    Then signature signs all content up to its own metadata and signature comment
    Then signature is appended at end of file
    And signature validates header and all content

  @REQ-SIG-021 @happy
  Scenario: Subsequent signatures in incremental signing
    Given a NovusPack package
    And a valid context
    And a package with first signature
    When subsequent signature is added
    Then signature signs all content up to that point including previous signatures
    Then signature metadata header and comment and data are appended
    And all previous signatures remain valid

  @REQ-SIG-021 @happy
  Scenario: Signature removal removes signature and all later signatures
    Given a NovusPack package
    And a valid context
    And a package with multiple signatures
    When RemoveSignature is called with index
    Then signature at index is removed
    And all later signatures are removed
    And remaining signatures remain valid

  @REQ-SIG-021 @error
  Scenario: Incremental signing handles errors
    Given a NovusPack package
    When signature operations fail
    Then appropriate errors are returned
    And errors follow structured error format
