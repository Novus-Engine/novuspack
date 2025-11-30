@domain:writing @m2 @REQ-WRITE-028 @REQ-WRITE-052 @spec(api_writing.md#45-use-cases)
Feature: Writing Use Cases

  @REQ-WRITE-028 @happy
  Scenario: Use case: Development workflow clears signatures to continue development
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then signatures are cleared
    And new unsigned file is created
    And development can continue on the new file
    And original signed file is preserved

  @REQ-WRITE-028 @happy
  Scenario: Use case: Package modification preserves content while clearing signatures
    Given an open NovusPack package
    And the package has signatures
    And the package needs modifications
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then package content is preserved (files, metadata, comments)
    And modifications can be made
    And new unsigned file contains all content
    And original signed file remains unchanged

  @REQ-WRITE-028 @happy
  Scenario: Use case: Signature management removes signatures before re-signing
    Given an open NovusPack package
    And the package has signatures
    And different keys are needed for re-signing
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then signatures are removed
    And package can be re-signed with different keys
    And signature workflow is supported

  @REQ-WRITE-028 @happy
  Scenario: Use case: Testing creates unsigned copies of signed packages
    Given an open NovusPack package
    And the package has signatures
    And unsigned copy is needed for testing
    When Write is called with clearSignatures flag set to true
    And a different filename is provided
    Then unsigned copy is created for testing
    And testing can proceed on unsigned copy
    And original signed file is preserved

  @REQ-WRITE-052 @happy
  Scenario: Use case: New package creation uses SafeWrite
    Given a new NovusPack package
    And no package file exists
    When Write is called with the target path
    Then SafeWrite is used for new package creation
    And new package file is created
    And package is ready for use

  @REQ-WRITE-052 @happy
  Scenario: Use case: Incremental file updates use FastWrite
    Given an open NovusPack package
    And an existing package file exists
    And individual files are being added or modified
    When Write is called with the target path
    Then FastWrite is used for incremental updates
    And only changed data is written
    And performance is optimized

  @REQ-WRITE-052 @happy
  Scenario: Use case: Critical operations use SafeWrite for maximum safety
    Given an open NovusPack package
    And data integrity is paramount
    When Write is called with the target path
    Then SafeWrite is used for critical operations
    And maximum safety and atomicity are provided
    And rollback capability exists

  @REQ-WRITE-052 @happy
  Scenario: Use case: Large package modifications use appropriate strategy
    Given an open NovusPack package
    And the package is large (>1GB)
    When Write is called with the target path
    Then appropriate write strategy is selected
    And streaming is used for large packages
    And memory usage is controlled

  @REQ-WRITE-052 @happy
  Scenario: Use case: Complete package rewrite uses SafeWrite
    Given an open NovusPack package
    And entire package content changes
    When Write is called with the target path
    Then SafeWrite is used for complete rewrite
    And complete replacement is performed
    And operation is atomic

  @REQ-WRITE-052 @happy
  Scenario: Use case: Defragmentation uses SafeWrite with guaranteed atomicity
    Given an open NovusPack package
    And package reorganization is needed
    When Write is called with the target path
    Then SafeWrite is used for defragmentation
    And complete package reorganization is performed
    And atomicity is guaranteed

  @REQ-WRITE-028 @error
  Scenario: Use case: Signed file write returns error when clearSignatures is false
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When Write is called with clearSignatures flag set to false
    Then SignedFileError is returned
    And error indicates signature protection
    And write operation is refused

  @REQ-WRITE-028 @error
  Scenario: Use case: Signed file write returns error when filename is not different
    Given an open NovusPack package
    And the package has signatures
    When Write is called with clearSignatures flag set to true
    And the same filename is provided
    Then SameFilenameError is returned
    And error indicates filename must be different
    And write operation is refused
