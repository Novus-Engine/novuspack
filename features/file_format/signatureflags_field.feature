@domain:file_format @m2 @REQ-FILEFMT-067 @spec(package_file_format.md#722-signatureflags-field)
Feature: SignatureFlags Field

  @REQ-FILEFMT-067 @happy
  Scenario: SignatureFlags field stores signature flags
    Given a signature block
    When SignatureFlags field is examined
    Then SignatureFlags is a 32-bit unsigned integer
    And SignatureFlags stores signature-specific metadata
    And SignatureFlags encodes signature options

  @REQ-FILEFMT-067 @happy
  Scenario: SignatureFlags reserved bits must be zero
    Given a signature block
    When SignatureFlags is examined
    Then bits 31-16 are reserved for future use
    And reserved bits must be 0
    And reserved bits enable future extensibility

  @REQ-FILEFMT-067 @happy
  Scenario: SignatureFlags features bits encode metadata
    Given a signature block
    When signature features bits are examined
    Then bit 15 (7 of features) indicates has timestamp
    And bit 14 (6 of features) indicates has metadata
    And bit 13 (5 of features) indicates has chain validation
    And bit 12 (4 of features) indicates has revocation
    And bit 11 (3 of features) indicates has expiration
    And bits 10-8 (2-0 of features) are reserved

  @REQ-FILEFMT-067 @happy
  Scenario: SignatureFlags status bits encode validation state
    Given a signature block
    When signature status bits are examined
    Then bit 7 indicates signature is valid
    And bit 6 indicates signature is verified
    And bit 5 indicates signature is trusted
    And bits 4-0 are reserved for future status flags

  @REQ-FILEFMT-067 @error
  Scenario: SignatureFlags with non-zero reserved bits is invalid
    Given a signature block
    And SignatureFlags has non-zero reserved bits (31-16)
    When signature is validated
    Then validation fails
    And structured validation error is returned
    And reserved bits violation is detected
