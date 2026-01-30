@domain:core @m1 @REQ-CORE-005 @spec(api_core.md#12-package-writer-interface)
Feature: PackageWriter interface methods

  @happy
  Scenario: AddFileFromMemory adds file to package from memory
    Given an open writable NovusPack package
    When AddFileFromMemory is called with path and data
    Then file is added to in-memory package
    And file is accessible via ReadFile
    And file index is updated in memory
    And changes are not written to disk until Write is called

  @happy
  Scenario: AddFileFromMemory respects AddFileOptions
    Given an open writable NovusPack package
    When AddFileFromMemory is called with AddFileOptions
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
    When Write is called with context only
    Then appropriate write strategy is selected
    And write operation completes successfully
    And changes are written to the Package's configured target path

  @happy
  Scenario: SafeWrite performs atomic write
    Given an open writable NovusPack package
    When SafeWrite is called with overwrite flag
    Then atomic write operation is performed
    And temp file is created in same directory as target
    And temp file is atomically renamed to target
    And package integrity is maintained
    And all in-memory changes are made durable on disk

  @happy
  Scenario: SafeWrite rejects overwrite when overwrite flag is false
    Given an open writable NovusPack package
    And target file already exists
    When SafeWrite is called with overwrite=false
    Then a structured validation error is returned
    And error type is ErrTypeValidation

  @error
  Scenario: SafeWrite fails for cross-filesystem operation
    Given an open writable NovusPack package
    And target is on different filesystem than temp directory
    When SafeWrite is called
    Then a structured error is returned
    And error indicates cross-filesystem operation not supported

  @happy
  Scenario: FastWrite performs in-place updates
    Given an open writable NovusPack package
    And target path matches opened package path
    When FastWrite is called
    Then in-place update is performed
    And write performance is optimized
    And package is updated efficiently
    And all in-memory changes are made durable on disk

  @error
  Scenario: FastWrite fails when target path differs from opened path
    Given an open writable NovusPack package
    And target path differs from opened package path
    When FastWrite is called
    Then a structured validation error is returned
    And error indicates target path must match opened path

  @error
  Scenario: FastWrite fails when target file does not exist
    Given an open writable NovusPack package
    And target file does not exist
    When FastWrite is called
    Then a structured validation error is returned
    And error indicates target file must exist for in-place updates

  @error
  Scenario: FastWrite can corrupt file on interruption
    Given an open writable NovusPack package
    When FastWrite is interrupted during write phase
    Then package file may be corrupted
    And recovery file is created if possible

  @error
  Scenario: AddFileFromMemory fails if package is read-only
    Given a read-only open NovusPack package
    When AddFileFromMemory is called
    Then a structured security error is returned
    And error type is ErrTypeSecurity

  @error
  Scenario: AddFileFromMemory fails for invalid package path
    Given an open writable NovusPack package
    When AddFileFromMemory is called with empty path
    Then a structured validation error is returned
    And error type is ErrTypeValidation

  @error
  Scenario: SafeWrite fails when attempting to overwrite signed package
    Given an open writable NovusPack package
    And package is signed
    When SafeWrite is called with overwrite=true
    Then a structured security error is returned
    And error type is ErrTypeSecurity

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
