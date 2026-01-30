@domain:core @m2 @REQ-CORE-138 @spec(api_core.md#1241-safewrite-purpose) @spec(api_writing.md#1-safewrite---atomic-package-writing)
Feature: SafeWrite purpose defines atomic write with temp file strategy

  @REQ-CORE-138 @happy
  Scenario: SafeWrite uses atomic write with temp file strategy
    Given a package opened for writing
    When SafeWrite is called
    Then the write is atomic
    And a temp file strategy is used to ensure atomicity
    And the purpose matches the SafeWrite specification
