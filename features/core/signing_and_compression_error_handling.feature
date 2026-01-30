@domain:core @m2 @REQ-CORE-037 @spec(api_core.md#5243-error-handling) @spec(api_writing.md#57-error-handling)
Feature: Signing and Compression Error Handling

  @REQ-CORE-037 @error
  Scenario: CompressSignedPackageError is returned when attempting to compress signed package
    Given an open NovusPack package
    And package has signatures
    And a valid context
    When compression operation is attempted
    Then CompressSignedPackageError is returned
    And error indicates signed package cannot be compressed
    And error follows structured error format

  @REQ-CORE-037 @error
  Scenario: All compression functions check for existing signatures
    Given an open NovusPack package
    And package has signatures
    And a valid context
    When compression functions are called
    Then validation checks for existing signatures
    And error is returned if signatures exist
    And validation prevents compression of signed packages

  @REQ-CORE-037 @happy
  Scenario: Clear workflow must be followed for compressing previously signed packages
    Given an open NovusPack package
    And package has signatures
    When clear workflow is followed
    Then first step is to clear signatures
    And second step is to compress package
    And third step is to re-sign if needed
    And workflow ensures proper ordering
