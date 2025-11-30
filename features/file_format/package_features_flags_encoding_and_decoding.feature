@domain:file_format @m1 @REQ-FILEFMT-008 @spec(package_file_format.md#25-package-features-flags)
Feature: Package features flags encoding and decoding

  @happy
  Scenario: Flags field structure is correctly encoded
    Given a NovusPack package header
    When flags are examined
    Then bits 31-24 are reserved and must be 0
    And bits 23-16 are reserved and must be 0
    And bits 15-8 encode package compression type
    And bits 7-0 encode package features

  @happy
  Scenario Outline: Package compression type is encoded in bits 15-8
    Given a NovusPack package
    When package compression type is set to <CompressionType>
    Then flags bits 15-8 equal <EncodedValue>
    And compression type can be decoded correctly

    Examples:
      | CompressionType | EncodedValue |
      | 0                | 0x00        |
      | 1                | 0x01        |
      | 2                | 0x02        |
      | 3                | 0x03        |

  @error
  Scenario: Invalid compression type values are rejected
    Given a NovusPack package
    When package compression type is set to a value greater than 3
    Then a structured invalid compression type error is returned

  @happy
  Scenario: Bit 7 indicates metadata-only package
    Given a NovusPack package with only special metadata files
    When the package is created
    Then flags bit 7 is set to 1
    And flags bit 7 indicates metadata-only package

  @happy
  Scenario: Bit 6 indicates presence of special metadata files
    Given a NovusPack package containing special metadata files
    When the package is created
    Then flags bit 6 is set to 1
    And flags bit 6 indicates special metadata files are present

  @happy
  Scenario: Bit 5 indicates per-file tags are present
    Given a NovusPack package with files having tags
    When the package is created
    Then flags bit 5 is set to 1
    And flags bit 5 indicates per-file tags are used

  @happy
  Scenario: Bit 4 indicates package comment is present
    Given a NovusPack package with a comment
    When the package is created
    Then flags bit 4 is set to 1
    And flags bit 4 indicates package has comment

  @happy
  Scenario: Bit 3 indicates extended attributes are present
    Given a NovusPack package with files having extended attributes
    When the package is created
    Then flags bit 3 is set to 1
    And flags bit 3 indicates extended attributes are used

  @happy
  Scenario: Bit 2 indicates encrypted files are present
    Given a NovusPack package with encrypted files
    When the package is created
    Then flags bit 2 is set to 1
    And flags bit 2 indicates package has encrypted files

  @happy
  Scenario: Bit 1 indicates compressed files are present
    Given a NovusPack package with compressed files
    When the package is created
    Then flags bit 1 is set to 1
    And flags bit 1 indicates package has compressed files

  @happy
  Scenario: Bit 0 indicates package is signed
    Given a NovusPack package with digital signatures
    When the first signature is added
    Then flags bit 0 is set to 1
    And flags bit 0 indicates package has signatures

  @error
  Scenario: Flags bit 0 must be set before adding first signature
    Given a NovusPack package without flags bit 0 set
    When a signature is attempted to be added
    Then a structured invalid flags error is returned
    And the signature addition fails

  @happy
  Scenario Outline: Multiple feature flags can be combined
    Given a NovusPack package
    When flags bits are set to <FeatureBits>
    Then flags encode <Features>
    And all feature flags are decoded correctly

    Examples:
      | FeatureBits | Features                                  |
      | 0x07        | Has signatures, compressed, and encrypted |
      | 0x0F        | All content features enabled              |
      | 0x38        | Metadata features enabled                 |

  @error
  Scenario: Reserved flag bits must be zero
    Given a NovusPack package header
    When flags bits 31-24 or 23-16 are non-zero
    Then a structured invalid header error is returned
    And header validation fails
