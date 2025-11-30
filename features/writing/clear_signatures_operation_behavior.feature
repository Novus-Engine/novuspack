@domain:writing @m2 @REQ-WRITE-026 @spec(api_writing.md#43-clear-signatures-behavior)
Feature: Clear-Signatures Operation Behavior

  @REQ-WRITE-026 @happy
  Scenario: Clear-signatures behavior creates new unsigned file when flag is true
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When Write is called with clearSignatures flag set to true
    Then new unsigned package file is created using SafeWrite
    And complete rewrite is performed
    And new file is unsigned

  @REQ-WRITE-026 @happy
  Scenario: Clear-signatures behavior requires different filename
    Given an open NovusPack package
    And the package has signatures
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then new unsigned file is created with different filename
    And original signed file is preserved
    And filename requirement is satisfied

  @REQ-WRITE-026 @happy
  Scenario: Clear-signatures behavior strips all signatures from new file
    Given an open NovusPack package
    And the package has signatures
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then all signatures are stripped from the new file
    And SignatureOffset is set to 0 in new file
    And new file has no signatures

  @REQ-WRITE-026 @happy
  Scenario: Clear-signatures behavior preserves all package content
    Given an open NovusPack package
    And the package has signatures
    And the package contains files, metadata, and comments
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then all package content is preserved (files, metadata, comments)
    And new unsigned file contains all original content
    And content integrity is maintained

  @REQ-WRITE-026 @happy
  Scenario: Clear-signatures behavior resets immutability in new file
    Given an open NovusPack package
    And the package has signatures
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then new file can be modified normally
    And new file is not immutable
    And immutability is reset

  @REQ-WRITE-026 @happy
  Scenario: Clear-signatures behavior always uses SafeWrite for signed files
    Given an open NovusPack package
    And the package has signatures
    When Write is called with clearSignatures flag set to true
    Then SafeWrite is always used
    And FastWrite is not used for signed files
    And complete rewrite is performed

  @REQ-WRITE-026 @error
  Scenario: Clear-signatures behavior returns error when filename is same
    Given an open NovusPack package
    And the package has signatures
    When Write is called with clearSignatures flag set to true
    And the same filename is provided
    Then SameFilenameError is returned
    And error indicates filename must be different
    And write operation is refused
