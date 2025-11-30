@domain:file_format @m1 @REQ-FILEFMT-010 @spec(package_file_format.md#29-signed-package-file-immutability-and-incremental-signatures)
Feature: Signed package file immutability enforcement

  @happy
  Scenario: Package becomes immutable after first signature
    Given a NovusPack package without signatures
    When the first signature is added
    Then SignatureOffset is set to a non-zero value
    And the entire file becomes immutable
    And flags bit 0 is set to 1

  @error
  Scenario: Header modifications are prohibited on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When a header field modification is attempted
    Then a structured immutability error is returned
    And the modification fails

  @error
  Scenario: File entry modifications are prohibited on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When a file entry modification is attempted
    Then a structured immutability error is returned
    And the modification fails

  @error
  Scenario: File data modifications are prohibited on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When file data modification is attempted
    Then a structured immutability error is returned
    And the modification fails

  @error
  Scenario: File index modifications are prohibited on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When file index modification is attempted
    Then a structured immutability error is returned
    And the modification fails

  @error
  Scenario: Package comment modifications are prohibited on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When package comment modification is attempted
    Then a structured immutability error is returned
    And the modification fails

  @happy
  Scenario: Read operations are allowed on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When file content is read
    Then the read operation succeeds
    And no immutability error is returned

  @happy
  Scenario: Signature addition is allowed on signed packages
    Given a signed NovusPack package with SignatureOffset > 0
    When an additional signature is added
    Then the signature is appended to the end of the file
    And existing signatures remain valid
    And SignatureOffset is updated

  @happy
  Scenario: Incremental signatures are appended sequentially
    Given a signed NovusPack package
    When a second signature is added
    Then the second signature is appended after the first
    And the first signature remains unchanged
    When a third signature is added
    Then the third signature is appended after the second
    And all previous signatures remain unchanged

  @happy
  Scenario: Each signature validates all content up to that point
    Given a NovusPack package with multiple signatures
    When a signature is validated
    Then the signature validates all content including previous signatures
    And the signature's metadata and comment are included

  @error
  Scenario: Unsigned packages allow all modifications
    Given an unsigned NovusPack package with SignatureOffset = 0
    When any modification is attempted
    Then the modification succeeds
    And no immutability error is returned
