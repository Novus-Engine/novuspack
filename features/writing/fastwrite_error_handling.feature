@domain:writing @m2 @REQ-WRITE-020 @spec(api_writing.md#25-fastwrite-error-handling)
Feature: FastWrite Error Handling

  @REQ-WRITE-020 @happy
  Scenario: FastWrite error handling validates existing package before modification
    Given an open NovusPack package
    And an existing package file exists
    When FastWrite is called with the target path
    Then existing package is validated before modification
    And package integrity is checked
    And validation prevents errors

  @REQ-WRITE-020 @happy
  Scenario: FastWrite error handling supports partial recovery from failures
    Given an open NovusPack package
    And an existing package file exists
    And FastWrite encounters a partial failure
    When FastWrite fails during execution
    Then partial recovery is possible
    And successfully updated entries are preserved
    And recovery mechanism works correctly

  @REQ-WRITE-020 @happy
  Scenario: FastWrite error handling tracks what was successfully updated
    Given an open NovusPack package
    And an existing package file exists
    And FastWrite encounters a failure
    When FastWrite fails during execution
    Then change tracking identifies successfully updated entries
    And tracking information is maintained
    And recovery can use tracking information

  @REQ-WRITE-020 @happy
  Scenario: FastWrite error handling falls back to SafeWrite on failure
    Given an open NovusPack package
    And an existing package file exists
    And FastWrite encounters an error
    When FastWrite fails during execution
    Then fallback to SafeWrite is triggered
    And SafeWrite completes the operation
    And complete rewrite is performed

  @REQ-WRITE-020 @error
  Scenario: FastWrite error handling returns error for signed packages
    Given an open NovusPack package
    And the package has signatures (SignatureOffset > 0)
    When FastWrite is called with the target path
    Then SignedFileError is returned
    And error indicates FastWrite cannot be used with signed packages
    And error follows structured error format

  @REQ-WRITE-020 @error
  Scenario: FastWrite error handling returns error for compressed packages
    Given an open NovusPack package
    And the package is compressed
    When FastWrite is called with the target path
    Then FastWriteOnCompressed error is returned
    And error indicates FastWrite cannot be used with compressed packages
    And error follows structured error format

  @REQ-WRITE-020 @error
  Scenario: FastWrite error handling returns error when entry validation fails
    Given an open NovusPack package
    And an existing package file exists
    And existing package is corrupted or invalid
    When FastWrite is called with the target path
    Then validation error is returned
    And error indicates package validation failure
    And error follows structured error format
