@domain:writing @m2 @REQ-WRITE-025 @spec(api_writing.md#42-clear-signatures-flag)
Feature: Clear Signatures Flag

  @REQ-WRITE-025 @happy
  Scenario: Clear-signatures flag allows signature removal
    Given a NovusPack package
    And a signed package
    When Write is called with clearSignatures flag set to true
    Then signature removal is allowed
    And new unsigned file is created
    And all signatures are stripped from new file

  @REQ-WRITE-025 @happy
  Scenario: Clear-signatures flag creates new unsigned file
    Given a NovusPack package
    And a signed package
    When Write is called with clearSignatures=true
    Then new file is created using SafeWrite (complete rewrite)
    And new filename must be different from current signed file
    And new file is unsigned
    And all package content is preserved

  @REQ-WRITE-025 @happy
  Scenario: Clear-signatures preserves package content
    Given a NovusPack package
    And a signed package
    When Write is called with clearSignatures=true
    Then all files are preserved
    And all metadata is preserved
    And all comments are preserved
    And immutability is reset (new file can be modified)

  @REQ-WRITE-025 @error
  Scenario: Clear-signatures requires different filename
    Given a NovusPack package
    And a signed package
    When Write is called with clearSignatures=true and same filename
    Then SameFilenameError is returned
    And error indicates filename requirement
    And error follows structured error format

  @REQ-WRITE-025 @error
  Scenario: Clear-signatures validates signed file state
    Given a NovusPack package
    And a corrupted signed package
    When Write is called with clearSignatures=true
    Then ValidationError is returned if file is corrupted
    And error indicates validation failure
    And error follows structured error format
