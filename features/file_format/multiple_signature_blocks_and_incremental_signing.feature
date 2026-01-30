@domain:file_format @m1 @signing @v2 @REQ-FILEFMT-024 @spec(package_file_format.md#7-digital-signatures-section-optional)
Feature: Multiple signature blocks and incremental signing

  @happy
  Scenario: Multiple signature blocks are discoverable
    Given a NovusPack file containing multiple signature blocks
    When the signature block directory is parsed
    Then all signature block types and offsets are reported

  @happy
  Scenario: Signatures are appended sequentially
    Given a NovusPack package with existing signature
    When additional signature is added
    Then new signature is appended after existing signature
    And existing signature remains unchanged
    And signature order is preserved

  @happy
  Scenario: Each signature validates content up to its point
    Given a NovusPack package with multiple signatures
    When signature validation is performed
    Then each signature validates all content up to its creation point
    And each signature includes previous signatures in validation
    And each signature validates its own metadata and comment

  @happy
  Scenario: Incremental signatures do not invalidate existing signatures
    Given a NovusPack package with signature
    When additional signature is added
    Then existing signature remains valid
    And existing signature is not modified
    And all signatures validate correctly

  @happy
  Scenario: Signature blocks support different signature types
    Given a NovusPack package
    When signatures of different types are added
    Then ML-DSA signature can be added
    And SLH-DSA signature can be added
    And PGP signature can be added
    And X.509 signature can be added
    And all signature types coexist

  @happy
  Scenario: Signature blocks are optional
    Given a NovusPack package without signatures
    When package structure is examined
    Then SignatureOffset equals 0
    And no signature blocks are present
    And package is valid without signatures

  @error
  Scenario: Signatures cannot be modified in place
    Given a NovusPack package with signature
    When signature modification is attempted
    Then modification fails
    And structured immutability error is returned
    And signature remains unchanged

  @happy
  Scenario: Signature index points to first signature
    Given a NovusPack package with signatures
    When SignatureOffset is examined
    Then SignatureOffset points to first signature block
    And subsequent signatures follow sequentially
    And signature chain is discoverable
