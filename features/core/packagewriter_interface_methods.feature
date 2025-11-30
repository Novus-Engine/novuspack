@domain:core @m1 @REQ-CORE-005 @spec(api_core.md#12-package-writer-interface)
Feature: PackageWriter interface methods

  @happy
  Scenario: WriteFile adds file to package
    Given an open writable NovusPack package
    When WriteFile is called with path and data
    Then file is added to package
    And file is accessible via ReadFile
    And file index is updated

  @happy
  Scenario: WriteFile respects AddFileOptions
    Given an open writable NovusPack package
    When WriteFile is called with AddFileOptions
    Then file is added with specified options
    And compression settings are applied
    And encryption settings are applied

  @happy
  Scenario: RemoveFile removes file from package
    Given an open writable NovusPack package with files
    When RemoveFile is called with a file path
    Then file is removed from package
    And file index is updated
    And file is no longer accessible

  @happy
  Scenario: Write performs general write operation
    Given an open writable NovusPack package
    When Write is called with path and options
    Then appropriate write strategy is selected
    And write operation completes successfully

  @happy
  Scenario: SafeWrite performs atomic write
    Given an open writable NovusPack package
    When SafeWrite is called
    Then atomic write operation is performed
    And temp file strategy is used
    And package integrity is maintained

  @happy
  Scenario: FastWrite performs in-place updates
    Given an open writable NovusPack package
    When FastWrite is called
    Then in-place update is performed
    And write performance is optimized
    And package is updated efficiently

  @error
  Scenario: WriteFile fails if package is read-only
    Given a read-only open NovusPack package
    When WriteFile is called
    Then a structured validation error is returned
    And error type is ErrTypeValidation

  @error
  Scenario: RemoveFile fails for non-existent file
    Given an open writable NovusPack package
    When RemoveFile is called with non-existent path
    Then a structured validation error is returned

  @error
  Scenario: Write operations respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When write operation is called
    Then a structured context error is returned

  @REQ-CORE-015 @REQ-CORE-018 @error
  Scenario: PackageReader methods validate input parameters
    Given an open NovusPack package
    When ReadFile is called with empty path
    Then structured validation error is returned
    And error indicates invalid path

  @REQ-CORE-015 @REQ-CORE-016 @error
  Scenario: PackageWriter methods respect context cancellation
    Given an open writable NovusPack package
    And a cancelled context
    When PackageWriter method is called
    Then structured context error is returned
    And error type is context cancellation
