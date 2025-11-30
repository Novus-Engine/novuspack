@domain:writing @m2 @REQ-WRITE-003 @spec(api_writing.md#3-write-strategy-selection)
Feature: Write strategy selection

  @happy
  Scenario: Strategy selection honors safety over performance
    Given a package with mixed write requirements
    When I select the write strategy
    Then safety should be prioritized over performance in the chosen strategy

  @happy
  Scenario: Write method selects SafeWrite for new packages
    Given a new package that does not exist
    When Write method is called
    Then SafeWrite is automatically selected
    And new package is created safely

  @happy
  Scenario: Write method selects SafeWrite for complete rewrites
    Given an existing package requiring complete rewrite
    When Write method is called
    Then SafeWrite is automatically selected
    And complete rewrite is performed safely

  @happy
  Scenario: Write method attempts FastWrite for existing unsigned packages
    Given an existing unsigned package
    When Write method is called with incremental changes
    Then FastWrite is attempted first
    And in-place update is performed if possible

  @happy
  Scenario: Write method falls back to SafeWrite if FastWrite fails
    Given an existing package
    And FastWrite operation fails
    When Write method is called
    Then fallback to SafeWrite occurs
    And write operation completes successfully

  @happy
  Scenario: Write method refuses signed packages without clearSignatures
    Given a signed package with SignatureOffset > 0
    When Write method is called without clearSignatures flag
    Then write operation is refused
    And structured immutability error is returned

  @happy
  Scenario: Write method uses SafeWrite for compressed packages
    Given a compressed package
    When Write method is called
    Then SafeWrite is selected
    And compressed package is written correctly

  @happy
  Scenario: Strategy selection considers package characteristics
    Given a package with specific characteristics
    When Write method is called
    Then strategy selection considers package size
    And strategy selection considers change scope
    And strategy selection considers package state
    And appropriate strategy is selected

  @error
  Scenario: Write method validates package state before selection
    Given a package in invalid state
    When Write method is called
    Then structured validation error is returned
    And strategy selection does not proceed

  @REQ-WRITE-008 @REQ-WRITE-009 @error
  Scenario: SelectWriteStrategy validates path parameter
    Given a package
    When SelectWriteStrategy is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-WRITE-008 @REQ-WRITE-011 @error
  Scenario: SelectWriteStrategy respects context cancellation
    Given a package
    And a cancelled context
    When SelectWriteStrategy is called
    Then structured context error is returned
    And error type is context cancellation
