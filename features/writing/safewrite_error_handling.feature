@domain:writing @m2 @REQ-WRITE-016 @spec(api_writing.md#16-safewrite-error-handling)
Feature: SafeWrite Error Handling

  @REQ-WRITE-016 @happy
  Scenario: SafeWrite error handling validates target directory exists
    Given an open NovusPack package
    When SafeWrite is called with the target path
    Then target directory is validated
    And directory existence is checked
    And validation prevents errors

  @REQ-WRITE-016 @happy
  Scenario: SafeWrite error handling automatically cleans up temp file on failure
    Given an open NovusPack package
    And write operation encounters an error
    When SafeWrite fails
    Then temporary file is automatically cleaned up
    And cleanup is performed on failure
    And no temp files are left behind

  @REQ-WRITE-016 @happy
  Scenario: SafeWrite error handling ensures no partial writes possible
    Given an open NovusPack package
    And write operation encounters an error
    When SafeWrite fails
    Then no partial writes occur
    And atomic rename is not performed
    And original file remains intact

  @REQ-WRITE-016 @happy
  Scenario: SafeWrite error handling propagates clear error messages
    Given an open NovusPack package
    And write operation encounters an error
    When SafeWrite fails
    Then clear error messages are returned for debugging
    And error details are provided
    And error follows structured error format

  @REQ-WRITE-016 @happy
  Scenario: SafeWrite error handling handles streaming failures gracefully
    Given an open NovusPack package
    And the package is large (>100MB)
    And streaming operation fails
    When SafeWrite encounters streaming error
    Then streaming failure is handled gracefully
    And cleanup is performed
    And error is returned with details

  @REQ-WRITE-016 @error
  Scenario: SafeWrite error handling returns error when directory validation fails
    Given an open NovusPack package
    And target directory does not exist
    And target directory cannot be created
    When SafeWrite is called with the target path
    Then directory validation error is returned
    And error indicates directory issue
    And error follows structured error format

  @REQ-WRITE-016 @error
  Scenario: SafeWrite error handling returns error when path is invalid
    Given an open NovusPack package
    And path parameter is invalid (empty, not normalized, not writable)
    When SafeWrite is called with invalid path
    Then path validation error is returned
    And error indicates path issue
    And error follows structured error format
