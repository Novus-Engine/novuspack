@domain:writing @m2 @REQ-WRITE-024 @REQ-WRITE-065 @spec(api_writing.md#421-signed-package-protection) @spec(api_writing.md#422-compression-configuration)
Feature: Signed File Protection

  @REQ-WRITE-024 @happy
  Scenario: Signed file protection prevents modification of signed packages
    Given a NovusPack package
    And a signed package (SignatureOffset > 0)
    When write operation is attempted without clearSignatures flag
    Then write operation is refused
    And signature invalidation is prevented
    And signed file protection prevents modification

  @REQ-WRITE-024 @happy
  Scenario: Signed file protection validates SignatureOffset
    Given a NovusPack package
    And a signed package (SignatureOffset > 0)
    When write operation is attempted
    Then SignatureOffset is checked
    And write operation is refused if SignatureOffset > 0
    And signed file protection is enforced

  @REQ-WRITE-024 @error
  Scenario: Signed file protection returns error for write attempts
    Given a NovusPack package
    And a signed package
    When write operation is attempted without clearSignatures flag
    Then SignedFileError is returned
    And error indicates signed file protection
    And error follows structured error format

  @REQ-WRITE-024 @happy
  Scenario: Signed file protection allows clearSignatures flag
    Given a NovusPack package
    And a signed package
    When write operation is attempted with clearSignatures flag
    Then write operation is allowed
    And new unsigned file is created
    And signatures are stripped from new file
