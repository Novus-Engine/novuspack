@domain:core @m2 @REQ-CORE-152 @spec(api_core.md#1257-fastwrite-concurrency) @spec(api_writing.md#2-fastwrite---in-place-package-updates)
Feature: FastWrite concurrency defines not safe for concurrent calls

  @REQ-CORE-152 @happy
  Scenario: FastWrite is not safe for concurrent calls
    Given a package opened for writing
    When FastWrite is called concurrently from multiple goroutines
    Then the behavior is undefined or errors are returned
    And callers must serialize FastWrite operations
    And the concurrency restriction is documented
