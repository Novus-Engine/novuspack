@domain:file_format @m2 @REQ-FILEFMT-066 @spec(package_file_format.md#721-signaturetype-field)
Feature: SignatureType Field

  @REQ-FILEFMT-066 @happy
  Scenario: SignatureType field stores signature type identifier
    Given a signature block
    When SignatureType field is examined
    Then SignatureType is a 32-bit unsigned integer
    And SignatureType identifies the signature algorithm used
    And SignatureType enables algorithm identification

  @REQ-FILEFMT-066 @happy
  Scenario: SignatureType value 0x01 identifies ML-DSA algorithm
    Given a signature block
    When SignatureType equals 0x01
    Then SignatureType represents ML-DSA algorithm
    And ML-DSA (Module-Lattice Digital Signature Algorithm) is identified
    And algorithm type is correctly encoded

  @REQ-FILEFMT-066 @happy
  Scenario: SignatureType value 0x02 identifies SLH-DSA algorithm
    Given a signature block
    When SignatureType equals 0x02
    Then SignatureType represents SLH-DSA algorithm
    And SLH-DSA (Stateless Hash-based Digital Signature Algorithm) is identified
    And algorithm type is correctly encoded

  @REQ-FILEFMT-066 @happy
  Scenario: SignatureType value 0x03 identifies PGP algorithm
    Given a signature block
    When SignatureType equals 0x03
    Then SignatureType represents PGP algorithm
    And PGP (Pretty Good Privacy) is identified
    And algorithm type is correctly encoded

  @REQ-FILEFMT-066 @happy
  Scenario: SignatureType value 0x04 identifies X.509 algorithm
    Given a signature block
    When SignatureType equals 0x04
    Then SignatureType represents X.509 algorithm
    And X.509 Certificate-based signature is identified
    And algorithm type is correctly encoded

  @REQ-FILEFMT-066 @happy
  Scenario: SignatureType values 0x05-0xFFFFFFFF are reserved
    Given a signature block
    When SignatureType is in range 0x05 to 0xFFFFFFFF
    Then SignatureType values are reserved for future signature types
    And reserved values enable future algorithm support
    And extensibility is supported
