@domain:core @m2 @REQ-CORE-164 @spec(api_core.md#1024-packageerror-is-method) @spec(api_core.md#packageerroris-method)
Feature: PackageError Is method implements error matching for error comparison

  @REQ-CORE-164 @happy
  Scenario: PackageError Is supports error matching
    Given a PackageError and a target error type
    When Is is called with the target
    Then error matching is performed as specified
    And errors.Is can be used for error comparison
    And the behavior matches the PackageError Is specification
