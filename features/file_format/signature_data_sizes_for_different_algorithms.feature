@domain:file_format @m1 @REQ-FILEFMT-023 @spec(package_file_format.md#73-signature-data-sizes)
Feature: Signature data sizes for different algorithms

  @happy
  Scenario: ML-DSA signature data size is correct
    Given a signature block with SignatureType 0x01
    When SignatureData size is examined
    Then SignatureData size is approximately 2,420-4,595 bytes
    And size depends on security level
    And size matches SignatureSize

  @happy
  Scenario: SLH-DSA signature data size is correct
    Given a signature block with SignatureType 0x02
    When SignatureData size is examined
    Then SignatureData size is approximately 7,856-17,088 bytes
    And size depends on security level
    And size matches SignatureSize

  @happy
  Scenario: PGP signature data size is correct
    Given a signature block with SignatureType 0x03
    When SignatureData size is examined
    Then SignatureData size is variable
    And typical size is 256-512 bytes
    And size matches SignatureSize

  @happy
  Scenario: X.509 signature data size is correct
    Given a signature block with SignatureType 0x04
    When SignatureData size is examined
    Then SignatureData size is variable
    And typical size is 256-4096 bytes
    And size matches SignatureSize

  @happy
  Scenario: Signature data size matches SignatureType requirements
    Given a signature block
    When SignatureData size is validated
    Then size is appropriate for SignatureType
    And size matches algorithm requirements
    And size is within expected range
