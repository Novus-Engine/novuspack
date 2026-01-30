@domain:core @m2 @REQ-CORE-144 @spec(api_core.md#1247-safewrite-concurrency) @spec(api_writing.md#1-safewrite---atomic-package-writing)
Feature: SafeWrite concurrency defines not safe for concurrent calls

  @REQ-CORE-144 @happy
  Scenario: SafeWrite is not safe for concurrent calls
    Given a package opened for writing
    When SafeWrite is called concurrently from multiple goroutines
    Then the behavior is undefined or errors are returned
    And callers must serialize SafeWrite operations
    And the concurrency restriction is documented
