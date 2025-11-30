@domain:file_format @m1 @REQ-FILEFMT-022 @spec(package_file_format.md#72-signature-types)
Feature: Signature types and encoding

  @happy
  Scenario Outline: SignatureType identifies supported algorithms
    Given a signature block
    When SignatureType equals <Type>
    Then SignatureType represents <Algorithm>

    Examples:
      | Type | Algorithm |
      | 0x01 | ML-DSA    |
      | 0x02 | SLH-DSA   |
      | 0x03 | PGP       |
      | 0x04 | X.509     |

  @happy
  Scenario: SignatureFlags bit encoding is correct
    Given a signature block
    When SignatureFlags is examined
    Then bits 31-16 are reserved (must be 0)
    And bits 15-8 encode signature features
    And bits 7-0 encode signature status
    And bit encoding follows specification

  @happy
  Scenario: SignatureFlags features bits encode metadata
    Given a signature block
    When signature features bits are examined
    Then bit 7 indicates has timestamp
    And bit 6 indicates has metadata
    And bit 5 indicates has chain validation
    And bit 4 indicates has revocation
    And bit 3 indicates has expiration

  @happy
  Scenario: SignatureFlags status bits encode validation state
    Given a signature block
    When signature status bits are examined
    Then bit 7 indicates valid
    And bit 6 indicates verified
    And bit 5 indicates trusted

  @happy
  Scenario: SignatureTimestamp uses Unix nanoseconds
    Given a signature block
    When SignatureTimestamp is examined
    Then timestamp is Unix nanoseconds format
    And timestamp indicates creation time
    And timestamp is within valid range

  @error
  Scenario: Invalid SignatureFlags reserved bits are rejected
    Given a signature block with non-zero reserved bits
    When signature is validated
    Then validation fails
    And structured validation error is returned
