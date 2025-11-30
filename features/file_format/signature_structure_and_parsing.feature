@domain:file_format @m1 @signing @REQ-FILEFMT-003 @REQ-FILEFMT-021 @spec(package_file_format.md#71-signature-structure)
Feature: Signature structure and parsing

  @happy
  Scenario: Signature block is discoverable
    Given a NovusPack file containing a signature block
    When the signature block directory is parsed
    Then the signature block type and offset are reported

  @happy
  Scenario: Signature structure contains all required fields
    Given a signature block
    When signature structure is examined
    Then SignatureType field is present (4 bytes)
    And SignatureSize field is present (4 bytes)
    And SignatureFlags field is present (4 bytes)
    And SignatureTimestamp field is present (4 bytes)
    And CommentLength field is present (2 bytes)
    And SignatureComment field is present (variable)
    And SignatureData field is present (variable)

  @happy
  Scenario: SignatureType identifies ML-DSA signature algorithm
    Given a signature block
    When SignatureType is examined
    Then SignatureType equals 0x01 for ML-DSA
    And SignatureType is a 32-bit unsigned integer

  @happy
  Scenario: SignatureType identifies SLH-DSA signature algorithm
    Given a signature block
    When SignatureType is examined
    Then SignatureType equals 0x02 for SLH-DSA
    And SignatureType is a 32-bit unsigned integer

  @happy
  Scenario: SignatureType identifies PGP signature algorithm
    Given a signature block
    When SignatureType is examined
    Then SignatureType equals 0x03 for PGP
    And SignatureType is a 32-bit unsigned integer

  @happy
  Scenario: SignatureType identifies X.509 signature algorithm
    Given a signature block
    When SignatureType is examined
    Then SignatureType equals 0x04 for X.509
    And SignatureType is a 32-bit unsigned integer

  @happy
  Scenario: SignatureSize indicates signature data length
    Given a signature block
    When SignatureSize is examined
    Then SignatureSize equals size of SignatureData
    And SignatureSize is a 32-bit unsigned integer
    And SignatureSize matches actual signature data

  @happy
  Scenario: SignatureFlags encodes signature metadata
    Given a signature block
    When SignatureFlags is examined
    Then bits 31-16 are reserved and must be 0
    And bits 15-8 encode signature features
    And bits 7-0 encode signature status
    And SignatureFlags is a 32-bit unsigned integer

  @happy
  Scenario: SignatureTimestamp records creation time
    Given a signature block
    When SignatureTimestamp is examined
    Then SignatureTimestamp is Unix nanoseconds
    And SignatureTimestamp indicates when signature was created
    And SignatureTimestamp is a 32-bit unsigned integer

  @happy
  Scenario: CommentLength indicates signature comment size
    Given a signature block
    When CommentLength is examined
    Then CommentLength equals length of SignatureComment including null terminator
    And CommentLength 0 indicates no comment
    And CommentLength is a 16-bit unsigned integer

  @happy
  Scenario: SignatureComment is UTF-8 null-terminated
    Given a signature block with comment
    When SignatureComment is examined
    Then SignatureComment is UTF-8 encoded
    And SignatureComment is null-terminated
    And comment length matches CommentLength

  @happy
  Scenario: SignatureData contains raw signature
    Given a signature block
    When SignatureData is examined
    Then SignatureData length matches SignatureSize
    And SignatureData contains algorithm-specific signature bytes
    And SignatureData format depends on SignatureType

  @happy
  Scenario: Signature validates all content up to creation point
    Given a signature block
    When signature validation is performed
    Then signature validates all package content up to signature
    And signature validates its own metadata
    And signature validates its comment

  @error
  Scenario: Invalid SignatureType is rejected
    Given a signature block with invalid SignatureType
    When signature is parsed
    Then parsing fails
    And structured validation error is returned

  @error
  Scenario: SignatureSize mismatch is detected
    Given a signature block where SignatureSize does not match actual data
    When signature is validated
    Then validation fails
    And structured corruption error is returned
