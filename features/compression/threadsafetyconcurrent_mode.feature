@domain:compression @m2 @REQ-COMPR-142 @spec(api_package_compression.md#813-threadsafetyconcurrent)
Feature: ThreadSafetyConcurrent Mode

  @REQ-COMPR-142 @happy
  Scenario: ThreadSafetyConcurrent supports concurrent read/write operations
    Given compression operations with ThreadSafetyConcurrent mode
    When concurrent operations are performed
    Then concurrent read/write operations are supported
    And multiple goroutines can access package concurrently
    And thread safety is provided

  @REQ-COMPR-142 @happy
  Scenario: ThreadSafetyConcurrent uses read-write mutex for optimal read performance
    Given compression operations with ThreadSafetyConcurrent mode
    When read operations are performed concurrently
    Then read-write mutex is used for optimal read performance
    And multiple readers can access simultaneously
    And write operations are protected

  @REQ-COMPR-142 @happy
  Scenario: ThreadSafetyConcurrent provides good performance for concurrent access
    Given compression operations with concurrent access
    When ThreadSafetyConcurrent mode is used
    Then good performance is achieved for concurrent access
    And read operations are optimized
    And write operations are properly synchronized
