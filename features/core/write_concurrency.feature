@domain:core @m2 @REQ-CORE-136 @spec(api_core.md#1237-write-concurrency) @spec(api_writing.md#1-safewrite---atomic-package-writing)
Feature: Write concurrency defines not safe for concurrent calls

  @REQ-CORE-136 @happy
  Scenario: Write is not safe for concurrent calls
    Given a package opened for writing
    When Write is called concurrently from multiple goroutines
    Then the behavior is undefined or errors are returned
    And callers must serialize write operations
    And the concurrency restriction is documented
