@domain:core @m2 @REQ-CORE-081 @spec(api_core.md#113-normalizepackagepath-function) @spec(api_core.md#normalizepackagepath-error-handling) @spec(api_core.md#normalizepackagepath-return-value)
Feature: NormalizePackagePath function normalizes package paths

  @REQ-CORE-081 @happy
  Scenario: NormalizePackagePath normalizes input paths for storage
    Given a package path provided by a caller
    When NormalizePackagePath is called with the path
    Then the path is normalized according to package path rules
    And errors are returned for invalid paths using structured errors
    And the return value is suitable for storage and lookup
