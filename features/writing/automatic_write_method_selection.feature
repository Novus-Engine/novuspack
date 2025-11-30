@domain:writing @m2 @REQ-WRITE-021 @REQ-WRITE-044 @spec(api_writing.md#31-automatic-selection-logic)
Feature: Automatic Write Method Selection

  @REQ-WRITE-021 @happy
  Scenario: Write method automatically selects SafeWrite for new packages
    Given an open NovusPack package
    And no package file exists at the target path
    When Write is called with the target path
    Then SafeWrite is automatically selected
    And a new package file is created
    And the operation completes successfully

  @REQ-WRITE-021 @happy
  Scenario: Write method automatically selects SafeWrite for signed packages
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When Write is called with the target path
    And clearSignatures flag is false
    Then write operation is refused
    And SignedFileError is returned
    And error indicates signature protection

  @REQ-WRITE-021 @happy
  Scenario: Write method automatically selects SafeWrite for compressed packages
    Given an open NovusPack package
    And the package is compressed (compression type in header flags)
    When Write is called with the target path
    Then SafeWrite is automatically selected
    And FastWrite is not attempted
    And compressed package is handled correctly

  @REQ-WRITE-021 @happy
  Scenario: Write method attempts FastWrite for existing unsigned packages
    Given an open NovusPack package
    And an existing package file at the target path
    And the package is unsigned (SignatureOffset = 0)
    And the package is uncompressed
    When Write is called with the target path
    Then FastWrite is attempted first
    And in-place update is performed if possible
    And operation completes successfully

  @REQ-WRITE-021 @happy
  Scenario: Write method falls back to SafeWrite when FastWrite fails
    Given an open NovusPack package
    And an existing package file at the target path
    And FastWrite fails due to an error condition
    When Write is called with the target path
    Then FastWrite failure triggers fallback
    And SafeWrite is automatically selected
    And complete rewrite is performed
    And operation completes successfully

  @REQ-WRITE-021 @happy
  Scenario: Write method selects SafeWrite for complete rewrite scenarios
    Given an open NovusPack package
    And the package requires complete rewrite
    When Write is called with the target path
    Then SafeWrite is automatically selected
    And complete rewrite is performed
    And operation completes successfully

  @REQ-WRITE-044 @happy
  Scenario: Automatic compression detection preserves compressed state
    Given an open NovusPack package
    And the package is compressed (compression type in header flags)
    When Write is called with the target path
    Then current compression state is detected
    And compressed package is written with compression preserved
    And header flags maintain compression type

  @REQ-WRITE-044 @happy
  Scenario: Automatic compression detection preserves uncompressed state
    Given an open NovusPack package
    And the package is uncompressed (compression type = 0)
    When Write is called with the target path
    Then current compression state is detected
    And uncompressed package is written without compression
    And header flags maintain no compression

  @REQ-WRITE-044 @happy
  Scenario: Automatic compression detection defaults to uncompressed for new packages
    Given a new NovusPack package
    And no compression state exists
    When Write is called with the target path
    And compressionType parameter is 0
    Then new package is created uncompressed by default
    And header flags indicate no compression
    And package is ready for use

  @REQ-WRITE-021 @error
  Scenario: Automatic selection returns error when signed package write is refused
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    And clearSignatures flag is false
    When Write is called with the target path
    Then SignedFileError is returned
    And error indicates signature protection prevents write
    And error follows structured error format

  @REQ-WRITE-021 @error
  Scenario: Automatic selection returns error when FastWrite fails and SafeWrite also fails
    Given an open NovusPack package
    And FastWrite fails due to an error condition
    And SafeWrite also fails due to an error condition
    When Write is called with the target path
    Then fallback to SafeWrite is attempted
    And SafeWrite error is returned
    And error follows structured error format
