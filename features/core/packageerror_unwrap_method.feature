@domain:core @m2 @REQ-CORE-163 @spec(api_core.md#1023-packageerror-unwrap-method) @spec(api_core.md#packageerrorunwrap-method)
Feature: PackageError Unwrap method returns underlying error for error unwrapping

  @REQ-CORE-163 @happy
  Scenario: PackageError Unwrap returns the underlying error
    Given a PackageError that wraps another error
    When Unwrap is called
    Then the underlying error is returned
    And error unwrapping is supported for compatibility
    And the behavior matches the PackageError Unwrap specification
